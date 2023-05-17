package msn

import (
	"context"
	"reflect"
	"sync"

	gql "github.com/graphql-go/graphql"
	"go.uber.org/zap"

	"github.com/greenboxal/aip/aip-forddb/pkg/apis/graphql"
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
)

type RealTimeService struct {
	logger *zap.SugaredLogger
	db     forddb.Database
	router *Router
}

func NewRealTimeService(
	logger *zap.SugaredLogger,
	db forddb.Database,
	router *Router,
) *RealTimeService {
	return &RealTimeService{
		logger: logger.Named("msn-rt"),
		db:     db,
		router: router,
	}
}

func (r *RealTimeService) BindResource(ctx graphql.BindingContext) {
	ctx.RegisterSubscription(&gql.Field{
		Name: "realTimeEvents",

		Type: ctx.LookupOutputType(typesystem.TypeOf(Event{})),

		Args: gql.FieldConfigArgument{
			"endpoint": &gql.ArgumentConfig{
				Type: gql.NewNonNull(gql.String),
			},
		},

		Subscribe: func(p gql.ResolveParams) (interface{}, error) {
			return r.SubscribeToEvents(p.Context, &SubscribeToEventsRequest{
				EndpointID: forddb.NewStringID[EndpointID](p.Args["endpoint"].(string)),
			})
		},
	})
}

func (r *RealTimeService) SubscribeToEvents(ctx context.Context, req *SubscribeToEventsRequest) (<-chan Event, error) {
	ch := make(chan Event, 128)
	sub := newSubscriber(r.router, ch)

	endpoint, err := forddb.Get[*Endpoint](ctx, r.db, req.EndpointID)

	if err != nil {
		return nil, err
	}

	_ = endpoint.GetResourceNode().Subscribe(forddb.ListenerFunc(func(
		id forddb.BasicResourceID,
		previous, current forddb.BasicResource,
	) {
		endpoint := current.(*Endpoint)

		for _, id := range endpoint.Subscriptions {
			sub.Subscribe(id)
		}
	}))

	for _, id := range endpoint.Subscriptions {
		sub.Subscribe(id)
	}

	go sub.Run(context.Background())

	//go func() {

	//	<-ctx.Done()

	//	cancel()
	//	sub.UnsubscribeAll()
	//}()

	return ch, nil
}

type subscriber struct {
	m             sync.RWMutex
	router        *Router
	subscriptions map[ChannelID]func()
	outbox        chan<- Event
	inbox         chan Event
}

func newSubscriber(router *Router, outbox chan<- Event) *subscriber {
	return &subscriber{
		outbox: outbox,
		router: router,

		inbox:         make(chan Event, 128),
		subscriptions: map[ChannelID]func(){},
	}
}

func (s *subscriber) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case event := <-s.inbox:
			s.outbox <- event
		}
	}
}

func (s *subscriber) Subscribe(channel ChannelID) {
	table := s.router.getTable(channel, true)

	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.subscriptions[channel]; ok {
		return
	}

	s.subscriptions[channel] = table.Subscribe(s.inbox)
}

func (s *subscriber) Unsubscribe(channel ChannelID) {
	s.m.Lock()
	defer s.m.Unlock()

	if f, ok := s.subscriptions[channel]; ok {
		f()

		delete(s.subscriptions, channel)
	}
}

func (s *subscriber) UnsubscribeAll() {
	s.m.Lock()
	defer s.m.Unlock()

	for _, f := range s.subscriptions {
		f()
	}

	s.subscriptions = map[ChannelID]func(){}
}

func (s *subscriber) prepareTargets(ctx context.Context, targets []reflect.SelectCase) []reflect.SelectCase {
	s.m.RLock()
	defer s.m.RUnlock()

	if len(s.subscriptions) > len(targets)+1 {
		targets = make([]reflect.SelectCase, 0, len(s.subscriptions)+1)
	} else {
		targets = targets[:0]
	}

	targets = append(targets, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ctx.Done()),
	})

	for _, v := range s.subscriptions {
		targets = append(targets, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(v),
		})
	}

	return targets
}
