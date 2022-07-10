package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simp *simplify) simplifyUnless(unless *ast.UnlessStatement) *ast2.If {
	body := make([]ast2.Node, 0, len(unless.Body))
	for _, node := range unless.Body {
		body = append(body, simp.simplifyNode(node))
	}
	root := &ast2.If{
		Condition: &ast2.Unary{
			Operator: ast2.Not,
			X:        simp.simplifyExpression(unless.Condition),
		},
		Body: body,
		Else: nil,
	}
	lastIf := root
	for _, elif := range unless.ElifBlocks {
		elifBody := make([]ast2.Node, 0, len(elif.Body))
		for _, node := range elif.Body {
			elifBody = append(elifBody, simp.simplifyNode(node))
		}
		newLastIf := &ast2.If{
			Condition: &ast2.Unary{
				Operator: ast2.Not,
				X:        simp.simplifyExpression(elif.Condition),
			},
			Body: elifBody,
			Else: nil,
		}
		lastIf.Else = []ast2.Node{newLastIf}
		lastIf = newLastIf
	}
	lastIf.Else = make([]ast2.Node, 0, len(unless.Else))
	for _, node := range unless.Else {
		lastIf.Else = append(lastIf.Else, simp.simplifyNode(node))
	}
	return root
}
