package compressors

import (
	"context"

	"github.com/greenboxal/aip/pkg/llm"
)

var simpleCompressorCompressPrompt = llm.NewTemplatePrompt(`
Summarize the text below in 1 sentence:

{{.input}}
`)

type SimpleCompressor struct {
	compressChain   llm.Chainable
	decompressChain llm.Chainable
}

func (s *SimpleCompressor) MaxTokens() int {
	return 0
}

func NewSimpleCompressor(model llm.LanguageModel) Compressor {
	compressChain := llm.Chain(
		llm.Predict(model, simpleCompressorCompressPrompt, llm.Output("uncompressed", llm.NoopParser())),
	)

	return &SimpleCompressor{
		compressChain: compressChain,
	}
}

func (s *SimpleCompressor) Compress(ctx context.Context, text string, options ...CompressionOption) (string, error) {
	pctx := llm.NewChainContext(ctx)

	pctx.SetInput("input", text)

	if err := s.compressChain.Run(pctx); err != nil {
		return "", err
	}

	return pctx.Output("compressed").(string), nil
}

func (s *SimpleCompressor) Decompress(ctx context.Context, text string, options ...CompressionOption) (string, error) {
	//TODO implement me
	panic("implement me")
}
