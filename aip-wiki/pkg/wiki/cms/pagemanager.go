package cms

import (
	"context"
	"net/url"
	"path"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/generators"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/models"
)

const PageGeneratorBaseUrl = "http://127.0.0.1:30100"

type PageManager struct {
	db    forddb.Database
	cache *ContentCache

	pageGenerator  *generators.PageGenerator
	imageGenerator *generators.ImageGenerator
}

func NewPageManager(
	db forddb.Database,
	cache *ContentCache,
	pageGenerator *generators.PageGenerator,
	imageGenerator *generators.ImageGenerator,
) *PageManager {
	return &PageManager{
		db:             db,
		cache:          cache,
		pageGenerator:  pageGenerator,
		imageGenerator: imageGenerator,
	}
}

func (pm *PageManager) GetPage(ctx context.Context, spec models.PageSpec) (*models.Page, error) {
	page, err := pm.cache.GetPage(ctx, spec)

	if err == forddb.ErrNotFound {
		return pm.GeneratePage(ctx, spec)
	} else if err != nil {
		return nil, err
	}

	return page, nil
}

func (pm *PageManager) GetImage(ctx context.Context, spec models.ImageSpec) (*models.Image, error) {
	page, err := pm.cache.GetImage(ctx, spec)

	if err == forddb.ErrNotFound {
		return pm.GenerateImage(ctx, spec)
	} else if err != nil {
		return nil, err
	}

	return page, nil
}

func (pm *PageManager) GeneratePage(ctx context.Context, spec models.PageSpec) (*models.Page, error) {
	id := models.BuildPageID(spec)

	page := &models.Page{}
	page.ID = id
	page.Spec = spec

	body, err := pm.pageGenerator.GetPage(ctx, spec)

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

	return pm.cache.PutPage(ctx, page)
}

func (pm *PageManager) GenerateImage(ctx context.Context, spec models.ImageSpec) (*models.Image, error) {
	id := models.BuildImageID(spec)

	status, err := pm.imageGenerator.GetImage(ctx, spec)

	if err != nil {
		return nil, err
	}

	image := &models.Image{
		Spec:   spec,
		Status: status,
	}

	image.ID = id

	return pm.cache.PutImage(ctx, image)
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
