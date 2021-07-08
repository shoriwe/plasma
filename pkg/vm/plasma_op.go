package vm

import (
	"errors"
	"fmt"
	errors2 "github.com/shoriwe/gplasma/pkg/errors"
)

func (p *Plasma) Execute(context *Context, bytecode *Bytecode) (Value, *Object) {
	if context == nil {
		context = NewContext()
		p.InitializeContext(context)
	}
	// defer fmt.Println(context.ObjectStack.HasNext())
	var executionError *Object
	var object Value
bytecodeExecutionLoop:
	for ; bytecode.HasNext(); {
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
			fmt.Println("Object:", object)
		*/

		switch code.Instruction.OpCode {
		// Literals
		case NewStringOP:
			object, executionError = p.newStringOP(context, code)
		case NewBytesOP:
			object, executionError = p.newBytesOP(context, code)
		case NewIntegerOP:
			object, executionError = p.newIntegerOP(context, code)
		case NewFloatOP:
			object, executionError = p.newFloatOP(context, code)
		case NewTrueBoolOP:
			object, executionError = p.newTrueBoolOP()
		case NewFalseBoolOP:
			object, executionError = p.newFalseBoolOP()
		case NewParenthesesOP:
			object, executionError = p.newParenthesesOP(context)
		case NewLambdaFunctionOP:
			object, executionError = p.newLambdaFunctionOP(context, bytecode, code)
		case GetNoneOP:
			object, executionError = p.getNoneOP()
		// Composite creation
		case NewTupleOP:
			object, executionError = p.newTupleOP(context, code)
		case NewArrayOP:
			object, executionError = p.newArrayOP(context, code)
		case NewHashOP:
			object, executionError = p.newHashOP(context, code)
		// Unary Expressions
		case NegateBitsOP:
			object, executionError = p.noArgsGetAndCall(context, NegBits)
		case BoolNegateOP:
			object, executionError = p.noArgsGetAndCall(context, Negate)
		case NegativeOP:
			object, executionError = p.noArgsGetAndCall(context, Negative)
		// Binary Expressions
		case AddOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, Add)
		case SubOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, Sub)
		case MulOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, Mul)
		case DivOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, Div)
		case FloorDivOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, FloorDiv)
		case ModOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, Mod)
		case PowOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, Pow)
		case BitXorOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, BitXor)
		case BitAndOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, BitAnd)
		case BitOrOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, BitOr)
		case BitLeftOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, BitLeft)
		case BitRightOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, BitRight)
		case AndOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, And)
		case OrOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, Or)
		case XorOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, Xor)
		case EqualsOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, Equals)
		case NotEqualsOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, NotEquals)
		case GreaterThanOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, GreaterThan)
		case LessThanOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, LessThan)
		case GreaterThanOrEqualOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, GreaterThanOrEqual)
		case LessThanOrEqualOP:
			object, executionError = p.leftBinaryExpressionFuncCall(context, LessThanOrEqual)
		case ContainsOP:
			// This operation is inverted, right is left and left is right
			leftHandSide := context.PopObject()
			rightHandSide := context.PopObject()
			context.PushObject(leftHandSide)
			context.PushObject(rightHandSide)
			object, executionError = p.leftBinaryExpressionFuncCall(context, Contains)
		// Other Expressions
		case GetIdentifierOP:
			object, executionError = p.getIdentifierOP(context, code)
		case IndexOP:
			object, executionError = p.indexOP(context)
		case SelectNameFromObjectOP:
			object, executionError = p.selectNameFromObjectOP(context, code)
		case MethodInvocationOP:
			object, executionError = p.methodInvocationOP(context, code)
		case NewIteratorOP:
			object, executionError = p.newIteratorOP(context, bytecode, code)
		// Assign Statement
		case AssignIdentifierOP:
			executionError = p.assignIdentifierOP(context, code)
		case AssignSelectorOP:
			executionError = p.assignSelectorOP(context, code)
		case AssignIndexOP:
			executionError = p.assignIndexOP(context)
		case ReturnOP:
			executionError = p.returnOP(context, code)
			break bytecodeExecutionLoop
		// Special Instructions
		case LoadFunctionArgumentsOP:
			executionError = p.loadFunctionArgumentsOP(context, code)
		case NewFunctionOP:
			executionError = p.newFunctionOP(context, bytecode, code)
		case IfJumpOP:
			executionError = p.ifJumpOP(context, bytecode, code)
		case UnlessJumpOP:
			executionError = p.unlessJumpOP(context, bytecode, code)
		case SetupLoopOP:
			executionError = p.setupForLoopOP(context, code)
		case PopLoopOP:
			context.LoopStack.Pop()
		case LoadForReloadOP:
			executionError = p.loadForLoopArguments(context)
		case UnpackForLoopOP:
			executionError = p.unpackForLoopOP(context, bytecode)
		case RedoOP, BreakOP, ContinueOP, JumpOP:
			executionError = p.jumpOP(bytecode, code)
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
			executionError = p.setupTryOP(context, bytecode, code)
		case PopTryOP:
			executionError = p.popTryOP(context)
		case ExceptOP:
			executionError = p.exceptOP(context, bytecode, code)
		case NewModuleOP:
			executionError = p.newModuleOP(context, bytecode, code)
		case RaiseOP:
			executionError = p.raiseOP(context)
		case NewClassOP:
			executionError = p.newClassOP(context, bytecode, code)
		case NewClassFunctionOP:
			executionError = p.newClassFunctionOP(context, bytecode, code)
		case CaseOP:
			executionError = p.caseOP(context, bytecode, code)
		default:
			panic(fmt.Sprintf("Unknown VM instruction %d", code.Instruction.OpCode))
		}
		if executionError != nil {
			// Handle Try and excepts
			if context.TryStack.HasNext() {
				p.tryJumpOP(context, bytecode, executionError)
				continue
			}
			return nil, executionError
		}
	}
	if context.ObjectStack.HasNext() {
		return context.PopObject(), nil
	}
	return p.GetNone(), nil
}

func (p *Plasma) setupTryOP(context *Context, bytecode *Bytecode, code Code) *Object {
	context.TryStack.Push(
		&TrySettings{
			StartIndex: bytecode.index,
			BodyLength: code.Value.(int),
		},
	)
	return nil
}

func (p *Plasma) popTryOP(context *Context) *Object {
	context.TryStack.Pop()
	return nil
}

func (p *Plasma) exceptOP(context *Context, bytecode *Bytecode, code Code) *Object {
	executionError := context.TryStack.Peek().LastError
	targets := context.PopObject()
	runtimeError := p.ForceMasterGetAny(RuntimeError).(*Type)
	if targets.GetLength() == 0 {
		context.TryStack.Peek().LastError = nil
		return nil
	}
	for _, target := range targets.GetContent() {
		if _, ok := target.(*Type); !ok {
			return p.NewInvalidTypeError(context, target.TypeName(), TypeName)
		}
		if !target.Implements(runtimeError) {
			return p.NewInvalidTypeError(context, target.TypeName(), RuntimeError)
		}
		if executionError.Implements(target.(*Type)) {
			// Assign the error to the receiver
			context.TryStack.Peek().LastError = nil
			context.PeekSymbolTable().Set(code.Value.([2]interface{})[0].(string), executionError)
			return nil
		}
	}
	// Jump to the next except, else or finally
	bytecode.Jump(code.Value.([2]interface{})[1].(int))
	return nil
}

func (p *Plasma) tryJumpOP(context *Context, bytecode *Bytecode, executionError Value) {
	bytecode.index = context.TryStack.Peek().StartIndex
	bytecode.Jump(context.TryStack.Peek().BodyLength + 1)
	context.TryStack.head.value.(*TrySettings).LastError = executionError
}

func (p *Plasma) newStringOP(context *Context, code Code) (Value, *Object) {
	value := code.Value.(string)
	stringObject := p.NewString(context, false, context.SymbolStack.Peek(), value)
	return stringObject, nil
}

func (p *Plasma) newBytesOP(context *Context, code Code) (Value, *Object) {
	value := code.Value.([]byte)
	return p.NewBytes(context, false, context.SymbolStack.Peek(), value), nil
}

func (p *Plasma) newIntegerOP(context *Context, code Code) (Value, *Object) {
	value := code.Value.(int64)
	return p.NewInteger(context, false, context.SymbolStack.Peek(), value), nil
}

func (p *Plasma) newFloatOP(context *Context, code Code) (Value, *Object) {
	value := code.Value.(float64)
	return p.NewFloat(context, false, context.SymbolStack.Peek(), value), nil
}

func (p *Plasma) newTrueBoolOP() (Value, *Object) {
	return p.GetTrue(), nil
}

func (p *Plasma) newFalseBoolOP() (Value, *Object) {
	return p.GetFalse(), nil
}

func (p *Plasma) getNoneOP() (Value, *Object) {
	return p.GetNone(), nil
}

func (p *Plasma) newTupleOP(context *Context, code Code) (Value, *Object) {
	numberOfValues := code.Value.(int)
	var values []Value
	for i := 0; i < numberOfValues; i++ {
		values = append(values, context.PopObject())
	}
	return p.NewTuple(context, false, context.PeekSymbolTable(), values), nil
}

func (p *Plasma) newArrayOP(context *Context, code Code) (Value, *Object) {
	numberOfValues := code.Value.(int)
	var values []Value
	for i := 0; i < numberOfValues; i++ {
		values = append(values, context.PopObject())
	}
	return p.NewArray(context, false, context.PeekSymbolTable(), values), nil
}

func (p *Plasma) newHashOP(context *Context, code Code) (Value, *Object) {
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
		return nil, p.NewObjectWithNameNotFoundError(context, hashTable.GetClass(p), Assign)
	}
	for _, keyValue := range keyValues {
		_, assignError := p.CallFunction(context, hashTableAssign, hashTable.SymbolTable(), keyValue.Key, keyValue.Value)
		if assignError != nil {
			return nil, assignError
		}
	}
	return hashTable, nil
}

func (p *Plasma) newParenthesesOP(context *Context) (Value, *Object) {
	return context.PopObject(), nil
}

func (p *Plasma) newLambdaFunctionOP(context *Context, bytecode *Bytecode, code Code) (Value, *Object) {
	functionInformation := code.Value.([2]int)
	codeLength := functionInformation[0]
	numberOfArguments := functionInformation[1]
	start := bytecode.index
	bytecode.index += codeLength
	end := bytecode.index
	functionCode := make([]Code, codeLength)
	copy(functionCode, bytecode.instructions[start:end])
	return p.NewFunction(context, false, context.PeekSymbolTable(), NewPlasmaFunction(numberOfArguments, functionCode)), nil
}

// Useful function to call those built ins that doesn't receive arguments of an object
func (p *Plasma) noArgsGetAndCall(context *Context, operationName string) (Value, *Object) {
	object := context.PopObject()
	operation, getError := object.Get(operationName)
	if getError != nil {
		return nil, p.NewObjectWithNameNotFoundError(context, object.GetClass(p), operationName)
	}
	return p.CallFunction(context, operation, object.SymbolTable())
}

// Function useful to cal object built-ins binary expression functions
func (p *Plasma) leftBinaryExpressionFuncCall(context *Context, operationName string) (Value, *Object) {
	leftHandSide := context.PopObject()
	rightHandSide := context.PopObject()
	operation, getError := leftHandSide.Get(operationName)
	if getError != nil {
		return p.rightBinaryExpressionFuncCall(context, leftHandSide, rightHandSide, operationName)
	}
	result, callError := p.CallFunction(context, operation, leftHandSide.SymbolTable(), rightHandSide)
	if callError != nil {
		return p.rightBinaryExpressionFuncCall(context, leftHandSide, rightHandSide, operationName)
	}
	return result, nil
}

func (p *Plasma) rightBinaryExpressionFuncCall(context *Context, leftHandSide Value, rightHandSide Value, operationName string) (Value, *Object) {
	operation, getError := rightHandSide.Get("Right" + operationName)
	if getError != nil {
		return nil, p.NewObjectWithNameNotFoundError(context, p.ForceMasterGetAny(ObjectName).(*Type), "Right"+operationName)
	}
	return p.CallFunction(context, operation, rightHandSide.SymbolTable(), leftHandSide)
}

func (p *Plasma) indexOP(context *Context) (Value, *Object) {
	index := context.PopObject()
	source := context.PopObject()
	indexOperation, getError := source.Get(Index)
	if getError != nil {
		return nil, p.NewObjectWithNameNotFoundError(context, source.GetClass(p), Index)
	}
	return p.CallFunction(context, indexOperation, source.SymbolTable(), index)
}

func (p *Plasma) selectNameFromObjectOP(context *Context, code Code) (Value, *Object) {
	name := code.Value.(string)
	source := context.PopObject()
	value, getError := source.Get(name)
	if getError != nil {
		return nil, p.NewObjectWithNameNotFoundError(context, source.GetClass(p), name)
	}
	return value, nil
}

func (p *Plasma) methodInvocationOP(context *Context, code Code) (Value, *Object) {
	numberOfArguments := code.Value.(int)
	function := context.PopObject()
	var arguments []Value
	for i := 0; i < numberOfArguments; i++ {
		if !context.ObjectStack.HasNext() {
			return nil, p.NewInvalidNumberOfArgumentsError(context, i, numberOfArguments)
		}
		arguments = append(arguments, context.PopObject())
	}
	var result Value
	var callError *Object
	if _, ok := function.(*Type); ok {
		result, callError = p.ConstructObject(context, function.(*Type), NewSymbolTable(function.SymbolTable().Parent))
		if callError != nil {
			return nil, callError
		}
		resultInitialize, getError := result.Get(Initialize)
		if getError != nil {
			return nil, p.NewObjectWithNameNotFoundError(context, result.GetClass(p), Initialize)
		}
		_, callError = p.CallFunction(context, resultInitialize, result.SymbolTable(), arguments...)
	} else {
		result, callError = p.CallFunction(context, function, NewSymbolTable(function.SymbolTable().Parent), arguments...)
	}
	return result, callError
}

func (p *Plasma) getIdentifierOP(context *Context, code Code) (Value, *Object) {
	value, getError := context.PeekSymbolTable().GetAny(code.Value.(string))
	if getError != nil {
		return nil, p.NewObjectWithNameNotFoundError(context, p.ForceMasterGetAny(ObjectName).GetClass(p), code.Value.(string))
	}
	return value, nil
}

func (p *Plasma) newIteratorOP(context *Context, bytecode *Bytecode, code Code) (Value, *Object) {
	source := context.PopObject()
	var iterSource Value
	var callError *Object
	if _, ok := source.(*Iterator); ok {
		iterSource = source
	} else {
		iter, getError := source.Get(Iter)
		if getError != nil {
			return nil, p.NewObjectWithNameNotFoundError(context, source.GetClass(p), Iter)
		}
		iterSource, callError = p.CallFunction(context, iter, source.SymbolTable())
		if callError != nil {
			return nil, callError
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
		func() Value {
			return p.NewFunction(context, false, generatorIterator.symbols,
				NewPlasmaClassFunction(generatorIterator, 0, nextCode),
			)
		},
	)
	generatorIterator.SetOnDemandSymbol(HasNext,
		func() Value {
			return p.NewFunction(context, false, generatorIterator.symbols,
				NewPlasmaClassFunction(generatorIterator, 0, hasNextCode),
			)
		},
	)
	return generatorIterator, nil
}

func (p *Plasma) setupForLoopOP(context *Context, code Code) *Object {
	source := context.ObjectStack.Pop()
	sourceNext, nextGetError := source.Get(Next)
	sourceHasNext, hasNextGetError := source.Get(HasNext)
	var sourceIter Value
	if nextGetError == nil && hasNextGetError == nil {
		sourceIter = source
	} else {
		sourceToIter, getError := source.Get(Iter)
		if getError != nil {
			return p.NewObjectWithNameNotFoundError(context, source.GetClass(p), Iter)
		}
		var callError *Object
		sourceIter, callError = p.CallFunction(context, sourceToIter, context.PeekSymbolTable())
		if callError != nil {
			return callError
		}
		sourceNext, nextGetError = sourceIter.Get(Next)
		sourceHasNext, hasNextGetError = sourceIter.Get(HasNext)
		if nextGetError != nil {
			return p.NewObjectWithNameNotFoundError(context, sourceIter.GetClass(p), Next)
		} else if hasNextGetError != nil {
			return p.NewObjectWithNameNotFoundError(context, sourceIter.GetClass(p), HasNext)
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
			MappedReceivers:   map[string]Value{},
		},
	)
	return nil
}

func (p *Plasma) loadForLoopArguments(context *Context) *Object {
	for name, value := range context.LoopStack.Peek().MappedReceivers {
		context.PeekSymbolTable().Set(name, value)
	}
	return nil
}

func (p *Plasma) unpackForArguments(context *Context, loopSettings *LoopSettings, result Value) *Object {
	if loopSettings.NumberOfReceivers == 1 {
		loopSettings.MappedReceivers[loopSettings.Receivers[0]] = result
	} else if _, ok := result.(*String); ok {
		if result.GetLength() != loopSettings.NumberOfReceivers {
			return p.NewInvalidNumberOfArgumentsError(context, loopSettings.NumberOfReceivers, result.GetLength())
		}
		for index, name := range loopSettings.Receivers {
			loopSettings.MappedReceivers[name] = p.NewString(context, false, context.PeekSymbolTable(), string(result.GetString()[index]))
		}
	} else if _, ok = result.(*Bytes); ok {
		if result.GetLength() != loopSettings.NumberOfReceivers {
			return p.NewInvalidNumberOfArgumentsError(context, loopSettings.NumberOfReceivers, result.GetLength())
		}
		for index, name := range loopSettings.Receivers {
			loopSettings.MappedReceivers[name] = p.NewInteger(context, false, context.PeekSymbolTable(), int64(result.GetBytes()[index]))
		}
	} else if _, ok = result.(*Array); ok {
		if result.GetLength() != loopSettings.NumberOfReceivers {
			return p.NewInvalidNumberOfArgumentsError(context, loopSettings.NumberOfReceivers, result.GetLength())
		}
		for index, name := range loopSettings.Receivers {
			loopSettings.MappedReceivers[name] = result.GetContent()[index]
		}
	} else if _, ok = result.(*Tuple); ok {
		if result.GetLength() != loopSettings.NumberOfReceivers {
			return p.NewInvalidNumberOfArgumentsError(context, loopSettings.NumberOfReceivers, result.GetLength())
		}
		for index, name := range loopSettings.Receivers {
			loopSettings.MappedReceivers[name] = result.GetContent()[index]
		}
	} else if _, ok = result.(*HashTable); ok {
		return p.NewGoRuntimeError(context, errors.New("HashTable doesn't support unpacking"))
	} else {
		var hasNext Value
		var next Value
		var resultAsIterator Value
		var hasNextGetError *errors2.Error
		var nextGetError *errors2.Error
		hasNext, hasNextGetError = result.Get(HasNext)
		next, nextGetError = result.Get(Next)
		if hasNextGetError != nil || nextGetError != nil {
			resultToIter, getError := result.Get(Iter)
			if getError != nil {
				return p.NewObjectWithNameNotFoundError(context, result.GetClass(p), Iter)
			}
			var callError *Object
			resultAsIterator, callError = p.CallFunction(context, resultToIter, result.SymbolTable())
			if callError != nil {
				return callError
			}
		} else {
			resultAsIterator = result
		}
		for index, name := range loopSettings.Receivers {
			hasNextResult, callError := p.CallFunction(context, hasNext, resultAsIterator.SymbolTable())
			if callError != nil {
				return callError
			}
			var hasNextResultBool bool
			hasNextResultBool, callError = p.QuickGetBool(context, hasNextResult)
			if !hasNextResultBool {
				return p.NewInvalidNumberOfArgumentsError(context, index, loopSettings.NumberOfReceivers)
			}
			var nextResult Value
			nextResult, callError = p.CallFunction(context, next, context.PeekSymbolTable())
			if callError != nil {
				return callError
			}
			context.PeekSymbolTable().Set(name, nextResult)
		}
	}
	return nil
}

func (p *Plasma) unpackForLoopOP(context *Context, bytecode *Bytecode) *Object {
	loopSettings := context.LoopStack.Peek()
	hasNext, callError := p.CallFunction(context, loopSettings.HasNext, loopSettings.Source.SymbolTable())
	if callError != nil {
		return callError
	}
	var hasNextBool bool
	hasNextBool, callError = p.QuickGetBool(context, hasNext)
	if callError != nil {
		return callError
	}
	if !hasNextBool {
		// Do the jump to outside the loop
		bytecode.Jump(loopSettings.Jump)
		return nil
	}
	var nextValue Value
	nextValue, callError = p.CallFunction(context, loopSettings.Next, loopSettings.Source.SymbolTable())
	if callError != nil {
		return nil
	}
	return p.unpackForArguments(context, loopSettings, nextValue)
}

// Assign Statement

func (p *Plasma) assignIdentifierOP(context *Context, code Code) *Object {
	identifier := code.Value.(string)
	context.PeekSymbolTable().Set(identifier, context.PopObject())
	return nil
}

func (p *Plasma) assignSelectorOP(context *Context, code Code) *Object {
	target := context.PopObject()
	value := context.PopObject()
	identifier := code.Value.(string)
	if target.IsBuiltIn() {
		return p.NewBuiltInSymbolProtectionError(context, identifier)
	}
	target.Set(identifier, value)
	return nil
}

func (p *Plasma) assignIndexOP(context *Context) *Object {
	index := context.PopObject()
	source := context.PopObject()
	value := context.PopObject()
	sourceAssign, getError := source.Get(Assign)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(context, source.GetClass(p), Assign)
	}
	_, callError := p.CallFunction(context, sourceAssign, context.PeekSymbolTable(), index, value)
	if callError != nil {
		return callError
	}
	return nil
}

func (p *Plasma) returnOP(context *Context, code Code) *Object {
	numberOfReturnValues := code.Value.(int)
	if numberOfReturnValues == 0 {
		context.PushObject(p.GetNone())
		return nil
	}
	if numberOfReturnValues == 1 {
		if !context.ObjectStack.HasNext() {
			return p.NewInvalidNumberOfArgumentsError(context, 1, numberOfReturnValues)
		}
		return nil
	}

	var values []Value
	for i := 0; i < numberOfReturnValues; i++ {
		if !context.ObjectStack.HasNext() {
			return p.NewInvalidNumberOfArgumentsError(context, i, numberOfReturnValues)
		}
		values = append(values, context.PopObject())
	}
	context.PushObject(p.NewTuple(context, false, context.PeekSymbolTable(), values))
	return nil
}

// Special Instructions

func (p *Plasma) loadFunctionArgumentsOP(context *Context, code Code) *Object {
	for _, argument := range code.Value.([]string) {
		value := context.PopObject()
		context.PeekSymbolTable().Set(argument, value)
	}
	return nil
}

func (p *Plasma) newFunctionOP(context *Context, bytecode *Bytecode, code Code) *Object {
	functionInformation := code.Value.([2]int)
	codeLength := functionInformation[0]
	numberOfArguments := functionInformation[1]
	start := bytecode.index
	bytecode.index += codeLength
	end := bytecode.index
	functionCode := make([]Code, codeLength)
	copy(functionCode, bytecode.instructions[start:end])
	context.PushObject(p.NewFunction(context, false, context.PeekSymbolTable(), NewPlasmaFunction(numberOfArguments, functionCode)))
	return nil
}

func (p *Plasma) jumpOP(bytecode *Bytecode, code Code) *Object {
	bytecode.index += code.Value.(int)
	return nil
}

func (p *Plasma) ifJumpOP(context *Context, bytecode *Bytecode, code Code) *Object {
	condition := context.ObjectStack.Pop()
	conditionBool, executionError := p.QuickGetBool(context, condition)
	if executionError != nil {
		return executionError
	}
	if !conditionBool {
		bytecode.index += code.Value.(int)
	}
	return nil
}

func (p *Plasma) unlessJumpOP(context *Context, bytecode *Bytecode, code Code) *Object {
	condition := context.ObjectStack.Pop()
	conditionBool, executionError := p.QuickGetBool(context, condition)
	if executionError != nil {
		return executionError
	}
	if conditionBool {
		bytecode.index += code.Value.(int)
	}
	return nil
}

type ModuleInformation struct {
	Name       string
	CodeLength int
}

func (p *Plasma) newModuleOP(context *Context, bytecode *Bytecode, code Code) *Object {
	moduleInformation := code.Value.(ModuleInformation)
	var moduleBody []Code
	for i := 0; i < moduleInformation.CodeLength; i++ {
		moduleBody = append(moduleBody, bytecode.Next())
	}
	module := p.NewModule(context, false, context.PeekSymbolTable())
	context.PushSymbolTable(module.SymbolTable())
	_, executionError := p.Execute(context, NewBytecodeFromArray(moduleBody))
	if executionError != nil {
		return executionError
	}
	context.PopSymbolTable()
	context.PeekSymbolTable().Set(moduleInformation.Name, module)
	return nil
}

type ClassInformation struct {
	Name       string
	BodyLength int
}

func (p *Plasma) newClassOP(context *Context, bytecode *Bytecode, code Code) *Object {
	classInformation := code.Value.(ClassInformation)
	rawSubClasses := context.PopObject().GetContent()
	var subClasses []*Type
	for _, subClass := range rawSubClasses {
		if _, ok := subClass.(*Type); !ok {
			return p.NewInvalidTypeError(context, subClass.TypeName(), TypeName)
		}
		subClasses = append(subClasses, subClass.(*Type))
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
	return nil
}

func (p *Plasma) newClassFunctionOP(context *Context, bytecode *Bytecode, code Code) *Object {
	functionInformation := code.Value.([2]int)
	codeLength := functionInformation[0]
	numberOfArguments := functionInformation[1]
	start := bytecode.index
	bytecode.index += codeLength
	end := bytecode.index
	functionCode := make([]Code, codeLength)
	copy(functionCode, bytecode.instructions[start:end])
	context.PushObject(p.NewFunction(context, false, context.PeekSymbolTable(), NewPlasmaClassFunction(context.PeekObject(), numberOfArguments, functionCode)))
	return nil
}

func (p *Plasma) raiseOP(context *Context) *Object {
	if _, ok := context.PeekObject().(*Object); !ok {
		return p.NewInvalidTypeError(context, context.PeekObject().TypeName(), RuntimeError)
	}
	if !context.PeekObject().Implements(p.ForceMasterGetAny(RuntimeError).(*Type)) {
		return p.NewInvalidTypeError(context, context.PeekObject().TypeName(), RuntimeError)
	}
	return context.PeekObject().(*Object)
}

func (p *Plasma) caseOP(context *Context, bytecode *Bytecode, code Code) *Object {
	references := context.PopObject()
	contains := p.ForceGetSelf(Contains, references)
	result, callError := p.CallFunction(context, contains, references.SymbolTable(), context.PeekObject())
	if callError != nil {
		return callError
	}
	boolResult, executionError := p.QuickGetBool(context, result)
	if executionError != nil {
		return executionError
	}
	if !boolResult {
		bytecode.index += code.Value.(int)
		return nil
	}
	context.PopObject()
	return nil
}
