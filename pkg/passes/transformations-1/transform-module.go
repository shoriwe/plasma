package transformations_1

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) Module(module *ast2.Module) []ast3.Node {
	body := make([]ast3.Node, 0, len(module.Body))
	for _, node := range module.Body {
		body = append(body, transform.Node(node)...)
	}
	class := &ast3.Identifier{
		Symbol: fmt.Sprintf("____module_%s", module.Name.Symbol),
	}
	tempClassAssignment := &ast3.Assignment{
		Left: class,
		Right: &ast3.Class{
			Body: body,
		},
	}
	moduleAssignment := &ast3.Assignment{
		Left: transform.Identifier(module.Name),
		Right: &ast3.Call{
			Function: class,
		},
	}
	return []ast3.Node{tempClassAssignment, moduleAssignment}
}
