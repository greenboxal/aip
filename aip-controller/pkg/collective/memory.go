package collective

import (
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-langchain/pkg/tokenizers"
)

type MemoryID struct {
	forddb.StringResourceID[*Memory] `ipld:",inline"`
}

type Memory struct {
	forddb.ResourceBase[MemoryID, *Memory] `json:"metadata"`

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

func (m *Memory) CalculateTokenCount(tokenizer tokenizers.BasicTokenizer) (int, error) {
	return tokenizer.Count(m.Data.Text)
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
