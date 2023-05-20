package objectstore

import (
	"context"

	"github.com/ipfs/go-cid"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
)

type IndexDefinition struct {
	Type    typesystem.Type
	Factory IndexFactory
	Fields  []IndexedFieldDefinition
}

type IndexedFieldDefinition struct {
	Field typesystem.Field
	Order forddb.SortOrder
}

type IndexedFieldExpression struct {
	Op    typesystem.OperatorName
	Value typesystem.Value
}

type ObjectIndex interface {
	Definition() IndexDefinition

	Add(ctx context.Context, obj *Object) error
	Remove(ctx context.Context, obj *Object) error

	Find(ctx context.Context, values []IndexedFieldExpression) ([]cid.Cid, error)

	Close() error
}

type IndexFactory interface {
	IndexFor(db forddb.Database, def IndexDefinition) (ObjectIndex, error)
}
