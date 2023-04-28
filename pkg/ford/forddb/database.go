package forddb

import (
	"reflect"
)

type Database interface {
	HasListeners

	List(typ ResourceTypeID) ([]BasicResource, error)
	Get(typ ResourceTypeID, id BasicResourceID) (BasicResource, error)
	Put(resource BasicResource) (BasicResource, error)
	Delete(resource BasicResource) (BasicResource, error)
}

func Put[T BasicResource](db Database, id T) (def T, _ error) {
	resource, err := db.Put(id)

	if err != nil {
		return def, err
	}

	if resource == nil {
		return def, ErrNotFound
	}

	return resource.(T), nil
}

func Get[T BasicResource](db Database, id ResourceID[T]) (def T, _ error) {
	typ := typeSystem.LookupByIDType(reflect.TypeOf(id))
	resource, err := db.Get(typ.ID(), id)

	if err != nil {
		return def, err
	}

	if resource == nil {
		return def, ErrNotFound
	}

	return resource.(T), nil
}
