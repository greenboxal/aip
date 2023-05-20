package graphql

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()
var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
var timeType = reflect.TypeOf((*time.Time)(nil)).Elem()
var rawResourceType = reflect.TypeOf((*forddb.RawResource)(nil)).Elem()
var cidType = reflect.TypeOf((*cid.Cid)(nil)).Elem()

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

			return graphql.NewList(elem)

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

				return graphql.NewList(kvType)
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

			return graphql.NewList(elem)

		case ipld.Kind_Map:
			if typ.PrimitiveKind() == typesystem.PrimitiveKindStruct {
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

					fields[f.Name] = f
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

				config := graphql.ObjectConfig{
					Name:   name,
					Fields: fields,
				}

				result = graphql.NewObject(config)
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

				return graphql.NewList(kvType)
			}

		default:
			panic("unknown primitive kind")
		}
	}

	q.outputTypeMap[typ] = result

	return result
}

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
