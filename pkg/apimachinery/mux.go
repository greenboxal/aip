package apimachinery

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	return mux
}
