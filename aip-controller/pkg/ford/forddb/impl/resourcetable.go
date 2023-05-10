package forddbimpl

import (
	"context"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/dgraph-io/ristretto"
	"github.com/hashicorp/go-multierror"
	"github.com/zyedidia/generic/mapset"

	forddb "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type resourceTable struct {
	forddb.HasListenersBase

	m sync.RWMutex

	db  *database
	typ forddb.TypeID

	cache               *ristretto.Cache
	persistentResources mapset.Set[forddb.BasicResourceID]
}

func newResourceTable(db *database, typ forddb.TypeID) (*resourceTable, error) {
	rt := &resourceTable{
		db:  db,
		typ: typ,

		persistentResources: mapset.New[forddb.BasicResourceID](),
	}

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.

		//KeyToHash: func(key interface{}) (uint64, uint64) {
		//	id := key.(forddb.BasicResourceID)

		//	th := id.BasicResourceType().GetResourceBasicID().Hash64()
		//	ih := id.Hash64()

		//	return th, ih
		//},

		OnExit: func(value interface{}) {
			rt.m.Lock()
			defer rt.m.Unlock()

			rs := value.(*resourceSlot)

			rt.persistentResources.Remove(rs.id)
		},

		Cost: func(value interface{}) int64 {
			rs := value.(*resourceSlot)

			return rs.getCost()
		},
	})

	if err != nil {
		return nil, err
	}

	rt.cache = cache

	return rt, nil
}

func (rt *resourceTable) ListKeys() ([]forddb.BasicResourceID, error) {
	rt.m.RLock()
	defer rt.m.RUnlock()

	resources := make([]forddb.BasicResourceID, 0, rt.persistentResources.Size())

	rt.persistentResources.Each(func(key forddb.BasicResourceID) {
		resources = append(resources, key)
	})

	return resources, nil
}

func (rt *resourceTable) List(ctx context.Context, opts forddb.ListOptions) ([]forddb.BasicResource, error) {
	var merr error

	if !rt.typ.Type().IsRuntimeOnly() {
		results, err := rt.db.storage.List(ctx, rt.typ, opts)

		if err != nil {
			return nil, errors.Wrap(err, "failed to list resources")
		}

		resources := make([]forddb.BasicResource, len(results))

		for i, v := range results {
			resource, err := forddb.Decode(v)

			if err != nil {
				return nil, err
			}

			resources[i] = resource
		}

		return resources, nil
	}

	keys, err := rt.ListKeys()

	if err != nil {
		return nil, err
	}

	resources := make([]forddb.BasicResource, 0, len(keys))

	for _, key := range keys {
		rs := rt.GetSlot(key, true)
		resource, err := rs.Get(ctx)

		if forddb.IsNotFound(err) {
			merr = multierror.Append(merr, err)
			continue
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

func (rt *resourceTable) GetSlot(id forddb.BasicResourceID, create bool) *resourceSlot {
	isNew := false

	defer func() {
		if isNew {
			rt.cache.Wait()
		}
	}()

	if create {
		rt.m.Lock()
		defer rt.m.Unlock()
	} else {
		rt.m.RLock()
		defer rt.m.RUnlock()
	}

	existing, hasExisting := rt.cache.Get(id.String())

	if hasExisting {
		return existing.(*resourceSlot)
	}

	if !create {
		return nil
	}

	rs := newResourceSlot(rt, id)

	rt.cache.Set(id.String(), rs, 0)
	rt.persistentResources.Put(id)

	isNew = true

	return rs
}

func (rt *resourceTable) notifyGet(rs *resourceSlot) {
	rt.db.notifyGet(rs)
}
