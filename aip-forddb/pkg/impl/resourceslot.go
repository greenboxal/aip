package forddbimpl

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type resourceSlot struct {
	forddb.HasListenersBase

	m    sync.RWMutex
	cond *sync.Cond

	table *resourceTable
	node  *resourceNode
	id    forddb.BasicResourceID

	isPinned bool
	hasValue bool

	lastRecord forddb.LogEntryRecord
	encoded    forddb.RawResource
	err        error
}

func newResourceSlot(table *resourceTable, id forddb.BasicResourceID) *resourceSlot {
	rs := &resourceSlot{
		table: table,
		id:    id,
	}

	rs.cond = sync.NewCond(&rs.m)
	rs.node = newResourceNode(rs)

	return rs
}

func (rs *resourceSlot) Get(ctx context.Context, options ...forddb.GetOption) (forddb.BasicResource, error) {
	raw, err := rs.doGet(ctx, true, true)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get resource")
	}

	if raw == nil {
		return nil, errors.Wrap(forddb.ErrNotFound, "failed to get resource")
	}

	result, err := forddb.Decode(raw)

	if err != nil {
		return nil, err
	}

	forddb.GetOrCreateResourceNode(result, func(resource forddb.BasicResource) forddb.ResourceNode {
		return rs.node
	})

	return result, nil
}

func (rs *resourceSlot) doGet(
	ctx context.Context,
	lock bool,
	wait bool,
) (raw forddb.RawResource, err error) {
	res, err := rs.table.db.storage.Get(ctx, rs.table.typ, rs.id, forddb.GetOptions{})

	if err != nil && !forddb.IsNotFound(err) {
		return nil, err
	}

	if err == nil {
		if lock {
			rs.m.Lock()
			defer rs.m.Unlock()

			lock = false
		}

		if rs.lastRecord.Version < res.GetResourceVersion() {
			rs.hasValue = true
			rs.encoded = res

			return rs.encoded, err
		}
	}

	wait = false

	if lock {
		rs.m.Lock()
		defer rs.m.Unlock()
	}

	if wait && !rs.hasValue && rs.err == nil {
		for !rs.hasValue && rs.err == nil {
			rs.table.notifyGet(rs)

			rs.cond.Wait()
		}
	}

	if rs.err != nil {
		return nil, rs.err
	}

	if !rs.hasValue {
		return nil, forddb.ErrNotFound
	}

	raw = rs.encoded

	return
}

func (rs *resourceSlot) Update(
	ctx context.Context,
	resource forddb.RawResource,
	opts forddb.PutOptions,
) (forddb.RawResource, error) {
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
	resource forddb.RawResource,
	opts forddb.PutOptions,
) (forddb.RawResource, forddb.RawResource, bool, error) {

	rs.m.Lock()
	defer rs.m.Unlock()

	current, err := rs.doGet(ctx, false, false)

	if forddb.IsNotFound(err) {
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
			if current.GetResourceVersion() != resource.GetResourceVersion() {
				return current, nil, false, forddb.ErrVersionMismatch
			}
		case forddb.OnConflictLatestWins:
			if current.GetResourceVersion() >= resource.GetResourceVersion() {
				return current, current, false, nil
			}
		case forddb.OnConflictReplace:
			// Nothing
		}
	}

	metadata := resource["metadata"].(map[string]interface{})
	metadata["version"] = resource.GetResourceVersion() + 1
	metadata["updated_at"] = time.Now()

	if resource.GetResourceMetadata().CreatedAt.IsZero() {
		metadata["created_at"] = metadata["updated_at"]
	}

	record, err := rs.table.db.log.Append(ctx, forddb.LogEntry{
		Kind:        forddb.LogEntryKindSet,
		Type:        rs.table.typ,
		ID:          rs.id.String(),
		Version:     resource.GetResourceVersion(),
		CurrentCid:  nil,
		PreviousCid: nil,
		Previous:    rs.encoded,
		Current:     resource,
	})

	if err != nil {
		return nil, nil, false, err
	}

	_, err = rs.table.db.storage.Put(ctx, resource, opts)

	if err != nil {
		return nil, nil, false, err
	}

	rs.lastRecord = record
	rs.encoded = resource

	rs.hasValue = true
	rs.err = nil

	return current, resource, true, nil
}

func (rs *resourceSlot) Delete(ctx context.Context) (forddb.RawResource, error) {
	previous, _, err := rs.doDelete(ctx)

	if err != nil {
		return nil, err
	}

	return previous, nil
}

func (rs *resourceSlot) doDelete(ctx context.Context) (forddb.RawResource, bool, error) {
	rs.cond.L.Lock()
	defer rs.cond.L.Unlock()

	value, err := rs.doGet(ctx, false, false)

	if err != nil {
		return nil, false, err
	}

	if !rs.table.typ.Type().IsRuntimeOnly() {
		_, err := rs.table.db.log.Append(ctx, forddb.LogEntry{
			Kind:        forddb.LogEntryKindDelete,
			Type:        rs.table.typ,
			ID:          rs.id.String(),
			Version:     value.GetResourceVersion(),
			CurrentCid:  nil,
			PreviousCid: nil,
			Previous:    value,
			Current:     nil,
		})

		if err != nil {
			return nil, false, err
		}
	}

	rs.encoded = nil
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
