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

func (m *Memory) Fork(clock, height uint64) Memory {
	return Memory{
		RootMemoryID:   m.RootMemoryID,
		BranchMemoryID: m.BranchMemoryID,
		ParentMemoryID: m.ID,
		Clock:          m.Clock + clock,
		Height:         m.Height + height,
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

func NewMemorySegment(memories ...Memory) *MemorySegment {
	return &MemorySegment{
		Memories: memories,
	}
}

func init() {
	forddb.DefineResourceType[MemoryID, *Memory]("memory")
	forddb.DefineResourceType[MemorySegmentID, *MemorySegment]("memory_segment")
}
