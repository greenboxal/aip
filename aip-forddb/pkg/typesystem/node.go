package typesystem

import (
	"encoding"
	"errors"
	"reflect"
	"strconv"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/schema"
)

type valueNode struct {
	v Value
}

func (n valueNode) Kind() datamodel.Kind {
	if n.v.Type().PrimitiveKind() == PrimitiveKindInterface {
		if n.v.Value().IsNil() {
			return datamodel.Kind_Null
		}

		if n.v.Value().Type().Implements(reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()) {
			return datamodel.Kind_String
		}

		if n.v.Value().Type().Implements(reflect.TypeOf((*encoding.BinaryMarshaler)(nil)).Elem()) {
			return datamodel.Kind_Bytes
		}

		return datamodel.Kind_Map
	}

	r := n.v.typ.IpldType().TypeKind().ActsLike()

	return r
}

func (n valueNode) LookupByString(key string) (datamodel.Node, error) {
	switch n.v.Type().PrimitiveKind() {
	case PrimitiveKindList:
		index, err := strconv.ParseInt(key, 10, 64)

		if err != nil {
			return nil, err
		}

		return n.LookupByIndex(index)

	case PrimitiveKindMap:
		v := n.v.Value().MapIndex(reflect.ValueOf(key))

		return ValueFrom(v).As(n.v.typ.Map().Value()).AsNode(), nil

	case PrimitiveKindStruct:
		st := n.v.typ.Struct()
		f := st.Field(key)

		if f == nil {
			return nil, errors.New("field not found")
		}

		return n.v.GetField(f).AsNode(), nil

	default:
		panic("invalid type")
	}
}

func (n valueNode) LookupByNode(key datamodel.Node) (datamodel.Node, error) {
	//TODO implement me
	panic("implement me")
}

func (n valueNode) LookupByIndex(idx int64) (datamodel.Node, error) {
	switch n.v.Type().PrimitiveKind() {
	case PrimitiveKindList:
		v := n.v.Value().Index(int(idx))

		return ValueFrom(v).As(n.v.typ.List().Elem()).AsNode(), nil

	case PrimitiveKindStruct:
		st := n.v.Type().Struct()
		f := st.FieldByIndex(int(idx))

		if f == nil {
			return nil, errors.New("field not found")
		}

		return n.v.GetField(f).AsNode(), nil

	default:
		panic("invalid type")
	}
}

func (n valueNode) LookupBySegment(seg datamodel.PathSegment) (datamodel.Node, error) {
	i, err := seg.Index()

	if err == nil && i >= 0 {
		return n.LookupByIndex(i)
	}

	return n.LookupByString(seg.String())
}

func (n valueNode) MapIterator() datamodel.MapIterator {
	switch n.v.typ.PrimitiveKind() {
	case PrimitiveKindMap:
		return &mapIterator{v: n.v}
	case PrimitiveKindStruct:
		return &structIterator{v: n.v, t: n.v.typ.(StructType)}
	case PrimitiveKindInterface:
		return &interfaceIterator{v: n.v}
	}

	panic("invalid type")
}

func (n valueNode) ListIterator() datamodel.ListIterator {
	return &listIterator{v: n.v}
}

func (n valueNode) Length() int64 {
	switch n.v.typ.PrimitiveKind() {
	case PrimitiveKindList:
		fallthrough
	case PrimitiveKindMap:
		fallthrough
	case PrimitiveKindBytes:
		fallthrough
	case PrimitiveKindString:
		return int64(n.v.Value().Len())
	case PrimitiveKindStruct:
		return int64(n.v.typ.(StructType).NumField())
	case PrimitiveKindInterface:
		if n.v.Value().IsNil() {
			return 0
		}

		c := 1
		v := ValueFrom(n.v.Value().Elem())
		l := v.AsNode().Length()

		if l >= 0 {
			c += int(l)
		}

		return int64(c)
	}

	return -1
}

func (n valueNode) IsAbsent() bool {
	return !n.v.Value().IsValid() || n.v.Value().IsZero()
}

func (n valueNode) IsNull() bool {
	if !n.v.Value().IsValid() {
		return true
	}

	switch n.v.Value().Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Map, reflect.Slice:
		return n.v.Value().IsNil()

	default:
		return false
	}
}

func (n valueNode) AsBool() (bool, error) {
	if n.v.Value().Kind() != reflect.Bool {
		return false, errors.New("cannot convert to bool")
	}

	return n.v.Value().Bool(), nil
}

func (n valueNode) AsInt() (int64, error) {
	if n.v.Type().PrimitiveKind() == PrimitiveKindUnsignedInt {
		return int64(n.v.Value().Uint()), nil
	}

	return n.v.Value().Int(), nil
}

func (n valueNode) AsFloat() (float64, error) {
	return n.v.Value().Float(), nil
}

func (n valueNode) AsString() (string, error) {
	if m, ok := TryCast[encoding.TextMarshaler](n.v.Value()); ok {
		str, err := m.MarshalText()

		if err != nil {
			return "", err
		}

		return string(str), nil
	}

	switch n.v.typ.PrimitiveKind() {
	case PrimitiveKindString:
		return n.v.Value().String(), nil

	case PrimitiveKindInt:
		return strconv.FormatInt(n.v.Value().Int(), 10), nil

	case PrimitiveKindUnsignedInt:
		return strconv.FormatUint(n.v.Value().Uint(), 10), nil

	case PrimitiveKindFloat:
		return strconv.FormatFloat(n.v.Value().Float(), 'f', -1, 64), nil

	case PrimitiveKindBytes:
		return string(n.v.Value().Bytes()), nil
	}

	return n.v.Value().String(), nil
}

func (n valueNode) AsBytes() ([]byte, error) {
	k := n.v.Value().Kind()

	if k == reflect.String {
		return []byte(n.v.Value().String()), nil
	} else if k == reflect.Slice && n.v.Value().Type().Elem().Kind() == reflect.Uint8 {
		return n.v.Value().Bytes(), nil
	}

	if m, ok := TryCast[encoding.BinaryMarshaler](n.v.Value()); ok {
		return m.MarshalBinary()
	}

	return nil, errors.New("cannot convert to bytes")
}

func (n valueNode) AsLink() (datamodel.Link, error) {
	return n.v.Value().Interface().(datamodel.Link), nil
}

func (n valueNode) Prototype() datamodel.NodePrototype {
	return n.v.Type().IpldPrototype()
}

func (n valueNode) Type() schema.Type {
	return n.v.Type().IpldType()
}

func (n valueNode) Representation() datamodel.Node {
	return reprNode{n}
}

type reprNode struct {
	valueNode
}

func reprStrategy(typ schema.Type) interface{} {
	// Can't use an interface check, as each method has a different result type.
	// TODO: consider inlining this type switch at each call site,
	// as the call sites need the underlying schema.Type too.
	switch typ := typ.(type) {
	case *schema.TypeStruct:
		return typ.RepresentationStrategy()
	case *schema.TypeUnion:
		return typ.RepresentationStrategy()
	case *schema.TypeEnum:
		return typ.RepresentationStrategy()
	}
	return nil
}

func (n reprNode) Kind() datamodel.Kind {
	st := n.v.Type().IpldType()
	switch reprStrategy(st).(type) {
	case schema.StructRepresentation_Stringjoin:
		return datamodel.Kind_String
	case schema.StructRepresentation_Map:
		return datamodel.Kind_Map
	case schema.StructRepresentation_Tuple:
		return datamodel.Kind_List
	case schema.UnionRepresentation_Keyed:
		return datamodel.Kind_Map
	case schema.UnionRepresentation_Stringprefix:
		return datamodel.Kind_String
	case schema.EnumRepresentation_Int:
		return datamodel.Kind_Int
	case schema.EnumRepresentation_String:
		return datamodel.Kind_String
	default:
		return n.valueNode.Kind()
	}
}

type listIterator struct {
	v     Value
	index int
}

func (l *listIterator) Next() (idx int64, value datamodel.Node, err error) {
	i := l.index

	if i >= l.v.Value().Len() {
		return -1, nil, datamodel.ErrIteratorOverread{}
	}

	l.index++

	v := l.v.Value().Index(i)

	return int64(i), ValueFrom(v).AsNode(), nil
}

func (l *listIterator) Done() bool {
	return l.index >= l.v.Value().Len()
}

type mapIterator struct {
	v     Value
	keys  []reflect.Value
	index int
}

func (m *mapIterator) Next() (key datamodel.Node, value datamodel.Node, err error) {
	if m.keys == nil {
		m.keys = m.v.Value().MapKeys()
		m.index = 0
	}

	i := m.index

	if i >= len(m.keys) {
		return nil, nil, datamodel.ErrIteratorOverread{}
	}

	m.index++

	k := m.keys[i]
	v := m.v.Value().MapIndex(k)

	return ValueFrom(k).AsNode(), ValueFrom(v).AsNode(), nil
}

func (m *mapIterator) Done() bool {
	if m.keys == nil {
		m.keys = m.v.Value().MapKeys()
		m.index = 0
	}

	return m.index >= len(m.keys)
}

type structIterator struct {
	v     Value
	t     StructType
	index int
}

func (m *structIterator) Next() (key datamodel.Node, value datamodel.Node, err error) {
	i := m.index

	if i >= m.t.NumField() {
		return nil, nil, datamodel.ErrIteratorOverread{}
	}

	m.index++

	f := m.t.FieldByIndex(i)
	v := f.Value(m.v)

	return ValueFrom(reflect.ValueOf(f.Name())).AsNode(), v.AsNode(), nil
}

func (m *structIterator) Done() bool {
	return m.index >= m.t.NumField()
}

type interfaceIterator struct {
	v  Value
	t  Type
	it datamodel.MapIterator
}

func (ii *interfaceIterator) Next() (key datamodel.Node, value datamodel.Node, err error) {
	if ii.t == nil {
		ii.t = TypeFrom(ii.v.Value().Elem().Type())

		if ii.t.PrimitiveKind() == PrimitiveKindStruct {
			ii.it = ValueFrom(ii.v.Value().Elem()).AsNode().MapIterator()
		}

		k := ValueOf("@type").AsNode()
		v := ValueOf(ii.t.Name().NormalizedFullNameWithArguments()).AsNode()

		return k, v, nil
	}

	if ii.it != nil {
		return ii.it.Next()
	}

	return nil, nil, datamodel.ErrIteratorOverread{}
}

func (ii *interfaceIterator) Done() bool {
	if ii.v.Value().IsNil() {
		return true
	}

	if ii.t == nil {
		return false
	}

	if ii.it != nil {
		return ii.it.Done()
	}

	return true
}
