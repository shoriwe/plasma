package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func simplifyIfOneLiner(if_ *ast.IfOneLinerExpression) *ast2.IfOneLiner {
	return &ast2.IfOneLiner{
		Condition: simplifyExpression(if_.Condition),
		Result:    simplifyExpression(if_.Result),
		Else:      simplifyExpression(if_.ElseResult),
	}
}

func simplifyUnlessOneLiner(unless *ast.UnlessOneLinerExpression) *ast2.IfOneLiner {
	return &ast2.IfOneLiner{
		Condition: &ast2.Unary{
			Operator: ast2.Not,
			X:        simplifyExpression(unless.Condition),
		},
		Result: simplifyExpression(unless.Result),
		Else:   simplifyExpression(unless.ElseResult),
	}
}
