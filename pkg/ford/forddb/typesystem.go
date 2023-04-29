package forddb

import (
	"reflect"
	"sync"

	"github.com/ipld/go-ipld-prime/schema"
)

var basicResourceType = reflect.TypeOf((*BasicResourceID)(nil)).Elem()
var basicResourceIdType = reflect.TypeOf((*BasicResourceID)(nil)).Elem()
var basicResourcePointerType = reflect.TypeOf((*BasicResourcePointer)(nil)).Elem()
var basicResourceSlotType = reflect.TypeOf((*BasicResourceSlot)(nil)).Elem()

var typeSystem = &ResourceTypeSystem{
	resourceTypes:   make(map[ResourceTypeID]BasicResourceType, 32),
	idTypeMap:       make(map[reflect.Type]BasicResourceType, 32),
	resourceTypeMap: make(map[reflect.Type]BasicResourceType, 32),
	typeSchemaCache: make(map[reflect.Type]schema.Type, 32),
}

var typeType = &resourceType[ResourceTypeID, BasicResourceType]{
	resourceType: reflect.TypeOf((*BasicResourceType)(nil)).Elem(),
	idType:       reflect.TypeOf((*ResourceTypeID)(nil)).Elem(),
}

func init() {
	typeType.isRuntimeOnly = true
	typeType.ResourceMetadata.Name = "type"
	typeType.ResourceMetadata.ID = typeType.MakeId(typeType.ResourceMetadata.Name).(ResourceTypeID)

	typeSystem.universe.Init()

	typeSystem.Register(typeType)
}

type ResourceTypeSystem struct {
	m sync.Mutex

	resourceTypes   map[ResourceTypeID]BasicResourceType
	idTypeMap       map[reflect.Type]BasicResourceType
	resourceTypeMap map[reflect.Type]BasicResourceType
	typeSchemaCache map[reflect.Type]schema.Type

	universe schema.TypeSystem
}

func (rts *ResourceTypeSystem) Register(t BasicResourceType) {
	rts.resourceTypes[t.ID()] = t
	rts.idTypeMap[t.IDType()] = t
	rts.resourceTypeMap[t.ResourceType()] = t

	t.initializeSchema(rts)
}

func (rts *ResourceTypeSystem) LookupByID(id ResourceTypeID) BasicResourceType {
	return rts.resourceTypes[id]
}

func (rts *ResourceTypeSystem) LookupByIDType(typ reflect.Type) BasicResourceType {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	return rts.idTypeMap[typ]
}

func (rts *ResourceTypeSystem) LookupByResourceType(typ reflect.Type) BasicResourceType {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	return rts.resourceTypeMap[typ]
}

func (rts *ResourceTypeSystem) ResourceTypes() []BasicResourceType {
	var result = make([]BasicResourceType, 0, len(rts.resourceTypes))

	for _, typ := range rts.resourceTypes {
		result = append(result, typ)
	}

	return result
}

func (rts *ResourceTypeSystem) Freeze() (*schema.TypeSystem, []error) {
	var result = make([]schema.Type, 0, len(rts.resourceTypes))

	for _, typ := range rts.resourceTypes {
		result = append(result, typ.SchemaResourceType())
	}

	return schema.SpawnTypeSystem(result...)
}

func (rts *ResourceTypeSystem) SchemaForType(typ reflect.Type) schema.Type {
	var result schema.Type

	typ = derefPointer(typ)

	if existing := rts.typeSchemaCache[typ]; existing != nil {
		return existing
	}

	switch typ.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		elem := rts.SchemaForType(typ.Elem())

		result = schema.SpawnList(elem.Name()+"List", elem.Name(), true)

	case reflect.Struct:
		result = rts.schemaForStruct(typ)

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
		name := normalizedTypeName(typ)

		if name == "" {
			name = typ.Kind().String()
		}

		result = schema.SpawnInt(name)

	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		name := normalizedTypeName(typ)

		if name == "" {
			name = typ.Kind().String()
		}

		result = schema.SpawnFloat(name)

	case reflect.Bool:
		name := normalizedTypeName(typ)

		if name == "" {
			name = typ.Kind().String()
		}

		result = schema.SpawnBool(name)

	case reflect.String:
		name := normalizedTypeName(typ)

		if name == "" {
			name = typ.Kind().String()
		}

		result = schema.SpawnString(name)

	default:
		panic("unsupported type")
	}

	rts.accumulate(typ, result)

	return result
}

func (rts *ResourceTypeSystem) schemaForStruct(typ reflect.Type) schema.Type {
	var fields []schema.StructField
	var repr schema.StructRepresentation
	name := normalizedTypeName(typ)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if !field.IsExported() {
			continue
		}

		fieldType := derefPointer(field.Type)

		if IsBasicResourcePointer(fieldType) {
			fieldSchemaType := rts.SchemaForType(field.Type)

			ref := schema.SpawnLinkReference("_"+fieldSchemaType.Name()+"Ptr", fieldSchemaType.Name())

			rts.accumulate(field.Type, ref)

			f := schema.SpawnStructField(
				field.Name,
				ref.Name(),
				false,
				field.Type.Kind() == reflect.Ptr || field.Type.Kind() == reflect.Interface,
			)

			fields = append(fields, f)
		} else if IsBasicResource(fieldType) {
			fieldSchemaType := rts.SchemaForType(field.Type)

			f := schema.SpawnStructField(
				field.Name,
				fieldSchemaType.Name(),
				false,
				field.Type.Kind() == reflect.Ptr || field.Type.Kind() == reflect.Interface,
			)

			fields = append(fields, f)
		} else {
			fieldTypeName := normalizedTypeName(field.Type)

			f := schema.SpawnStructField(
				field.Name,
				fieldTypeName,
				false,
				field.Type.Kind() == reflect.Ptr || field.Type.Kind() == reflect.Interface,
			)

			fields = append(fields, f)
		}
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
