package chat

import (
	"errors"
	"io"
	"os"

	"go.uber.org/zap"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
	"github.com/greenboxal/aip/aip-langchain/pkg/chain"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm"
)

const MemoryContextKey chain.ContextKey[Memory] = "MemoryKey"

type predictChain struct {
	logger *zap.SugaredLogger

	prompt Prompt
	model  LanguageModel

	memory  *chain.ContextKey[Memory]
	outputs []chain.OutputParser

	opts  ChatOptions
	debug bool
}

func (p *predictChain) Run(ctx chain.ChainContext) error {
	var result Message
	var memory Memory

	if p.memory != nil {
		memory = chain.Input[Memory](ctx, *p.memory)
	}

	if memory != nil {
		_, err := memory.Load(ctx)

		if err != nil {
			return err
		}
	}

	prompt, err := p.prompt.Build(ctx)

	if err != nil {
		return err
	}

	if p.debug {
		_, _ = os.Stdout.Write([]byte("Request:\n\n"))
		_, _ = os.Stdout.Write([]byte(prompt.AsText()))
	}

	shouldStream := p.debug

	opts := p.buildOptions()

	if shouldStream {
		var entry MessageEntry

		stream, err := p.model.PredictChatStream(ctx.Context(), prompt, opts...)

		if err != nil {
			return err
		}

		if p.debug {
			_, _ = os.Stdout.Write([]byte("\n"))
		}

		for {
			frag, err := stream.Recv()

			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				return err
			}

			entry.Text += frag.Delta

			if p.debug {
				_, _ = os.Stdout.Write([]byte(frag.Delta))
			}
		}

		if p.debug {
			_, _ = os.Stdout.Write([]byte("\n"))
		}

		entry.Role = msn.RoleAI

		result.Entries = append(result.Entries, entry)
	} else {
		result, err = p.model.PredictChat(ctx.Context(), prompt, opts...)

		if err != nil {
			return err
		}
	}

	if memory != nil {
		if err := memory.Append(ctx, prompt); err != nil {
			return err
		}

		if err := memory.Append(ctx, result); err != nil {
			return err
		}
	}

	for _, output := range p.outputs {
		err := output.Parse(ctx, result.Entries[0].Text)

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *predictChain) buildOptions() []llm.PredictOption {
	opts := []llm.PredictOption{
		llm.WithMaxTokens(p.opts.MaxTokens),
	}

	if p.opts.AutoMaxTokens {
		opts = append(opts, llm.WithAutoMaxTokens(p.opts.MaxTokens))
	}

	if p.opts.Functions != nil {
		opts = append(opts, llm.WithFunctions(p.opts.Functions))
	}

	if p.opts.AllowFunctionCall {
		opts = append(opts, llm.WithAllowFunctionCall())
	}

	if p.opts.ForceFunctionCall != nil {
		opts = append(opts, llm.WithForceFunctionCall(*p.opts.ForceFunctionCall))
	}

	return opts
}

func CompletionMessageParser(key chain.ContextKey[Message]) chain.OutputParser {
	return chain.OutputParserFunc(func(ctx chain.ChainContext, result string) error {
		msg := Compose(Entry(msn.RoleAI, result))

		ctx.SetOutput(key, msg)

		return nil
	})
}

type ChatOptions struct {
	Temperature   float32
	StopSequences []string

	OutputParsers []chain.OutputParser
	ChatMemory    *chain.ContextKey[Memory]

	MaxTokens     int
	AutoMaxTokens bool

	Functions         map[string]llm.FunctionDeclaration
	AllowFunctionCall bool
	ForceFunctionCall *string
}

type ChatOption func(*ChatOptions)

func WithStopSequences(stopSequences ...string) ChatOption {
	return func(opts *ChatOptions) {
		opts.StopSequences = stopSequences
	}
}

func WithTemperature(temperature float32) ChatOption {
	return func(opts *ChatOptions) {
		opts.Temperature = temperature
	}
}

func WithTemperatureScale(scale float32) ChatOption {
	return func(opts *ChatOptions) {
		opts.Temperature *= scale

		if opts.Temperature > 1 {
			opts.Temperature = 1
		}
	}
}

func WithFunctions(functions []llm.FunctionDeclaration) ChatOption {
	return func(opts *ChatOptions) {
		if opts.Functions == nil {
			opts.Functions = make(map[string]llm.FunctionDeclaration)
		}

		for _, v := range functions {
			opts.Functions[v.Name] = v
		}
	}
}

func WithAllowFunctionCall() ChatOption {
	return func(opts *ChatOptions) {
		opts.AllowFunctionCall = true
	}
}

func WithForceFunctionCall(functionName string) ChatOption {
	return func(opts *ChatOptions) {
		opts.AllowFunctionCall = true
		opts.ForceFunctionCall = &functionName
	}
}

func WithFunction(fn llm.FunctionDeclaration) ChatOption {
	return func(opts *ChatOptions) {
		if opts.Functions == nil {
			opts.Functions = make(map[string]llm.FunctionDeclaration)
		}

		opts.Functions[fn.Name] = fn
	}
}

func NewChatOptions(options ...ChatOption) (result ChatOptions) {
	for _, opt := range options {
		opt(&result)
	}

	return
}

func WithMaxTokens(maxTokens int) ChatOption {
	return func(options *ChatOptions) {
		options.MaxTokens = maxTokens
	}
}

func WithAuthMaxTokens(targetMaxTokens int) ChatOption {
	return func(options *ChatOptions) {
		options.MaxTokens = targetMaxTokens
		options.AutoMaxTokens = true
	}
}

func WithOutputParsers(parsers ...chain.OutputParser) ChatOption {
	return func(options *ChatOptions) {
		options.OutputParsers = append(options.OutputParsers, parsers...)
	}
}

func WithChatMemory(memory chain.ContextKey[Memory]) ChatOption {
	return func(options *ChatOptions) {
		options.ChatMemory = &memory
	}
}

func Predict(model LanguageModel, prompt Prompt, options ...ChatOption) chain.Handler {
	opts := NewChatOptions(options...)

	if len(opts.OutputParsers) == 0 {
		opts.OutputParsers = []chain.OutputParser{
			CompletionMessageParser(ChatReplyContextKey),
		}
	}

	if opts.MaxTokens == 0 {
		opts.MaxTokens = model.MaxTokens() / 2
	}

	handler := &predictChain{
		model:   model,
		prompt:  prompt,
		outputs: opts.OutputParsers,
		memory:  opts.ChatMemory,
		opts:    opts,

		debug: true,
	}

	return chain.New(
		chain.WithName("chat.Predict"),
		chain.WithHandler(handler),
	)
}
