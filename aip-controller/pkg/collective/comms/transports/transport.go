package transports

import (
	"context"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
)

type Transport interface {
	Subscribe(channel string) error

	Incoming() <-chan msn.Message
	RouteMessage(ctx context.Context, msg msn.Message) error

	Close() error
}
