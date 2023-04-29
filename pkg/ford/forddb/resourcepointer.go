package forddb

type BasicResourcePointer interface {
	BasicType() BasicResourceType
	BasicID() BasicResourceID
}

type BasicResourceSlot interface {
	BasicResourcePointer

	GetResource() BasicResource
	SetResource(resource BasicResource)
}

type ResourceSlot[ID ResourceID[T], T Resource[ID]] struct {
	ResourcePointer[ID, T]

	Resource T `json:"-"`
}

func (r *ResourceSlot[ID, T]) GetResource() BasicResource {
	return r.Resource
}

func (r *ResourceSlot[ID, T]) SetResource(resource BasicResource) {
	r.Resource = resource.(T)
}

type ResourcePointer[ID ResourceID[T], T Resource[ID]] struct {
	BasicResourcePointer

	TypeID ResourceTypeID `json:"type_id"`
	ID     ID             `json:"resource_id"`
}

func (rp *ResourcePointer[ID, T]) Type() ResourceType[ID, T] {
	return rp.Type().(ResourceType[ID, T])
}

func (rp *ResourcePointer[ID, T]) BasicType() BasicResourceType {
	return rp.Type()
}

func (rp *ResourcePointer[ID, T]) BasicID() BasicResourceID {
	return rp.ID
}
