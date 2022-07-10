package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simp *simplify) simplifyFunction(f *ast.FunctionDefinitionStatement) *ast2.FunctionDefinition {
	arguments := make([]*ast2.Identifier, 0, len(f.Arguments))
	for _, argument := range f.Arguments {
		arguments = append(arguments, simp.simplifyIdentifier(argument))
	}
	body := make([]ast2.Node, 0, len(f.Body))
	for _, node := range f.Body {
		body = append(body, simp.simplifyNode(node))
	}
	return &ast2.FunctionDefinition{
		Name:      simp.simplifyIdentifier(f.Name),
		Arguments: arguments,
		Body:      body,
	}
}
