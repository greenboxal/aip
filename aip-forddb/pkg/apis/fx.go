package apis

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/apis/memorylink"
	"github.com/greenboxal/aip/aip-forddb/pkg/apis/graphql"
	"github.com/greenboxal/aip/aip-forddb/pkg/apis/rest"
	"github.com/greenboxal/aip/aip-forddb/pkg/apis/rpc"
	"github.com/greenboxal/aip/aip-forddb/pkg/apis/supervisor"
)

var Module = fx.Module(
	"apis/v1",

	graphql.Module,
	rest.Module,
	rpc.Module,

	rpc.ProvideRpcService[*supervisor.SupervisorAPI](supervisor.NewSupervisorApi, "supervisor"),
	rpc.ProvideRpcService[*memorylink.MemoryLink](memorylink.NewMemoryLink, "memlink"),
)
