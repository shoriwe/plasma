package vm

import (
	"fmt"
	"github.com/fatih/color"
	"strconv"
)

func (p *Plasma) Execute(context *Context, bytecode *Bytecode) (*Value, bool) {
	if context == nil {
		context = p.NewContext()
	}
	var executionError *Value
	for bytecode.HasNext() {
		code := bytecode.Next()

		if code.Line != 0 {
			fmt.Println(color.GreenString(strconv.Itoa(code.Line)), instructionNames[code.Instruction.OpCode], code.Value)
		} else {
			fmt.Println(color.RedString("UL"), instructionNames[code.Instruction.OpCode], code.Value)
		}

		switch code.Instruction.OpCode {
		case GetFalseOP:
			context.LastObject = p.GetFalse()
		case GetTrueOP:
			context.LastObject = p.GetTrue()
		case GetNoneOP:
			context.LastObject = p.GetNone()
		case NewStringOP:
			executionError = p.newStringOP(context, code.Value.(string))
		case NewBytesOP:
			executionError = p.newBytesOP(context, code.Value.([]uint8))
		case NewIntegerOP:
			executionError = p.newIntegerOP(context, code.Value.(int64))
		case NewFloatOP:
			executionError = p.newFloatOP(context, code.Value.(float64))
		case NewArrayOP:
			executionError = p.newArrayOP(context, code.Value.(int))
		case NewTupleOP:
			executionError = p.newTupleOP(context, code.Value.(int))
		case NewHashOP:
			executionError = p.newHashTableOP(context, code.Value.(int))
		case UnaryOP:
			executionError = p.unaryOP(context, code.Value.(uint8))
		case BinaryOP:
			executionError = p.binaryOP(context, code.Value.(uint8))
		case MethodInvocationOP:
			executionError = p.methodInvocationOP(context, code.Value.(int))
		case GetIdentifierOP:
			executionError = p.getIdentifierOP(context, code.Value.(string))
		case SelectNameFromObjectOP:
			executionError = p.selectNameFromObjectOP(context, code.Value.(string))
		case IndexOP:
			executionError = p.indexOP(context)
		case PushOP:
			executionError = p.pushOP(context)
		case AssignIdentifierOP:
			executionError = p.assignIdentifierOP(context, code.Value.(string))
		case NewClassOP:
			executionError = p.newClassOP(context, bytecode, code.Value.(ClassInformation))
		case NewClassFunctionOP:
			executionError = p.newClassFunctionOP(context, bytecode, code.Value.(FunctionInformation))
		case NewFunctionOP:
			executionError = p.newFunctionOP(context, bytecode, code.Value.(FunctionInformation))
		default:
			panic(instructionNames[code.Instruction.OpCode])
		}
		if executionError != nil {
			// Do Something with the error
			panic(executionError.GetClass(p).Name)
		}
	}
	return p.GetNone(), true
}
