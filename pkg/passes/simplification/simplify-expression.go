package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"reflect"
)

func SimplifyExpression(expression ast.IExpression) ast.IExpression {
	switch expr := expression.(type) {
	case *ast.Identifier:
		return expr
	case *ast.BasicLiteralExpression:
		return simplifyBasicLiteral(expr)
	case *ast.ArrayExpression:
		return simplifyArray(expr)
	case *ast.TupleExpression:
		return simplifyTuple(expr)
	case *ast.HashExpression:
		return simplifyHash(expr)
	case *ast.BinaryExpression:
		return simplifyBinary(expr)
	case *ast.UnaryExpression:
		return simplifyUnary(expr)
	case *ast.ParenthesesExpression:
		return SimplifyExpression(expr.X)
	case *ast.LambdaExpression:
		return simplifyLambda(expr)
	case *ast.GeneratorExpression:
		return simplifyGenerator(expr)
	case *ast.SelectorExpression:
		return expr
	case *ast.MethodInvocationExpression:
		return simplifyMethodInvocation(expr)
	case *ast.IndexExpression:
		return simplifyIndex(expr)
	case *ast.IfOneLinerExpression:
		return simplifyIfOneLiner(expr)
	case *ast.UnlessOneLinerExpression:
		return simplifyUnlessOneLiner(expr)
	default:
		panic(reflect.TypeOf(expr))
	}
}
