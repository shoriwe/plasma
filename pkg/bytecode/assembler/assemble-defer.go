package assembler

import (
	"github.com/shoriwe/plasma/pkg/ast3"
	"github.com/shoriwe/plasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/plasma/pkg/common"
)

func (a *assembler) Defer(defer_ *ast3.Defer) []byte {
	expression := a.Expression(defer_.X)
	result := []byte{opcodes.Defer}
	result = append(result, common.IntToBytes(len(expression))...)
	result = append(result, expression...)
	return result
}
