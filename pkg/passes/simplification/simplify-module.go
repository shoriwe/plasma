package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func simplifyModule(module *ast.ModuleStatement) *ast2.Module {
	body := make([]ast2.Node, 0, len(module.Body))
	for _, node := range module.Body {
		body = append(body, simplifyNode(node))
	}
	return &ast2.Module{
		Name: simplifyIdentifier(module.Name),
		Body: body,
	}
}
