package forddbimpl

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/greenboxal/aip/pkg/ford/forddb"
	"github.com/greenboxal/aip/pkg/ford/forddb/logstore"
)

type resourceSlot struct {
	m    sync.RWMutex
	cond *sync.Cond

	table *resourceTable
	id    forddb.BasicResourceID

	isPinned bool
	hasValue bool

	lastRecord logstore.LogEntryRecord
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

	if raw == nil {
		return nil, forddb.ErrNotFound
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

	rs.table.notifyGet(rs)

	for wait && !rs.hasValue && rs.err == nil {
		rs.cond.Wait()
	}

	if rs.err != nil {
		return nil, nil, rs.err
	}

	if !rs.hasValue {
		return nil, nil, forddb.ErrNotFound
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

func (rs *resourceSlot) Update(
	ctx context.Context,
	resource forddb.BasicResource,
	opts forddb.PutOptions,
) (forddb.BasicResource, error) {
	_, current, changed, err := rs.doUpdate(ctx, resource, opts)

	if err != nil {
		return nil, err
	}

	if changed {
		rs.cond.Broadcast()
	}

	return current, nil
}

func (rs *resourceSlot) doUpdate(
	ctx context.Context,
	resource forddb.BasicResource,
	opts forddb.PutOptions,
) (forddb.BasicResource, forddb.BasicResource, bool, error) {
	rs.cond.L.Lock()
	defer rs.cond.L.Unlock()

	_, current, err := rs.doGet(ctx, false, false, true)

	if err == forddb.ErrNotFound {
		current = nil
	} else if err != nil {
		return nil, nil, false, err
	}

	if current != nil {
		if reflect.DeepEqual(current, resource) {
			return current, resource, false, nil
		}

		switch opts.OnConflict {
		case forddb.OnConflictError:
			return current, nil, false, forddb.ErrVersionMismatch
		case forddb.OnConflictOptimistic:
			if current.GetVersion() != resource.GetVersion() {
				return current, nil, false, forddb.ErrVersionMismatch
			}
		case forddb.OnConflictLatestWins:
			if current.GetVersion() >= resource.GetVersion() {
				return current, current, false, nil
			}
		case forddb.OnConflictReplace:
			// Nothing
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

		record, err := rs.table.db.log.Append(ctx, logstore.LogEntry{
			Kind:           logstore.LogEntryKindSet,
			Type:           rs.table.typ,
			ID:             rs.id,
			Version:        meta.Version,
			CurrentCid:     nil,
			PreviousCid:    nil,
			Previous:       rs.encoded,
			Current:        encoded,
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
	rs.err = nil

	return current, resource, true, nil
}

func (rs *resourceSlot) Delete(ctx context.Context) (forddb.BasicResource, error) {
	previous, _, err := rs.doDelete(ctx)

	if err != nil {
		return nil, err
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
		_, err := rs.table.db.log.Append(ctx, logstore.LogEntry{
			Kind:           logstore.LogEntryKindDelete,
			Type:           rs.table.typ,
			ID:             rs.id,
			Version:        value.GetVersion(),
			CurrentCid:     nil,
			PreviousCid:    nil,
			Previous:       raw,
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

func (rs *resourceSlot) setError(err error) {
	rs.m.Lock()
	defer rs.m.Unlock()

	rs.err = err
	rs.cond.Broadcast()
}
