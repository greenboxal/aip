package forddbimpl

import (
	"context"
	"sync"

	"github.com/jbenet/goprocess"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-forddb/pkg/objectstore"
)

type database struct {
	forddb.HasListenersBase

	m       sync.RWMutex
	log     forddb.LogStore
	storage objectstore.ObjectStore

	resources map[forddb.TypeID]*resourceTable

	objectFetcher      *objectFetcher
	objectFetchProcess goprocess.Process

	eventDispatcher        *eventDispatcher
	eventDispatcherProcess goprocess.Process
}

func NewDatabase(
	logStore forddb.LogStore,
	storage objectstore.ObjectStore,
) forddb.Database {
	db := &database{
		log:     logStore,
		storage: storage,

		resources: make(map[forddb.TypeID]*resourceTable),
	}

	db.objectFetcher = newObjectFetcher(db)
	db.objectFetchProcess = goprocess.Go(db.objectFetcher.Run)

	db.eventDispatcher = newEventDispatcher(db)
	db.eventDispatcherProcess = goprocess.Go(db.eventDispatcher.Run)

	return db
}

func (db *database) LogStore() forddb.LogStore {
	return db.log
}

func (db *database) List(
	ctx context.Context,
	typ forddb.TypeID,
	options ...forddb.ListOption,
) ([]forddb.BasicResource, error) {
	opts := forddb.NewListOptions(typ, options...)

	rt := db.GetTable(typ, true)

	if rt == nil {
		return nil, nil
	}

	return rt.List(ctx, opts)
}

func (db *database) Get(
	ctx context.Context,
	typ forddb.TypeID,
	id forddb.BasicResourceID,
	options ...forddb.GetOption,
) (forddb.BasicResource, error) {
	slot := db.GetSlot(typ, id, true)

	if slot == nil {
		return nil, forddb.ErrNotFound
	}

	return slot.Get(ctx, options...)
}

func (db *database) Put(
	ctx context.Context,
	resource forddb.BasicResource,
	options ...forddb.PutOption,
) (forddb.BasicResource, error) {
	resource.OnBeforeSave(resource)

	raw, err := forddb.Encode(resource)

	if err != nil {
		return nil, err
	}

	opts := forddb.NewPutOptions(options...)
	slot := db.GetSlot(resource.GetResourceTypeID(), resource.GetResourceBasicID(), true)

	result, err := slot.Update(ctx, raw, opts)

	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return forddb.Decode(result)
}

func (db *database) Delete(
	ctx context.Context,
	resource forddb.BasicResource,
	options ...forddb.DeleteOption,
) (forddb.BasicResource, error) {
	slot := db.GetSlot(resource.GetResourceTypeID(), resource.GetResourceBasicID(), true)

	if slot == nil {
		return nil, forddb.ErrNotFound
	}

	result, err := slot.Delete(ctx)

	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return forddb.Decode(result)
}

func (db *database) GetSlot(typ forddb.TypeID, id forddb.BasicResourceID, create bool) *resourceSlot {
	tb := db.GetTable(typ, create)

	if tb == nil {
		return nil
	}

	slot := tb.GetSlot(id, create)

	if slot == nil {
		return nil
	}

	return slot
}

func (db *database) GetTable(typ forddb.TypeID, create bool) *resourceTable {
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

func (db *database) Close() error {
	db.m.Lock()
	defer db.m.Unlock()

	if err := db.storage.Close(); err != nil {
		return err
	}

	if err := db.log.Close(); err != nil {
		return err
	}

	return nil
}

func (db *database) notifyGet(rs *resourceSlot) {
	//db.objectFetcher.requestCh <- fetchResourceRequest{
	//	storage: db.storage,
	//	slot:    rs,
	//}
}
