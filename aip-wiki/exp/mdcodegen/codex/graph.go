package codex

import (
	"gonum.org/v1/gonum/graph/multi"
)

type Graph struct {
	*multi.DirectedGraph
}

func NewGraph() *Graph {
	return &Graph{
		DirectedGraph: multi.NewDirectedGraph(),
	}
}

func (g *Graph) AddNode(n Element) {
	n.Node().attachToGraph(g)
	g.AddNode(n)
}

func (g *Graph) RemoveNode(n Element) {
	g.RemoveNode(n)
	n.Node().detachFromGraph(nil)
}
