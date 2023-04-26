package supervisor

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/greenboxal/aip/pkg/collective"
)

type Config struct {
	ID string

	Program string
	Args    []string
}

type Supervisor struct {
	ctx    context.Context
	cancel context.CancelFunc

	config *Config

	process *Process
	port    collective.Port
}

func NewSupervisor(ctx context.Context, config *Config, port collective.Port) (*Supervisor, error) {
	ctx, cancel := context.WithCancel(ctx)

	proc, err := NewProcess(ctx, config.Program, config.Args...)

	if err != nil {
		cancel()
		return nil, err
	}

	return &Supervisor{
		ctx:    ctx,
		cancel: cancel,

		config:  config,
		process: proc,
		port:    port,
	}, nil
}

func (s *Supervisor) Run() error {
	wg, gctx := errgroup.WithContext(s.ctx)

	wg.Go(func() error {
		return s.process.Run()
	})

	wg.Go(func() error {
		for {
			select {
			case <-gctx.Done():
				return gctx.Err()

			case msg := <-s.port.Incoming():
				if msg.From == s.config.ID {
					continue
				}

				s.process.Incoming() <- msg

			case msg := <-s.process.Outgoing():
				if err := s.port.Send(gctx, msg); err != nil {
					return err
				}
			}
		}
	})

	return wg.Wait()
}

func (s *Supervisor) Close() error {
	s.cancel()

	return s.process.Close()
}
