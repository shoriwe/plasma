package vm

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/gplasma/pkg/common"
	magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"
	special_symbols "github.com/shoriwe/gplasma/pkg/common/special-symbols"
)

func (ctx *context) pushCode(bytecode []byte) {
	ctx.code.Push(
		&contextCode{
			bytecode: bytecode,
			index:    0,
			onExit:   &common.ListStack[[]byte]{},
		},
	)
}
func (ctx *context) popBlock() {
	// Pop bytecode
	ctxCode := ctx.code.Pop()
	// Load on exit code
	for ctxCode.onExit.HasNext() {
		ctx.pushCode(ctxCode.onExit.Pop())
	}
	// Update symbols
	if ctx.currentSymbols.call != nil {
		ctx.currentSymbols = ctx.currentSymbols.call
	} else {
		ctx.currentSymbols = ctx.currentSymbols.Parent
	}
}

func (plasma *Plasma) prepareClassInitCode(classInfo *ClassInfo) {
	var result []byte
	for _, base := range classInfo.Bases {
		if base.typeId != ClassId {
			panic("no type received as base for class")
		}
		baseClassInfo := base.GetClassInfo()
		if !baseClassInfo.prepared {
			plasma.prepareClassInitCode(baseClassInfo)
		}
		result = append(result, baseClassInfo.Bytecode...)
	}
	result = append(result, classInfo.Bytecode...)
	classInfo.prepared = true
	classInfo.Bytecode = result
}

func (plasma *Plasma) do(ctx *context) {
	ctxCode := ctx.code.Peek()
	instruction := ctxCode.bytecode[ctxCode.index]
	switch instruction {
	case opcodes.Push:
		ctxCode.index++
		ctx.stack.Push(ctx.register)
	case opcodes.Pop:
		ctxCode.index++
		ctx.register = ctx.stack.Pop()
	case opcodes.IdentifierAssign:
		ctxCode.index++
		symbolLength := common.BytesToInt(ctxCode.bytecode[ctxCode.index : ctxCode.index+8])
		ctxCode.index += 8
		symbol := string(ctxCode.bytecode[ctxCode.index : ctxCode.index+symbolLength])
		ctxCode.index += symbolLength
		ctx.currentSymbols.Set(symbol, ctx.stack.Pop())
	case opcodes.SelectorAssign:
		ctxCode.index++
		symbolLength := common.BytesToInt(ctxCode.bytecode[ctxCode.index : ctxCode.index+8])
		ctxCode.index += 8
		symbol := string(ctxCode.bytecode[ctxCode.index : ctxCode.index+symbolLength])
		ctxCode.index += symbolLength
		selector := ctx.stack.Pop()
		selector.Set(symbol, ctx.stack.Pop())
	case opcodes.Label:
		ctxCode.index += 9 // OP + Label
	case opcodes.Jump:
		ctxCode.index += common.BytesToInt(ctxCode.bytecode[1+ctxCode.index : 9+ctxCode.index])
	case opcodes.IfJump:
		if ctx.stack.Pop().Bool() {
			ctxCode.index += common.BytesToInt(ctxCode.bytecode[1+ctxCode.index : 9+ctxCode.index])
		} else {
			ctxCode.index++
		}
	case opcodes.Return:
		ctx.register = ctx.stack.Pop()
		ctx.popBlock()
	case opcodes.DeleteIdentifier:
		ctxCode.index++
		symbolLength := common.BytesToInt(ctxCode.bytecode[ctxCode.index : ctxCode.index+8])
		ctxCode.index += 8
		symbol := string(ctxCode.bytecode[ctxCode.index : ctxCode.index+symbolLength])
		ctxCode.index += symbolLength
		delError := ctx.currentSymbols.Del(symbol)
		if delError != nil {
			panic(delError)
		}
	case opcodes.DeleteSelector:
		ctxCode.index++
		symbolLength := common.BytesToInt(ctxCode.bytecode[ctxCode.index : ctxCode.index+8])
		ctxCode.index += 8
		symbol := string(ctxCode.bytecode[ctxCode.index : ctxCode.index+symbolLength])
		ctxCode.index += symbolLength
		selector := ctx.stack.Pop()
		delError := selector.Del(symbol)
		if delError != nil {
			panic(delError)
		}
	case opcodes.Defer:
		ctxCode.index++
		exprLength := common.BytesToInt(ctxCode.bytecode[ctxCode.index : ctxCode.index+8])
		ctxCode.index += 8
		onExitCode := ctxCode.bytecode[ctxCode.index : ctxCode.index+exprLength]
		ctxCode.index += exprLength
		ctxCode.onExit.Push(onExitCode)
	case opcodes.EnterBlock:
		ctxCode.index++
		ctx.currentSymbols = NewSymbols(ctx.currentSymbols)
	case opcodes.ExitBlock:
		ctxCode.index++
		ctx.popBlock()
	case opcodes.NewFunction:
		ctxCode.index++
		numberOfArgument := common.BytesToInt(ctxCode.bytecode[ctxCode.index : ctxCode.index+8])
		ctxCode.index += 8
		arguments := make([]string, numberOfArgument)
		for i := int64(0); i < numberOfArgument; i++ {
			symbolLength := common.BytesToInt(ctxCode.bytecode[ctxCode.index : ctxCode.index+8])
			ctxCode.index += 8
			symbol := string(ctxCode.bytecode[ctxCode.index : ctxCode.index+symbolLength])
			ctxCode.index += symbolLength
			arguments = append(arguments, symbol)
		}
		bytecodeLength := common.BytesToInt(ctxCode.bytecode[ctxCode.index : ctxCode.index+8])
		ctxCode.index += 8
		bytecode := ctxCode.bytecode[ctxCode.index : ctxCode.index+bytecodeLength]
		ctxCode.index += bytecodeLength
		funcInfo := FuncInfo{
			Arguments: arguments,
			Bytecode:  bytecode,
		}
		funcObject := plasma.NewValue(ctx.currentSymbols, FunctionId, plasma.function)
		funcObject.SetAny(funcInfo)
		ctx.register = funcObject
	case opcodes.NewClass:
		ctxCode.index++
		numberOfBases := common.BytesToInt(ctxCode.bytecode[ctxCode.index : ctxCode.index+8])
		ctxCode.index += 8
		bodyLength := common.BytesToInt(ctxCode.bytecode[ctxCode.index : ctxCode.index+8])
		ctxCode.index += 8
		body := ctxCode.bytecode[ctxCode.index : ctxCode.index+bodyLength]
		ctxCode.index += bodyLength
		// Get bases
		bases := make([]*Value, numberOfBases)
		for i := numberOfBases - 1; i >= 0; i-- {
			bases[i] = ctx.stack.Pop()
		}
		classInfo := &ClassInfo{
			Bases:    bases,
			Bytecode: body,
		}
		classObject := plasma.NewValue(ctx.currentSymbols, ClassId, plasma.class)
		classObject.SetAny(classInfo)
		ctx.register = classObject
	case opcodes.Call:
		ctxCode.index++
		function := ctx.stack.Pop()
		numberOfArguments := common.BytesToInt(ctxCode.bytecode[ctxCode.index : ctxCode.index+8])
		ctxCode.index += 8
		arguments := make([]*Value, numberOfArguments)
		for i := numberOfArguments - 1; i >= 0; i-- {
			arguments[i] = ctx.stack.Pop()
		}
		var callError error
		tries := 0
	doCall:
		if tries == MaxDoCallSearch {
			panic("infinite nested __call__")
		}
		switch function.typeId {
		case BuiltInFunctionId, BuiltInClassId:
			ctx.register, callError = function.Call(arguments...)
			if callError != nil {
				panic(callError)
			}
		case FunctionId:
			funcInfo := function.GetFuncInfo()
			// Push new symbol table based on the function
			newSymbols := NewSymbols(function.vtable)
			newSymbols.call = ctx.currentSymbols
			ctx.currentSymbols = newSymbols
			if int64(len(funcInfo.Arguments)) != numberOfArguments {
				panic("invalid number of argument for function call")
			}
			// Load arguments
			for index, argument := range funcInfo.Arguments {
				ctx.currentSymbols.Set(argument, arguments[index])
			}
			// Push code
			ctx.pushCode(funcInfo.Bytecode)
		case ClassId:
			classInfo := function.GetClassInfo()
			if !classInfo.prepared {
				plasma.prepareClassInitCode(classInfo)
			}
			// Instantiate object
			object := plasma.NewValue(function.vtable, ValueId, plasma.value)
			object.class = function
			object.Set(special_symbols.Self, object)
			// Push object
			ctx.stack.Push(object)
			for _, argument := range arguments {
				ctx.stack.Push(argument)
			}
			// Push class code
			classCode := make([]byte, 0, len(classInfo.Bytecode))
			classCode = append(classCode, classInfo.Bytecode...)
			// inject init code: object.__init__(arguments...)
			classCode = append(classCode, opcodes.Identifier)
			classCode = append(classCode, common.IntToBytes(len(magic_functions.Init))...)
			classCode = append(classCode, magic_functions.Init...)
			classCode = append(classCode, opcodes.Push)
			classCode = append(classCode, opcodes.Call)
			classCode = append(classCode, common.IntToBytes(numberOfArguments)...)
			// Inject pop object to register
			classCode = append(classCode, opcodes.Pop)
			// Load code
			ctx.pushCode(classCode)
			newSymbols := object.vtable
			newSymbols.call = ctx.currentSymbols
			ctx.currentSymbols = newSymbols
		default: // __call__
			call, getError := function.Get(magic_functions.Call)
			if getError != nil {
				panic(getError)
			}
			function = call
			tries++
			goto doCall
		}
	case opcodes.NewArray:
		ctxCode.index++
		numberOfValues := common.BytesToInt(ctxCode.bytecode[ctxCode.index : ctxCode.index+8])
		ctxCode.index += 8
		values := make([]*Value, numberOfValues)
		for i := numberOfValues - 1; i >= 0; i-- {
			values[i] = ctx.stack.Pop()
		}
		ctx.register = plasma.NewArray(values)
	case opcodes.NewTuple:
		ctxCode.index++
		numberOfValues := common.BytesToInt(ctxCode.bytecode[ctxCode.index : ctxCode.index+8])
		ctxCode.index += 8
		values := make([]*Value, numberOfValues)
		for i := numberOfValues - 1; i >= 0; i-- {
			values[i] = ctx.stack.Pop()
		}
		ctx.register = plasma.NewTuple(values)
	case opcodes.NewHash:
		ctxCode.index++
		numberOfValues := common.BytesToInt(ctxCode.bytecode[ctxCode.index : ctxCode.index+8])
		ctxCode.index += 8
		hash := plasma.NewInternalHash()
		for i := numberOfValues - 1; i >= 0; i-- {
			key := ctx.stack.Pop()
			value := ctx.stack.Pop()
			setError := hash.Set(key, value)
			if setError != nil {
				panic(setError)
			}
		}
		ctx.register = plasma.NewHash(hash)
	case opcodes.Identifier:
		ctxCode.index++
		symbolLength := common.BytesToInt(ctxCode.bytecode[ctxCode.index : ctxCode.index+8])
		ctxCode.index += 8
		symbol := string(ctxCode.bytecode[ctxCode.index : ctxCode.index+symbolLength])
		var getError error
		ctx.register, getError = ctx.currentSymbols.Get(symbol)
		if getError != nil {
			panic(getError)
		}
	case opcodes.Integer:
		ctxCode.index++
		value := common.BytesToInt(ctxCode.bytecode[ctxCode.index : ctxCode.index+8])
		ctxCode.index += 8
		ctx.register = plasma.NewInt(value)
	case opcodes.Float:
		ctxCode.index++
		value := common.BytesToFloat(ctxCode.bytecode[ctxCode.index : ctxCode.index+8])
		ctxCode.index += 8
		ctx.register = plasma.NewFloat(value)
	case opcodes.String:
		ctxCode.index++
		stringLength := common.BytesToInt(ctxCode.bytecode[ctxCode.index : ctxCode.index+8])
		ctxCode.index += 8
		contents := ctxCode.bytecode[ctxCode.index : ctxCode.index+stringLength]
		ctxCode.index += stringLength
		ctx.register = plasma.NewString(contents)
	case opcodes.Bytes:
		ctxCode.index++
		stringLength := common.BytesToInt(ctxCode.bytecode[ctxCode.index : ctxCode.index+8])
		ctxCode.index += 8
		contents := ctxCode.bytecode[ctxCode.index : ctxCode.index+stringLength]
		ctxCode.index += stringLength
		ctx.register = plasma.NewBytes(contents)
	case opcodes.True:
		ctx.register = plasma.true
	case opcodes.False:
		ctx.register = plasma.false
	case opcodes.None:
		ctx.register = plasma.none
	case opcodes.Selector:
		ctxCode.index++
		symbolLength := common.BytesToInt(ctxCode.bytecode[ctxCode.index : ctxCode.index+8])
		ctxCode.index += 8
		symbol := string(ctxCode.bytecode[ctxCode.index : ctxCode.index+symbolLength])
		selector := ctx.stack.Pop()
		var getError error
		ctx.register, getError = selector.Get(symbol)
		if getError != nil {
			panic(getError)
		}
	case opcodes.Super:
		break // TODO: Implement me!
	default:
		panic(fmt.Sprintf("unknown opcode %d", instruction))
	}
}
