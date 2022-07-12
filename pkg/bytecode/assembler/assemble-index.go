package assembler

import (
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
)

func (a *assembler) Index(index *ast3.Index) []byte {
	var result []byte
	result = append(result, a.Expression(index.Index)...)
	result = append(result, opcodes.Push)
	result = append(result, a.Expression(index.Source)...)
	result = append(result, opcodes.Push)
	result = append(result, opcodes.Index)
	return result
}
