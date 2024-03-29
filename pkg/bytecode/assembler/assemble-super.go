package assembler

import (
	"github.com/shoriwe/plasma/pkg/ast3"
	"github.com/shoriwe/plasma/pkg/bytecode/opcodes"
)

func (a *assembler) Super(super *ast3.Super) []byte {
	var result []byte
	result = append(result, a.Expression(super.X)...)
	result = append(result, opcodes.Push)
	result = append(result, opcodes.Super)
	return result
}
