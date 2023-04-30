package milvus

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/pkg/indexing"
)

var Module = fx.Module(
	"milvus",

	fx.Provide(NewStorage),
)

func WithMilvusStorage() fx.Option {
	return fx.Options(
		Module,

		fx.Provide(func(s *Storage) indexing.Storage {
			return s
		}),
	)
}
