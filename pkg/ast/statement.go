package ast

import (
	"github.com/shoriwe/gplasma/pkg/lexer"
)

type (
	IStatement interface {
		S()
		Node
	}

	AssignStatement struct {
		IStatement
		LeftHandSide   IExpression // Identifiers or Selectors
		AssignOperator *lexer.Token
		RightHandSide  IExpression
	}

	DoWhileStatement struct {
		IStatement
		Condition IExpression
		Body      []Node
	}

	WhileLoopStatement struct {
		IStatement
		Condition IExpression
		Body      []Node
	}

	UntilLoopStatement struct {
		IStatement
		Condition IExpression
		Body      []Node
	}

	ForLoopStatement struct {
		IStatement
		Receivers []*Identifier
		Source    IExpression
		Body      []Node
	}

	ElifBlock struct {
		Condition IExpression
		Body      []Node
	}

	IfStatement struct {
		IStatement
		Condition  IExpression
		Body       []Node
		ElifBlocks []ElifBlock
		Else       []Node
	}

	UnlessStatement struct {
		IStatement
		Condition  IExpression
		Body       []Node
		ElifBlocks []ElifBlock
		Else       []Node
	}

	CaseBlock struct {
		Cases []IExpression
		Body  []Node
	}

	SwitchStatement struct {
		IStatement
		Target     IExpression
		CaseBlocks []*CaseBlock
		Default    []Node
	}

	ModuleStatement struct {
		IStatement
		Name *Identifier
		Body []Node
	}

	FunctionDefinitionStatement struct {
		IStatement
		Name      *Identifier
		Arguments []*Identifier
		Body      []Node
	}

	GeneratorDefinitionStatement struct {
		IStatement
		Name      *Identifier
		Arguments []*Identifier
		Body      []Node
	}

	InterfaceStatement struct {
		IStatement
		Name              *Identifier
		Bases             []IExpression
		MethodDefinitions []*FunctionDefinitionStatement
	}

	ClassStatement struct {
		IStatement
		Name  *Identifier
		Bases []IExpression // Identifiers and selectors
		Body  []Node
	}

	ExceptBlock struct {
		Targets     []IExpression
		CaptureName *Identifier
		Body        []Node
	}

	TryStatement struct {
		IStatement
		Body         []Node
		ExceptBlocks []*ExceptBlock
		Else         []Node
		Finally      []Node
	}

	BeginStatement struct {
		IStatement
		Body []Node
	}

	EndStatement struct {
		IStatement
		Body []Node
	}

	ReturnStatement struct {
		IStatement
		Results []IExpression
	}

	YieldStatement struct {
		IStatement
		Results []IExpression
	}

	ContinueStatement struct {
		IStatement
	}

	BreakStatement struct {
		IStatement
	}

	RedoStatement struct {
		IStatement
	}

	PassStatement struct {
		IStatement
	}

	RaiseStatement struct {
		IStatement
		X IExpression
	}

	RequireStatement struct {
		IStatement
		X IExpression
	}

	DeleteStatement struct {
		IStatement
		X IExpression
	}
)
