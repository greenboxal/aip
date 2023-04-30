package ford

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/pkg/ford/agent"
	"github.com/greenboxal/aip/pkg/ford/reconcilers"
	"github.com/greenboxal/aip/pkg/ford/reconciliation"
	"github.com/greenboxal/aip/pkg/indexing"
)

var Module = fx.Module(
	"ford",

	fx.Provide(NewManager),
	fx.Provide(agent.NewManager),

	fx.Provide(func(m *Manager) indexing.Index {
		return m.Index()
	}),

	reconciliation.WithReconciler[*reconcilers.TaskReconciler](reconcilers.NewTaskReconciler),
	reconciliation.WithReconciler[*reconcilers.AgentReconciler](reconcilers.NewAgentReconciler),
	reconciliation.WithReconciler[*reconcilers.PortReconciler](reconcilers.NewPortReconciler),
)
