package chat

import (
	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
	chain2 "github.com/greenboxal/aip/aip-langchain/pkg/llm/chain"
)

const MemoryContextKey chain2.ContextKey[Memory] = "MemoryKey"

type predictChain struct {
	model   LanguageModel
	prompt  Prompt
	outputs []chain2.OutputParser
	memory  *chain2.ContextKey[Memory]
}

func (p *predictChain) Run(ctx chain2.ChainContext) error {
	var memory Memory

	if p.memory != nil {
		memory = chain2.Input[Memory](ctx, *p.memory)
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
		err := output.Parse(ctx, result.Entries[0].Text)

		if err != nil {
			return err
		}
	}

	return nil
}

func CompletionMessageParser(key chain2.ContextKey[Message]) chain2.OutputParser {
	return chain2.OutputParserFunc(func(ctx chain2.ChainContext, result string) error {
		msg := Compose(Entry(msn.RoleAI, result))

		ctx.SetOutput(key, msg)

		return nil
	})
}

type ChatOptions struct {
	OutputParsers []chain2.OutputParser
	ChatMemory    *chain2.ContextKey[Memory]
}

type ChatOption func(*ChatOptions)

func NewChatOptions(options ...ChatOption) (result ChatOptions) {
	for _, opt := range options {
		opt(&result)
	}

	return
}

func WithOutputParsers(parsers ...chain2.OutputParser) ChatOption {
	return func(options *ChatOptions) {
		options.OutputParsers = append(options.OutputParsers, parsers...)
	}
}

func WithChatMemory(memory chain2.ContextKey[Memory]) ChatOption {
	return func(options *ChatOptions) {
		options.ChatMemory = &memory
	}
}

func Predict(model LanguageModel, prompt Prompt, options ...ChatOption) chain2.Handler {
	opts := NewChatOptions(options...)

	if len(opts.OutputParsers) == 0 {
		opts.OutputParsers = []chain2.OutputParser{
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
