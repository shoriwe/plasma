package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) Block(block *ast2.Block) []ast3.Node {
	body := make([]ast3.Node, 0, len(block.Body))
	for _, node := range block.Body {
		body = append(body, transform.Node(node)...)
	}
	return []ast3.Node{
		&ast3.Block{
			Body: body,
		},
	}
}
