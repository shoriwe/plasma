package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) Require(require *ast2.Require) []ast3.Node {
	return []ast3.Node{
		&ast3.Require{
			X: transform.Expression(require.X),
		},
	}
}
