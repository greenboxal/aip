package reducers

import (
	"github.com/greenboxal/aip/pkg/indexing"
)

type SummaryDiffusionReducer struct {
	Summarizer Summarizer

	MaxChunkTokens int
	MaxTotalTokens int
	MaxOverlap     int
	MaxDepth       int
}

func (s *SummaryDiffusionReducer) ReduceSegment(ctx *indexing.ReducerContext) (*indexing.MemorySegment, error) {
	var sessionStack []indexing.Session

	depth := 0
	currentSession := ctx.Session
	currentRoot := ctx.Segment
	overlapFactor := 1 + s.MaxOverlap

	for depth < s.MaxDepth {
		split, totalTokens, err := SplitSegment(ctx.Segment, s.MaxChunkTokens, s.MaxOverlap)

		if err != nil {
			return nil, err
		}

		if totalTokens < s.MaxTotalTokens {
			break
		}

		factor := totalTokens / s.MaxTotalTokens * overlapFactor
		segments := split.PartitionEven(factor)
		memories := make([]indexing.Memory, len(segments))

		for i, segment := range segments {
			summary, err := s.Summarizer.Summarize(ctx.Context, segment, ctx.Hint, s.MaxTotalTokens)

			if err != nil {
				return nil, err
			}

			memories[i] = indexing.NewMemory(ctx.Session.MemoryAddress(), summary)
		}

		currentRoot = indexing.NewMemorySegment(memories...)
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
