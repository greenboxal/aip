package ford

import (
	"context"

	"github.com/greenboxal/aip/pkg/collective"
	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type Manager struct {
	db forddb.Database
}

func NewManager(db forddb.Database) *Manager {
	return &Manager{
		db: db,
	}
}

func (m *Manager) StartAgent(ctx context.Context, task *collective.Agent) error {
	return nil
}

func (m *Manager) StopAgent(ctx context.Context, previous *collective.Agent) error {
	return nil
}
