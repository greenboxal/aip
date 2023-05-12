package forddb

import (
	"github.com/ipld/go-ipld-prime"
)

type BasicResourcePointer interface {
	BasicType() BasicResourceType
	BasicID() BasicResourceID
	AsLink() ipld.Link
}

type ResourcePointer[ID ResourceID[T], T Resource[ID]] struct {
	BasicResourcePointer

	TypeID TypeID    `json:"type_id"`
	ID     ID        `json:"resource_id"`
	Link   ipld.Link `json:"resource_link"`
}

func (rp *ResourcePointer[ID, T]) Type() ResourceType[ID, T] {
	return rp.Type().(ResourceType[ID, T])
}
func (rp *ResourcePointer[ID, T]) BasicType() BasicResourceType { return rp.Type() }
func (rp *ResourcePointer[ID, T]) BasicID() BasicResourceID     { return rp.ID }
func (rp *ResourcePointer[ID, T]) AsLink() ipld.Link            { return rp.Link }

type BasicResourceSlot interface {
	BasicResourcePointer

	GetResource() BasicResource
	SetResource(resource BasicResource)
}

type ResourceSlot[ID ResourceID[T], T Resource[ID]] struct {
	ResourcePointer[ID, T]

	Resource T `json:"-"`
}

func (r *ResourceSlot[ID, T]) GetResource() BasicResource         { return r.Resource }
func (r *ResourceSlot[ID, T]) SetResource(resource BasicResource) { r.Resource = resource.(T) }
