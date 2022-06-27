package ast

import (
	"github.com/shoriwe/gplasma/pkg/lexer"
)

type Statement interface {
	S()
	Node
}

type AssignStatement struct {
	Statement
	LeftHandSide   IExpression // Identifiers or Selectors
	AssignOperator *lexer.Token
	RightHandSide  IExpression
}

type DoWhileStatement struct {
	Statement
	Condition IExpression
	Body      []Node
}

type WhileLoopStatement struct {
	Statement
	Condition IExpression
	Body      []Node
}

type UntilLoopStatement struct {
	Statement
	Condition IExpression
	Body      []Node
}

type ForLoopStatement struct {
	Statement
	Receivers []*Identifier
	Source    IExpression
	Body      []Node
}

type IfStatement struct {
	Statement
	Condition IExpression
	Body      []Node
	Else      []Node
}

type UnlessStatement struct {
	Statement
	Condition IExpression
	Body      []Node
	Else      []Node
}

type CaseBlock struct {
	Cases []IExpression
	Body  []Node
}

type SwitchStatement struct {
	Statement
	Target     IExpression
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
	Bases             []IExpression
	MethodDefinitions []*FunctionDefinitionStatement
}

type ClassStatement struct {
	Statement
	Name  *Identifier
	Bases []IExpression // Identifiers and selectors
	Body  []Node
}

type ExceptBlock struct {
	Targets     []IExpression
	CaptureName *Identifier
	Body        []Node
}

type RaiseStatement struct {
	Statement
	X IExpression
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
	Results []IExpression
}

type YieldStatement struct {
	Statement
	Results []IExpression
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
