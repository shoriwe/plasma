package ast

import "github.com/shoriwe/gplasma/pkg/compiler/lexer"

type Expression interface {
	E()
	Node
}

type ArrayExpression struct {
	Expression
	Values []Expression
}

type TupleExpression struct {
	Expression
	Values []Expression
}

type KeyValue struct {
	Key   Expression
	Value Expression
}

type HashExpression struct {
	Expression
	Values []*KeyValue
}

type Identifier struct {
	Expression
	Token *lexer.Token
}

type BasicLiteralExpression struct {
	Expression
	Token       *lexer.Token
	Kind        uint8
	DirectValue uint8
}

type BinaryExpression struct {
	Expression
	LeftHandSide  Expression
	Operator      *lexer.Token
	RightHandSide Expression
}

type UnaryExpression struct {
	Expression
	Operator *lexer.Token
	X        Expression
}

type ParenthesesExpression struct {
	Expression
	X Expression
}

type LambdaExpression struct {
	Expression
	Arguments []*Identifier
	Code      Expression
}

type GeneratorExpression struct {
	Expression
	Operation Expression
	Receivers []*Identifier
	Source    Expression
}

type SelectorExpression struct {
	Expression
	X          Expression
	Identifier *Identifier
}

type MethodInvocationExpression struct {
	Expression
	Function  Expression
	Arguments []Expression
}

type IndexExpression struct {
	Expression
	Source Expression
	Index  Expression
}

type IfOneLinerExpression struct {
	Expression
	Result     Expression
	Condition  Expression
	ElseResult Expression
}

type UnlessOneLinerExpression struct {
	Expression
	Result     Expression
	Condition  Expression
	ElseResult Expression
}
