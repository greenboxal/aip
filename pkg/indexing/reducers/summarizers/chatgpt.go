package summarizers

import (
	"context"
	"errors"
	"fmt"

	"github.com/sashabaranov/go-openai"

	"github.com/greenboxal/aip/pkg/indexing"
)

type ChatGptSummarizer struct {
	Client *openai.Client

	Model string
}

func (gs *ChatGptSummarizer) MaxTokens() int {
	return 4096 - 100
}

func (gs *ChatGptSummarizer) Summarize(
	ctx context.Context,
	document string,
	options ...SummarizeOption,
) (indexing.MemoryData, error) {
	opts := SummarizeOptions{
		Temperature: 0.0,
		MaxTokens:   256,
	}

	for _, opt := range options {
		opt(&opts)
	}

	result, err := gs.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       gs.Model,
		Temperature: opts.Temperature,
		MaxTokens:   opts.MaxTokens,

		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You are an AI Assistant trained to best summarize and extract deep connections between related documents.",
			},
			{
				Role: "user",
				Content: fmt.Sprintf(
					"Please summarize the document excerpts below optimizing the signal to noise ratio with regards to the following context: %s\nDocuments:%s",
					opts.ContextHint,
					document,
				),
			},
			{
				Role:    "assistant",
				Content: "",
			},
		},
	})

	if err != nil {
		return indexing.MemoryData{}, err
	}

	if len(result.Choices) == 0 {
		return indexing.MemoryData{}, errors.New("no choices")
	}

	choice := result.Choices[0]

	return indexing.NewMemoryData(choice.Message.Content), nil
}
