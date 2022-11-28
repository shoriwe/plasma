package transformations_1

import (
	"github.com/shoriwe/plasma/pkg/ast2"
	"github.com/shoriwe/plasma/pkg/ast3"
)

func (transform *transformPass) Defer(def *ast2.Defer) []ast3.Node {
	return []ast3.Node{
		&ast3.Defer{
			X: transform.Expression(def.X),
		},
	}
}
