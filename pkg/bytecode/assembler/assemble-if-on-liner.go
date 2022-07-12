package assembler

import (
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/gplasma/pkg/common"
)

func (a *assembler) IfOneLiner(ifOneLiner *ast3.IfOneLiner) []byte {
	if_ := a.Expression(ifOneLiner.Result)
	condition := a.Expression(ifOneLiner.Condition)
	else_ := a.Expression(ifOneLiner.Else)
	var result []byte
	result = append(result, condition...)
	result = append(result, opcodes.Push)
	result = append(result, opcodes.IfOneLiner)
	result = append(result, common.IntToBytes(len(if_))...)
	result = append(result, common.IntToBytes(len(else_))...)
	result = append(result, if_...)
	result = append(result, else_...)
	return result
}
