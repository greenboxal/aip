package reconcilers

import (
	"context"

	"go.uber.org/zap"

	"github.com/greenboxal/aip/aip-controller/pkg/collective"
	"github.com/greenboxal/aip/aip-controller/pkg/ford/agent"
	forddb2 "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	"github.com/greenboxal/aip/aip-controller/pkg/ford/reconciliation"
)

type AgentReconciler struct {
	*reconciliation.ReconcilerBase[collective.AgentID, *collective.Agent]

	logger  *zap.SugaredLogger
	manager *agent.Manager
	db      forddb2.Database
}

func NewAgentReconciler(
	logger *zap.SugaredLogger,
	db forddb2.Database,
	manager *agent.Manager,
) *AgentReconciler {
	ar := &AgentReconciler{
		logger: logger.Named("agent-reconciler"),

		db:      db,
		manager: manager,
	}

	ar.ReconcilerBase = reconciliation.NewReconciler(
		ar.logger,
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

				return forddb2.Put(tr.db, current)
			}

			current.Status.State = collective.AgentStateScheduled

			tr.logger.Infow("agent scheduled", "id", current.ID)
		}
	}

	return forddb2.Put(tr.db, current)
}