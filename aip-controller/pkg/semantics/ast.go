package semantics

type AstNode interface {
}

type UnaryOperationType int

const (
	UnOpInvalid UnaryOperationType = iota
	UnOpNegate
	UnOpSquare
	UnOpRoot
)

type BinaryOperationType int

const (
	BinOpInvalid BinaryOperationType = iota
	BinOpAdd
	BinOpSub
)

type UnaryOperation interface {
	AstNode

	Op() UnaryOperationType
	Operand() AstNode
}

type BinaryOperation interface {
	AstNode

	Left() AstNode
	Right() AstNode

	Op() BinaryOperationType
}
