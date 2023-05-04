package wiki

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/apimachinery"
)

var Module = fx.Module(
	"wiki",

	fx.Provide(NewWiki),

	apimachinery.ProvideHttpService[*Router](NewRouter, "/", apimachinery.WithStripPrefix()),
)
