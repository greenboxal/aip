package llm

import "context"

type PredictOptions struct {
	MaxTokens     int
	Temperature   float32
	StopSequences []string
}

type PredictOption func(opts *PredictOptions)

func WithStopSequences(stopSequences ...string) PredictOption {
	return func(opts *PredictOptions) {
		opts.StopSequences = stopSequences
	}
}

func WithTemperature(temperature float32) PredictOption {
	return func(opts *PredictOptions) {
		opts.Temperature = temperature
	}
}

func WithTemperatureScale(scale float32) PredictOption {
	return func(opts *PredictOptions) {
		opts.Temperature *= scale

		if opts.Temperature > 1 {
			opts.Temperature = 1
		}
	}
}

func WithMaxTokens(maxTokens int) PredictOption {
	return func(opts *PredictOptions) {
		opts.MaxTokens = maxTokens
	}
}

func NewPredictOptions(options ...PredictOption) PredictOptions {
	opts := PredictOptions{}

	for _, opt := range options {
		opt(&opts)
	}

	return opts
}

type LanguageModel interface {
	MaxTokens() int

	Predict(ctx context.Context, prompt string, options ...PredictOption) (string, error)
}
