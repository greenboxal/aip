package transports

import (
	"context"

	"github.com/greenboxal/aip/aip-controller/pkg/collective"
)

type Port interface {
	Subscribe(channel string) error

	Incoming() <-chan collective.Message

	Send(ctx context.Context, msg collective.Message) error

	Close() error
}