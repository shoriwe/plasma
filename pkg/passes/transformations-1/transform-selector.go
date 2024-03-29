package transformations_1

import (
	"github.com/shoriwe/plasma/pkg/ast2"
	"github.com/shoriwe/plasma/pkg/ast3"
)

func (transform *transformPass) Selector(selector *ast2.Selector) *ast3.Selector {
	return &ast3.Selector{
		X:          transform.Expression(selector.X),
		Identifier: transform.Identifier(selector.Identifier),
	}
}
