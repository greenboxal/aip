package agent

import (
	"context"

	"github.com/greenboxal/aip/pkg/collective"
	"github.com/greenboxal/aip/pkg/collective/comms"
	"github.com/greenboxal/aip/pkg/collective/supervisor"
	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type Manager struct {
	db         forddb.Database
	routing    *comms.Routing
	supervisor *supervisor.Manager
}

func NewManager(db forddb.Database, routing *comms.Routing, sup *supervisor.Manager) *Manager {
	return &Manager{
		db:         db,
		routing:    routing,
		supervisor: sup,
	}
}

func (m *Manager) StartAgent(ctx context.Context, task *collective.Agent) error {
	port, err := m.routing.AddPort(task.Spec.PortID)

	if err != nil {
		return err
	}

	if err := port.Subscribe(task.Name); err != nil {
		return err
	}

	if err := port.Subscribe("aip-bod-room"); err != nil {
		return err
	}

	args := []string{
		"-i", task.Name,
	}

	args = append(args, task.Spec.ExtraArgs...)

	_, err = m.supervisor.Spawn(
		supervisor.WithID(task.Name),
		supervisor.WithArgs(args...),
		supervisor.WithPort(port),
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) StopAgent(ctx context.Context, previous *collective.Agent) error {
	return nil
}
