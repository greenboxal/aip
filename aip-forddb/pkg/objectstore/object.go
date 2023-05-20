package objectstore

import (
	"github.com/ipld/go-ipld-prime"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
)

type ObjectID string

type Object struct {
	ipld.Node
}

func (o *Object) GetResourceMetadata() (forddb.ResourceMetadata, error) {
	node, err := o.LookupByString("@metadata")

	if err != nil {
		return nil, err
	}

	res, ok := typesystem.TryUnwrap[forddb.ResourceMetadata](node)

	if !ok {
		return nil, typesystem.ErrInvalidType
	}

	return res, nil
}
