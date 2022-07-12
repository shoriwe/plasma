package assembler

import (
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
)

func (a *assembler) Require(require *ast3.Require) []byte {
	result := a.Expression(require.X)
	result = append(result, opcodes.Push, opcodes.Require)
	return result
}
