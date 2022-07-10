package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simp *simplify) simplifyModule(module *ast.ModuleStatement) *ast2.Module {
	body := make([]ast2.Node, 0, len(module.Body))
	for _, node := range module.Body {
		body = append(body, simp.simplifyNode(node))
	}
	return &ast2.Module{
		Name: simp.simplifyIdentifier(module.Name),
		Body: body,
	}
}
