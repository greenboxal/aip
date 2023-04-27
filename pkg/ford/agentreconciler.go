package ford

import (
	"context"

	"go.uber.org/zap"

	"github.com/greenboxal/aip/pkg/collective"
	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type AgentReconciler struct {
	logger *zap.SugaredLogger

	manager *Manager

	db forddb.Database
}

func NewAgentReconciler(
	logger *zap.SugaredLogger,
	db forddb.Database,
	manager *Manager,
) *AgentReconciler {
	ar := &AgentReconciler{
		logger:  logger.Named("agent-reconciler"),
		db:      db,
		manager: manager,
	}

	db.AddListener(
		forddb.TypedListenerFunc[collective.AgentID, *collective.Agent](
			func(id collective.AgentID, previous, current *collective.Agent) {
				_, err := ar.Reconcile(context.Background(), previous, current)

				if err != nil {
					logger.Error(err)
				}
			},
		),
	)

	return ar
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
