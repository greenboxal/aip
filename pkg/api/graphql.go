package api

import (
	"reflect"
	"regexp"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type GraphQL struct {
	*handler.Handler

	db forddb.Database

	schema  graphql.Schema
	typeMap map[reflect.Type]graphql.Output
}

func NewGraphQL(db forddb.Database) *GraphQL {
	gql := &GraphQL{
		db:      db,
		typeMap: map[reflect.Type]graphql.Output{},
	}

	gql.initializeTypeSystem()

	gql.Handler = handler.New(&handler.Config{
		Schema:     &gql.schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	})

	return gql
}

func (q *GraphQL) initializeTypeSystem() {
	q.typeMap[reflect.TypeOf((*time.Time)(nil)).Elem()] = graphql.DateTime

	ts := forddb.TypeSystem()
	fields := graphql.Fields{}

	config := graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: fields,
	}

	resourceTypes := ts.ResourceTypes()

	for i := range resourceTypes {
		typ := resourceTypes[i]

		if typ.IsRuntimeOnly() {
			continue
		}

		name := typ.Name()
		name = regexp.MustCompile("[^_a-zA-Z0-9]").ReplaceAllString(name, "_")

		gqlType := q.lookupType(typ)

		if gqlType == nil {
			panic("no gql type for " + name)
		}

		fields[name] = &graphql.Field{
			Type: gqlType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},

			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				idQuery, ok := p.Args["name"].(string)

				if !ok {
					return nil, nil
				}

				id := typ.MakeId(idQuery)

				return q.db.Get(typ.ID(), id)
			},
		}

		fields[name+"List"] = &graphql.Field{
			Type: graphql.NewList(gqlType),

			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return q.db.List(typ.ID())
			},
		}
	}

	rootQuery := graphql.NewObject(config)

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})

	if err != nil {
		panic(err)
	}

	q.schema = schema
}

func (q *GraphQL) lookupType(typ forddb.BasicResourceType) graphql.Output {
	return q.lookupTypeFromReflection(typ.ResourceType())
}

func (q *GraphQL) lookupTypeFromReflection(typ reflect.Type) (result graphql.Output) {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if existing, ok := q.typeMap[typ]; ok {
		return existing
	}

	fields := graphql.Fields{}

	switch typ.Kind() {
	case reflect.Uint64:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int:
		result = graphql.Int

	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		result = graphql.Float

	case reflect.String:
		result = graphql.String
	case reflect.Bool:
		result = graphql.Boolean

	case reflect.Array:
		elem := q.lookupTypeFromReflection(typ.Elem())

		result = graphql.NewList(elem)

	case reflect.Slice:
		elem := q.lookupTypeFromReflection(typ.Elem())

		result = graphql.NewList(elem)

	case reflect.Struct:
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)

			if !field.IsExported() {
				continue
			}

			name := field.Name
			name = regexp.MustCompile("[^_a-zA-Z0-9]").ReplaceAllString(name, "_")

			def := &graphql.Field{
				Name: name,
				Type: q.lookupTypeFromReflection(field.Type),
			}

			if def.Type == nil {
				continue
			}

			fields[def.Name] = def
		}

		name := typ.Name()
		name = regexp.MustCompile("[^_a-zA-Z0-9]").ReplaceAllString(name, "_")

		config := graphql.ObjectConfig{
			Name:   name,
			Fields: fields,
		}

		result = graphql.NewObject(config)

	default:
		result = nil
	}

	q.typeMap[typ] = result

	return
}
