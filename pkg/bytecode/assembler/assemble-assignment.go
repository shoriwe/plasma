package assembler

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/gplasma/pkg/common"
	"reflect"
)

func (a *assembler) Assignment(assign *ast3.Assignment) []byte {
	result := a.Expression(assign.Right)
	result = append(result, opcodes.Push)
	switch left := assign.Left.(type) {
	case *ast3.Identifier:
		result = append(result, opcodes.IdentifierAssign)
		result = append(result, common.IntToBytes(len(left.Symbol))...)
		result = append(result, []byte(left.Symbol)...)
	case *ast3.Selector:
		result = append(result, a.Expression(left.X)...)
		result = append(result, opcodes.Push)
		result = append(result, opcodes.SelectorAssign)
		result = append(result, common.IntToBytes(len(left.Identifier.Symbol))...)
		result = append(result, []byte(left.Identifier.Symbol)...)
	case *ast3.Index:
		result = append(result, a.Expression(left.Source)...)
		result = append(result, opcodes.Push)
		result = append(result, a.Expression(left.Index)...)
		result = append(result, opcodes.Push)
		result = append(result, opcodes.IndexAssign)
	default:
		panic(fmt.Sprintf("unknown left hand side type %s", reflect.TypeOf(left).String()))
	}
	return result
}
