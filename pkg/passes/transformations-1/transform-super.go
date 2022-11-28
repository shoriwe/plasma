package transformations_1

import (
	"github.com/shoriwe/plasma/pkg/ast2"
	"github.com/shoriwe/plasma/pkg/ast3"
)

func (transform *transformPass) Super(super *ast2.Super) *ast3.Super {
	return &ast3.Super{
		X: transform.Expression(super.X),
	}
}
