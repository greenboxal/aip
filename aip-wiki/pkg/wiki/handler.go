package wiki

import (
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Router struct {
	chi.Router

	wiki *Wiki
}

func NewRouter(wiki *Wiki) *Router {
	r := &Router{}
	r.wiki = wiki
	r.Router = chi.NewRouter()

	r.NotFound(r.handlePage)

	return r
}

var SanitizePageSlugRegex = regexp.MustCompile(`[^a-zA-Z0-9\-_]`)

func (r *Router) handlePage(writer http.ResponseWriter, request *http.Request) {
	var pageSettings PageSettings
	var siteSettings SiteSettings
	var imageSettings ImageSettings

	url := request.URL
	host := request.Header.Get("Host")

	siteSettings.Title = "BolsoWiki"
	siteSettings.BaseUrl = "http://127.0.0.1:30100/wiki"

	pageSettings.Title = SanitizePageSlugRegex.ReplaceAllString(url.Path, " ")
	pageSettings.Voice = url.Query().Get("voice")
	pageSettings.Language = url.Query().Get("lang")

	if pageSettings.Voice == "" {
		pageSettings.Voice = "Jair Bolsonaro"
	}

	if strings.HasSuffix(host, "desciclo.ai") {
		if pageSettings.Language == "" {
			pageSettings.Language = "Portuguese"
		}
	} else if strings.HasSuffix(host, "uncyclo.ai") {
		if pageSettings.Language == "" {
			pageSettings.Language = "English"
		}
	} else {
		if pageSettings.Language == "" {
			pageSettings.Language = "English"
		}
	}

	extension := path.Ext(url.Path)
	isImage := false

	switch extension {
	case ".txt":
		pageSettings.Format = "text/plain"
	case ".md":
		pageSettings.Format = "text/markdown"
	case ".html":
		pageSettings.Format = "text/html"
	case ".jpg", ".png", ".gif", ".svg", ".ico", ".webp", ".bmp", ".tiff", ".tif":
		pageSettings.Format = "HTML"
		isImage = true
	default:
		if url.Query().Has("format") {
			pageSettings.Format = url.Query().Get("format")
		} else {
			pageSettings.Format = "text/html"
		}
	}

	if isImage {
		imageSettings.Path = url.Path

		imageUrl, err := r.wiki.GetImage(request.Context(), siteSettings, pageSettings, imageSettings)

		if err != nil {
			panic(err)
		}

		writer.Header().Set("Location", imageUrl)
		writer.WriteHeader(http.StatusMovedPermanently)
	} else {
		pageContents, err := r.wiki.GetPage(request.Context(), siteSettings, pageSettings)

		if err != nil {
			panic(err)
		}

		writer.Header().Set("Content-Type", pageSettings.Format)
		writer.WriteHeader(http.StatusOK)
		writer.Write(pageContents)
	}
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}
