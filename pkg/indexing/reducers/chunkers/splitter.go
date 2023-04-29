package chunkers

import (
	"github.com/greenboxal/aip/pkg/indexing"
	"github.com/greenboxal/aip/pkg/utils"
)

func SplitSegment(segment *indexing.MemorySegment, maxTokens int, maxOverlap int) (*indexing.MemorySegment, int, error) {
	index := 0
	totalTokens := 0
	memories := make([]indexing.Memory, 0, len(segment.Memories))

	for _, m := range segment.Memories {
		chunks, err := SplitTextIntoChunks(m.Data.Text, maxTokens, maxOverlap)

		if err != nil {
			return nil, totalTokens, err
		}

		baseIndex := index
		index += len(chunks)

		memories = utils.Grow(memories, len(chunks))

		for i, chunk := range chunks {
			pointer := indexing.Absolute(
				indexing.Anchors(m.RootMemoryID, m.BranchMemoryID, m.ParentMemoryID),
				indexing.Relative(m.Clock, m.Height),
			)

			data := indexing.NewMemoryData(chunk)

			memories[baseIndex+i] = indexing.NewMemory(pointer, data)

			totalTokens += len(chunk)
		}
	}

	return indexing.NewMemorySegment(memories...), totalTokens, nil
}
