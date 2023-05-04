package forddb

import (
	"reflect"
	"sync"

	"github.com/ipld/go-ipld-prime/schema"
)

func newResourceType[ID ResourceID[T], T Resource[ID]](name string) *resourceTypeImpl[ID, T] {
	rt := &resourceTypeImpl[ID, T]{}

	idTemplate := reflect.New(reflect.TypeOf((*ID)(nil)).Elem()).Interface().(BasicResourceID)

	idTyp := DerefType[ID]()
	resourceTyp := DerefType[T]()

	rt.idType = newBasicType(
		KindId,
		idTemplate.PrimitiveKind(),
		name+"ID",
		idTyp,
		false,
	)

	rt.basicTypeImpl = newBasicType(
		KindResource,
		PrimitiveKindStruct,
		name,
		resourceTyp,
		false,
	)

	return rt
}

type resourceTypeImpl[ID ResourceID[T], T Resource[ID]] struct {
	*basicTypeImpl

	m sync.Mutex

	idType       BasicType
	idSchemaType schema.Type
	idPrototype  schema.TypedPrototype
}

func (rt *resourceTypeImpl[ID, T]) CreateID(name string) BasicResourceID {
	idValue := reflect.New(rt.idType.RuntimeType())

	idValue.Interface().(IStringResourceID).SetValueString(name)

	return idValue.Elem().Interface().(BasicResourceID)
}

func (rt *resourceTypeImpl[ID, T]) IDType() BasicType {
	return rt.idType
}

func (rt *resourceTypeImpl[ID, T]) ResourceType() BasicType {
	return rt.basicTypeImpl
}

func (rt *resourceTypeImpl[ID, T]) SetRuntimeOnly() {
	rt.isRuntimeOnly = true
}
