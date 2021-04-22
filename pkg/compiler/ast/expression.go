package ast

type Expression interface {
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
	Values []KeyValue
}

type Identifier struct {
	Expression
	String string
}

type BasicLiteralExpression struct {
	Expression
	String string
	Kind   int
}

type BinaryExpression struct {
	Expression
	LeftHandSide  Expression
	Operator      string
	RightHandSide Expression
}

type UnaryExpression struct {
	Expression
	Operator string
	X        Expression
}

type ParenthesesExpression struct {
	Expression
	X Expression
}

type LambdaExpression struct {
	Expression
	Arguments []Identifier
	Code      Expression
}

type GeneratorExpression struct {
	Expression
	Operation Expression
	Variables []Identifier
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
	Index  [2]Expression
}

type AwaitExpression struct {
	Expression
	X Expression
}

type OneLineIfExpression struct {
	Expression
	Result     Expression
	Condition  Expression
	ElseResult Expression
}

type OneLineUnlessExpression struct {
	Expression
	Result     Expression
	Condition  Expression
	ElseResult Expression
}
