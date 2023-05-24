package vectorstore

import (
	"context"

	"github.com/greenboxal/aip/aip-langchain/pkg/llm"
)

type Retriever interface {
	Dimensions() int

	Search(
		ctx context.Context,
		query string,
		options ...SearchDocumentOption,
	) (*SearchResult, error)

	SimilaritySearch(
		ctx context.Context,
		embeddings []float32,
		opts *SearchDocumentOptions,
	) (*SearchResult, error)
}

type SearchHit struct {
	DocumentChunkReference `json:"reference"`

	Score float32 `json:"score"`

	Embeddings []float32 `json:"embeddings"`
	Content    string    `json:"content"`
}

type SearchResult struct {
	Hits []SearchHit `json:"hits"`
}

type SearchDocumentOptions struct {
	MaxHits int

	Embedder llm.Embedder

	ReturnHitContents   bool
	ReturnHitEmbeddings bool
}

type SearchDocumentOption func(options *SearchDocumentOptions)

func NewSearchDocumentOptions(opts ...SearchDocumentOption) *SearchDocumentOptions {
	options := &SearchDocumentOptions{
		MaxHits: 10,
	}

	for _, opt := range opts {
		opt(options)
	}

	return options
}

func WithSearchEmbedder(embedder llm.Embedder) SearchDocumentOption {
	return func(options *SearchDocumentOptions) {
		options.Embedder = embedder
	}
}

func WithReturnHitContents() SearchDocumentOption {
	return func(options *SearchDocumentOptions) {
		options.ReturnHitContents = true
	}
}

func WithReturnHitEmbeddings() SearchDocumentOption {
	return func(options *SearchDocumentOptions) {
		options.ReturnHitEmbeddings = true
	}
}

func WithMaxHits(maxHits int) SearchDocumentOption {
	return func(options *SearchDocumentOptions) {
		options.MaxHits = maxHits
	}
}
