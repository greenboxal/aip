package apis

import (
	"go.uber.org/fx"

	apimachinery2 "github.com/greenboxal/aip/aip-controller/pkg/apimachinery"
	"github.com/greenboxal/aip/aip-controller/pkg/apis/graphql"
	"github.com/greenboxal/aip/aip-controller/pkg/apis/memorylink"
	"github.com/greenboxal/aip/aip-controller/pkg/apis/rest"
	"github.com/greenboxal/aip/aip-controller/pkg/apis/supervisor"
)

var Module = fx.Module(
	"apis/v1",

	apimachinery2.ProvideHttpService[*apimachinery2.RpcService](apimachinery2.NewRpcService, "/v1/rpc"),
	apimachinery2.ProvideHttpService[*apimachinery2.Docs](apimachinery2.NewDocs, "/v1/docs"),

	apimachinery2.ProvideHttpService[*graphql.GraphQL](graphql.NewGraphQL, "/v1/graphql"),
	apimachinery2.ProvideHttpService[*rest.ResourcesAPI](rest.NewResourcesAPI, "/v1"),

	apimachinery2.ProvideRpcService[*supervisor.SupervisorAPI](supervisor.NewSupervisorApi, "supervisor"),
	apimachinery2.ProvideRpcService[*memorylink.MemoryLink](memorylink.NewMemoryLink, "memlink"),
)
