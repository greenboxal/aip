package forddb

import (
	"reflect"
	"time"
)

type IResourceMetadata interface {
	GetResourceID() BasicResourceID
	GetType() ResourceTypeID
	GetVersion() uint64
}

type BasicResource interface {
	IResourceMetadata

	GetMetadata() *BasicResourceMetadata

	onBeforeSerialize()
}

type Resource[ID BasicResourceID] interface {
	BasicResource

	GetID() ID
}

type BasicResourceMetadata struct {
	Kind      ResourceTypeID `json:"kind"`
	Namespace string         `json:"namespace"`
	Name      string         `json:"name"`
	Version   uint64         `json:"version"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
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

func (r *BasicResourceMetadata) GetVersion() uint64 {
	return r.Version
}

func (r *ResourceMetadata[ID, T]) onBeforeSerialize() {
	r.Kind = r.GetType()
}
