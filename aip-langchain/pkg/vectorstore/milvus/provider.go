package milvus

import (
	"go.uber.org/zap"

	"github.com/greenboxal/aip/aip-langchain/pkg/vectorstore"
)

type Provider struct {
	logger *zap.SugaredLogger

	milvus *Milvus
}

func NewProvider(logger *zap.SugaredLogger, milvus *Milvus) *Provider {
	return &Provider{
		logger: logger.Named("milvus-index-provider"),
		milvus: milvus,
	}
}

func (p *Provider) Collection(name string, dim int) vectorstore.Collection {
	return newCollection(p.logger, p.milvus, name, dim)
}

func (p *Provider) Close() error {
	return nil
}

type Collection struct {
	logger *zap.SugaredLogger
	milvus *Milvus
	name   string
	dim    int
}
