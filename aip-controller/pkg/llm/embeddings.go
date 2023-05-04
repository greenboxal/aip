package llm

import "context"

type Embedder interface {
	MaxTokensPerChunk() int

	GetEmbeddings(ctx context.Context, chunks []string) ([]Embeddings, error)
}

type Embeddings struct {
	Dim        int
	Embeddings []float32
}

func (e Embeddings) Len() int        { return len(e.Embeddings) }
func (e Embeddings) TokenCount() int { return len(e.Embeddings) / e.Dim }
