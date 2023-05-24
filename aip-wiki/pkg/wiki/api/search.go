package api

import (
	"context"
	"errors"

	"github.com/greenboxal/aip/aip-langchain/pkg/vectorstore"
)

type Search struct {
	provider vectorstore.Provider
}

func NewSearch(
	provider vectorstore.Provider,
) *Search {
	return &Search{
		provider: provider,
	}
}

type SearchRequest struct {
	Collections []string `json:"collections"`
	MaxHits     int      `json:"max_hits"`

	ReturnHitContents   bool `json:"return_hit_contents"`
	ReturnHitEmbeddings bool `json:"return_hit_embeddings"`

	Query string `json:"query"`
}

type SearchResponse struct {
	Hits []vectorstore.SearchHit `json:"hits"`
}

func (sa *Search) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	var options []vectorstore.SearchDocumentOption

	if len(req.Collections) != 1 {
		return nil, errors.New("multi-shard search is not supported yet")
	}

	collectionName := req.Collections[0]
	col := sa.provider.Collection(collectionName, 1536)

	if req.ReturnHitEmbeddings {
		options = append(options, vectorstore.WithReturnHitEmbeddings())
	}

	if req.ReturnHitContents {
		options = append(options, vectorstore.WithReturnHitContents())
	}

	result, err := col.Search(ctx, req.Query, options...)

	if err != nil {
		return nil, err
	}

	return &SearchResponse{
		Hits: result.Hits,
	}, nil
}
