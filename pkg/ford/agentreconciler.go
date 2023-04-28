package ford

import (
	"context"

	"go.uber.org/zap"

	"github.com/greenboxal/aip/pkg/collective"
	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type AgentReconciler struct {
	*ReconcilerBase[collective.AgentID, *collective.Agent]

	manager *Manager
	db      forddb.Database
}

func NewAgentReconciler(
	logger *zap.SugaredLogger,
	db forddb.Database,
	manager *Manager,
) *AgentReconciler {
	ar := &AgentReconciler{

		db:      db,
		manager: manager,
	}

	ar.ReconcilerBase = NewReconciler(
		logger.Named("agent-reconciler"),
		db,
		ar.Reconcile,
	)

	return ar
}

func (tr *AgentReconciler) Reconcile(
	ctx context.Context,
	id collective.AgentID,
	previous *collective.Agent,
	current *collective.Agent,
) (*collective.Agent, error) {
	if current == nil && previous != nil {
		tr.logger.Infow("deleting agent", "id", previous.ID)

		if err := tr.manager.StopAgent(ctx, previous); err != nil {
			return nil, err
		}

		return nil, nil
	}

	if current.Status.State == "" {
		current.Status.State = collective.AgentStateCreated
	}

	if previous == nil || previous.Status.State != current.Status.State {
		switch current.Status.State {
		case collective.AgentStateCreated:
			current.Status.State = collective.AgentStatePending

			tr.logger.Infow("agent created", "id", current.ID)

		case collective.AgentStatePending:
			if err := tr.manager.StartAgent(ctx, current); err != nil {
				current.Status.State = collective.AgentStateFailed
				current.Status.LastError = err.Error()

				return forddb.Put(tr.db, current)
			}

			current.Status.State = collective.AgentStateScheduled

			tr.logger.Infow("agent scheduled", "id", current.ID)
		}
	}

	return forddb.Put(tr.db, current)
}
