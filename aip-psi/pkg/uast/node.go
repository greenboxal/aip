package uast

import "github.com/greenboxal/aip/aip-psi/pkg/psi"

type NodeList []Node

type Node interface {
	Element() psi.Element

	Children() NodeList
}

type EmitContext interface {
	Push(node Node)
}

type Renderer interface {
	Render(ctx EmitContext, node Node) (bool, error)
}
