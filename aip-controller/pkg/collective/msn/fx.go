package msn

import (
	"go.uber.org/fx"

	api2 "github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-sdk/pkg/apis/graphql"
	"github.com/greenboxal/aip/aip-sdk/pkg/apis/rpc"
)

var Module = fx.Module(
	"collective/msn",

	fx.Provide(NewRouter),
	fx.Provide(NewRealTimeService),
	fx.Provide(NewService),
	fx.Provide(newEventDispatcher),
	fx.Provide(newApi),

	rpc.BindRpcService[*Service]("msn-router"),

	graphql.WithBinding[*RealTimeService](),

	fx.Invoke(func(dispatcher *eventDispatcher) {}),
)

func init() {
	api2.DefineResourceType[MessageID, *Message]("message")
	api2.DefineResourceType[EndpointID, *Endpoint]("endpoint")
	api2.DefineResourceType[ChannelID, *Channel]("channel")
}
