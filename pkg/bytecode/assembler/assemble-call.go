package assembler

import (
	"github.com/shoriwe/plasma/pkg/ast3"
	"github.com/shoriwe/plasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/plasma/pkg/common"
)

func (a *assembler) Call(call *ast3.Call) []byte {
	var result []byte
	for _, argument := range call.Arguments {
		result = append(result, a.Expression(argument)...)
		result = append(result, opcodes.Push)
	}
	result = append(result, a.Expression(call.Function)...)
	result = append(result, opcodes.Push)
	result = append(result, opcodes.Call)
	result = append(result, common.IntToBytes(len(call.Arguments))...)
	return result
}
