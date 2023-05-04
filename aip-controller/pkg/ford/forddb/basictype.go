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
	ResourceMetadata[ResourceTypeID, BasicResourceType] `json:"metadata"`

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

	t.ResourceMetadata.ID = NewStringID[ResourceTypeID](name)
	t.ResourceMetadata.Name = name

	t.kind = kind
	t.primitiveKind = primitiveKind
	t.typ = typ
	t.fieldMap = make(map[string]BasicField, 32)

	t.isRuntimeOnly = isRuntimeOnly

	return t
}

func (bt *basicTypeImpl) TypeSystem() *ResourceTypeSystem    { return bt.universe }
func (bt *basicTypeImpl) GetID() ResourceTypeID              { return bt.ResourceMetadata.ID }
func (bt *basicTypeImpl) Name() string                       { return bt.ResourceMetadata.Name }
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
	bt.universe = ts

	if bt.typ.Kind() == reflect.Struct {
		for i := 0; i < bt.typ.NumField(); i++ {
			f := bt.typ.Field(i)

			if !f.IsExported() {
				continue
			}

			fieldName := f.Name
			fieldType := TypeSystem().LookupByType(f.Type)

			if fieldType == nil {
				fieldType = TypeSystem().LookupByType(f.Type)
				panic("field type not found")
			}

			tag, ok := f.Tag.Lookup("json")

			if !ok {
				continue
			}

			tagParts := strings.Split(tag, ",")
			fieldName = tagParts[0]

			if f.Anonymous {
				fieldStructType := fieldType.(*basicTypeImpl)

				for _, field := range fieldStructType.fields {
					bt.fields = append(bt.fields, field)
					bt.fieldMap[field.Name()] = field
				}
			} else {
				field := NewReflectedField(fieldName, bt, fieldType, f)

				bt.fields = append(bt.fields, field)
				bt.fieldMap[field.Name()] = field
			}
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
