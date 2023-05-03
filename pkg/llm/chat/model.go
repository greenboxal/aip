package chat

import (
	"context"

	"github.com/greenboxal/aip/pkg/llm"
)

type LanguageModel interface {
	MaxTokens() int

	PredictChat(ctx context.Context, msg Message, options ...llm.PredictOption) (Message, error)
}
