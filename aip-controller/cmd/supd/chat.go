package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"

	chain "github.com/greenboxal/aip/aip-controller/pkg/llm/chain"
	chat "github.com/greenboxal/aip/aip-controller/pkg/llm/chat"
	compressors2 "github.com/greenboxal/aip/aip-controller/pkg/llm/compressors"
	openai "github.com/greenboxal/aip/aip-controller/pkg/llm/providers/openai"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/tokenizers"
)

type ChatHandler struct {
	Input  io.Reader
	Output io.Writer

	model      *openai.ChatLanguageModel
	tokenizer  *tokenizers.TikTokenTokenizer
	compressor compressors2.Compressor

	ctx   chain.ChainContext
	chain chain.Chain
}

var ChatPrompt = chat.ComposeTemplate(
	chat.HistoryFromContext(chat.ChatHistoryContextKey),
	chat.EntryTemplate(chat.RoleUser, chain.TemplateFromContext(chat.ChatReplyContextKey)),
	chat.EntryTemplate(chat.RoleAI, chain.Static("")),
)

func NewChatHandler(client *openai.Client) (*ChatHandler, error) {
	model := &openai.ChatLanguageModel{
		Client: client,
		Model:  "gpt-3.5-turbo",
	}

	tokenizer, err := tokenizers.TikTokenForModel(openai.AdaEmbeddingV2.String())

	if err != nil {
		return nil, err
	}

	compressor := compressors2.NewSimpleCompressor(model, tokenizer)

	return &ChatHandler{
		Input:  os.Stdin,
		Output: os.Stdout,

		model:      model,
		tokenizer:  tokenizer,
		compressor: compressor,

		chain: chain.Compose(
			chat.Predict(model, ChatPrompt),

			chain.MapContext(
				chain.TransformInput(chat.ChatReplyContextKey, compressors2.CompressionInputKey, func(msg chat.Message) string {
					return msg.String()
				}),
			),

			compressors2.CompressorChain(compressor),
		),
	}, nil
}

func (ch *ChatHandler) Run(proc goprocess.Process) {
	history := chat.Message{}
	inputStream := bufio.NewReader(ch.Input)

	ctx := goprocessctx.OnClosingContext(proc)
	ch.ctx = chain.NewChainContext(ctx)

	for {
		select {
		case <-ctx.Done():
			return

		default:
		}

		entry := chat.MessageEntry{
			Role:    chat.RoleUser,
			Content: "",
		}

		_, err := fmt.Fprintf(ch.Output, "%s", entry)

		if err != nil {
			panic(err)
		}

		entry.Content, err = inputStream.ReadString('\n')

		if err != nil {
			panic(err)
		}

		ch.ctx.Flip()
		ch.ctx.SetInput(chat.ChatReplyContextKey, entry)
		ch.ctx.SetInput(chat.ChatHistoryContextKey, history)

		if err := ch.chain.Run(ch.ctx); err != nil {
			_, err = fmt.Fprintf(ch.Output, "ERROR: %s\n", err)

			if err != nil {
				panic(err)
			}
		}

		result := chain.Output(ch.ctx, chat.ChatReplyContextKey)

		history.Entries = append(history.Entries, entry)
		history.Entries = append(history.Entries, result.Entries...)

		for _, replies := range result.Entries {
			_, err = fmt.Fprintf(ch.Output, "%s\n", replies)

			if err != nil {
				panic(err)
			}
		}

		compressed := chain.Output(ch.ctx, compressors2.CompressionOutputKey)

		_, err = fmt.Fprintf(ch.Output, "COMPRESSED: %s\n", compressed)

		if err != nil {
			panic(err)
		}
	}
}
