package wiki

import (
	"html/template"

	"github.com/greenboxal/aip/aip-controller/pkg/llm/chain"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/chat"
)

var PageSettingsKey chain.ContextKey[PageSettings] = "PageSettings"
var SiteSettingsKey chain.ContextKey[SiteSettings] = "SiteSettings"
var PageLayoutKey chain.ContextKey[*template.Template] = "PageLayout"
var PageContentKey chain.ContextKey[string] = "PageContent"

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

var PageGeneratorPrompt = chat.ComposeTemplate(
	chat.EntryTemplate(
		chat.RoleSystem,
		chain.NewTemplatePrompt(`
You are an AI assistant specialized in generating Wiki-style satirical content in the voice of {{.PageSettings.Voice}}.
Be as expressive as possible. Use as many curse words as you can. Be as funny as you can.

You should follow the following rules:

* Output {{.PageSettings.Format}}.
* Base URL is {{.SiteSettings.BaseUrl}}
* Add real references and citations in the footer
* Add inline links to other pages within the text of article. Be sure to link every occurrence of any term that deserves its own page.
* Add at least one image. Add alt tags to all images with a proper description. All image URLs should include the description of the image.
* Make it pretty, using modern design and typography
* Make it funny. Use satire, irony, exaggeration, ridicule, or the like, in exposing, denouncing, or deriding vice, folly, etc.
* Write it in {{.PageSettings.Language}}.
`,
			chain.WithRequiredInput(PageSettingsKey, SiteSettingsKey),
		),
	),

	chat.EntryTemplate(
		chat.RoleUser,
		chain.NewTemplatePrompt(
			`Generate a Wiki style page about "{{.PageSettings.Title}}", including HTML tags.`,
			chain.WithRequiredInput(PageSettingsKey, SiteSettingsKey),
		),
	),

	chat.EntryTemplate(chat.RoleAI, chain.Static("")),
)

var LinkEnricherPrompt = chat.ComposeTemplate(
	chat.EntryTemplate(
		chat.RoleSystem,
		chain.NewTemplatePrompt(`
You are an AI assistant specialized in analyzing the most relevant topics in an article and adding links to the HTML page for the article.
Return the full HTML page with links added.
`,
			chain.WithRequiredInput(PageSettingsKey, SiteSettingsKey, PageContentKey),
		),
	),

	chat.EntryTemplate(
		chat.RoleUser,
		chain.NewTemplatePrompt(
			`Insert links in the HTML below that was generated for a Wiki style page about "{{.PageSettings.Title}}"
{{.PageContent}}"`,
			chain.WithRequiredInput(PageSettingsKey, SiteSettingsKey, PageContentKey),
		),
	),

	chat.EntryTemplate(chat.RoleAI, chain.Static("")),
)

func GeneratedHtmlParser(key chain.ContextKey[string]) chain.OutputParser {
	return chain.OutputParserFunc(func(ctx chain.ChainContext, result string) error {
		ctx.SetOutput(key, result)

		return nil
	})
}

func GoTemplateParser(key chain.ContextKey[*template.Template]) chain.OutputParser {
	return chain.OutputParserFunc(func(ctx chain.ChainContext, result string) error {
		t, err := template.New("template").Parse(result)

		if err != nil {
			return err
		}

		ctx.SetOutput(key, t)

		return nil
	})
}
