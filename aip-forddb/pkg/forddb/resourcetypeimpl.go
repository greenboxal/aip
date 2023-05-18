package forddb

import (
	"reflect"
	"sync"

	"github.com/ipld/go-ipld-prime/schema"

	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
)

func newResourceType[ID ResourceID[T], T Resource[ID]](name string) *resourceType[ID, T] {
	rt := &resourceType[ID, T]{}

	idTemplate := reflect.New(reflect.TypeOf((*ID)(nil)).Elem()).Interface().(BasicResourceID)

	idTyp := typesystem.TypeFrom(DerefType[ID]())
	resourceTyp := typesystem.TypeFrom(DerefType[T]())

	rt.name = ResourceTypeNameFromSingular(name)

	idTypeMetadata := TypeMetadata{
		Kind:          KindId,
		PrimitiveKind: idTemplate.PrimitiveKind(),
		Name:          name + "ID",
	}

	typeMetadata := TypeMetadata{
		Kind:          KindResource,
		PrimitiveKind: typesystem.PrimitiveKindStruct,
		Name:          name,
	}

	rt.idType = newBasicType(idTyp, idTypeMetadata)
	rt.basicType = newBasicType(resourceTyp, typeMetadata)

	return rt
}

type resourceType[ID ResourceID[T], T Resource[ID]] struct {
	*basicType

	m sync.Mutex

	name             ResourceTypeName
	idType           BasicType
	idSchemaType     schema.Type
	idPrototype      schema.TypedPrototype
	filterableFields []FilterableField
}

func (rt *resourceType[ID, T]) FilterableFields() []FilterableField {
	return rt.filterableFields
}

func (rt *resourceType[ID, T]) ResourceName() ResourceTypeName {
	return rt.name
}

func (rt *resourceType[ID, T]) CreateID(name string) BasicResourceID {
	idValue := reflect.New(rt.idType.RuntimeType())

	idValue.Interface().(IStringResourceID).SetValueString(name)

	return idValue.Elem().Interface().(BasicResourceID)
}

func (rt *resourceType[ID, T]) IDType() BasicType {
	return rt.idType
}

func (rt *resourceType[ID, T]) ResourceType() BasicType {
	return rt.basicType
}

func (rt *resourceType[ID, T]) SetRuntimeOnly() {
	rt.metadata.IsRuntimeOnly = true
}

func (rt *resourceType[ID, T]) Initialize(ts *ResourceTypeSystem) {
	rt.basicType.Initialize(ts)

	if st, ok := rt.Type.(typesystem.StructType); ok {
		for i := 0; i < st.NumField(); i++ {
			var filterable FilterableField

			f := st.FieldByIndex(i)

			filterable.Field = f

			switch f.Type().PrimitiveKind() {
			case typesystem.PrimitiveKindString:
				filterable.Operators = []string{"==", "!="}

			case typesystem.PrimitiveKindBoolean:
				filterable.Operators = []string{"==", "!="}

			case typesystem.PrimitiveKindInt:
				fallthrough
			case typesystem.PrimitiveKindUnsignedInt:
				fallthrough
			case typesystem.PrimitiveKindFloat:
				filterable.Operators = []string{"==", "!=", "<", "<=", ">", ">="}

			default:

				continue
			}

			rt.filterableFields = append(rt.filterableFields, filterable)
		}
	}
}
