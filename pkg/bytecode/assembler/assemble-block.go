package assembler

import (
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
)

func (a *assembler) Block(block *ast3.Block) []byte {
	body := make([]byte, 0, len(block.Body))
	body = append(body, opcodes.EnterBlock)
	for _, node := range block.Body {
		switch node.(type) {
		case *ast3.Yield, *ast3.Return:
			body = append(body, opcodes.ExitBlock)
		default:
			break
		}
		body = append(body, a.assemble(node)...)
	}
	return body
}
