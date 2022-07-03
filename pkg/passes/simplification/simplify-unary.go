package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/common"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func evalUnaryExpression(x *ast.BasicLiteralExpression, operator *lexer.Token) ast.IExpression {
	switch operator.DirectValue {
	case lexer.NegateBits:
		switch {
		case literalIsInteger(x):
			newValue := common.IntegerNegateBitsToken(x.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		}
	}
	return &ast.UnaryExpression{
		Operator: operator,
		X:        x,
	}
}

func simplifyUnary(unary *ast.UnaryExpression) ast.IExpression {
	newX := SimplifyExpression(unary.X)
	switch x := newX.(type) {
	case *ast.BasicLiteralExpression:
		return evalUnaryExpression(x, unary.Operator)
	}
	return &ast.UnaryExpression{
		Operator: unary.Operator,
		X:        newX,
	}
}
