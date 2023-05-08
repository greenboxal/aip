package graphql

import (
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/greenboxal/aip/aip-controller/pkg/apis/rpc"
	forddb "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	"github.com/greenboxal/aip/aip-controller/pkg/utils"
)

type GraphQL struct {
	*handler.Handler

	db forddb.Database

	schema graphql.Schema

	outputTypeMap map[reflect.Type]graphql.Output
	inputTypeMap  map[reflect.Type]graphql.Input
	rpcBindings   utils.BindingRegistry[rpc.RpcServiceBinding]

	rootQueryFields graphql.Fields
	rootQueryConfig graphql.ObjectConfig
	rootQuery       *graphql.Object

	rootMutationFields graphql.Fields
	rootMutationConfig graphql.ObjectConfig
	rootMutation       *graphql.Object
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
	rpcBindings utils.BindingRegistry[rpc.RpcServiceBinding],
) *GraphQL {
	gql := &GraphQL{
		db:          db,
		rpcBindings: rpcBindings,

		outputTypeMap: map[reflect.Type]graphql.Output{},
		inputTypeMap:  map[reflect.Type]graphql.Input{},
	}

	gql.initializeTypeSystem()

	for _, binding := range resourceBindings.Bindings() {
		binding.BindResource(gql)
	}

	gql.buildTypeSystem()

	gql.Handler = handler.New(&handler.Config{
		Schema:     &gql.schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	})

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
