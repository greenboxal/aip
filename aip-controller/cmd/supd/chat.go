package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
	chain2 "github.com/greenboxal/aip/aip-langchain/pkg/llm/chain"
	chat2 "github.com/greenboxal/aip/aip-langchain/pkg/llm/chat"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/compressors"
	openai2 "github.com/greenboxal/aip/aip-langchain/pkg/llm/providers/openai"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/tokenizers"
)

type ChatHandler struct {
	Input  io.Reader
	Output io.Writer

	model      *openai2.ChatLanguageModel
	tokenizer  *tokenizers.TikTokenTokenizer
	compressor compressors.Compressor

	ctx   chain2.ChainContext
	chain chain2.Handler
}

var ChatPrompt = chat2.ComposeTemplate(
	chat2.HistoryFromContext(chat2.ChatHistoryContextKey),
	chat2.EntryTemplate(msn.RoleUser, chain2.TemplateFromContext(chat2.ChatReplyContextKey)),
	chat2.EntryTemplate(msn.RoleAI, chain2.Static("")),
)

func NewChatHandler(client *openai2.Client) (*ChatHandler, error) {
	model := &openai2.ChatLanguageModel{
		Client: client,
		Model:  "gpt-3.5-turbo",
	}

	tokenizer, err := tokenizers.TikTokenForModel(openai2.AdaEmbeddingV2.String())

	if err != nil {
		return nil, err
	}

	compressor := compressors.NewSimpleCompressor(model, tokenizer)

	return &ChatHandler{
		Input:  os.Stdin,
		Output: os.Stdout,

		model:      model,
		tokenizer:  tokenizer,
		compressor: compressor,

		chain: chain2.Sequential(
			chat2.Predict(model, ChatPrompt),

			chain2.MapContext(
				chain2.TransformInput(chat2.ChatReplyContextKey, compressors.CompressionInputKey, func(msg chat2.Message) string {
					return msg.String()
				}),
			),

			compressors.CompressorChain(compressor),
		),
	}, nil
}

func (ch *ChatHandler) Run(proc goprocess.Process) {
	history := chat2.Message{}
	inputStream := bufio.NewReader(ch.Input)

	ctx := goprocessctx.OnClosingContext(proc)
	ch.ctx = chain2.NewChainContext(ctx)

	for {
		select {
		case <-ctx.Done():
			return

		default:
		}

		entry := chat2.MessageEntry{
			Role: msn.RoleUser,
			Text: "",
		}

		_, err := fmt.Fprintf(ch.Output, "%s", entry)

		if err != nil {
			panic(err)
		}

		entry.Text, err = inputStream.ReadString('\n')

		if err != nil {
			panic(err)
		}

		ch.ctx.Flip()
		ch.ctx.SetInput(chat2.ChatReplyContextKey, entry)
		ch.ctx.SetInput(chat2.ChatHistoryContextKey, history)

		if err := ch.chain.Run(ch.ctx); err != nil {
			_, err = fmt.Fprintf(ch.Output, "ERROR: %s\n", err)

			if err != nil {
				panic(err)
			}
		}

		result := chain2.Output(ch.ctx, chat2.ChatReplyContextKey)

		history.Entries = append(history.Entries, entry)
		history.Entries = append(history.Entries, result.Entries...)

		for _, replies := range result.Entries {
			_, err = fmt.Fprintf(ch.Output, "%s\n", replies)

			if err != nil {
				panic(err)
			}
		}

		compressed := chain2.Output(ch.ctx, compressors.CompressionOutputKey)

		_, err = fmt.Fprintf(ch.Output, "COMPRESSED: %s\n", compressed)

		if err != nil {
			panic(err)
		}
	}
}
