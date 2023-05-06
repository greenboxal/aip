package reconciliation

import (
	"context"

	"github.com/modern-go/reflect2"
	"go.uber.org/zap"

	forddb2 "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type Reconciler interface {
	Run(ctx context.Context)
	Enqueue(id forddb2.BasicResourceID)
	AsListener() forddb2.Listener
	Close() error
}

type ReconcilerBase[ID forddb2.ResourceID[T], T forddb2.Resource[ID]] struct {
	logger  *zap.SugaredLogger
	db      forddb2.Database
	cache   map[string]T
	handler ReconcilerHandler[ID, T]
	ch      chan ID
}

type ReconcilerHandler[ID forddb2.ResourceID[T], T forddb2.Resource[ID]] func(ctx context.Context, id ID, previous, current T) (T, error)

func NewReconciler[ID forddb2.ResourceID[T], T forddb2.Resource[ID]](
	logger *zap.SugaredLogger,
	db forddb2.Database,
	handler ReconcilerHandler[ID, T],
) *ReconcilerBase[ID, T] {
	r := &ReconcilerBase[ID, T]{
		logger:  logger,
		handler: handler,
		db:      db,

		cache: map[string]T{},
		ch:    make(chan ID, 128),
	}

	return r
}

func (r *ReconcilerBase[ID, T]) Run(ctx context.Context) {
	for id := range r.ch {
		previous := r.cache[id.String()]

		current, err := forddb2.Get[T](r.db, id)

		if err != nil {
			r.logger.Error(err)
			continue
		}

		if reflect2.IsNil(previous) || reflect2.IsNil(current) || current.GetResourceVersion() > previous.GetResourceVersion() {
			encoded, err := forddb2.Encode(current)

			if err != nil {
				r.logger.Error(err)
				continue
			}

			decoded, err := forddb2.Decode(encoded)

			if err != nil {
				r.logger.Error(err)
				continue
			}

			_, err = r.handler(ctx, id, previous, decoded.(T))

			if err != nil {
				r.logger.Error(err)
			}
		}

		r.cache[id.String()] = current
	}
}

func (r *ReconcilerBase[ID, T]) AsListener() forddb2.Listener {
	return forddb2.TypedListenerFunc[ID, T](func(id ID, previous, current T) {
		r.Enqueue(id)
	})
}

func (r *ReconcilerBase[ID, T]) Enqueue(id forddb2.BasicResourceID) {
	tid, ok := id.(ID)

	if !ok {
		return
	}

	r.ch <- tid
}

func (r *ReconcilerBase[ID, T]) Close() error {
	return nil
}
