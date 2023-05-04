package graphql

import (
	"reflect"
	"regexp"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	forddb2 "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type GraphQL struct {
	*handler.Handler

	db forddb2.Database

	schema graphql.Schema

	typeMap map[reflect.Type]graphql.Output
}

func NewGraphQL(db forddb2.Database) *GraphQL {
	gql := &GraphQL{
		db: db,

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

	// TODO: Freeze and use IPLD type system
	ts := forddb2.TypeSystem()

	fields := graphql.Fields{
		"empty": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "", nil
			},
		},
	}

	config := graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: fields,
	}

	resourceTypes := ts.ResourceTypes()

	for _, typ := range resourceTypes {
		if typ.Kind() != forddb2.KindResource {
			continue
		}

		if typ.IsRuntimeOnly() {
			continue
		}

		q.compileResource(fields, typ)
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

func (q *GraphQL) compileResource(fields graphql.Fields, typ forddb2.BasicResourceType) {
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

			id := typ.CreateID(idQuery)

			return q.db.Get(p.Context, typ.GetID(), id)
		},
	}

	fields[name+"List"] = &graphql.Field{
		Type: graphql.NewList(gqlType),

		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return q.db.List(p.Context, typ.GetID())
		},
	}
}

func (q *GraphQL) lookupType(typ forddb2.BasicType) graphql.Output {
	var result graphql.Output

	if existing, ok := q.typeMap[typ.RuntimeType()]; ok {
		return existing
	}

	switch typ.PrimitiveKind() {
	case forddb2.PrimitiveKindBoolean:
		result = graphql.Boolean
	case forddb2.PrimitiveKindString:
		result = graphql.String
	case forddb2.PrimitiveKindBytes:
		result = graphql.String
	case forddb2.PrimitiveKindFloat:
		result = graphql.Float
	case forddb2.PrimitiveKindInt:
		result = graphql.Int
	case forddb2.PrimitiveKindUnsignedInt:
		result = graphql.Int

	case forddb2.PrimitiveKindList:
		elem := q.lookupTypeFromReflection(typ.RuntimeType().Elem())

		return graphql.NewList(elem)

	case forddb2.PrimitiveKindStruct:
		fields := graphql.Fields{}

		for _, field := range typ.Fields() {
			fieldType := field.BasicType()

			f := &graphql.Field{
				Name: field.Name(),
				Type: q.lookupType(fieldType),
			}

			fields[f.Name] = f
		}

		if len(fields) == 0 {
			result = graphql.String
			break
		}

		name := typ.Name()
		name = regexp.MustCompile("[^_a-zA-Z0-9]").ReplaceAllString(name, "_")

		config := graphql.ObjectConfig{
			Name:   name,
			Fields: fields,
		}

		result = graphql.NewObject(config)

	default:
		panic("unknown primitive kind")
	}

	q.typeMap[typ.RuntimeType()] = result

	return result
}

func (q *GraphQL) lookupTypeFromReflection(typ reflect.Type) (result graphql.Output) {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if existing, ok := q.typeMap[typ]; ok {
		return existing
	}

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
		if typ.Elem().Kind() == reflect.Uint8 {
			result = graphql.String
			break
		}

		elem := q.lookupTypeFromReflection(typ.Elem())

		result = graphql.NewList(elem)

	case reflect.Struct:
		fields := graphql.Fields{}

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
