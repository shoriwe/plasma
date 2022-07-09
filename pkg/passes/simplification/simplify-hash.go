package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func simplifyHash(hash *ast.HashExpression) *ast2.Hash {
	values := make([]*ast2.KeyValue, 0, len(hash.Values))
	for _, value := range hash.Values {
		values = append(values, &ast2.KeyValue{
			Key:   simplifyExpression(value.Key),
			Value: simplifyExpression(value.Value),
		})
	}
	return &ast2.Hash{
		Values: values,
	}
}
