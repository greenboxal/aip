package apimachinery

import (
	"context"
	"reflect"

	"github.com/swaggest/jsonrpc"
	"github.com/swaggest/usecase"

	"github.com/greenboxal/aip/aip-controller/pkg/utils"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()
var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()

type RpcService struct {
	*jsonrpc.Handler
}

func NewRpcService() *RpcService {
	handler := &jsonrpc.Handler{}

	apiSchema := jsonrpc.OpenAPI{}
	apiSchema.Reflector().SpecEns().Info.Title = "supd"
	apiSchema.Reflector().SpecEns().Info.Version = "v1.0.0"

	apiSchema.Reflector().InterceptDefName(func(t reflect.Type, defaultDefName string) string {
		return utils.GetParsedTypeName(t).NormalizedFullNameWithArguments()
	})

	handler.OpenAPI = &apiSchema
	handler.Validator = &jsonrpc.JSONSchemaValidator{}
	handler.SkipResultValidation = true

	return &RpcService{Handler: handler}
}

func mustRegister(srv *jsonrpc.Handler, name string, target any) {
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
		} else {
			inType = mtyp.In(0)
		}

		if mtyp.NumOut() == 2 {
			if !mtyp.Out(1).ConvertibleTo(errorType) {
				continue
			}

			hasError = true
			outType = mtyp.Out(0)
		} else if mtyp.NumOut() > 0 {
			outType = mtyp.Out(0)
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

		u := usecase.NewIOI(
			reflect.New(inType).Interface(),
			reflect.New(outType).Interface(),
			func(ctx context.Context, input, output interface{}) error {
				var args [2]reflect.Value

				if hasCtx {
					args[0] = reflect.ValueOf(ctx)
					args[1] = reflect.ValueOf(input)
				} else {
					args[0] = reflect.ValueOf(input)
				}

				result := mi.Call(args[:mtyp.NumIn()])

				if hasError && result[1].IsValid() {
					err := result[1].Interface()

					if err != nil {
						return err.(error)
					}
				}

				if len(result) > 0 {
					if outType == errorType {
						if result[0].IsNil() {
							return nil
						}

						return result[0].Interface().(error)
					} else {
						v := result[0]

						if v.IsValid() {
							for v.Kind() == reflect.Ptr {
								v = v.Elem()
							}

							reflect.ValueOf(output).Elem().Set(v)
						}
					}
				}

				return nil
			},
		)

		u.SetName(name + "." + m.Name)

		srv.Add(u)
	}
}