package local

import (
	"context"

	"github.com/greenboxal/aip/pkg/collective"
)

type Port struct {
	name     string
	local    *Transport
	incoming chan collective.Message
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
