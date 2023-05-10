package msn

import (
	"context"

	"github.com/jbenet/goprocess"
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

const eventDispatcherStreamID = "msn/eventdispatcher"

type eventDispatcher struct {
	forddb.LogConsumer

	db     forddb.Database
	router *Router
}

func newEventDispatcher(
	lc fx.Lifecycle,
	db forddb.Database,
	router *Router,
) *eventDispatcher {
	ed := &eventDispatcher{}
	ed.db = db
	ed.router = router
	ed.LogStore = db.LogStore()
	ed.StreamID = eventDispatcherStreamID
	ed.Handler = ed.processLogRecord

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			goprocess.Go(ed.Run)

			return nil
		},
	})

	return ed
}

func (ed *eventDispatcher) processLogRecord(ctx context.Context, record *forddb.LogEntryRecord) error {
	if record.Type != (MessageID{}).BasicResourceType().GetResourceID() {
		return nil
	}

	switch record.Kind {
	case forddb.LogEntryKindSet:
		current, err := forddb.Decode(record.Current)

		if err != nil {
			return err
		}

		ed.router.Dispatch(Event{
			MessageEvent: &MessageEvent{
				Message: *current.(*Message),
			},
		})
	}

	return nil
}
