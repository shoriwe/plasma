package ast

type Statement interface {
	Node
}

type WhileLoopStatement struct {
	Statement
	Name      Identifier
	Condition Expression
	Body      []Node
}

type UntilLoopStatement struct {
	Statement
	Name      Identifier
	Condition Expression
	Body      []Node
}

type ForLoopStatement struct {
	Statement
	Name      Identifier
	Receivers []Identifier
	Source    Expression
	Body      []Node
}

type ElifBlock struct {
	Condition Expression
	Body      []Node
}

type IfStatement struct {
	Statement
	Condition  Expression
	Body       []Node
	ElifBlocks []ElifBlock
	Else       []Node
}

type UnlessStatement struct {
	Statement
	Condition  Expression
	Body       []Node
	ElifBlocks []ElifBlock
	Else       []Node
}

type CaseBlock struct {
	Cases []Expression
	Body  []Node
}

type SwitchStatement struct {
	Statement
	Name       Identifier
	Target     Expression
	CaseBlocks []CaseBlock
	Else       []Node
}

type ModuleStatement struct {
	Statement
	Name Identifier
	Body []Node
}

type FunctionDefinitionStatement struct {
	Statement
	Name      Identifier
	Arguments []Expression
	Body      []Node
}

type AsyncFunctionDefinitionStatement struct {
	Statement
	Name      Identifier
	Arguments []Expression
	Body      []Node
}

type StructStatement struct {
	Statement
	Name   Identifier
	Fields []Identifier
}

type InterfaceStatement struct {
	Statement
	Name                   Identifier
	Bases                  []Identifier
	MethodDefinitions      []FunctionDefinitionStatement
	AsyncMethodDefinitions []AsyncFunctionDefinitionStatement
}

type ClassStatement struct {
	Statement
	Name  Identifier
	Bases []Identifier
	Body  []Node
}

type EnumStatement struct {
	Statement
	Name            Identifier
	EnumIdentifiers []Expression
}

type ExceptBlock struct {
	Targets     []Identifier
	CaptureName Identifier
	Body        []Node
}

type TryStatement struct {
	Statement
	Body         []Node
	ExceptBlocks []ExceptBlock
	Else         []Node
	Finally      []Node
}

type BeginStatement struct {
	Statement
	Body []Node
}

type EndStatement struct {
	Statement
	Body []Node
}
