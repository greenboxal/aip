package reducers

import (
	"github.com/greenboxal/aip/pkg/indexing"
)

func SplitSegment(segment *indexing.MemorySegment, maxTokens int, maxOverlap int) (*indexing.MemorySegment, int, error) {
	index := 0
	totalTokens := 0
	memories := make([]indexing.Memory, 0, len(segment.Memories))

	for _, m := range segment.Memories {
		chunks, err := SplitTextIntoChunks(string(m.Data.Data), maxTokens, maxOverlap)

		if err != nil {
			return nil, totalTokens, err
		}

		baseIndex := index
		index += len(chunks)

		memories = Grow(memories, len(chunks))

		for i, chunk := range chunks {
			pointer := indexing.Absolute(
				indexing.Anchors(m.RootMemoryID, m.BranchMemoryID, m.ParentMemoryID),
				indexing.Relative(m.Clock, m.Height),
			)

			data := indexing.NewMemoryData([]byte(chunk))

			memories[baseIndex+i] = indexing.NewMemory(pointer, data)

			totalTokens += len(chunk)
		}
	}

	return indexing.NewMemorySegment(memories...), totalTokens, nil
}
