package assembler

import (
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
)

func (a *assembler) Jump(jump *ast3.Jump) []byte {
	lAccess, found := a.labels[jump.Target.Code]
	if !found {
		lAccess = &labelAccess{
			index:         0,
			jumpIndexes:   nil,
			ifJumpIndexes: nil,
		}
		a.labels[jump.Target.Code] = lAccess
	}
	lAccess.jumpIndexes = append(lAccess.jumpIndexes, a.bytecodeIndex)
	return []byte{opcodes.Jump, 0, 0, 0, 0, 0, 0, 0, 0}
}

func (a *assembler) ContinueJump(jump *ast3.ContinueJump) []byte {
	lAccess, found := a.labels[jump.Target.Code]
	if !found {
		lAccess = &labelAccess{
			index:         0,
			jumpIndexes:   nil,
			ifJumpIndexes: nil,
		}
		a.labels[jump.Target.Code] = lAccess
	}
	lAccess.jumpIndexes = append(lAccess.jumpIndexes, a.bytecodeIndex)
	return []byte{opcodes.Jump, 0, 0, 0, 0, 0, 0, 0, 0}
}

func (a *assembler) BreakJump(jump *ast3.BreakJump) []byte {
	lAccess, found := a.labels[jump.Target.Code]
	if !found {
		lAccess = &labelAccess{
			index:         0,
			jumpIndexes:   nil,
			ifJumpIndexes: nil,
		}
		a.labels[jump.Target.Code] = lAccess
	}
	lAccess.jumpIndexes = append(lAccess.jumpIndexes, a.bytecodeIndex)
	return []byte{opcodes.Jump, 0, 0, 0, 0, 0, 0, 0, 0}
}

func (a *assembler) IfJump(jump *ast3.IfJump) []byte {
	index := a.bytecodeIndex
	condition := a.Expression(jump.Condition)
	index += 1 + len(condition)

	lAccess, found := a.labels[jump.Target.Code]
	if !found {
		lAccess = &labelAccess{
			index:         0,
			jumpIndexes:   nil,
			ifJumpIndexes: nil,
		}
		a.labels[jump.Target.Code] = lAccess
	}
	lAccess.ifJumpIndexes = append(lAccess.jumpIndexes, index)

	result := append(condition, opcodes.Push, opcodes.IfJump, 0, 0, 0, 0, 0, 0, 0, 0)
	return result
}
