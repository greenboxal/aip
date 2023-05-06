package graphql

import (
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/stoewer/go-strcase"

	"github.com/greenboxal/aip/aip-controller/pkg/apimachinery"
	forddb "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	"github.com/greenboxal/aip/aip-controller/pkg/utils"
)

type GraphQL struct {
	*handler.Handler

	db forddb.Database

	schema graphql.Schema

	typeMap     map[reflect.Type]graphql.Output
	rpcBindings utils.BindingRegistry[apimachinery.RpcServiceBinding]

	rootQueryFields graphql.Fields
	rootQueryConfig graphql.ObjectConfig
	rootQuery       *graphql.Object

	rootMutationFields graphql.Fields
	rootMutationConfig graphql.ObjectConfig
	rootMutation       *graphql.Object
}

func NewGraphQL(
	db forddb.Database,
	rpcBindings utils.BindingRegistry[apimachinery.RpcServiceBinding],
) *GraphQL {
	gql := &GraphQL{
		db:          db,
		rpcBindings: rpcBindings,

		typeMap: map[reflect.Type]graphql.Output{},
	}

	gql.initializeTypeSystem()
	gql.initializeRpcMutations()
	gql.buildTypeSystem()

	gql.Handler = handler.New(&handler.Config{
		Schema:     &gql.schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	})

	return gql
}

func (q *GraphQL) initializeRpcMutations() {
	for _, binding := range q.rpcBindings.Bindings() {
		q.compileRpcMutation(binding)
	}
}

func (q *GraphQL) compileRpcMutation(binding apimachinery.RpcServiceBinding) {
	target := binding.Implementation()
	value := reflect.ValueOf(target)
	typ := value.Type()

	for i := 0; i < typ.NumMethod(); i++ {
		var inType reflect.Type
		var outType reflect.Type

		m := typ.Method(i)
		mi := value.Method(m.Index)
		mtyp := mi.Type()

		hasCtx := false
		hasError := false

		if !m.IsExported() {
			continue
		}

		if mtyp.NumIn() == 2 {
			if !mtyp.In(0).ConvertibleTo(contextType) {
				continue
			}

			hasCtx = true
			inType = mtyp.In(1)
		} else if mtyp.NumIn() == 1 {
			inType = mtyp.In(0)
		} else {
			continue
		}

		if mtyp.NumOut() == 2 {
			if !mtyp.Out(1).ConvertibleTo(errorType) {
				continue
			}

			hasError = true
			outType = mtyp.Out(0)
		} else if mtyp.NumOut() == 1 {
			outType = mtyp.Out(0)
		} else {
			continue
		}

		if inType == nil {
			inType = reflect.TypeOf(struct{}{})
		}

		if outType == nil {
			outType = reflect.TypeOf(struct{}{})
		}

		for inType.Kind() == reflect.Ptr {
			inType = inType.Elem()
		}

		for outType.Kind() == reflect.Ptr {
			outType = outType.Elem()
		}

		if inType.Kind() == reflect.Interface {
			continue
		}

		if outType.Kind() == reflect.Interface {
			continue
		}

		name := strcase.LowerCamelCase(binding.Name()) + strcase.UpperCamelCase(m.Name)
		args := graphql.FieldConfigArgument{}

		if inType != nil {
			args["args"] = &graphql.ArgumentConfig{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "empty",
					Fields: graphql.Fields{
						"empty": &graphql.Field{
							Type: graphql.String,
						},
					},
				}),
			}
		} else {
			args["args"] = &graphql.ArgumentConfig{
				Type: q.lookupType(forddb.TypeSystem().LookupByType(inType)),
			}
		}

		q.rootMutationFields[name] = &graphql.Field{
			Type: q.lookupType(forddb.TypeSystem().LookupByType(outType)),

			Args: map[string]*graphql.ArgumentConfig{},

			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var args [2]reflect.Value

				input := p.Args["args"]

				if hasCtx {
					args[0] = reflect.ValueOf(p.Context)
					args[1] = reflect.ValueOf(input)
				} else {
					args[0] = reflect.ValueOf(input)
				}

				result := mi.Call(args[:mtyp.NumIn()])

				if hasError && result[1].IsValid() {
					err := result[1].Interface()

					if err != nil {
						return nil, err.(error)
					}
				}

				if len(result) > 0 {
					if outType == errorType {
						if result[0].IsNil() {
							return nil, nil
						}

						return nil, result[0].Interface().(error)
					} else {
						v := result[0]

						if v.IsValid() {
							for v.Kind() == reflect.Ptr {
								v = v.Elem()
							}

							return v.Interface(), nil
						}
					}
				}

				return nil, nil
			},
		}
	}
}

func prepareResource(res forddb.BasicResource) (forddb.RawResource, error) {
	raw, err := forddb.Encode(res)

	if err != nil {
		return nil, err
	}

	raw["id"] = res.GetResourceBasicID().String()

	return raw, nil
}
