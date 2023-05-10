package forddb

import (
	"reflect"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb/nodebinder"
)

type BasicType interface {
	BasicResource

	GetResourceID() TypeID
	Name() string
	Kind() Kind
	PrimitiveKind() PrimitiveKind

	Metadata() TypeMetadata

	RuntimeType() reflect.Type
	IsRuntimeOnly() bool

	CreateInstance() any

	SchemaType() schema.Type
	SchemaPrototype() schema.TypedPrototype
	SchemaLinkPrototype() ipld.LinkPrototype

	TypeSystem() *ResourceTypeSystem

	//Encode(resource any) (RawResource, error)
	//Decode(resource RawResource) (any, error)

	NumFields() int
	Fields() []BasicField
	FieldByName(name string) BasicField
	FieldByIndex(index int) BasicField

	Initialize(ts *ResourceTypeSystem, options ...nodebinder.Option)
}

type Type[T any] interface {
	BasicType
}
