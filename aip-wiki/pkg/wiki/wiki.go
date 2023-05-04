package wiki

import (
	"context"

	"github.com/greenboxal/aip/aip-controller/pkg/llm/chain"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/chat"
	openai "github.com/greenboxal/aip/aip-controller/pkg/llm/providers/openai"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/tokenizers"
)

func NewWiki(client *openai.Client) (*Wiki, error) {
	var err error

	w := &Wiki{}
	w.client = client

	w.model = &openai.ChatLanguageModel{
		Client: client,
		Model:  "gpt-3.5-turbo",
	}

	w.tokenizer, err = tokenizers.TikTokenForModel(openai.AdaEmbeddingV2.String())

	if err != nil {
		return nil, err
	}

	w.layoutChain = chain.Compose(
		chat.Predict(w.model, LayoutGeneratorPrompt, GoTemplateParser(PageLayoutKey)),
	)

	w.contentChain = chain.Compose(
		chat.Predict(w.model, PageGeneratorPrompt, GeneratedHtmlParser(PageContentKey)),
	)

	w.imageChain = chain.Compose(
		chat.Predict(w.model, ImageGeneratorPrompt, GeneratedHtmlParser(ImagePromptKey)),
		chat.Predict(w.model, AnonymizerPrompt, GeneratedHtmlParser(ImagePromptKey)),
	)

	return w, nil
}

type Wiki struct {
	client *openai.Client

	model     *openai.ChatLanguageModel
	tokenizer *tokenizers.TikTokenTokenizer

	layoutChain  chain.Chain
	contentChain chain.Chain
	imageChain   chain.Chain
}

func (w *Wiki) GetPage(ctx context.Context, siteSettings SiteSettings, pageSettings PageSettings) ([]byte, error) {
	cctx := chain.NewChainContext(ctx)

	cctx.SetInput(SiteSettingsKey, siteSettings)
	cctx.SetInput(PageSettingsKey, pageSettings)

	if err := w.contentChain.Run(cctx); err != nil {
		return nil, err
	}

	pageContent := chain.Output(cctx, PageContentKey)

	return []byte(pageContent), nil
}

func (w *Wiki) GetImage(ctx context.Context, siteSettings SiteSettings, pageSettings PageSettings, imageSettings ImageSettings) (string, error) {
	cctx := chain.NewChainContext(ctx)

	cctx.SetInput(SiteSettingsKey, siteSettings)
	cctx.SetInput(PageSettingsKey, pageSettings)
	cctx.SetInput(ImageSettingsKey, imageSettings)

	if err := w.imageChain.Run(cctx); err != nil {
		return "", err
	}

	prompt := chain.Output(cctx, ImagePromptKey)

	result, err := w.client.CreateImage(ctx, openai.ImageRequest{
		N:              1,
		Size:           "1024x1024",
		ResponseFormat: "url",
		Prompt:         prompt,
	})

	if err != nil {
		return "", err
	}

	return result.Data[0].URL, nil
}
