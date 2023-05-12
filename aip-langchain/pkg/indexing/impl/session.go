package impl

import (
	"context"
	"sync"

	collective2 "github.com/greenboxal/aip/aip-controller/pkg/collective"
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	indexing22 "github.com/greenboxal/aip/aip-langchain/pkg/indexing"
)

type Session struct {
	*Iterator

	index   *Index
	context indexing22.IndexContext
	options indexing22.SessionOptions

	commitMutex sync.RWMutex
	logMutex    sync.RWMutex
	log         []collective2.Memory
}

func NewSession(ctx context.Context, index *Index, options indexing22.SessionOptions) (*Session, error) {
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

func (s *Session) Options() indexing22.SessionOptions {
	return s.options
}

func (s *Session) Branch(ctx context.Context, clock, height int64) (indexing22.Session, error) {
	s.commitMutex.RLock()
	defer s.commitMutex.RUnlock()

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

func (s *Session) Fork(ctx context.Context) (indexing22.Session, error) {
	return s.Branch(ctx, 0, 1)
}

func (s *Session) Split(ctx context.Context) (indexing22.Session, error) {
	return s.Branch(ctx, 0, 0)
}

func (s *Session) Push(data collective2.MemoryData) (collective2.Memory, error) {
	var head collective2.Memory

	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	if s.current == nil {
		head = collective2.NewMemory(s.MemoryAddress(), data)
	} else {
		head = s.current.Fork(1, 0)
		head.Data = data
	}

	headPtr := s.appendToLog(head)

	if err := s.setCurrent(headPtr); err != nil {
		return collective2.Memory{}, err
	}

	return head, nil
}

func (s *Session) UpdateMemoryData(data collective2.MemoryData) error {
	_, err := s.Push(data)

	if err != nil {
		return err
	}

	return nil
}

func (s *Session) Discard() {
	s.commitMutex.Lock()
	defer s.commitMutex.Unlock()

	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.discardLog()

	err := s.set(
		s.rootMemoryID,
		s.branchMemoryID,
		s.branchMemoryID,
		s.branchMemoryID,
		s.currentClock,
		s.currentHeight,
		s.current,
	)

	if err != nil {
		panic(err)
	}
}

func (s *Session) Merge() error {
	s.commitMutex.Lock()
	defer s.commitMutex.Unlock()

	doMerge := func() (*collective2.MemorySegment, collective2.Memory) {
		s.logMutex.Lock()
		defer s.logMutex.Unlock()

		if len(s.log) == 0 {
			return nil, collective2.Memory{}
		}

		head := s.current.Fork(1, -1)

		s.appendToLog(head)

		targets := s.cloneLog()

		s.discardLog()

		return collective2.NewMemorySegment(targets...), head
	}

	segment, head := doMerge()

	if segment == nil {
		return nil
	}

	rctx := &indexing22.ReducerContext{
		Context: s.context.Context(),
		Session: s,
		Segment: segment,
	}

	reduced, err := s.index.cfg.Reducer.ReduceSegment(rctx)

	if err != nil {
		return err
	}

	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	if err := s.context.AppendSegment(reduced); err != nil {
		return err
	}

	mergeTargetPtr := s.appendToLog(head)

	return s.setCurrent(mergeTargetPtr)
}

func (s *Session) appendToLog(memory collective2.Memory) *collective2.Memory {
	index := len(s.log)
	s.log = append(s.log, memory)

	return &s.log[index]
}

func (s *Session) discardLog() {
	s.log = s.log[0:0]
}

func (s *Session) cloneLog() []collective2.Memory {
	return append([]collective2.Memory(nil), s.log...)
}
