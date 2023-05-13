package forddb

import (
	"reflect"
	"strings"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/schema"
	"github.com/multiformats/go-multicodec"
	"github.com/multiformats/go-multihash"

	"github.com/greenboxal/aip/aip-forddb/pkg/impl/nodebinder"
)

type basicType struct {
	ResourceBase[TypeID, BasicResourceType] `json:"metadata"`

	metadata TypeMetadata

	fields   []BasicField
	fieldMap map[string]BasicField

	typ         reflect.Type
	schemaType  schema.Type
	schemaProto schema.TypedPrototype

	universe *ResourceTypeSystem
}

var _ BasicType = (*basicType)(nil)

func newBasicType(
	typ reflect.Type,
	metadata TypeMetadata,
) *basicType {
	t := &basicType{}

	metadata.ID = TypeID(metadata.Name)
	t.ID = metadata.ID
	t.typ = typ
	t.fieldMap = make(map[string]BasicField, 32)
	t.metadata = metadata

	return t
}

func (bt *basicType) TypeSystem() *ResourceTypeSystem    { return bt.universe }
func (bt *basicType) GetResourceID() TypeID              { return bt.ResourceBase.ID }
func (bt *basicType) Name() string                       { return bt.metadata.Name }
func (bt *basicType) Kind() Kind                         { return bt.metadata.Kind }
func (bt *basicType) PrimitiveKind() PrimitiveKind       { return bt.metadata.PrimitiveKind }
func (bt *basicType) RuntimeType() reflect.Type          { return bt.typ }
func (bt *basicType) Metadata() TypeMetadata             { return bt.metadata }
func (bt *basicType) IsRuntimeOnly() bool                { return bt.metadata.IsRuntimeOnly }
func (bt *basicType) CreateInstance() any                { return reflect.New(bt.typ).Interface() }
func (bt *basicType) NumFields() int                     { return len(bt.fields) }
func (bt *basicType) Fields() []BasicField               { return bt.fields }
func (bt *basicType) FieldByName(name string) BasicField { return bt.fieldMap[name] }

func (bt *basicType) FieldByIndex(index int) BasicField {
	if index < 0 || index >= len(bt.fields) {
		return nil
	}

	return bt.fields[index]
}
func (bt *basicType) SchemaType() schema.Type {
	if bt.metadata.IsRuntimeOnly {
		return nil
	}

	return bt.schemaType
}

func (bt *basicType) SchemaPrototype() schema.TypedPrototype {
	if bt.metadata.IsRuntimeOnly {
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

func (bt *basicType) Initialize(ts *ResourceTypeSystem, options ...nodebinder.Option) {
	var walkFields func(typ reflect.Type, indexBase []int)

	bt.universe = ts

	walkFields = func(typ reflect.Type, indexBase []int) {
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)

			nestedF := f
			nestedF.Index = append([]int{}, indexBase...)
			nestedF.Index = append([]int{}, f.Index...)

			if !f.IsExported() {
				continue
			}

			fieldName := f.Name
			taggedName := ""

			tag, ok := f.Tag.Lookup("json")

			if !ok && !f.Anonymous {
				continue
			}

			if tag == "-" {
				continue
			}

			if tag != "" {
				parts := strings.Split(tag, ",")
				taggedName = parts[0]
				fieldName = taggedName
			}

			actualType := f.Type

			for actualType.Kind() == reflect.Ptr {
				actualType = actualType.Elem()
			}

			fieldType := TypeSystem().LookupByType(f.Type)

			if fieldType == nil {
				TypeSystem().LookupByType(f.Type)
				panic("field type not found")
			}

			if f.Anonymous && fieldType.Kind() != KindId && actualType.Kind() == reflect.Struct && taggedName == "" {
				walkFields(f.Type, nestedF.Index)
			} else {
				field := NewReflectedField(fieldName, bt, fieldType, nestedF)

				bt.fields = append(bt.fields, field)
				bt.fieldMap[field.Name()] = field
			}
		}
	}

	if bt.typ.Kind() == reflect.Struct {
		walkFields(bt.typ, nil)
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
