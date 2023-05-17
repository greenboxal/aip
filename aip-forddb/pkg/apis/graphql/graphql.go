package graphql

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/eientei/wsgraphql/v1"
	"github.com/eientei/wsgraphql/v1/compat/gorillaws"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
	"github.com/greenboxal/aip/aip-sdk/pkg/utils"
)

type GraphQL struct {
	chi.Router

	db forddb.Database

	schema graphql.Schema

	outputTypeMap map[typesystem.Type]graphql.Output
	inputTypeMap  map[typesystem.Type]graphql.Input

	rootQueryFields graphql.Fields
	rootQueryConfig graphql.ObjectConfig
	rootQuery       *graphql.Object

	rootMutationFields graphql.Fields
	rootMutationConfig graphql.ObjectConfig
	rootMutation       *graphql.Object

	subscriptionFields graphql.Fields
	subscriptionConfig graphql.ObjectConfig
	subscription       *graphql.Object
}

func NewGraphQL(
	db forddb.Database,
	resourceBindings utils.BindingRegistry[ResourceBinding],
) *GraphQL {
	gql := &GraphQL{
		Router: chi.NewRouter(),

		db: db,

		outputTypeMap: map[typesystem.Type]graphql.Output{},
		inputTypeMap:  map[typesystem.Type]graphql.Input{},
	}

	gql.initializeTypeSystem()

	for _, binding := range resourceBindings.Bindings() {
		binding.BindResource(gql)
	}

	gql.buildTypeSystem()

	wsHandler, err := wsgraphql.NewServer(
		gql.schema,
		wsgraphql.WithProtocol(wsgraphql.WebsocketSubprotocolGraphqlTransportWS),
		wsgraphql.WithKeepalive(30*time.Second),
		wsgraphql.WithUpgrader(gorillaws.Wrap(&websocket.Upgrader{
			CheckOrigin:  func(r *http.Request) bool { return true },
			Subprotocols: []string{wsgraphql.WebsocketSubprotocolGraphqlTransportWS.String()},
		})),
	)

	if err != nil {
		panic(err)
	}

	httpHandler := handler.New(&handler.Config{
		Schema:     &gql.schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	})

	gql.Mount("/ws", wsHandler)
	gql.Mount("/", httpHandler)

	return gql
}

func prepareResource(res forddb.BasicResource) (any, error) {
	var raw map[string]interface{}

	data, err := ipld.Encode(typesystem.Wrap(res), dagjson.Encode)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	raw["id"] = res.GetResourceBasicID().String()

	return raw, nil
}
