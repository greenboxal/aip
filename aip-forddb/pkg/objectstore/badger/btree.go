package badger

import (
	"context"
	"fmt"
	"strings"

	"github.com/dgraph-io/badger/v4"
	"github.com/ipfs/go-cid"
	"github.com/samber/lo"

	"github.com/greenboxal/aip/aip-forddb/pkg/objectstore"
	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
)

type BtreeIndex struct {
	db  *badger.DB
	def objectstore.IndexDefinition

	prefix string
}

func NewBtreeIndex(db *badger.DB, def objectstore.IndexDefinition) *BtreeIndex {
	fieldNames := strings.Join(lo.Map(def.Fields, func(f objectstore.IndexedFieldDefinition, _ int) string {
		return f.Field.Name()
	}), "+")

	prefix := fmt.Sprintf(
		"/indexes/btree/%s/%s",
		def.Type.Name().NormalizedFullNameWithArguments(),
		fieldNames,
	)

	return &BtreeIndex{
		db:     db,
		def:    def,
		prefix: prefix,
	}
}

func (b *BtreeIndex) keyForObject(obj *objectstore.Object) ([]byte, error) {
	v := typesystem.ValueOf(obj.Node)
	fields := make([]string, len(b.def.Fields))

	for i, f := range b.def.Fields {
		val, err := f.Field.Resolve(v).AsNode().AsString()

		if err != nil {
			return nil, err
		}

		fields[i] = val
	}

	key := strings.Join(fields, "\000")

	return []byte(fmt.Sprintf("%s/%s", b.prefix, key)), nil
}

func (b *BtreeIndex) valueForObject(obj *objectstore.Object) ([]byte, error) {
	meta, err := obj.GetResourceMetadata()

	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("%s", meta.GetResourceBasicID())), nil
}

func (b *BtreeIndex) Definition() objectstore.IndexDefinition {
	return b.def
}

func (b *BtreeIndex) Add(ctx context.Context, obj *objectstore.Object) error {
	key, err := b.keyForObject(obj)

	if err != nil {
		return err
	}

	value, err := b.valueForObject(obj)

	if err != nil {
		return err
	}

	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (b *BtreeIndex) Remove(ctx context.Context, obj *objectstore.Object) error {
	key, err := b.keyForObject(obj)

	if err != nil {
		return err
	}

	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func (b *BtreeIndex) Find(
	ctx context.Context,
	values []objectstore.IndexedFieldExpression,
) (result []cid.Cid, err error) {
	valuesStr := make([]string, len(values))

	for i, v := range values {
		str, err := v.Value.AsNode().AsString()

		if err != nil {
			return nil, err
		}

		valuesStr[i] = str
	}

	err = b.db.View(func(txn *badger.Txn) error {

		return nil
	})

	return
}

func (b *BtreeIndex) Close() error {
	return nil
}
