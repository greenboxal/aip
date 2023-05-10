package graphql

import (
	"net/http"
	"reflect"
	"time"

	"github.com/eientei/wsgraphql/v1"
	"github.com/eientei/wsgraphql/v1/compat/gorillaws"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	"github.com/greenboxal/aip/aip-controller/pkg/utils"
)

type GraphQL struct {
	chi.Router

	handler   *handler.Handler
	wsHandler http.Handler

	db forddb.Database

	schema graphql.Schema

	outputTypeMap map[reflect.Type]graphql.Output
	inputTypeMap  map[reflect.Type]graphql.Input

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

func (q *GraphQL) RegisterTypeMapping(t reflect.Type, input graphql.Input, output graphql.Output) {
	if input != nil {
		q.inputTypeMap[t] = input
	}

	if output != nil {
		q.outputTypeMap[t] = output
	}
}

func NewGraphQL(
	db forddb.Database,
	resourceBindings utils.BindingRegistry[ResourceBinding],
) *GraphQL {
	gql := &GraphQL{
		Router: chi.NewRouter(),

		db: db,

		outputTypeMap: map[reflect.Type]graphql.Output{},
		inputTypeMap:  map[reflect.Type]graphql.Input{},
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

	gql.wsHandler = wsHandler

	gql.handler = handler.New(&handler.Config{
		Schema:     &gql.schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	})

	gql.Mount("/ws", gql.wsHandler)
	gql.Mount("/", gql.handler)

	return gql
}

func prepareResource(res forddb.BasicResource) (forddb.RawResource, error) {
	raw, err := forddb.Encode(res)

	if err != nil {
		return nil, err
	}

	raw["id"] = res.GetResourceBasicID().String()

	return raw, nil
}
