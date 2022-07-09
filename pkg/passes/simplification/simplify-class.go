package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func simplifyClass(class *ast.ClassStatement) *ast2.Class {
	bases := make([]ast2.Expression, 0, len(class.Bases))
	for _, base := range class.Bases {
		bases = append(bases, simplifyExpression(base))
	}
	body := make([]ast2.Node, 0, len(class.Body))
	for _, node := range class.Body {
		body = append(body, simplifyNode(node))
	}
	return &ast2.Class{
		Name:  simplifyIdentifier(class.Name),
		Bases: bases,
		Body:  body,
	}
}
