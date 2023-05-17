package msn

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-forddb/pkg/apis/graphql"
	"github.com/greenboxal/aip/aip-forddb/pkg/apis/rpc"
	api2 "github.com/greenboxal/aip/aip-forddb/pkg/forddb"
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
