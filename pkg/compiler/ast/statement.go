package ast

import (
	"github.com/shoriwe/gplasma/pkg/compiler/lexer"
)

type Statement interface {
	S()
	Node
}

type AssignStatement struct {
	Statement
	LeftHandSide   Expression // Identifiers or Selectors
	AssignOperator *lexer.Token
	RightHandSide  Expression
}

type DeferStatement struct {
	Statement
	X *MethodInvocationExpression
}

type DoWhileStatement struct {
	Statement
	Condition Expression
	Body      []Node
}
type WhileLoopStatement struct {
	Statement
	Condition Expression
	Body      []Node
}

type UntilLoopStatement struct {
	Statement
	Condition Expression
	Body      []Node
}

type ForLoopStatement struct {
	Statement
	Receivers []*Identifier
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
	ElifBlocks []*ElifBlock
	Else       []Node
}

type UnlessStatement struct {
	Statement
	Condition  Expression
	Body       []Node
	ElifBlocks []*ElifBlock
	Else       []Node
}

type CaseBlock struct {
	Cases []Expression
	Body  []Node
}

type SwitchStatement struct {
	Statement
	Target     Expression
	CaseBlocks []*CaseBlock
	Default    []Node
}

type ModuleStatement struct {
	Statement
	Name *Identifier
	Body []Node
}

type FunctionDefinitionStatement struct {
	Statement
	Name      *Identifier
	Arguments []*Identifier
	Body      []Node
}

type InterfaceStatement struct {
	Statement
	Name              *Identifier
	Bases             []Expression
	MethodDefinitions []*FunctionDefinitionStatement
}

type ClassStatement struct {
	Statement
	Name  *Identifier
	Bases []Expression // Identifiers and selectors
	Body  []Node
}

type ExceptBlock struct {
	Targets     []Expression
	CaptureName *Identifier
	Body        []Node
}

type RaiseStatement struct {
	Statement
	X Expression
}
type TryStatement struct {
	Statement
	Body         []Node
	ExceptBlocks []*ExceptBlock
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

type ReturnStatement struct {
	Statement
	Results []Expression
}

type YieldStatement struct {
	Statement
	Results []Expression
}

type SuperInvocationStatement struct {
	Statement
	Arguments []Expression
}

type ContinueStatement struct {
	Statement
}

type BreakStatement struct {
	Statement
}

type RedoStatement struct {
	Statement
}

type PassStatement struct {
	Statement
}
