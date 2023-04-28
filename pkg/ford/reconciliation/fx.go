package reconciliation

import (
	"context"

	"go.uber.org/fx"

	"github.com/greenboxal/aip/pkg/ford/forddb"
)

func WithReconciler[R Reconciler](iface any) fx.Option {
	return fx.Options(
		fx.Provide(iface),

		fx.Invoke(func(db forddb.Database, r R) {
			db.AddListener(r.AsListener())

			go r.Run(context.Background())
		}),
	)
}
