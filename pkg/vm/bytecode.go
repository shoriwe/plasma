package vm

import (
	"math/bits"
)

type Instruction struct {
	OpCode uint8
	Line   int
}

func NewInstruction(opCode uint8) Instruction {
	return Instruction{
		OpCode: opCode,
	}
}

type Code struct {
	Instruction Instruction
	Value       interface{}
	Line        int
}

func NewCode(opCode uint8, line int, value interface{}) Code {
	return Code{
		Instruction: NewInstruction(opCode),
		Value:       value,
		Line:        line,
	}
}

type bytecodeNode struct {
	code Code
	next *bytecodeNode
}

type Bytecode struct {
	currentInstruction *bytecodeNode
	length             uint
}

func (bytecode *Bytecode) HasNext() bool {
	return bytecode.length != 0
}

func (bytecode *Bytecode) Next() Code {
	result := bytecode.currentInstruction
	bytecode.currentInstruction = bytecode.currentInstruction.next
	bytecode.length--
	return result.code
}

func (bytecode *Bytecode) Push(code Code) {
	if bytecode.length == bits.UintSize {
		panic("Bytecode reached its maximum capacity which is platform uint")
	}
	bytecode.currentInstruction = &bytecodeNode{
		code: code,
		next: bytecode.currentInstruction,
	}
	bytecode.length++
}

func NewBytecodeFromArray(codes []Code) *Bytecode {
	result := &Bytecode{
		currentInstruction: nil,
		length:             0,
	}
	codesLength := len(codes)
	for i := codesLength - 1; i > -1; i-- {
		result.Push(codes[i])
	}
	return result
}
