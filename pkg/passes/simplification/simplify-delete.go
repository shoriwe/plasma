package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simp *simplify) simplifyDelete(del *ast.DeleteStatement) *ast2.Delete {
	var x ast2.Assignable
	switch dx := del.X.(type) {
	case *ast.Identifier:
		x = simp.simplifyIdentifier(dx)
	case *ast.IndexExpression:
		x = simp.simplifyIndex(dx)
	case *ast.SelectorExpression:
		x = simp.simplifySelector(dx)
	default:
		panic("unknown selector type")
	}
	return &ast2.Delete{
		X: x,
	}
}
