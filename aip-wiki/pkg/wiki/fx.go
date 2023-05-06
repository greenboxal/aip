package wiki

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/apimachinery"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/cms"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/generators"
)

var Module = fx.Module(
	"wiki",

	fx.Provide(NewWiki),
	fx.Provide(cms.NewPageManager),
	fx.Provide(cms.NewContentCache),
	fx.Provide(generators.NewImageGenerator),
	fx.Provide(generators.NewPageGenerator),

	apimachinery.ProvideHttpService[*Router](NewRouter, "/", apimachinery.WithStripPrefix()),
)
