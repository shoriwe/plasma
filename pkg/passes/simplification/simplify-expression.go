package simplification

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
	"reflect"
)

func simplifyExpression(expr ast.Expression) ast2.Expression {
	if expr == nil {
		return &ast2.None{}
	}
	switch e := expr.(type) {
	case *ast.ArrayExpression:
		return simplifyArray(e)
	case *ast.TupleExpression:
		return simplifyTuple(e)
	case *ast.HashExpression:
		return simplifyHash(e)
	case *ast.Identifier:
		return simplifyIdentifier(e)
	case *ast.BasicLiteralExpression:
		return simplifyLiteral(e)
	case *ast.BinaryExpression:
		return simplifyBinary(e)
	case *ast.UnaryExpression:
		return simplifyUnary(e)
	case *ast.ParenthesesExpression:
		return simplifyParentheses(e)
	case *ast.LambdaExpression:
		return simplifyLambda(e)
	case *ast.GeneratorExpression:
		return simplifyGeneratorExpr(e)
	case *ast.SelectorExpression:
		return simplifySelector(e)
	case *ast.MethodInvocationExpression:
		return simplifyCall(e)
	case *ast.IndexExpression:
		return simplifyIndex(e)
	case *ast.IfOneLinerExpression:
		return simplifyIfOneLiner(e)
	case *ast.UnlessOneLinerExpression:
		return simplifyUnlessOneLiner(e)
	case *ast.SuperExpression:
		return simplifySuper(e)
	default:
		panic(fmt.Sprintf("unknown expression type %s", reflect.TypeOf(expr).String()))
	}
}
