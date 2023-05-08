package forddb

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/multiformats/go-multibase"
	"github.com/multiformats/go-multihash"
)

type ResourceMetadata interface {
	GetResourceMetadata() *Metadata
	GetResourceBasicID() BasicResourceID
	GetResourceTypeID() TypeID
	GetResourceVersion() uint64
}

type Metadata struct {
	Kind      TypeID        `json:"kind"`
	Scope     ResourceScope `json:"scope"`
	Namespace string        `json:"namespace"`
	Version   uint64        `json:"version"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (r *Metadata) GetResourceMetadata() *Metadata { return r }
func (r *Metadata) GetResourceVersion() uint64     { return r.Version }
func (r *Metadata) GetResourceTypeID() TypeID {
	return r.Kind
}

type ResourceBase[ID ResourceID[T], T Resource[ID]] struct {
	ID ID `json:"id"`
	Metadata
}

func (r *ResourceBase[ID, T]) GetResourceID() ID                   { return r.ID }
func (r *ResourceBase[ID, T]) GetResourceBasicID() BasicResourceID { return r.ID }

func (r *ResourceBase[ID, T]) GetResourceTypeID() TypeID {
	if r.Kind == "" {
		r.Kind = r.ID.BasicResourceType().GetResourceID()
	}

	return r.Kind
}

func (r *ResourceBase[ID, T]) GetResourceType() ResourceType[ID, T] {
	return TypeSystem().LookupByResourceType(reflect.TypeOf((*T)(nil)).Elem()).(ResourceType[ID, T])
}

func (r *ResourceBase[ID, T]) OnBeforeSave(self BasicResource) {
	r.Kind = self.GetResourceTypeID()
}

type BasicContentAddressedResource interface {
	GetContentAddressableRoot() any
}

type ContentAddressedResource[ID BasicResourceID] interface {
	Resource[ID]

	BasicContentAddressedResource
}

type ContentAddressedResourceBase[ID ResourceID[T], T ContentAddressedResource[ID]] struct {
	ResourceBase[ID, T]
}

func (r *ContentAddressedResourceBase[ID, T]) OnBeforeSave(self BasicResource) {
	car := self.(BasicContentAddressedResource)

	r.ID = CreateContentAddressableID[ID](car.GetContentAddressableRoot())

	r.ResourceBase.OnBeforeSave(self)
}

func CreateContentAddressableID[ID BasicResourceID](spec any) ID {
	data, err := json.Marshal(spec)

	if err != nil {
		panic(err)
	}

	h, err := multihash.Sum(data, multihash.SHA1, -1)

	if err != nil {
		panic(err)
	}

	b, err := multibase.Encode(multibase.Base36, h)

	if err != nil {
		panic(err)
	}

	return NewStringID[ID](b)
}
