package forddbimpl

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"

	"github.com/greenboxal/aip/pkg/ford/forddb"
)

const objectIndexerStreamID = "forddb/objectindexer"

type objectIndexer struct {
	forddb.LogConsumer

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

func (oi *objectIndexer) processLogRecord(ctx context.Context, record *forddb.LogEntryRecord) error {
	switch record.Kind {
	case forddb.LogEntryKindSet:
		if err := oi.prefetchObject(ctx, record, record.Version > 1, true); err != nil {
			return err
		}

		if _, err := oi.db.storage.Put(ctx, record.CachedCurrent); err != nil {
			return err
		}

		forddb.FireListeners(&oi.db.HasListenersBase, record.ID, record.CachedPrevious, record.CachedCurrent)

	case forddb.LogEntryKindDelete:
		if err := oi.prefetchObject(ctx, record, true, false); err != nil {
			return err
		}

		if _, err := oi.db.storage.Delete(ctx, record.CachedPrevious); err != nil {
			return err
		}

		forddb.FireListeners(&oi.db.HasListenersBase, record.ID, record.CachedPrevious, record.CachedCurrent)
	}

	return nil
}

func (oi *objectIndexer) prefetchObject(
	ctx context.Context,
	record *forddb.LogEntryRecord,
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
		res, err := oi.db.Get(ctx, id.BasicResourceType().GetID(), id)

		if err == nil && res.GetVersion() == version {
			*dst = res
		}
	}

	// TODO: Implement
	if *dst == nil && cid != nil {
		return nil
	}

	// FIXME: figure out what this is supposed to do, like scanning the log?

	if *dst == nil {
		return errors.New("could not fetch object")
	}

	return nil
}
