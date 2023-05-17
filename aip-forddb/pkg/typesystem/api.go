package typesystem

import (
	"reflect"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/schema"
)

type Type interface {
	Name() TypeName
	PrimitiveKind() PrimitiveKind
	RuntimeType() reflect.Type

	IpldType() schema.Type
	IpldPrimitive() ipld.NodePrototype
	IpldPrototype() schema.TypedPrototype
	IpldRepresentationKind() datamodel.Kind

	Struct() StructType
	List() ListType
	Map() MapType
}

type StructType interface {
	Type

	NumField() int
	Field(name string) Field
	FieldByIndex(index int) Field
}

type ListType interface {
	Type

	Elem() Type
}

type MapType interface {
	Type

	Key() Type
	Value() Type
}

type Field interface {
	Name() string
	Type() Type
	DeclaringType() StructType

	Value(v Value) Value
}

func Universe() *TypeSystem {
	return globalTypeSystem
}

func TypeOf(v interface{}) Type {
	return Universe().LookupByType(reflect.TypeOf(v))
}

func TypeFrom(t reflect.Type) Type {
	return Universe().LookupByType(t)
}

func ValueOf(v interface{}) Value {
	return Value{
		typ: TypeOf(v),
		v:   reflect.ValueOf(v),
	}
}

func ValueFrom(v reflect.Value) Value {
	return Value{
		typ: TypeFrom(v.Type()),
		v:   v,
	}
}

func New(t Type) Value {
	return Value{
		typ: t,
		v:   reflect.New(t.RuntimeType()).Elem(),
	}
}

func MakeList(t Type, length, cap int) Value {
	sliceType := reflect.SliceOf(t.RuntimeType())

	v := Value{
		v: reflect.MakeSlice(sliceType, length, cap),
	}

	v.typ = TypeFrom(v.v.Type())

	return v
}

func MakeMap(k, v Type, length int) Value {
	rt := reflect.MapOf(k.RuntimeType(), v.RuntimeType())
	t := TypeFrom(rt)

	return Value{
		typ: t,
		v:   reflect.MakeMapWithSize(rt, length),
	}
}

func Wrap(v any) schema.TypedNode {
	return ValueOf(v).AsNode().(schema.TypedNode)
}

func Unwrap(v ipld.Node) any {
	val := v.(valueNode).v.Value()

	if !val.IsValid() {
		return nil
	}

	return val.Interface()
}
