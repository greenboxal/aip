package wiki

import (
	"context"

	"github.com/jbenet/goprocess"
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/llm/providers/openai"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/tokenizers"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/indexer"
)

func NewWiki(
	lc fx.Lifecycle,
	client *openai.Client,
	pageIndexer *indexer.PageIndexer,
) (*Wiki, error) {
	var err error

	w := &Wiki{}
	w.client = client
	w.pageIndexer = pageIndexer

	w.model = &openai.ChatLanguageModel{
		Client: client,
		Model:  "gpt-3.5-turbo",
	}

	w.tokenizer, err = tokenizers.TikTokenForModel(openai.AdaEmbeddingV2.String())

	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return w.Start(ctx)
		},
	})

	return w, nil
}

type Wiki struct {
	client *openai.Client

	model     *openai.ChatLanguageModel
	tokenizer *tokenizers.TikTokenTokenizer

	pageIndexer *indexer.PageIndexer
}

func (w *Wiki) Start(ctx context.Context) error {
	goprocess.Go(w.pageIndexer.Run)

	return nil
}
