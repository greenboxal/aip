package codex

import (
	"golang.org/x/exp/slices"
	"gonum.org/v1/gonum/graph"
)

type NodeType string

type EdgeType interface {
	Name() string

	ReverseType() EdgeType
	Reverse(e *Edge) *Edge
}

type Edge struct {
	g    *Graph
	id   int64
	typ  EdgeType
	from *Node
	to   *Node
}

func (e *Edge) To() graph.Node {
	return e.to
}

func (e *Edge) From() graph.Node {
	return e.from
}

func (e *Edge) ID() int64 {
	return e.id
}

func (e *Edge) ReversedLine() graph.Line {
	rt := e.typ.ReverseType()

	if rt == nil {
		panic("not supported")
	}

	return e.typ.Reverse(e)
}

func (e *Edge) attachToGraph(g *Graph) {
	if e.g == g {
		return
	}

	if g == nil {
		e.detachFromGraph(nil)
		return
	}

	if e.g != nil {
		panic("node already attached to a graph")
	}

	e.g = g
	e.id = g.NewLine(e.from, e.to).ID()

	e.from.attachToGraph(g)
	e.to.attachToGraph(g)
}

func (e *Edge) detachFromGraph(g *Graph) {
	if e.g == nil {
		return
	}

	if e.g != g {
		return
	}

	e.from.detachFromGraph(g)
	e.to.detachFromGraph(g)

	e.g = nil
}

type EdgeID struct {
	Type  EdgeType
	Index int
}

type Node struct {
	g        *Graph
	id       int64
	typ      NodeType
	parent   *Node
	children []Element
	edges    map[EdgeID]*Edge
	self     Element
}

func (n *Node) Children() []Element {
	return n.children
}

func (n *Node) ID() int64 {
	return n.id
}

func (n *Node) Node() *Node {
	return n
}

func (n *Node) NodeType() NodeType {
	return n.typ
}

func (n *Node) Parent() Element {
	return n.parent.self
}

func (n *Node) SetParent(parent Element) {
	parentNode := parent.Node()

	if n.parent == parentNode {
		return
	}

	if n.parent != nil {
		n.parent.removeChildNode(n)
	}

	n.parent = parentNode

	if n.parent != nil {
		n.attachToGraph(n.parent.g)
	} else {
		n.detachFromGraph(nil)
	}
}

func (n *Node) Edges() map[EdgeID]*Edge {
	return n.edges
}

func (n *Node) Edge(id EdgeID) *Edge {
	return n.edges[id]
}

func (n *Node) attachToGraph(g *Graph) {
	if n.g == g {
		return
	}

	if g == nil {
		n.detachFromGraph(nil)
		return
	}

	if n.g != nil {
		panic("node already attached to a graph")
	}

	n.g = g
	n.id = g.NewNode().ID()

	for _, e := range n.edges {
		e.attachToGraph(g)
	}
}

func (n *Node) detachFromGraph(g *Graph) {
	if n.g == nil {
		return
	}

	if n.g != g {
		return
	}

	for _, e := range n.edges {
		e.detachFromGraph(n.g)
	}

	n.g = nil
}

func (n *Node) addChildNode(child *Node) {
	idx := slices.Index(n.children, child)

	if idx != -1 {
		return
	}

	n.children = append(n.children, child)
}

func (n *Node) removeChildNode(child *Node) {
	idx := slices.Index(n.children, child)

	if idx == -1 {
		return
	}

	n.children = slices.Delete(n.children, idx, idx+1)
}
