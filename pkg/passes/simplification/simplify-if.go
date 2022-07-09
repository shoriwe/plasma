package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func simplifyIf(if_ *ast.IfStatement) *ast2.If {
	body := make([]ast2.Node, 0, len(if_.Body))
	for _, node := range if_.Body {
		body = append(body, simplifyNode(node))
	}
	root := &ast2.If{
		Condition: simplifyExpression(if_.Condition),
		Body:      body,
		Else:      nil,
	}
	lastIf := root
	for _, elif := range if_.ElifBlocks {
		elifBody := make([]ast2.Node, 0, len(elif.Body))
		for _, node := range elif.Body {
			elifBody = append(elifBody, simplifyNode(node))
		}
		newLastIf := &ast2.If{
			Condition: simplifyExpression(elif.Condition),
			Body:      elifBody,
			Else:      nil,
		}
		lastIf.Else = []ast2.Node{newLastIf}
		lastIf = newLastIf
	}
	lastIf.Else = make([]ast2.Node, 0, len(if_.Else))
	for _, node := range if_.Else {
		lastIf.Else = append(lastIf.Else, simplifyNode(node))
	}
	return root
}
