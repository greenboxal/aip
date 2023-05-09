package forddbimpl

import (
	"context"

	"github.com/ipfs/go-cid"

	forddb "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	logstore2 "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb/logstore"
)

const objectIndexerStreamID = "forddb/objectindexer"

type objectIndexer struct {
	logstore2.LogConsumer

	db *database
}

func newObjectIndexer(db *database) *objectIndexer {
	oi := &objectIndexer{}
	oi.db = db
	oi.LogStore = db.log
	oi.StreamID = objectIndexerStreamID
	oi.Handler = oi.processLogRecord

	return oi
}

func (oi *objectIndexer) processLogRecord(ctx context.Context, record *logstore2.LogEntryRecord) error {
	switch record.Kind {
	case logstore2.LogEntryKindSet:
		if err := oi.prefetchObject(ctx, record, record.Version > 1, true); err != nil {
			if !forddb.IsNotFound(err) {
				return err
			}
		}

		if _, err := oi.db.storage.Put(ctx, record.Current, forddb.PutOptions{
			OnConflict: forddb.OnConflictReplace,
		}); err != nil {
			return err
		}

		forddb.FireListeners(&oi.db.HasListenersBase, record.ID, record.CachedPrevious, record.CachedCurrent)

	case logstore2.LogEntryKindDelete:
		if err := oi.prefetchObject(ctx, record, true, false); err != nil {
			if !forddb.IsNotFound(err) {
				return err
			}
		}

		if _, err := oi.db.storage.Delete(ctx, record.Previous, forddb.DeleteOptions{}); err != nil {
			return err
		}

		forddb.FireListeners(&oi.db.HasListenersBase, record.ID, record.CachedPrevious, record.CachedCurrent)
	}

	return nil
}

func (oi *objectIndexer) prefetchObject(
	ctx context.Context,
	record *logstore2.LogEntryRecord,
	previous bool,
	current bool,
) error {
	if current {
		err := oi.doPrefetchObject(
			ctx,
			&record.CachedCurrent,
			record.Current,
			record.CurrentCid,
			record.Version,
			record.ID,
		)

		if err != nil {
			return err
		}
	}

	if previous {
		err := oi.doPrefetchObject(
			ctx,
			&record.CachedPrevious,
			record.Previous,
			record.PreviousCid,
			record.Version-1,
			record.ID,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (oi *objectIndexer) doPrefetchObject(
	ctx context.Context,
	dst *forddb.BasicResource,
	raw forddb.RawResource,
	cid *cid.Cid,
	version uint64,
	id forddb.BasicResourceID,
) error {
	if *dst == nil && raw != nil {
		res, err := forddb.Decode(raw)

		if err != nil {
			return err
		}

		*dst = res
	}

	if *dst == nil {
		res, err := oi.db.Get(ctx, id.BasicResourceType().GetResourceID(), id)

		if err == nil && res.GetResourceVersion() == version {
			*dst = res
		}
	}

	// TODO: Implement
	if *dst == nil && cid != nil {
		return nil
	}

	// FIXME: figure out what this is supposed to do, like scanning the log?

	if *dst == nil {
		return forddb.ErrNotFound
	}

	return nil
}
