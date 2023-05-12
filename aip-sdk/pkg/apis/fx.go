package apis

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/apis/memorylink"
	"github.com/greenboxal/aip/aip-sdk/pkg/apis/graphql"
	"github.com/greenboxal/aip/aip-sdk/pkg/apis/rest"
	rpc2 "github.com/greenboxal/aip/aip-sdk/pkg/apis/rpc"
	"github.com/greenboxal/aip/aip-sdk/pkg/apis/supervisor"
)

var Module = fx.Module(
	"apis/v1",

	graphql.Module,
	rest.Module,
	rpc2.Module,

	rpc2.ProvideRpcService[*supervisor.SupervisorAPI](supervisor.NewSupervisorApi, "supervisor"),
	rpc2.ProvideRpcService[*memorylink.MemoryLink](memorylink.NewMemoryLink, "memlink"),
)
