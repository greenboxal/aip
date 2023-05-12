package forddbimpl

import (
	"context"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type resourceNode struct {
	slot *resourceSlot
}

func (r *resourceNode) Subscribe(listener forddb.Listener) func() {
	return r.slot.Subscribe(listener)
}

func (r *resourceNode) ResourceID() forddb.BasicResourceID {
	return r.slot.id
}

func (r *resourceNode) ResourceType() forddb.TypeID {
	return r.slot.table.typ
}

func newResourceNode(slot *resourceSlot) *resourceNode {
	return &resourceNode{
		slot: slot,
	}
}

func (r *resourceNode) Get(ctx context.Context) (forddb.BasicResource, error) {
	return r.slot.Get(ctx)
}
