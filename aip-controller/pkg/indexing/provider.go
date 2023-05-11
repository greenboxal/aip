package indexing

import (
	"context"

	"github.com/greenboxal/aip/aip-controller/pkg/llm"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/chunkers"
)

type Provider interface {
	IndexDocument(ctx context.Context, document *Document, options ...IndexDocumentOption) (*IndexedDocument, error)

	Search(ctx context.Context, query string, options ...SearchDocumentOption) (*SearchResult, error)
}

type Document struct {
	DocumentReference

	Content string
}

type DocumentReference struct {
	ID   string
	Type string
}

type DocumentChunkReference struct {
	ID    string
	Type  string
	Chunk int
}

type DocumentChunk struct {
	DocumentChunkReference

	Content    string
	Embeddings llm.Embeddings
}

type SearchHit struct {
	DocumentChunkReference

	Score float32

	Embeddings []float32
	Content    string
}

type SearchResult struct {
	Hits []SearchHit
}

type IndexedDocument struct {
	DocumentReference

	Chunks []DocumentChunkReference
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
