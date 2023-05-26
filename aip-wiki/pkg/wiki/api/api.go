package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/go-github/v52/github"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-langchain/pkg/chunkers"
	"github.com/greenboxal/aip/aip-langchain/pkg/providers/openai"
	"github.com/greenboxal/aip/aip-langchain/pkg/tokenizers"
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

	api.Post("/", func(writer http.ResponseWriter, request *http.Request) {
		res, err := api.HandleAPI(request)

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

func (a *API) HandleAPI(req *http.Request) (any, error) {
	var payload github.PushEvent

	body, err := io.ReadAll(req.Body)
	// body.commmits[0].message

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	commit := &Commit{Message: *payload.Commits[0].Message}
	commit.ID = forddb.NewStringID[CommitID](*payload.Commits[0].ID)

	fmt.Fprintf(os.Stderr, "Adding %#v to forddb\n", commit)

	commit, err = forddb.Put(req.Context(), a.db, commit)

	if err != nil {
		return nil, err
	}

	return "Hello world", nil
}

type CommitID struct {
	forddb.StringResourceID[*Commit] `ipld:",inline"`
}

type Commit struct {
	forddb.ResourceBase[CommitID, *Commit] `json:"metadata"`

	Message string `json:"message"`
}

func init() {
	forddb.DefineResourceType[CommitID, *Commit]("commit")
}
