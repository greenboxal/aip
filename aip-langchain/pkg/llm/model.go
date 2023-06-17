package llm

import "context"

type FunctionDeclaration struct {
	Name        string
	Description string
	Parameters  []any
}

type PredictOptions struct {
	Temperature   float32
	StopSequences []string

	MaxTokens     int
	AutoMaxTokens bool

	Functions map[string]FunctionDeclaration
}

type PredictOption func(opts *PredictOptions)

func WithFunctions(functions []FunctionDeclaration) PredictOption {
	return func(opts *PredictOptions) {
		if opts.Functions == nil {
			opts.Functions = make(map[string]FunctionDeclaration)
		}

		for _, v := range functions {
			opts.Functions[v.Name] = v
		}
	}
}

func WithFunction(fn FunctionDeclaration) PredictOption {
	return func(opts *PredictOptions) {
		if opts.Functions == nil {
			opts.Functions = make(map[string]FunctionDeclaration)
		}

		opts.Functions[fn.Name] = fn
	}
}

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

func WithAutoMaxTokens() PredictOption {
	return func(opts *PredictOptions) {
		opts.AutoMaxTokens = true
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
