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

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb/nodebinder"
)

type basicTypeImpl struct {
	ResourceBase[TypeID, BasicResourceType] `json:"metadata"`

	name          string
	kind          Kind
	primitiveKind PrimitiveKind

	fields   []BasicField
	fieldMap map[string]BasicField

	typ         reflect.Type
	schemaType  schema.Type
	schemaProto schema.TypedPrototype

	isRuntimeOnly bool

	universe *ResourceTypeSystem
}

var _ BasicType = (*basicTypeImpl)(nil)

func newBasicType(
	kind Kind,
	primitiveKind PrimitiveKind,
	name string,
	typ reflect.Type,
	isRuntimeOnly bool,
) *basicTypeImpl {
	t := &basicTypeImpl{}

	t.ResourceBase.ID = NewStringID[TypeID](name)

	t.name = name
	t.kind = kind
	t.primitiveKind = primitiveKind
	t.typ = typ
	t.fieldMap = make(map[string]BasicField, 32)

	t.isRuntimeOnly = isRuntimeOnly

	return t
}

func (bt *basicTypeImpl) TypeSystem() *ResourceTypeSystem    { return bt.universe }
func (bt *basicTypeImpl) GetResourceID() TypeID              { return bt.ResourceBase.ID }
func (bt *basicTypeImpl) Name() string                       { return bt.name }
func (bt *basicTypeImpl) Kind() Kind                         { return bt.kind }
func (bt *basicTypeImpl) PrimitiveKind() PrimitiveKind       { return bt.primitiveKind }
func (bt *basicTypeImpl) RuntimeType() reflect.Type          { return bt.typ }
func (bt *basicTypeImpl) IsRuntimeOnly() bool                { return bt.isRuntimeOnly }
func (bt *basicTypeImpl) CreateInstance() any                { return reflect.New(bt.typ).Interface() }
func (bt *basicTypeImpl) NumFields() int                     { return len(bt.fields) }
func (bt *basicTypeImpl) Fields() []BasicField               { return bt.fields }
func (bt *basicTypeImpl) FieldByName(name string) BasicField { return bt.fieldMap[name] }

func (bt *basicTypeImpl) FieldByIndex(index int) BasicField {
	if index < 0 || index >= len(bt.fields) {
		return nil
	}

	return bt.fields[index]
}
func (bt *basicTypeImpl) SchemaType() schema.Type {
	if bt.isRuntimeOnly {
		return nil
	}

	return bt.schemaType
}

func (bt *basicTypeImpl) SchemaPrototype() schema.TypedPrototype {
	if bt.isRuntimeOnly {
		return nil
	}

	return bt.schemaProto
}

func (bt *basicTypeImpl) SchemaLinkPrototype() ipld.LinkPrototype {
	return cidlink.LinkPrototype{
		Prefix: cid.Prefix{
			Version:  1,
			Codec:    uint64(multicodec.Raw),
			MhType:   multihash.SHA2_256,
			MhLength: 32,
		},
	}
}

func (bt *basicTypeImpl) Initialize(ts *ResourceTypeSystem, options ...nodebinder.Option) {
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
				field := NewReflectedField(fieldName, bt, fieldType, f)

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
