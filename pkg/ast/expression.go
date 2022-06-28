package ast

import (
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

type (
	IExpression interface {
		Node
		E()
	}

	ArrayExpression struct {
		IExpression
		Values []IExpression
	}

	TupleExpression struct {
		IExpression
		Values []IExpression
	}

	KeyValue struct {
		Key   IExpression
		Value IExpression
	}

	HashExpression struct {
		IExpression
		Values []*KeyValue
	}

	Identifier struct {
		IExpression
		Token *lexer2.Token
	}

	BasicLiteralExpression struct {
		IExpression
		Token       *lexer2.Token
		Kind        lexer2.Kind
		DirectValue lexer2.DirectValue
	}

	BinaryExpression struct {
		IExpression
		LeftHandSide  IExpression
		Operator      *lexer2.Token
		RightHandSide IExpression
	}

	UnaryExpression struct {
		IExpression
		Operator *lexer2.Token
		X        IExpression
	}

	ParenthesesExpression struct {
		IExpression
		X IExpression
	}

	LambdaExpression struct {
		IExpression
		Arguments []*Identifier
		Code      IExpression
	}

	GeneratorExpression struct {
		IExpression
		Operation IExpression
		Receivers []*Identifier
		Source    IExpression
	}

	SelectorExpression struct {
		IExpression
		X          IExpression
		Identifier *Identifier
	}

	MethodInvocationExpression struct {
		IExpression
		Function  IExpression
		Arguments []IExpression
	}

	IndexExpression struct {
		IExpression
		Source IExpression
		Index  IExpression
	}

	IfOneLinerExpression struct {
		IExpression
		Result     IExpression
		Condition  IExpression
		ElseResult IExpression
	}

	UnlessOneLinerExpression struct {
		IExpression
		Result     IExpression
		Condition  IExpression
		ElseResult IExpression
	}

	SuperExpression struct {
		IExpression
		X IExpression
	}
)
