package graphql

import (
	"context"
	"reflect"
	"regexp"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/samber/lo"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()
var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()

func (q *GraphQL) initializeTypeSystem() {
	q.typeMap[reflect.TypeOf((*time.Time)(nil)).Elem()] = graphql.DateTime

	// TODO: Freeze and use IPLD type system
	ts := forddb.TypeSystem()

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

	resourceTypes := ts.ResourceTypes()

	for _, typ := range resourceTypes {
		if typ.Kind() != forddb.KindResource {
			continue
		}

		if typ.IsRuntimeOnly() {
			continue
		}

		q.compileResource(q.rootQueryFields, typ)
	}
}

func (q *GraphQL) buildTypeSystem() {
	var types []graphql.Type

	for _, typ := range q.typeMap {
		types = append(types, typ)
	}

	if len(q.rootMutationFields) == 0 {
		q.rootMutationFields["noop"] = &graphql.Field{
			Type: graphql.Boolean,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return true, nil
			},
		}
	}

	if len(q.rootQueryFields) == 0 {
		q.rootQueryFields["noop"] = &graphql.Field{
			Type: graphql.Boolean,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return true, nil
			},
		}
	}

	q.rootQuery = graphql.NewObject(q.rootQueryConfig)
	q.rootMutation = graphql.NewObject(q.rootMutationConfig)

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Types: types,

		Query:    q.rootQuery,
		Mutation: q.rootMutation,
	})

	if err != nil {
		panic(err)
	}

	q.schema = schema
}

func (q *GraphQL) compileResource(fields graphql.Fields, typ forddb.BasicResourceType) {
	name := typ.Name()
	name = regexp.MustCompile("[^_a-zA-Z0-9]").ReplaceAllString(name, "_")

	gqlType := q.lookupType(typ)

	if gqlType == nil {
		panic("no gql type for " + name)
	}

	getByIdName := typ.ResourceName().ToTitle()
	allName := "all" + typ.ResourceName().ToTitlePlural()
	allMetaName := "_all" + typ.ResourceName().ToTitlePlural() + "Meta"

	fields[getByIdName] = &graphql.Field{
		Type: gqlType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},

		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			idQuery, ok := p.Args["id"].(string)

			if !ok {
				return nil, nil
			}

			id := typ.CreateID(idQuery)
			res, err := q.db.Get(p.Context, typ.GetResourceID(), id)

			if err != nil {
				return nil, err
			}

			raw, err := prepareResource(res)

			if err != nil {
				return nil, err
			}

			return raw, nil
		},
	}

	filterType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: typ.ResourceName().ToTitle() + "Filter",
		Fields: graphql.InputObjectConfigFieldMap{
			"q": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"id": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	})

	fields[allName] = &graphql.Field{
		Type: graphql.NewList(gqlType),

		Args: graphql.FieldConfigArgument{
			"page": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"perPage": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"sortField": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"sortOrder": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"filter": &graphql.ArgumentConfig{
				Type: filterType,
			},
		},

		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			pageIndex := 0
			perPage := 10

			if pageVal, ok := p.Args["page"]; ok {
				pageIndex = pageVal.(int)
			}

			if perPageVal, ok := p.Args["perPage"]; ok {
				perPage = perPageVal.(int)
			}

			results, err := q.db.List(
				p.Context,
				typ.GetResourceID(),
				forddb.WithOffset(pageIndex*perPage),
				forddb.WithLimit(perPage),
			)

			if err != nil {
				return nil, err
			}

			return lo.Map(results, func(item forddb.BasicResource, _index int) forddb.RawResource {
				raw, err := prepareResource(item)

				if err != nil {
					panic(err)
				}

				return raw
			}), nil
		},
	}

	fields[allMetaName] = &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name: typ.ResourceName().Name + "ListMetadata",
			Fields: graphql.Fields{
				"count": &graphql.Field{
					Type: graphql.Int,
				},
			},
		}),

		Args: graphql.FieldConfigArgument{
			"page": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"perPage": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"sortField": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"sortOrder": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"filter": &graphql.ArgumentConfig{
				Type: filterType,
			},
		},

		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			res, err := q.db.List(p.Context, typ.GetResourceID())

			if err != nil {
				return nil, err
			}

			return map[string]interface{}{
				"count": len(res),
			}, nil
		},
	}
}

func (q *GraphQL) lookupType(typ forddb.BasicType) graphql.Output {
	var result graphql.Output

	if existing, ok := q.typeMap[typ.RuntimeType()]; ok {
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

		case forddb.PrimitiveKindList:
			elemType := forddb.TypeSystem().LookupByType(typ.RuntimeType().Elem())
			elem := q.lookupType(elemType)

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
					Type: q.lookupType(fieldType),
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

	q.typeMap[typ.RuntimeType()] = result

	return result
}
