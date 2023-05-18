package api

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-forddb/pkg/apimachinery"
)

var Module = fx.Module(
	"api",

	apimachinery.ProvideHttpService[*API](NewAPI, "/v1/test"),
)
