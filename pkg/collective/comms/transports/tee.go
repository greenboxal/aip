package transports

import (
	"context"
	"reflect"
	"sync"

	"go.uber.org/multierr"
	"golang.org/x/exp/slices"

	"github.com/greenboxal/aip/pkg/collective"
)

func Tee(transports ...Transport) *TeeTransport {
	return &TeeTransport{
		targets: transports,
	}
}

type TeeTransport struct {
	m          sync.Mutex
	targets    []Transport
	incomingCh chan collective.Message
}

func (t *TeeTransport) Subscribe(channel string) error {
	var merr error

	for _, target := range t.targets {
		if err := target.Subscribe(channel); err != nil {
			merr = multierr.Append(merr, err)
		}
	}

	return merr
}

func (t *TeeTransport) Incoming() <-chan collective.Message {
	if t.incomingCh == nil {
		t.m.Lock()
		defer t.m.Unlock()

		if t.incomingCh != nil {
			return t.incomingCh
		}

		t.incomingCh = make(chan collective.Message, 16)

		go func() {
			cases := make([]reflect.SelectCase, 0, len(t.targets))

			for _, target := range t.targets {
				ch := target.Incoming()

				if ch == nil {
					continue
				}

				cases = append(cases, reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(ch),
				})
			}

			for len(cases) > 0 {
				i, value, ok := reflect.Select(cases)

				if !ok {
					cases = slices.Delete(cases, i, i+1)
					continue
				}

				t.incomingCh <- value.Interface().(collective.Message)
			}
		}()
	}

	return t.incomingCh
}

func (t *TeeTransport) RouteMessage(ctx context.Context, msg collective.Message) error {
	var merr error

	for _, target := range t.targets {
		if err := target.RouteMessage(ctx, msg); err != nil {
			merr = multierr.Append(merr, err)
		}
	}

	return merr
}

func (t *TeeTransport) Close() error {
	return nil
}
