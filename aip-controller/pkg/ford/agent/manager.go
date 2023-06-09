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

	collective2 "github.com/greenboxal/aip/aip-controller/pkg/collective"
	"github.com/greenboxal/aip/aip-controller/pkg/collective/comms"
	supervisor2 "github.com/greenboxal/aip/aip-controller/pkg/collective/supervisor"
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-sdk/pkg/config"
)

type Manager struct {
	logger *zap.SugaredLogger

	db         forddb.Database
	rsm        *config.ResourceManager
	routing    *comms.Routing
	supervisor *supervisor2.Manager

	aidCounter atomic.Uint64

	procDir string
}

func NewManager(
	logger *zap.SugaredLogger,
	rsm *config.ResourceManager,
	db forddb.Database,
	routing *comms.Routing,
	sup *supervisor2.Manager,
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

func (m *Manager) StartAgent(ctx context.Context, task *collective2.Agent) error {
	aid := m.aidCounter.Add(1)
	procDir := path.Join(m.procDir, fmt.Sprintf("%d", aid))

	if err := os.MkdirAll(procDir, 0755); err != nil {
		return err
	}

	port, err := m.routing.AddPort(task.Spec.PortID)

	if err != nil {
		return err
	}

	if err := port.Subscribe(task.ID.String()); err != nil {
		return err
	}

	if err := port.Subscribe("aip-bod-room"); err != nil {
		return err
	}

	profile, err := forddb.Get[*collective2.Profile](ctx, m.db, task.Spec.ProfileID)

	if err != nil {
		return err
	}

	tmpProfilePath := path.Join(procDir, "profile.json")
	tmpProfile, err := os.OpenFile(tmpProfilePath, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	if err := SerializeTo(tmpProfile, profile); err != nil {
		return err
	}

	if err := tmpProfile.Close(); err != nil {
		return err
	}

	args := []string{
		"-i", task.ID.String(),
		"-p", tmpProfile.Name(),
	}

	err = os.WriteFile(path.Join(procDir, "args"), []byte(strings.Join(args, "\000")), 0644)

	if err != nil {
		return err
	}

	args = append(args, task.Spec.ExtraArgs...)

	p, err := m.supervisor.Spawn(
		supervisor2.WithID(task.ID.String()),
		supervisor2.WithArgs(args...),
		supervisor2.WithPort(port),
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

func SerializeTo(profile *os.File, profile2 *collective2.Profile) error {
	panic("implement me")
}

func (m *Manager) StopAgent(ctx context.Context, previous *collective2.Agent) error {
	return nil
}
