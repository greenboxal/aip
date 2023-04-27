package forddb

import (
	"errors"
	"reflect"
)

type Database interface {
	HasListeners

	List(typ ResourceTypeID) ([]BasicResource, error)
	Get(typ ResourceTypeID, id BasicResourceID) (BasicResource, error)
	UpdateOrCreate(resource BasicResource) (BasicResource, error)
	Delete(resource BasicResource) (BasicResource, error)
}

var ErrVersionMismatch = errors.New("version mismatch")

func CreateOrUpdate[T BasicResource](db Database, id T) (def T, _ error) {
	resource, err := db.UpdateOrCreate(id)

	if err != nil {
		return def, err
	}

	if resource == nil {
		return def, nil
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
		return def, nil
	}

	return resource.(T), nil
}
