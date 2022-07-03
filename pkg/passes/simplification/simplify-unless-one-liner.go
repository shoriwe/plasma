package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func simplifyUnlessOneLiner(u *ast.UnlessOneLinerExpression) *ast.IfOneLinerExpression {
	return &ast.IfOneLinerExpression{
		Result: SimplifyExpression(u.Result),
		Condition: SimplifyExpression(&ast.UnaryExpression{
			Operator: &lexer.Token{
				Contents:    []rune(lexer.NotString),
				DirectValue: lexer.Not,
				Kind:        lexer.Comparator,
			},
			X: u.Condition,
		}),
		ElseResult: SimplifyExpression(u.ElseResult),
	}
}
