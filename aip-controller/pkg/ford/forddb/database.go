package forddb

import (
	"context"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/ast"
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
	AsExpr() *vm.Program
	AsAst() ast.Node
}

type ProgramFilter struct{ *vm.Program }

func (p ProgramFilter) AsExpr() *vm.Program {
	return p.Program
}

func (p ProgramFilter) AsAst() ast.Node {
	return p.Program.Node
}

type QueryOptions struct {
	ResourceType    TypeID
	ReadConsistency ReadConsistencyLevel

	FilterExpression Filter
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
