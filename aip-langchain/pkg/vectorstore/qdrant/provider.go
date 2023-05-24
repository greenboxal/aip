package qdrant

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/greenboxal/aip/aip-langchain/pkg/vectorstore"
)

type Config struct {
	Endpoint string `mapstructure:"endpoint"`
}

type Provider struct {
	logger *zap.SugaredLogger

	config Config
}

func NewProvider(logger *zap.SugaredLogger, lc fx.Lifecycle, config *Config) *Provider {
	p := &Provider{
		logger: logger,
		config: *config,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return nil
		},

		OnStop: func(ctx context.Context) error {
			return p.Close()
		},
	})

	return p
}

func (p *Provider) Collection(name string, dim int) vectorstore.Collection {
	return newCollection(p, name, dim)
}

func (p *Provider) Close() error {
	return nil
}
