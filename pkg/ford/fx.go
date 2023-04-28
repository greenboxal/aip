package ford

import (
	"context"

	"go.uber.org/fx"

	"github.com/greenboxal/aip/pkg/ford/forddb"
)

var Module = fx.Module(
	"ford",

	fx.Provide(NewManager),

	WithReconciler[*TaskReconciler](NewTaskReconciler),
	WithReconciler[*AgentReconciler](NewAgentReconciler),
	WithReconciler[*PortReconciler](NewPortReconciler),

	WithInMemoryDatabase(),
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

func WithInMemoryDatabase() fx.Option {
	return fx.Provide(func() forddb.Database {
		return forddb.NewInMemory()
	})
}
