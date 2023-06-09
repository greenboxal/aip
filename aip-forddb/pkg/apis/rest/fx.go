package rest

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-forddb/pkg/apimachinery"
)

var Module = fx.Module(
	"apis/rest",

	apimachinery.ProvideHttpService[*ResourcesAPI](NewResourcesAPI, "/v1"),
)
