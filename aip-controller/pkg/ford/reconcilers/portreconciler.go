package reconcilers

import (
	"context"

	"go.uber.org/zap"

	"github.com/greenboxal/aip/aip-controller/pkg/collective"
	"github.com/greenboxal/aip/aip-controller/pkg/ford/reconciliation"
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type PortReconciler struct {
	*reconciliation.ReconcilerBase[collective.PortID, *collective.Port]

	logger *zap.SugaredLogger
	db     forddb.Database
}

func NewPortReconciler(
	logger *zap.SugaredLogger,
	db forddb.Database,
) *PortReconciler {
	pr := &PortReconciler{
		logger: logger.Named("port-reconciler"),

		db: db,
	}

	pr.ReconcilerBase = reconciliation.NewReconciler(
		pr.logger,
		db,
		"port-reconciler",
		pr.Reconcile,
	)

	return pr
}

func (pr *PortReconciler) Reconcile(
	ctx context.Context,
	id collective.PortID,
	previous *collective.Port,
	current *collective.Port,
) (*collective.Port, error) {
	var err error

	if current == nil && previous != nil {
		pr.logger.Info("port deleted", "task_id", current.ID)

		return nil, nil
	}

	pr.logger.Info("entering reconciliation loop", "port_id", current.ID)

	if current.Status.State == "" {
		current.Status.State = collective.PortCreated
	}

	if previous == nil || previous.Status.State != current.Status.State {
		current, err = forddb.Put(ctx, pr.db, current)

		if err != nil {
			return nil, err
		}
	}

	return current, nil
}
