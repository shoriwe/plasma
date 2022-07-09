package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func simplifyDoWhile(do *ast.DoWhileStatement) *ast2.DoWhile {
	var (
		body                      = make([]ast2.Node, 0, len(do.Body))
		condition ast2.Expression = simplifyExpression(do.Condition)
	)
	for _, node := range do.Body {
		body = append(body, simplifyNode(node))
	}
	return &ast2.DoWhile{
		Body:      body,
		Condition: condition,
	}
}
