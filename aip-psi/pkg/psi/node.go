package psi

import "github.com/greenboxal/aip/aip-forddb/pkg/forddb"

type Node interface {
	forddb.Resource[ElementID]

	Children() []ElementID
	AsElement() *Element
}

type ElementID struct {
	forddb.StringResourceID[Node] `ipld:",inline"`
}

type Element struct {
	forddb.ResourceBase[ElementID, Node] `json:"metadata"`

	Children []forddb.ResourceLink[ElementID] `json:"children"`
}

var _ Node = (*Element)(nil)

func (e *Element) AsElement() *Element {
	return e
}
