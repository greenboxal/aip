package faiss

import (
	"context"
	"sync"

	"github.com/DataIntelligenceCrew/go-faiss"

	"github.com/greenboxal/aip/aip-langchain/pkg/vectorstore"
)

func NewIndex(dim int) (*Index, error) {
	var err error

	idx := &Index{}
	idx.dim = dim
	idx.chunks = make(map[int64]vectorstore.DocumentChunk)

	idx.index, err = faiss.NewIndexFlatIP(dim)

	if err != nil {
		return nil, err
	}

	return idx, nil
}

type Index struct {
	m sync.RWMutex

	dim   int
	index faiss.Index

	chunks map[int64]vectorstore.DocumentChunk
}

var _ vectorstore.Index = (*Index)(nil)

func (idx *Index) Dimensions() int {
	return idx.dim
}

func (idx *Index) IndexDocument(
	ctx context.Context,
	document *vectorstore.Document,
	options ...vectorstore.IndexDocumentOption,
) (*vectorstore.IndexedDocument, error) {
	opts := vectorstore.NewIndexDocumentOptions(options...)

	textChunks, err := opts.Chunker.SplitTextIntoStrings(
		ctx,
		document.Content,
		opts.MaxChunkSize,
		opts.ChunkOverlap,
	)

	if err != nil {
		return nil, err
	}

	documentChunks := make([]*vectorstore.DocumentChunk, len(textChunks))

	result := &vectorstore.IndexedDocument{}
	result.DocumentReference = document.DocumentReference
	result.Chunks = make([]vectorstore.DocumentChunkReference, len(textChunks))

	embeddings, err := opts.Embedder.GetEmbeddings(ctx, textChunks)

	if err != nil {
		return nil, err
	}

	for idx, content := range textChunks {
		ref := vectorstore.DocumentChunkReference{
			ID:    document.ID,
			Type:  document.Type,
			Chunk: idx,
		}

		documentChunks[idx] = &vectorstore.DocumentChunk{
			DocumentChunkReference: ref,
			Content:                content,
			Embeddings:             embeddings[idx],
		}

		result.Chunks[idx] = ref
	}

	for _, chunk := range documentChunks {
		err = idx.IndexChunk(ctx, *chunk)

		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (idx *Index) Search(
	ctx context.Context,
	query string,
	options ...vectorstore.SearchDocumentOption,
) (*vectorstore.SearchResult, error) {
	opts := vectorstore.NewSearchDocumentOptions(options...)
	embeddings, err := opts.Embedder.GetEmbeddings(ctx, []string{query})

	if err != nil {
		return nil, err
	}

	return idx.SimilaritySearch(ctx, embeddings[0].Embeddings, opts)
}

func (idx *Index) IndexChunk(ctx context.Context, chunk vectorstore.DocumentChunk) error {
	idx.m.Lock()
	defer idx.m.Unlock()

	index := idx.index.Ntotal()

	err := idx.index.Add(chunk.Embeddings.Embeddings)

	if err != nil {
		return err
	}

	idx.chunks[index] = chunk

	return nil
}

func (idx *Index) SimilaritySearch(
	ctx context.Context,
	embeddings []float32,
	opts *vectorstore.SearchDocumentOptions,
) (*vectorstore.SearchResult, error) {
	idx.m.RLock()
	defer idx.m.RUnlock()

	distances, labels, err := idx.index.Search(embeddings, int64(opts.MaxHits))

	if err != nil {
		return nil, err
	}

	result := &vectorstore.SearchResult{}

	result.Hits = make([]vectorstore.SearchHit, len(labels))

	for i, label := range labels {
		distance := distances[i]

		chunk := idx.chunks[label]

		hit := vectorstore.SearchHit{
			DocumentChunkReference: chunk.DocumentChunkReference,
			Score:                  distance,
		}

		if opts.ReturnHitContents {
			hit.Content = chunk.Content
		}

		if opts.ReturnHitEmbeddings {
			hit.Embeddings = chunk.Embeddings.Embeddings
		}

		result.Hits[i] = hit
	}

	return result, nil
}

func (idx *Index) Close() error {
	if idx.index != nil {
		idx.index.Delete()
		idx.index = nil
	}

	return nil
}
