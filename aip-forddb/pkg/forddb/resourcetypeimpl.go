package forddb

import (
	"reflect"
	"sync"

	"github.com/ipld/go-ipld-prime/schema"

	"github.com/greenboxal/aip/aip-forddb/pkg/impl/nodebinder"
)

func newResourceType[ID ResourceID[T], T Resource[ID]](name string) *resourceType[ID, T] {
	rt := &resourceType[ID, T]{}

	idTemplate := reflect.New(reflect.TypeOf((*ID)(nil)).Elem()).Interface().(BasicResourceID)

	idTyp := DerefType[ID]()
	resourceTyp := DerefType[T]()

	rt.name = ResourceTypeNameFromSingular(name)

	idTypeMetadata := TypeMetadata{
		Kind:          KindId,
		PrimitiveKind: idTemplate.PrimitiveKind(),
		Name:          name + "ID",
	}

	typeMetadata := TypeMetadata{
		Kind:          KindResource,
		PrimitiveKind: PrimitiveKindStruct,
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

func (rt *resourceType[ID, T]) Initialize(ts *ResourceTypeSystem, options ...nodebinder.Option) {
	rt.basicType.Initialize(ts, options...)

	for _, f := range rt.fields {
		var filterable FilterableField

		filterable.Field = f

		switch f.BasicType().PrimitiveKind() {
		case PrimitiveKindString:
			filterable.Operators = []string{"==", "!="}

		case PrimitiveKindBoolean:
			filterable.Operators = []string{"==", "!="}

		case PrimitiveKindInt:
			fallthrough
		case PrimitiveKindUnsignedInt:
			fallthrough
		case PrimitiveKindFloat:
			filterable.Operators = []string{"==", "!=", "<", "<=", ">", ">="}

		default:
			if f.BasicType().Kind() == KindId {
				filterable.Operators = []string{"==", "!="}
			} else {
				continue
			}

			continue
		}

		rt.filterableFields = append(rt.filterableFields, filterable)
	}
}
