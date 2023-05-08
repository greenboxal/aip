package reconciliation

import (
	"context"
	"reflect"

	"github.com/modern-go/reflect2"
	"go.uber.org/zap"

	forddb "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type Reconciler interface {
	Run(ctx context.Context)
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

func (r *ReconcilerBase[ID, T]) Run(ctx context.Context) {
	consumer := &forddb.LogConsumer{
		LogStore: r.db.LogStore(),
		StreamID: r.id,
		Handler: func(ctx context.Context, record *forddb.LogEntryRecord) error {
			if record.Type != r.resourceType {
				return nil
			}

			id := forddb.NewStringID[ID](record.ID)
			previous := r.cache[id.String()]

			current, err := forddb.Get[T](ctx, r.db, id)

			if err != nil {
				return err
			}

			if reflect2.IsNil(previous) || reflect2.IsNil(current) || current.GetResourceVersion() > previous.GetResourceVersion() {
				encoded, err := forddb.Encode(current)

				if err != nil {
					return err
				}

				decoded, err := forddb.Decode(encoded)

				if err != nil {
					return err
				}

				_, err = r.handler(ctx, id, previous, decoded.(T))

				if err != nil {
					r.logger.Error(err)
				}
			}

			r.cache[id.String()] = current

			return nil
		},
	}

	if err := consumer.Run(ctx); err != nil {
		panic(err)
	}
}

func (r *ReconcilerBase[ID, T]) Close() error {
	return nil
}
