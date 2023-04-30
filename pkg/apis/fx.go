package apis

import (
	"go.uber.org/fx"

	apimachinery "github.com/greenboxal/aip/pkg/apimachinery"
	"github.com/greenboxal/aip/pkg/apis/graphql"
	"github.com/greenboxal/aip/pkg/apis/memorylink"
	"github.com/greenboxal/aip/pkg/apis/rest"
	"github.com/greenboxal/aip/pkg/apis/supervisor"
)

var Module = fx.Module(
	"apis/v1",

	apimachinery.ProvideHttpService[*graphql.GraphQL](graphql.NewGraphQL, "/v1/graphql"),
	apimachinery.ProvideHttpService[*rest.ResourcesAPI](rest.NewResourcesAPI, "/v1"),

	apimachinery.ProvideRpcService[*supervisor.SupervisorAPI](supervisor.NewSupervisorApi, "supervisor"),
	apimachinery.ProvideRpcService[*memorylink.MemoryLink](memorylink.NewMemoryLink, "memlink"),
)
