package ast

type Node interface {
}

type Program struct {
	Node
	Begin *BeginStatement
	End   *EndStatement
	Body  []Node
}
