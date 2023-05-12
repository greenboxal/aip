package inmemory

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/greenboxal/aip/aip-controller/pkg/collective"
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type InMemoryDatabase struct {
	forddb.HasListenersBase

	m         sync.RWMutex
	resources map[forddb.TypeID]*resourceTable
}

func (db *InMemoryDatabase) AppendSegment(ctx context.Context, segment *collective.MemorySegment) error {
	//TODO implement me
	panic("implement me")
}

func (db *InMemoryDatabase) GetSegment(ctx context.Context, id collective.MemorySegmentID) (*collective.MemorySegment, error) {
	//TODO implement me
	panic("implement me")
}

func (db *InMemoryDatabase) GetMemory(ctx context.Context, id collective.MemoryID) (*collective.Memory, error) {
	//TODO implement me
	panic("implement me")
}

func (db *InMemoryDatabase) Close() error {
	return nil
}

func NewInMemory() *InMemoryDatabase {
	db := &InMemoryDatabase{
		resources: make(map[forddb.TypeID]*resourceTable),
	}

	return db
}

func (db *InMemoryDatabase) List(
	ctx context.Context,
	typ forddb.TypeID,
	opts forddb.ListOptions,
) ([]forddb.BasicResource, error) {
	rt := db.getTable(typ, false)

	if rt == nil {
		return nil, nil
	}

	return rt.List()
}

func (db *InMemoryDatabase) Get(
	ctx context.Context,
	typ forddb.TypeID,
	id forddb.BasicResourceID,
	opts forddb.GetOptions,
) (forddb.BasicResource, error) {
	slot := db.getSlot(typ, id, false, false)

	if slot == nil {
		return nil, forddb.ErrNotFound
	}

	return slot.Get()
}

func (db *InMemoryDatabase) Put(
	ctx context.Context,
	resource forddb.BasicResource,
	opts forddb.PutOptions,
) (forddb.BasicResource, error) {
	slot := db.getSlot(resource.GetResourceTypeID(), resource.GetResourceBasicID(), true, false)

	if slot == nil {
		return nil, forddb.ErrNotFound
	}

	return slot.Update(resource)
}

func (db *InMemoryDatabase) Delete(
	ctx context.Context,
	resource forddb.BasicResource,
	opts forddb.DeleteOptions,
) (forddb.BasicResource, error) {
	slot := db.getSlot(resource.GetResourceTypeID(), resource.GetResourceBasicID(), true, false)

	if slot == nil {
		return nil, forddb.ErrNotFound
	}

	return slot.Delete()
}

func (db *InMemoryDatabase) getSlot(typ forddb.TypeID, id forddb.BasicResourceID, create, lock bool) *resourceSlot {
	tb := db.getTable(typ, create)

	if tb == nil {
		return nil
	}

	slot := tb.getSlot(id, create, lock)

	if slot == nil {
		return nil
	}

	return slot
}

func (db *InMemoryDatabase) getTable(typ forddb.TypeID, create bool) *resourceTable {
	db.m.Lock()
	defer db.m.Unlock()

	if existing := db.resources[typ]; existing != nil {
		return existing
	}

	if !create {
		return nil
	}

	rt := &resourceTable{
		db:  db,
		typ: typ,

		resources: make(map[forddb.BasicResourceID]*resourceSlot, 32),
	}

	db.resources[typ] = rt

	return rt
}

type resourceTable struct {
	forddb.HasListenersBase

	m         sync.RWMutex
	db        *InMemoryDatabase
	typ       forddb.TypeID
	resources map[forddb.BasicResourceID]*resourceSlot
}

func (rt *resourceTable) getSlot(id forddb.BasicResourceID, create bool, lock bool) *resourceSlot {
	rt.m.Lock()
	defer rt.m.Unlock()

	if existing := rt.resources[id]; existing != nil {
		return existing
	}

	if !create {
		return nil
	}

	rs := &resourceSlot{
		table: rt,
		id:    id,
	}

	rt.resources[id] = rs

	if lock {
		rs.Lock()
	}

	return rs
}

func (rt *resourceTable) List() ([]forddb.BasicResource, error) {
	rt.m.RLock()
	defer rt.m.RUnlock()

	resources := make([]forddb.BasicResource, 0, len(rt.resources))

	for _, v := range rt.resources {
		if !v.hasValue {
			continue
		}

		value, err := v.Get()

		if err != nil {
			return nil, err
		}

		resources = append(resources, value)
	}

	return resources, nil
}

type resourceSlot struct {
	sync.RWMutex
	forddb.HasListenersBase

	table *resourceTable
	id    forddb.BasicResourceID

	hasValue bool

	encoded forddb.RawResource
	value   forddb.BasicResource
}

func (s *resourceSlot) Get() (forddb.BasicResource, error) {
	s.RLock()
	defer s.RUnlock()

	if !s.hasValue {
		return nil, forddb.ErrNotFound
	}

	if s.table.typ.Type().IsRuntimeOnly() {
		return s.value, nil
	}

	instance, err := forddb.Decode(s.encoded)

	if err != nil {
		return nil, err
	}

	return instance.(forddb.BasicResource), nil
}

func (s *resourceSlot) Update(resource forddb.BasicResource) (forddb.BasicResource, error) {
	old, current, changed, err := s.doUpdate(resource)

	if err != nil {
		return nil, err
	}

	if changed {
		forddb.FireListeners(&s.HasListenersBase, s.id, old, current)
		forddb.FireListeners(&s.table.HasListenersBase, s.id, old, current)
		forddb.FireListeners(&s.table.db.HasListenersBase, s.id, old, current)
	}

	return current, nil
}

func (s *resourceSlot) doUpdate(resource forddb.BasicResource) (forddb.BasicResource, forddb.BasicResource, bool, error) {
	s.RLock()
	defer s.RUnlock()

	current, err := s.Get()

	if forddb.IsNotFound(err) {
		current = nil
	} else if err != nil {
		return nil, nil, false, err
	}

	if current != nil {
		if current.GetResourceVersion() != resource.GetResourceVersion() {
			return nil, nil, false, forddb.ErrVersionMismatch
		}

		if current != nil {
			if reflect.DeepEqual(current, resource) {
				return current, resource, false, nil
			}
		}
	}

	meta := resource.GetResourceMetadata()

	meta.Version += 1
	meta.UpdatedAt = time.Now()

	if !s.hasValue {
		meta.CreatedAt = meta.UpdatedAt
	}

	if s.table.typ.Type().IsRuntimeOnly() {
		s.value = resource
	} else {
		encoded, err := forddb.Encode(resource)

		if err != nil {
			return nil, nil, false, err
		}

		s.encoded = encoded
	}

	s.hasValue = true

	return current, resource, true, nil
}

func (s *resourceSlot) Delete() (forddb.BasicResource, error) {
	previous, ok, err := s.doDelete()

	if err != nil {
		return nil, err
	}

	if ok {
		forddb.FireListeners(&s.HasListenersBase, s.id, previous, nil)
		forddb.FireListeners(&s.table.HasListenersBase, s.id, previous, nil)
		forddb.FireListeners(&s.table.db.HasListenersBase, s.id, previous, nil)
	}

	return previous, nil
}

func (s *resourceSlot) doDelete() (forddb.BasicResource, bool, error) {
	value, err := s.Get()

	if err != nil {
		return nil, false, err
	}

	s.encoded = nil
	s.value = nil
	s.hasValue = false

	return value, true, nil
}
