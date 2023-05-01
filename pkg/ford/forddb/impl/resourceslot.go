package forddbimpl

import (
	"context"
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

	lastRecord forddb.LogEntryRecord
	encoded    forddb.RawResource
	value      forddb.BasicResource
	err        error
}

func newResourceSlot(table *resourceTable, id forddb.BasicResourceID) *resourceSlot {
	rs := &resourceSlot{
		table: table,
		id:    id,
	}

	rs.cond = sync.NewCond(&rs.m)

	return rs
}

func (rs *resourceSlot) Get(ctx context.Context) (forddb.BasicResource, error) {
	raw, res, err := rs.doGet(ctx, true, true, false)

	if err != nil {
		return nil, err
	}

	if res != nil {
		return res, nil
	}

	return forddb.Decode(raw)
}

func (rs *resourceSlot) doGet(
	ctx context.Context,
	lock bool,
	wait bool,
	decode bool,
) (raw forddb.RawResource, res forddb.BasicResource, err error) {
	if lock {
		rs.m.RLock()
		defer rs.m.RUnlock()
	}

	if !rs.hasValue && wait {
		if !rs.hasValue {
			rs.cond.Wait()
		}
	}

	if !rs.hasValue {
		return nil, nil, forddb.ErrNotFound
	}

	if rs.err != nil {
		return nil, nil, rs.err
	}

	if rs.table.typ.Type().IsRuntimeOnly() {
		res = rs.value
	} else {
		raw = rs.encoded

		if decode {
			res, err = forddb.Decode(rs.encoded)

			if err != nil {
				return
			}
		}
	}

	return
}

func (rs *resourceSlot) Update(ctx context.Context, resource forddb.BasicResource) (forddb.BasicResource, error) {
	old, current, changed, err := rs.doUpdate(ctx, resource)

	if err != nil {
		return nil, err
	}

	if changed {
		rs.cond.Broadcast()

		forddb.FireListeners(&rs.HasListenersBase, rs.id, old, current)
		forddb.FireListeners(&rs.table.HasListenersBase, rs.id, old, current)
		forddb.FireListeners(&rs.table.db.HasListenersBase, rs.id, old, current)
	}

	return current, nil
}

func (rs *resourceSlot) doUpdate(ctx context.Context, resource forddb.BasicResource) (forddb.BasicResource, forddb.BasicResource, bool, error) {
	rs.cond.L.Lock()
	defer rs.cond.L.Unlock()

	_, current, err := rs.doGet(ctx, false, false, true)

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

	if !rs.hasValue {
		meta.CreatedAt = meta.UpdatedAt
	}

	if rs.table.typ.Type().IsRuntimeOnly() {
		rs.value = resource
	} else {
		encoded, err := forddb.Encode(resource)

		if err != nil {
			return nil, nil, false, err
		}

		previous := rs.encoded

		record, err := rs.table.db.log.Append(ctx, forddb.LogEntry{
			Kind:           forddb.LogEntryKindSet,
			Type:           rs.table.typ,
			ID:             rs.id,
			Version:        meta.Version,
			CurrentCid:     nil,
			PreviousCid:    nil,
			Previous:       &previous,
			Current:        &encoded,
			CachedPrevious: current,
			CachedCurrent:  resource,
		})

		if err != nil {
			return nil, nil, false, err
		}

		rs.lastRecord = record
		rs.encoded = encoded
	}

	rs.hasValue = true

	return current, resource, true, nil
}

func (rs *resourceSlot) Delete(ctx context.Context) (forddb.BasicResource, error) {
	previous, ok, err := rs.doDelete(ctx)

	if err != nil {
		return nil, err
	}

	if ok {
		forddb.FireListeners(&rs.HasListenersBase, rs.id, previous, nil)
		forddb.FireListeners(&rs.table.HasListenersBase, rs.id, previous, nil)
		forddb.FireListeners(&rs.table.db.HasListenersBase, rs.id, previous, nil)
	}

	return previous, nil
}

func (rs *resourceSlot) doDelete(ctx context.Context) (forddb.BasicResource, bool, error) {
	rs.cond.L.Lock()
	defer rs.cond.L.Unlock()

	raw, value, err := rs.doGet(ctx, false, false, true)

	if err != nil {
		return nil, false, err
	}

	if !rs.table.typ.Type().IsRuntimeOnly() {
		_, err := rs.table.db.log.Append(ctx, forddb.LogEntry{
			Kind:           forddb.LogEntryKindDelete,
			Type:           rs.table.typ,
			ID:             rs.id,
			Version:        value.GetVersion(),
			CurrentCid:     nil,
			PreviousCid:    nil,
			Previous:       &raw,
			Current:        nil,
			CachedPrevious: value,
			CachedCurrent:  nil,
		})

		if err != nil {
			return nil, false, err
		}
	}

	rs.encoded = nil
	rs.value = nil
	rs.hasValue = false

	return value, true, nil
}

func (rs *resourceSlot) getCost() int64 {
	cost := int64(1)

	if rs.isPinned {
		cost = -cost
	}

	return cost
}
