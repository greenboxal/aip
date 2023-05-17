package forddb

import (
	"reflect"

	"github.com/greenboxal/aip/aip-forddb/pkg/impl/nodebinder"
	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
)

type BasicType interface {
	BasicResource

	ActualType() typesystem.Type

	GetResourceID() TypeID
	Kind() Kind

	Metadata() TypeMetadata

	RuntimeType() reflect.Type
	IsRuntimeOnly() bool

	CreateInstance() any

	TypeSystem() *ResourceTypeSystem

	Initialize(ts *ResourceTypeSystem, options ...nodebinder.Option)
}

type Type[T any] interface {
	BasicType
}
