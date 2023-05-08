package graphql

import (
	"regexp"

	"github.com/graphql-go/graphql"
	"github.com/samber/lo"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type DatabaseResourceBinding struct {
	db forddb.Database
}

func NewDatabaseResourceBinding(
	db forddb.Database,
) *DatabaseResourceBinding {
	return &DatabaseResourceBinding{
		db: db,
	}
}

func (r *DatabaseResourceBinding) BindResource(ctx BindingContext) {
	ts := forddb.TypeSystem()
	resourceTypes := ts.ResourceTypes()

	for _, typ := range resourceTypes {
		if typ.Kind() != forddb.KindResource {
			continue
		}

		if typ.IsRuntimeOnly() {
			continue
		}

		r.compileResource(ctx, typ)
	}
}

func (r *DatabaseResourceBinding) compileResource(ctx BindingContext, typ forddb.BasicResourceType) {
	name := typ.Name()
	name = regexp.MustCompile("[^_a-zA-Z0-9]").ReplaceAllString(name, "_")

	gqlType := ctx.LookupOutputType(typ)

	if gqlType == nil {
		panic("no gql type for " + name)
	}

	getByIdName := typ.ResourceName().ToTitle()
	getAllName := "all" + typ.ResourceName().ToTitlePlural()
	getAllMetaName := "_all" + typ.ResourceName().ToTitlePlural() + "Meta"

	getById := &graphql.Field{
		Name: getByIdName,
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
			res, err := r.db.Get(p.Context, typ.GetResourceTypeID(), id)

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
			"ids": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(graphql.String),
			},
		},
	})

	getAll := &graphql.Field{
		Name: getAllName,
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
			var options []forddb.ListOption

			pageIndex := 0
			perPage := -1

			if pageVal, ok := p.Args["page"]; ok {
				pageIndex = pageVal.(int)
			}

			if perPageVal, ok := p.Args["perPage"]; ok {
				perPage = perPageVal.(int)
			}

			if sortFieldVal, ok := p.Args["sortField"]; ok {
				sortField := sortFieldVal.(string)
				sortOrder := forddb.Asc

				if sortOrderVal, ok := p.Args["sortOrder"]; ok {
					sortOrder = forddb.SortOrder(sortOrderVal.(string))
				}

				options = append(options, forddb.WithSortField(sortField, sortOrder))
			}

			if filterVal, ok := p.Args["filter"]; ok {
				var ids []forddb.BasicResourceID

				filter := filterVal.(map[string]interface{})

				if id, ok := filter["id"]; ok {
					ids = append(ids, typ.CreateID(id.(string)))
				}

				if stringIds, ok := filter["ids"]; ok {
					for _, id := range stringIds.([]interface{}) {
						ids = append(ids, typ.CreateID(id.(string)))
					}
				}

				if len(ids) > 0 {
					options = append(options, forddb.WithResourceIDs(ids...))
				}
			}

			if perPage != -1 {
				options = append(options,
					forddb.WithOffset(pageIndex*perPage),
					forddb.WithLimit(perPage),
				)
			}

			results, err := r.db.List(
				p.Context,
				typ.GetResourceID(),
				options...,
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

	getAllMeta := &graphql.Field{
		Name: getAllMetaName,

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
			res, err := r.db.List(p.Context, typ.GetResourceID())

			if err != nil {
				return nil, err
			}

			return map[string]interface{}{
				"count": len(res),
			}, nil
		},
	}

	ctx.RegisterQuery(getById, getAll, getAllMeta)
}
