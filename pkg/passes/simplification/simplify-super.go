package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simp *simplify) simplifySuper(super *ast.SuperExpression) *ast2.Super {
	return &ast2.Super{
		X: simp.simplifyExpression(super.X),
	}
}
