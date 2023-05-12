package generators

import (
	"html/template"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
	chain2 "github.com/greenboxal/aip/aip-langchain/pkg/llm/chain"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/chat"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/memory"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/models"
)

var BasePageKey chain2.ContextKey[*models.Page] = "BasePage"
var PageSettingsKey chain2.ContextKey[models.PageSpec] = "PageSettings"
var SiteSettingsKey chain2.ContextKey[SiteSettings] = "SiteSettings"
var PageLayoutKey chain2.ContextKey[*template.Template] = "PageLayout"
var PageContentKey chain2.ContextKey[string] = "PageContent"
var AttentionContextKey chain2.ContextKey[string] = "AttentionContext"

type SiteSettings struct {
	Title       string
	Description string
	BaseUrl     string
}

type PageSettings struct {
	Format      string
	Layout      string
	Title       string
	Description string
	Voice       string
	Language    string
}

var PageGeneratorHeader = chat.EntryTemplate(
	msn.RoleSystem,
	chain2.NewTemplatePrompt(`
You are an AI assistant specialized in generating Wiki-style satirical content in the voice of {{.PageSettings.Voice}}.
Be as expressive as possible. Use as many curse words as you can. Be as funny as you can.

You should follow the following rules:

* Output valid Markdown
* Base URL is {{.SiteSettings.BaseUrl}}
* Add real references and citations in the footer
* Add inline links to other pages within the text of article. Be sure to link every occurrence of any term that deserves its own page.
* Add at least one image. Add alt tags to all images with a proper description. All image URLs should include the description of the image.
* Make it pretty, using modern design and typography
* Make it funny. Use satire, irony, exaggeration, ridicule, or the like, in exposing, denouncing, or deriding vice, folly, etc.
* Write it in {{.PageSettings.Language}}.
`,
		chain2.WithRequiredInput(PageSettingsKey, SiteSettingsKey),
	),
)

var PageGeneratorPrompt = chat.ComposeTemplate(
	PageGeneratorHeader,

	/*chat.EntryTemplate(
			chat.RoleSystem,
			chain.NewTemplatePrompt(`
	Attention Context: {{.AttentionContext}}
	`, chain.WithRequiredInput(AttentionContextKey)),
		),*/

	chat.HistoryFromContext(memory.ContextualMemoryKey),

	chat.EntryTemplate(
		msn.RoleUser,
		chain2.NewTemplatePrompt(
			`Generate a Wiki style page about "{{.PageSettings.Title}}".`,
			chain2.WithRequiredInput(PageSettingsKey),
		),
	),

	chat.EntryTemplate(msn.RoleAI, chain2.Static("")),
)

var PageEditorPrompt = chat.ComposeTemplate(
	PageGeneratorHeader,

	chat.HistoryFromContext(memory.ContextualMemoryKey),

	chat.EntryTemplate(
		msn.RoleUser,
		chain2.NewTemplatePrompt(
			`Improve the Wiki page by expanding the topic about "{{.PageSettings.Title}}, and correcting anything you believe is wrong.
Remember to leave notes at the bottom of the page explaining what you did and why. Write everything in {{.PageSettings.Language}}.`,
			chain2.WithRequiredInput(PageSettingsKey),
		),
	),

	chat.EntryTemplate(msn.RoleAI, chain2.Static("")),
)

var LinkEnricherPrompt = chat.ComposeTemplate(
	chat.EntryTemplate(
		msn.RoleSystem,
		chain2.NewTemplatePrompt(`
You are an AI assistant specialized in analyzing the most relevant topics in an article and adding links to the HTML page for the article.
Return the full HTML page with links added.
`,
			chain2.WithRequiredInput(PageSettingsKey, SiteSettingsKey, PageContentKey),
		),
	),

	chat.EntryTemplate(
		msn.RoleUser,
		chain2.NewTemplatePrompt(
			`Insert links in the HTML below that was generated for a Wiki style page about "{{.PageSettings.Title}}"
{{.PageContent}}"`,
			chain2.WithRequiredInput(PageSettingsKey, SiteSettingsKey, PageContentKey),
		),
	),

	chat.EntryTemplate(msn.RoleAI, chain2.Static("")),
)
