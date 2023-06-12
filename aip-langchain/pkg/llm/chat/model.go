package chat

import (
	"context"

	"github.com/greenboxal/aip/aip-langchain/pkg/llm"
)

type MessageStream interface {
	Recv() (MessageFragment, error)
	Close() error
}

type MessageFragment struct {
	MessageIndex int
	Delta        string
}

type LanguageModel interface {
	MaxTokens() int

	PredictChat(ctx context.Context, msg Message, options ...llm.PredictOption) (Message, error)
	PredictChatStream(ctx context.Context, msg Message, options ...llm.PredictOption) (MessageStream, error)
}
