package forddb

import (
	"reflect"
	"time"
)

type IResourceMetadata interface {
	GetResourceID() BasicResourceID
	GetType() ResourceTypeID
	GetVersion() int
}

type BasicResource interface {
	IResourceMetadata

	GetMetadata() *BasicResourceMetadata
}

type Resource[ID BasicResourceID] interface {
	BasicResource

	GetID() ID
}

type BasicResourceMetadata struct {
	Namespace string    `json:"namespace"`
	Name      string    `json:"name"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResourceMetadata[ID ResourceID[T], T Resource[ID]] struct {
	BasicResourceMetadata

	ID ID `json:"id"`
}

func (r *ResourceMetadata[ID, T]) GetMetadata() *BasicResourceMetadata {
	return &r.BasicResourceMetadata
}

func (r *ResourceMetadata[ID, T]) GetResourceID() BasicResourceID {
	return r.ID
}

func (r *ResourceMetadata[ID, T]) GetID() ID {
	return r.ID
}

func (r *ResourceMetadata[ID, T]) GetType() ResourceTypeID {
	t := typeSystem.LookupByResourceType(reflect.TypeOf((*T)(nil)).Elem())

	if t == nil {
		panic("resource type not found")
	}

	return t.GetID()
}

func (r *BasicResourceMetadata) GetVersion() int {
	return r.Version
}
