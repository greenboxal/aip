package forddb

import (
	"github.com/ipld/go-ipld-prime/datamodel"
)

type BasicResourceLink interface {
	LinkedResourceType() BasicResourceType
	LinkedResourceID() BasicResourceID
}

type ResourceLink[T BasicResourceID] struct {
	ID T `json:"id"`
}

func (r ResourceLink[T]) LinkedResourceType() BasicResourceType {
	return r.ID.BasicResourceType()
}

func (r ResourceLink[T]) LinkedResourceID() BasicResourceID {
	return r.ID
}

func (r ResourceLink[T]) Prototype() datamodel.LinkPrototype {
	return ResourceLinkPrototype[T]{}
}

func (r ResourceLink[T]) String() string {
	return r.ID.String()
}

func (r ResourceLink[T]) Binary() string {
	return r.ID.String()
}

type ResourceLinkPrototype[T BasicResourceID] struct{}

func (r ResourceLinkPrototype[T]) BuildLink(hashsum []byte) datamodel.Link {
	return ResourceLink[T]{ID: NewStringID[T](string(hashsum))}
}
