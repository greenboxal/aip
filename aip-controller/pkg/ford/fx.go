package ford

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/agent"
	reconcilers2 "github.com/greenboxal/aip/aip-controller/pkg/ford/reconcilers"
	"github.com/greenboxal/aip/aip-controller/pkg/ford/reconciliation"
	"github.com/greenboxal/aip/aip-langchain/pkg/indexing"
)

var Module = fx.Module(
	"ford",

	reconciliation.Module,

	fx.Provide(NewManager),
	fx.Provide(agent.NewManager),

	fx.Provide(func(m *Manager) indexing.Index {
		return m.Index()
	}),

	reconciliation.WithReconciler[*reconcilers2.TaskReconciler](reconcilers2.NewTaskReconciler),
	reconciliation.WithReconciler[*reconcilers2.AgentReconciler](reconcilers2.NewAgentReconciler),
	reconciliation.WithReconciler[*reconcilers2.PortReconciler](reconcilers2.NewPortReconciler),
)
