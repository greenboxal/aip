package forddbimpl

import (
	"context"
	"sync"

	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"

	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type database struct {
	forddb.HasListenersBase

	m       sync.RWMutex
	log     forddb.LogStore
	storage forddb.Storage

	resources map[forddb.ResourceTypeID]*resourceTable

	objectIndexer        *objectIndexer
	objectIndexerProcess goprocess.Process

	objectFetcher      *objectFetcher
	objectFetchProcess goprocess.Process
}

func NewDatabase(logStore forddb.LogStore, storage forddb.Storage) forddb.Database {
	db := &database{
		log:     logStore,
		storage: storage,

		resources: make(map[forddb.ResourceTypeID]*resourceTable),
	}

	db.objectIndexer = newObjectIndexer(db)
	db.objectIndexerProcess = goprocess.Go(func(proc goprocess.Process) {
		ctx := goprocessctx.OnClosingContext(proc)

		if err := db.objectIndexer.Run(ctx); err != nil {
			panic(err)
		}
	})

	db.objectFetcher = newObjectFetcher(db)
	db.objectFetchProcess = goprocess.Go(func(proc goprocess.Process) {
		ctx := goprocessctx.OnClosingContext(proc)

		if err := db.objectIndexer.Run(ctx); err != nil {
			panic(err)
		}
	})

	// Index all resource types
	for _, typ := range forddb.TypeSystem().ResourceTypes() {
		if _, err := db.Put(context.Background(), typ); err != nil {
			panic(err)
		}
	}

	return db
}

func (db *database) List(ctx context.Context, typ forddb.ResourceTypeID) ([]forddb.BasicResource, error) {
	rt := db.GetTable(typ, false)

	if rt == nil {
		return nil, nil
	}

	return rt.List(ctx)
}

func (db *database) Get(ctx context.Context, typ forddb.ResourceTypeID, id forddb.BasicResourceID) (forddb.BasicResource, error) {
	slot := db.GetSlot(typ, id, false)

	if slot == nil {
		return nil, forddb.ErrNotFound
	}

	return slot.Get(ctx)
}

func (db *database) Put(ctx context.Context, resource forddb.BasicResource, options ...forddb.PutOption) (forddb.BasicResource, error) {
	opts := forddb.NewPutOptions(options...)
	slot := db.GetSlot(resource.GetType(), resource.GetResourceID(), true)

	if slot == nil {
		return nil, forddb.ErrNotFound
	}

	return slot.Update(ctx, resource, opts)
}

func (db *database) Delete(ctx context.Context, resource forddb.BasicResource) (forddb.BasicResource, error) {
	slot := db.GetSlot(resource.GetType(), resource.GetResourceID(), true)

	if slot == nil {
		return nil, forddb.ErrNotFound
	}

	return slot.Delete(ctx)
}

func (db *database) GetSlot(typ forddb.ResourceTypeID, id forddb.BasicResourceID, create bool) *resourceSlot {
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

func (db *database) GetTable(typ forddb.ResourceTypeID, create bool) *resourceTable {
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

	if err := db.objectIndexerProcess.Close(); err != nil {
		return err
	}

	if err := db.storage.Close(); err != nil {
		return err
	}

	if err := db.log.Close(); err != nil {
		return err
	}

	return nil
}

func (db *database) notifyGet(rs *resourceSlot) {
	db.objectFetcher.requestCh <- fetchResourceRequest{
		storage: db.storage,
		slot:    rs,
	}
}
