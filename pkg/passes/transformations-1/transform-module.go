package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) Module(module *ast2.Module) []ast3.Node {
	body := make([]ast3.Node, 0, len(module.Body))
	for _, node := range module.Body {
		body = append(body, transform.Node(node)...)
	}
	moduleAssignment := &ast3.Assignment{
		Left: transform.Identifier(module.Name),
		Right: &ast3.Class{
			Body: body,
		},
	}
	return []ast3.Node{moduleAssignment}
}
