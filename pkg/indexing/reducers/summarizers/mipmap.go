package summarizers

import (
	"github.com/greenboxal/aip/pkg/indexing"
	"github.com/greenboxal/aip/pkg/indexing/llm"
	"github.com/greenboxal/aip/pkg/indexing/reducers/chunkers"
)

type MipMapSummarizer struct {
	Summarizer Summarizer
	Tokenizer  llm.BasicTokenizer

	MinTokens int
	MaxTokens int
	MaxLevels int
}

type mipMapSummarizerContext struct {
	*indexing.ReducerContext

	MipMapSummarizer

	currentSegment           *indexing.MemorySegment
	currentSegmentTokenCount int
	currentDepth             int
}

func (ctx *mipMapSummarizerContext) setCurrentSegment(ms *indexing.MemorySegment) error {
	count, err := ms.CalculateTokenCount(ctx.Tokenizer)

	if err != nil {
		return err
	}

	ctx.currentSegment = ms
	ctx.currentSegmentTokenCount = count

	return nil
}

func (ctx *mipMapSummarizerContext) shouldReduce() bool {
	return ctx.currentSegment != nil && ctx.currentSegmentTokenCount > ctx.MinTokens && ctx.currentDepth < ctx.MaxLevels
}

func (ctx *mipMapSummarizerContext) reduceRound() error {
	ctx.currentDepth++

	targetSegment := ctx.currentSegment
	targetTokenCount := ctx.currentSegmentTokenCount / 2

	if ctx.currentSegmentTokenCount > ctx.Summarizer.MaxTokens() {
		split, _, err := chunkers.SplitSegment(ctx.currentSegment, ctx.Summarizer.MaxTokens()/2, 0)

		if err != nil {
			return err
		}

		targetSegment = split
		targetTokenCount = ctx.Summarizer.MaxTokens()
	}

	if targetTokenCount < ctx.MinTokens {
		targetTokenCount = ctx.MinTokens
	}

	if targetTokenCount > ctx.MaxTokens {
		targetTokenCount = ctx.MaxTokens
	}

	memories := make([]indexing.Memory, len(targetSegment.Memories))

	for i, m := range targetSegment.Memories {
		result, err := ctx.Summarizer.Summarize(
			ctx.Context,
			m.Data.Text,
			WithMaxTokens(targetTokenCount),
			WithContextHint(ctx.Hint),
		)

		if err != nil {
			return err
		}

		reduced, err := ctx.Session.Push(result)

		if err != nil {
			return err
		}

		memories[i] = reduced
	}

	newSegment := indexing.NewMemorySegment(memories...)

	return ctx.setCurrentSegment(newSegment)
}

func (ctx *mipMapSummarizerContext) Reduce() (*indexing.MemorySegment, error) {
	if err := ctx.setCurrentSegment(ctx.Segment); err != nil {
		return nil, err
	}

	for ctx.shouldReduce() {
		if err := ctx.reduceRound(); err != nil {
			return nil, err
		}
	}

	return ctx.currentSegment, nil
}

func (m *MipMapSummarizer) ReduceSegment(ctx *indexing.ReducerContext) (*indexing.MemorySegment, error) {
	sctx := &mipMapSummarizerContext{
		ReducerContext:   ctx,
		MipMapSummarizer: *m,
	}

	return sctx.Reduce()
}
