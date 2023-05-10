package transports

import (
	"context"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
)

type Port interface {
	Subscribe(channel string) error

	Incoming() <-chan msn.Message

	Send(ctx context.Context, msg msn.Message) error

	Close() error
}
