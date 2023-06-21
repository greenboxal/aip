package llm

import "context"

type LanguageModel interface {
	MaxTokens() int

	Predict(ctx context.Context, prompt string, options ...PredictOption) (string, error)
}

func WithFunctions(functions map[string]FunctionDeclaration) PredictOption {
	return func(opts *PredictOptions) {
		if opts.Functions == nil {
			opts.Functions = make(map[string]FunctionDeclaration)
		}

		for _, v := range functions {
			opts.Functions[v.Name] = v
		}
	}
}

type PredictOption func(opts *PredictOptions)

type PredictOptions struct {
	Temperature   float32
	StopSequences []string

	MaxTokens     int
	AutoMaxTokens bool

	Functions         map[string]FunctionDeclaration
	AllowFunctionCall bool
	ForceFunctionCall *string
}

func NewPredictOptions(options ...PredictOption) PredictOptions {
	opts := PredictOptions{}

	for _, opt := range options {
		opt(&opts)
	}

	return opts
}

func WithAllowFunctionCall() PredictOption {
	return func(opts *PredictOptions) {
		opts.AllowFunctionCall = true
	}
}

func WithForceFunctionCall(functionName string) PredictOption {
	return func(opts *PredictOptions) {
		opts.AllowFunctionCall = true
		opts.ForceFunctionCall = &functionName
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

func WithAutoMaxTokens(targetMaxTokens int) PredictOption {
	return func(opts *PredictOptions) {
		opts.MaxTokens = targetMaxTokens
		opts.AutoMaxTokens = true
	}
}
