package openai

import (
	"context"

	"github.com/sashabaranov/go-openai"

	"github.com/greenboxal/aip/aip-langchain/pkg/llm"
)

type EmbeddingModel = openai.EmbeddingModel

const (
	AdaSimilarity         = openai.AdaSimilarity
	BabbageSimilarity     = openai.BabbageSimilarity
	CurieSimilarity       = openai.CurieSimilarity
	DavinciSimilarity     = openai.DavinciSimilarity
	AdaSearchDocument     = openai.AdaSearchDocument
	AdaSearchQuery        = openai.AdaSearchQuery
	BabbageSearchDocument = openai.BabbageSearchDocument
	BabbageSearchQuery    = openai.BabbageSearchQuery
	CurieSearchDocument   = openai.CurieSearchDocument
	CurieSearchQuery      = openai.CurieSearchQuery
	DavinciSearchDocument = openai.DavinciSearchDocument
	DavinciSearchQuery    = openai.DavinciSearchQuery
	AdaCodeSearchCode     = openai.AdaCodeSearchCode
	AdaCodeSearchText     = openai.AdaCodeSearchText
	BabbageCodeSearchCode = openai.BabbageCodeSearchCode
	BabbageCodeSearchText = openai.BabbageCodeSearchText
	AdaEmbeddingV2        = openai.AdaEmbeddingV2
)

type Embedder struct {
	Client *openai.Client
	Model  openai.EmbeddingModel
}

func (o *Embedder) MaxTokensPerChunk() int {
	return 2048
}

func (o *Embedder) GetEmbeddings(ctx context.Context, chunks []string) ([]llm.Embeddings, error) {
	embeddings := make([]llm.Embeddings, len(chunks))

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
