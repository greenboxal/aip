package psi

import "github.com/ipld/go-ipld-prime"

type NodeList struct {
	Nodes []Node
}

type Node interface {
	Attributes() ipld.Node
	Children() NodeList
}
