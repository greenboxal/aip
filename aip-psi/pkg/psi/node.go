package psi

import "github.com/ipld/go-ipld-prime"

type ElementList struct {
	Nodes []Element
}

type Element interface {
	Attributes() ipld.Node
	Children() ElementList
}
