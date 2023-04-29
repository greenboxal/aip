package reducers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"

	"github.com/greenboxal/aip/pkg/indexing"
)

type Summarizer interface {
	Summarize(
		ctx context.Context,
		segment *indexing.MemorySegment,
		contextHint string,
		maxTokens int,
	) (indexing.MemoryData, error)
}

type ChatGptSummarizer struct {
	Client *openai.Client
	Model  string
}

func (gs *ChatGptSummarizer) Summarize(
	ctx context.Context,
	segment *indexing.MemorySegment,
	contextHint string,
	maxTokens int,
) (indexing.MemoryData, error) {
	documents := make([]string, len(segment.Memories))

	for i := range segment.Memories {
		documents[i] = string(segment.Memories[i].Data.Data)
	}

	joinedDocuments := strings.Join(documents, "\n\n")

	result, err := gs.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       gs.Model,
		Temperature: 0.70,
		MaxTokens:   maxTokens,

		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You are an AI Assistant trained to best summarize and extract deep connections between related documents.",
			},
			{
				Role: "user",
				Content: fmt.Sprintf(
					"Please summarize the document excerpts below optimizing the signal to noise ratio with regards to the following context: %s\nDocuments:%s",
					contextHint,
					joinedDocuments,
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

	return indexing.NewMemoryData([]byte(choice.Message.Content)), nil
}
