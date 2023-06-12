package chat

import (
	"errors"
	"io"
	"os"

	"go.uber.org/zap"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
	"github.com/greenboxal/aip/aip-langchain/pkg/chain"
)

const MemoryContextKey chain.ContextKey[Memory] = "MemoryKey"

type predictChain struct {
	prompt Prompt
	model  LanguageModel

	memory  *chain.ContextKey[Memory]
	outputs []chain.OutputParser

	debug  bool
	logger *zap.SugaredLogger
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

	if shouldStream {
		var entry MessageEntry

		stream, err := p.model.PredictChatStream(ctx.Context(), prompt)

		if err != nil {
			return err
		}

		if p.debug {
			_, _ = os.Stdout.Write([]byte("Reply:\n\n"))
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
		result, err = p.model.PredictChat(ctx.Context(), prompt)

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

func CompletionMessageParser(key chain.ContextKey[Message]) chain.OutputParser {
	return chain.OutputParserFunc(func(ctx chain.ChainContext, result string) error {
		msg := Compose(Entry(msn.RoleAI, result))

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

func Predict(model LanguageModel, prompt Prompt, options ...ChatOption) chain.Handler {
	opts := NewChatOptions(options...)

	if len(opts.OutputParsers) == 0 {
		opts.OutputParsers = []chain.OutputParser{
			CompletionMessageParser(ChatReplyContextKey),
		}
	}

	handler := &predictChain{
		model:   model,
		prompt:  prompt,
		outputs: opts.OutputParsers,
		memory:  opts.ChatMemory,

		debug: true,
	}

	return chain.New(
		chain.WithName("chat.Predict"),
		chain.WithHandler(handler),
	)
}
