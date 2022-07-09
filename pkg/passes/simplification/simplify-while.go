package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func simplifyWhile(while *ast.WhileLoopStatement) *ast2.While {
	var (
		body                      = make([]ast2.Node, 0, len(while.Body))
		condition ast2.Expression = simplifyExpression(while.Condition)
	)
	for _, node := range while.Body {
		body = append(body, simplifyNode(node))
	}
	return &ast2.While{
		Body:      body,
		Condition: condition,
	}
}
