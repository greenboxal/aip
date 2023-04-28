package ford

import (
	"context"

	"go.uber.org/zap"

	"github.com/greenboxal/aip/pkg/collective"
	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type PortReconciler struct {
	db     forddb.Database
	logger *zap.SugaredLogger

	cache map[collective.PortID]*collective.Port
}

func NewPortReconciler(
	logger *zap.SugaredLogger,
	db forddb.Database,
) *PortReconciler {
	pr := &PortReconciler{
		logger: logger.Named("port-reconciler"),

		db: db,

		cache: map[collective.PortID]*collective.Port{},
	}

	watchCh := make(chan collective.PortID, 128)

	db.AddListener(
		forddb.TypedListenerFunc[collective.PortID, *collective.Port](
			func(id collective.PortID, previous, current *collective.Port) {
				watchCh <- id
			},
		),
	)

	go func() {
		for id := range watchCh {
			previous := pr.cache[id]

			current, err := forddb.Get[*collective.Port](db, id)

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

				_, err = pr.Reconcile(context.Background(), previous, decoded.(*collective.Port))

				if err != nil {
					logger.Error(err)
				}
			}

			pr.cache[id] = current
		}
	}()

	return pr
}

func (pr *PortReconciler) Reconcile(
	ctx context.Context,
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
		current, err = forddb.CreateOrUpdate(pr.db, current)

		if err != nil {
			return nil, err
		}
	}

	return current, nil
}
