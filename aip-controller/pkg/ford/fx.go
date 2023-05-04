package ford

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/config"
	"github.com/greenboxal/aip/aip-controller/pkg/ford/agent"
	forddbimpl2 "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb/impl"
	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb/logstore"
	reconcilers2 "github.com/greenboxal/aip/aip-controller/pkg/ford/reconcilers"
	"github.com/greenboxal/aip/aip-controller/pkg/ford/reconciliation"
	"github.com/greenboxal/aip/aip-controller/pkg/indexing"
)

var Module = fx.Module(
	"ford",

	fx.Provide(NewManager),
	fx.Provide(agent.NewManager),
	fx.Provide(forddbimpl2.NewDatabase),
	fx.Provide(forddbimpl2.NewFileLogStore),

	fx.Provide(func(rsm *config.ResourceManager) (logstore.LogStore, error) {
		//path := rsm.GetDataDirectory("log")
		//fss, err := forddbimpl.NewFileLogStore(path)

		//if err != nil {
		//	return nil, err
		//}

		//return fss, nil

		return forddbimpl2.NewMemoryLogStore(), nil
	}),

	fx.Provide(func(m *Manager) indexing.Index {
		return m.Index()
	}),

	reconciliation.WithReconciler[*reconcilers2.TaskReconciler](reconcilers2.NewTaskReconciler),
	reconciliation.WithReconciler[*reconcilers2.AgentReconciler](reconcilers2.NewAgentReconciler),
	reconciliation.WithReconciler[*reconcilers2.PortReconciler](reconcilers2.NewPortReconciler),
)