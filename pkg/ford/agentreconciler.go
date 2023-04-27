package ford

import (
	"context"

	"github.com/greenboxal/aip/pkg/collective"
	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type AgentReconciler struct {
	manager *Manager

	db forddb.Database
}

func NewAgentReconciler(
	db forddb.Database,
	manager *Manager,
) *AgentReconciler {
	return &AgentReconciler{
		db:      db,
		manager: manager,
	}
}

func (tr *AgentReconciler) Reconcile(ctx context.Context, previous, current *collective.Agent) (*collective.Agent, error) {
	if current == nil && previous != nil {
		if err := tr.manager.StopAgent(ctx, previous); err != nil {
			return nil, err
		}

		return nil, nil
	}

	if current.Status.State == "" {
		current.Status.State = collective.AgentStateCreated
	}

	if previous.Status.State != current.Status.State {
		switch current.Status.State {
		case collective.AgentStateCreated:
			current.Status.State = collective.AgentStatePending

		case collective.AgentStatePending:
			if err := tr.manager.StartAgent(ctx, current); err != nil {
				current.Status.State = collective.AgentStateFailed
				current.Status.LastError = err.Error()

				return forddb.CreateOrUpdate(tr.db, current)
			}

			current.Status.State = collective.AgentStateScheduled
		}
	}

	return forddb.CreateOrUpdate(tr.db, current)
}
