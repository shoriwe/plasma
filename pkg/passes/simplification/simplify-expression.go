package simplification

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
	"reflect"
)

func (simp *simplify) simplifyExpression(expr ast.Expression) ast2.Expression {
	if expr == nil {
		return &ast2.None{}
	}
	switch e := expr.(type) {
	case *ast.ArrayExpression:
		return simp.simplifyArray(e)
	case *ast.TupleExpression:
		return simp.simplifyTuple(e)
	case *ast.HashExpression:
		return simp.simplifyHash(e)
	case *ast.Identifier:
		return simp.simplifyIdentifier(e)
	case *ast.BasicLiteralExpression:
		return simp.simplifyLiteral(e)
	case *ast.BinaryExpression:
		return simp.simplifyBinary(e)
	case *ast.UnaryExpression:
		return simp.simplifyUnary(e)
	case *ast.ParenthesesExpression:
		return simp.simplifyParentheses(e)
	case *ast.LambdaExpression:
		return simp.simplifyLambda(e)
	case *ast.GeneratorExpression:
		return simp.simplifyGeneratorExpr(e)
	case *ast.SelectorExpression:
		return simp.simplifySelector(e)
	case *ast.MethodInvocationExpression:
		return simp.simplifyCall(e)
	case *ast.IndexExpression:
		return simp.simplifyIndex(e)
	case *ast.IfOneLinerExpression:
		return simp.simplifyIfOneLiner(e)
	case *ast.UnlessOneLinerExpression:
		return simp.simplifyUnlessOneLiner(e)
	case *ast.SuperExpression:
		return simp.simplifySuper(e)
	default:
		panic(fmt.Sprintf("unknown expression type %s", reflect.TypeOf(expr).String()))
	}
}
