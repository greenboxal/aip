package generators

import (
	"context"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
	"github.com/greenboxal/aip/aip-langchain/pkg/chain"
	chat2 "github.com/greenboxal/aip/aip-langchain/pkg/llm/chat"
	"github.com/greenboxal/aip/aip-langchain/pkg/memoryctx"
	"github.com/greenboxal/aip/aip-langchain/pkg/providers/openai"
	"github.com/greenboxal/aip/aip-langchain/pkg/tokenizers"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/models"
)

var ImageSettingsKey chain.ContextKey[ImageSettings] = "ImageSettings"
var ImagePromptKey chain.ContextKey[string] = "ImagePrompt"

type ImageSettings struct {
	Prompt string
	Path   string
}

var ImageGeneratorPrompt = chat2.ComposeTemplate(
	chat2.EntryTemplate(
		msn.RoleSystem,
		chain.NewTemplatePrompt(`
You are an AI assistant specialized in generating prompts for images for a Wiki-style satirical content in the voice of {{.PageSettings.Voice}}.
Be as funny as possible but don't use any curse words or aggressive language.

Be as short as possible.
`,
			chain.WithRequiredInput(PageSettingsKey, SiteSettingsKey, ImageSettingsKey),
		),
	),

	chat2.EntryTemplate(
		msn.RoleUser,
		chain.NewTemplatePrompt(
			`Generate a prompt for image for a Wiki style page about "{{.PageSettings.Title}}". The image is named {{.ImageSettings.Path}}.`,
			chain.WithRequiredInput(PageSettingsKey, SiteSettingsKey, ImageSettingsKey),
		),
	),

	chat2.EntryTemplate(msn.RoleAI, chain.Static("")),
)

type ImageGenerator struct {
	client *openai.Client

	model     *openai.ChatLanguageModel
	tokenizer *tokenizers.TikTokenTokenizer

	imageChain chain.Handler
}

func NewImageGenerator(client *openai.Client) (*ImageGenerator, error) {
	w := &ImageGenerator{}
	w.client = client

	w.model = &openai.ChatLanguageModel{
		Client: client,
		Model:  "gpt-3.5-turbo",
	}

	w.tokenizer = tokenizers.TikTokenForModel(openai.AdaEmbeddingV2.String())

	w.imageChain = chain.Sequential(
		chat2.Predict(
			w.model,
			ImageGeneratorPrompt,
			chat2.WithChatMemory(chat2.MemoryContextKey),
			chat2.WithOutputParsers(
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
	cctx.SetInput(chat2.MemoryContextKey, chatMemory)

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
