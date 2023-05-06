package wiki

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/apimachinery"
	"github.com/greenboxal/aip/aip-controller/pkg/config"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/cms"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/generators"
)

var Module = fx.Module(
	"wiki",

	config.RegisterConfig[cms.FileManagerConfig]("wiki.filemanager"),

	fx.Provide(NewWiki),
	fx.Provide(cms.NewPageManager),
	fx.Provide(generators.NewContentCache),
	fx.Provide(cms.NewFileManager),
	fx.Provide(cms.NewService),
	fx.Provide(generators.NewImageGenerator),
	fx.Provide(generators.NewPageGenerator),

	apimachinery.BindRpcService[*cms.Service]("wiki"),
	apimachinery.BindRpcService[*cms.PageManager]("wikiPageManager"),
	apimachinery.BindRpcService[*generators.ContentCache]("wikiContentCache"),
	apimachinery.BindRpcService[*cms.FileManager]("wikiFileManager"),
	apimachinery.BindRpcService[*generators.ImageGenerator]("wikiImageGenerator"),
	apimachinery.BindRpcService[*generators.PageGenerator]("wikiPageGenerator"),

	apimachinery.ProvideHttpService[*Router](NewRouter, "/", apimachinery.WithStripPrefix()),
)
