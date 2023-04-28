package ford

import (
	"context"

	"go.uber.org/zap"

	"github.com/greenboxal/aip/pkg/collective"
	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type PortReconciler struct {
	*ReconcilerBase[collective.PortID, *collective.Port]

	db forddb.Database
}

func NewPortReconciler(
	logger *zap.SugaredLogger,
	db forddb.Database,
) *PortReconciler {
	pr := &PortReconciler{
		db: db,
	}

	pr.ReconcilerBase = NewReconciler(
		logger.Named("port-reconciler"),
		db,
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
		current, err = forddb.Put(pr.db, current)

		if err != nil {
			return nil, err
		}
	}

	return current, nil
}
