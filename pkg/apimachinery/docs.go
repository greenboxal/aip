package apimachinery

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/swaggest/jsonrpc"
	"github.com/swaggest/swgui"
	"github.com/swaggest/swgui/v3cdn"
)

type Docs struct {
	chi.Router
}

func NewDocs(rpc *RpcService) *Docs {
	mux := &Docs{Router: chi.NewMux()}

	mux.Method(http.MethodGet, "/docs/openapi.json", rpc.OpenAPI)

	mux.Mount("/docs", v3cdn.NewHandlerWithConfig(swgui.Config{
		Title:       "AIP SUPD",
		SwaggerJSON: "/docs/openapi.json",
		BasePath:    "/docs",
		SettingsUI:  jsonrpc.SwguiSettings(nil, "/rpc/v1"),
	}))

	return mux
}
