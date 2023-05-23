package firestore

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/ast"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type Config struct {
	ProjectID    string
	CollectionID string
}

func (c *Config) SetDefaults() {
	if c.ProjectID == "" {
		c.ProjectID = "uncyclo-385820"
	}

	if c.CollectionID == "" {

	}
}

type Storage struct {
	config *Config
	client *firestore.Client
}

func NewStorage(config *Config) (*Storage, error) {
	ctx := context.Background()

	conf := &firebase.Config{ProjectID: config.ProjectID}
	app, err := firebase.NewApp(ctx, conf)

	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)

	if err != nil {
		return nil, err
	}

	return &Storage{
		client: client,
	}, nil
}

func (s *Storage) List(
	ctx context.Context,
	typ forddb.TypeID,
	opts forddb.ListOptions,
) ([]forddb.RawResource, error) {
	collection := typ.Name()

	col := s.client.Collection(collection)

	query := col.Query

	if opts.Offset > 0 {
		query = col.Offset(opts.Offset)
	}

	if opts.Limit > 0 {
		query = col.Limit(opts.Limit)
	}

	for _, item := range opts.SortFields {
		itemPath := item.Path

		if itemPath == "id" {
			itemPath = "metadata.id"
		}

		if item.Order == forddb.Asc {
			query = query.OrderBy(itemPath, firestore.Asc)
		} else {
			query = query.OrderBy(itemPath, firestore.Desc)
		}
	}

	if opts.FilterExpression != nil {
		node := opts.FilterExpression.AsAst()
		conditions, err := parseConditions(node, opts.FilterParameters)

		if err != nil {
			return nil, err
		}

		for _, condition := range conditions {
			query = query.WherePath(condition.Path, condition.Op, condition.Val)
		}
	}

	iterator := query.Documents(ctx)

	all, err := iterator.GetAll()

	if err != nil {
		return nil, err
	}

	raws := make([]forddb.RawResource, len(all))

	for i, v := range all {
		data := v.Data()

		serialized, err := json.Marshal(data)

		if err != nil {
			return nil, err
		}

		node, err := ipld.DecodeUsingPrototype(serialized, dagjson.Decode, typ.Type().ActualType().IpldPrototype())

		if err != nil {
			return nil, err
		}

		raws[i] = node
	}

	return raws, nil
}

type parsedCondition struct {
	Path firestore.FieldPath
	Op   string
	Val  interface{}
}

func parseConditions(node ast.Node, args any) ([]parsedCondition, error) {
	var walkExpression func(node ast.Node) error
	var walkValue func(node ast.Node) (any, error)
	var conditions []parsedCondition

	walkPath := func(node ast.Node, root string) (firestore.FieldPath, error) {
		var path firestore.FieldPath

		for node != nil {
			switch n := node.(type) {
			case *ast.IdentifierNode:
				if n.Value == root {
					node = nil
				} else {
					return nil, fmt.Errorf("unsupported node: %#v", n)
				}

			case *ast.MemberNode:
				var name string

				switch v := n.Property.(type) {
				case *ast.IdentifierNode:
					name = v.Value

				case *ast.IntegerNode:
					name = strconv.Itoa(v.Value)

				case *ast.StringNode:
					name = v.Value

				default:
					return nil, fmt.Errorf("unsupported node: %#v", n)
				}

				path = slices.Insert(path, 0, name)

				node = n.Node

			default:
				return nil, fmt.Errorf("unsupported node: %#v", n)
			}
		}

		return path, nil
	}

	walkValue = func(node ast.Node) (any, error) {
		switch n := node.(type) {
		case *ast.ConstantNode:
			if m, ok := n.Value.(map[string]struct{}); ok {
				return maps.Keys(m), nil
			}

			return n.Value, nil

		case *ast.MemberNode:
			path, err := walkPath(node, "args")

			if err != nil {
				return nil, err
			}

			s := strings.Join(path, ".")

			return expr.Eval(s, args)

		case *ast.IdentifierNode:
			path, err := walkPath(node, "args")

			if err != nil {
				return nil, err
			}

			s := strings.Join(path, ".")

			return expr.Eval(s, args)

		case *ast.StringNode:
			return n.Value, nil

		case *ast.IntegerNode:
			return n.Value, nil

		case *ast.FloatNode:
			return n.Value, nil

		case *ast.BoolNode:
			return n.Value, nil

		case *ast.ArrayNode:
			var values []any

			for _, item := range n.Nodes {
				value, err := walkValue(item)

				if err != nil {
					return nil, err
				}

				values = append(values, value)
			}

			return values, nil

		case *ast.MapNode:
			var values map[any]any

			for i := 0; i < len(n.Pairs); i += 2 {
				k, err := walkValue(n.Pairs[i])

				if err != nil {
					return nil, err
				}

				v, err := walkValue(n.Pairs[i+1])

				if err != nil {
					return nil, err
				}

				values[k] = v
			}

			return values, nil

		default:
			return nil, fmt.Errorf("unsupported node: %#v", node)
		}
	}

	walkBinOp := func(node *ast.BinaryNode) error {
		path, err := walkPath(node.Left, "resource")

		if err != nil {
			return err
		}

		value, err := walkValue(node.Right)

		if err != nil {
			return err
		}

		conditions = append(conditions, parsedCondition{
			Path: path,
			Op:   node.Operator,
			Val:  value,
		})

		return nil
	}

	walkUnOp := func(node *ast.UnaryNode) error {
		return fmt.Errorf("unsupported node: %#v", node)
	}

	walkExpression = func(node ast.Node) error {
		switch n := node.(type) {
		case *ast.BinaryNode:
			if n.Operator == "&&" {
				if err := walkExpression(n.Left); err != nil {
					return err
				}

				if err := walkExpression(n.Left); err != nil {
					return err
				}

				return nil
			} else {
				return walkBinOp(n)
			}

		case *ast.UnaryNode:
			return walkUnOp(n)

		default:
			return fmt.Errorf("unsupported node: %#v", node)
		}
	}

	if err := walkExpression(node); err != nil {
		return nil, err
	}

	return conditions, nil
}

func (s *Storage) Get(
	ctx context.Context,
	typ forddb.TypeID,
	id forddb.BasicResourceID,
	opts forddb.GetOptions,
) (forddb.RawResource, error) {
	collection := typ.Name()

	col := s.client.Collection(collection)
	doc := col.Doc(id.String())

	result, err := doc.Get(ctx)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, forddb.ErrNotFound
		}

		return nil, err
	}

	fields := result.Data()

	// Fixme
	delete(fields, "kind")

	serialized, err := json.Marshal(fields)

	if err != nil {
		return nil, err
	}

	node, err := ipld.DecodeUsingPrototype(serialized, dagjson.Decode, typ.Type().ActualType().IpldPrototype())

	if err != nil {
		return nil, err
	}

	return node, nil
}

func (s *Storage) Put(
	ctx context.Context,
	resource forddb.RawResource,
	opts forddb.PutOptions,
) (forddb.RawResource, error) {
	var fields map[string]interface{}
	var preconditions []firestore.Precondition

	serialized, err := ipld.Encode(resource, dagjson.Encode)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(serialized, &fields); err != nil {
		return nil, err
	}

	unk := forddb.UnknownResource{RawResource: resource}
	collection := unk.GetResourceTypeID().Name()
	col := s.client.Collection(collection)
	doc := col.Doc(unk.GetResourceBasicID().String())

	version := unk.GetResourceVersion()

	if version >= 0 {
		if r, err := doc.Set(ctx, fields); err != nil {
			return nil, err
		} else {
			_ = r
		}
	} else {
		switch opts.OnConflict {
		case forddb.OnConflictReplace:
			if r, err := doc.Set(ctx, fields); err != nil {
				return nil, err
			} else {
				_ = r
			}

		default:
			updates := make([]firestore.Update, 0, len(fields))

			for k, v := range fields {
				updates = append(updates, firestore.Update{
					Path:  k,
					Value: v,
				})
			}

			/*if !metadata.UpdatedAt.IsZero() && metadata.Version > 1 {
				preconditions = append(preconditions, firestore.LastUpdateTime(metadata.UpdatedAt))
			}*/

			if r, err := doc.Update(ctx, updates, preconditions...); err != nil {
				return nil, err
			} else {
				_ = r
			}
		}
	}

	return resource, nil
}

func (s *Storage) Delete(
	ctx context.Context,
	resource forddb.RawResource,
	opts forddb.DeleteOptions,
) (forddb.RawResource, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) Close() error {
	return s.client.Close()
}
