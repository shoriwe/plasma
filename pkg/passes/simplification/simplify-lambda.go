package simplification

import "github.com/shoriwe/gplasma/pkg/ast"

func simplifyLambda(lambda *ast.LambdaExpression) *ast.LambdaExpression {
	newArguments := make([]*ast.Identifier, 0, len(lambda.Arguments))
	for _, argument := range lambda.Arguments {
		newArguments = append(newArguments, SimplifyExpression(argument).(*ast.Identifier))
	}
	return &ast.LambdaExpression{
		Arguments: newArguments,
		Code:      SimplifyExpression(lambda.Code),
	}
}
