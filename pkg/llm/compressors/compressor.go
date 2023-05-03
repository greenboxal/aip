package compressors

import (
	"context"

	"github.com/greenboxal/aip/pkg/llm/documents"
)

type CompressionOption func(options *CompressionOptions)

func WithMaxTokens(maxTokens int) CompressionOption {
	return func(options *CompressionOptions) {
		options.MaxTokens = maxTokens
	}
}

type CompressionOptions struct {
	MaxTokens int
}

func NewCompressionOptions(options ...CompressionOption) (result CompressionOptions) {
	for _, opt := range options {
		opt(&result)
	}
	
	return
}

type Compressor interface {
	MaxTokens() int

	Compress(ctx context.Context, text string, options ...CompressionOption) (string, error)
	Decompress(ctx context.Context, text string, options ...CompressionOption) (string, error)
}

type DocumentCompressor interface {
	MaxTokens() int

	CompressDocument(ctx context.Context, text documents.Document, options ...CompressionOption) (documents.Document, error)
	DecompressDocument(ctx context.Context, text documents.Document, options ...CompressionOption) (documents.Document, error)
}
