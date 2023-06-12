package dom

import (
	"golang.org/x/exp/slices"
)

type BasicNode struct {
	id       string
	kind     NodeType
	self     Node
	parent   Node
	children []Node
}

var _ Node = (*BasicNode)(nil)

func (n *BasicNode) Initialize(self Node, id string, kind NodeType) {
	n.self = self
	n.id = id
	n.kind = kind
}

func (n *BasicNode) BasicNode() *BasicNode    { return n }
func (n *BasicNode) NodeID() string           { return n.id }
func (n *BasicNode) NodeType() NodeType       { return n.kind }
func (n *BasicNode) NodeImplementation() Node { return n.self }
func (n *BasicNode) Parent() Node             { return n.parent }
func (n *BasicNode) HasChildNodes() bool      { return len(n.children) > 0 }
func (n *BasicNode) ChildNodes() []Node       { return n.children }
func (n *BasicNode) Contains(child Node) bool { return slices.Contains(n.children, child) }

func (n *BasicNode) InsertBefore(newNode Node, referenceNode Node) {
	referenceIndex := slices.Index(n.children, referenceNode)

	if referenceIndex == -1 {
		panic("reference element not found")
	}

	n.InsertAt(newNode, referenceIndex)
}

func (n *BasicNode) InsertAfter(newNode Node, referenceNode Node) {
	referenceIndex := slices.Index(n.children, referenceNode)

	if referenceIndex == -1 {
		panic("reference element not found")
	}

	n.InsertAt(newNode, referenceIndex+1)
}

func (n *BasicNode) InsertAt(node Node, idx int) {
	existing := slices.Index(n.children, node)

	n.children = slices.Insert(n.children, idx, node)

	if existing != -1 {
		n.children = slices.Delete(n.children, existing, existing+1)
	} else {
		node.BasicNode().attachFromParent(n)
	}
}

func (n *BasicNode) ReplaceChild(oldNode Node, newNode Node) {
	idx := slices.Index(n.children, oldNode)

	if idx == -1 {
		panic("old node not found")
	}

	n.children[idx] = newNode
}

func (n *BasicNode) AppendChild(child Node) {
	n.InsertAt(child, len(n.children))
}

func (n *BasicNode) RemoveChild(child Node) {
	idx := slices.Index(n.children, child)

	if idx == -1 {
		panic("not not found")
	}

	n.children = slices.Delete(n.children, idx, idx+1)

	child.BasicNode().detachFromParent()
}

func (n *BasicNode) CloneNode() Node {
	cloned := &BasicNode{}

	cloned.Initialize(cloned, n.id, n.kind)

	for _, child := range n.children {
		clonedChild := child.CloneNode()

		cloned.AppendChild(clonedChild)
	}

	return cloned
}

func (n *BasicNode) attachFromParent(parent Node) {
	n.parent = nil
}

func (n *BasicNode) detachFromParent() {
	n.parent = nil
}
