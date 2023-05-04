package forddbimpl

import (
	"context"
	"sync"

	"github.com/dgraph-io/ristretto"
	"github.com/hashicorp/go-multierror"
	"github.com/zyedidia/generic/mapset"

	forddb2 "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type resourceTable struct {
	m sync.RWMutex

	db  *database
	typ forddb2.ResourceTypeID

	cache               *ristretto.Cache
	persistentResources mapset.Set[forddb2.BasicResourceID]
}

func newResourceTable(db *database, typ forddb2.ResourceTypeID) (*resourceTable, error) {
	rt := &resourceTable{
		db:  db,
		typ: typ,

		persistentResources: mapset.New[forddb2.BasicResourceID](),
	}

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.

		//KeyToHash: func(key interface{}) (uint64, uint64) {
		//	id := key.(forddb.BasicResourceID)

		//	th := id.BasicResourceType().GetResourceID().Hash64()
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

func (rt *resourceTable) ListKeys() ([]forddb2.BasicResourceID, error) {
	rt.m.RLock()
	defer rt.m.RUnlock()

	resources := make([]forddb2.BasicResourceID, 0, rt.persistentResources.Size())

	rt.persistentResources.Each(func(key forddb2.BasicResourceID) {
		resources = append(resources, key)
	})

	return resources, nil
}

func (rt *resourceTable) List(ctx context.Context) ([]forddb2.BasicResource, error) {
	var merr error

	keys, err := rt.ListKeys()

	if err != nil {
		return nil, err
	}

	resources := make([]forddb2.BasicResource, 0, len(keys))

	for _, key := range keys {
		v := rt.GetSlot(key, false)

		if v == nil {
			continue
		}

		if !v.hasValue {
			continue
		}

		_, value, err := v.doGet(ctx, false, false, true)

		if err == forddb2.ErrNotFound {
			continue
		} else if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		resources = append(resources, value)
	}

	return resources, nil
}

func (rt *resourceTable) GetSlot(id forddb2.BasicResourceID, create bool) *resourceSlot {
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
