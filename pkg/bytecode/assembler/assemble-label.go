package assembler

import "github.com/shoriwe/gplasma/pkg/ast3"

func (a *assembler) Label(label *ast3.Label) []byte {
	lAccess, found := a.labels[label.Code]
	if !found {
		lAccess = &labelAccess{
			index:         0,
			jumpIndexes:   nil,
			ifJumpIndexes: nil,
		}
		a.labels[label.Code] = lAccess
	}
	lAccess.index = a.bytecodeIndex
	return nil
}
