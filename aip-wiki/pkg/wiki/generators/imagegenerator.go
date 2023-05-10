package generators

import (
	"context"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/chain"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/chat"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/memoryctx"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/providers/openai"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/tokenizers"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/models"
)

var ImageSettingsKey chain.ContextKey[ImageSettings] = "ImageSettings"
var ImagePromptKey chain.ContextKey[string] = "ImagePrompt"

type ImageSettings struct {
	Prompt string
	Path   string
}

var ImageGeneratorPrompt = chat.ComposeTemplate(
	chat.EntryTemplate(
		msn.RoleSystem,
		chain.NewTemplatePrompt(`
You are an AI assistant specialized in generating prompts for images for a Wiki-style satirical content in the voice of {{.PageSettings.Voice}}.
Be as funny as possible but don't use any curse words or aggressive language.

Be as short as possible.
`,
			chain.WithRequiredInput(PageSettingsKey, SiteSettingsKey, ImageSettingsKey),
		),
	),

	chat.EntryTemplate(
		msn.RoleUser,
		chain.NewTemplatePrompt(
			`Generate a prompt for image for a Wiki style page about "{{.PageSettings.Title}}". The image is named {{.ImageSettings.Path}}.`,
			chain.WithRequiredInput(PageSettingsKey, SiteSettingsKey, ImageSettingsKey),
		),
	),

	chat.EntryTemplate(msn.RoleAI, chain.Static("")),
)

type ImageGenerator struct {
	client *openai.Client

	model     *openai.ChatLanguageModel
	tokenizer *tokenizers.TikTokenTokenizer

	imageChain chain.Chain
}

func NewImageGenerator(client *openai.Client) (*ImageGenerator, error) {
	var err error

	w := &ImageGenerator{}
	w.client = client

	w.model = &openai.ChatLanguageModel{
		Client: client,
		Model:  "gpt-3.5-turbo",
	}

	w.tokenizer, err = tokenizers.TikTokenForModel(openai.AdaEmbeddingV2.String())

	if err != nil {
		return nil, err
	}

	w.imageChain = chain.Compose(
		chat.Predict(
			w.model,
			ImageGeneratorPrompt,
			chat.WithChatMemory(chat.MemoryContextKey),
			chat.WithOutputParsers(
				GeneratedHtmlParser(ImagePromptKey),
			),
		),
	)

	return w, nil
}

func (ig *ImageGenerator) GetImage(
	ctx context.Context,
	imageSettings models.ImageSpec,
) (models.ImageStatus, error) {
	// FIXME: ?
	pageSettings := models.PageSpec{}
	siteSettings := SiteSettings{}

	cctx := chain.NewChainContext(ctx)

	chatMemory := memoryctx.GetMemory(ctx)

	cctx.SetInput(SiteSettingsKey, siteSettings)
	cctx.SetInput(PageSettingsKey, pageSettings)
	cctx.SetInput(ImageSettingsKey, imageSettings)
	cctx.SetInput(chat.MemoryContextKey, chatMemory)

	if err := ig.imageChain.Run(cctx); err != nil {
		return models.ImageStatus{}, err
	}

	prompt := chain.Output(cctx, ImagePromptKey)

	result, err := ig.client.CreateImage(ctx, openai.ImageRequest{
		N:              1,
		Size:           "1024x1024",
		ResponseFormat: "url",
		Prompt:         prompt,
	})

	if err != nil {
		return models.ImageStatus{}, err
	}

	return models.ImageStatus{
		URL:    result.Data[0].URL,
		Prompt: prompt,
	}, nil
}
