package indexing

import (
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"

	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type MemoryDataID struct {
	forddb.CidResourceID[*Memory]
}

type MemoryData struct {
	forddb.ResourceMetadata[MemoryDataID, *MemoryData] `json:"metadata"`

	Cid  cid.Cid `json:"cid"`
	Data []byte  `json:"data"`
}

func NewMemoryData(data []byte) MemoryData {
	h, err := multihash.Sum(data, multihash.SHA2_256, -1)

	if err != nil {
		panic(err)
	}

	return MemoryData{
		Cid:  cid.NewCidV1(cid.Raw, h),
		Data: data,
	}
}

type MemoryID struct {
	forddb.StringResourceID[*Memory]
}

type Memory struct {
	forddb.ResourceMetadata[MemoryID, *Memory] `json:"metadata"`

	RootMemoryID   MemoryID   `json:"root_memory_id"`
	BranchMemoryID MemoryID   `json:"branch_memory_id"`
	ParentMemoryID MemoryID   `json:"parent_memory_id"`
	Clock          uint64     `json:"clock"`
	Height         uint64     `json:"height"`
	Data           MemoryData `json:"data"`
}

func NewMemory(pointer MemoryAbsoluteAddress, data MemoryData) Memory {
	return Memory{
		RootMemoryID:   pointer.GetRootMemoryID(),
		BranchMemoryID: pointer.GetBranchMemoryID(),
		ParentMemoryID: pointer.GetParentMemoryID(),
		Clock:          pointer.GetClock(),
		Height:         pointer.GetHeight(),
		Data:           data,
	}
}

func (m *Memory) GetMemoryID() MemoryID {
	return m.ID
}

func (m *Memory) GetRootMemoryID() MemoryID {
	return m.RootMemoryID
}

func (m *Memory) GetBranchMemoryID() MemoryID {
	return m.BranchMemoryID
}

func (m *Memory) GetParentMemoryID() MemoryID {
	return m.ParentMemoryID
}

func (m *Memory) GetMemoryAddress() MemoryAbsoluteAddress {
	return Absolute(
		Anchors(m.RootMemoryID, m.BranchMemoryID, m.ParentMemoryID),
		Relative(m.Clock, m.Height),
	)
}

func (m *Memory) Fork(clock, height int64) Memory {
	return Memory{
		RootMemoryID:   m.RootMemoryID,
		BranchMemoryID: m.BranchMemoryID,
		ParentMemoryID: m.ID,
		Clock:          m.Clock + uint64(clock),
		Height:         m.Height + uint64(height),
		Data:           m.Data,
	}
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
			min(i*partitionSize, len(ms.Memories)),
			min((i+1)*partitionSize, len(ms.Memories)),
		)
	}

	return partitions
}

func NewMemorySegment(memories ...Memory) *MemorySegment {
	return &MemorySegment{
		Memories: memories,
	}
}

func init() {
	forddb.DefineResourceType[MemoryID, *Memory]("memory")
	forddb.DefineResourceType[MemorySegmentID, *MemorySegment]("memory_segment")
}

func min(i int, i2 int) int {
	if i < i2 {
		return i
	}

	return i2
}
