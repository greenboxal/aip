package ford

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/pkg/config"
	"github.com/greenboxal/aip/pkg/ford/agent"
	"github.com/greenboxal/aip/pkg/ford/forddb"
	forddbimpl "github.com/greenboxal/aip/pkg/ford/forddb/impl"
	"github.com/greenboxal/aip/pkg/ford/reconcilers"
	"github.com/greenboxal/aip/pkg/ford/reconciliation"
	"github.com/greenboxal/aip/pkg/indexing"
)

var Module = fx.Module(
	"ford",

	fx.Provide(NewManager),
	fx.Provide(agent.NewManager),
	fx.Provide(forddbimpl.NewDatabase),
	fx.Provide(forddbimpl.NewFileLogStore),

	fx.Provide(func(rsm *config.ResourceManager) (forddb.LogStore, error) {
		//path := rsm.GetDataDirectory("log")
		//fss, err := forddbimpl.NewFileLogStore(path)

		//if err != nil {
		//	return nil, err
		//}

		//return fss, nil
	
		return forddbimpl.NewMemoryLogStore(), nil
	}),

	fx.Provide(func(m *Manager) indexing.Index {
		return m.Index()
	}),

	reconciliation.WithReconciler[*reconcilers.TaskReconciler](reconcilers.NewTaskReconciler),
	reconciliation.WithReconciler[*reconcilers.AgentReconciler](reconcilers.NewAgentReconciler),
	reconciliation.WithReconciler[*reconcilers.PortReconciler](reconcilers.NewPortReconciler),
)
