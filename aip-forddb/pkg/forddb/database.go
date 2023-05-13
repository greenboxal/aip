package forddb

import (
	"context"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/compiler"
	"github.com/antonmedv/expr/conf"
	"github.com/antonmedv/expr/file"
	"github.com/antonmedv/expr/optimizer"
	"github.com/antonmedv/expr/parser"
	"github.com/antonmedv/expr/vm"
)

type Database interface {
	HasListeners

	LogStore() LogStore

	List(ctx context.Context, typ TypeID, options ...ListOption) ([]BasicResource, error)
	Get(ctx context.Context, typ TypeID, id BasicResourceID, options ...GetOption) (BasicResource, error)
	Put(ctx context.Context, resource BasicResource, options ...PutOption) (BasicResource, error)
	Delete(ctx context.Context, resource BasicResource, options ...DeleteOption) (BasicResource, error)
}

type Filter interface {
	Evaluate(value any, args any) (bool, error)

	AsExpr() *vm.Program
	AsAst() ast.Node

	String() string
}

type ProgramFilter struct{ *vm.Program }

func (p ProgramFilter) AsExpr() *vm.Program {
	return p.Program
}

func (p ProgramFilter) AsAst() ast.Node {
	return p.Program.Node
}

func (p ProgramFilter) String() string {
	return p.Program.Disassemble()
}

func (p ProgramFilter) Evaluate(value any, args any) (bool, error) {
	env := map[string]interface{}{
		"resource": value,
		"args":     args,
	}

	result, err := expr.Run(p.Program, env)

	if err != nil {
		return false, err
	}

	return result.(bool), nil
}

type QueryOptions struct {
	ResourceType    TypeID
	ReadConsistency ReadConsistencyLevel

	FilterExpression Filter
	FilterParameters map[string]interface{}
}

type QueryOption func(opts *QueryOptions)

func WithReadConsistency(level ReadConsistencyLevel) QueryOption {
	return func(opts *QueryOptions) {
		opts.ReadConsistency = level
	}
}

type GetOptions struct {
	QueryOptions
}

func NewGetOptions(typ TypeID, opts ...GetOption) GetOptions {
	var result GetOptions

	result.ResourceType = typ

	for _, opt := range opts {
		opt(&result)
	}

	return result
}

type GetOption func(opts *GetOptions)

func WithGetQueryOptions(options ...QueryOption) GetOption {
	return func(opts *GetOptions) {
		for _, opt := range options {
			opt(&opts.QueryOptions)
		}
	}
}

func WithFilterParameters(parameters map[string]interface{}) QueryOption {
	return func(opts *QueryOptions) {
		if opts.FilterParameters == nil {
			opts.FilterParameters = map[string]interface{}{}
		}

		for k, v := range parameters {
			opts.FilterParameters[k] = v
		}
	}
}

func WithFilterParameter(name string, value any) QueryOption {
	return WithFilterParameters(map[string]interface{}{
		name: value,
	})
}

func WithFilterExpressionNode(node ast.Node) QueryOption {
	return func(opts *QueryOptions) {
		var err error

		tree := &parser.Tree{
			Node: node,
		}

		config := conf.CreateNew()
		config.Check()

		if len(config.Operators) > 0 {
			config.Visitors = append(config.Visitors, &conf.OperatorPatcher{
				Operators: config.Operators,
				Types:     config.Types,
			})
		}

		if len(config.Visitors) > 0 {
			for _, v := range config.Visitors {
				// We need to perform types check, because some visitors may rely on
				// types information available in the tree.
				_, _ = checker.Check(tree, config)
				ast.Walk(&tree.Node, v)
			}
			_, err = checker.Check(tree, config)
			if err != nil {
				panic(err)
			}
		} else {
			_, err = checker.Check(tree, config)
			if err != nil {
				panic(err)
			}
		}

		if config.Optimize {
			err = optimizer.Optimize(&tree.Node, config)
			if err != nil {
				if fileError, ok := err.(*file.Error); ok {
					panic(fileError.Bind(tree.Source))
				}

				panic(err)
			}
		}

		program, err := compiler.Compile(tree, config)
		if err != nil {
			panic(err)
		}

		opts.FilterExpression = ProgramFilter{program}
	}
}
func WithFilterExpression(q string) QueryOption {
	return func(opts *QueryOptions) {
		program, err := expr.Compile(q)

		if err != nil {
			panic(err)
		}

		opts.FilterExpression = ProgramFilter{program}
	}
}

type DeleteOptions struct {
}

type DeleteOption func(opts *DeleteOptions)

type SortField struct {
	Path  string    `json:"field"`
	Order SortOrder `json:"order"`
}

type ListOptions struct {
	QueryOptions

	Offset     int
	Limit      int
	SortFields []SortField

	ResourceIDs []BasicResourceID
}

func NewListOptions(typ TypeID, opts ...ListOption) ListOptions {
	var result ListOptions

	result.ResourceType = typ

	for _, opt := range opts {
		opt(&result)
	}

	return result
}

type ListOption func(opts *ListOptions)

func WithListQueryOptions(options ...QueryOption) ListOption {
	return func(opts *ListOptions) {
		for _, opt := range options {
			opt(&opts.QueryOptions)
		}
	}
}

func WithSortField(field string, order SortOrder) ListOption {
	return func(opts *ListOptions) {
		opts.SortFields = append(opts.SortFields, SortField{
			Path:  field,
			Order: order,
		})
	}
}

func WithSortFields(fields ...SortField) ListOption {
	return func(opts *ListOptions) {
		opts.SortFields = append(opts.SortFields, fields...)
	}
}

func WithResourceIDs(ids ...BasicResourceID) ListOption {
	return func(opts *ListOptions) {
		opts.ResourceIDs = append(opts.ResourceIDs, ids...)
	}
}

func WithOffset(offset int) ListOption {
	return func(opts *ListOptions) {
		opts.Offset = offset
	}
}

func WithLimit(limit int) ListOption {
	return func(opts *ListOptions) {
		opts.Limit = limit
	}
}

type PutOptions struct {
	OnConflict OnConflict
}

type PutOption func(opts *PutOptions)

func WithOnConflict(onConflict OnConflict) PutOption {
	return func(opts *PutOptions) {
		opts.OnConflict = onConflict
	}
}

func NewPutOptions(options ...PutOption) PutOptions {
	opts := PutOptions{
		OnConflict: OnConflictOptimistic,
	}

	for _, option := range options {
		option(&opts)
	}

	return opts
}
