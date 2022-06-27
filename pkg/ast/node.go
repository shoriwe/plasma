package ast

type Node interface {
	N()
}

type Program struct {
	Node
	Begin *BeginStatement
	End   *EndStatement
	Body  []Node
}
