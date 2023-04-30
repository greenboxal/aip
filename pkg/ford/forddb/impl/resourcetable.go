package forddbimpl

import (
	"sync"

	"github.com/dgraph-io/ristretto"
	"github.com/hashicorp/go-multierror"
	"github.com/zyedidia/generic/mapset"

	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type resourceTable struct {
	forddb.HasListenersBase

	m sync.RWMutex

	db  *database
	typ forddb.ResourceTypeID

	cache               *ristretto.Cache
	persistentResources mapset.Set[forddb.BasicResourceID]
}

func newResourceTable(db *database, typ forddb.ResourceTypeID) (*resourceTable, error) {
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

func (rt *resourceTable) List() ([]forddb.BasicResource, error) {
	var merr error

	rt.m.RLock()
	defer rt.m.RUnlock()

	resources := make([]forddb.BasicResource, 0, rt.persistentResources.Size())

	rt.persistentResources.Each(func(key forddb.BasicResourceID) {
		v := rt.getSlot(key, false)

		if v == nil {
			return
		}

		value, err := v.doGet(false)

		if err == forddb.ErrNotFound {
			return
		} else if err != nil {
			merr = multierror.Append(merr, err)
			return
		}

		resources = append(resources, value)
	})

	return resources, nil
}

func (rt *resourceTable) getSlot(id forddb.BasicResourceID, create bool) *resourceSlot {
	isNew := false

	defer func() {
		if isNew {
			rt.cache.Wait()
		}
	}()

	rt.m.Lock()
	defer rt.m.Unlock()

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