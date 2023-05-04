package openai

import (
	"context"

	"github.com/sashabaranov/go-openai"

	"github.com/greenboxal/aip/aip-controller/pkg/llm"
	chat2 "github.com/greenboxal/aip/aip-controller/pkg/llm/chat"
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

func (lm *LanguageModel) PredictChat(ctx context.Context, msg chat2.Message, options ...llm.PredictOption) (chat2.Message, error) {
	result, err := lm.Predict(ctx, msg.String(), options...)

	if err != nil {
		return chat2.Message{}, err
	}

	return chat2.Message{
		Entries: []chat2.MessageEntry{
			{
				Role:    chat2.RoleAI,
				Content: result,
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
	msg := chat2.Compose(
		chat2.Entry(chat2.RoleUser, prompt),
		chat2.Entry(chat2.RoleAI, ""),
	)

	result, err := lm.PredictChat(ctx, msg, options...)

	if err != nil {
		return "", nil
	}

	return result.Entries[0].Content, nil
}

func (lm *ChatLanguageModel) PredictChat(ctx context.Context, msg chat2.Message, options ...llm.PredictOption) (chat2.Message, error) {
	opts := llm.NewPredictOptions(options...)
	messages := make([]openai.ChatCompletionMessage, len(msg.Entries))

	for i, m := range msg.Entries {
		var role string

		switch m.Role {
		case chat2.RoleSystem:
			role = "system"
		case chat2.RoleUser:
			role = "user"
		case chat2.RoleAI:
			role = "assistant"
		default:
			panic("unknown role")
		}

		messages[i] = openai.ChatCompletionMessage{
			Role:    role,
			Content: m.Content,
		}
	}

	result, err := lm.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       lm.Model,
		Temperature: opts.Temperature,
		MaxTokens:   opts.MaxTokens,
		Messages:    messages,
	})

	if err != nil {
		return chat2.Message{}, nil
	}

	entries := make([]chat2.MessageEntry, len(result.Choices))

	for i, c := range result.Choices {
		var role chat2.Role

		switch c.Message.Role {
		case "system":
			role = chat2.RoleSystem
		case "user":
			role = chat2.RoleUser
		case "assistant":
			role = chat2.RoleAI
		default:
			panic("unknown role")
		}

		entries[i] = chat2.MessageEntry{
			Role:    role,
			Name:    c.Message.Name,
			Content: c.Message.Content,
		}
	}

	return chat2.Compose(entries...), nil
}
