package reconciliation

import (
	"github.com/jbenet/goprocess"
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/apis/graphql"
)

var Module = fx.Module(
	"reconciliation",

	graphql.ProvideBinding[*reconcilersApi](newReconcilersApi),
)

func WithReconciler[R Reconciler](iface any) fx.Option {
	return fx.Options(
		fx.Provide(iface),

		fx.Invoke(func(r R) {
			goprocess.Go(r.Run)
		}),
	)
}
