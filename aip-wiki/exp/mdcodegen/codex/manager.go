package codex

type Manager struct {
}

type Scope struct {
	Node

	Roots []Element
}

type Module struct {
	Node

	Types []*Type
}

type Type struct {
	Node

	Methods    []*Method
	Properties []*Property
}

type Method struct {
	Node
}

type Property struct {
	Node
}
