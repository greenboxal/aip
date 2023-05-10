package openai

import (
	"context"

	"github.com/sashabaranov/go-openai"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
	"github.com/greenboxal/aip/aip-controller/pkg/llm"
	chat "github.com/greenboxal/aip/aip-controller/pkg/llm/chat"
)

type LanguageModel struct {
	Client *Client

	Model       string
	Temperature float32
}

func (lm *LanguageModel) MaxTokens() int {
	// FIXME: Wrong
	return 4096
}

func (lm *LanguageModel) Predict(ctx context.Context, prompt string, options ...llm.PredictOption) (string, error) {
	opts := llm.NewPredictOptions(options...)

	result, err := lm.Client.CreateCompletion(ctx, openai.CompletionRequest{
		Model:       lm.Model,
		Temperature: opts.Temperature,
		MaxTokens:   opts.MaxTokens,
		Prompt:      prompt,
	})

	if err != nil {
		return "", err
	}

	return result.Choices[0].Text, nil
}

func (lm *LanguageModel) PredictChat(ctx context.Context, msg chat.Message, options ...llm.PredictOption) (chat.Message, error) {
	result, err := lm.Predict(ctx, msg.String(), options...)

	if err != nil {
		return chat.Message{}, err
	}

	return chat.Message{
		Entries: []chat.MessageEntry{
			{
				Role: msn.RoleAI,
				Text: result,
			},
		},
	}, nil
}

type ChatLanguageModel struct {
	Client *Client

	Model       string
	Temperature float32
}

func (lm *ChatLanguageModel) MaxTokens() int {
	// FIXME: Return actual number
	return 4096
}

func (lm *ChatLanguageModel) Predict(ctx context.Context, prompt string, options ...llm.PredictOption) (string, error) {
	msg := chat.Compose(
		chat.Entry(msn.RoleUser, prompt),
		chat.Entry(msn.RoleAI, ""),
	)

	result, err := lm.PredictChat(ctx, msg, options...)

	if err != nil {
		return "", nil
	}

	return result.Entries[0].Text, nil
}

func (lm *ChatLanguageModel) PredictChat(ctx context.Context, msg chat.Message, options ...llm.PredictOption) (chat.Message, error) {
	opts := llm.NewPredictOptions(options...)
	messages := make([]openai.ChatCompletionMessage, len(msg.Entries))

	for i, m := range msg.Entries {
		var role string

		switch m.Role {
		case msn.RoleSystem:
			role = "system"
		case msn.RoleUser:
			role = "user"
		case msn.RoleAI:
			role = "assistant"
		default:
			panic("unknown role")
		}

		messages[i] = openai.ChatCompletionMessage{
			Role:    role,
			Content: m.Text,
		}
	}

	result, err := lm.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       lm.Model,
		Temperature: opts.Temperature,
		MaxTokens:   opts.MaxTokens,
		Messages:    messages,
	})

	if err != nil {
		return chat.Message{}, nil
	}

	entries := make([]chat.MessageEntry, len(result.Choices))

	for i, c := range result.Choices {
		var role msn.Role

		switch c.Message.Role {
		case "system":
			role = msn.RoleSystem
		case "user":
			role = msn.RoleUser
		case "assistant":
			role = msn.RoleAI
		default:
			panic("unknown role")
		}

		entries[i] = chat.MessageEntry{
			Role: role,
			Name: c.Message.Name,
			Text: c.Message.Content,
		}
	}

	return chat.Compose(entries...), nil
}
