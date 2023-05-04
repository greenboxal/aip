package local

import (
	"context"
	"sync"
	"time"

	"github.com/greenboxal/aip/aip-controller/pkg/collective"
)

const SeenMessageTimeout = 60 * time.Second

type Port struct {
	m       sync.Mutex
	seenSet map[string]time.Time

	name     string
	local    *Transport
	incoming chan collective.Message
}

func NewPort(local *Transport, name string) *Port {
	return &Port{
		name:     name,
		local:    local,
		incoming: make(chan collective.Message, 16),
		seenSet:  map[string]time.Time{},
	}
}

func (p *Port) Subscribe(channel string) error {
	return p.local.subscribeLocal(p, channel)
}

func (p *Port) Incoming() <-chan collective.Message {
	return p.incoming
}

func (p *Port) Send(ctx context.Context, msg collective.Message) error {
	return p.local.RouteMessage(ctx, msg)
}

func (p *Port) Close() error {
	close(p.incoming)
	p.incoming = nil

	p.local.removePort(p)

	return nil
}

func (p *Port) routeMessage(ctx context.Context, msg collective.Message) {
	if !p.isMessageVisible(msg) {
		return
	}

	p.incoming <- msg
}

func (p *Port) isMessageVisible(msg collective.Message) bool {
	defer p.cleanSeenSet()

	p.m.Lock()
	defer p.m.Unlock()

	_, ok := p.seenSet[msg.ID]

	if ok {
		return false
	}

	p.seenSet[msg.ID] = time.Now()

	return true
}

func (p *Port) cleanSeenSet() {
	p.m.Lock()
	defer p.m.Unlock()

	for k, v := range p.seenSet {
		if time.Since(v) > SeenMessageTimeout {
			delete(p.seenSet, k)
		}
	}
}
