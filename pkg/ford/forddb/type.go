package forddb

import (
	"reflect"
	"sync"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
	"github.com/multiformats/go-multicodec"
	"github.com/multiformats/go-multihash"
)

type BasicResourceType interface {
	BasicResource

	GetID() ResourceTypeID
	Name() string
	IDType() reflect.Type
	ResourceType() reflect.Type
	IsRuntimeOnly() bool

	CreateInstance() BasicResource
	CreateID(name string) BasicResourceID

	SchemaIdType() schema.Type
	SchemaIdPrototype() schema.TypedPrototype

	SchemaResourceType() schema.Type
	SchemaResourcePrototype() schema.TypedPrototype

	SchemaLinkPrototype() ipld.LinkPrototype

	TypeSystem() *ResourceTypeSystem

	initializeSchema(ts *ResourceTypeSystem, options ...bindnode.Option)
}

type ResourceType[ID ResourceID[T], T Resource[ID]] interface {
	BasicResourceType
}

func LookupTypeByName(name string) BasicResourceType {
	return typeSystem.LookupByID(NewStringID[ResourceTypeID](name))
}

func DefineResourceType[ID ResourceID[T], T Resource[ID]](name string) ResourceType[ID, T] {
	t := &resourceType[ID, T]{
		idType:       derefType[ID](),
		resourceType: derefType[T](),
	}

	t.ResourceMetadata.ID = NewStringID[ResourceTypeID](name)
	t.ResourceMetadata.Name = name

	typeSystem.Register(t)

	return t
}

type resourceType[ID ResourceID[T], T Resource[ID]] struct {
	ResourceMetadata[ResourceTypeID, BasicResourceType]

	idType       reflect.Type
	resourceType reflect.Type

	isRuntimeOnly bool

	idSchemaType schema.Type
	idPrototype  schema.TypedPrototype

	resourceSchemaType schema.Type
	resourcePrototype  schema.TypedPrototype

	m        sync.Mutex
	universe *ResourceTypeSystem
}

func (r *resourceType[ID, T]) SchemaLinkPrototype() ipld.LinkPrototype {
	return cidlink.LinkPrototype{
		Prefix: cid.Prefix{
			Version:  1,
			Codec:    uint64(multicodec.Raw),
			MhType:   multihash.SHA2_256,
			MhLength: 32,
		},
	}
}

func (r *resourceType[ID, T]) TypeSystem() *ResourceTypeSystem {
	return r.universe
}

func (r *resourceType[ID, T]) SchemaResourceType() schema.Type {
	if r.isRuntimeOnly {
		return nil
	}

	if r.resourceSchemaType == nil {
		r.m.Lock()
		defer r.m.Unlock()

		if r.resourceSchemaType == nil {
			r.resourceSchemaType = r.universe.SchemaForType(r.resourceType)
		}
	}

	return r.resourceSchemaType
}

func (r *resourceType[ID, T]) SchemaResourcePrototype() schema.TypedPrototype {
	if r.isRuntimeOnly {
		return nil
	}

	if r.resourcePrototype == nil {
		r.m.Lock()
		defer r.m.Unlock()

		if r.resourcePrototype == nil {
			r.resourcePrototype = r.universe.MakePrototype(r.resourceType, r.resourceSchemaType)
		}
	}

	return r.resourcePrototype
}

func (r *resourceType[ID, T]) SchemaIdType() schema.Type {
	if r.resourceSchemaType == nil {
		r.m.Lock()
		defer r.m.Unlock()

		if r.idSchemaType == nil {
			r.idSchemaType = r.universe.SchemaForType(r.idType)
		}
	}

	return r.idSchemaType
}

func (r *resourceType[ID, T]) SchemaIdPrototype() schema.TypedPrototype {
	if r.idPrototype == nil {
		r.m.Lock()
		defer r.m.Unlock()

		if r.idPrototype == nil {
			r.idPrototype = r.universe.MakePrototype(r.idType, r.idSchemaType)
		}
	}

	return r.idPrototype
}

func (r *resourceType[ID, T]) CreateInstance() BasicResource {
	return reflect.New(r.resourceType).Interface().(BasicResource)
}

func (r *resourceType[ID, T]) CreateID(name string) BasicResourceID {
	idValue := reflect.New(r.idType)

	idValue.Interface().(IStringResourceID).setValueString(name)

	return idValue.Elem().Interface().(BasicResourceID)
}

func (r *resourceType[ID, T]) IsRuntimeOnly() bool {
	return r.isRuntimeOnly
}

func (r *resourceType[ID, T]) GetID() ResourceTypeID {
	return r.ResourceMetadata.ID
}

func (r *resourceType[ID, T]) Name() string {
	return r.ResourceMetadata.Name
}

func (r *resourceType[ID, T]) IDType() reflect.Type {
	return r.idType
}

func (r *resourceType[ID, T]) ResourceType() reflect.Type {
	return r.resourceType
}

func (r *resourceType[ID, T]) initializeSchema(ts *ResourceTypeSystem, options ...bindnode.Option) {
	r.universe = ts

	r.SchemaResourceType()
	r.SchemaIdType()
}

func derefType[T any]() reflect.Type {
	t := reflect.TypeOf((*T)(nil)).Elem()

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t
}

func derefPointer(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t
}
