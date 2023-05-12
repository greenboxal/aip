package chat

import (
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/chain"
)

type Memory interface {
	Load(ctx chain.ChainContext) (Message, error)
	Append(ctx chain.ChainContext, msg Message) error
}

type InMemoryChatMemory struct {
	ContextKey chain.ContextKey[Message]

	Messages []MessageEntry
}

func NewInMemoryHistory(key chain.ContextKey[Message]) *InMemoryChatMemory {
	return &InMemoryChatMemory{ContextKey: key}
}

func (i *InMemoryChatMemory) Load(ctx chain.ChainContext) (Message, error) {
	msg := Message{Entries: i.Messages}

	ctx.SetInput(i.ContextKey, msg)

	return msg, nil
}

func (i *InMemoryChatMemory) Append(ctx chain.ChainContext, msg Message) error {
	i.Messages = append(i.Messages, msg.Entries...)

	return nil
}
