package simplification

import "github.com/shoriwe/gplasma/pkg/ast"

func simplifyHash(hash *ast.HashExpression) *ast.HashExpression {
	newContents := make([]*ast.KeyValue, 0, len(hash.Values))
	for _, keyValue := range hash.Values {
		newKey := SimplifyExpression(keyValue.Key)
		newValue := SimplifyExpression(keyValue.Value)
		newContents = append(newContents,
			&ast.KeyValue{
				Key:   newKey,
				Value: newValue,
			},
		)
	}
	return &ast.HashExpression{
		Values: newContents,
	}
}
