package daemon

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/supervisor"
)

var Module = fx.Module(
	"daemon",

	fx.Provide(supervisor.NewManager),

	fx.Provide(NewDaemon),
)
