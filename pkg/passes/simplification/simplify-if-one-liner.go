package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simp *simplify) simplifyIfOneLiner(if_ *ast.IfOneLinerExpression) *ast2.IfOneLiner {
	return &ast2.IfOneLiner{
		Condition: simp.simplifyExpression(if_.Condition),
		Result:    simp.simplifyExpression(if_.Result),
		Else:      simp.simplifyExpression(if_.ElseResult),
	}
}

func (simp *simplify) simplifyUnlessOneLiner(unless *ast.UnlessOneLinerExpression) *ast2.IfOneLiner {
	return &ast2.IfOneLiner{
		Condition: &ast2.Unary{
			Operator: ast2.Not,
			X:        simp.simplifyExpression(unless.Condition),
		},
		Result: simp.simplifyExpression(unless.Result),
		Else:   simp.simplifyExpression(unless.ElseResult),
	}
}
