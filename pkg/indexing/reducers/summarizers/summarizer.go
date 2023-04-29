package summarizers

import (
	"context"

	"github.com/greenboxal/aip/pkg/indexing"
)

type SummarizeOption func(options *SummarizeOptions)

func WithContextHint(contextHint string) SummarizeOption {
	return func(options *SummarizeOptions) {
		options.ContextHint = contextHint
	}
}

func WithMaxTokens(maxTokens int) SummarizeOption {
	return func(options *SummarizeOptions) {
		options.MaxTokens = maxTokens
	}
}

func WithTemperature(temperature float32) SummarizeOption {
	return func(options *SummarizeOptions) {
		options.Temperature = temperature
	}
}

type SummarizeOptions struct {
	ContextHint string

	MaxTokens   int
	Temperature float32
}

type Summarizer interface {
	MaxTokens() int

	Summarize(
		ctx context.Context,
		document string,
		options ...SummarizeOption,
	) (indexing.MemoryData, error)
}
