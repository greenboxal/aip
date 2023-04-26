package collective

import "context"

type Port interface {
	Subscribe(channel string) error

	Incoming() <-chan Message

	Send(ctx context.Context, msg Message) error

	Close() error
}
