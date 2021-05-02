package vm

import (
	"github.com/shoriwe/gruby/pkg/errors"
	"reflect"
)

type Plasma struct {
	symbolTable      *SymbolTable
	symbolTableStack *Stack
	objectStack      *Stack
	code             []interface{}
	codeLength       int
	cursor           int
}

func (plasma *Plasma) basicBinaryOP(leftOP string, rightOP string) *errors.Error {
	plasma.cursor++
	leftHandSide := plasma.objectStack.Pop()
	rightHandSide := plasma.objectStack.Pop()
	operation, getError := GetAttribute(leftHandSide.(Object), leftOP, false)
	var result Object
	var opError *errors.Error
	isRight := false
	if getError != nil {
		var getError2 *errors.Error
		operation, getError2 = GetAttribute(rightHandSide.(Object), rightOP, false)
		if getError2 != nil {
			return getError
		}
		isRight = true
	}
	switch operation.(type) {
	case func(Object) (Object, *errors.Error):
		if isRight {
			result, opError = operation.(func(Object) (Object, *errors.Error))(leftHandSide.(Object))
		} else {
			result, opError = operation.(func(Object) (Object, *errors.Error))(rightHandSide.(Object))
		}
	case *Function:
		if isRight {
			result, opError = operation.(*Function).Call(rightHandSide.(Object))
		} else {
			result, opError = operation.(*Function).Call(leftHandSide.(Object))
		}
	default:
		return NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
	}
	if opError != nil {
		return opError
	}
	plasma.objectStack.Push(result)
	return nil
}

func (plasma *Plasma) additionOP() *errors.Error {
	return plasma.basicBinaryOP(Addition, RightAddition)
}

func (plasma *Plasma) subtractionOP() *errors.Error {
	return plasma.basicBinaryOP(Subtraction, RightSubtraction)
}

func (plasma *Plasma) divisionOP() *errors.Error {
	return plasma.basicBinaryOP(Division, RightDivision)
}

func (plasma *Plasma) pushOP() *errors.Error {
	plasma.cursor++
	plasma.objectStack.Push(plasma.code[plasma.cursor])
	plasma.cursor++
	return nil
}

func (plasma *Plasma) Initialize(code []interface{}) {
	plasma.code = code
	plasma.codeLength = len(code)
}

func (plasma *Plasma) Execute() (Object, *errors.Error) {
	for ; plasma.cursor < plasma.codeLength; {
		if _, ok := plasma.code[plasma.cursor].(uint); !ok {
			return nil, NewRuntimeError(UnknownOP, "Unknown instruction type")
		}
		var instructionExecError *errors.Error
		switch plasma.code[plasma.cursor].(uint) {
		case AddOP:
			instructionExecError = plasma.additionOP()
		case SubOP:
			instructionExecError = plasma.subtractionOP()
		case DivOP:
			instructionExecError = plasma.divisionOP()
		case PushOP:
			instructionExecError = plasma.pushOP()
		case ReturnOP:
			plasma.cursor++
			if plasma.objectStack.IsEmpty() {
				return NewNone(), nil
			}
			return plasma.objectStack.Pop().(Object), nil
		default:
			instructionExecError = errors.New(UnknownLine, "Unknown Operation", UnknownOP)
		}
		if instructionExecError != nil {
			return nil, instructionExecError
		}
	}
	return NewNone(), nil
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
