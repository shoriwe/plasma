package simplification

import (
	"github.com/shoriwe/plasma/pkg/ast"
	"github.com/shoriwe/plasma/pkg/ast2"
)

func (simplify *simplifyPass) Index(index *ast.IndexExpression) *ast2.Index {
	return &ast2.Index{
		Source: simplify.Expression(index.Source),
		Index:  simplify.Expression(index.Index),
	}
}
