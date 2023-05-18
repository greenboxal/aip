package milvus

import (
	"context"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"go.uber.org/zap"

	"github.com/greenboxal/aip/aip-langchain/pkg/vectorstore"
)

type IndexProvider struct {
	logger *zap.SugaredLogger
	milvus *Milvus
	dim    int
}

func NewIndexProvider(
	logger *zap.SugaredLogger,
	milvus *Milvus,
) *IndexProvider {
	return &IndexProvider{
		logger: logger.Named("milvus-index-provider"),
		milvus: milvus,
		dim:    1536,
	}
}

var _ vectorstore.Index = (*IndexProvider)(nil)

func (i *IndexProvider) Dimensions() int {
	return i.dim
}

func (i *IndexProvider) IndexDocument(
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

	err = i.milvus.IndexDocumentChunks(ctx, documentChunks)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (i *IndexProvider) Search(
	ctx context.Context,
	query string,
	options ...vectorstore.SearchDocumentOption,
) (*vectorstore.SearchResult, error) {
	opts := vectorstore.NewSearchDocumentOptions(options...)

	embeddings, err := opts.Embedder.GetEmbeddings(ctx, []string{query})

	if err != nil {
		return nil, err
	}

	return i.SimilaritySearch(ctx, embeddings[0].Embeddings, opts)
}

func (i *IndexProvider) SimilaritySearch(ctx context.Context, embeddings []float32, opts *vectorstore.SearchDocumentOptions) (*vectorstore.SearchResult, error) {
	vectors := []entity.Vector{entity.FloatVector(embeddings)}
	outputFields := []string{"document_id", "document_type", "chunk_id"}

	if opts.ReturnHitEmbeddings {
		outputFields = append(outputFields, "chunk_embeddings")
	}

	if opts.ReturnHitContents {
		outputFields = append(outputFields, "chunk_content")
	}

	sp, err := entity.NewIndexFlatSearchParam()

	if err != nil {
		return nil, err
	}

	hits, err := i.milvus.Client().Search(
		ctx,
		"global_index",
		[]string{"_default"},
		"_id > 0",
		outputFields,
		vectors,
		"chunk_embeddings",
		entity.IP,
		opts.MaxHits,
		sp,
	)

	if err != nil {
		return nil, err
	}

	result := &vectorstore.SearchResult{}
	result.Hits = make([]vectorstore.SearchHit, 0, len(hits))

	for _, h := range hits {
		for i := 0; i < h.ResultCount; i++ {
			hit := vectorstore.SearchHit{
				Score: h.Scores[0],
			}

			for _, field := range h.Fields {
				switch field.Name() {
				case "document_id":
					hit.ID = field.(*entity.ColumnVarChar).Data()[i]

				case "document_type":
					hit.Type = field.(*entity.ColumnVarChar).Data()[i]

				case "chunk_id":
					hit.Chunk = int(field.(*entity.ColumnInt32).Data()[i])

				case "chunk_content":
					hit.Content = field.(*entity.ColumnVarChar).Data()[i]

				case "chunk_embeddings":
					hit.Embeddings = field.(*entity.ColumnFloatVector).Data()[i]
				}
			}

			result.Hits = append(result.Hits, hit)
		}
	}

	return result, nil
}

func (i *IndexProvider) Close() error {
	return nil
}
