package forddbimpl

import (
	"reflect"
	"sync"
	"time"

	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type resourceSlot struct {
	forddb.HasListenersBase

	m    sync.RWMutex
	cond *sync.Cond

	table *resourceTable
	id    forddb.BasicResourceID

	isPinned bool
	hasValue bool

	encoded forddb.RawResource
	value   forddb.BasicResource
	err     error
}

func newResourceSlot(table *resourceTable, id forddb.BasicResourceID) *resourceSlot {
	rs := &resourceSlot{
		table: table,
		id:    id,
	}

	rs.cond = sync.NewCond(&rs.m)

	return rs
}

func (s *resourceSlot) Get() (forddb.BasicResource, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.doGet(true)
}

func (s *resourceSlot) doGet(wait bool) (forddb.BasicResource, error) {
	if !s.hasValue && wait {
		if !s.hasValue {
			s.cond.Wait()
		}
	}

	if !s.hasValue {
		return nil, forddb.ErrNotFound
	}

	if s.err != nil {
		return nil, s.err
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
		s.cond.Broadcast()

		forddb.FireListeners(&s.HasListenersBase, s.id, old, current)
		forddb.FireListeners(&s.table.HasListenersBase, s.id, old, current)
		forddb.FireListeners(&s.table.db.HasListenersBase, s.id, old, current)
	}

	return current, nil
}

func (s *resourceSlot) doUpdate(resource forddb.BasicResource) (forddb.BasicResource, forddb.BasicResource, bool, error) {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	current, err := s.doGet(false)

	if err == forddb.ErrNotFound {
		current = nil
	} else if err != nil {
		return nil, nil, false, err
	}

	if current != nil {
		if current.GetVersion() != resource.GetVersion() {
			return nil, nil, false, forddb.ErrVersionMismatch
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
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	value, err := s.doGet(false)

	if err != nil {
		return nil, false, err
	}

	s.encoded = nil
	s.value = nil
	s.hasValue = false

	return value, true, nil
}

func (s *resourceSlot) getCost() int64 {
	cost := int64(1)

	if s.isPinned {
		cost = -cost
	}

	return cost
}
