package forddb

import (
	"encoding/json"
	"reflect"
	"time"
)

type BasicResourceID interface {
	String() string
	MarshalJSON() ([]byte, error)
}

type ResourceID[T BasicResource] interface {
	BasicResourceID
}

type StringResourceID[T BasicResource] string

func (s StringResourceID[T]) BasicResourceID() BasicResourceID {
	return s
}

func (s StringResourceID[T]) String() string {
	return string(s)
}

func (s StringResourceID[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

func (s *StringResourceID[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*string)(s))
}

func (s *StringResourceID[T]) setValue(value string) {
	*s = StringResourceID[T](value)
}

type stringResourceID interface {
	setValue(value string)
}

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
	return typeSystem.LookupByResourceType(reflect.TypeOf((*T)(nil)).Elem()).ID()
}

func (r *BasicResourceMetadata) GetVersion() int {
	return r.Version
}
