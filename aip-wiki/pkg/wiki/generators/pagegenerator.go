package generators

import (
	"context"
	"net/url"
	"path"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	chain2 "github.com/greenboxal/aip/aip-langchain/pkg/llm/chain"
	chat2 "github.com/greenboxal/aip/aip-langchain/pkg/llm/chat"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/memory"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/memoryctx"
	openai2 "github.com/greenboxal/aip/aip-langchain/pkg/llm/providers/openai"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/tokenizers"
	"github.com/greenboxal/aip/aip-langchain/pkg/vectorstore"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/models"
)

const PageGeneratorBaseUrl = "http://127.0.0.1:30100"
const PageGeneratorWikiUrl = "http://127.0.0.1:30100"

type PageGenerator struct {
	client *openai2.Client

	cache     *ContentCache
	model     *openai2.ChatLanguageModel
	tokenizer *tokenizers.TikTokenTokenizer

	contentChain chain2.Handler
	editorChain  chain2.Handler
}

func NewPageGenerator(
	client *openai2.Client,
	cache *ContentCache,
	index vectorstore.VectorStore,
	oai *openai2.Client,
) (*PageGenerator, error) {
	var err error

	w := &PageGenerator{}
	w.client = client
	w.cache = cache

	w.model = &openai2.ChatLanguageModel{
		Client: client,
		Model:  "gpt-3.5-turbo",
	}

	w.tokenizer, err = tokenizers.TikTokenForModel(openai2.AdaEmbeddingV2.String())

	if err != nil {
		return nil, err
	}

	contextualMemory := &memory.ContextualMemory{
		HistoryKey: chat2.MemoryContextKey,
		InputKey:   chat2.ChatReplyContextKey,
		ContextKey: memory.ContextualMemoryKey,

		Index: index,

		Embedder: &openai2.Embedder{
			Client: oai,
			Model:  openai2.AdaEmbeddingV2,
		},
	}

	w.contentChain = chain2.Sequential(
		chain2.Func(func(ctx chain2.ChainContext) error {
			page := chain2.Input(ctx, PageSettingsKey)

			return contextualMemory.LoadFor(ctx, page.Title)
		}),

		chat2.Predict(
			w.model,
			PageGeneratorPrompt,
			chat2.WithChatMemory(chat2.MemoryContextKey),
			chat2.WithOutputParsers(
				GeneratedHtmlParser(PageContentKey),
			),
		),
	)

	w.editorChain = chain2.Sequential(
		chain2.Func(func(ctx chain2.ChainContext) error {
			page := chain2.Input(ctx, PageSettingsKey)

			return contextualMemory.LoadFor(ctx, page.Title)
		}),

		chat2.Predict(
			w.model,
			PageEditorPrompt,
			chat2.WithChatMemory(chat2.MemoryContextKey),
			chat2.WithOutputParsers(
				GeneratedHtmlParser(PageContentKey),
			),
		),
	)

	return w, nil
}

func (pg *PageGenerator) GetPage(
	ctx context.Context,
	pageSettings models.PageSpec,
) ([]byte, error) {
	siteSettings := SiteSettings{
		BaseUrl: "http://127.0.0.1:30100",
	}

	chatMemory := memoryctx.GetMemory(ctx)

	cctx := chain2.NewChainContext(ctx)

	cctx.SetInput(SiteSettingsKey, siteSettings)
	cctx.SetInput(PageSettingsKey, pageSettings)
	cctx.SetInput(chat2.MemoryContextKey, chatMemory)

	if pageSettings.BasePage.IsEmpty() {
		if err := pg.contentChain.Run(cctx); err != nil {
			return nil, err
		}

		pageContent := chain2.Output(cctx, PageContentKey)

		return []byte(pageContent), nil
	} else {
		basePage, err := pg.cache.GetPageByID(ctx, pageSettings.BasePage)

		if err != nil {
			return nil, err
		}

		cctx.SetInput(BasePageKey, basePage)

		if err := pg.editorChain.Run(cctx); err != nil {
			return nil, err
		}

		pageContent := chain2.Output(cctx, PageContentKey)

		return []byte(pageContent), nil
	}
}

func (pg *PageGenerator) GeneratePage(ctx context.Context, spec models.PageSpec) (*models.Page, error) {
	id := models.BuildPageID(spec)

	page := &models.Page{}
	page.ID = id
	page.Spec = spec

	body, err := pg.GetPage(ctx, spec)

	if err != nil {
		return nil, err
	}

	md := ParseMarkdown(body)

	ast.WalkFunc(md, func(node ast.Node, entering bool) ast.WalkStatus {
		switch n := node.(type) {
		case *ast.Link:
			link := models.PageLink{
				Title: string(n.Title),
				To:    string(n.Destination),
			}

			if strings.HasPrefix(link.To, PageGeneratorWikiUrl) {
				link.To = strings.TrimPrefix(link.To, PageGeneratorWikiUrl)
			}

			if strings.HasPrefix(link.To, PageGeneratorBaseUrl) {
				link.To = strings.TrimPrefix(link.To, PageGeneratorBaseUrl)
			}

			page.Status.Links = append(page.Status.Links, link)

			n.Destination = []byte(link.To)

		case *ast.Image:
			image := models.PageImage{
				Title:  string(n.Title),
				Source: string(n.Destination),
			}

			if strings.HasPrefix(image.Source, PageGeneratorBaseUrl) {
				image.Source = strings.TrimPrefix(image.Source, PageGeneratorBaseUrl)
				image.Source = path.Join("/images/"+url.QueryEscape(image.Title), image.Source)
			}

			page.Status.Images = append(page.Status.Images, image)

			n.Destination = []byte(image.Source)
		}

		return ast.GoToNext
	})

	page.Status.Markdown = string(body)
	page.Status.HTML = string(RenderMarkdownToHtml(md))

	return pg.cache.PutPage(ctx, page)
}

func ParseMarkdown(md []byte) ast.Node {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)

	return p.Parse(md)
}

func RenderMarkdownToHtml(doc ast.Node) []byte {
	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}
