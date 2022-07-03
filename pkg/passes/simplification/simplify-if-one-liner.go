package simplification

import "github.com/shoriwe/gplasma/pkg/ast"

func simplifyIfOneLiner(i *ast.IfOneLinerExpression) *ast.IfOneLinerExpression {
	return &ast.IfOneLinerExpression{
		Result:     SimplifyExpression(i.Result),
		Condition:  SimplifyExpression(i.Condition),
		ElseResult: SimplifyExpression(i.ElseResult),
	}
}
