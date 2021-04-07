package ast

type Expression interface {
	Node
}

type Identifier struct {
	Expression
	String string
}

type BasicLiteralExpression struct {
	Expression
	String string
	Kind   rune
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
	Identifier Identifier
}

type MethodInvocationExpression struct {
	Expression
	Function  Expression
	Arguments []Expression
}

type IndexExpression struct {
	Expression
	Source Expression
	Index  [3]Expression
}

type GoExpression struct {
	Expression
	FunctionInvocation Expression
}

type ReturnExpression struct {
	Expression
	Results []Expression
}

type YieldExpression struct {
	Expression
	Results []Expression
}

type SuperInvocationExpression struct {
	Expression
	Arguments []Expression
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

type RetryExpression struct {
	Expression
	Target Identifier
}

type BreakExpression struct {
	Expression
	Target Identifier
}

type RedoExpression struct {
	Expression
	Target Identifier
}

type PassExpression struct {
	Expression
}
