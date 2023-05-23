package graphql

import (
	"context"
	"sync"
	"time"

	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"
	"go.uber.org/fx"
	"golang.org/x/exp/slices"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-forddb/pkg/logstore"
)

type ResourceEventType string

const (
	ResourceEventTypeInvalid ResourceEventType = ""
	ResourceEventTypeCreated ResourceEventType = "created"
	ResourceEventTypeUpdated ResourceEventType = "updated"
	ResourceEventTypeDeleted ResourceEventType = "deleted"
)

type ResourceEvent struct {
	Type ResourceEventType `json:"type"`

	Payload ResourceEventPayload `json:"payload"`
}

type ResourceEventPayload struct {
	IDs []string `json:"ids"`
}

type SubscriptionManager struct {
	db forddb.Database

	m             sync.RWMutex
	subscriptions map[forddb.BasicResourceType][]chan ResourceEvent
}

func NewSubscriptionManager(
	lc fx.Lifecycle,
	db forddb.Database,
) *SubscriptionManager {
	sm := &SubscriptionManager{
		db:            db,
		subscriptions: map[forddb.BasicResourceType][]chan ResourceEvent{},
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			goprocess.Go(sm.Run)

			return nil
		},
	})

	return sm
}

func (sm *SubscriptionManager) Run(proc goprocess.Process) {
	ctx := goprocessctx.OnClosingContext(proc)

	iterator := sm.db.LogStore().Iterator(forddb.WithBlockingIterator())

	for {
		if !iterator.Next(ctx) {
			if err := iterator.SetLSN(ctx, forddb.MakeLSN(logstore.FileSegmentBaseSeekToHead, time.Now())); err != nil {
				return
			}
			continue
		}

		err := iterator.Error()

		if err != nil {
			panic(iterator.Error())
		}

		record := iterator.Record()

		sm.dispatch(record)
	}
}

func (sm *SubscriptionManager) Subscribe(
	ctx context.Context,
	typ forddb.BasicResourceType,
) (<-chan ResourceEvent, error) {
	ch := make(chan ResourceEvent, 128)

	sm.addSubscription(typ, ch)

	go func() {
		//defer close(ch)
		//defer sm.removeSubscription(typ, ch)

		<-ctx.Done()
	}()

	return ch, nil
}

func (sm *SubscriptionManager) addSubscription(typ forddb.BasicResourceType, ch chan ResourceEvent) {
	sm.m.Lock()
	defer sm.m.Unlock()

	subs := sm.subscriptions[typ]

	index := slices.Index(subs, ch)

	if index != -1 {
		return
	}

	index = slices.Index(subs, nil)

	if index != -1 {
		subs[index] = ch
	} else {
		subs = append(subs, ch)
	}

	sm.subscriptions[typ] = subs
}

func (sm *SubscriptionManager) removeSubscription(typ forddb.BasicResourceType, ch chan ResourceEvent) {
	sm.m.Lock()
	defer sm.m.Unlock()

	subs := sm.subscriptions[typ]

	if len(subs) == 0 {
		return
	}

	index := slices.Index(subs, ch)

	if index == -1 {
		return
	}

	subs = slices.Delete(subs, index, index+1)

	sm.subscriptions[typ] = subs
}

func (sm *SubscriptionManager) dispatch(record *forddb.LogEntryRecord) {
	event := ResourceEvent{}

	switch record.Kind {
	case forddb.LogEntryKindSet:
		if record.Previous == nil {
			event.Type = ResourceEventTypeCreated
		} else {
			event.Type = ResourceEventTypeUpdated
		}
	case forddb.LogEntryKindDelete:
		event.Type = ResourceEventTypeDeleted
	}

	event.Payload.IDs = []string{record.ID}

	sm.m.RLock()
	defer sm.m.RUnlock()

	subs := sm.subscriptions[record.Type.Type()]

	if len(subs) == 0 {
		return
	}

	for _, sub := range subs {
		if sub == nil {
			continue
		}

		sub <- event
	}
}
