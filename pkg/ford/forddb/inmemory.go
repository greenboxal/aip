package forddb

import (
	"encoding/json"
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
		if _, err := db.UpdateOrCreate(typ); err != nil {
			panic(err)
		}
	}

	return db
}

func (i *inMemoryDatabase) List(typ ResourceTypeID) ([]BasicResource, error) {
	rt := i.getTable(typ, false)

	if rt == nil {
		return nil, nil
	}

	return rt.List()
}

func (i *inMemoryDatabase) Get(typ ResourceTypeID, id BasicResourceID) (BasicResource, error) {
	slot := i.getSlot(typ, id, false, false)

	if slot == nil {
		return nil, nil
	}

	return slot.Get()
}

func (i *inMemoryDatabase) UpdateOrCreate(resource BasicResource) (BasicResource, error) {
	slot := i.getSlot(resource.GetType(), resource.GetResourceID(), true, false)

	if slot == nil {
		return nil, nil
	}

	return slot.Update(resource)
}

func (i *inMemoryDatabase) Delete(resource BasicResource) (BasicResource, error) {
	slot := i.getSlot(resource.GetType(), resource.GetResourceID(), true, false)

	if slot == nil {
		return nil, nil
	}

	return slot.Delete()
}

func (i *inMemoryDatabase) getSlot(typ ResourceTypeID, id BasicResourceID, create, lock bool) *resourceSlot {
	tb := i.getTable(typ, create)

	if tb == nil {
		return nil
	}

	slot := tb.getSlot(id, create, lock)

	if slot == nil {
		return nil
	}

	return slot
}

func (i *inMemoryDatabase) getTable(typ ResourceTypeID, create bool) *resourceTable {
	i.m.Lock()
	defer i.m.Unlock()

	if existing := i.resources[typ]; existing != nil {
		return existing
	}

	if !create {
		return nil
	}

	rt := &resourceTable{
		typ: typ,

		resources: make(map[BasicResourceID]*resourceSlot, 32),
	}

	i.resources[typ] = rt

	return rt
}

type resourceTable struct {
	HasListenersBase

	m         sync.RWMutex
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

	hasValue bool

	version    int
	serialized []byte
	value      BasicResource
}

func (s *resourceSlot) Get() (BasicResource, error) {
	s.RLock()
	defer s.RUnlock()

	if !s.hasValue {
		return nil, nil
	}

	if s.table.typ.Type().IsRuntimeOnly() {
		return s.value, nil
	}

	instance := reflect.New(s.table.typ.Type().ResourceType()).Interface()

	if err := json.Unmarshal(s.serialized, instance); err != nil {
		return nil, err
	}

	return instance.(BasicResource), nil
}

func (s *resourceSlot) Update(resource BasicResource) (BasicResource, error) {
	s.RLock()
	defer s.RUnlock()

	if s.hasValue && s.version != resource.GetVersion() {
		return nil, ErrVersionMismatch
	}

	current, err := s.Get()

	if err != nil {
		return nil, err
	}

	if current != nil {
		if reflect.DeepEqual(current, resource) {
			return resource, nil
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
		data, err := json.Marshal(resource)

		if err != nil {
			return nil, err
		}

		s.serialized = data
	}

	s.version = resource.GetVersion()
	s.hasValue = true

	return resource, nil
}

func (s *resourceSlot) Delete() (BasicResource, error) {
	value, err := s.Get()

	if err != nil {
		return nil, err
	}

	s.serialized = nil
	s.value = nil
	s.hasValue = false

	return value, nil
}
