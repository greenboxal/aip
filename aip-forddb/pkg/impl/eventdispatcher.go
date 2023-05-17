package forddbimpl

import (
	"context"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

const eventDispatcherStreamID = "forddb/eventdispatcher"

type eventDispatcher struct {
	forddb.LogConsumer

	db *database
}

func newEventDispatcher(db *database) *eventDispatcher {
	oi := &eventDispatcher{}
	oi.db = db
	oi.LogStore = db.log
	oi.StreamID = eventDispatcherStreamID
	oi.Handler = oi.processLogRecord

	return oi
}

func (ed *eventDispatcher) processLogRecord(ctx context.Context, record *forddb.LogEntryRecord) error {
	switch record.Kind {
	case forddb.LogEntryKindDelete:
		fallthrough
	case forddb.LogEntryKindSet:
		id := record.Type.Type().CreateID(record.ID)
		slot := ed.db.GetSlot(record.Type, id, false)

		if slot == nil {
			break
		}

		forddb.FireListeners(&slot.HasListenersBase, id, record.Previous, record.Current)
		forddb.FireListeners(&slot.table.HasListenersBase, id, record.Previous, record.Current)
		forddb.FireListeners(&ed.db.HasListenersBase, id, record.Previous, record.Current)
	}

	return nil
}
