package typesystem

import (
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/ipld/go-ipld-prime/schema"

	"github.com/greenboxal/aip/pkg/ford/forddb"
	"github.com/greenboxal/aip/pkg/ford/forddb/nodebinder"
	"github.com/greenboxal/aip/pkg/utils"
)

var typeSystem = NewResourceTypeSystem()

func TypeSystem() *ResourceTypeSystem {
	return typeSystem
}

type ResourceTypeSystem struct {
	m sync.Mutex

	resourceTypes   map[forddb.ResourceTypeID]forddb.BasicResourceType
	idTypeMap       map[reflect.Type]forddb.BasicResourceType
	resourceTypeMap map[reflect.Type]forddb.BasicResourceType
	typeMap         map[reflect.Type]forddb.BasicType
	typeSchemaCache map[reflect.Type]schema.Type

	universe        schema.TypeSystem
	bindNodeOptions []nodebinder.Option
}

func NewResourceTypeSystem() *ResourceTypeSystem {
	typeType := NewResourceType[forddb.ResourceTypeID, forddb.BasicResourceType]("type")
	typeType.SetRuntimeOnly()

	rts := &ResourceTypeSystem{
		resourceTypes:   make(map[forddb.ResourceTypeID]forddb.BasicResourceType, 32),
		idTypeMap:       make(map[reflect.Type]forddb.BasicResourceType, 32),
		typeMap:         make(map[reflect.Type]forddb.BasicType, 32),
		resourceTypeMap: make(map[reflect.Type]forddb.BasicResourceType, 32),
		typeSchemaCache: make(map[reflect.Type]schema.Type, 32),
	}

	rts.universe.Init()

	rts.Register(typeType)

	rts.bindNodeOptions = append(rts.bindNodeOptions, nodebinder.TypedIntConverter(
		(*time.Time)(nil),
		func(i int64) (interface{}, error) {
			return time.Unix(0, i), nil
		},
		func(i interface{}) (int64, error) {
			t := i.(time.Time)

			return t.UnixNano(), nil
		},
	))

	rts.accumulate(reflect.TypeOf((*time.Time)(nil)).Elem(), schema.SpawnInt("timeTime"))

	return rts
}

func (rts *ResourceTypeSystem) Register(t forddb.BasicResourceType) {
	if rts.doRegister(t, true) {
		t.Initialize(rts, rts.bindNodeOptions...)
	}
}
func (rts *ResourceTypeSystem) doRegister(t forddb.BasicType, lock bool) bool {
	if lock {
		rts.m.Lock()
		defer rts.m.Unlock()
	}

	if _, ok := rts.typeMap[t.RuntimeType()]; ok {
		return false
	}

	rts.typeMap[t.RuntimeType()] = t

	if t, ok := t.(forddb.BasicResourceType); ok {
		rts.resourceTypes[t.GetID()] = t
		rts.resourceTypeMap[t.RuntimeType()] = t
		rts.idTypeMap[t.IDType().RuntimeType()] = t
	}

	if t.Kind() == forddb.KindId {
		rts.bindNodeOptions = append(
			rts.bindNodeOptions,

			nodebinder.TypedStringConverter(
				reflect.New(t.RuntimeType()),
				func(s string) (interface{}, error) {
					idVal := reflect.New(t.RuntimeType())
					idStr := idVal.Interface().(forddb.IStringResourceID)
					idStr.SetValueString(s)
					return idVal.Elem().Interface(), nil
				},
				func(i interface{}) (string, error) {
					return i.(forddb.BasicResourceID).String(), nil
				},
			),
		)
	}

	return true
}

func (rts *ResourceTypeSystem) LookupByID(id forddb.ResourceTypeID) forddb.BasicResourceType {
	return rts.resourceTypes[id]
}

func (rts *ResourceTypeSystem) LookupByIDType(typ reflect.Type) forddb.BasicResourceType {
	typ = DerefPointer(typ)

	return rts.idTypeMap[typ]
}

func (rts *ResourceTypeSystem) LookupByResourceType(typ reflect.Type) forddb.BasicResourceType {
	typ = DerefPointer(typ)

	return rts.resourceTypeMap[typ]
}

func (rts *ResourceTypeSystem) LookupByType(typ reflect.Type) (result forddb.BasicType) {
	isNew := false

	if IsBasicResource(typ) {
		return rts.LookupByResourceType(typ)
	}

	defer func() {
		if isNew {
			result.Initialize(rts, rts.bindNodeOptions...)
		}
	}()

	typ = DerefPointer(typ)

	if existing := rts.typeMap[typ]; existing != nil {
		return existing
	}

	rts.m.Lock()
	defer rts.m.Unlock()

	if existing := rts.typeMap[typ]; existing != nil {
		return existing
	}

	name := utils.GetParsedTypeName(typ).NormalizedFullNameWithArguments()

	kind := forddb.KindValue

	if IsBasicResourceId(typ) {
		kind = forddb.KindId
	}

	result = NewBasicType(
		kind,
		getTypePrimitiveKind(typ),
		name,
		typ,
		false,
	)

	isNew = true

	rts.doRegister(result, false)

	return
}

func (rts *ResourceTypeSystem) ResourceTypes() []forddb.BasicResourceType {
	var result = make([]forddb.BasicResourceType, 0, len(rts.resourceTypes))

	for _, typ := range rts.resourceTypes {
		result = append(result, typ)
	}

	return result
}

func (rts *ResourceTypeSystem) Freeze() (*schema.TypeSystem, []error) {
	var result = make([]schema.Type, 0, len(rts.resourceTypes))

	for _, typ := range rts.typeMap {
		result = append(result, typ.SchemaType())
	}

	return schema.SpawnTypeSystem(result...)
}

func (rts *ResourceTypeSystem) SchemaForType(typ reflect.Type) schema.Type {
	var result schema.Type

	typ = DerefPointer(typ)

	if existing := rts.typeSchemaCache[typ]; existing != nil {
		return existing
	}

	name := utils.GetParsedTypeName(typ).NormalizedFullNameWithArguments()

	if name == "" {
		name = typ.Kind().String()
	}

	switch typ.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		elem := rts.SchemaForType(typ.Elem())

		result = schema.SpawnList(elem.Name()+"List", elem.Name(), false)

	case reflect.Struct:
		if IsBasicResourceId(typ) {
			result = rts.schemaForId(typ)
		} else {
			result = rts.schemaForStruct(typ)
		}

	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		result = schema.SpawnInt(name)

	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		result = schema.SpawnFloat(name)

	case reflect.Bool:
		result = schema.SpawnBool(name)

	case reflect.String:
		result = schema.SpawnString(name)

	default:
		panic("unsupported type")
	}

	rts.accumulate(typ, result)

	return result
}

func (rts *ResourceTypeSystem) schemaForStruct(typ reflect.Type) schema.Type {
	var fields []schema.StructField

	name := utils.GetParsedTypeName(typ).NormalizedFullNameWithArguments()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if !field.IsExported() {
			continue
		}

		fieldSchemaType := rts.SchemaForType(field.Type)

		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			if fieldSchemaType, ok := fieldSchemaType.(*schema.TypeStruct); ok {
				fields = append(fields, fieldSchemaType.Fields()...)
				continue
			}
		}

		fieldName := field.Name

		if tag, ok := field.Tag.Lookup("json"); ok {
			parts := strings.SplitN(tag, ",", 2)
			fieldName = parts[0]
		}

		f := schema.SpawnStructField(
			fieldName,
			fieldSchemaType.Name(),
			false,
			field.Type.Kind() == reflect.Ptr || field.Type.Kind() == reflect.Interface,
		)

		fields = append(fields, f)
	}

	var repr schema.StructRepresentation

	if IsBasicResourceId(typ) {
		repr = schema.SpawnStructRepresentationStringjoin("")
	}

	return schema.SpawnStruct(name, fields, repr)
}

func (rts *ResourceTypeSystem) accumulate(typ reflect.Type, ref schema.Type) {
	rts.m.Lock()
	defer rts.m.Unlock()

	rts.typeSchemaCache[typ] = ref

	if existing := rts.universe.TypeByName(ref.Name()); existing != nil {
		if existing != ref {
			panic("duplicate type name")
		}

		return
	}

	rts.universe.Accumulate(ref)
}

func (rts *ResourceTypeSystem) MakePrototype(typ reflect.Type, schemaType schema.Type) schema.TypedPrototype {
	return nodebinder.Prototype(reflect.New(typ).Interface(), schemaType, rts.bindNodeOptions...)
}

func (rts *ResourceTypeSystem) schemaForId(typ reflect.Type) schema.Type {
	return rts.schemaForStruct(typ)
}

func getTypePrimitiveKind(typ reflect.Type) forddb.PrimitiveKind {
	switch typ.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		elem := typ.Elem()

		if elem.Kind() == reflect.Uint8 && elem.Name() == "" {
			return forddb.PrimitiveKindBytes
		} else {
			return forddb.PrimitiveKindList
		}

	case reflect.Struct:
		return forddb.PrimitiveKindStruct

	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		return forddb.PrimitiveKindInt
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		return forddb.PrimitiveKindUnsignedInt

	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return forddb.PrimitiveKindFloat

	case reflect.Bool:
		return forddb.PrimitiveKindBoolean

	case reflect.Map:
		return forddb.PrimitiveKindMap

	case reflect.String:
		return forddb.PrimitiveKindString

	default:
		panic("unsupported type")
	}
}
