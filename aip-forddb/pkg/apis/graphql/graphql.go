package graphql

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/eientei/wsgraphql/v1"
	"github.com/eientei/wsgraphql/v1/compat/gorillaws"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/samber/lo"

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

	gql.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "loaders", &Loaders{
				db:      db,
				loaders: map[typesystem.Type]*dataloader.Loader{},
			})

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	})

	gql.Mount("/ws", wsHandler)
	gql.Mount("/", httpHandler)

	return gql
}

type Loaders struct {
	m       sync.Mutex
	db      forddb.Database
	loaders map[typesystem.Type]*dataloader.Loader
}

func (l *Loaders) Get(typ typesystem.Type) *dataloader.Loader {
	l.m.Lock()
	defer l.m.Unlock()

	if existing := l.loaders[typ]; existing != nil {
		return existing
	}

	resourceType := forddb.TypeSystem().LookupByResourceType(typ.RuntimeType())

	loader := dataloader.NewBatchedLoader(func(
		ctx context.Context,
		keys dataloader.Keys,
	) (results []*dataloader.Result) {
		var wg sync.WaitGroup

		results = make([]*dataloader.Result, len(keys))

		targetIds := lo.Map(keys, func(key dataloader.Key, _ int) string {
			return key.String()
		})

		for baseIndex := 0; baseIndex < len(keys); baseIndex += 30 {
			wg.Add(1)

			baseIndex := baseIndex
			batchIds := targetIds[baseIndex:utils.Min(baseIndex+30, len(keys))]

			go func() {
				defer wg.Done()

				resById := make(map[string]any, len(keys))

				res, err := l.db.List(
					ctx,
					resourceType.GetResourceID(),
					forddb.WithListQueryOptions(
						forddb.WithFilterExpression("resource.metadata.id in args.ids"),
						forddb.WithFilterParameter("ids", batchIds),
					),
				)

				for _, r := range res {
					resById[r.GetResourceBasicID().String()] = r
				}

				for i, id := range batchIds {
					results[baseIndex+i] = &dataloader.Result{Data: resById[id], Error: err}
				}
			}()
		}

		wg.Wait()

		return
	})

	l.loaders[typ] = loader

	return loader
}
