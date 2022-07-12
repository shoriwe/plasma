package assembler

import (
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
)

func (a *assembler) Jump(jump *ast3.Jump) []byte {
	_, found := a.labels[jump.Target.Code]
	if !found {
		a.labels[jump.Target.Code] = -1
	}
	a.jumpsIndexes[a.bytecodeIndex] = jump.Target.Code
	return []byte{opcodes.Jump}
}

func (a *assembler) ContinueJump(jump *ast3.ContinueJump) []byte {
	_, found := a.labels[jump.Target.Code]
	if !found {
		a.labels[jump.Target.Code] = -1
	}
	a.jumpsIndexes[a.bytecodeIndex] = jump.Target.Code
	return []byte{opcodes.Jump}
}

func (a *assembler) BreakJump(jump *ast3.BreakJump) []byte {
	_, found := a.labels[jump.Target.Code]
	if !found {
		a.labels[jump.Target.Code] = -1
	}
	a.jumpsIndexes[a.bytecodeIndex] = jump.Target.Code
	return []byte{opcodes.Jump}
}

func (a *assembler) IfJump(jump *ast3.IfJump) []byte {
	index := a.bytecodeIndex
	condition := a.Expression(jump.Condition)
	index += 1 + len(condition)

	_, found := a.labels[jump.Target.Code]
	if !found {
		a.labels[jump.Target.Code] = -1
	}
	a.jumpsIndexes[index] = jump.Target.Code

	result := append(condition, opcodes.Push, opcodes.IfJump)
	return result
}
