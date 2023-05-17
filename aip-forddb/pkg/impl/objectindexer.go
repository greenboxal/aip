package forddbimpl

import (
	"context"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
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
		encoded, err := forddb.Encode(record.Current)

		if err != nil {
			return err
		}

		if _, err := oi.db.storage.Put(ctx, encoded, forddb.PutOptions{
			OnConflict: forddb.OnConflictReplace,
		}); err != nil {
			return err
		}

	case forddb.LogEntryKindDelete:
		encoded, err := forddb.Encode(record.Previous)

		if err != nil {
			return err
		}

		if _, err := oi.db.storage.Delete(ctx, encoded, forddb.DeleteOptions{}); err != nil {
			return err
		}
	}

	return nil
}
