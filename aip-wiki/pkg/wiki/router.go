package wiki

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/cms"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/models"
	"github.com/greenboxal/aip/aip-wiki/public"
)

type Router2 struct {
	chi.Router

	db forddb.Database
}

func NewRouter2(wiki *Wiki, pm *cms.PageManager) *Router2 {
	r := &Router2{}
	r.Router = chi.NewRouter()

	assets := http.FileServer(http.FS(public.Content()))

	r.Handle("/static/*", http.StripPrefix("static", assets))

	r.Get("/*", func(writer http.ResponseWriter, request *http.Request) {
		err := r.handleRequest(writer, request)

		if err != nil {
			panic(err)
		}
	})

	return r
}

func (r *Router2) handleRequest(writer http.ResponseWriter, request *http.Request) error {
	domain, err := r.getDomain(request)

	if err != nil {
		return err
	}

	route, err := r.getRoute(request, domain)

	if err != nil {
		return err
	}

	page, err := r.getPage(request, domain, route)

	if err != nil {
		return err
	}

	// FIXME
	writer.WriteHeader(200)
	_, _ = writer.Write([]byte(page.Status.HTML))

	return nil
}

func (r *Router2) getDomain(request *http.Request) (*models.Domain, error) {
	id := forddb.NewStringID[models.DomainID](request.Host)

	return forddb.Get[*models.Domain](request.Context(), r.db, id)
}

func (r *Router2) getRoute(request *http.Request, domain *models.Domain) (*models.RouteBinding, error) {
	result, err := forddb.List[*models.RouteBinding](
		request.Context(),
		r.db,
		models.RouteBindingType,
		forddb.WithListQueryOptions(
			forddb.WithFilterExpression("resource.domain_id == args.domain_id"),
			forddb.WithFilterParameter("domain_id", domain.ID),
		),
	)

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, forddb.ErrNotFound
	}

	return result[0], nil
}

func (r *Router2) getPage(
	request *http.Request,
	domain *models.Domain,
	route *models.RouteBinding,
) (*models.Page, error) {
	return nil, nil
}
