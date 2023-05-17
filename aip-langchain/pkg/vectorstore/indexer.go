package vectorstore

import (
	"context"

	"github.com/greenboxal/aip/aip-langchain/pkg/chunkers"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm"
)

type Indexer interface {
	IndexDocument(ctx context.Context, document *Document, options ...IndexDocumentOption) (*IndexedDocument, error)
}

type IndexDocumentOptions struct {
	Chunker      chunkers.Chunker
	MaxChunkSize int
	ChunkOverlap int

	Embedder llm.Embedder
}

type IndexDocumentOption func(options *IndexDocumentOptions)

func NewIndexDocumentOptions(opts ...IndexDocumentOption) *IndexDocumentOptions {
	options := &IndexDocumentOptions{
		Chunker:      chunkers.TikToken{},
		MaxChunkSize: 1024,
	}

	for _, opt := range opts {
		opt(options)
	}

	return options
}

func WithChunker(chunker chunkers.Chunker) IndexDocumentOption {
	return func(options *IndexDocumentOptions) {
		options.Chunker = chunker
	}
}

func WithMaxChunkSize(maxChunkSize int) IndexDocumentOption {
	return func(options *IndexDocumentOptions) {
		options.MaxChunkSize = maxChunkSize
	}
}

func WithChunkOverlap(chunkOverlap int) IndexDocumentOption {
	return func(options *IndexDocumentOptions) {
		options.ChunkOverlap = chunkOverlap
	}
}

func WithIndexEmbedder(embedder llm.Embedder) IndexDocumentOption {
	return func(options *IndexDocumentOptions) {
		options.Embedder = embedder
	}
}
