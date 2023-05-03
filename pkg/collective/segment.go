package collective

import (
	"github.com/greenboxal/aip/pkg/ford/forddb"
	"github.com/greenboxal/aip/pkg/llm/tokenizers"
	"github.com/greenboxal/aip/pkg/utils"
)

type MemorySet interface {
	Slice(start, end int) *MemorySegment
	PartitionEven(count int) []*MemorySegment
}

type MemorySegmentID struct {
	forddb.StringResourceID[*Memory]
}

type MemorySegment struct {
	forddb.ResourceMetadata[MemorySegmentID, *MemorySegment] `json:"metadata"`

	Memories []Memory
}

func (ms *MemorySegment) Slice(start, end int) *MemorySegment {
	return NewMemorySegment(ms.Memories[start:end]...)
}

func (ms *MemorySegment) PartitionEven(count int) []*MemorySegment {
	partitionSize := len(ms.Memories) / count
	partitions := make([]*MemorySegment, count)

	for i := range partitions {
		partitions[i] = ms.Slice(
			utils.Min(i*partitionSize, len(ms.Memories)),
			utils.Min((i+1)*partitionSize, len(ms.Memories)),
		)
	}

	return partitions
}

func (ms *MemorySegment) CalculateTokenCount(tokenizer tokenizers.BasicTokenizer) (int, error) {
	total := 0

	for _, m := range ms.Memories {
		tokens, err := m.CalculateTokenCount(tokenizer)

		if err != nil {
			return 0, err
		}

		total += tokens
	}

	return total, nil
}

func NewMemorySegment(memories ...Memory) *MemorySegment {
	return &MemorySegment{
		Memories: memories,
	}
}
