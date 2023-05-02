package apimachinery

import (
	"net/http"
	"net/http/pprof"

	"github.com/arl/statsviz"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type RootMux struct {
	chi.Router
}

func NewRootMux() *RootMux {
	mux := &RootMux{
		Router: chi.NewRouter(),
	}

	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(cors.AllowAll().Handler)
	mux.Use(middleware.AllowContentEncoding("gzip"))
	mux.Use(middleware.Heartbeat("/ping"))

	statsVizMux := http.NewServeMux()
	if err := statsviz.Register(statsVizMux); err != nil {
		panic(err)
	}

	mux.Mount("/metrics", promhttp.Handler())
	mux.Mount("/debug/statsviz", statsVizMux)

	mux.Get("/debug/pprof/", pprof.Index)
	mux.Get("/debug/pprof/cmdline", pprof.Cmdline)
	mux.Get("/debug/pprof/profile", pprof.Profile)
	mux.Get("/debug/pprof/symbol", pprof.Symbol)
	mux.Get("/debug/pprof/trace", pprof.Trace)

	return mux
}
