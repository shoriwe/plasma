package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func simplifyLambda(lambda *ast.LambdaExpression) *ast2.Lambda {
	arguments := make([]*ast2.Identifier, 0, len(lambda.Arguments))
	for _, argument := range lambda.Arguments {
		arguments = append(arguments, simplifyIdentifier(argument))
	}
	return &ast2.Lambda{
		Arguments: arguments,
		Result:    simplifyExpression(lambda.Code),
	}
}
