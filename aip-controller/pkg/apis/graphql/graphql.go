package graphql

import (
	"encoding/json"
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/mashingan/smapping"
	"github.com/stoewer/go-strcase"

	"github.com/greenboxal/aip/aip-controller/pkg/apimachinery"
	forddb "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	"github.com/greenboxal/aip/aip-controller/pkg/utils"
)

type GraphQL struct {
	*handler.Handler

	db forddb.Database

	schema graphql.Schema

	outputTypeMap map[reflect.Type]graphql.Output
	inputTypeMap  map[reflect.Type]graphql.Input
	rpcBindings   utils.BindingRegistry[apimachinery.RpcServiceBinding]

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

		outputTypeMap: map[reflect.Type]graphql.Output{},
		inputTypeMap:  map[reflect.Type]graphql.Input{},
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

		if inType == nil {
			continue
		}

		inBasicType := forddb.TypeSystem().LookupByType(inType)

		if inBasicType == nil {
			continue
		}

		args := graphql.FieldConfigArgument{}

		if inBasicType.PrimitiveKind() == forddb.PrimitiveKindStruct {
			for _, f := range inBasicType.Fields() {
				args[f.Name()] = &graphql.ArgumentConfig{
					Type: q.lookupInputType(f.BasicType()),
				}
			}
		} else {
			args["arg"] = &graphql.ArgumentConfig{
				Type: q.lookupInputType(inBasicType),
			}
		}

		if len(args) == 0 {
			args = nil
		}

		q.rootMutationFields[name] = &graphql.Field{
			Type: q.lookupOutputType(forddb.TypeSystem().LookupByType(outType)),

			Args: args,

			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var args [2]reflect.Value

				data, err := json.Marshal(p.Args)

				if err != nil {
					return nil, err
				}

				input := inBasicType.CreateInstance()

				if err := json.Unmarshal(data, &input); err != nil {
					return nil, err
				}

				inputVal := reflect.ValueOf(input)

				if inType.Kind() == reflect.Ptr && inputVal.Kind() != reflect.Ptr {
					inputVal = inputVal.Addr()
				} else if inType.Kind() != reflect.Ptr && inputVal.Kind() == reflect.Ptr {
					for inputVal.Kind() == reflect.Ptr {
						inputVal = inputVal.Elem()
					}
				}

				if hasCtx {
					args[0] = reflect.ValueOf(p.Context)
					args[1] = inputVal
				} else {
					args[0] = inputVal
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
							if v.Type().ConvertibleTo(outType) {
								v = v.Convert(outType)
							}

							if forddb.IsBasicResource(v.Type()) {
								resource := v.Interface().(forddb.BasicResource)
								mapped, err := forddb.Encode(resource)

								if err != nil {
									return nil, err
								}

								mapped["id"] = resource.GetResourceBasicID().String()

								return mapped, nil
							} else {
								return smapping.MapTags(v.Interface(), "json"), nil
							}
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
