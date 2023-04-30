package chunkers

import (
	"github.com/greenboxal/aip/pkg/collective"
	"github.com/greenboxal/aip/pkg/utils"
)

func SplitSegment(segment *collective.MemorySegment, maxTokens int, maxOverlap int) (*collective.MemorySegment, int, error) {
	index := 0
	totalTokens := 0
	memories := make([]collective.Memory, 0, len(segment.Memories))

	for _, m := range segment.Memories {
		chunks, err := SplitTextIntoChunks(m.Data.Text, maxTokens, maxOverlap)

		if err != nil {
			return nil, totalTokens, err
		}

		baseIndex := index
		index += len(chunks)

		memories = utils.Grow(memories, len(chunks))

		for i, chunk := range chunks {
			pointer := collective.Absolute(
				collective.Anchors(m.RootMemoryID, m.BranchMemoryID, m.ParentMemoryID),
				collective.Relative(m.Clock, m.Height),
			)

			data := collective.NewMemoryData(chunk)

			memories[baseIndex+i] = collective.NewMemory(pointer, data)

			totalTokens += len(chunk)
		}
	}

	return collective.NewMemorySegment(memories...), totalTokens, nil
}
