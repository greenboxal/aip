package wiki

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-forddb/pkg/apimachinery"
	"github.com/greenboxal/aip/aip-forddb/pkg/apis/rpc"
	"github.com/greenboxal/aip/aip-sdk/pkg/config"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/api"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/cms"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/generators"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/indexer"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/jobs"
)

var Module = fx.Module(
	"wiki",

	config.RegisterConfig[cms.FileManagerConfig]("wiki.filemanager"),

	jobs.Module,
	api.Module,

	fx.Provide(NewWiki),
	fx.Provide(cms.NewPageManager),
	fx.Provide(cms.NewFileManager),
	fx.Provide(cms.NewService),
	fx.Provide(indexer.NewPageIndexer),
	fx.Provide(indexer.NewCommitIndexer),
	fx.Provide(generators.NewContentCache),
	fx.Provide(generators.NewImageGenerator),
	fx.Provide(generators.NewPageGenerator),

	rpc.BindRpcService[*cms.Service]("wiki"),
	rpc.BindRpcService[*cms.PageManager]("wikiPageManager"),
	rpc.BindRpcService[*generators.ContentCache]("wikiContentCache"),
	rpc.BindRpcService[*cms.FileManager]("wikiFileManager"),
	rpc.BindRpcService[*generators.ImageGenerator]("wikiImageGenerator"),
	rpc.BindRpcService[*generators.PageGenerator]("wikiPageGenerator"),

	apimachinery.ProvideHttpService[*Router](NewRouter, "/", apimachinery.WithStripPrefix()),
)
