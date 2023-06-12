package forddbimpl

import (
	"context"

	"github.com/ipld/go-ipld-prime"
)

type ObjectID = ipld.Link

type Object interface {
	ObjectID() ObjectID
	ObjectLink() ipld.Link
}

type Tree interface {
	Object
}

type Resource interface {
	Object

	AsNode() ipld.Node
}

type TreeEntry struct {
	Name string `json:"name"`

	Node ipld.Node `json:"node"`
	Link ipld.Link `json:"link"`
}

type TreeIterator interface {
	Next() bool
	Entry() TreeEntry
}

type NodeGetter interface {
	Get(ctx context.Context, link ipld.Link) (ipld.Node, error)
}
