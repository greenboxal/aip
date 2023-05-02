package forddb

import (
	"reflect"
	"strings"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/schema"
	"github.com/multiformats/go-multicodec"
	"github.com/multiformats/go-multihash"

	"github.com/greenboxal/aip/pkg/ford/forddb/nodebinder"
)

type Kind int

const (
	KindInvalid Kind = iota
	KindId
	KindResource
	KindValue
	KindPointer
)

type BasicType interface {
	BasicResource

	GetID() ResourceTypeID
	Name() string
	Kind() Kind

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

	initializeSchema(ts *ResourceTypeSystem, options ...nodebinder.Option)
}

type Type[T any] interface {
	BasicType
}

type basicType struct {
	ResourceMetadata[ResourceTypeID, BasicResourceType] `json:"metadata"`

	kind     Kind
	fields   []BasicField
	fieldMap map[string]BasicField

	typ         reflect.Type
	schemaType  schema.Type
	schemaProto schema.TypedPrototype

	isRuntimeOnly bool

	universe *ResourceTypeSystem
}

var _ BasicType = (*basicType)(nil)

func newBasicType(kind Kind, name string, typ reflect.Type, isRuntimeOnly bool) *basicType {
	t := &basicType{}

	t.ResourceMetadata.ID = NewStringID[ResourceTypeID](name)
	t.ResourceMetadata.Name = name

	t.kind = kind
	t.typ = typ
	t.fieldMap = make(map[string]BasicField, 32)

	t.isRuntimeOnly = isRuntimeOnly

	return t
}

func (bt *basicType) GetID() ResourceTypeID {
	return bt.ResourceMetadata.ID
}

func (bt *basicType) Name() string {
	return bt.ResourceMetadata.Name
}

func (bt *basicType) Kind() Kind {
	return bt.kind
}

func (bt *basicType) RuntimeType() reflect.Type {
	return bt.typ
}

func (bt *basicType) IsRuntimeOnly() bool {
	return bt.isRuntimeOnly
}

func (bt *basicType) NumFields() int {
	return len(bt.fields)
}

func (bt *basicType) Fields() []BasicField {
	return bt.fields
}

func (bt *basicType) FieldByName(name string) BasicField {
	return bt.fieldMap[name]
}

func (bt *basicType) FieldByIndex(index int) BasicField {
	if index < 0 || index >= len(bt.fields) {
		return nil
	}

	return bt.fields[index]
}

func (bt *basicType) CreateInstance() any {
	return reflect.New(bt.typ).Interface()
}

func (bt *basicType) SchemaType() schema.Type {
	if bt.isRuntimeOnly {
		return nil
	}

	return bt.schemaType
}

func (bt *basicType) SchemaPrototype() schema.TypedPrototype {
	if bt.isRuntimeOnly {
		return nil
	}

	return bt.schemaProto
}

func (bt *basicType) SchemaLinkPrototype() ipld.LinkPrototype {
	return cidlink.LinkPrototype{
		Prefix: cid.Prefix{
			Version:  1,
			Codec:    uint64(multicodec.Raw),
			MhType:   multihash.SHA2_256,
			MhLength: 32,
		},
	}
}

func (bt *basicType) TypeSystem() *ResourceTypeSystem {
	return bt.universe
}

//func (bt *basicType) Encode(resource any) (RawResource, error) {
//	node := nodebinder.Wrap(resource, bt.SchemaType(), bt.universe.bindNodeOptions...)
//
//	return node, nil
//}
//
//func (bt *basicType) Decode(resource RawResource) (any, error) {
//	res := nodebinder.Unwrap(resource)
//
//	if res == nil {
//		encoded, err := ipld.Encode(resource, dagjson.Encode)
//
//		if err != nil {
//			return nil, err
//		}
//
//		res, err = ipld.DecodeUsingPrototype(encoded, dagjson.Decode, bt.SchemaPrototype())
//
//		if err != nil {
//			return nil, err
//		}
//
//		res = nodebinder.Unwrap(resource)
//	}
//
//	return res, nil
//}

func (bt *basicType) initializeSchema(ts *ResourceTypeSystem, options ...nodebinder.Option) {
	bt.universe = ts

	if bt.typ.Kind() == reflect.Struct {
		for i := 0; i < bt.typ.NumField(); i++ {
			f := bt.typ.Field(i)

			if !f.IsExported() {
				continue
			}

			fieldName := f.Name
			fieldType := typeSystem.LookupByType(f.Type)

			tag, ok := f.Tag.Lookup("json")

			if !ok {
				continue
			}

			tagParts := strings.Split(tag, ",")
			fieldName = tagParts[0]

			field := NewReflectedField(fieldName, bt, fieldType, f)

			bt.fields = append(bt.fields, field)
			bt.fieldMap[field.Name()] = field
		}
	}

	//if !bt.isRuntimeOnly {
	//	bt.schemaType = bt.universe.SchemaForType(bt.typ)

	//	if bt.Kind() == KindId {
	//		bt.schemaProto = &resourceIdPrototype{typ: bt}
	//	} else {
	//		bt.schemaProto = bt.universe.MakePrototype(bt.typ, bt.schemaType)
	//	}
	//}
}

type resourceIdPrototype struct {
	typ BasicType
}

func (r *resourceIdPrototype) NewBuilder() datamodel.NodeBuilder {
	return &resourceIdBuilder{proto: r}
}

func (r *resourceIdPrototype) Type() schema.Type {
	return r.typ.SchemaType()
}

func (r *resourceIdPrototype) Representation() datamodel.NodePrototype {
	return r
}

type resourceIdBuilder struct {
	proto *resourceIdPrototype
	value reflect.Value
}

func (r *resourceIdBuilder) BeginMap(sizeHint int64) (datamodel.MapAssembler, error) {
	panic("not supported")
}

func (r *resourceIdBuilder) BeginList(sizeHint int64) (datamodel.ListAssembler, error) {
	panic("not supported")
}

func (r *resourceIdBuilder) AssignNull() error {
	panic("not supported")
}

func (r *resourceIdBuilder) AssignBool(b bool) error {
	panic("not supported")
}

func (r *resourceIdBuilder) AssignInt(i int64) error {
	panic("not supported")
}

func (r *resourceIdBuilder) AssignFloat(f float64) error {
	panic("not supported")
}

func (r *resourceIdBuilder) AssignString(s string) error {
	id := r.value.Interface().(IStringResourceID)

	id.setValueString(s)

	return nil
}

func (r *resourceIdBuilder) AssignBytes(bytes []byte) error {
	panic("not supported")
}

func (r *resourceIdBuilder) AssignLink(link datamodel.Link) error {
	panic("not supported")
}

func (r *resourceIdBuilder) AssignNode(node datamodel.Node) error {
	panic("not supported")
}

func (r *resourceIdBuilder) Prototype() datamodel.NodePrototype {
	return r.proto
}

func (r *resourceIdBuilder) Build() datamodel.Node {
	return &resourceIdNode{
		typ:   r.proto.typ,
		value: r.value,
	}
}

func (r *resourceIdBuilder) Reset() {
}

type resourceIdNode struct {
	typ   BasicType
	value reflect.Value
}

func (r *resourceIdNode) Kind() datamodel.Kind {
	return datamodel.Kind_String
}

func (r *resourceIdNode) LookupByString(key string) (datamodel.Node, error) {
	//TODO implement me
	panic("implement me")
}

func (r *resourceIdNode) LookupByNode(key datamodel.Node) (datamodel.Node, error) {
	//TODO implement me
	panic("implement me")
}

func (r *resourceIdNode) LookupByIndex(idx int64) (datamodel.Node, error) {
	//TODO implement me
	panic("implement me")
}

func (r *resourceIdNode) LookupBySegment(seg datamodel.PathSegment) (datamodel.Node, error) {
	//TODO implement me
	panic("implement me")
}

func (r *resourceIdNode) MapIterator() datamodel.MapIterator {
	//TODO implement me
	panic("implement me")
}

func (r *resourceIdNode) ListIterator() datamodel.ListIterator {
	//TODO implement me
	panic("implement me")
}

func (r *resourceIdNode) Length() int64 {
	//TODO implement me
	panic("implement me")
}

func (r *resourceIdNode) IsAbsent() bool {
	//TODO implement me
	panic("implement me")
}

func (r *resourceIdNode) IsNull() bool {
	//TODO implement me
	panic("implement me")
}

func (r *resourceIdNode) AsBool() (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *resourceIdNode) AsInt() (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (r *resourceIdNode) AsFloat() (float64, error) {
	//TODO implement me
	panic("implement me")
}

func (r *resourceIdNode) AsString() (string, error) {
	id := r.value.Interface().(IStringResourceID)

	return id.String(), nil
}

func (r *resourceIdNode) AsBytes() ([]byte, error) {
	id := r.value.Interface().(IStringResourceID)

	return id.MarshalBinary()
}

func (r *resourceIdNode) AsLink() (datamodel.Link, error) {
	//TODO implement me
	panic("implement me")
}

func (r *resourceIdNode) Prototype() datamodel.NodePrototype {
	return r.typ.SchemaPrototype()
}

func (r *resourceIdNode) Type() schema.Type {
	return r.typ.SchemaType()
}

func (r *resourceIdNode) Representation() datamodel.Node {
	return r
}
