package impl

import (
	"github.com/greenboxal/aip/pkg/collective"
	"github.com/greenboxal/aip/pkg/indexing"
)

type Iterator struct {
	index   *Index
	context indexing.IndexContext

	rootMemoryID    collective.MemoryID
	branchMemoryID  collective.MemoryID
	parentMemoryID  collective.MemoryID
	currentMemoryID collective.MemoryID

	currentHeight uint64
	currentClock  uint64

	current *collective.Memory
}

func NewIterator(index *Index, ctx indexing.IndexContext) *Iterator {
	return &Iterator{
		index:   index,
		context: ctx,
	}
}

func (s *Iterator) Index() indexing.Index {
	return s.index
}

func (s *Iterator) MemoryAddress() collective.MemoryAbsoluteAddress {
	return collective.Absolute(
		collective.Anchors(s.rootMemoryID, s.branchMemoryID, s.parentMemoryID),
		collective.Relative(s.currentClock, s.currentHeight),
	)
}

func (s *Iterator) RootMemoryID() collective.MemoryID {
	return s.rootMemoryID
}

func (s *Iterator) BranchMemoryID() collective.MemoryID {
	return s.branchMemoryID
}

func (s *Iterator) ParentMemoryID() collective.MemoryID {
	return s.parentMemoryID
}

func (s *Iterator) CurrentMemoryID() collective.MemoryID {
	return s.currentMemoryID
}

func (s *Iterator) SeekRelative(offset int) error {
	return s.setCoordinates(s.currentHeight, s.currentClock+uint64(offset))
}

func (s *Iterator) SeekTo(id collective.MemoryID) error {
	m, err := s.context.ResolveMemory(id)

	if err != nil {
		return err
	}

	return s.setCurrent(m)
}

func (s *Iterator) GetCurrentMemory() collective.Memory {
	if s.current == nil {
		panic("no current memory")
	}

	return *s.current
}

func (s *Iterator) GetCurrentMemoryData() collective.MemoryData {
	if s.current == nil {
		return collective.MemoryData{}
	}

	return s.current.Data
}

func (s *Iterator) setCurrent(m *collective.Memory) error {
	return s.set(m.RootMemoryID, m.BranchMemoryID, m.ParentMemoryID, m.ID, m.Height, m.Clock, m)
}

func (s *Iterator) set(root, branch, parent, current collective.MemoryID, height, clock uint64, m *collective.Memory) error {
	s.current = m
	s.currentHeight = height
	s.currentClock = clock
	s.rootMemoryID = root
	s.branchMemoryID = branch
	s.parentMemoryID = parent
	s.currentMemoryID = current

	return s.checkCurrent()
}

func (s *Iterator) setPointers(root, branch, parent, current collective.MemoryID) error {
	s.rootMemoryID = root
	s.branchMemoryID = branch
	s.parentMemoryID = parent
	s.currentMemoryID = current

	return s.checkCurrent()
}

func (s *Iterator) setCoordinates(height, clock uint64) error {
	s.currentHeight = height
	s.currentClock = clock

	return s.checkCurrent()
}

func (s *Iterator) checkCurrent() error {
	if !s.isValid() {
		if err := s.invalidate(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Iterator) invalidate() error {
	query := indexing.ResolveOptions{
		Context: s.context.Context(),

		BaseMemoryID:   s.currentMemoryID,
		TargetMemoryID: nil,

		BaseMemoryHint: s.current,
	}

	memory, err := s.context.ResolveMemoryRelative(query)

	if err != nil {
		return err
	}

	s.current = memory

	return nil
}

func (s *Iterator) isValid() bool {
	if s.current == nil {
		return false
	}

	if s.currentMemoryID != s.current.ID {
		return false
	}

	if s.parentMemoryID != s.current.ParentMemoryID {
		return false
	}

	if s.branchMemoryID != s.current.BranchMemoryID {
		return false
	}

	if s.rootMemoryID != s.current.RootMemoryID {
		return false
	}

	if s.currentHeight != s.current.Height {
		return false
	}

	if s.currentClock != s.current.Clock {
		return false
	}

	return true
}
