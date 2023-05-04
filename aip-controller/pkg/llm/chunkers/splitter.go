package chunkers

import (
	collective2 "github.com/greenboxal/aip/aip-controller/pkg/collective"
	"github.com/greenboxal/aip/aip-controller/pkg/utils"
)

func SplitSegment(segment *collective2.MemorySegment, maxTokens int, maxOverlap int) (*collective2.MemorySegment, int, error) {
	index := 0
	totalTokens := 0
	memories := make([]collective2.Memory, 0, len(segment.Memories))

	for _, m := range segment.Memories {
		chunks, err := SplitTextIntoChunks(m.Data.Text, maxTokens, maxOverlap)

		if err != nil {
			return nil, totalTokens, err
		}

		baseIndex := index
		index += len(chunks)

		memories = utils.Grow(memories, len(chunks))

		for i, chunk := range chunks {
			pointer := collective2.Absolute(
				collective2.Anchors(m.RootMemoryID, m.BranchMemoryID, m.ParentMemoryID),
				collective2.Relative(m.Clock, m.Height),
			)

			data := collective2.NewMemoryData(chunk)

			memories[baseIndex+i] = collective2.NewMemory(pointer, data)

			totalTokens += len(chunk)
		}
	}

	return collective2.NewMemorySegment(memories...), totalTokens, nil
}
