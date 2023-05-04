package transports

import (
	"context"

	"github.com/greenboxal/aip/aip-controller/pkg/collective"
)

type Transport interface {
	Subscribe(channel string) error

	Incoming() <-chan collective.Message
	RouteMessage(ctx context.Context, msg collective.Message) error

	Close() error
}
