package vectorstore

import "context"

type Provider interface {
	Dimensions() int

	IndexChunk(
		ctx context.Context,
		chunk DocumentChunk,
	) error

	SimilaritySearch(
		ctx context.Context,
		embeddings []float32,
		opts *SearchDocumentOptions,
	) (*SearchResult, error)

	Close() error
}
