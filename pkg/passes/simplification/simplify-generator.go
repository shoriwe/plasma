package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simp *simplify) simplifyGeneratorExpr(generator *ast.GeneratorExpression) *ast2.Generator {
	receivers := make([]*ast2.Identifier, 0, len(generator.Receivers))
	for _, receiver := range generator.Receivers {
		receivers = append(receivers, simp.simplifyIdentifier(receiver))
	}
	return &ast2.Generator{
		Operation: simp.simplifyExpression(generator.Operation),
		Receivers: receivers,
		Source:    simp.simplifyExpression(generator.Source),
	}
}

func (simp *simplify) simplifyGeneratorDef(generator *ast.GeneratorDefinitionStatement) *ast2.GeneratorDefinition {
	arguments := make([]*ast2.Identifier, 0, len(generator.Arguments))
	for _, argument := range generator.Arguments {
		arguments = append(arguments, simp.simplifyIdentifier(argument))
	}
	body := make([]ast2.Node, 0, len(generator.Body))
	for _, node := range generator.Body {
		body = append(body, simp.simplifyNode(node))
	}
	return &ast2.GeneratorDefinition{
		Name:      simp.simplifyIdentifier(generator.Name),
		Arguments: arguments,
		Body:      body,
	}
}
