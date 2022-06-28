package ast

import (
	"github.com/shoriwe/gplasma/pkg/lexer"
)

type (
	Statement interface {
		S()
		Node
	}

	AssignStatement struct {
		Statement
		LeftHandSide   IExpression // Identifiers or Selectors
		AssignOperator *lexer.Token
		RightHandSide  IExpression
	}

	DoWhileStatement struct {
		Statement
		Condition IExpression
		Body      []Node
	}

	WhileLoopStatement struct {
		Statement
		Condition IExpression
		Body      []Node
	}

	UntilLoopStatement struct {
		Statement
		Condition IExpression
		Body      []Node
	}

	ForLoopStatement struct {
		Statement
		Receivers []*Identifier
		Source    IExpression
		Body      []Node
	}

	ElifBlock struct {
		Condition IExpression
		Body      []Node
	}

	IfStatement struct {
		Statement
		Condition  IExpression
		Body       []Node
		ElifBlocks []ElifBlock
		Else       []Node
	}

	UnlessStatement struct {
		Statement
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
		Statement
		Target     IExpression
		CaseBlocks []*CaseBlock
		Default    []Node
	}

	ModuleStatement struct {
		Statement
		Name *Identifier
		Body []Node
	}

	FunctionDefinitionStatement struct {
		Statement
		Name      *Identifier
		Arguments []*Identifier
		Body      []Node
	}

	GeneratorDefinitionStatement struct {
		Statement
		Name      *Identifier
		Arguments []*Identifier
		Body      []Node
	}

	InterfaceStatement struct {
		Statement
		Name              *Identifier
		Bases             []IExpression
		MethodDefinitions []*FunctionDefinitionStatement
	}

	ClassStatement struct {
		Statement
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
		Statement
		Body         []Node
		ExceptBlocks []*ExceptBlock
		Else         []Node
		Finally      []Node
	}

	BeginStatement struct {
		Statement
		Body []Node
	}

	EndStatement struct {
		Statement
		Body []Node
	}

	ReturnStatement struct {
		Statement
		Results []IExpression
	}

	YieldStatement struct {
		Statement
		Results []IExpression
	}

	ContinueStatement struct {
		Statement
	}

	BreakStatement struct {
		Statement
	}

	RedoStatement struct {
		Statement
	}

	PassStatement struct {
		Statement
	}

	RaiseStatement struct {
		Statement
		X IExpression
	}

	RequireStatement struct {
		Statement
		X IExpression
	}

	DeleteStatement struct {
		Statement
		X IExpression
	}
)
