package generators

import (
	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/chain"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/chat"
)

var AnonymizerPrompt = chat.ComposeTemplate(
	chat.EntryTemplate(
		msn.RoleSystem,
		chain.NewTemplatePrompt(`
You are an AI assistant specialized in anonymizing generated content so it can pass the content guidelines for generating images.
`,
			chain.WithRequiredInput(PageSettingsKey, SiteSettingsKey, ImageSettingsKey, ImagePromptKey),
		),
	),

	chat.EntryTemplate(
		msn.RoleUser,
		chain.NewTemplatePrompt(
			`Anonymize the below content that was generated for a Wiki style page about "{{.PageSettings.Title}}
{{.ImagePrompt}}"`,
			chain.WithRequiredInput(PageSettingsKey, SiteSettingsKey, ImageSettingsKey, ImagePromptKey),
		),
	),

	chat.EntryTemplate(msn.RoleAI, chain.Static("")),
)
