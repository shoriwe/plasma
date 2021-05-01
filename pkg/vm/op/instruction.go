package op

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/errors"
	"github.com/shoriwe/gruby/pkg/vm/utils"
	"github.com/shoriwe/gruby/pkg/vm/vm-errors"
)

const (
	InstructionNotFound = "InstructionNotFound"
)

type InstructionSet struct {
	instructions map[uint]InstructionOP
}

func (instructionSet *InstructionSet) Execute(op uint, stack *utils.Stack) *errors.Error {
	instruction, ok := instructionSet.instructions[op]
	if !ok {
		return vm_errors.NewRuntimeError(InstructionNotFound, fmt.Sprintf("Instruction %d not found", op))
	}
	return instruction(stack)
}

func (instructionSet *InstructionSet) Merge(otherInstructionSet *InstructionSet) *InstructionSet {
	result := map[uint]InstructionOP{}
	for op, instruction := range otherInstructionSet.instructions {
		result[op] = instruction
	}
	for op, instruction := range instructionSet.instructions {
		result[op] = instruction
	}
	return &InstructionSet{instructions: result}
}
