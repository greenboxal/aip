package comms

import (
	"context"
	"os"

	"go.uber.org/fx"

	transports2 "github.com/greenboxal/aip/aip-controller/pkg/collective/comms/transports"
	"github.com/greenboxal/aip/aip-controller/pkg/collective/comms/transports/local"
)

type Routing struct {
	*local.Transport
}

func NewRouting(
	lc fx.Lifecycle,
) *Routing {
	gw := transports2.Tee(
		&transports2.StdioTransport{Stdout: os.Stdout},

		//slack,
	)

	r := &Routing{
		Transport: local.NewTransport(gw),
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return r.Close()
		},
	})

	return r
}
