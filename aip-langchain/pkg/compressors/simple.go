package compressors

import (
	"context"

	chain "github.com/greenboxal/aip/aip-langchain/pkg/chain"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm"
	"github.com/greenboxal/aip/aip-langchain/pkg/tokenizers"
)

var simpleCompressorCompressPrompt = chain.NewTemplatePrompt(`
Summarize the text below:

{{.input}}
`,
	chain.WithRequiredInput(CompressionInputKey),
	chain.WithRequiredOutput(CompressionOutputKey),
)

type SimpleCompressor struct {
	model     llm.LanguageModel
	tokenizer tokenizers.BasicTokenizer

	compressChain   chain.Handler
	decompressChain chain.Handler
}

func (s *SimpleCompressor) MaxInputTokens() int {
	return s.model.MaxTokens() - simpleCompressorCompressPrompt.GetEmptyTokenCount(s.tokenizer)
}

func NewSimpleCompressor(model llm.LanguageModel, tokenizer tokenizers.BasicTokenizer) Compressor {
	compressChain := chain.Sequential(
		chain.Predict(model, simpleCompressorCompressPrompt),
	)

	return &SimpleCompressor{
		model:         model,
		tokenizer:     tokenizer,
		compressChain: compressChain,
	}
}

func (s *SimpleCompressor) Compress(ctx context.Context, text string, options ...CompressionOption) (string, error) {
	return s.runSinglePass(ctx, s.compressChain, text, options...)
}

func (s *SimpleCompressor) Decompress(ctx context.Context, text string, options ...CompressionOption) (string, error) {
	return s.runSinglePass(ctx, s.decompressChain, text, options...)
}

func (s *SimpleCompressor) runSinglePass(
	ctx context.Context,
	target chain.Handler,
	text string,
	options ...CompressionOption,
) (string, error) {
	_ = NewCompressionOptions(options...)

	pctx := chain.NewChainContext(ctx)

	pctx.SetInput(CompressionInputKey, text)

	if err := target.Run(pctx); err != nil {
		return "", err
	}

	return chain.Output(pctx, CompressionOutputKey), nil
}
