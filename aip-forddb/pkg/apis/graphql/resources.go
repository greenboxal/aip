package graphql

import (
	"context"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/parser"
	"github.com/graphql-go/graphql"
	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"
	"github.com/samber/lo"
	"go.uber.org/fx"
	"golang.org/x/exp/slices"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-forddb/pkg/logstore"
	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
)

type ResourceEventType string

const (
	ResourceEventTypeInvalid ResourceEventType = ""
	ResourceEventTypeCreated ResourceEventType = "created"
	ResourceEventTypeUpdated ResourceEventType = "updated"
	ResourceEventTypeDeleted ResourceEventType = "deleted"
)

type ResourceEvent struct {
	Type ResourceEventType `json:"type"`

	Payload ResourceEventPayload `json:"payload"`
}

type ResourceEventPayload struct {
	IDs []string `json:"ids"`
}

type DatabaseResourceBinding struct {
	db forddb.Database
	sm *SubscriptionManager
}

func NewDatabaseResourceBinding(
	db forddb.Database,
	sm *SubscriptionManager,
) *DatabaseResourceBinding {
	return &DatabaseResourceBinding{
		db: db,
		sm: sm,
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

	resourceChangedEventType := ctx.LookupOutputType(typesystem.TypeOf(&ResourceEvent{}))

	ctx.RegisterSubscription(&graphql.Field{
		Name: "resourceChanged",
		Type: resourceChangedEventType,

		Args: graphql.FieldConfigArgument{
			"resourceType": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},

		Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
			resourceTypeName := p.Args["resourceType"].(string)
			resourceType := forddb.LookupTypeByName(resourceTypeName)

			return r.sm.Subscribe(p.Context, resourceType)
		},
	})
}

func (r *DatabaseResourceBinding) compileResource(ctx BindingContext, typ forddb.BasicResourceType) {
	name := typ.ActualType().Name().Name
	name = regexp.MustCompile("[^_a-zA-Z0-9]").ReplaceAllString(name, "_")

	gqlType := ctx.LookupOutputType(typ.ActualType())

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
			res, err := r.db.Get(p.Context, id.BasicResourceType().GetResourceID(), id)

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

	filterFields := graphql.InputObjectConfigFieldMap{
		"q": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"id": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"ids": &graphql.InputObjectFieldConfig{
			Type: graphql.NewList(graphql.String),
		},
	}

	for _, f := range typ.FilterableFields() {
		for _, op := range f.Operators {
			var suffix string

			if op != "==" {
				mapped, ok := operatorMap[op]

				if !ok {
					continue
				}

				suffix = "_" + mapped
			}

			filterFields[f.Field.Name()+suffix] = &graphql.InputObjectFieldConfig{
				Type: ctx.LookupInputType(f.Field.Type()),
			}
		}
	}

	filterType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:   typ.ResourceName().ToTitle() + "Filter",
		Fields: filterFields,
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
				var q string

				filter := filterVal.(map[string]interface{})

				if id, ok := filter["id"]; ok {
					ids = append(ids, typ.CreateID(id.(string)))
				}

				if stringIds, ok := filter["ids"]; ok {
					for _, id := range stringIds.([]interface{}) {
						ids = append(ids, typ.CreateID(id.(string)))
					}
				}

				if qVal, ok := filter["q"]; ok {
					q = qVal.(string)
				}

				var filters []ast.Node

				for key, val := range filter {
					targetOp := "=="
					fieldName := key
					components := strings.Split(key, "_")

					if key == "id" || key == "ids" || key == "q" {
						continue
					}

					if len(components) > 1 {
						op := components[len(components)-1]
						mappedOp, ok := reverseOperatorMap[op]

						if !ok {
							continue
						}

						targetOp = mappedOp
						fieldName = strings.Join(components[:len(components)-1], "_")
					}

					filters = append(filters, &ast.BinaryNode{
						Left: &ast.MemberNode{
							Name: fieldName,

							Property: &ast.IdentifierNode{Value: fieldName},
							Node:     &ast.IdentifierNode{Value: "resource"},
						},

						Operator: targetOp,

						Right: &ast.ConstantNode{
							Value: val,
						},
					})
				}

				if q != "" {
					parsed, err := parser.Parse(q)

					if err != nil {
						return nil, err
					}

					filters = append(filters, parsed.Node)
				}

				for _, id := range ids {
					filters = append(filters, &ast.BinaryNode{
						Left: &ast.MemberNode{
							Name:     "id",
							Property: &ast.IdentifierNode{Value: "id"},

							Node: &ast.MemberNode{
								Name:     "metadata",
								Property: &ast.IdentifierNode{Value: "metadata"},
								Node:     &ast.IdentifierNode{Value: "resource"},
							},
						},

						Operator: "==",
						Right:    &ast.StringNode{Value: id.String()},
					})
				}

				if len(filters) > 0 {
					var rootNode ast.Node

					for i, v := range filters {
						if i == 0 {
							rootNode = v

							continue
						}

						rootNode = &ast.BinaryNode{
							Left:     rootNode,
							Operator: "&&",
							Right:    v,
						}
					}

					options = append(options, forddb.WithListQueryOptions(forddb.WithFilterExpressionNode(rootNode)))
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

			return lo.Map(results, func(item forddb.BasicResource, _index int) any {
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

type SubscriptionManager struct {
	db forddb.Database

	m             sync.RWMutex
	subscriptions map[forddb.BasicResourceType][]chan ResourceEvent
}

func NewSubscriptionManager(
	lc fx.Lifecycle,
	db forddb.Database,
) *SubscriptionManager {
	sm := &SubscriptionManager{
		db:            db,
		subscriptions: map[forddb.BasicResourceType][]chan ResourceEvent{},
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			goprocess.Go(sm.Run)

			return nil
		},
	})

	return sm
}

func (sm *SubscriptionManager) Run(proc goprocess.Process) {
	ctx := goprocessctx.OnClosingContext(proc)

	iterator := sm.db.LogStore().Iterator(forddb.WithBlockingIterator())

	for {
		if !iterator.Next(ctx) {
			if err := iterator.SetLSN(ctx, forddb.MakeLSN(logstore.FileSegmentBaseSeekToHead, time.Now())); err != nil {
				return
			}
			continue
		}

		err := iterator.Error()

		if err != nil {
			panic(iterator.Error())
		}

		record := iterator.Record()

		sm.dispatch(record)
	}
}

func (sm *SubscriptionManager) Subscribe(
	ctx context.Context,
	typ forddb.BasicResourceType,
) (<-chan ResourceEvent, error) {
	ch := make(chan ResourceEvent, 128)

	sm.addSubscription(typ, ch)

	go func() {
		//defer close(ch)
		//defer sm.removeSubscription(typ, ch)

		<-ctx.Done()
	}()

	return ch, nil
}

func (sm *SubscriptionManager) addSubscription(typ forddb.BasicResourceType, ch chan ResourceEvent) {
	sm.m.Lock()
	defer sm.m.Unlock()

	subs := sm.subscriptions[typ]

	index := slices.Index(subs, ch)

	if index != -1 {
		return
	}

	index = slices.Index(subs, nil)

	if index != -1 {
		subs[index] = ch
	} else {
		subs = append(subs, ch)
	}

	sm.subscriptions[typ] = subs
}

func (sm *SubscriptionManager) removeSubscription(typ forddb.BasicResourceType, ch chan ResourceEvent) {
	sm.m.Lock()
	defer sm.m.Unlock()

	subs := sm.subscriptions[typ]

	if len(subs) == 0 {
		return
	}

	index := slices.Index(subs, ch)

	if index == -1 {
		return
	}

	subs = slices.Delete(subs, index, index+1)

	sm.subscriptions[typ] = subs
}

func (sm *SubscriptionManager) dispatch(record *forddb.LogEntryRecord) {
	event := ResourceEvent{}

	switch record.Kind {
	case forddb.LogEntryKindSet:
		if record.Previous == nil {
			event.Type = ResourceEventTypeCreated
		} else {
			event.Type = ResourceEventTypeUpdated
		}
	case forddb.LogEntryKindDelete:
		event.Type = ResourceEventTypeDeleted
	}

	event.Payload.IDs = []string{record.ID}

	sm.m.RLock()
	defer sm.m.RUnlock()

	subs := sm.subscriptions[record.Type.Type()]

	if len(subs) == 0 {
		return
	}

	for _, sub := range subs {
		if sub == nil {
			continue
		}

		sub <- event
	}
}

var operatorMap = map[string]string{
	"==": "eq",
	"!=": "neq",
	"<":  "lt",
	"<=": "lte",
	">":  "gt",
	">=": "gte",
}

var reverseOperatorMap = map[string]string{}

func init() {
	for k, v := range operatorMap {
		reverseOperatorMap[v] = k
	}
}
