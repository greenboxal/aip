package indexing

import "context"

type SessionOptions struct {
	Namespace       string
	ReadOnly        bool
	RootMemoryID    MemoryID
	BranchMemoryID  MemoryID
	InitialMemoryID MemoryID
}

type Session interface {
	Iterator

	Options() SessionOptions

	Branch(ctx context.Context, clock, height int64) (Session, error)
	Fork(ctx context.Context) (Session, error)
	Split(ctx context.Context) (Session, error)
	Push(data MemoryData) (Memory, error)

	UpdateMemoryData(data MemoryData) error

	Discard() error
	Merge() error
}
