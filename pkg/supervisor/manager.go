package supervisor

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/greenboxal/aip/pkg/collective"
)

type Manager struct {
	ctx    context.Context
	cancel context.CancelFunc

	logger *zap.SugaredLogger
}

func NewManager(lc fx.Lifecycle, logger *zap.SugaredLogger) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	m := &Manager{
		logger: logger.Named("supervisor"),

		ctx:    ctx,
		cancel: cancel,
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return m.Close()
		},
	})

	return m
}

func (m *Manager) Spawn(config *Config, port collective.Port) (*Supervisor, error) {
	sup, err := NewSupervisor(m.ctx, config, port)

	if err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			if e := recover(); e != nil {
				_ = sup.Close()

				m.logger.Errorw("supervisor panic", "id", config.ID, "error", e)
			}
		}()

		m.logger.Infow("supervisor started", "id", config.ID)

		if err := sup.Run(); err != nil {
			m.logger.Errorw("supervisor panic", "id", config.ID, "error", err)
		}
	}()

	return sup, nil
}

func (m *Manager) Close() error {
	m.cancel()

	return nil
}
