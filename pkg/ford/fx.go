package ford

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/pkg/ford/forddb"
)

var Module = fx.Module(
	"ford",

	fx.Provide(NewDatabase),
	fx.Provide(NewManager),
	fx.Provide(NewTaskReconciler),
	fx.Provide(NewAgentReconciler),

	fx.Invoke(func(tr *TaskReconciler, ar *AgentReconciler) {}),
)

func NewDatabase() forddb.Database {
	return forddb.NewInMemory()
}
