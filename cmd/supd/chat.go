package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"

	"github.com/greenboxal/aip/pkg/llm"
	"github.com/greenboxal/aip/pkg/llm/chat"
	"github.com/greenboxal/aip/pkg/llm/providers/openai"
)

type ChatHandler struct {
	Input  io.Reader
	Output io.Writer

	ctx   llm.ChainContext
	chain llm.Chainable
}

var ChatPrompt = chat.ComposeTemplate(
	chat.HistoryFromContext(chat.ChatHistoryContextKey),
	chat.EntryTemplate("", chat.RoleUser, llm.TemplateFromContext(chat.ChatReplyContextKey)),
)

func NewChatHandler(client *openai.Client) *ChatHandler {
	model := &openai.ChatLanguageModel{
		Client: client,
		Model:  "gpt-3.5-turbo",
	}

	return &ChatHandler{
		Input:  os.Stdin,
		Output: os.Stdout,

		chain: llm.Chain(
			chat.Predict(model, ChatPrompt),
		),
	}
}

func (ch *ChatHandler) Run(proc goprocess.Process) {
	ctx := goprocessctx.OnClosingContext(proc)

	ch.ctx = llm.NewChainContext(ctx)

	inputStream := bufio.NewReader(ch.Input)

	history := chat.Message{}

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

		result := llm.GetOutput(ch.ctx, chat.ChatReplyContextKey)

		history.Entries = append(history.Entries, entry)
		history.Entries = append(history.Entries, result.Entries...)

		for _, replies := range result.Entries {
			_, err = fmt.Fprintf(ch.Output, "%s\n", replies)

			if err != nil {
				panic(err)
			}
		}
	}
}
