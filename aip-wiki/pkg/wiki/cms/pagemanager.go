package cms

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/multiformats/go-multihash"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/generators"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/models"
)

const PageGeneratorBaseUrl = "http://127.0.0.1:30100"
const PageGeneratorWikiUrl = "http://127.0.0.1:30100"

type PageManager struct {
	db forddb.Database

	fm    *FileManager
	cache *generators.ContentCache

	pageGenerator  *generators.PageGenerator
	imageGenerator *generators.ImageGenerator
}

func NewPageManager(
	db forddb.Database,
	fm *FileManager,
	cache *generators.ContentCache,
	pageGenerator *generators.PageGenerator,
	imageGenerator *generators.ImageGenerator,
) *PageManager {
	return &PageManager{
		db:             db,
		fm:             fm,
		cache:          cache,
		pageGenerator:  pageGenerator,
		imageGenerator: imageGenerator,
	}
}

func (pm *PageManager) GetPageByID(ctx context.Context, id models.PageID) (*models.Page, error) {
	return forddb.Get[*models.Page](ctx, pm.db, id)
}

func (pm *PageManager) GetPage(ctx context.Context, spec models.PageSpec) (*models.Page, error) {
	page, err := pm.cache.GetPage(ctx, spec)

	if forddb.IsNotFound(err) {
		return pm.GeneratePage(ctx, spec)
	} else if err != nil {
		return nil, err
	}

	return page, nil
}

func (pm *PageManager) GetImage(ctx context.Context, spec models.ImageSpec) (*models.Image, error) {
	page, err := pm.cache.GetImage(ctx, spec)

	if forddb.IsNotFound(err) {
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

	return pm.cache.PutPage(ctx, page)
}

func (pm *PageManager) GenerateImage(ctx context.Context, spec models.ImageSpec) (*models.Image, error) {
	id := models.BuildImageID(spec)

	status, err := pm.imageGenerator.GetImage(ctx, spec)

	if err != nil {
		return nil, err
	}

	response, err := http.Get(status.URL)

	if err != nil {
		return nil, err
	}

	tempFileName := fmt.Sprintf("temp-%s-%d", id.String(), time.Now().UnixNano())

	writer := pm.fm.OpenWriter(ctx, tempFileName)
	reader := io.TeeReader(response.Body, writer)

	h, err := multihash.SumStream(reader, multihash.SHA2_256, -1)

	if err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	fileName := "images/" + h.B58String() + path.Ext(status.URL)

	if err := pm.fm.Rename(ctx, tempFileName, fileName); err != nil {
		return nil, err
	}

	status.URL = "https://cdn.desciclo.ai/" + fileName

	result := &models.Image{
		Spec:   spec,
		Status: status,
	}

	result.ID = id

	return pm.cache.PutImage(ctx, result)
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
