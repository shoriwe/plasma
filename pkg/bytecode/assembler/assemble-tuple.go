package assembler

import (
	"github.com/shoriwe/plasma/pkg/ast3"
	"github.com/shoriwe/plasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/plasma/pkg/common"
)

func (a *assembler) Tuple(tuple *ast3.Tuple) []byte {
	var result []byte
	for _, value := range tuple.Values {
		result = append(result, a.Expression(value)...)
		result = append(result, opcodes.Push)
	}
	result = append(result, opcodes.NewTuple)
	result = append(result, common.IntToBytes(len(tuple.Values))...)
	return result
}
