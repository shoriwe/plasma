package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
)

func simplifyMethodInvocation(mi *ast.MethodInvocationExpression) *ast.MethodInvocationExpression {
	newArguments := make([]ast.IExpression, 0, len(mi.Arguments))
	for _, argument := range mi.Arguments {
		newArguments = append(newArguments, SimplifyExpression(argument))
	}
	return &ast.MethodInvocationExpression{
		Function:  SimplifyExpression(mi.Function),
		Arguments: newArguments,
	}
}
