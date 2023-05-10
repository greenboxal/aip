package graphql

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/ipfs/go-cid"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
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

func (q *GraphQL) RegisterTypeMapping(t reflect.Type, input graphql.Input, output graphql.Output) {
	if input != nil {
		q.inputTypeMap[t] = input
	}

	if output != nil {
		q.outputTypeMap[t] = output
	}
}

func (q *GraphQL) LookupInputType(typ forddb.BasicType) graphql.Input {
	var result graphql.Input

	if typ == nil {
		panic("typ is null")
	}

	if existing, ok := q.inputTypeMap[typ.RuntimeType()]; ok {
		return existing
	}

	if typ.Kind() == forddb.KindId {
		result = graphql.String
	} else {
		switch typ.PrimitiveKind() {
		case forddb.PrimitiveKindBoolean:
			result = graphql.Boolean
		case forddb.PrimitiveKindString:
			result = graphql.String
		case forddb.PrimitiveKindBytes:
			result = graphql.String
		case forddb.PrimitiveKindFloat:
			result = graphql.Float
		case forddb.PrimitiveKindInt:
			result = graphql.Int
		case forddb.PrimitiveKindUnsignedInt:
			result = graphql.Int

		case forddb.PrimitiveKindList:
			elemType := forddb.TypeSystem().LookupByType(typ.RuntimeType().Elem())
			elem := q.LookupInputType(elemType)

			return graphql.NewList(elem)

		case forddb.PrimitiveKindStruct:
			fields := graphql.InputObjectConfigFieldMap{}

			for _, field := range typ.Fields() {
				fieldType := field.BasicType()

				f := &graphql.InputObjectFieldConfig{
					Type: q.LookupInputType(fieldType),
				}

				fields[field.Name()] = f
			}

			if len(fields) == 0 {
				result = graphql.String
				break
			}

			name := typ.Name()

			if typ, ok := typ.(forddb.BasicResourceType); ok {
				name = typ.ResourceName().ToTitle()
			}

			name = regexp.MustCompile("[^_a-zA-Z0-9]").ReplaceAllString(name, "_")

			config := graphql.InputObjectConfig{
				Name:   "In" + name,
				Fields: fields,
			}

			result = graphql.NewInputObject(config)

		default:
			panic("unknown primitive kind")
		}
	}

	if result == nil {
		panic("result is null")
	}

	q.inputTypeMap[typ.RuntimeType()] = result

	return result
}

func (q *GraphQL) LookupOutputType(typ forddb.BasicType) graphql.Output {
	var result graphql.Output

	if existing, ok := q.outputTypeMap[typ.RuntimeType()]; ok {
		return existing
	}

	if forddb.IsBasicResourceId(typ.RuntimeType()) {
		result = graphql.String
	} else {
		switch typ.PrimitiveKind() {
		case forddb.PrimitiveKindBoolean:
			result = graphql.Boolean
		case forddb.PrimitiveKindString:
			result = graphql.String
		case forddb.PrimitiveKindBytes:
			result = graphql.String
		case forddb.PrimitiveKindFloat:
			result = graphql.Float
		case forddb.PrimitiveKindInt:
			result = graphql.Int
		case forddb.PrimitiveKindUnsignedInt:
			result = graphql.Int

		case forddb.PrimitiveKindMap:
			keyType := forddb.TypeSystem().LookupByType(typ.RuntimeType().Key())
			keyOutputType := q.LookupOutputType(keyType)

			valueType := forddb.TypeSystem().LookupByType(typ.RuntimeType().Elem())
			valueOutputType := q.LookupOutputType(valueType)

			kvType := graphql.NewObject(graphql.ObjectConfig{
				Name: fmt.Sprintf("Map_%s_%s", keyType.Name(), valueType.Name()),
				Fields: graphql.Fields{
					"key":   &graphql.Field{Name: "key", Type: keyOutputType},
					"value": &graphql.Field{Name: "value", Type: valueOutputType},
				},
			})

			return graphql.NewList(kvType)

		case forddb.PrimitiveKindList:
			elemType := forddb.TypeSystem().LookupByType(typ.RuntimeType().Elem())
			elem := q.LookupOutputType(elemType)

			return graphql.NewList(elem)

		case forddb.PrimitiveKindStruct:
			fields := graphql.Fields{
				"id": &graphql.Field{
					Name: "id",
					Type: graphql.String,
				},
			}

			for _, field := range typ.Fields() {
				fieldType := field.BasicType()

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

			name := typ.Name()

			if typ, ok := typ.(forddb.BasicResourceType); ok {
				name = typ.ResourceName().ToTitle()
			}

			name = regexp.MustCompile("[^_a-zA-Z0-9]").ReplaceAllString(name, "_")

			config := graphql.ObjectConfig{
				Name:   name,
				Fields: fields,
			}

			result = graphql.NewObject(config)

		default:
			panic("unknown primitive kind")
		}
	}

	q.outputTypeMap[typ.RuntimeType()] = result

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
