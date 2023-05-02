package forddb

import (
	"reflect"
	"sync"

	"github.com/ipld/go-ipld-prime/schema"

	"github.com/greenboxal/aip/pkg/ford/forddb/nodebinder"
)

type BasicResourceType interface {
	BasicType

	IDType() BasicType
	ResourceType() BasicType

	CreateID(name string) BasicResourceID
}

type ResourceType[ID ResourceID[T], T Resource[ID]] interface {
	Type[T]

	BasicResourceType
}

func DefineResourceType[ID ResourceID[T], T Resource[ID]](name string) ResourceType[ID, T] {
	t := newResourceType[ID, T](name)

	typeSystem.Register(t)

	return t
}

func newResourceType[ID ResourceID[T], T Resource[ID]](name string) *resourceType[ID, T] {
	rt := &resourceType[ID, T]{}

	rt.idType = newBasicType(KindId, name+"ID", derefType[ID](), false)
	rt.basicType = newBasicType(KindResource, name, derefType[T](), false)

	return rt
}

type resourceType[ID ResourceID[T], T Resource[ID]] struct {
	*basicType

	m sync.Mutex

	idType       BasicType
	idSchemaType schema.Type
	idPrototype  schema.TypedPrototype
}

func (rt *resourceType[ID, T]) CreateID(name string) BasicResourceID {
	idValue := reflect.New(rt.idType.RuntimeType())

	idValue.Interface().(IStringResourceID).setValueString(name)

	return idValue.Elem().Interface().(BasicResourceID)
}

func (rt *resourceType[ID, T]) IDType() BasicType {
	return rt.idType
}

func (rt *resourceType[ID, T]) ResourceType() BasicType {
	return rt.basicType
}

func (rt *resourceType[ID, T]) initializeSchema(ts *ResourceTypeSystem, options ...nodebinder.Option) {

	rt.basicType.initializeSchema(ts, options...)
}
