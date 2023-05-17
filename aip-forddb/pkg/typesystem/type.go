package typesystem

import (
	"reflect"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/schema"
)

type typeInitializer interface {
	initialize(ts *TypeSystem)
}

type basicType struct {
	self Type

	name                   TypeName
	primitiveKind          PrimitiveKind
	runtimeType            reflect.Type
	ipldType               schema.Type
	ipldPrimitive          ipld.NodePrototype
	ipldPrototype          schema.TypedPrototype
	ipldRepresentationKind datamodel.Kind
}

func (bt *basicType) Name() TypeName                         { return bt.name }
func (bt *basicType) PrimitiveKind() PrimitiveKind           { return bt.primitiveKind }
func (bt *basicType) RuntimeType() reflect.Type              { return bt.runtimeType }
func (bt *basicType) Struct() StructType                     { return bt.self.(StructType) }
func (bt *basicType) List() ListType                         { return bt.self.(ListType) }
func (bt *basicType) Map() MapType                           { return bt.self.(MapType) }
func (bt *basicType) IpldType() schema.Type                  { return bt.ipldType }
func (bt *basicType) IpldPrimitive() ipld.NodePrototype      { return bt.ipldPrimitive }
func (bt *basicType) IpldPrototype() schema.TypedPrototype   { return bt.ipldPrototype }
func (bt *basicType) IpldRepresentationKind() datamodel.Kind { return bt.ipldRepresentationKind }

type scalarType struct {
	basicType
}

type interfaceType struct {
	basicType
}

type structType struct {
	basicType

	fields   []Field
	fieldMap map[string]Field
}

func (st *structType) NumField() int                { return len(st.fields) }
func (st *structType) Field(name string) Field      { return st.fieldMap[name] }
func (st *structType) FieldByIndex(index int) Field { return st.fields[index] }

type listType struct {
	basicType

	elem Type
}

func (lt *listType) Elem() Type { return lt.elem }

type mapType struct {
	basicType

	key Type
	val Type
}

func (mt *mapType) Key() Type   { return mt.key }
func (mt *mapType) Value() Type { return mt.val }

type field struct {
	declaringType StructType
	name          string
	typ           Type
	runtimeField  reflect.StructField
}

func (f *field) Name() string              { return f.name }
func (f *field) Type() Type                { return f.typ }
func (f *field) DeclaringType() StructType { return f.declaringType }

func (f *field) Value(receiver Value) Value {
	v := reflect.Indirect(receiver.Value()).FieldByIndex(f.runtimeField.Index)

	return ValueFrom(v).As(f.typ)
}
