package dom

type NodeType string

type Node interface {
	BasicNode() *BasicNode

	NodeID() string
	NodeType() NodeType
	NodeImplementation() Node

	Parent() Node

	ChildNodes() []Node

	HasChildNodes() bool
	Contains(child Node) bool

	InsertBefore(newNode Node, referenceNode Node)
	InsertAfter(newNode Node, referenceNode Node)
	ReplaceChild(oldNode Node, newNode Node)
	AppendChild(child Node)
	RemoveChild(child Node)

	CloneNode() Node
}

type Document interface {
	Node

	RootElement() Element
}

type Element interface {
	Node

	Document() Document

	Children() []Element
}
