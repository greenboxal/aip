package forddb

import (
	"context"
	"reflect"
)

type Database interface {
	HasListeners

	List(ctx context.Context, typ ResourceTypeID) ([]BasicResource, error)
	Get(ctx context.Context, typ ResourceTypeID, id BasicResourceID) (BasicResource, error)
	Put(ctx context.Context, resource BasicResource) (BasicResource, error)
	Delete(ctx context.Context, resource BasicResource) (BasicResource, error)
}

func Put[T BasicResource](db Database, id T) (def T, _ error) {
	resource, err := db.Put(context.TODO(), id)

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
	resource, err := db.Get(context.TODO(), typ.ID(), id)

	if err != nil {
		return def, err
	}

	if resource == nil {
		return def, ErrNotFound
	}

	return resource.(T), nil
}
