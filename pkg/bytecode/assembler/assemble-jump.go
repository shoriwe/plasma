package assembler

import (
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
)

func (a *assembler) Jump(jump *ast3.Jump) []byte {
	return []byte{opcodes.Jump, 0, 0, 0, 0, 0, 0, 0, 0}
}

func (a *assembler) ContinueJump(jump *ast3.ContinueJump) []byte {
	return []byte{opcodes.Jump, 0, 0, 0, 0, 0, 0, 0, 0}
}

func (a *assembler) BreakJump(jump *ast3.BreakJump) []byte {
	return []byte{opcodes.Jump, 0, 0, 0, 0, 0, 0, 0, 0}
}

func (a *assembler) IfJump(jump *ast3.IfJump) []byte {
	result := a.Expression(jump.Condition)
	result = append(result, opcodes.Push, opcodes.IfJump, 0, 0, 0, 0, 0, 0, 0, 0)
	return result
}
