package forddbimpl

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-forddb/pkg/logstore"
	"github.com/greenboxal/aip/aip-sdk/pkg/config"
)

var Module = fx.Module(
	"forddb/impl",

	fx.Provide(NewDatabase),

	fx.Provide(func(rsm *config.ResourceManager) (forddb.LogStore, error) {
		//path := rsm.GetDataDirectory("log")
		//fss, err := logstore.NewOldFileLogStore(path)

		//if err != nil {
		//	return nil, err
		//}

		//return fss, nil

		return logstore.NewMemoryLogStore(), nil
	}),
)
