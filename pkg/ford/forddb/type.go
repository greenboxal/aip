package forddb

import (
	"encoding/json"
	"reflect"
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

func derefType[T any]() reflect.Type {
	t := reflect.TypeOf((*T)(nil)).Elem()

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t
}
