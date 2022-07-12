package assembler

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast3"
	"reflect"
)

type (
	labelAccess struct {
		index         int
		jumpIndexes   []int
		ifJumpIndexes []int
	}
	assembler struct {
		bytecodeIndex int
		labels        map[int]*labelAccess
	}
)

func newAssembler() *assembler {
	return &assembler{
		bytecodeIndex: 0,
		labels:        map[int]*labelAccess{},
	}
}

func (a *assembler) assemble(node ast3.Node) []byte {
	switch n := node.(type) {
	case ast3.Statement:
		return a.Statement(n)
	case ast3.Expression:
		return a.Expression(n)
	default:
		panic(fmt.Sprintf("unknown node type %s", reflect.TypeOf(n).String()))
	}
}

func Assemble(program ast3.Program) []byte {
	bytecode := make([]byte, 0, len(program))
	a := newAssembler()
	for _, node := range program {
		chunk := append(bytecode, a.assemble(node)...)
		a.bytecodeIndex += len(chunk)
		bytecode = append(bytecode, chunk...)
	}
	return bytecode
}
