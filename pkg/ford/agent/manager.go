package agent

import (
	"context"
	"os"

	"go.uber.org/zap"

	"github.com/greenboxal/aip/pkg/collective"
	"github.com/greenboxal/aip/pkg/collective/comms"
	"github.com/greenboxal/aip/pkg/collective/supervisor"
	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type Manager struct {
	logger     *zap.SugaredLogger
	db         forddb.Database
	routing    *comms.Routing
	supervisor *supervisor.Manager
}

func NewManager(
	logger *zap.SugaredLogger,
	db forddb.Database,
	routing *comms.Routing,
	sup *supervisor.Manager,
) *Manager {
	return &Manager{
		logger: logger.Named("agent-manager"),

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

	profile, err := forddb.Get[*collective.Profile](m.db, task.Spec.ProfileID)

	if err != nil {
		return err
	}

	tmpProfile, err := os.CreateTemp(os.TempDir(), "aip-profile-")

	if err != nil {
		return err
	}

	if err := forddb.SerializeTo(tmpProfile, forddb.Json, profile); err != nil {
		return err
	}

	args := []string{
		"-i", task.Name,
		"-p", tmpProfile.Name(),
	}

	args = append(args, task.Spec.ExtraArgs...)

	m.logger.Infow("args", "args", args)

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
