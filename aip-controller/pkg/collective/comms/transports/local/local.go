package local

import (
	"context"
	"errors"
	"sync"

	transports2 "github.com/greenboxal/aip/aip-controller/pkg/collective/comms/transports"
	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
)

type Transport struct {
	m              sync.RWMutex
	knownPorts     map[string]*Port
	subscriptions  map[string]map[string]struct{}
	defaultGateway transports2.Transport
}

func NewTransport(defaultGateway transports2.Transport) *Transport {
	t := &Transport{
		knownPorts:     map[string]*Port{},
		subscriptions:  map[string]map[string]struct{}{},
		defaultGateway: defaultGateway,
	}

	go func() {
		ch := t.defaultGateway.Incoming()

		if ch == nil {
			return
		}

		for msg := range ch {
			_ = t.routeMessage(context.Background(), msg, false)
		}
	}()

	return t
}

func (t *Transport) Subscribe(channel string) error {
	return t.defaultGateway.Subscribe(channel)
}

func (t *Transport) Incoming() <-chan msn.Message {
	return nil
}

func (t *Transport) AddPort(name string) (transports2.Port, error) {
	t.m.Lock()
	defer t.m.Unlock()

	if existing := t.knownPorts[name]; existing != nil {
		return nil, errors.New("port already exists")
	}

	p := NewPort(t, name)

	t.knownPorts[name] = p

	return p, nil
}

func (t *Transport) RouteMessage(ctx context.Context, msg msn.Message) error {
	return t.routeMessage(ctx, msg, true)
}

func (t *Transport) Close() error {
	return nil
}

func (t *Transport) routeMessage(ctx context.Context, msg msn.Message, allowExternal bool) error {
	t.routeSubscriptions(ctx, msg)

	if msg.ChannelID.String() != msg.From.String() {
		p := t.getPort(msg.ChannelID.String())

		if p != nil {
			p.routeMessage(ctx, msg)

			return nil
		}
	}

	if t.defaultGateway == nil || !allowExternal {
		return errors.New("no route to channel")
	}

	return t.defaultGateway.RouteMessage(ctx, msg)
}

func (t *Transport) getPort(name string) *Port {
	t.m.RLock()
	defer t.m.RUnlock()

	return t.knownPorts[name]
}

func (t *Transport) removePort(p *Port) {
	t.m.Lock()
	defer t.m.Unlock()

	existing := t.knownPorts[p.name]

	if existing == p {
		delete(t.knownPorts, p.name)
	}
}

func (t *Transport) subscribeLocal(p *Port, channel string) error {
	t.m.Lock()
	defer t.m.Unlock()

	subs := t.subscriptions[channel]

	if subs == nil {
		subs = map[string]struct{}{}
		t.subscriptions[channel] = subs
	}

	subs[p.name] = struct{}{}

	return t.Subscribe(channel)
}

func (t *Transport) routeSubscriptions(ctx context.Context, msg msn.Message) {
	t.m.RLock()
	defer t.m.RUnlock()

	subs := t.subscriptions[msg.ChannelID.String()]

	if subs == nil {
		return
	}

	for key := range subs {
		if key == msg.From.String() {
			continue
		}

		p := t.knownPorts[key]

		if p == nil {
			continue
		}

		p.routeMessage(ctx, msg)
	}
}
