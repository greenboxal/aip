package impl

import (
	"context"
	"sync"

	"github.com/greenboxal/aip/pkg/ford/forddb"
	"github.com/greenboxal/aip/pkg/indexing"
)

type Session struct {
	*Iterator

	index   *Index
	context indexing.IndexContext
	options indexing.SessionOptions

	logMutex sync.RWMutex
	log      []indexing.Memory
}

func NewSession(ctx context.Context, index *Index, options indexing.SessionOptions) (*Session, error) {
	ictx, err := index.CreateContext(ctx)

	if err != nil {
		return nil, err
	}

	sess := &Session{
		index:   index,
		context: ictx,
		options: options,
	}

	sess.Iterator = NewIterator(index, ictx)

	sess.currentMemoryID = options.InitialMemoryID
	sess.rootMemoryID = options.RootMemoryID
	sess.branchMemoryID = options.BranchMemoryID

	return sess, nil
}

func (s *Session) Options() indexing.SessionOptions {
	return s.options
}

func (s *Session) Branch(ctx context.Context, clock, height int64) (indexing.Session, error) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	if s.current == nil {
		return nil, forddb.ErrNotFound
	}

	head := s.current.Fork(clock, height)

	branch, err := NewSession(ctx, s.index, s.options)

	if err != nil {
		return nil, err
	}

	branch.log = s.cloneLog()

	branchPtr := branch.appendToLog(head)

	if err = branch.setCurrent(branchPtr); err != nil {
		return nil, err
	}

	return nil, err
}

func (s *Session) Fork(ctx context.Context) (indexing.Session, error) {
	return s.Branch(ctx, 0, 1)
}

func (s *Session) Split(ctx context.Context) (indexing.Session, error) {
	return s.Branch(ctx, 0, 0)
}

func (s *Session) Push(data indexing.MemoryData) (indexing.Memory, error) {
	var head indexing.Memory

	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	if s.current == nil {
		head = indexing.NewMemory(s.MemoryAddress(), data)
	} else {
		head = s.current.Fork(1, 0)
		head.Data = data
	}

	headPtr := s.appendToLog(head)

	if err := s.setCurrent(headPtr); err != nil {
		return indexing.Memory{}, err
	}

	return head, nil
}

func (s *Session) UpdateMemoryData(data indexing.MemoryData) error {
	_, err := s.Push(data)

	if err != nil {
		return err
	}

	return nil
}

func (s *Session) Discard() error {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.discardLog()

	return s.set(
		s.rootMemoryID,
		s.branchMemoryID,
		s.branchMemoryID,
		s.branchMemoryID,
		s.currentClock,
		s.currentHeight,
		s.current,
	)
}

func (s *Session) Merge() error {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	head := s.current.Fork(1, -1)

	s.appendToLog(head)

	targets := s.cloneLog()
	segment := indexing.NewMemorySegment(targets...)

	rctx := &indexing.ReducerContext{
		Context: s.context.Context(),
		Session: s,
		Segment: segment,
	}

	reduced, err := s.index.cfg.Reducer.ReduceSegment(rctx)

	if err != nil {
		return err
	}

	if err := s.context.AppendSegment(reduced); err != nil {
		return err
	}

	s.discardLog()

	mergeTargetPtr := s.appendToLog(head)

	return s.setCurrent(mergeTargetPtr)
}

func (s *Session) appendToLog(memory indexing.Memory) *indexing.Memory {
	index := len(s.log)
	s.log = append(s.log, memory)

	return &s.log[index]
}

func (s *Session) discardLog() {
	s.log = s.log[0:0]
}

func (s *Session) cloneLog() []indexing.Memory {
	return append([]indexing.Memory(nil), s.log...)
}
