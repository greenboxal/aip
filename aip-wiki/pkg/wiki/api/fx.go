package api

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-forddb/pkg/apimachinery"
	"github.com/greenboxal/aip/aip-forddb/pkg/apis/rpc"
)

var Module = fx.Module(
	"api",

	rpc.ProvideRpcService[*Search](NewSearch, "search"),

	apimachinery.ProvideHttpService[*API](NewAPI, "/v1/test"),
)
