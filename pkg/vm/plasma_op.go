package vm

import (
	"errors"
	"fmt"
	errors2 "github.com/shoriwe/gplasma/pkg/errors"
)

func (p *Plasma) Execute(context *Context, bytecode *Bytecode) (*Value, bool) {
	if context == nil {
		context = NewContext()
		p.InitializeContext(context)
	}
	// defer fmt.Println(context.ObjectStack.HasNext())
	var success bool
	var object *Value
bytecodeExecutionLoop:
	for bytecode.HasNext() {
		code := bytecode.Next()

		/*
			if code.Line != 0 {
				fmt.Println(color.GreenString(strconv.Itoa(code.Line)), code.Instruction, code.Value)
			} else {
				fmt.Println(color.RedString("UL"), code.Instruction, code.Value)
			}
			if context.ObjectStack.head != nil {
				fmt.Println("Head:", context.ObjectStack.head.value)
			}
			fmt.Println(" Value:", object)
		*/

		switch code.Instruction.OpCode {
		// Literals
		case NewStringOP:
			object, success = p.newStringOP(context, code)
		case NewBytesOP:
			object, success = p.newBytesOP(context, code)
		case NewIntegerOP:
			object, success = p.newIntegerOP(context, code)
		case NewFloatOP:
			object, success = p.newFloatOP(context, code)
		case NewTrueBoolOP:
			object, success = p.newTrueBoolOP()
		case NewFalseBoolOP:
			object, success = p.newFalseBoolOP()
		case NewLambdaFunctionOP:
			object, success = p.newLambdaFunctionOP(context, bytecode, code)
		case GetNoneOP:
			object, success = p.getNoneOP()
		// Composite creation
		case NewTupleOP:
			object, success = p.newTupleOP(context, code)
		case NewArrayOP:
			object, success = p.newArrayOP(context, code)
		case NewHashOP:
			object, success = p.newHashOP(context, code)
		// Unary Expressions
		case NegateBitsOP:
			object, success = p.noArgsGetAndCall(context, NegBits)
		case BoolNegateOP:
			object, success = p.noArgsGetAndCall(context, Negate)
		case NegativeOP:
			object, success = p.noArgsGetAndCall(context, Negative)
		// Binary Expressions
		case AddOP:
			object, success = p.leftBinaryExpressionFuncCall(context, Add)
		case SubOP:
			object, success = p.leftBinaryExpressionFuncCall(context, Sub)
		case MulOP:
			object, success = p.leftBinaryExpressionFuncCall(context, Mul)
		case DivOP:
			object, success = p.leftBinaryExpressionFuncCall(context, Div)
		case FloorDivOP:
			object, success = p.leftBinaryExpressionFuncCall(context, FloorDiv)
		case ModOP:
			object, success = p.leftBinaryExpressionFuncCall(context, Mod)
		case PowOP:
			object, success = p.leftBinaryExpressionFuncCall(context, Pow)
		case BitXorOP:
			object, success = p.leftBinaryExpressionFuncCall(context, BitXor)
		case BitAndOP:
			object, success = p.leftBinaryExpressionFuncCall(context, BitAnd)
		case BitOrOP:
			object, success = p.leftBinaryExpressionFuncCall(context, BitOr)
		case BitLeftOP:
			object, success = p.leftBinaryExpressionFuncCall(context, BitLeft)
		case BitRightOP:
			object, success = p.leftBinaryExpressionFuncCall(context, BitRight)
		case AndOP:
			object, success = p.leftBinaryExpressionFuncCall(context, And)
		case OrOP:
			object, success = p.leftBinaryExpressionFuncCall(context, Or)
		case XorOP:
			object, success = p.leftBinaryExpressionFuncCall(context, Xor)
		case EqualsOP:
			object, success = p.leftBinaryExpressionFuncCall(context, Equals)
		case NotEqualsOP:
			object, success = p.leftBinaryExpressionFuncCall(context, NotEquals)
		case GreaterThanOP:
			object, success = p.leftBinaryExpressionFuncCall(context, GreaterThan)
		case LessThanOP:
			object, success = p.leftBinaryExpressionFuncCall(context, LessThan)
		case GreaterThanOrEqualOP:
			object, success = p.leftBinaryExpressionFuncCall(context, GreaterThanOrEqual)
		case LessThanOrEqualOP:
			object, success = p.leftBinaryExpressionFuncCall(context, LessThanOrEqual)
		case ContainsOP:
			// This operation is inverted, right is left and left is right
			leftHandSide := context.PopObject()
			rightHandSide := context.PopObject()
			context.PushObject(leftHandSide)
			context.PushObject(rightHandSide)
			object, success = p.leftBinaryExpressionFuncCall(context, Contains)
		// Other Expressions
		case GetIdentifierOP:
			object, success = p.getIdentifierOP(context, code)
		case IndexOP:
			object, success = p.indexOP(context)
		case SelectNameFromObjectOP:
			object, success = p.selectNameFromObjectOP(context, code)
		case MethodInvocationOP:
			object, success = p.methodInvocationOP(context, code)
		case NewIteratorOP:
			object, success = p.newIteratorOP(context, bytecode, code)
		// Assign Statement
		case AssignIdentifierOP:
			success = p.assignIdentifierOP(context, code)
		case AssignSelectorOP:
			object, success = p.assignSelectorOP(context, code)
		case AssignIndexOP:
			object, success = p.assignIndexOP(context)
		case ReturnOP:
			object, success = p.returnOP(context, code)
			break bytecodeExecutionLoop
		// Special Instructions
		case LoadFunctionArgumentsOP:
			success = p.loadFunctionArgumentsOP(context, code)
		case NewFunctionOP:
			success = p.newFunctionOP(context, bytecode, code)
		case IfJumpOP:
			object, success = p.ifJumpOP(context, bytecode, code)
		case UnlessJumpOP:
			object, success = p.unlessJumpOP(context, bytecode, code)
		case SetupLoopOP:
			object, success = p.setupForLoopOP(context, code)
		case PopLoopOP:
			context.LoopStack.Pop()
		case LoadForReloadOP:
			success = p.loadForLoopArguments(context)
		case UnpackForLoopOP:
			object, success = p.unpackForLoopOP(context, bytecode)
		case RedoOP, BreakOP, ContinueOP, JumpOP:
			success = p.jumpOP(bytecode, code)
		case PushOP:
			if object != nil {
				context.ObjectStack.Push(object)
				object = nil
			}
		case PopOP:
			context.ObjectStack.Pop()
		case NOP:
			break
		case SetupTryOP:
			success = p.setupTryOP(context, bytecode, code)
		case PopTryOP:
			success = p.popTryOP(context)
		case ExceptOP:
			object, success = p.exceptOP(context, bytecode, code)
		case NewModuleOP:
			object, success = p.newModuleOP(context, bytecode, code)
		case RaiseOP:
			object, success = p.raiseOP(context)
		case NewClassOP:
			object, success = p.newClassOP(context, bytecode, code)
		case NewClassFunctionOP:
			success = p.newClassFunctionOP(context, bytecode, code)
		case CaseOP:
			object, success = p.caseOP(context, bytecode, code)
		default:
			panic(fmt.Sprintf("Unknown VM instruction %d", code.Instruction.OpCode))
		}
		if !success {
			// Handle Try and excepts
			if context.TryStack.HasNext() {
				p.tryJumpOP(context, bytecode, object)
				continue
			}
			return object, false
		}
	}
	if context.ObjectStack.HasNext() {
		return context.PopObject(), true
	}
	return p.GetNone(), true
}

func (p *Plasma) setupTryOP(context *Context, bytecode *Bytecode, code Code) bool {
	context.TryStack.Push(
		&TrySettings{
			StartIndex: bytecode.index,
			BodyLength: code.Value.(int),
		},
	)
	return true
}

func (p *Plasma) popTryOP(context *Context) bool {
	context.TryStack.Pop()
	return true
}

func (p *Plasma) exceptOP(context *Context, bytecode *Bytecode, code Code) (*Value, bool) {
	failed := context.TryStack.Peek().LastError
	targets := context.PopObject()
	runtimeError := p.ForceMasterGetAny(RuntimeError)
	if targets.GetLength() == 0 {
		context.TryStack.Peek().LastError = nil
		return nil, true
	}
	for _, target := range targets.GetContent() {
		if !target.IsTypeById(TypeId) {
			return p.NewInvalidTypeError(context, target.TypeName(), TypeName), false
		}
		if !target.Implements(runtimeError) {
			return p.NewInvalidTypeError(context, target.TypeName(), RuntimeError), false
		}
		if failed.Implements(target) {
			// Assign the error to the receiver
			context.TryStack.Peek().LastError = nil
			context.PeekSymbolTable().Set(code.Value.([2]interface{})[0].(string), failed)
			return nil, true
		}
	}
	// Jump to the next except, else or finally
	bytecode.Jump(code.Value.([2]interface{})[1].(int))
	return nil, true
}

func (p *Plasma) tryJumpOP(context *Context, bytecode *Bytecode, failed *Value) {
	bytecode.index = context.TryStack.Peek().StartIndex
	bytecode.Jump(context.TryStack.Peek().BodyLength + 1)
	context.TryStack.head.value.(*TrySettings).LastError = failed
}

func (p *Plasma) newStringOP(context *Context, code Code) (*Value, bool) {
	value := code.Value.(string)
	stringObject := p.NewString(context, false, context.SymbolStack.Peek(), value)
	return stringObject, true
}

func (p *Plasma) newBytesOP(context *Context, code Code) (*Value, bool) {
	value := code.Value.([]byte)
	return p.NewBytes(context, false, context.SymbolStack.Peek(), value), true
}

func (p *Plasma) newIntegerOP(context *Context, code Code) (*Value, bool) {
	value := code.Value.(int64)
	return p.NewInteger(context, false, context.SymbolStack.Peek(), value), true
}

func (p *Plasma) newFloatOP(context *Context, code Code) (*Value, bool) {
	value := code.Value.(float64)
	return p.NewFloat(context, false, context.SymbolStack.Peek(), value), true
}

func (p *Plasma) newTrueBoolOP() (*Value, bool) {
	return p.GetTrue(), true
}

func (p *Plasma) newFalseBoolOP() (*Value, bool) {
	return p.GetFalse(), true
}

func (p *Plasma) getNoneOP() (*Value, bool) {
	return p.GetNone(), true
}

func (p *Plasma) newTupleOP(context *Context, code Code) (*Value, bool) {
	numberOfValues := code.Value.(int)
	var values []*Value
	for i := 0; i < numberOfValues; i++ {
		values = append(values, context.PopObject())
	}
	return p.NewTuple(context, false, context.PeekSymbolTable(), values), true
}

func (p *Plasma) newArrayOP(context *Context, code Code) (*Value, bool) {
	numberOfValues := code.Value.(int)
	var values []*Value
	for i := 0; i < numberOfValues; i++ {
		values = append(values, context.PopObject())
	}
	return p.NewArray(context, false, context.PeekSymbolTable(), values), true
}

func (p *Plasma) newHashOP(context *Context, code Code) (*Value, bool) {
	numberOfValues := code.Value.(int)
	var keyValues []*KeyValue
	for i := 0; i < numberOfValues; i++ {

		key := context.PopObject()
		value := context.PopObject()
		keyValues = append(keyValues, &KeyValue{
			Key:   key,
			Value: value,
		})
	}
	hashTable := p.NewHashTable(context, false, context.PeekSymbolTable(), map[int64][]*KeyValue{}, numberOfValues)
	hashTableAssign, getError := hashTable.Get(Assign)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(context, hashTable.GetClass(p), Assign), false
	}
	for _, keyValue := range keyValues {
		assignError, success := p.CallFunction(context, hashTableAssign, keyValue.Key, keyValue.Value)
		if !success {
			return assignError, false
		}
	}
	return hashTable, true
}

func (p *Plasma) newLambdaFunctionOP(context *Context, bytecode *Bytecode, code Code) (*Value, bool) {
	functionInformation := code.Value.([2]int)
	codeLength := functionInformation[0]
	numberOfArguments := functionInformation[1]
	start := bytecode.index
	bytecode.index += codeLength
	end := bytecode.index
	functionCode := make([]Code, codeLength)
	copy(functionCode, bytecode.instructions[start:end])
	return p.NewFunction(context, false, context.PeekSymbolTable(), NewPlasmaFunction(numberOfArguments, functionCode)), true
}

// Useful function to call those built ins that doesn't receive arguments of an object
func (p *Plasma) noArgsGetAndCall(context *Context, operationName string) (*Value, bool) {
	object := context.PopObject()
	operation, getError := object.Get(operationName)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(context, object.GetClass(p), operationName), false
	}
	return p.CallFunction(context, operation)
}

// Function useful to cal object built-ins binary expression functions
func (p *Plasma) leftBinaryExpressionFuncCall(context *Context, operationName string) (*Value, bool) {
	leftHandSide := context.PopObject()
	rightHandSide := context.PopObject()
	operation, getError := leftHandSide.Get(operationName)
	if getError != nil {
		return p.rightBinaryExpressionFuncCall(context, leftHandSide, rightHandSide, operationName)
	}
	result, success := p.CallFunction(context, operation, rightHandSide)
	if !success {
		return p.rightBinaryExpressionFuncCall(context, leftHandSide, rightHandSide, operationName)
	}
	return result, true
}

func (p *Plasma) rightBinaryExpressionFuncCall(context *Context, leftHandSide *Value, rightHandSide *Value, operationName string) (*Value, bool) {
	operation, getError := rightHandSide.Get("Right" + operationName)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(context, p.ForceMasterGetAny(ValueName), "Right"+operationName), false
	}
	return p.CallFunction(context, operation, leftHandSide)
}

func (p *Plasma) indexOP(context *Context) (*Value, bool) {
	index := context.PopObject()
	source := context.PopObject()
	indexOperation, getError := source.Get(Index)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(context, source.GetClass(p), Index), false
	}
	return p.CallFunction(context, indexOperation, index)
}

func (p *Plasma) selectNameFromObjectOP(context *Context, code Code) (*Value, bool) {
	name := code.Value.(string)
	source := context.PopObject()
	value, getError := source.Get(name)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(context, source.GetClass(p), name), false
	}
	return value, true
}

func (p *Plasma) methodInvocationOP(context *Context, code Code) (*Value, bool) {
	numberOfArguments := code.Value.(int)
	function := context.PopObject()
	var arguments []*Value
	for i := 0; i < numberOfArguments; i++ {
		if !context.ObjectStack.HasNext() {
			return p.NewInvalidNumberOfArgumentsError(context, i, numberOfArguments), false
		}
		arguments = append(arguments, context.PopObject())
	}
	return p.CallFunction(context, function, arguments...)
}

func (p *Plasma) getIdentifierOP(context *Context, code Code) (*Value, bool) {
	value, getError := context.PeekSymbolTable().GetAny(code.Value.(string))
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(context, p.ForceMasterGetAny(ValueName).GetClass(p), code.Value.(string)), false
	}
	return value, true
}

func (p *Plasma) newIteratorOP(context *Context, bytecode *Bytecode, code Code) (*Value, bool) {
	source := context.PopObject()
	var iterSource *Value
	var success bool
	if source.IsTypeById(IteratorId) {
		iterSource = source
	} else {
		iter, getError := source.Get(Iter)
		if getError != nil {
			return p.NewObjectWithNameNotFoundError(context, source.GetClass(p), Iter), true
		}
		iterSource, success = p.CallFunction(context, iter)
		if !success {
			return iterSource, false
		}
	}
	generatorIterator := p.NewIterator(context, false, context.PeekSymbolTable())
	generatorIterator.Set(Source, iterSource)

	hasNextCodeLength, nextCodeLength := code.Value.([2]int)[0], code.Value.([2]int)[1]
	var hasNextCode []Code
	for i := 0; i < hasNextCodeLength; i++ {
		hasNextCode = append(hasNextCode, bytecode.Next())
	}
	var nextCode []Code
	for i := 0; i < nextCodeLength; i++ {
		nextCode = append(nextCode, bytecode.Next())
	}
	generatorIterator.SetOnDemandSymbol(Next,
		func() *Value {
			return p.NewFunction(context, false, generatorIterator.symbols,
				NewPlasmaClassFunction(generatorIterator, 0, nextCode),
			)
		},
	)
	generatorIterator.SetOnDemandSymbol(HasNext,
		func() *Value {
			return p.NewFunction(context, false, generatorIterator.symbols,
				NewPlasmaClassFunction(generatorIterator, 0, hasNextCode),
			)
		},
	)
	return generatorIterator, true
}

func (p *Plasma) setupForLoopOP(context *Context, code Code) (*Value, bool) {
	source := context.ObjectStack.Pop()
	sourceNext, nextGetError := source.Get(Next)
	sourceHasNext, hasNextGetError := source.Get(HasNext)
	var sourceIter *Value
	if nextGetError == nil && hasNextGetError == nil {
		sourceIter = source
	} else {
		sourceToIter, getError := source.Get(Iter)
		if getError != nil {
			return p.NewObjectWithNameNotFoundError(context, source.GetClass(p), Iter), false
		}
		var success bool
		sourceIter, success = p.CallFunction(context, sourceToIter)
		if !success {
			return sourceIter, false
		}
		sourceNext, nextGetError = sourceIter.Get(Next)
		sourceHasNext, hasNextGetError = sourceIter.Get(HasNext)
		if nextGetError != nil {
			return p.NewObjectWithNameNotFoundError(context, sourceIter.GetClass(p), Next), false
		} else if hasNextGetError != nil {
			return p.NewObjectWithNameNotFoundError(context, sourceIter.GetClass(p), HasNext), false
		}
	}

	context.LoopStack.Push(
		&LoopSettings{
			Receivers:         code.Value.([2]interface{})[0].([]string),
			NumberOfReceivers: len(code.Value.([2]interface{})[0].([]string)),
			Source:            sourceIter,
			Next:              sourceNext,
			HasNext:           sourceHasNext,
			Jump:              code.Value.([2]interface{})[1].(int),
			MappedReceivers:   map[string]*Value{},
		},
	)
	return nil, true
}

func (p *Plasma) loadForLoopArguments(context *Context) bool {
	for name, value := range context.LoopStack.Peek().MappedReceivers {
		context.PeekSymbolTable().Set(name, value)
	}
	return true
}

func (p *Plasma) unpackForArguments(context *Context, loopSettings *LoopSettings, result *Value) (*Value, bool) {
	if loopSettings.NumberOfReceivers == 1 {
		loopSettings.MappedReceivers[loopSettings.Receivers[0]] = result
	} else if result.IsTypeById(StringId) {
		if result.GetLength() != loopSettings.NumberOfReceivers {
			return p.NewInvalidNumberOfArgumentsError(context, loopSettings.NumberOfReceivers, result.GetLength()), false
		}
		for index, name := range loopSettings.Receivers {
			loopSettings.MappedReceivers[name] = p.NewString(context, false, context.PeekSymbolTable(), string(result.GetString()[index]))
		}
	} else if result.IsTypeById(BytesId) {
		if result.GetLength() != loopSettings.NumberOfReceivers {
			return p.NewInvalidNumberOfArgumentsError(context, loopSettings.NumberOfReceivers, result.GetLength()), false
		}
		for index, name := range loopSettings.Receivers {
			loopSettings.MappedReceivers[name] = p.NewInteger(context, false, context.PeekSymbolTable(), int64(result.GetBytes()[index]))
		}
	} else if result.IsTypeById(ArrayId) {
		if result.GetLength() != loopSettings.NumberOfReceivers {
			return p.NewInvalidNumberOfArgumentsError(context, loopSettings.NumberOfReceivers, result.GetLength()), false
		}
		for index, name := range loopSettings.Receivers {
			loopSettings.MappedReceivers[name] = result.GetContent()[index]
		}
	} else if result.IsTypeById(TupleId) {
		if result.GetLength() != loopSettings.NumberOfReceivers {
			return p.NewInvalidNumberOfArgumentsError(context, loopSettings.NumberOfReceivers, result.GetLength()), false
		}
		for index, name := range loopSettings.Receivers {
			loopSettings.MappedReceivers[name] = result.GetContent()[index]
		}
	} else if !result.IsTypeById(HashTableId) {
		return p.NewGoRuntimeError(context, errors.New("HashTable doesn't support unpacking")), false
	} else {
		var hasNext *Value
		var next *Value
		var resultAsIterator *Value
		var hasNextGetError *errors2.Error
		var nextGetError *errors2.Error
		hasNext, hasNextGetError = result.Get(HasNext)
		next, nextGetError = result.Get(Next)
		if hasNextGetError != nil || nextGetError != nil {
			resultToIter, getError := result.Get(Iter)
			if getError != nil {
				return p.NewObjectWithNameNotFoundError(context, result.GetClass(p), Iter), false
			}
			var success bool
			resultAsIterator, success = p.CallFunction(context, resultToIter)
			if !success {
				return resultAsIterator, false
			}
		} else {
			resultAsIterator = result
		}
		for index, name := range loopSettings.Receivers {
			hasNextResult, success := p.CallFunction(context, hasNext)
			if !success {
				return hasNextResult, false
			}
			hasNextResultBool, callError := p.QuickGetBool(context, hasNextResult)
			if callError != nil {
				return callError, false
			}
			if !hasNextResultBool {
				return p.NewInvalidNumberOfArgumentsError(context, index, loopSettings.NumberOfReceivers), false
			}
			var nextResult *Value
			nextResult, success = p.CallFunction(context, next)
			if !success {
				return nextResult, false
			}
			context.PeekSymbolTable().Set(name, nextResult)
		}
	}
	return nil, true
}

func (p *Plasma) unpackForLoopOP(context *Context, bytecode *Bytecode) (*Value, bool) {
	loopSettings := context.LoopStack.Peek()
	hasNext, success := p.CallFunction(context, loopSettings.HasNext)
	if !success {
		return hasNext, false
	}
	hasNextBool, callError := p.QuickGetBool(context, hasNext)
	if callError != nil {
		return callError, false
	}
	if !hasNextBool {
		// Do the jump to outside the loop
		bytecode.Jump(loopSettings.Jump)
		return nil, true
	}
	var nextValue *Value
	nextValue, success = p.CallFunction(context, loopSettings.Next)
	if !success {
		return nextValue, false
	}
	return p.unpackForArguments(context, loopSettings, nextValue)
}

// Assign Statement

func (p *Plasma) assignIdentifierOP(context *Context, code Code) bool {
	identifier := code.Value.(string)
	context.PeekSymbolTable().Set(identifier, context.PopObject())
	return true
}

func (p *Plasma) assignSelectorOP(context *Context, code Code) (*Value, bool) {
	target := context.PopObject()
	value := context.PopObject()
	identifier := code.Value.(string)
	if target.IsBuiltIn() {
		return p.NewBuiltInSymbolProtectionError(context, identifier), false
	}
	target.Set(identifier, value)
	return nil, true
}

func (p *Plasma) assignIndexOP(context *Context) (*Value, bool) {
	index := context.PopObject()
	source := context.PopObject()
	value := context.PopObject()
	sourceAssign, getError := source.Get(Assign)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(context, source.GetClass(p), Assign), false
	}
	callError, success := p.CallFunction(context, sourceAssign, index, value)
	if !success {
		return callError, false
	}
	return nil, true
}

func (p *Plasma) returnOP(context *Context, code Code) (*Value, bool) {
	numberOfReturnValues := code.Value.(int)
	if numberOfReturnValues == 0 {
		context.PushObject(p.GetNone())
		return nil, true
	}
	if numberOfReturnValues == 1 {
		if !context.ObjectStack.HasNext() {
			return p.NewInvalidNumberOfArgumentsError(context, 1, numberOfReturnValues), false
		}
		return nil, true
	}

	var values []*Value
	for i := 0; i < numberOfReturnValues; i++ {
		if !context.ObjectStack.HasNext() {
			return p.NewInvalidNumberOfArgumentsError(context, i, numberOfReturnValues), false
		}
		values = append(values, context.PopObject())
	}
	context.PushObject(p.NewTuple(context, false, context.PeekSymbolTable(), values))
	return nil, true
}

// Special Instructions

func (p *Plasma) loadFunctionArgumentsOP(context *Context, code Code) bool {
	for _, argument := range code.Value.([]string) {
		value := context.PopObject()
		context.PeekSymbolTable().Set(argument, value)
	}
	return true
}

func (p *Plasma) newFunctionOP(context *Context, bytecode *Bytecode, code Code) bool {
	functionInformation := code.Value.([2]int)
	codeLength := functionInformation[0]
	numberOfArguments := functionInformation[1]
	start := bytecode.index
	bytecode.index += codeLength
	end := bytecode.index
	functionCode := make([]Code, codeLength)
	copy(functionCode, bytecode.instructions[start:end])
	context.PushObject(p.NewFunction(context, false, context.PeekSymbolTable(), NewPlasmaFunction(numberOfArguments, functionCode)))
	return true
}

func (p *Plasma) jumpOP(bytecode *Bytecode, code Code) bool {
	bytecode.index += code.Value.(int)
	return true
}

func (p *Plasma) ifJumpOP(context *Context, bytecode *Bytecode, code Code) (*Value, bool) {
	condition := context.ObjectStack.Pop()
	conditionBool, callError := p.QuickGetBool(context, condition)
	if callError != nil {
		return callError, false
	}
	if !conditionBool {
		bytecode.index += code.Value.(int)
	}
	return nil, true
}

func (p *Plasma) unlessJumpOP(context *Context, bytecode *Bytecode, code Code) (*Value, bool) {
	condition := context.ObjectStack.Pop()
	conditionBool, callError := p.QuickGetBool(context, condition)
	if callError != nil {
		return callError, false
	}
	if conditionBool {
		bytecode.index += code.Value.(int)
	}
	return nil, true
}

type ModuleInformation struct {
	Name       string
	CodeLength int
}

func (p *Plasma) newModuleOP(context *Context, bytecode *Bytecode, code Code) (*Value, bool) {
	moduleInformation := code.Value.(ModuleInformation)
	var moduleBody []Code
	for i := 0; i < moduleInformation.CodeLength; i++ {
		moduleBody = append(moduleBody, bytecode.Next())
	}
	module := p.NewModule(context, false, context.PeekSymbolTable())
	context.PushSymbolTable(module.SymbolTable())
	callError, success := p.Execute(context, NewBytecodeFromArray(moduleBody))
	if !success {
		return callError, false
	}
	context.PopSymbolTable()
	context.PeekSymbolTable().Set(moduleInformation.Name, module)
	return nil, true
}

type ClassInformation struct {
	Name       string
	BodyLength int
}

func (p *Plasma) newClassOP(context *Context, bytecode *Bytecode, code Code) (*Value, bool) {
	classInformation := code.Value.(ClassInformation)
	rawSubClasses := context.PopObject().GetContent()
	var subClasses []*Value
	for _, subClass := range rawSubClasses {
		if !subClass.IsTypeById(TypeId) {
			return p.NewInvalidTypeError(context, subClass.TypeName(), TypeName), false
		}
		subClasses = append(subClasses, subClass)
	}

	var classBody []Code
	for i := 0; i < classInformation.BodyLength; i++ {
		classBody = append(classBody, bytecode.Next())
	}
	class := p.NewType(
		context,
		false,
		classInformation.Name,
		context.PeekSymbolTable(),
		subClasses,
		NewPlasmaConstructor(classBody),
	)
	context.PeekSymbolTable().Set(classInformation.Name, class)
	return nil, true
}

func (p *Plasma) newClassFunctionOP(context *Context, bytecode *Bytecode, code Code) bool {
	functionInformation := code.Value.([2]int)
	codeLength := functionInformation[0]
	numberOfArguments := functionInformation[1]
	start := bytecode.index
	bytecode.index += codeLength
	end := bytecode.index
	functionCode := make([]Code, codeLength)
	copy(functionCode, bytecode.instructions[start:end])
	context.PushObject(p.NewFunction(context, false, context.PeekSymbolTable(), NewPlasmaClassFunction(context.PeekObject(), numberOfArguments, functionCode)))
	return true
}

func (p *Plasma) raiseOP(context *Context) (*Value, bool) {
	if !context.PeekObject().IsTypeById(TypeId) {
		return p.NewInvalidTypeError(context, context.PeekObject().TypeName(), RuntimeError), false
	}
	if !context.PeekObject().Implements(p.ForceMasterGetAny(RuntimeError)) {
		return p.NewInvalidTypeError(context, context.PeekObject().TypeName(), RuntimeError), false
	}
	return context.PeekObject(), true
}

func (p *Plasma) caseOP(context *Context, bytecode *Bytecode, code Code) (*Value, bool) {
	references := context.PopObject()
	contains := p.ForceGetSelf(Contains, references)
	result, success := p.CallFunction(context, contains, context.PeekObject())
	if !success {
		return result, false
	}
	boolResult, callError := p.QuickGetBool(context, result)
	if callError != nil {
		return callError, false
	}
	if !boolResult {
		bytecode.index += code.Value.(int)
		return nil, true
	}
	context.PopObject()
	return nil, true
}
