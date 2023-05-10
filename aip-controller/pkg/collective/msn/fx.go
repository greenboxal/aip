package msn

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/apis/graphql"
	"github.com/greenboxal/aip/aip-controller/pkg/apis/rpc"
	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

var Module = fx.Module(
	"collective/msn",

	fx.Provide(NewRouter),
	fx.Provide(NewRealTimeService),
	fx.Provide(NewService),
	fx.Provide(newEventDispatcher),

	rpc.BindRpcService[*Service]("msn-router"),

	graphql.WithBinding[*RealTimeService](),

	fx.Invoke(func(dispatcher *eventDispatcher) {}),
)

func init() {
	forddb.DefineResourceType[MessageID, *Message]("message")
	forddb.DefineResourceType[EndpointID, *Endpoint]("endpoint")
	forddb.DefineResourceType[ChannelID, *Channel]("channel")
}
