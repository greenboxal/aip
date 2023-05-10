package chat

import (
	"github.com/greenboxal/aip/aip-controller/pkg/llm/chain"
)

const MemoryContextKey chain.ContextKey[Memory] = "MemoryKey"

type predictChain struct {
	model   LanguageModel
	prompt  Prompt
	outputs []chain.OutputParser
	memory  *chain.ContextKey[Memory]
}

func (p *predictChain) Run(ctx chain.ChainContext) error {
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

	result, err := p.model.PredictChat(ctx.Context(), prompt)

	if err != nil {
		return err
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
		err := output.Parse(ctx, result.Entries[0].Content)

		if err != nil {
			return err
		}
	}

	return nil
}

func CompletionMessageParser(key chain.ContextKey[Message]) chain.OutputParser {
	return chain.OutputParserFunc(func(ctx chain.ChainContext, result string) error {
		msg := Compose(Entry(RoleAI, result))

		ctx.SetOutput(key, msg)

		return nil
	})
}

type ChatOptions struct {
	OutputParsers []chain.OutputParser
	ChatMemory    *chain.ContextKey[Memory]
}

type ChatOption func(*ChatOptions)

func NewChatOptions(options ...ChatOption) (result ChatOptions) {
	for _, opt := range options {
		opt(&result)
	}

	return
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

func Predict(model LanguageModel, prompt Prompt, options ...ChatOption) chain.Chain {
	opts := NewChatOptions(options...)

	if len(opts.OutputParsers) == 0 {
		opts.OutputParsers = []chain.OutputParser{
			CompletionMessageParser(ChatReplyContextKey),
		}
	}

	return &predictChain{
		model:   model,
		prompt:  prompt,
		outputs: opts.OutputParsers,
		memory:  opts.ChatMemory,
	}
}
