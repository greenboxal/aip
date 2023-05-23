package graphql

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	"github.com/hashicorp/go-multierror"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/samber/lo"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
)

var stringType = reflect.TypeOf((*string)(nil)).Elem()
var errorType = reflect.TypeOf((*error)(nil)).Elem()
var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
var timeType = reflect.TypeOf((*time.Time)(nil)).Elem()
var rawResourceType = reflect.TypeOf((*forddb.RawResource)(nil)).Elem()
var cidType = reflect.TypeOf((*cid.Cid)(nil)).Elem()

func (q *GraphQL) RegisterQuery(fields ...*graphql.Field) {
	for _, field := range fields {
		q.rootQueryFields[field.Name] = field
	}
}

func (q *GraphQL) RegisterMutation(fields ...*graphql.Field) {
	for _, field := range fields {
		q.rootMutationFields[field.Name] = field
	}
}

func (q *GraphQL) RegisterSubscription(fields ...*graphql.Field) {
	for _, field := range fields {
		q.subscriptionFields[field.Name] = field
	}
}

func (q *GraphQL) initializeTypeSystem() {
	q.RegisterTypeMapping(timeType, graphql.DateTime, graphql.DateTime)
	q.RegisterTypeMapping(cidType, graphql.String, graphql.String)
	q.RegisterTypeMapping(rawResourceType, graphql.String, graphql.String)

	q.rootQueryFields = graphql.Fields{}
	q.rootQueryConfig = graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: q.rootQueryFields,
	}

	q.rootMutationFields = graphql.Fields{}
	q.rootMutationConfig = graphql.ObjectConfig{
		Name:   "Mutations",
		Fields: q.rootMutationFields,
	}

	q.subscriptionFields = graphql.Fields{}
	q.subscriptionConfig = graphql.ObjectConfig{
		Name:   "Subscriptions",
		Fields: q.subscriptionFields,
	}
}

func (q *GraphQL) buildTypeSystem() {
	var types []graphql.Type

	for _, typ := range q.outputTypeMap {
		types = append(types, typ)
	}

	q.rootQuery = graphql.NewObject(q.rootQueryConfig)
	q.rootMutation = graphql.NewObject(q.rootMutationConfig)
	q.subscription = graphql.NewObject(q.subscriptionConfig)

	schemaConfig := graphql.SchemaConfig{
		Types: types,
	}

	if len(q.rootMutationFields) > 0 {
		schemaConfig.Mutation = q.rootMutation
	}

	if len(q.rootQueryFields) > 0 {
		schemaConfig.Query = q.rootQuery
	}

	if len(q.subscriptionFields) > 0 {
		schemaConfig.Subscription = q.subscription
	}

	schema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		panic(err)
	}

	q.schema = schema
}

func (q *GraphQL) RegisterTypeMapping(typ reflect.Type, input graphql.Input, output graphql.Output) {
	t := typesystem.TypeOf(typ)

	if input != nil {
		q.inputTypeMap[t] = input
	}

	if output != nil {
		q.outputTypeMap[t] = output
	}
}

func (q *GraphQL) LookupInputType(typ typesystem.Type) graphql.Input {
	var result graphql.Input

	if typ == nil {
		panic("typ is null")
	}

	if existing, ok := q.inputTypeMap[typ]; ok {
		return existing
	}

	if forddb.IsBasicResourceId(typ.RuntimeType()) {
		result = graphql.String
	} else if typ.PrimitiveKind() == typesystem.PrimitiveKindInterface {
		result = graphql.String
	} else {
		result = graphql.String
		switch typ.IpldRepresentationKind() {
		case ipld.Kind_Bool:
			result = graphql.Boolean
		case ipld.Kind_String:
			result = graphql.String
		case ipld.Kind_Bytes:
			result = graphql.String
		case ipld.Kind_Float:
			result = graphql.Float
		case ipld.Kind_Int:
			result = graphql.Int

		case ipld.Kind_List:
			lt := typ.List()
			elem := q.LookupInputType(lt.Elem())

			result = graphql.NewList(elem)

		case ipld.Kind_Map:
			if typ.PrimitiveKind() == typesystem.PrimitiveKindStruct {
				fields := graphql.InputObjectConfigFieldMap{}

				st := typ.Struct()

				for i := 0; i < st.NumField(); i++ {
					field := st.FieldByIndex(i)
					fieldType := field.Type()

					f := &graphql.InputObjectFieldConfig{
						Type: q.LookupInputType(fieldType),
					}

					fields[field.Name()] = f
				}

				if len(fields) == 0 {
					result = graphql.String
					break
				}

				name := typ.Name().NormalizedFullNameWithArguments()

				resTyp := forddb.TypeSystem().LookupByType(typ.RuntimeType())

				if typ, ok := resTyp.(forddb.BasicResourceType); ok {
					name = typ.ResourceName().ToTitle()
				}

				name = regexp.MustCompile("[^_a-zA-Z0-9]").ReplaceAllString(name, "_")

				config := graphql.InputObjectConfig{
					Name:   "In" + name,
					Fields: fields,
				}

				result = graphql.NewInputObject(config)
			} else {
				mt := typ.Map()

				keyOutputType := q.LookupInputType(mt.Key())
				valueOutputType := q.LookupInputType(mt.Map())

				kvType := graphql.NewInputObject(graphql.InputObjectConfig{
					Name: fmt.Sprintf(
						"Map_%s_%s",
						mt.Key().Name().NormalizedFullNameWithArguments(),
						mt.Value().Name().NormalizedFullNameWithArguments(),
					),

					Fields: graphql.Fields{
						"key":   &graphql.Field{Name: "key", Type: keyOutputType},
						"value": &graphql.Field{Name: "value", Type: valueOutputType},
					},
				})

				result = graphql.NewList(kvType)
			}

		default:
			panic("unknown primitive kind")
		}
	}

	if result == nil {
		panic("result is null")
	}

	q.inputTypeMap[typ] = result

	return result
}

func (q *GraphQL) LookupOutputType(typ typesystem.Type) graphql.Output {
	var result graphql.Output

	if existing, ok := q.outputTypeMap[typ]; ok {
		return existing
	}

	if forddb.IsBasicResourceId(typ.RuntimeType()) {
		result = graphql.String
	} else if typ.PrimitiveKind() == typesystem.PrimitiveKindInterface {
		result = graphql.String
	} else {
		switch typ.IpldRepresentationKind() {
		case ipld.Kind_Bool:
			result = graphql.Boolean
		case ipld.Kind_String:
			result = graphql.String
		case ipld.Kind_Bytes:
			result = graphql.String
		case ipld.Kind_Float:
			result = graphql.Float
		case ipld.Kind_Int:
			result = graphql.Int
		case ipld.Kind_Link:
			result = graphql.String

		case ipld.Kind_List:
			lt := typ.List()
			elem := q.LookupOutputType(lt.Elem())

			result = graphql.NewList(elem)

		case ipld.Kind_Map:
			if typ.PrimitiveKind() == typesystem.PrimitiveKindStruct {
				result = q.compileOutputStruct(typ)
			} else {
				mt := typ.Map()

				keyOutputType := q.LookupOutputType(mt.Key())
				valueOutputType := q.LookupOutputType(mt.Map())

				kvType := graphql.NewObject(graphql.ObjectConfig{
					Name: fmt.Sprintf(
						"Map_%s_%s",
						mt.Key().Name().NormalizedFullNameWithArguments(),
						mt.Value().Name().NormalizedFullNameWithArguments(),
					),

					Fields: graphql.Fields{
						"key":   &graphql.Field{Name: "key", Type: keyOutputType},
						"value": &graphql.Field{Name: "value", Type: valueOutputType},
					},
				})

				result = graphql.NewList(kvType)
			}

		default:
			panic("unknown primitive kind")
		}
	}

	q.outputTypeMap[typ] = result

	return result
}

func (q *GraphQL) compileOutputStruct(typ typesystem.Type) graphql.Output {

	name := typ.Name().NormalizedFullNameWithArguments()

	resTyp := forddb.TypeSystem().LookupByType(typ.RuntimeType())

	if typ, ok := resTyp.(forddb.BasicResourceType); ok {
		name = typ.ResourceName().ToTitle()
	}

	name = regexp.MustCompile("[^_a-zA-Z0-9]").ReplaceAllString(name, "_")

	config := graphql.ObjectConfig{
		Name: name,

		Fields: graphql.FieldsThunk(func() graphql.Fields {
			fields := graphql.Fields{
				"id": &graphql.Field{
					Name: "id",
					Type: graphql.String,
				},
			}

			st := typ.Struct()

			for i := 0; i < st.NumField(); i++ {
				field := st.FieldByIndex(i)
				fieldType := field.Type()

				f := &graphql.Field{
					Name: field.Name(),
					Type: q.LookupOutputType(fieldType),
				}

				if strings.HasSuffix(f.Name, "_id") && forddb.IsBasicResourceId(fieldType.RuntimeType()) {
					expandedFieldType := forddb.TypeSystem().LookupByIDType(fieldType.RuntimeType())
					expandedFieldGqlType := q.LookupOutputType(expandedFieldType.ActualType())

					expandedField := &graphql.Field{
						Name: strings.TrimSuffix(f.Name, "_id"),
						Type: expandedFieldGqlType,

						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							receiver := reflect.ValueOf(p.Source)
							loaders := p.Context.Value("loaders").(*Loaders)
							loader := loaders.Get(expandedFieldType.ActualType())

							fieldValue := receiver.MapIndex(reflect.ValueOf(name))

							if !fieldValue.IsValid() {
								return nil, nil
							}

							if fieldValue.CanInterface() && reflect.DeepEqual(fieldValue.Interface(), reflect.Zero(fieldValue.Type()).Interface()) {
								return nil, nil
							}

							v, ok := typesystem.TryCast[forddb.BasicResourceID](fieldValue)

							if !ok {
								return nil, errors.New("invalid id")
							}

							ch := make(chan dataLoaderSingleResult, 1)
							thunk := loader.Load(p.Context, dataloader.StringKey(v.String()))

							go func() {
								res, err := thunk()

								ch <- dataLoaderSingleResult{res.(forddb.BasicResource), err}
							}()

							return func() (any, error) {
								result := <-ch

								if result.Error != nil {
									return nil, result.Error
								}

								return prepareResource(result.Value)
							}, nil
						},
					}

					fields[expandedField.Name] = expandedField
				} else if strings.HasSuffix(f.Name, "_ids") && fieldType.PrimitiveKind() == typesystem.PrimitiveKindList {
					lt := fieldType.List()

					if forddb.IsBasicResourceId(lt.Elem().RuntimeType()) {
						expandedFieldType := forddb.TypeSystem().LookupByIDType(lt.Elem().RuntimeType())
						expandedFieldGqlType := q.LookupOutputType(expandedFieldType.ActualType())

						expandedField := &graphql.Field{
							Name: strings.TrimSuffix(f.Name, "_ids"),
							Type: graphql.NewList(expandedFieldGqlType),

							Resolve: func(p graphql.ResolveParams) (interface{}, error) {
								var targetIds []dataloader.Key

								receiver := reflect.ValueOf(p.Source)
								loaders := p.Context.Value("loaders").(*Loaders)
								loader := loaders.Get(expandedFieldType.ActualType())

								fieldValue := receiver.MapIndex(reflect.ValueOf(name))

								if fieldValue.IsValid() && !fieldValue.IsNil() {
									targetIds = make([]dataloader.Key, fieldValue.Len())

									for i := 0; i < fieldValue.Len(); i++ {
										v, ok := typesystem.TryCast[forddb.BasicResourceID](fieldValue.Index(i))

										if !ok {
											return nil, errors.New("invalid id")
										}

										targetIds[i] = dataloader.StringKey(v.String())
									}
								} else {
									return nil, nil
								}

								ch := make(chan dataLoaderManyResult, 1)
								thunk := loader.LoadMany(p.Context, targetIds)

								go func() {
									var merr error

									res, err := thunk()

									for _, e := range err {
										merr = multierror.Append(merr, e)
									}

									resources := lo.Map(res, func(item any, index int) forddb.BasicResource {
										return item.(forddb.BasicResource)
									})

									ch <- dataLoaderManyResult{resources, merr}
								}()

								return func() (any, error) {
									result := <-ch

									if result.Error != nil {
										return nil, result.Error
									}

									return mapResources(result.Value)
								}, nil
							},
						}

						fields[expandedField.Name] = expandedField
					}
				}

				fields[f.Name] = f
			}

			return fields
		}),
	}

	return graphql.NewObject(config)
}

type dataLoaderSingleResult struct {
	Value forddb.BasicResource
	Error error
}

type dataLoaderManyResult struct {
	Value []forddb.BasicResource
	Error error
}
