package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simp *simplify) simplifyDefer(d *ast.DeferStatement) *ast2.Defer {
	return &ast2.Defer{
		X: simp.simplifyExpression(d.X),
	}
}
