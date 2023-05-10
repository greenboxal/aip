package reducers

import (
	"errors"

	collective2 "github.com/greenboxal/aip/aip-controller/pkg/collective"
	indexing2 "github.com/greenboxal/aip/aip-controller/pkg/indexing"
	"github.com/greenboxal/aip/aip-controller/pkg/llm"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/chunkers"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/summarizers"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/tokenizers"
)

type SummaryDiffusionReducer struct {
	Summarizer summarizers.Summarizer
	Embedder   llm.Embedder
	Tokenizer  tokenizers.BasicTokenizer

	MinTokens int
	MaxTokens int

	MaxChunkTokens int
	MaxTotalTokens int
	MaxOverlap     int
	MaxDepth       int
}

func (s *SummaryDiffusionReducer) ReduceSegment(ctx *indexing2.ReducerContext) (*collective2.MemorySegment, error) {
	sctx := &summaryDiffusionReducerContext{
		ReducerContext: ctx,
		reducer:        s,
	}

	return sctx.Reduce()
}

type summaryDiffusionReducerContext struct {
	*indexing2.ReducerContext

	reducer *SummaryDiffusionReducer

	currentSegment           *collective2.MemorySegment
	currentSegmentTokenCount int
	currentDepth             int
}

func (ctx *summaryDiffusionReducerContext) setCurrentSegment(ms *collective2.MemorySegment) error {
	count, err := ms.CalculateTokenCount(ctx.reducer.Tokenizer)

	if err != nil {
		return err
	}

	ctx.currentSegment = ms
	ctx.currentSegmentTokenCount = count

	return nil
}

func (ctx *summaryDiffusionReducerContext) shouldReduce() bool {
	return ctx.currentSegment != nil && ctx.currentSegmentTokenCount > ctx.reducer.MinTokens
}

func (ctx *summaryDiffusionReducerContext) reduceRound() error {
	ctx.currentDepth++

	return errors.New("not implemented")
}

func (ctx *summaryDiffusionReducerContext) Reduce() (*collective2.MemorySegment, error) {
	var sessionStack []indexing2.Session

	if err := ctx.setCurrentSegment(ctx.Segment); err != nil {
		return nil, err
	}

	for ctx.shouldReduce() {
		if err := ctx.reduceRound(); err != nil {
			return nil, err
		}
	}

	s := ctx.reducer

	depth := 0
	currentSession := ctx.Session
	currentRoot := ctx.Segment
	overlapFactor := 1 + s.MaxOverlap

	// FIXME: mipmap instead
	for depth < s.MaxDepth {
		split, totalTokens, err := chunkers.SplitSegment(ctx.Segment, s.MaxChunkTokens, s.MaxOverlap)

		if err != nil {
			return nil, err
		}

		if totalTokens < s.MaxTotalTokens {
			break
		}

		factor := totalTokens / s.MaxTotalTokens * overlapFactor
		segments := split.PartitionEven(factor)
		memories := make([]collective2.Memory, len(segments))

		for i := range segments {
			summary, err := s.Summarizer.Summarize(
				ctx.Context,
				"",
				summarizers.WithContextHint(ctx.Hint),
				summarizers.WithMaxTokens(s.MaxTotalTokens),
			)

			if err != nil {
				return nil, err
			}

			memories[i] = collective2.NewMemory(ctx.Session.MemoryAddress(), summary)
		}

		currentRoot = collective2.NewMemorySegment(memories...)
		currentSession, err = currentSession.Fork(ctx.Context)
		sessionStack = append(sessionStack, currentSession)

		if err != nil {
			return nil, err
		}

		depth++
	}

	for i := len(sessionStack) - 1; i >= 0; i-- {
		if err := sessionStack[i].Merge(); err != nil {
			return nil, err
		}
	}

	return currentRoot, nil
}
