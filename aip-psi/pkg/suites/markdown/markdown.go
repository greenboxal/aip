package markdown

import (
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"

	"github.com/greenboxal/aip/aip-psi/pkg/psi"
	"github.com/greenboxal/aip/aip-psi/pkg/uast"
)

type Leaf struct {
	leaf *ast.Leaf
}

func (c *Leaf) Element() psi.Element {
	return nil
}

func (c *Leaf) Children() uast.NodeList {
	return uast.NodeList{}
}

type Container struct {
	container *ast.Container
	children  uast.NodeList
}

func (c *Container) Element() psi.Element {
	return nil
}

func (c *Container) Children() uast.NodeList {
	return c.children
}

func Parse(s string) uast.Node {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	n := p.Parse([]byte(s))

	return ParseNode(n)
}

func ParseNode(n ast.Node) uast.Node {
	if c := n.AsContainer(); c != nil {
		children := make([]uast.Node, len(c.Children))

		for i, child := range c.Children {
			children[i] = ParseNode(child)
		}

		return &Container{
			container: c,
			children:  children,
		}
	}

	return &Leaf{
		leaf: n.AsLeaf(),
	}
}
