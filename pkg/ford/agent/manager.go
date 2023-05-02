package agent

import (
	"context"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"sync/atomic"

	"go.uber.org/zap"

	"github.com/greenboxal/aip/pkg/collective"
	"github.com/greenboxal/aip/pkg/collective/comms"
	"github.com/greenboxal/aip/pkg/collective/supervisor"
	"github.com/greenboxal/aip/pkg/config"
	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type Manager struct {
	logger *zap.SugaredLogger

	db         forddb.Database
	rsm        *config.ResourceManager
	routing    *comms.Routing
	supervisor *supervisor.Manager

	aidCounter atomic.Uint64

	procDir string
}

func NewManager(
	logger *zap.SugaredLogger,
	rsm *config.ResourceManager,
	db forddb.Database,
	routing *comms.Routing,
	sup *supervisor.Manager,
) (*Manager, error) {
	procDir := rsm.GetProcDirectory("agents")

	_ = os.RemoveAll(procDir)

	if err := os.MkdirAll(procDir, 0755); err != nil {
		return nil, err
	}

	return &Manager{
		logger: logger.Named("agent-manager"),

		db:         db,
		rsm:        rsm,
		routing:    routing,
		supervisor: sup,

		procDir: procDir,
	}, nil
}

func (m *Manager) StartAgent(ctx context.Context, task *collective.Agent) error {
	aid := m.aidCounter.Add(1)
	procDir := path.Join(m.procDir, fmt.Sprintf("%d", aid))

	if err := os.MkdirAll(procDir, 0755); err != nil {
		return err
	}

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

	tmpProfilePath := path.Join(procDir, "profile.json")
	tmpProfile, err := os.OpenFile(tmpProfilePath, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	if err := forddb.SerializeTo(tmpProfile, forddb.Json, profile); err != nil {
		return err
	}

	if err := tmpProfile.Close(); err != nil {
		return err
	}

	args := []string{
		"-i", task.Name,
		"-p", tmpProfile.Name(),
	}

	err = os.WriteFile(path.Join(procDir, "args"), []byte(strings.Join(args, "\000")), 0644)

	if err != nil {
		return err
	}

	args = append(args, task.Spec.ExtraArgs...)

	p, err := m.supervisor.Spawn(
		supervisor.WithID(task.Name),
		supervisor.WithArgs(args...),
		supervisor.WithPort(port),
	)

	err = os.WriteFile(path.Join(procDir, "pid"), []byte(strconv.Itoa(p.Process.Pid())), 0644)

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) StopAgent(ctx context.Context, previous *collective.Agent) error {
	return nil
}
