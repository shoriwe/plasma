package simplification

import "github.com/shoriwe/gplasma/pkg/ast"

func simplifyGenerator(generator *ast.GeneratorExpression) *ast.GeneratorExpression {
	newReceivers := make([]*ast.Identifier, 0, len(generator.Receivers))
	for _, receiver := range generator.Receivers {
		newReceivers = append(newReceivers, SimplifyExpression(receiver).(*ast.Identifier))
	}
	return &ast.GeneratorExpression{
		Operation: SimplifyExpression(generator.Operation),
		Receivers: newReceivers,
		Source:    SimplifyExpression(generator.Source),
	}
}
