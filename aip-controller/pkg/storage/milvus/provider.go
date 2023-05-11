package milvus

import (
	"context"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"go.uber.org/zap"

	"github.com/greenboxal/aip/aip-controller/pkg/indexing"
)

type IndexProvider struct {
	logger *zap.SugaredLogger
	milvus *Milvus
}

func NewIndexProvider(
	logger *zap.SugaredLogger,
	milvus *Milvus,
) *IndexProvider {
	return &IndexProvider{
		logger: logger.Named("milvus-index-provider"),
		milvus: milvus,
	}
}

func (i *IndexProvider) IndexDocument(
	ctx context.Context,
	document *indexing.Document,
	options ...indexing.IndexDocumentOption,
) (*indexing.IndexedDocument, error) {
	opts := indexing.NewIndexDocumentOptions(options...)

	textChunks, err := opts.Chunker.SplitTextIntoChunks(
		ctx,
		document.Content,
		opts.MaxChunkSize,
		opts.ChunkOverlap,
	)

	if err != nil {
		return nil, err
	}

	documentChunks := make([]*indexing.DocumentChunk, len(textChunks))

	result := &indexing.IndexedDocument{}
	result.DocumentReference = document.DocumentReference
	result.Chunks = make([]indexing.DocumentChunkReference, len(textChunks))

	embeddings, err := opts.Embedder.GetEmbeddings(ctx, textChunks)

	if err != nil {
		return nil, err
	}

	for idx, content := range textChunks {
		ref := indexing.DocumentChunkReference{
			ID:    document.ID,
			Type:  document.Type,
			Chunk: idx,
		}

		documentChunks[idx] = &indexing.DocumentChunk{
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
	options ...indexing.SearchDocumentOption,
) (*indexing.SearchResult, error) {
	opts := indexing.NewSearchDocumentOptions(options...)
	outputFields := []string{"document_id", "document_type", "chunk_id"}

	if opts.ReturnHitEmbeddings {
		outputFields = append(outputFields, "chunk_embeddings")
	}

	if opts.ReturnHitContents {
		outputFields = append(outputFields, "chunk_content")
	}

	embeddings, err := opts.Embedder.GetEmbeddings(ctx, []string{query})

	if err != nil {
		return nil, err
	}

	vectors := []entity.Vector{entity.FloatVector(embeddings[0].Embeddings)}

	sp, err := entity.NewIndexFlatSearchParam()

	if err != nil {
		return nil, err
	}

	hits, err := i.milvus.client.Search(
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

	result := &indexing.SearchResult{}
	result.Hits = make([]indexing.SearchHit, 0, len(hits))

	for _, h := range hits {
		for i := 0; i < h.ResultCount; i++ {
			hit := indexing.SearchHit{
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
