package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simp *simplify) simplifyInterface(i *ast.InterfaceStatement) *ast2.Class {
	bases := make([]ast2.Expression, 0, len(i.Bases))
	for _, base := range i.Bases {
		bases = append(bases, simp.simplifyExpression(base))
	}
	body := make([]ast2.Node, 0, len(i.MethodDefinitions))
	for _, methodDefinition := range i.MethodDefinitions {
		body = append(body, simp.simplifyFunction(methodDefinition))
	}
	return &ast2.Class{
		Name:  simp.simplifyIdentifier(i.Name),
		Bases: bases,
		Body:  body,
	}
}
