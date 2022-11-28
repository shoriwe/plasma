package assembler

import (
	"github.com/shoriwe/plasma/pkg/ast3"
	"github.com/shoriwe/plasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/plasma/pkg/common"
)

func (a *assembler) Array(array *ast3.Array) []byte {
	var result []byte
	for _, value := range array.Values {
		result = append(result, a.Expression(value)...)
		result = append(result, opcodes.Push)
	}
	result = append(result, opcodes.NewArray)
	result = append(result, common.IntToBytes(len(array.Values))...)
	return result
}
