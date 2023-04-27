package daemon

import (
	"context"
	"os"

	"go.uber.org/fx"

	"github.com/greenboxal/aip/pkg/collective/transports"
	"github.com/greenboxal/aip/pkg/collective/transports/local"
	"github.com/greenboxal/aip/pkg/collective/transports/pubsub"
	"github.com/greenboxal/aip/pkg/collective/transports/slack"
)

type Routing struct {
	*local.Transport
}

func NewRouting(
	lc fx.Lifecycle,
	slack *slack.Transport,
	pubsub *pubsub.Transport,
) *Routing {
	gw := transports.Tee(
		&transports.StdioTransport{Stdout: os.Stdout},

		slack,
		pubsub,
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
