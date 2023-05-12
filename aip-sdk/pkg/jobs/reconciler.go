package jobs

import (
	"context"

	"go.uber.org/zap"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/reconciliation"
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type Reconciler struct {
	*reconciliation.ReconcilerBase[JobID, *Job]

	logger     *zap.SugaredLogger
	supervisor *Supervisor
}

func NewReconciler(
	logger *zap.SugaredLogger,
	db forddb.Database,
	supervisor *Supervisor,
) *Reconciler {
	r := &Reconciler{}
	r.logger = logger.Named("job-reconciler")
	r.supervisor = supervisor

	r.ReconcilerBase = reconciliation.NewReconciler[JobID, *Job](
		r.logger,
		db,
		"job-reconciler",
		r.reconcile,
	)

	return r
}

func (r *Reconciler) reconcile(ctx context.Context, id JobID, previous *Job, current *Job) (*Job, error) {
	if _, err := r.supervisor.CheckJobState(current); err != nil {
		return nil, err
	}

	return current, nil
}
