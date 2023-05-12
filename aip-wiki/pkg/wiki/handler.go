package wiki

import (
	"html/template"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/cms"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/generators"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/models"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/templates"
	"github.com/greenboxal/aip/aip-wiki/public"
)

type Router struct {
	chi.Router

	pm   *cms.PageManager
	wiki *Wiki
}

func NewRouter(wiki *Wiki, pm *cms.PageManager) *Router {
	r := &Router{}
	r.pm = pm
	r.wiki = wiki
	r.Router = chi.NewRouter()

	assets := http.FileServer(http.FS(public.Content()))

	r.Handle("/static/*", http.StripPrefix("static", assets))
	r.Get("/*", r.handle)

	return r
}

var SanitizePageSlugRegex = regexp.MustCompile(`[^a-zA-Z0-9\-_]`)

func (r *Router) getPageSettings(request *http.Request) (models.PageSpec, error) {
	var pageSettings models.PageSpec
	var siteSettings generators.SiteSettings

	url := request.URL
	host := request.Header.Get("Host")

	if idx := strings.LastIndex(host, ":"); idx != -1 {
		host = host[:idx]
	}

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

	if basePageId := url.Query().Get("basePageId"); basePageId != "" {
		pageSettings.BasePage = forddb.NewStringID[models.PageID](basePageId)
	}

	return pageSettings, nil
}

func (r *Router) handle(writer http.ResponseWriter, request *http.Request) {
	requestUrl := request.URL
	extension := path.Ext(requestUrl.Path)
	isImage := false

	pageSpec, err := r.getPageSettings(request)

	if err != nil {
		panic(err)
	}

	switch extension {
	case ".txt":
		pageSpec.Format = "text/plain"
	case ".md":
		pageSpec.Format = "text/markdown"
	case ".html":
		pageSpec.Format = "text/html"
	case ".jpg", ".png", ".gif", ".svg", ".ico", ".webp", ".bmp", ".tiff", ".tif":
		pageSpec.Format = "HTML"
		isImage = true
	default:
		if requestUrl.Query().Has("format") {
			pageSpec.Format = requestUrl.Query().Get("format")
		} else {
			pageSpec.Format = "text/html"
		}
	}

	if isImage {
		spec := models.ImageSpec{
			Path: path.Join(requestUrl.Path, requestUrl.RawQuery),
		}

		imageUrl, err := r.pm.GetImage(request.Context(), spec)

		if err != nil {
			panic(err)
		}

		writer.Header().Set("Location", imageUrl.Status.URL)
		writer.WriteHeader(http.StatusMovedPermanently)
	} else {
		var page *models.Page

		isCanonicalPath := false
		pathComponents := strings.Split(requestUrl.Path, "/")

		if len(pathComponents) >= 2 && pathComponents[0] == "wiki" {
			pageIdStr := pathComponents[1]
			pageId := forddb.NewStringID[models.PageID](pageIdStr)

			p, err := r.pm.GetPageByID(request.Context(), pageId)

			if err != nil {
				panic(err)
			}

			page = p
			isCanonicalPath = true
		}

		if page == nil {
			page, err = r.pm.GetPage(request.Context(), pageSpec)

			if err != nil {
				panic(err)
			}

			isCanonicalPath = false
		}

		writer.Header().Set("Content-Type", pageSpec.Format+";charset=UTF-8")

		if isCanonicalPath {
			writer.WriteHeader(http.StatusOK)
		} else {
			canonicalPath := "/wiki/" + page.ID.String() + "/" + url.PathEscape(page.Spec.Title)

			writer.Header().Set("Location", canonicalPath)
			writer.WriteHeader(http.StatusTemporaryRedirect)
		}

		err = templates.Templates().ExecuteTemplate(writer, "index.tmpl.html", map[string]interface{}{
			"Page":     page,
			"PageHTML": template.HTML(page.Status.HTML),
		})

		if err != nil {
			panic(err)
		}
	}
}
