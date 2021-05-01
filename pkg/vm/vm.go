package vm

import (
	"github.com/shoriwe/gruby/pkg/errors"
	"github.com/shoriwe/gruby/pkg/vm/object"
	"github.com/shoriwe/gruby/pkg/vm/utils"
	vm_errors "github.com/shoriwe/gruby/pkg/vm/vm-errors"
)

type Plasma struct {
	symbolTable *utils.SymbolTable
	stack       *utils.Stack
	code        []interface{}
	cursor      int
}

func (plasma *Plasma) additionOP() *errors.Error {
	plasma.cursor++
	leftHandSide := plasma.stack.Pop()
	rightHandSide := plasma.stack.Pop()
	result, opError := leftHandSide.(object.Object).Addition(rightHandSide.(object.Object))
	if opError != nil {
		return opError
	}
	plasma.stack.Push(result)
	return nil
}

func (plasma *Plasma) pushOP() *errors.Error {
	plasma.cursor++
	plasma.stack.Push(plasma.code[plasma.cursor])
	plasma.cursor++
	return nil
}

func (plasma *Plasma) executeOP() *errors.Error {
	switch plasma.code[plasma.cursor].(uint) {
	case AddOP:
		return plasma.additionOP()
	case PushOP:
		return plasma.pushOP()
	}
	return errors.New(vm_errors.UnknownLine, "Unknown Operation", vm_errors.UnknownOP)
}

func (plasma *Plasma) Execute(code []interface{}) (object.Object, *errors.Error) {
	plasma.code = code
	codeLength := len(code)
	for ; plasma.cursor < codeLength; {
		switch plasma.code[plasma.cursor].(type) {
		case uint:
			opExecutionError := plasma.executeOP()
			if opExecutionError != nil {
				return nil, opExecutionError
			}
		default:
			return nil, vm_errors.NewRuntimeError(vm_errors.UnknownOP, "Unknown instruction type")
		}
	}
	if !plasma.stack.IsEmpty() {
		return plasma.stack.Pop().(object.Object), nil
	}
	return nil, nil
}

func NewPlasmaVM(symbolTable *utils.SymbolTable) *Plasma {
	if symbolTable == nil {
		symbolTable = utils.NewSymbolTable()
	}
	return &Plasma{
		symbolTable: symbolTable,
		stack:       utils.NewStack(),
		cursor:      0,
	}
}
