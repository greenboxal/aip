package supervisor

import (
	"context"
	"sync"

	"github.com/jbenet/goprocess"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Manager struct {
	m sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc

	logger *zap.SugaredLogger

	children map[string]*Supervisor
}

func NewManager(lc fx.Lifecycle, logger *zap.SugaredLogger) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	m := &Manager{
		logger: logger.Named("supervisor"),

		ctx:    ctx,
		cancel: cancel,

		children: map[string]*Supervisor{},
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return m.Close()
		},
	})

	return m
}

func (m *Manager) Child(child string) *Supervisor {
	m.m.RLock()
	defer m.m.RUnlock()

	return m.children[child]
}

func (m *Manager) Children() []*Supervisor {
	m.m.RLock()
	defer m.m.RUnlock()

	result := make([]*Supervisor, 0, len(m.children))

	for _, sup := range m.children {
		result = append(result, sup)
	}

	return result
}

func (m *Manager) Spawn(options ...SupervisorOption) (*Supervisor, error) {
	config := NewConfig(options...)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	m.m.Lock()
	defer m.m.Unlock()

	if _, ok := m.children[config.ID]; ok {
		return nil, ErrSupervisorAlreadyExists
	}

	sup, err := NewSupervisor(m.logger.Named("sup-child-"+config.ID), config)

	if err != nil {
		return nil, err
	}

	m.children[config.ID] = sup

	goprocess.Go(func(proc goprocess.Process) {
		defer func() {
			if e := recover(); e != nil {
				_ = sup.Close()

				m.logger.Errorw("supervisor panic", "id", config.ID, "error", e)
			}
		}()

		m.logger.Infow("supervisor started", "id", config.ID)

		if err := sup.Run(proc); err != nil {
			m.logger.Errorw("supervisor panic", "id", config.ID, "error", err)
		}
	})

	return sup, nil
}

func (m *Manager) Close() error {
	m.cancel()

	return nil
}
