package graphql

import (
	"reflect"

	"github.com/graphql-go/graphql"
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-forddb/pkg/apimachinery"
	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
	"github.com/greenboxal/aip/aip-sdk/pkg/utils"
)

var Module = fx.Module(
	"apis/graphql",

	fx.Provide(NewSubscriptionManager),

	apimachinery.ProvideHttpService[*GraphQL](NewGraphQL, "/v1/graphql"),

	utils.WithBindingRegistry[ResourceBinding]("graphql-resource-bindings"),

	ProvideBinding[*RpcResourceBinding](NewRpcResourceBinding),
	ProvideBinding[*DatabaseResourceBinding](NewDatabaseResourceBinding),
)

type BindingContext interface {
	RegisterTypeMapping(t reflect.Type, input graphql.Input, output graphql.Output)

	RegisterQuery(field ...*graphql.Field)
	RegisterMutation(field ...*graphql.Field)
	RegisterSubscription(fields ...*graphql.Field)

	LookupInputType(typ typesystem.Type) graphql.Input
	LookupOutputType(typ typesystem.Type) graphql.Output
}

type ResourceBinding interface {
	BindResource(ctx BindingContext)
}

func ProvideBinding[T ResourceBinding](constructor any) fx.Option {
	return fx.Options(
		fx.Provide(constructor),

		WithBinding[T](),
	)
}

func WithBinding[T ResourceBinding]() fx.Option {
	return utils.WithBinding[ResourceBinding]("graphql-resource-bindings", func(t T) ResourceBinding {
		return t
	})
}
