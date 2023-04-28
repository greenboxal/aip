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

	db    forddb.Database
	cache map[collective.AgentID]*collective.Agent
}

func NewAgentReconciler(
	logger *zap.SugaredLogger,
	db forddb.Database,
	manager *Manager,
) *AgentReconciler {
	ar := &AgentReconciler{
		logger: logger.Named("agent-reconciler"),

		db:      db,
		manager: manager,

		cache: map[collective.AgentID]*collective.Agent{},
	}

	watchCh := make(chan collective.AgentID, 128)

	db.AddListener(
		forddb.TypedListenerFunc[collective.AgentID, *collective.Agent](
			func(id collective.AgentID, previous, current *collective.Agent) {
				watchCh <- id
			},
		),
	)

	go func() {
		for id := range watchCh {
			previous := ar.cache[id]

			current, err := forddb.Get[*collective.Agent](db, id)

			if err != nil {
				logger.Error(err)
				continue
			}

			if previous == nil || current == nil || current.GetVersion() > previous.GetVersion() {
				encoded, err := forddb.Encode(current)

				if err != nil {
					logger.Error(err)
					continue
				}

				decoded, err := forddb.Decode(encoded)

				if err != nil {
					logger.Error(err)
					continue
				}

				_, err = ar.Reconcile(context.Background(), previous, decoded.(*collective.Agent))

				if err != nil {
					logger.Error(err)
				}
			}

			ar.cache[id] = current
		}
	}()

	return ar
}

func (tr *AgentReconciler) Reconcile(ctx context.Context, previous, current *collective.Agent) (*collective.Agent, error) {
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

				return forddb.CreateOrUpdate(tr.db, current)
			}

			current.Status.State = collective.AgentStateScheduled

			tr.logger.Infow("agent scheduled", "id", current.ID)
		}
	}

	return forddb.CreateOrUpdate(tr.db, current)
}
