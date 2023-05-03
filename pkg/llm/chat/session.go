package chat

import (
	"context"

	"github.com/greenboxal/aip/pkg/llm"
)

type Session struct {
	Model LanguageModel

	history []MessageEntry
}

func NewSession(model LanguageModel) *Session {
	return &Session{
		Model: model,
	}
}

func (s *Session) AppendHistory(ctx context.Context, entry ...MessageEntry) {
	s.history = append(s.history, entry...)
}

func (s *Session) Predict(ctx context.Context, prompt Prompt, inputs ...llm.ChainInput) (result Message, pctx llm.ChainContext, err error) {
	pctx = llm.NewChainContext(ctx)

	pctx.SetInput(ChatHistoryContextKey, Message{Entries: s.history})

	for _, input := range inputs {
		pctx.SetInput(input.Key, input.Value)
	}

	promptText, err := prompt.Build(pctx)

	if err != nil {
		return result, nil, err
	}

	result, err = s.Model.PredictChat(ctx, promptText)

	return
}
