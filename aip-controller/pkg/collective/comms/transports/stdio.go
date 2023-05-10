package transports

import (
	"context"
	"fmt"
	"io"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
)

type StdioTransport struct {
	Stdout io.Writer
}

func (t *StdioTransport) Subscribe(channel string) error {
	return nil
}

func (t *StdioTransport) Incoming() <-chan msn.Message {
	return nil
}

func (s *StdioTransport) RouteMessage(ctx context.Context, msg msn.Message) error {
	_, err := fmt.Fprintf(s.Stdout, "[%s] %s: %s\n", msg.Channel, msg.From, msg.Text)

	return err
}

func (s *StdioTransport) Close() error {
	return nil
}
