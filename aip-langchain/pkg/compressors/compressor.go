package compressors

import (
	"context"

	chain2 "github.com/greenboxal/aip/aip-langchain/pkg/chain"
	"github.com/greenboxal/aip/aip-langchain/pkg/documents"
)

const CompressionInputKey chain2.ContextKey[string] = "input"
const CompressionOutputKey chain2.ContextKey[string] = "output"

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
	MaxInputTokens() int

	Compress(ctx context.Context, text string, options ...CompressionOption) (string, error)
	Decompress(ctx context.Context, text string, options ...CompressionOption) (string, error)
}

type DocumentCompressor interface {
	MaxTokens() int

	CompressDocument(ctx context.Context, text documents.Document, options ...CompressionOption) (documents.Document, error)
	DecompressDocument(ctx context.Context, text documents.Document, options ...CompressionOption) (documents.Document, error)
}

func CompressorChain(compressor Compressor) chain2.Handler {
	return chain2.Func(func(ctx chain2.ChainContext) error {
		input := chain2.Input(ctx, CompressionInputKey)

		output, err := compressor.Compress(ctx.Context(), input)

		if err != nil {
			return err
		}

		ctx.SetOutput(CompressionOutputKey, output)

		return nil
	})
}
