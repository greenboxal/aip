package ford

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/pkg/ford/agent"
	"github.com/greenboxal/aip/pkg/ford/forddb"
	"github.com/greenboxal/aip/pkg/ford/reconcilers"
	"github.com/greenboxal/aip/pkg/ford/reconciliation"
)

var Module = fx.Module(
	"ford",

	fx.Provide(agent.NewManager),

	reconciliation.WithReconciler[*reconcilers.TaskReconciler](reconcilers.NewTaskReconciler),
	reconciliation.WithReconciler[*reconcilers.AgentReconciler](reconcilers.NewAgentReconciler),
	reconciliation.WithReconciler[*reconcilers.PortReconciler](reconcilers.NewPortReconciler),

	WithInMemoryDatabase(),
)

func WithInMemoryDatabase() fx.Option {
	return fx.Provide(func() forddb.Database {
		return forddb.NewInMemory()
	})
}
