package assembler

import (
	"github.com/shoriwe/plasma/pkg/ast3"
	"github.com/shoriwe/plasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/plasma/pkg/common"
)

func (a *assembler) Label(label *ast3.Label) []byte {
	result := []byte{opcodes.Label}
	result = append(result, common.IntToBytes(label.Code)...)
	return result
}
