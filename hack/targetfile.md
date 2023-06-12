# @package: forddbimpl

## @file: diff.go

```go
package forddbimpl

import (
	"strconv"

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
	Name string    `json:"name"`
	Link ipld.Link `json:"link"`
}

type TreeIterator interface {
	Next() bool
	Entry() TreeEntry
}

type TreeDiff struct {
	Objects map[string]ObjectDiff `json:"objects"`
}

type ObjectDiff struct {
	Old     ObjectID       `json:"old_id"`
	New     ObjectID       `json:"new_id"`
	Changes []ObjectChange `json:"changes"`
}

type ObjectChange struct {
	Path string    `json:"path"`
	Old  ipld.Node `json:"old"`
	New  ipld.Node `json:"new"`
}

// Diff compares two trees and returns a diff between them.
// The diff is a map of object IDs to a diff of the object.
// The diff of an object is a list of changes to the object.
// The changes are a path to the changed node and the old and new values.
// The path is a dot-separated list of keys to traverse the object.
// The old and new values are the values at the path.
func Diff(oldTree, newTree Tree, getter NodeGetter) (TreeDiff, error) {
	// TODO: Implement
}
```