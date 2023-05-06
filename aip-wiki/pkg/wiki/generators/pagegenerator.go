package generators

import (
	"context"

	"github.com/greenboxal/aip/aip-controller/pkg/llm/chain"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/chat"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/providers/openai"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/tokenizers"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/models"
)

type PageGenerator struct {
	client *openai.Client

	model     *openai.ChatLanguageModel
	tokenizer *tokenizers.TikTokenTokenizer

	contentChain chain.Chain
}

func NewPageGenerator(client *openai.Client) (*PageGenerator, error) {
	var err error

	w := &PageGenerator{}
	w.client = client

	w.model = &openai.ChatLanguageModel{
		Client: client,
		Model:  "gpt-3.5-turbo",
	}

	w.tokenizer, err = tokenizers.TikTokenForModel(openai.AdaEmbeddingV2.String())

	if err != nil {
		return nil, err
	}

	w.contentChain = chain.Compose(
		chat.Predict(w.model, PageGeneratorPrompt, GeneratedHtmlParser(PageContentKey)),
	)

	return w, nil
}

func (pg *PageGenerator) GetPage(
	ctx context.Context,
	pageSettings models.PageSpec,
) ([]byte, error) {
	cctx := chain.NewChainContext(ctx)

	cctx.SetInput(SiteSettingsKey, SiteSettings{
		BaseUrl: "http://127.0.0.1:30100",
	})

	cctx.SetInput(PageSettingsKey, pageSettings)

	if err := pg.contentChain.Run(cctx); err != nil {
		return nil, err
	}

	pageContent := chain.Output(cctx, PageContentKey)

	return []byte(pageContent), nil
}
