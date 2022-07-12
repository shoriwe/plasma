package assembler

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/gplasma/pkg/common"
	"reflect"
)

type (
	assembler struct {
		bytecodeIndex int
		labels        map[int]int
		jumpsIndexes  map[int]int
	}
)

func newAssembler() *assembler {
	return &assembler{
		bytecodeIndex: 0,
		labels:        map[int]int{},
		jumpsIndexes:  map[int]int{},
	}
}

func (a *assembler) assemble(node ast3.Node) []byte {
	if node == nil {
		return nil
	}
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
		chunk := a.assemble(node)
		a.bytecodeIndex += len(chunk)
		bytecode = append(bytecode, chunk...)
	}
	finalBytecode := make([]byte, 0, len(bytecode))
	for index, operation := range bytecode {
		finalBytecode = append(finalBytecode, operation)
		if labelCode, found := a.jumpsIndexes[index]; found && operation == opcodes.Jump || operation == opcodes.IfJump {
			labelIndex := a.labels[labelCode]
			// fmt.Printf("(%d): %d - %d = %d\n", operation, labelIndex, index, labelIndex-index)
			finalBytecode = append(finalBytecode, common.IntToBytes(labelIndex-index)...)
		}
	}
	return finalBytecode
}
