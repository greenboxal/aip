package reconciliation

import (
	"context"
	"reflect"

	"github.com/jbenet/goprocess"
	"github.com/modern-go/reflect2"
	"go.uber.org/zap"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type Reconciler interface {
	Run(proc goprocess.Process)
	Close() error
}

type ReconcilerBase[ID forddb.ResourceID[T], T forddb.Resource[ID]] struct {
	logger *zap.SugaredLogger

	db      forddb.Database
	cache   map[string]T
	handler ReconcilerHandler[ID, T]
	stream  *forddb.LogStream

	id           forddb.LogStreamID
	resourceType forddb.BasicResourceID
}

type ReconcilerHandler[ID forddb.ResourceID[T], T forddb.Resource[ID]] func(ctx context.Context, id ID, previous, current T) (T, error)

func NewReconciler[ID forddb.ResourceID[T], T forddb.Resource[ID]](
	logger *zap.SugaredLogger,
	db forddb.Database,
	id forddb.LogStreamID,
	handler ReconcilerHandler[ID, T],
) *ReconcilerBase[ID, T] {
	r := &ReconcilerBase[ID, T]{
		logger:       logger,
		handler:      handler,
		db:           db,
		id:           id,
		resourceType: forddb.TypeSystem().LookupByIDType(reflect.TypeOf((*ID)(nil)).Elem()).GetResourceBasicID(),

		cache: map[string]T{},
	}

	return r
}

func (r *ReconcilerBase[ID, T]) Run(proc goprocess.Process) {
	consumer := &forddb.LogConsumer{
		LogStore: r.db.LogStore(),
		StreamID: r.id,
		Handler: func(ctx context.Context, record *forddb.LogEntryRecord) error {
			var previous, current T

			if record.Type != r.resourceType {
				return nil
			}

			if record.Current != nil {
				current = record.Current.(T)
			}

			if record.Previous != nil {
				previous = record.Previous.(T)
			}

			if reflect2.IsNil(previous) || reflect2.IsNil(current) || current.GetResourceVersion() > previous.GetResourceVersion() {
				_, err := r.handler(ctx, current.GetResourceBasicID().(ID), previous, current)

				if err != nil {
					r.logger.Error(err)
				}
			}

			return nil
		},
	}

	consumer.Run(proc)
}

func (r *ReconcilerBase[ID, T]) Close() error {
	return nil
}
