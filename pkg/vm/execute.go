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
		if context.ObjectStack.head != nil {
			current := context.ObjectStack.head
			for ; current != nil; current = current.next {
				fmt.Println(current.value.(*Value).GetClass(p).Name)
			}
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
		case LoadFunctionArgumentsOP:
			executionError = p.loadFunctionArgumentsOP(context, code.Value.([]string))
		case ReturnOP:
			returnResult := p.returnOP(context, code.Value.(int))
			context.ReturnState()
			return returnResult, true
		case IfOneLinerOP:
			executionError = p.ifOneLinerOP(context, bytecode, code.Value.(ConditionInformation))
		case UnlessOneLinerOP:
			executionError = p.unlessOneLinerOP(context, bytecode, code.Value.(ConditionInformation))
		case AssignSelectorOP:
			executionError = p.assignSelectorOP(context, code.Value.(string))
		case AssignIndexOP:
			executionError = p.assignIndexOP(context)
		case IfOP:
			executionError = p.ifOP(context, bytecode, code.Value.(ConditionInformation))
			switch context.LastState {
			case BreakState, RedoState, ContinueState:
				return p.GetNone(), true
			case ReturnState:
				result := context.LastObject
				context.LastObject = nil
				return result, true
			}
		case UnlessOP:
			executionError = p.unlessOP(context, bytecode, code.Value.(ConditionInformation))
			switch context.LastState {
			case BreakState, RedoState, ContinueState:
				return p.GetNone(), true
			case ReturnState:
				result := context.LastObject
				context.LastObject = nil
				return result, true
			}
		case ForLoopOP:
			executionError = p.forLoopOP(context, bytecode, code.Value.(LoopInformation))
			if context.LastState == ReturnState {
				result := context.LastObject
				context.LastObject = nil
				return result, true
			}
		case NewGeneratorOP:
			executionError = p.newGeneratorOP(context, code.Value.(int))
		case WhileLoopOP:
			executionError = p.whileLoopOP(context, bytecode, code.Value.(LoopInformation))
			if context.LastState == ReturnState {
				result := context.LastObject
				context.LastObject = nil
				return result, true
			}
		case DoWhileLoopOP:
			executionError = p.doWhileLoopOP(context, bytecode, code.Value.(LoopInformation))
			if context.LastState == ReturnState {
				result := context.LastObject
				context.LastObject = nil
				return result, true
			}
		case UntilLoopOP:
			executionError = p.untilLoopOP(context, bytecode, code.Value.(LoopInformation))
			if context.LastState == ReturnState {
				result := context.LastObject
				context.LastObject = nil
				return result, true
			}
		case BreakOP:
			context.BreakState()
			return p.GetNone(), true
		case ContinueOP:
			context.ContinueState()
			return p.GetNone(), true
		case RedoOP:
			context.RedoState()
			return p.GetNone(), true
		default:
			panic(instructionNames[code.Instruction.OpCode])
		}
		if executionError != nil {
			// Error state?
			// Do Something with the error
			toString, getError := executionError.Get(p, context, ToString)
			if getError != nil {
				return toString, false
			}
			asString, callError := p.CallFunction(context, toString)
			if !callError {
				return asString, false
			}
			panic(asString.String)
		}
	}
	context.NoState()
	if context.ObjectStack.HasNext() {
		return context.PopObject(), true
	}
	return p.GetNone(), true
}
