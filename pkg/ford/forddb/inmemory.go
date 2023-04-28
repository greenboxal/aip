package forddb

import (
	"reflect"
	"sync"
	"time"
)

type inMemoryDatabase struct {
	HasListenersBase

	m         sync.RWMutex
	resources map[ResourceTypeID]*resourceTable
}

func NewInMemory() Database {
	db := &inMemoryDatabase{
		resources: make(map[ResourceTypeID]*resourceTable),
	}

	for _, typ := range typeSystem.resourceTypes {
		if _, err := db.Put(typ); err != nil {
			panic(err)
		}
	}

	return db
}

func (db *inMemoryDatabase) List(typ ResourceTypeID) ([]BasicResource, error) {
	rt := db.getTable(typ, false)

	if rt == nil {
		return nil, nil
	}

	return rt.List()
}

func (db *inMemoryDatabase) Get(typ ResourceTypeID, id BasicResourceID) (BasicResource, error) {
	slot := db.getSlot(typ, id, false, false)

	if slot == nil {
		return nil, ErrNotFound
	}

	return slot.Get()
}

func (db *inMemoryDatabase) Put(resource BasicResource) (BasicResource, error) {
	slot := db.getSlot(resource.GetType(), resource.GetResourceID(), true, false)

	if slot == nil {
		return nil, ErrNotFound
	}

	return slot.Update(resource)
}

func (db *inMemoryDatabase) Delete(resource BasicResource) (BasicResource, error) {
	slot := db.getSlot(resource.GetType(), resource.GetResourceID(), true, false)

	if slot == nil {
		return nil, ErrNotFound
	}

	return slot.Delete()
}

func (db *inMemoryDatabase) getSlot(typ ResourceTypeID, id BasicResourceID, create, lock bool) *resourceSlot {
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

func (db *inMemoryDatabase) getTable(typ ResourceTypeID, create bool) *resourceTable {
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

		resources: make(map[BasicResourceID]*resourceSlot, 32),
	}

	db.resources[typ] = rt

	return rt
}

type resourceTable struct {
	HasListenersBase

	m         sync.RWMutex
	db        *inMemoryDatabase
	typ       ResourceTypeID
	resources map[BasicResourceID]*resourceSlot
}

func (rt *resourceTable) getSlot(id BasicResourceID, create bool, lock bool) *resourceSlot {
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

func (rt *resourceTable) List() ([]BasicResource, error) {
	rt.m.RLock()
	defer rt.m.RUnlock()

	resources := make([]BasicResource, 0, len(rt.resources))

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
	HasListenersBase

	table *resourceTable
	id    BasicResourceID

	hasValue bool

	encoded RawResource
	value   BasicResource
}

func (s *resourceSlot) Get() (BasicResource, error) {
	s.RLock()
	defer s.RUnlock()

	if !s.hasValue {
		return nil, ErrNotFound
	}

	if s.table.typ.Type().IsRuntimeOnly() {
		return s.value, nil
	}

	instance, err := Decode(s.encoded)

	if err != nil {
		return nil, err
	}

	return instance.(BasicResource), nil
}

func (s *resourceSlot) Update(resource BasicResource) (BasicResource, error) {
	old, current, changed, err := s.doUpdate(resource)

	if err != nil {
		return nil, err
	}

	if changed {
		FireListeners(&s.HasListenersBase, s.id, old, current)
		FireListeners(&s.table.HasListenersBase, s.id, old, current)
		FireListeners(&s.table.db.HasListenersBase, s.id, old, current)
	}

	return current, nil
}

func (s *resourceSlot) doUpdate(resource BasicResource) (BasicResource, BasicResource, bool, error) {
	s.RLock()
	defer s.RUnlock()

	current, err := s.Get()

	if err == ErrNotFound {
		current = nil
	} else if err != nil {
		return nil, nil, false, err
	}

	if current != nil {
		if current.GetVersion() != resource.GetVersion() {
			return nil, nil, false, ErrVersionMismatch
		}

		if current != nil {
			if reflect.DeepEqual(current, resource) {
				return current, resource, false, nil
			}
		}
	}

	meta := resource.GetMetadata()

	meta.Version += 1
	meta.UpdatedAt = time.Now()

	if !s.hasValue {
		meta.CreatedAt = meta.UpdatedAt
	}

	if s.table.typ.Type().IsRuntimeOnly() {
		s.value = resource
	} else {
		encoded, err := Encode(resource)

		if err != nil {
			return nil, nil, false, err
		}

		s.encoded = encoded
	}

	s.hasValue = true

	return current, resource, true, nil
}

func (s *resourceSlot) Delete() (BasicResource, error) {
	previous, ok, err := s.doDelete()

	if err != nil {
		return nil, err
	}

	if ok {
		FireListeners(&s.HasListenersBase, s.id, previous, nil)
		FireListeners(&s.table.HasListenersBase, s.id, previous, nil)
		FireListeners(&s.table.db.HasListenersBase, s.id, previous, nil)
	}

	return previous, nil
}

func (s *resourceSlot) doDelete() (BasicResource, bool, error) {
	value, err := s.Get()

	if err != nil {
		return nil, false, err
	}

	s.encoded = nil
	s.value = nil
	s.hasValue = false

	return value, true, nil
}
