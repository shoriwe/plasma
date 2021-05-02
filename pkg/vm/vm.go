package vm

import (
	"github.com/shoriwe/gruby/pkg/errors"
	"github.com/shoriwe/gruby/pkg/vm/runtime"
	"reflect"
)

type Plasma struct {
	symbolTable      *runtime.SymbolTable
	symbolTableStack *runtime.Stack
	objectStack      *runtime.Stack
	code             []interface{}
	codeLength       int
	cursor           int
}

func (plasma *Plasma) basicBinaryOP(leftOP string, rightOP string) *errors.Error {
	plasma.cursor++
	leftHandSide := plasma.objectStack.Pop()
	rightHandSide := plasma.objectStack.Pop()
	operation, getError := runtime.GetAttribute(leftHandSide.(runtime.Object), leftOP, false)
	var result runtime.Object
	var opError *errors.Error
	isRight := false
	if getError != nil {
		var getError2 *errors.Error
		operation, getError2 = runtime.GetAttribute(rightHandSide.(runtime.Object), rightOP, false)
		if getError2 != nil {
			return getError
		}
		isRight = true
	}
	switch operation.(type) {
	case func(runtime.Object) (runtime.Object, *errors.Error):
		if isRight {
			result, opError = operation.(func(runtime.Object) (runtime.Object, *errors.Error))(leftHandSide.(runtime.Object))
		} else {
			result, opError = operation.(func(runtime.Object) (runtime.Object, *errors.Error))(rightHandSide.(runtime.Object))
		}
	case *runtime.Function:
		if isRight {
			result, opError = operation.(*runtime.Function).Call(rightHandSide.(runtime.Object))
		} else {
			result, opError = operation.(*runtime.Function).Call(leftHandSide.(runtime.Object))
		}
	default:
		return runtime.NewTypeError(runtime.FunctionTypeString, reflect.TypeOf(operation).String())
	}
	if opError != nil {
		return opError
	}
	plasma.objectStack.Push(result)
	return nil
}

func (plasma *Plasma) additionOP() *errors.Error {
	return plasma.basicBinaryOP(runtime.AdditionName, runtime.RightAdditionName)
}

func (plasma *Plasma) subtractionOP() *errors.Error {
	return plasma.basicBinaryOP(runtime.SubtractionName, runtime.RightSubtractionName)
}

func (plasma *Plasma) divisionOP() *errors.Error {
	return plasma.basicBinaryOP(runtime.DivisionName, runtime.RightDivisionName)
}

func (plasma *Plasma) modulusOP() *errors.Error {
	return plasma.basicBinaryOP(runtime.ModulusName, runtime.RightModulusName)
}

func (plasma *Plasma) multiplicationOP() *errors.Error {
	return plasma.basicBinaryOP(runtime.MultiplicationName, runtime.RightMultiplicationName)
}

func (plasma *Plasma) powerOfOP() *errors.Error {
	return plasma.basicBinaryOP(runtime.PowerOfName, runtime.RightPowerOfName)
}

func (plasma *Plasma) floorDivisionOP() *errors.Error {
	return plasma.basicBinaryOP(runtime.FloorDivisionName, runtime.RightFloorDivisionName)
}

func (plasma *Plasma) bitwiseLeftOP() *errors.Error {
	return plasma.basicBinaryOP(runtime.FloorDivisionName, runtime.RightFloorDivisionName)
}

func (plasma *Plasma) bitwiseRightOP() *errors.Error {
	return plasma.basicBinaryOP(runtime.FloorDivisionName, runtime.RightFloorDivisionName)
}

func (plasma *Plasma) bitwiseAndOP() *errors.Error {
	return plasma.basicBinaryOP(runtime.FloorDivisionName, runtime.RightFloorDivisionName)
}

func (plasma *Plasma) bitwiseOrOP() *errors.Error {
	return plasma.basicBinaryOP(runtime.FloorDivisionName, runtime.RightFloorDivisionName)
}

func (plasma *Plasma) bitwiseXorOP() *errors.Error {
	return plasma.basicBinaryOP(runtime.FloorDivisionName, runtime.RightFloorDivisionName)
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

func (plasma *Plasma) Execute() (runtime.Object, *errors.Error) {
	for ; plasma.cursor < plasma.codeLength; {
		if _, ok := plasma.code[plasma.cursor].(uint); !ok {
			return nil, runtime.NewRuntimeError(runtime.UnknownOP, "Unknown instruction type")
		}
		var instructionExecError *errors.Error
		switch plasma.code[plasma.cursor].(uint) {
		case runtime.AddOP:
			instructionExecError = plasma.additionOP()
		case runtime.SubOP:
			instructionExecError = plasma.subtractionOP()
		case runtime.DivOP:
			instructionExecError = plasma.divisionOP()
		case runtime.ModOP:
			instructionExecError = plasma.modulusOP()
		case runtime.MulOP:
			instructionExecError = plasma.multiplicationOP()
		case runtime.PowOP:
			instructionExecError = plasma.powerOfOP()
		case runtime.FloorDivOP:
			instructionExecError = plasma.floorDivisionOP()
		case runtime.PushOP:
			instructionExecError = plasma.pushOP()
		case runtime.BitwiseLeft:
			instructionExecError = plasma.bitwiseLeftOP()
		case runtime.BitwiseRight:
			instructionExecError = plasma.bitwiseRightOP()
		case runtime.BitwiseAnd:
			instructionExecError = plasma.bitwiseAndOP()
		case runtime.BitwiseOr:
			instructionExecError = plasma.bitwiseAndOP()
		case runtime.BitwiseXor:
			instructionExecError = plasma.bitwiseXorOP()
		case runtime.ReturnOP:
			plasma.cursor++
			if plasma.objectStack.IsEmpty() {
				return runtime.NewNone(), nil
			}
			return plasma.objectStack.Pop().(runtime.Object), nil
		default:
			instructionExecError = errors.New(runtime.UnknownLine, "Unknown Operation", runtime.UnknownOP)
		}
		if instructionExecError != nil {
			return nil, instructionExecError
		}
	}
	return runtime.NewNone(), nil
}

func NewPlasmaVM(symbolTable *runtime.SymbolTable) *Plasma {
	if symbolTable == nil {
		symbolTable = runtime.NewSymbolTable(nil)
	}
	return &Plasma{
		symbolTable:      symbolTable,
		symbolTableStack: runtime.NewStack(),
		objectStack:      runtime.NewStack(),
		cursor:           0,
	}
}
