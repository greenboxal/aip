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
	outputStream := bufio.NewWriter(ch.Output)

	history := chat.Message{}

	for {
		select {
		case <-ctx.Done():
			return

		default:
		}

		line, err := inputStream.ReadString('\n')

		if err != nil {
			panic(err)
		}

		entry := chat.MessageEntry{
			Role:    chat.RoleUser,
			Content: line,
		}

		ch.ctx.Flip()
		ch.ctx.SetInput(chat.ChatReplyContextKey, entry)
		ch.ctx.SetInput(chat.ChatHistoryContextKey, history)

		if err := ch.chain.Run(ch.ctx); err != nil {
			_, err = outputStream.WriteString(fmt.Sprintf("ERROR: %s\n", err))

			if err != nil {
				panic(err)
			}
		}

		result := llm.GetOutput(ch.ctx, chat.ChatReplyContextKey)

		newEntries := append([]chat.MessageEntry{entry}, result.Entries...)
		history.Entries = append(history.Entries, newEntries...)

		for _, entry := range newEntries {
			_, err = fmt.Fprintf(ch.Output, "%s\n", entry)

			if err != nil {
				panic(err)
			}
		}
	}
}
