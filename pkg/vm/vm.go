package vm

import (
	"github.com/shoriwe/gruby/pkg/errors"
	"github.com/shoriwe/gruby/pkg/vm/runtime"
)

type Plasma struct {
	symbolTable      *runtime.SymbolTable
	symbolTableStack *runtime.Stack
	objectStack      *runtime.Stack
	code             []interface{}
	codeLength       int
	cursor           int
}

func (plasma *Plasma) basicBinaryOP(leftOPName string, rightOPName string) *errors.Error {
	plasma.cursor++
	leftHandSide := plasma.objectStack.Pop()
	rightHandSide := plasma.objectStack.Pop()
	result, callError := runtime.BasicBinaryOP(leftOPName, rightOPName, leftHandSide.(runtime.Object), rightHandSide.(runtime.Object))
	if callError != nil {
		return callError
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
	return plasma.basicBinaryOP(runtime.BitwiseLeftName, runtime.RightBitwiseLeftName)
}

func (plasma *Plasma) bitwiseRightOP() *errors.Error {
	return plasma.basicBinaryOP(runtime.BitwiseRightName, runtime.RightBitwiseRightName)
}

func (plasma *Plasma) bitwiseAndOP() *errors.Error {
	return plasma.basicBinaryOP(runtime.BitwiseAndName, runtime.RightBitwiseAndName)
}

func (plasma *Plasma) bitwiseOrOP() *errors.Error {
	return plasma.basicBinaryOP(runtime.BitwiseOrName, runtime.RightBitwiseOrName)
}

func (plasma *Plasma) bitwiseXorOP() *errors.Error {
	return plasma.basicBinaryOP(runtime.BitwiseXorName, runtime.RightBitwiseXorName)
}

func (plasma *Plasma) andOP() *errors.Error {
	return plasma.basicBinaryOP(runtime.AndName, runtime.RightAndName)
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
		case runtime.BitwiseLeftOP:
			instructionExecError = plasma.bitwiseLeftOP()
		case runtime.BitwiseRightOP:
			instructionExecError = plasma.bitwiseRightOP()
		case runtime.BitwiseAndOP:
			instructionExecError = plasma.bitwiseAndOP()
		case runtime.BitwiseOrOP:
			instructionExecError = plasma.bitwiseAndOP()
		case runtime.BitwiseXorOP:
			instructionExecError = plasma.bitwiseXorOP()
		case runtime.AndOP:
			instructionExecError = plasma.andOP()
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
