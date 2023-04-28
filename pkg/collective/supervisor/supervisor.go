package supervisor

import (
	"context"

	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/greenboxal/aip/pkg/collective"
)

type Supervisor struct {
	Config  *Config
	Process *Process
}

func NewSupervisor(logger *zap.SugaredLogger, config *Config) (*Supervisor, error) {
	sup := &Supervisor{
		Config: config,
	}

	program := "/usr/bin/env"

	args := []string{
		"bash",
		"-c",
		`set -eu; if [ -f .env ]; then source .env; fi; exec python -m aip ipc "$@"`,
		"--",
	}

	args = append(args, config.Args...)

	proc, err := NewProcess(logger, func(m collective.Message) {
		if err := config.Port.Send(context.Background(), m); err != nil {
			logger.Error(err)
		}
	}, program, args...)

	if err != nil {
		return nil, err
	}

	sup.Process = proc

	return sup, nil
}

func (s *Supervisor) Run(proc goprocess.Process) error {
	ctx := goprocessctx.OnClosingContext(proc)

	wg, gctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		return s.Process.Run(gctx)
	})

	wg.Go(func() error {
		for {
			select {
			case <-gctx.Done():
				return gctx.Err()

			case msg := <-s.Config.Port.Incoming():
				if msg.From == s.Config.ID {
					continue
				}

				if err := s.Process.Send(msg); err != nil {
					return err
				}
			}
		}
	})

	return wg.Wait()
}

func (s *Supervisor) Close() error {
	return s.Process.Close()
}
