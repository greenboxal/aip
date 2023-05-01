package llm

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type Embedder interface {
	MaxTokensPerChunk() int

	GetEmbeddings(ctx context.Context, chunks []string) ([]Embeddings, error)
}

type OpenAIEmbedder struct {
	Client *openai.Client
	Model  openai.EmbeddingModel
}

func (o *OpenAIEmbedder) MaxTokensPerChunk() int {
	return 2048
}

func (o *OpenAIEmbedder) GetEmbeddings(ctx context.Context, chunks []string) ([]Embeddings, error) {
	embeddings := make([]Embeddings, len(chunks))

	result, err := o.Client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: chunks,
		Model: o.Model,
	})

	for i, v := range result.Data {
		embeddings[i].Dim = 1536
		embeddings[i].Embeddings = v.Embedding
	}

	if err != nil {
		return nil, err
	}

	return embeddings, nil
}
