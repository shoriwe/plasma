package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simp *simplify) simplifyRequire(require *ast.RequireStatement) *ast2.Require {
	return &ast2.Require{
		X: simp.simplifyExpression(require.X),
	}
}
