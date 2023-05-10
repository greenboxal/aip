package jobs

import (
	"context"

	"github.com/google/uuid"
	"github.com/jbenet/goprocess"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type Manager struct {
	logger     *zap.SugaredLogger
	db         forddb.Database
	reconciler *Reconciler
	supervisor *Supervisor
}

func NewManager(
	lc fx.Lifecycle,
	logger *zap.SugaredLogger,
	db forddb.Database,
	supervisor *Supervisor,
	reconciler *Reconciler,
) *Manager {
	m := &Manager{
		logger:     logger.Named("job-manager"),
		db:         db,
		supervisor: supervisor,
		reconciler: reconciler,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return m.Start(ctx)
		},
	})

	return m
}

func (m *Manager) DispatchEphemeral(ctx context.Context, spec JobSpec) (JobHandle, error) {
	job := &Job{}

	job.ID = forddb.NewStringID[JobID](uuid.New().String())
	job.Spec = spec
	job.Status.State = JobStateScheduled

	job, err := forddb.Put(ctx, m.db, job)

	if err != nil {
		return nil, err
	}

	return m.supervisor.CheckJobState(job)
}

func (m *Manager) Start(ctx context.Context) error {
	goprocess.Go(m.reconciler.Run)

	return nil
}
