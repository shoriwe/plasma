package vm

import (
	"fmt"
)

func (p *Plasma) Execute(context *Context, bytecode *Bytecode) (Value, *Object) {
	if context == nil {
		context = NewContext()
		p.InitializeContext(context)
	}
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
			if context.MemoryStack.head != nil {
				fmt.Println("Head:", context.MemoryStack.head.value)
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
		case IfOP:
			executionError = p.ifOP(context, code)
		case IfOneLinerOP:
			object, executionError = p.ifOneLinerOP(context, code)
		case UnlessOP:
			executionError = p.unlessOP(context, code)
		case UnlessOneLinerOP:
			object, executionError = p.unlessOneLinerOP(context, code)
		// Special Instructions
		case LoadFunctionArgumentsOP:
			executionError = p.loadFunctionArgumentsOP(context, code)
		case NewFunctionOP:
			executionError = p.newFunctionOP(context, bytecode, code)
		case JumpOP:
			executionError = p.jumpOP(bytecode, code)
		case RedoOP:
			// executionError = p.jumpOP(bytecode, code)
			context.StateStack.Peek().Action = Redo
			return nil, nil
		case BreakOP:
			// executionError = p.jumpOP(bytecode, code)
			context.StateStack.Peek().Action = Break
			return nil, nil
		case ContinueOP:
			// executionError = p.jumpOP(bytecode, code)
			return nil, nil
		case PushOP:
			if object != nil {
				context.MemoryStack.Push(object)
				object = nil
			}
		case PopOP:
			context.MemoryStack.Pop()
		case NOP:
			break
		case DoWhileLoop:
			executionError = p.setupDoWhileLoop(context, bytecode, code)
		case WhileLoop:
			executionError = p.setupWhileLoop(context, bytecode, code)
		case ForLoopOP:
			executionError = p.setupForLoopOP(context, bytecode, code)
		case TryOP:
			executionError = p.tryOP(context, code)
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
			return nil, executionError
		}
	}
	if context.MemoryStack.HasNext() {
		return context.PopObject(), nil
	}
	return p.GetNone(), nil
}

func (p *Plasma) newStringOP(context *Context, code Code) (Value, *Object) {
	value := code.Value.(string)
	stringObject := p.NewString(context, false, context.SymbolTableStack.Peek(), value)
	return stringObject, nil
}

func (p *Plasma) newBytesOP(context *Context, code Code) (Value, *Object) {
	value := code.Value.([]byte)
	return p.NewBytes(context, false, context.SymbolTableStack.Peek(), value), nil
}

func (p *Plasma) newIntegerOP(context *Context, code Code) (Value, *Object) {
	value := code.Value.(int64)
	return p.NewInteger(context, false, context.SymbolTableStack.Peek(), value), nil
}

func (p *Plasma) newFloatOP(context *Context, code Code) (Value, *Object) {
	value := code.Value.(float64)
	return p.NewFloat(context, false, context.SymbolTableStack.Peek(), value), nil
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
		if !context.MemoryStack.HasNext() {
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
	generatorIterator.Set(Next,
		p.NewFunction(context, false, generatorIterator.symbols,
			NewPlasmaClassFunction(generatorIterator, 0, nextCode),
		),
	)
	generatorIterator.Set(HasNext,
		p.NewFunction(context, false, generatorIterator.symbols,
			NewPlasmaClassFunction(generatorIterator, 0, hasNextCode),
		),
	)
	return generatorIterator, nil
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
		if !context.MemoryStack.HasNext() {
			return p.NewInvalidNumberOfArgumentsError(context, 1, numberOfReturnValues)
		}
		return nil
	}

	var values []Value
	for i := 0; i < numberOfReturnValues; i++ {
		if !context.MemoryStack.HasNext() {
			return p.NewInvalidNumberOfArgumentsError(context, i, numberOfReturnValues)
		}
		values = append(values, context.PopObject())
	}
	context.PushObject(p.NewTuple(context, false, context.PeekSymbolTable(), values))
	return nil
}

func (p *Plasma) ifOP(context *Context, code Code) *Object {
	ifInformation := code.Value.(*IfInformation)
	condition, executionError := p.Execute(context, NewBytecodeFromArray(ifInformation.Condition))
	if executionError != nil {
		return executionError
	}
	var conditionBool bool
	conditionBool, executionError = p.QuickGetBool(context, condition)
	if conditionBool {
		_, executionError = p.Execute(context, NewBytecodeFromArray(ifInformation.Body))
		return executionError
	}
	for _, elif := range ifInformation.ElifBlocks {
		condition, executionError = p.Execute(context, NewBytecodeFromArray(elif.Condition))
		if executionError != nil {
			return executionError
		}
		conditionBool, executionError = p.QuickGetBool(context, condition)
		if conditionBool {
			_, executionError = p.Execute(context, NewBytecodeFromArray(elif.Body))
			return executionError
		}
	}
	_, executionError = p.Execute(context, NewBytecodeFromArray(ifInformation.Else))
	return executionError
}

func (p *Plasma) ifOneLinerOP(context *Context, code Code) (Value, *Object) {
	ifInformation := code.Value.(*IfInformation)
	condition, executionError := p.Execute(context, NewBytecodeFromArray(ifInformation.Condition))
	if executionError != nil {
		return nil, executionError
	}
	var conditionBool bool
	conditionBool, executionError = p.QuickGetBool(context, condition)
	if conditionBool {
		return p.Execute(context, NewBytecodeFromArray(ifInformation.Body))
	} else if ifInformation.Else != nil {
		return p.Execute(context, NewBytecodeFromArray(ifInformation.Else))
	}
	return p.GetNone(), nil
}

func (p *Plasma) unlessOP(context *Context, code Code) *Object {
	ifInformation := code.Value.(*IfInformation)
	condition, executionError := p.Execute(context, NewBytecodeFromArray(ifInformation.Condition))
	if executionError != nil {
		return executionError
	}
	var conditionBool bool
	conditionBool, executionError = p.QuickGetBool(context, condition)
	if !conditionBool {
		_, executionError = p.Execute(context, NewBytecodeFromArray(ifInformation.Body))
		return executionError
	}
	for _, elif := range ifInformation.ElifBlocks {
		condition, executionError = p.Execute(context, NewBytecodeFromArray(elif.Condition))
		if executionError != nil {
			return executionError
		}
		conditionBool, executionError = p.QuickGetBool(context, condition)
		if !conditionBool {
			_, executionError = p.Execute(context, NewBytecodeFromArray(elif.Body))
			return executionError
		}
	}
	_, executionError = p.Execute(context, NewBytecodeFromArray(ifInformation.Else))
	return executionError
}

func (p *Plasma) unlessOneLinerOP(context *Context, code Code) (Value, *Object) {
	ifInformation := code.Value.(*IfInformation)
	condition, executionError := p.Execute(context, NewBytecodeFromArray(ifInformation.Condition))
	if executionError != nil {
		return nil, executionError
	}
	var conditionBool bool
	conditionBool, executionError = p.QuickGetBool(context, condition)
	if !conditionBool {
		return p.Execute(context, NewBytecodeFromArray(ifInformation.Body))
	} else if ifInformation.Else != nil {
		return p.Execute(context, NewBytecodeFromArray(ifInformation.Else))
	}
	return p.GetNone(), nil
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

func (p *Plasma) setupDoWhileLoop(context *Context, bytecode *Bytecode, code Code) *Object {
	condition := NewBytecodeFromArray(bytecode.NextN(code.Value.([2]int)[0]))
	body := NewBytecodeFromArray(bytecode.NextN(code.Value.([2]int)[1]))
	doWhileLoopEntry := &stateEntry{
		Action: NoAction,
	}
	context.StateStack.Push(
		doWhileLoopEntry,
	)
	defer context.StateStack.Pop()
loop:
	for {
	redoLocation:
		// Execute the body
		_, executionError := p.Execute(context, body)
		body.index = 0
		if executionError != nil {
			return executionError
		}
		// Check continue, redo and break
		switch doWhileLoopEntry.Action {
		case Break:
			break loop
		case Redo:
			doWhileLoopEntry.Action = NoAction
			goto redoLocation
		}

		// Evaluate the condition
		var result Value
		result, executionError = p.Execute(context, condition)
		condition.index = 0
		if executionError != nil {
			return executionError
		}
		var conditionBool bool
		conditionBool, executionError = p.QuickGetBool(context, result)
		if !conditionBool {
			break
		}
	}
	return nil
}

func (p *Plasma) setupWhileLoop(context *Context, bytecode *Bytecode, code Code) *Object {
	condition := NewBytecodeFromArray(bytecode.NextN(code.Value.([2]int)[0]))
	body := NewBytecodeFromArray(bytecode.NextN(code.Value.([2]int)[1]))
	whileLoopEntry := &stateEntry{
		Action: NoAction,
	}
	context.StateStack.Push(
		whileLoopEntry,
	)
	defer context.StateStack.Pop()
loop:
	for {
		// First Evaluate the condition
		result, executionError := p.Execute(context, condition)
		condition.index = 0
		if executionError != nil {
			return executionError
		}

		var conditionBool bool
		conditionBool, executionError = p.QuickGetBool(context, result)
		if !conditionBool {
			break
		}
	redoLocation:
		// Execute the body
		_, executionError = p.Execute(context, body)
		body.index = 0
		if executionError != nil {
			return executionError
		}
		switch whileLoopEntry.Action {
		case Break:
			break loop
		case Redo:
			whileLoopEntry.Action = NoAction
			goto redoLocation
		}
	}
	return nil
}

func (p *Plasma) reloadForLoopContext(context *Context, loopContext *map[string]Value, numberOfReceivers int, receivers []string, sourceHasNext Value, sourceNext Value) (bool, *Object) {
	hasNextObject, callError := p.CallFunction(context, sourceHasNext, sourceHasNext.SymbolTable())
	if callError != nil {
		return false, callError
	}
	hasNextObjectBool, executionError := p.QuickGetBool(context, hasNextObject)
	if executionError != nil {
		return false, executionError
	}
	if !hasNextObjectBool {
		return false, nil
	}
	var value Value
	value, callError = p.CallFunction(context, sourceNext, sourceNext.SymbolTable())
	if callError != nil {
		return false, callError
	}
	if numberOfReceivers == 1 {
		(*loopContext)[receivers[0]] = value
		return true, nil
	}
	// Unpack it as first calling to iter
	valueIterFunc, getError := value.Get(Iter)
	if getError != nil {
		return false, p.NewObjectWithNameNotFoundError(context, value.GetClass(p), Iter)
	}
	var valueAsIter Value
	valueAsIter, callError = p.CallFunction(context, valueIterFunc, valueIterFunc.SymbolTable())
	if callError != nil {
		return false, callError
	}
	var valueAsIterHasNext Value
	valueAsIterHasNext, getError = valueAsIter.Get(HasNext)
	if getError != nil {
		return false, p.NewObjectWithNameNotFoundError(context, valueAsIter.GetClass(p), HasNext)
	}
	var valueAsIterNext Value
	valueAsIterNext, getError = valueAsIter.Get(Next)
	if getError != nil {
		return false, p.NewObjectWithNameNotFoundError(context, valueAsIter.GetClass(p), Next)
	}
	for index, receiver := range receivers {
		hasNextObject, callError = p.CallFunction(context, valueAsIterHasNext, valueAsIterHasNext.SymbolTable())
		if callError != nil {
			return false, callError
		}
		hasNextObjectBool, executionError = p.QuickGetBool(context, hasNextObject)
		if executionError != nil {
			return false, executionError
		}
		if !hasNextObjectBool {
			return false, p.NewInvalidNumberOfArgumentsError(context, numberOfReceivers, index+1)
		}
		value, callError = p.CallFunction(context, valueAsIterNext, valueAsIterNext.SymbolTable())
		if callError != nil {
			return false, callError
		}
		(*loopContext)[receiver] = value
	}
	return true, nil
}

func (p *Plasma) setupForLoopOP(context *Context, bytecode *Bytecode, code Code) *Object {
	source := context.PopObject()
	sourceIter, getError := source.Get(Iter)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(context, source.GetClass(p), Iter)
	}
	sourceAsIter, callError := p.CallFunction(context, sourceIter, sourceIter.SymbolTable())
	if callError != nil {
		return callError
	}
	loopSettings := code.Value.(ForLoopSettings)

	bodyBytecode := NewBytecodeFromArray(bytecode.NextN(loopSettings.BodyLength))
	receivers := loopSettings.Receivers
	forLoopEntry := &stateEntry{
		Action: NoAction,
	}
	context.StateStack.Push(
		forLoopEntry,
	)
	defer context.StateStack.Pop()
	loopContext := map[string]Value{}
	var sourceHasNext Value
	sourceHasNext, getError = sourceAsIter.Get(HasNext)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(context, sourceAsIter.GetClass(p), HasNext)
	}
	var sourceNext Value
	sourceNext, getError = sourceAsIter.Get(Next)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(context, sourceAsIter.GetClass(p), Next)
	}
	receiversLength := len(receivers)
loop:
	for {
		// Update receivers
		// Check if the iteration can continue
		hasNext, loadSymbolsError := p.reloadForLoopContext(context, &loopContext, receiversLength, receivers, sourceHasNext, sourceNext)
		if loadSymbolsError != nil {
			return loadSymbolsError
		}
		if !hasNext {
			break
		}
	redoLocation:
		// Load the receivers
		for receiver, object := range loopContext {
			context.PeekSymbolTable().Set(receiver, object)
		}
		// Execute body
		_, bodyExecutionError := p.Execute(context, bodyBytecode)
		bodyBytecode.index = 0
		// If fail return return error
		if bodyExecutionError != nil {
			return bodyExecutionError
		}
		// Check continue, redo and break
		switch forLoopEntry.Action {
		case Break:
			break loop
		case Redo:
			forLoopEntry.Action = NoAction
			goto redoLocation
		}
	}
	return nil
}

func (p *Plasma) executeFinally(context *Context, finally []Code) *Object {
	if finally != nil {
		_, executionError := p.Execute(context, NewBytecodeFromArray(finally))
		return executionError
	}
	return nil
}

func (p *Plasma) tryOP(context *Context, code Code) *Object {
	tryInformation := code.Value.(*TryInformation)
	_, executionError := p.Execute(context, NewBytecodeFromArray(tryInformation.Body))
	if executionError == nil {
		return p.executeFinally(context, tryInformation.Finally)
	}
	var targetError Value
	var executionError2 *Object
	for _, exceptBlock := range tryInformation.ExceptBlocks {
		if exceptBlock.TargetErrors != nil {
			for _, targetErrorCode := range exceptBlock.TargetErrors {
				targetError, executionError2 = p.Execute(context, NewBytecodeFromArray(targetErrorCode))
				if executionError2 != nil {
					return executionError
				}
				if !targetError.Implements(p.ForceMasterGetAny(RuntimeError).(*Type)) {
					return p.NewInvalidTypeError(context, targetError.TypeName(), RuntimeError)
				}
				if executionError.Implements(targetError.(*Type)) {
					context.PeekSymbolTable().Set(exceptBlock.Receiver, executionError)
					_, executionError2 = p.Execute(context, NewBytecodeFromArray(exceptBlock.Body))
					if executionError2 == nil {
						return p.executeFinally(context, tryInformation.Finally)
					}
					return executionError2
				}
			}
		} else {
			context.PeekSymbolTable().Set(exceptBlock.Receiver, executionError)
			_, executionError2 = p.Execute(context, NewBytecodeFromArray(exceptBlock.Body))
			if executionError2 == nil {
				return p.executeFinally(context, tryInformation.Finally)
			}
			return executionError2
		}
	}
	if tryInformation.Else != nil {
		_, executionError2 = p.Execute(context, NewBytecodeFromArray(tryInformation.Else))
		if executionError2 != nil {
			return executionError2
		}
	}
	return p.executeFinally(context, tryInformation.Finally)
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
