package llm

import (
	"context"

	"gonum.org/v1/gonum/mat"
)

type Embedder interface {
	MaxTokensPerChunk() int

	GetEmbeddings(ctx context.Context, chunks []string) ([]Embedding, error)
}

type Embedding struct {
	Embeddings []float32
}

func (e Embedding) Dim() int { return len(e.Embeddings) }

func (e Embedding) Float32() []float32 {
	return e.Embeddings
}

func (e Embedding) Float64() []float64 {
	floats := make([]float64, len(e.Embeddings))

	for i, v := range e.Embeddings {
		floats[i] = float64(v)
	}

	return floats
}

func (e Embedding) Vector() *mat.VecDense {
	return mat.NewVecDense(e.Dim(), e.Float64())
}
