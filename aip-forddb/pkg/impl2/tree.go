package forddbimpl

import (
	"github.com/ipfs/go-unixfsnode/directory"
	"github.com/ipfs/go-unixfsnode/iter"
)

type MutableTree struct {
	entries map[string]TreeEntry
}

type ImmutableTree struct {
	substrate directory.UnixFSBasicDir
}

type treeIterator struct {
	it *iter.UnixFSDir__Itr
}

func (t *treeIterator) Next() bool {
	return !t.it.Done()
}

func (t *treeIterator) Entry() TreeEntry {
	k, v := t.it.Next()

	return TreeEntry{
		Name: k.String(),
		Link: v.Link(),
	}
}

func (t *ImmutableTree) Iterator() TreeIterator {
	return &treeIterator{
		it: t.substrate.Iterator(),
	}
}
