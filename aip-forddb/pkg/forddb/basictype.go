package forddb

import (
	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
)

type basicType struct {
	ResourceBase[TypeID, BasicResourceType] `json:"metadata"`
	typesystem.Type

	metadata TypeMetadata

	universe *ResourceTypeSystem
}

func (bt *basicType) ActualType() typesystem.Type {
	return bt.Type
}

var _ BasicType = (*basicType)(nil)

func newBasicType(
	typ typesystem.Type,
	metadata TypeMetadata,
) *basicType {
	t := &basicType{}

	metadata.ID = TypeID(metadata.Name)
	t.ID = metadata.ID
	t.Type = typ
	t.metadata = metadata

	return t
}

func (bt *basicType) TypeSystem() *ResourceTypeSystem { return bt.universe }
func (bt *basicType) GetResourceID() TypeID           { return bt.ResourceBase.ID }
func (bt *basicType) Kind() Kind                      { return bt.metadata.Kind }
func (bt *basicType) Metadata() TypeMetadata          { return bt.metadata }
func (bt *basicType) IsRuntimeOnly() bool             { return bt.metadata.IsRuntimeOnly }
func (bt *basicType) CreateInstance() any             { return typesystem.New(bt.Type).Value().Interface() }

func (bt *basicType) Initialize(ts *ResourceTypeSystem) {
	bt.universe = ts
}
