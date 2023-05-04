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

	mux.Method(http.MethodGet, "/openapi.json", rpc.OpenAPI)

	mux.Mount("/", v3cdn.NewHandlerWithConfig(swgui.Config{
		Title:       "AIP SUPD",
		SwaggerJSON: "/v1/docs/openapi.json",
		BasePath:    "/v1/docs",
		SettingsUI:  jsonrpc.SwguiSettings(nil, "/v1/rpc"),
	}))

	return mux
}
