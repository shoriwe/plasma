package vm

import (
	"github.com/shoriwe/gruby/pkg/errors"
)

type Plasma struct {
	symbolTable      *SymbolTable
	symbolTableStack *Stack
	objectStack      *Stack
	code             []interface{}
	cursor           int
}

func (plasma *Plasma) additionOP() *errors.Error {
	plasma.cursor++
	leftHandSide := plasma.objectStack.Pop()
	rightHandSide := plasma.objectStack.Pop()
	operation, getError := GetAttribute(leftHandSide.(Object), Addition)
	var result Object
	var opError *errors.Error
	if getError != nil {
		var getError2 *errors.Error
		operation, getError2 = GetAttribute(rightHandSide.(Object), RightAddition)
		if getError2 != nil {
			return getError
		}
		result, opError = operation.(func(Object) (Object, *errors.Error))(leftHandSide.(Object))
	} else {
		result, opError = operation.(func(Object) (Object, *errors.Error))(rightHandSide.(Object))
	}
	if opError != nil {
		return opError
	}
	plasma.objectStack.Push(result)
	return nil
}

func (plasma *Plasma) pushOP() *errors.Error {
	plasma.cursor++
	plasma.objectStack.Push(plasma.code[plasma.cursor])
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
	return errors.New(UnknownLine, "Unknown Operation", UnknownOP)
}

func (plasma *Plasma) Execute(code []interface{}) (Object, *errors.Error) {
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
			return nil, NewRuntimeError(UnknownOP, "Unknown instruction type")
		}
	}
	if !plasma.objectStack.IsEmpty() {
		return plasma.objectStack.Pop().(Object), nil
	}
	return nil, nil
}

func NewPlasmaVM(symbolTable *SymbolTable) *Plasma {
	if symbolTable == nil {
		symbolTable = NewSymbolTable(nil)
	}
	return &Plasma{
		symbolTable:      symbolTable,
		symbolTableStack: NewStack(),
		objectStack:      NewStack(),
		cursor:           0,
	}
}
