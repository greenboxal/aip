package qdrant

import (
	"context"

	"github.com/greenboxal/aip/aip-langchain/pkg/vectorstore"
)

type Collection struct {
	p    *Provider
	name string
	dim  int
}

func newCollection(p *Provider, name string, dim int) *Collection {
	return &Collection{
		p:    p,
		name: name,
		dim:  dim,
	}
}

func (c *Collection) Dimensions() int {
	return c.dim
}

func (c *Collection) IndexDocument(ctx context.Context, document *vectorstore.Document, options ...vectorstore.IndexDocumentOption) (*vectorstore.IndexedDocument, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Collection) Search(ctx context.Context, query string, options ...vectorstore.SearchDocumentOption) (*vectorstore.SearchResult, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Collection) SimilaritySearch(ctx context.Context, embeddings []float32, opts *vectorstore.SearchDocumentOptions) (*vectorstore.SearchResult, error) {
	//TODO implement me
	panic("implement me")
}
