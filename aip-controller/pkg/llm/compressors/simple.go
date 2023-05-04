package compressors

import (
	"context"

	"github.com/greenboxal/aip/aip-controller/pkg/llm"
	chain2 "github.com/greenboxal/aip/aip-controller/pkg/llm/chain"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/tokenizers"
)

var simpleCompressorCompressPrompt = chain2.NewTemplatePrompt(`
Summarize the text below:

{{.input}}
`,
	chain2.WithRequiredInput(CompressionInputKey),
	chain2.WithRequiredOutput(CompressionOutputKey),
)

type SimpleCompressor struct {
	model     llm.LanguageModel
	tokenizer tokenizers.BasicTokenizer

	compressChain   chain2.Chain
	decompressChain chain2.Chain
}

func (s *SimpleCompressor) MaxInputTokens() int {
	return s.model.MaxTokens() - simpleCompressorCompressPrompt.GetEmptyTokenCount(s.tokenizer)
}

func NewSimpleCompressor(model llm.LanguageModel, tokenizer tokenizers.BasicTokenizer) Compressor {
	compressChain := chain2.Compose(
		chain2.Predict(model, simpleCompressorCompressPrompt),
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
	target chain2.Chain,
	text string,
	options ...CompressionOption,
) (string, error) {
	_ = NewCompressionOptions(options...)

	pctx := chain2.NewChainContext(ctx)

	pctx.SetInput(CompressionInputKey, text)

	if err := target.Run(pctx); err != nil {
		return "", err
	}

	return chain2.Output(pctx, CompressionOutputKey), nil
}
