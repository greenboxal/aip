package forddb

import (
	"encoding/json"
	"reflect"

	"github.com/ipld/go-ipld-prime/schema"
)

type ResourceTypeID string

func (s ResourceTypeID) BasicResourceID() BasicResourceID {
	return s
}

func (s ResourceTypeID) String() string {
	return string(s)
}

func (s ResourceTypeID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

func (s *ResourceTypeID) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*string)(s))
}

func (s *ResourceTypeID) setValue(value string) {
	*s = ResourceTypeID(value)
}

func (i ResourceTypeID) Name() string {
	return i.String()
}

func (i ResourceTypeID) Type() BasicResourceType {
	return typeSystem.LookupByID(i)
}

type BasicResourceType interface {
	BasicResource

	GetID() ResourceTypeID

	ID() ResourceTypeID
	Name() string
	IDType() reflect.Type
	ResourceType() reflect.Type
	IsRuntimeOnly() bool

	New() BasicResource
	MakeId(name string) BasicResourceID

	SchemaIdType() schema.Type
	SchemaIdPrototype() schema.TypedPrototype

	SchemaResourceType() schema.Type
	SchemaResourcePrototype() schema.TypedPrototype

	initializeSchema(ts *ResourceTypeSystem)
}

type ResourceType[ID ResourceID[T], T Resource[ID]] interface {
	BasicResourceType
}

func LookupTypeByName(name string) BasicResourceType {
	return typeSystem.LookupByID(typeType.MakeId(name).(ResourceTypeID))
}

func DefineResourceType[ID ResourceID[T], T Resource[ID]](name string) ResourceType[ID, T] {
	t := &resourceType[ID, T]{
		idType:       derefType[ID](),
		resourceType: derefType[T](),
	}

	t.ResourceMetadata.ID = typeType.MakeId(name).(ResourceTypeID)
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
}

func (r *resourceType[ID, T]) SchemaResourceType() schema.Type {
	return r.resourceSchemaType
}

func (r *resourceType[ID, T]) SchemaResourcePrototype() schema.TypedPrototype {
	return r.resourcePrototype
}

func (r *resourceType[ID, T]) SchemaIdType() schema.Type {
	return r.idSchemaType
}

func (r *resourceType[ID, T]) SchemaIdPrototype() schema.TypedPrototype {
	return r.idPrototype
}

func (r *resourceType[ID, T]) New() BasicResource {
	return reflect.New(r.resourceType).Interface().(BasicResource)
}

func (r *resourceType[ID, T]) MakeId(name string) BasicResourceID {
	idValue := reflect.New(r.idType)

	idValue.Interface().(stringResourceID).setValue(name)

	return idValue.Elem().Interface().(BasicResourceID)
}

func (r *resourceType[ID, T]) IsRuntimeOnly() bool {
	return r.isRuntimeOnly
}

func (r *resourceType[ID, T]) ID() ResourceTypeID {
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

func (r *resourceType[ID, T]) initializeSchema(ts *ResourceTypeSystem) {
	//idType := ts.SchemaForType(r.idType)
	//idPrototype := bindnode.Prototype((*T)(nil), idType)

	//r.idSchemaType = idType
	//r.idPrototype = idPrototype

	//schemaType := ts.SchemaForType(r.resourceType)
	//resourcePrototype := bindnode.Prototype((*T)(nil), schemaType)

	//r.resourceSchemaType = schemaType
	//r.resourcePrototype = resourcePrototype
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

type stringResourceID interface {
	setValue(value string)
}
