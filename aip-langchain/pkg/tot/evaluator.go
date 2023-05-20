package tot

import (
	"context"

	"github.com/greenboxal/aip/aip-langchain/pkg/chunkers"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/chat"
	"github.com/greenboxal/aip/aip-langchain/pkg/tokenizers"
)

type Plan struct {
	Nodes []*Node
}

type Generator interface {
	Generate(ctx context.Context, parent *Node) (Plan, error)
}

type Scorer interface {
	Score(ctx context.Context) (float32, error)
}

type Context struct {
	Evaluator *Evaluator
	Context   context.Context
	Root      *Node
}

type Evaluator struct {
	Model     chat.LanguageModel
	Tokenizer tokenizers.BasicTokenizer
	Chunker   chunkers.Chunker
	Embedder  llm.Embedder

	Generator Generator

	Parallelism int
	MaxDepth    int
}

func (e *Evaluator) Evaluate(ctx context.Context, root *Node) (Plan, error) {
	return Plan{}, nil
}
