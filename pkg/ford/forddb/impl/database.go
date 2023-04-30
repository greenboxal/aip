package forddbimpl

import (
	"context"
	"sync"

	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type database struct {
	forddb.HasListenersBase

	m         sync.RWMutex
	resources map[forddb.ResourceTypeID]*resourceTable
}

func NewDatabase(storage forddb.Storage) forddb.Database {
	db := &database{
		resources: make(map[forddb.ResourceTypeID]*resourceTable),
	}

	// FIXME: Shouldn't happen here
	for _, typ := range forddb.TypeSystem().ResourceTypes() {
		if _, err := db.Put(context.Background(), typ); err != nil {
			panic(err)
		}
	}

	return db
}

func (db *database) List(ctx context.Context, typ forddb.ResourceTypeID) ([]forddb.BasicResource, error) {
	rt := db.getTable(typ, false)

	if rt == nil {
		return nil, nil
	}

	return rt.List()
}

func (db *database) Get(ctx context.Context, typ forddb.ResourceTypeID, id forddb.BasicResourceID) (forddb.BasicResource, error) {
	slot := db.getSlot(typ, id, false)

	if slot == nil {
		return nil, forddb.ErrNotFound
	}

	return slot.Get()
}

func (db *database) Put(ctx context.Context, resource forddb.BasicResource) (forddb.BasicResource, error) {
	slot := db.getSlot(resource.GetType(), resource.GetResourceID(), true)

	if slot == nil {
		return nil, forddb.ErrNotFound
	}

	return slot.Update(resource)
}

func (db *database) Delete(ctx context.Context, resource forddb.BasicResource) (forddb.BasicResource, error) {
	slot := db.getSlot(resource.GetType(), resource.GetResourceID(), true)

	if slot == nil {
		return nil, forddb.ErrNotFound
	}

	return slot.Delete()
}

func (db *database) getSlot(typ forddb.ResourceTypeID, id forddb.BasicResourceID, create bool) *resourceSlot {
	tb := db.getTable(typ, create)

	if tb == nil {
		return nil
	}

	slot := tb.getSlot(id, create)

	if slot == nil {
		return nil
	}

	return slot
}

func (db *database) getTable(typ forddb.ResourceTypeID, create bool) *resourceTable {
	db.m.Lock()
	defer db.m.Unlock()

	if existing := db.resources[typ]; existing != nil {
		return existing
	}

	if !create {
		return nil
	}

	rt, err := newResourceTable(db, typ)

	if err != nil {
		panic(err)
	}

	db.resources[typ] = rt

	return rt
}
