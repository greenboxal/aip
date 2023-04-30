package indexing

import (
	"context"

	"github.com/greenboxal/aip/pkg/collective"
)

type SessionOptions struct {
	Namespace       string
	ReadOnly        bool
	RootMemoryID    collective.MemoryID
	BranchMemoryID  collective.MemoryID
	InitialMemoryID collective.MemoryID
}

type Session interface {
	Iterator

	Options() SessionOptions

	Branch(ctx context.Context, clock, height int64) (Session, error)
	Fork(ctx context.Context) (Session, error)
	Split(ctx context.Context) (Session, error)
	Push(data collective.MemoryData) (collective.Memory, error)

	UpdateMemoryData(data collective.MemoryData) error

	Discard()
	Merge() error
}
