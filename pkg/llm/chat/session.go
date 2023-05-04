package chat

import (
	"context"

	"github.com/greenboxal/aip/pkg/llm/chain"
)

type Session struct {
	Model LanguageModel

	ctx     chain.ChainContext
	history []MessageEntry
}

func NewSession(ctx chain.ChainContext, model LanguageModel) *Session {
	return &Session{
		Model: model,

		ctx: ctx,
	}
}

func (s *Session) History() Message                    { return Message{Entries: s.history} }
func (s *Session) AppendHistory(entry ...MessageEntry) { s.history = append(s.history, entry...) }

func (s *Session) Predict(ctx context.Context, prompt Prompt) (result Message, pctx chain.ChainContext, err error) {
	pctx.SetInput(ChatHistoryContextKey, Message{Entries: s.history})

	promptText, err := prompt.Build(pctx)

	if err != nil {
		return result, nil, err
	}

	result, err = s.Model.PredictChat(ctx, promptText)

	if err != nil {
		return result, nil, err
	}

	s.AppendHistory(result.Entries...)

	return
}
