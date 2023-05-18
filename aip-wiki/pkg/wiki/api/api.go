package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/samber/lo"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-langchain/pkg/chunkers"
	"github.com/greenboxal/aip/aip-langchain/pkg/providers/openai"
	"github.com/greenboxal/aip/aip-langchain/pkg/tokenizers"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/models"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/sema"
)

type API struct {
	chi.Router

	db forddb.Database

	sctx sema.SemanticContext
}

func NewAPI(db forddb.Database, oai *openai.Client) *API {
	api := &API{}

	api.Router = chi.NewRouter()
	api.db = db

	api.sctx.Tokenizer = tokenizers.TikTokenForModel("gpt-3.5-turbo")
	api.sctx.Chunker = chunkers.TikToken{}

	api.sctx.Embedder = &openai.Embedder{
		Model:  openai.AdaEmbeddingV2,
		Client: oai,
	}

	api.sctx.Model = &openai.ChatLanguageModel{
		Model:  "gpt-3.5-turbo",
		Client: oai,
	}

	sema.InitializeSemanticContext(&api.sctx)

	api.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		res, err := api.HandleAPI(request.Context(), request.URL.Query())

		if err != nil {
			panic(err)
		}

		data, err := json.Marshal(res)

		if err != nil {
			panic(err)
		}

		writer.WriteHeader(http.StatusOK)

		_, _ = writer.Write(data)
	})

	return api
}

func (a *API) HandleAPI(ctx context.Context, query url.Values) (any, error) {
	basePageIds := lo.Map(query["base_page_id"], func(s string, _ int) models.PageID {
		return forddb.NewStringID[models.PageID](s)
	})

	basePages := lo.Map(basePageIds, func(id models.PageID, _ int) *models.Page {
		basePage, err := forddb.Get[*models.Page](ctx, a.db, id)

		if err != nil {
			panic(err)
		}

		return basePage
	})

	baseNodes := lo.Map(basePages, func(p *models.Page, _ int) *sema.SemanticNode {
		node := &sema.SemanticNode{}
		node.ID = forddb.NewStringID[sema.SemanticNodeID](p.ID.String())
		node.Status.Value = a.sctx.Unit(a.sctx.Content(p.Status.Markdown))
		return node
	})

	for _, n := range baseNodes {
		if err := a.sctx.Append(ctx, n); err != nil {
			return nil, err
		}
	}

	rootNode := baseNodes[0]

	result, err := a.sctx.Refine(ctx, rootNode)

	if err != nil {
		return nil, err
	}

	return result, nil
}
