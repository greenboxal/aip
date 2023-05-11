package indexing2

import (
	"context"

	collective2 "github.com/greenboxal/aip/aip-controller/pkg/collective"
)

type SessionOptions struct {
	Namespace       string
	ReadOnly        bool
	RootMemoryID    collective2.MemoryID
	BranchMemoryID  collective2.MemoryID
	InitialMemoryID collective2.MemoryID
}

type Session interface {
	Iterator

	Options() SessionOptions

	Branch(ctx context.Context, clock, height int64) (Session, error)
	Fork(ctx context.Context) (Session, error)
	Split(ctx context.Context) (Session, error)
	Push(data collective2.MemoryData) (collective2.Memory, error)

	UpdateMemoryData(data collective2.MemoryData) error

	Discard()
	Merge() error
}
