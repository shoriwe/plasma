package ast

import (
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

type IExpression interface {
	Node
	E()
}

type ArrayExpression struct {
	IExpression
	Values []IExpression
}

type TupleExpression struct {
	IExpression
	Values []IExpression
}

type KeyValue struct {
	Key   IExpression
	Value IExpression
}

type HashExpression struct {
	IExpression
	Values []*KeyValue
}

type Identifier struct {
	IExpression
	Token *lexer2.Token
}

type BasicLiteralExpression struct {
	IExpression
	Token       *lexer2.Token
	Kind        lexer2.Kind
	DirectValue lexer2.DirectValue
}

type BinaryExpression struct {
	IExpression
	LeftHandSide  IExpression
	Operator      *lexer2.Token
	RightHandSide IExpression
}

type UnaryExpression struct {
	IExpression
	Operator *lexer2.Token
	X        IExpression
}

type ParenthesesExpression struct {
	IExpression
	X IExpression
}

type LambdaExpression struct {
	IExpression
	Arguments []*Identifier
	Code      IExpression
}

type GeneratorExpression struct {
	IExpression
	Operation IExpression
	Receivers []*Identifier
	Source    IExpression
}

type SelectorExpression struct {
	IExpression
	X          IExpression
	Identifier *Identifier
}

type MethodInvocationExpression struct {
	IExpression
	Function  IExpression
	Arguments []IExpression
}

type IndexExpression struct {
	IExpression
	Source IExpression
	Index  IExpression
}

type IfOneLinerExpression struct {
	IExpression
	Result     IExpression
	Condition  IExpression
	ElseResult IExpression
}

type UnlessOneLinerExpression struct {
	IExpression
	Result     IExpression
	Condition  IExpression
	ElseResult IExpression
}
