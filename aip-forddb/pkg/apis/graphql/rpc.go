package graphql

import (
	"encoding/json"
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/stoewer/go-strcase"

	"github.com/greenboxal/aip/aip-forddb/pkg/apis/rpc"
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
	"github.com/greenboxal/aip/aip-sdk/pkg/utils"
)

type RpcResourceBinding struct {
	rpcBindings utils.BindingRegistry[rpc.RpcServiceBinding]
}

func NewRpcResourceBinding(
	rpcBindings utils.BindingRegistry[rpc.RpcServiceBinding],
) *RpcResourceBinding {
	return &RpcResourceBinding{
		rpcBindings: rpcBindings,
	}
}

func (r *RpcResourceBinding) BindResource(ctx BindingContext) {
	for _, binding := range r.rpcBindings.Bindings() {
		r.compileRpcMutation(ctx, binding)
	}
}

func (r *RpcResourceBinding) compileRpcMutation(ctx BindingContext, binding rpc.RpcServiceBinding) {
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

		inTypePtr := inType
		outTypePtr := inType

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

		if inBasicType.ActualType().PrimitiveKind() == typesystem.PrimitiveKindStruct {
			st := inBasicType.ActualType().(typesystem.StructType)

			for i := 0; i < st.NumField(); i++ {
				f := st.FieldByIndex(i)

				args[f.Name()] = &graphql.ArgumentConfig{
					Type: ctx.LookupInputType(f.Type()),
				}
			}
		} else {
			args["arg"] = &graphql.ArgumentConfig{
				Type: ctx.LookupInputType(inBasicType.ActualType()),
			}
		}

		if len(args) == 0 {
			args = nil
		}

		field := &graphql.Field{
			Name: name,

			Type: ctx.LookupOutputType(typesystem.TypeFrom(outType)),

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

				if inTypePtr.Kind() == reflect.Ptr && inputVal.Kind() != reflect.Ptr {
					inputVal = inputVal.Addr()
				} else if inTypePtr.Kind() != reflect.Ptr && inputVal.Kind() == reflect.Ptr {
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
					if outTypePtr == errorType {
						if result[0].IsNil() {
							return nil, nil
						}

						return nil, result[0].Interface().(error)
					} else {
						v := result[0]

						if v.IsValid() {
							if v.Type().ConvertibleTo(outTypePtr) {
								v = v.Convert(outTypePtr)
							}

							return v.Interface(), nil
						}
					}
				}

				return nil, nil
			},
		}

		ctx.RegisterMutation(field)
	}
}
