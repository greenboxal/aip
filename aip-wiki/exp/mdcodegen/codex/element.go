package codex

type Element interface {
	ID() int64
	Node() *Node

	Parent() Element

	Children() []Element
}
