package forddbimpl

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
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

	latestRecord forddb.LogEntryRecord
	latestValue  forddb.RawResource

	lastError error
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

	result, err := forddb.DecodeAs[forddb.BasicResource](raw)

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
	resource, err := rs.table.db.storage.Get(ctx, rs.table.typ, rs.id, forddb.GetOptions{})

	if err != nil && !forddb.IsNotFound(err) {
		return nil, err
	}

	if err == nil {
		if lock {
			rs.m.Lock()
			defer rs.m.Unlock()

			lock = false
		}

		unknown := forddb.UnknownResource{RawResource: resource}

		if rs.latestRecord.Version <= unknown.GetResourceVersion() {
			rs.hasValue = true
			rs.latestValue = resource

			return rs.latestValue, err
		}
	}

	wait = false

	if lock {
		rs.m.Lock()
		defer rs.m.Unlock()
	}

	if wait && !rs.hasValue && rs.lastError == nil {
		for !rs.hasValue && rs.lastError == nil {
			rs.table.notifyGet(rs)

			rs.cond.Wait()
		}
	}

	if rs.lastError != nil {
		return nil, rs.lastError
	}

	if !rs.hasValue {
		return nil, forddb.ErrNotFound
	}

	raw = rs.latestValue

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

	unknown := forddb.UnknownResource{RawResource: resource}
	metadata := unknown.GetResourceMetadata()

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

		currentUnk := forddb.UnknownResource{RawResource: current}

		switch opts.OnConflict {
		case forddb.OnConflictError:
			return current, nil, false, forddb.ErrVersionMismatch
		case forddb.OnConflictOptimistic:
			if currentUnk.GetResourceVersion() != metadata.GetResourceVersion() {
				return current, nil, false, forddb.ErrVersionMismatch
			}
		case forddb.OnConflictLatestWins:
			if currentUnk.GetResourceVersion() >= metadata.GetResourceVersion() {
				return current, current, false, nil
			}
		case forddb.OnConflictReplace:
			// Nothing
		}
	}

	newMetadata := *metadata
	newMetadata.Version++
	newMetadata.UpdatedAt = time.Now()

	if newMetadata.CreatedAt.IsZero() {
		newMetadata.CreatedAt = newMetadata.UpdatedAt
	}

	builder := resource.Prototype().NewBuilder()

	if err := builder.AssignNode(resource); err != nil {
		return current, nil, false, err
	}

	root, err := builder.BeginMap(-1)
	if err != nil {
		panic(err)
	}

	metadataNode, err := root.AssembleEntry("metadata")
	if err != nil {
		panic(err)
	}

	if err := metadataNode.AssignNode(typesystem.Wrap(newMetadata)); err != nil {
		panic(err)
	}

	resource = builder.Build()

	res, err := forddb.DecodeAs[forddb.BasicResource](resource)

	if err != nil {
		return current, nil, false, err
	}

	record, err := rs.table.db.log.Append(ctx, forddb.LogEntry{
		Kind:     forddb.LogEntryKindSet,
		Type:     rs.table.typ,
		ID:       rs.id.String(),
		Version:  metadata.Version,
		Previous: rs.latestRecord.Current,
		Current:  res,
	})

	if err != nil {
		return nil, nil, false, err
	}

	_, err = rs.table.db.storage.Put(ctx, resource, opts)

	if err != nil {
		return nil, nil, false, err
	}

	rs.latestRecord = record
	rs.latestValue = resource

	rs.hasValue = true
	rs.lastError = nil

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

	res, err := forddb.DecodeAs[forddb.BasicResource](value)

	if err != nil {
		return nil, false, err
	}

	if !rs.table.typ.Type().IsRuntimeOnly() {
		_, err := rs.table.db.log.Append(ctx, forddb.LogEntry{
			Kind:     forddb.LogEntryKindDelete,
			Type:     rs.table.typ,
			ID:       rs.id.String(),
			Version:  res.GetResourceVersion(),
			Previous: res,
			Current:  nil,
		})

		if err != nil {
			return nil, false, err
		}
	}

	rs.latestValue = nil
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

	rs.lastError = err
	rs.cond.Broadcast()
}
