package forddb

import (
	"reflect"
)

type ResourceBase[ID ResourceID[T], T Resource[ID]] struct {
	resourceNode ResourceNode

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

func (r *ResourceBase[ID, T]) GetResourceNode() ResourceNode {
	return r.resourceNode
}

func (r *ResourceBase[ID, T]) setResourceNode(node ResourceNode) {
	r.resourceNode = node
}

func (r *ResourceBase[ID, T]) GetResourceType() ResourceType[ID, T] {
	return TypeSystem().LookupByResourceType(reflect.TypeOf((*T)(nil)).Elem()).(ResourceType[ID, T])
}

func (r *ResourceBase[ID, T]) OnBeforeSave(self BasicResource) {
	r.Kind = self.GetResourceTypeID()
}
