package wiki

import (
	"github.com/greenboxal/aip/aip-controller/pkg/llm/chain"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/chat"
)

var ImageSettingsKey chain.ContextKey[ImageSettings] = "ImageSettings"
var ImagePromptKey chain.ContextKey[string] = "ImagePrompt"

type ImageSettings struct {
	Prompt string
	Path   string
}

var ImageGeneratorPrompt = chat.ComposeTemplate(
	chat.EntryTemplate(
		chat.RoleSystem,
		chain.NewTemplatePrompt(`
You are an AI assistant specialized in generating prompts for images for a Wiki-style satirical content in the voice of {{.PageSettings.Voice}}.
Be as funny as possible but don't use any curse words or aggressive language.

Be as short as possible.
`,
			chain.WithRequiredInput(PageSettingsKey, SiteSettingsKey, ImageSettingsKey),
		),
	),

	chat.EntryTemplate(
		chat.RoleUser,
		chain.NewTemplatePrompt(
			`Generate a prompt for image for a Wiki style page about "{{.PageSettings.Title}}". The image is named {{.ImageSettings.Path}}.`,
			chain.WithRequiredInput(PageSettingsKey, SiteSettingsKey, ImageSettingsKey),
		),
	),

	chat.EntryTemplate(chat.RoleAI, chain.Static("")),
)
