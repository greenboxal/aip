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
	var err error

	switch record.Kind {
	case forddb.LogEntryKindDelete:
		fallthrough
	case forddb.LogEntryKindSet:
		id := record.Type.Type().CreateID(record.ID)
		slot := ed.db.GetSlot(record.Type, id, false)

		if slot == nil {
			break
		}

		var previous, current forddb.BasicResource

		if record.Previous != nil {
			previous, err = forddb.Decode(record.Previous)

			if err != nil {
				return err
			}
		}

		if record.Current != nil {
			current, err = forddb.Decode(record.Current)

			if err != nil {
				return err
			}
		}

		forddb.FireListeners(&slot.HasListenersBase, id, previous, current)
		forddb.FireListeners(&slot.table.HasListenersBase, id, previous, current)
		forddb.FireListeners(&ed.db.HasListenersBase, id, previous, current)
	}

	return nil
}
