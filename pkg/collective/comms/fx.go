package comms

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/pkg/collective/comms/transports/pubsub"
	"github.com/greenboxal/aip/pkg/collective/comms/transports/slack"
)

var Module = fx.Module(
	"comms",

	fx.Provide(NewRouting),
	fx.Provide(NewManager),

	fx.Provide(slack.NewTransport),
	fx.Provide(pubsub.NewTransport),
)
