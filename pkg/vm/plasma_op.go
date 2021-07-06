package vm

import (
	"fmt"
)

func (p *Plasma) Execute(bytecode *Bytecode) (Value, *Object) {
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
			if p.MemoryStack.head != nil {
				fmt.Println("Head:", p.MemoryStack.head.value)
			}
			fmt.Println("Object:", object)
		*/

		switch code.Instruction.OpCode {
		// Literals
		case NewStringOP:
			object, executionError = p.newStringOP(code)
		case NewBytesOP:
			object, executionError = p.newBytesOP(code)
		case NewIntegerOP:
			object, executionError = p.newIntegerOP(code)
		case NewFloatOP:
			object, executionError = p.newFloatOP(code)
		case NewTrueBoolOP:
			object, executionError = p.newTrueBoolOP()
		case NewFalseBoolOP:
			object, executionError = p.newFalseBoolOP()
		case NewParenthesesOP:
			object, executionError = p.newParenthesesOP()
		case NewLambdaFunctionOP:
			object, executionError = p.newLambdaFunctionOP(bytecode, code)
		case GetNoneOP:
			object, executionError = p.getNoneOP()
		// Composite creation
		case NewTupleOP:
			object, executionError = p.newTupleOP(code)
		case NewArrayOP:
			object, executionError = p.newArrayOP(code)
		case NewHashOP:
			object, executionError = p.newHashOP(code)
		// Unary Expressions
		case NegateBitsOP:
			object, executionError = p.noArgsGetAndCall(NegBits)
		case BoolNegateOP:
			object, executionError = p.noArgsGetAndCall(Negate)
		case NegativeOP:
			object, executionError = p.noArgsGetAndCall(Negative)
		// Binary Expressions
		case AddOP:
			object, executionError = p.leftBinaryExpressionFuncCall(Add)
		case SubOP:
			object, executionError = p.leftBinaryExpressionFuncCall(Sub)
		case MulOP:
			object, executionError = p.leftBinaryExpressionFuncCall(Mul)
		case DivOP:
			object, executionError = p.leftBinaryExpressionFuncCall(Div)
		case FloorDivOP:
			object, executionError = p.leftBinaryExpressionFuncCall(FloorDiv)
		case ModOP:
			object, executionError = p.leftBinaryExpressionFuncCall(Mod)
		case PowOP:
			object, executionError = p.leftBinaryExpressionFuncCall(Pow)
		case BitXorOP:
			object, executionError = p.leftBinaryExpressionFuncCall(BitXor)
		case BitAndOP:
			object, executionError = p.leftBinaryExpressionFuncCall(BitAnd)
		case BitOrOP:
			object, executionError = p.leftBinaryExpressionFuncCall(BitOr)
		case BitLeftOP:
			object, executionError = p.leftBinaryExpressionFuncCall(BitLeft)
		case BitRightOP:
			object, executionError = p.leftBinaryExpressionFuncCall(BitRight)
		case AndOP:
			object, executionError = p.leftBinaryExpressionFuncCall(And)
		case OrOP:
			object, executionError = p.leftBinaryExpressionFuncCall(Or)
		case XorOP:
			object, executionError = p.leftBinaryExpressionFuncCall(Xor)
		case EqualsOP:
			object, executionError = p.leftBinaryExpressionFuncCall(Equals)
		case NotEqualsOP:
			object, executionError = p.leftBinaryExpressionFuncCall(NotEquals)
		case GreaterThanOP:
			object, executionError = p.leftBinaryExpressionFuncCall(GreaterThan)
		case LessThanOP:
			object, executionError = p.leftBinaryExpressionFuncCall(LessThan)
		case GreaterThanOrEqualOP:
			object, executionError = p.leftBinaryExpressionFuncCall(GreaterThanOrEqual)
		case LessThanOrEqualOP:
			object, executionError = p.leftBinaryExpressionFuncCall(LessThanOrEqual)
		case ContainsOP:
			// This operation is inverted, right is left and left is right
			leftHandSide := p.PopObject()
			rightHandSide := p.PopObject()
			p.PushObject(leftHandSide)
			p.PushObject(rightHandSide)
			object, executionError = p.leftBinaryExpressionFuncCall(Contains)
		// Other Expressions
		case GetIdentifierOP:
			object, executionError = p.getIdentifierOP(code)
		case IndexOP:
			object, executionError = p.indexOP()
		case SelectNameFromObjectOP:
			object, executionError = p.selectNameFromObjectOP(code)
		case MethodInvocationOP:
			object, executionError = p.methodInvocationOP(code)
		case NewIteratorOP:
			object, executionError = p.newIteratorOP(bytecode, code)
		// Assign Statement
		case AssignIdentifierOP:
			executionError = p.assignIdentifierOP(code)
		case AssignSelectorOP:
			executionError = p.assignSelectorOP(code)
		case AssignIndexOP:
			executionError = p.assignIndexOP()
		case ReturnOP:
			executionError = p.returnOP(code)
			break bytecodeExecutionLoop
		case IfOP:
			executionError = p.ifOP(code)
		case IfOneLinerOP:
			object, executionError = p.ifOneLinerOP(code)
		case UnlessOP:
			executionError = p.unlessOP(code)
		case UnlessOneLinerOP:
			object, executionError = p.unlessOneLinerOP(code)
		// Special Instructions
		case LoadFunctionArgumentsOP:
			executionError = p.loadFunctionArgumentsOP(code)
		case NewFunctionOP:
			executionError = p.newFunctionOP(bytecode, code)
		case JumpOP:
			executionError = p.jumpOP(bytecode, code)
		case RedoOP:
			executionError = p.jumpOP(bytecode, code)
		case BreakOP:
			executionError = p.jumpOP(bytecode, code)
		case ContinueOP:
			executionError = p.jumpOP(bytecode, code)
		case PushOP:
			if object != nil {
				p.MemoryStack.Push(object)
				object = nil
			}
		case PopOP:
			p.MemoryStack.Pop()
		case NOP:
			break
		case DoWhileLoop:
			executionError = p.setupDoWhileLoop(bytecode, code)
		case WhileLoop:
			executionError = p.setupWhileLoop(bytecode, code)
		case ForLoopOP:
			executionError = p.setupForLoopOP(bytecode, code)
		case TryOP:
			executionError = p.tryOP(code)
		case NewModuleOP:
			executionError = p.newModuleOP(bytecode, code)
		case RaiseOP:
			executionError = p.raiseOP()
		case NewClassOP:
			executionError = p.newClassOP(bytecode, code)
		case NewClassFunctionOP:
			executionError = p.newClassFunctionOP(bytecode, code)
		case CaseOP:
			executionError = p.caseOP(bytecode, code)
		default:
			panic(fmt.Sprintf("Unknown VM instruction %d", code.Instruction.OpCode))
		}
		if executionError != nil {
			return nil, executionError
		}
	}
	if p.MemoryStack.HasNext() {
		return p.PopObject(), nil
	}
	return p.GetNone(), nil
}

func (p *Plasma) newStringOP(code Code) (Value, *Object) {
	value := code.Value.(string)
	stringObject := p.NewString(false, p.SymbolTableStack.Peek(), value)
	return stringObject, nil
}

func (p *Plasma) newBytesOP(code Code) (Value, *Object) {
	value := code.Value.([]byte)
	return p.NewBytes(false, p.SymbolTableStack.Peek(), value), nil
}

func (p *Plasma) newIntegerOP(code Code) (Value, *Object) {
	value := code.Value.(int64)
	return p.NewInteger(false, p.SymbolTableStack.Peek(), value), nil
}

func (p *Plasma) newFloatOP(code Code) (Value, *Object) {
	value := code.Value.(float64)
	return p.NewFloat(false, p.SymbolTableStack.Peek(), value), nil
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

func (p *Plasma) newTupleOP(code Code) (Value, *Object) {
	numberOfValues := code.Value.(int)
	var values []Value
	for i := 0; i < numberOfValues; i++ {
		values = append(values, p.PopObject())
	}
	return p.NewTuple(false, p.PeekSymbolTable(), values), nil
}

func (p *Plasma) newArrayOP(code Code) (Value, *Object) {
	numberOfValues := code.Value.(int)
	var values []Value
	for i := 0; i < numberOfValues; i++ {
		values = append(values, p.PopObject())
	}
	return p.NewArray(false, p.PeekSymbolTable(), values), nil
}

func (p *Plasma) newHashOP(code Code) (Value, *Object) {
	numberOfValues := code.Value.(int)
	var keyValues []*KeyValue
	for i := 0; i < numberOfValues; i++ {

		key := p.PopObject()
		value := p.PopObject()
		keyValues = append(keyValues, &KeyValue{
			Key:   key,
			Value: value,
		})
	}
	hashTable := p.NewHashTable(false, p.PeekSymbolTable(), map[int64][]*KeyValue{}, numberOfValues)
	hashTableAssign, getError := hashTable.Get(Assign)
	if getError != nil {
		return nil, p.NewObjectWithNameNotFoundError(hashTable.GetClass(p), Assign)
	}
	for _, keyValue := range keyValues {
		_, assignError := p.CallFunction(hashTableAssign, hashTable.SymbolTable(), keyValue.Key, keyValue.Value)
		if assignError != nil {
			return nil, assignError
		}
	}
	return hashTable, nil
}

func (p *Plasma) newParenthesesOP() (Value, *Object) {
	return p.PopObject(), nil
}

func (p *Plasma) newLambdaFunctionOP(bytecode *Bytecode, code Code) (Value, *Object) {
	functionInformation := code.Value.([2]int)
	codeLength := functionInformation[0]
	numberOfArguments := functionInformation[1]
	start := bytecode.index
	bytecode.index += codeLength
	end := bytecode.index
	functionCode := make([]Code, codeLength)
	copy(functionCode, bytecode.instructions[start:end])
	return p.NewFunction(false, p.PeekSymbolTable(), NewPlasmaFunction(numberOfArguments, functionCode)), nil
}

// Useful function to call those built ins that doesn't receive arguments of an object
func (p *Plasma) noArgsGetAndCall(operationName string) (Value, *Object) {
	object := p.PopObject()
	operation, getError := object.Get(operationName)
	if getError != nil {
		return nil, p.NewObjectWithNameNotFoundError(object.GetClass(p), operationName)
	}
	return p.CallFunction(operation, object.SymbolTable())
}

// Function useful to cal object built-ins binary expression functions
func (p *Plasma) leftBinaryExpressionFuncCall(operationName string) (Value, *Object) {
	leftHandSide := p.PopObject()
	rightHandSide := p.PopObject()
	operation, getError := leftHandSide.Get(operationName)
	if getError != nil {
		return p.rightBinaryExpressionFuncCall(leftHandSide, rightHandSide, operationName)
	}
	result, callError := p.CallFunction(operation, leftHandSide.SymbolTable(), rightHandSide)
	if callError != nil {
		return p.rightBinaryExpressionFuncCall(leftHandSide, rightHandSide, operationName)
	}
	return result, nil
}

func (p *Plasma) rightBinaryExpressionFuncCall(leftHandSide Value, rightHandSide Value, operationName string) (Value, *Object) {
	operation, getError := rightHandSide.Get("Right" + operationName)
	if getError != nil {
		return nil, p.NewObjectWithNameNotFoundError(p.ForceMasterGetAny(ObjectName).(*Type), "Right"+operationName)
	}
	return p.CallFunction(operation, rightHandSide.SymbolTable(), leftHandSide)
}

func (p *Plasma) indexOP() (Value, *Object) {
	index := p.PopObject()
	source := p.PopObject()
	indexOperation, getError := source.Get(Index)
	if getError != nil {
		return nil, p.NewObjectWithNameNotFoundError(source.GetClass(p), Index)
	}
	return p.CallFunction(indexOperation, source.SymbolTable(), index)
}

func (p *Plasma) selectNameFromObjectOP(code Code) (Value, *Object) {
	name := code.Value.(string)
	source := p.PopObject()
	value, getError := source.Get(name)
	if getError != nil {
		return nil, p.NewObjectWithNameNotFoundError(source.GetClass(p), name)
	}
	return value, nil
}

func (p *Plasma) methodInvocationOP(code Code) (Value, *Object) {
	numberOfArguments := code.Value.(int)
	function := p.PopObject()
	var arguments []Value
	for i := 0; i < numberOfArguments; i++ {
		if !p.MemoryStack.HasNext() {
			return nil, p.NewInvalidNumberOfArgumentsError(i, numberOfArguments)
		}
		arguments = append(arguments, p.PopObject())
	}
	var result Value
	var callError *Object
	if _, ok := function.(*Type); ok {
		result, callError = p.ConstructObject(function.(*Type), NewSymbolTable(function.SymbolTable().Parent))
		if callError != nil {
			return nil, callError
		}
		resultInitialize, getError := result.Get(Initialize)
		if getError != nil {
			return nil, p.NewObjectWithNameNotFoundError(result.GetClass(p), Initialize)
		}
		_, callError = p.CallFunction(resultInitialize, result.SymbolTable(), arguments...)
	} else {
		result, callError = p.CallFunction(function, NewSymbolTable(function.SymbolTable().Parent), arguments...)
	}
	return result, callError
}

func (p *Plasma) getIdentifierOP(code Code) (Value, *Object) {
	value, getError := p.PeekSymbolTable().GetAny(code.Value.(string))
	if getError != nil {
		return nil, p.NewObjectWithNameNotFoundError(p.ForceMasterGetAny(ObjectName).GetClass(p), code.Value.(string))
	}
	return value, nil
}

func (p *Plasma) newIteratorOP(bytecode *Bytecode, code Code) (Value, *Object) {
	source := p.PopObject()
	var iterSource Value
	var callError *Object
	if _, ok := source.(*Iterator); ok {
		iterSource = source
	} else {
		iter, getError := source.Get(Iter)
		if getError != nil {
			return nil, p.NewObjectWithNameNotFoundError(source.GetClass(p), Iter)
		}
		iterSource, callError = p.CallFunction(iter, source.SymbolTable())
		if callError != nil {
			return nil, callError
		}
	}
	generatorIterator := p.NewIterator(false, p.PeekSymbolTable())
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
		p.NewFunction(false, generatorIterator.symbols,
			NewPlasmaClassFunction(generatorIterator, 0, nextCode),
		),
	)
	generatorIterator.Set(HasNext,
		p.NewFunction(false, generatorIterator.symbols,
			NewPlasmaClassFunction(generatorIterator, 0, hasNextCode),
		),
	)
	return generatorIterator, nil
}

// Assign Statement

func (p *Plasma) assignIdentifierOP(code Code) *Object {
	identifier := code.Value.(string)
	p.PeekSymbolTable().Set(identifier, p.PopObject())
	return nil
}

func (p *Plasma) assignSelectorOP(code Code) *Object {
	target := p.PopObject()
	value := p.PopObject()
	identifier := code.Value.(string)
	if target.IsBuiltIn() {
		return p.NewBuiltInSymbolProtectionError(identifier)
	}
	target.Set(identifier, value)
	return nil
}

func (p *Plasma) assignIndexOP() *Object {
	index := p.PopObject()
	source := p.PopObject()
	value := p.PopObject()
	sourceAssign, getError := source.Get(Assign)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(source.GetClass(p), Assign)
	}
	_, callError := p.CallFunction(sourceAssign, p.PeekSymbolTable(), index, value)
	if callError != nil {
		return callError
	}
	return nil
}

func (p *Plasma) returnOP(code Code) *Object {
	numberOfReturnValues := code.Value.(int)
	if numberOfReturnValues == 0 {
		p.PushObject(p.GetNone())
		return nil
	}
	if numberOfReturnValues == 1 {
		if !p.MemoryStack.HasNext() {
			return p.NewInvalidNumberOfArgumentsError(1, numberOfReturnValues)
		}
		return nil
	}

	var values []Value
	for i := 0; i < numberOfReturnValues; i++ {
		if !p.MemoryStack.HasNext() {
			return p.NewInvalidNumberOfArgumentsError(i, numberOfReturnValues)
		}
		values = append(values, p.PopObject())
	}
	p.PushObject(p.NewTuple(false, p.PeekSymbolTable(), values))
	return nil
}

func (p *Plasma) ifOP(code Code) *Object {
	ifInformation := code.Value.(*IfInformation)
	condition, executionError := p.Execute(NewBytecodeFromArray(ifInformation.Condition))
	if executionError != nil {
		return executionError
	}
	var conditionBool bool
	conditionBool, executionError = p.QuickGetBool(condition)
	if conditionBool {
		_, executionError = p.Execute(NewBytecodeFromArray(ifInformation.Body))
		return executionError
	}
	for _, elif := range ifInformation.ElifBlocks {
		condition, executionError = p.Execute(NewBytecodeFromArray(elif.Condition))
		if executionError != nil {
			return executionError
		}
		conditionBool, executionError = p.QuickGetBool(condition)
		if conditionBool {
			_, executionError = p.Execute(NewBytecodeFromArray(elif.Body))
			return executionError
		}
	}
	_, executionError = p.Execute(NewBytecodeFromArray(ifInformation.Else))
	return executionError
}

func (p *Plasma) ifOneLinerOP(code Code) (Value, *Object) {
	ifInformation := code.Value.(*IfInformation)
	condition, executionError := p.Execute(NewBytecodeFromArray(ifInformation.Condition))
	if executionError != nil {
		return nil, executionError
	}
	var conditionBool bool
	conditionBool, executionError = p.QuickGetBool(condition)
	if conditionBool {
		return p.Execute(NewBytecodeFromArray(ifInformation.Body))
	} else if ifInformation.Else != nil {
		return p.Execute(NewBytecodeFromArray(ifInformation.Else))
	}
	return p.GetNone(), nil
}

func (p *Plasma) unlessOP(code Code) *Object {
	ifInformation := code.Value.(*IfInformation)
	condition, executionError := p.Execute(NewBytecodeFromArray(ifInformation.Condition))
	if executionError != nil {
		return executionError
	}
	var conditionBool bool
	conditionBool, executionError = p.QuickGetBool(condition)
	if !conditionBool {
		_, executionError = p.Execute(NewBytecodeFromArray(ifInformation.Body))
		return executionError
	}
	for _, elif := range ifInformation.ElifBlocks {
		condition, executionError = p.Execute(NewBytecodeFromArray(elif.Condition))
		if executionError != nil {
			return executionError
		}
		conditionBool, executionError = p.QuickGetBool(condition)
		if !conditionBool {
			_, executionError = p.Execute(NewBytecodeFromArray(elif.Body))
			return executionError
		}
	}
	_, executionError = p.Execute(NewBytecodeFromArray(ifInformation.Else))
	return executionError
}

func (p *Plasma) unlessOneLinerOP(code Code) (Value, *Object) {
	ifInformation := code.Value.(*IfInformation)
	condition, executionError := p.Execute(NewBytecodeFromArray(ifInformation.Condition))
	if executionError != nil {
		return nil, executionError
	}
	var conditionBool bool
	conditionBool, executionError = p.QuickGetBool(condition)
	if !conditionBool {
		return p.Execute(NewBytecodeFromArray(ifInformation.Body))
	} else if ifInformation.Else != nil {
		return p.Execute(NewBytecodeFromArray(ifInformation.Else))
	}
	return p.GetNone(), nil
}

// Special Instructions

func (p *Plasma) loadFunctionArgumentsOP(code Code) *Object {
	for _, argument := range code.Value.([]string) {
		value := p.PopObject()
		p.PeekSymbolTable().Set(argument, value)
	}
	return nil
}

func (p *Plasma) newFunctionOP(bytecode *Bytecode, code Code) *Object {
	functionInformation := code.Value.([2]int)
	codeLength := functionInformation[0]
	numberOfArguments := functionInformation[1]
	start := bytecode.index
	bytecode.index += codeLength
	end := bytecode.index
	functionCode := make([]Code, codeLength)
	copy(functionCode, bytecode.instructions[start:end])
	p.PushObject(p.NewFunction(false, p.PeekSymbolTable(), NewPlasmaFunction(numberOfArguments, functionCode)))
	return nil
}

func (p *Plasma) jumpOP(bytecode *Bytecode, code Code) *Object {
	bytecode.index += code.Value.(int)
	return nil
}

func (p *Plasma) setupDoWhileLoop(bytecode *Bytecode, code Code) *Object {
	condition := NewBytecodeFromArray(bytecode.NextN(code.Value.([2]int)[0]))
	body := NewBytecodeFromArray(bytecode.NextN(code.Value.([2]int)[1]))
	doWhileLoopEntry := &loopEntry{
		Action: NoAction,
	}
	p.LoopStack.Push(
		doWhileLoopEntry,
	)
	defer p.LoopStack.Pop()
loop:
	for {
	redoLocation:
		// Execute the body
		_, executionError := p.Execute(body)
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
		result, executionError = p.Execute(condition)
		condition.index = 0
		if executionError != nil {
			return executionError
		}
		var conditionBool bool
		conditionBool, executionError = p.QuickGetBool(result)
		if !conditionBool {
			break
		}
	}
	return nil
}

func (p *Plasma) setupWhileLoop(bytecode *Bytecode, code Code) *Object {
	condition := NewBytecodeFromArray(bytecode.NextN(code.Value.([2]int)[0]))
	body := NewBytecodeFromArray(bytecode.NextN(code.Value.([2]int)[1]))
	whileLoopEntry := &loopEntry{
		Action: NoAction,
	}
	p.LoopStack.Push(
		whileLoopEntry,
	)
	defer p.LoopStack.Pop()
loop:
	for {
		// First Evaluate the condition
		result, executionError := p.Execute(condition)
		condition.index = 0
		if executionError != nil {
			return executionError
		}

		var conditionBool bool
		conditionBool, executionError = p.QuickGetBool(result)
		if !conditionBool {
			break
		}
	redoLocation:
		// Execute the body
		_, executionError = p.Execute(body)
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

func (p *Plasma) reloadForLoopContext(context *map[string]Value, numberOfReceivers int, receivers []string, sourceHasNext Value, sourceNext Value) (bool, *Object) {
	hasNextObject, callError := p.CallFunction(sourceHasNext, sourceHasNext.SymbolTable())
	if callError != nil {
		return false, callError
	}
	hasNextObjectBool, executionError := p.QuickGetBool(hasNextObject)
	if executionError != nil {
		return false, executionError
	}
	if !hasNextObjectBool {
		return false, nil
	}
	var value Value
	value, callError = p.CallFunction(sourceNext, sourceNext.SymbolTable())
	if callError != nil {
		return false, callError
	}
	if numberOfReceivers == 1 {
		(*context)[receivers[0]] = value
		return true, nil
	}
	// Unpack it as first calling to iter
	valueIterFunc, getError := value.Get(Iter)
	if getError != nil {
		return false, p.NewObjectWithNameNotFoundError(value.GetClass(p), Iter)
	}
	var valueAsIter Value
	valueAsIter, callError = p.CallFunction(valueIterFunc, valueIterFunc.SymbolTable())
	if callError != nil {
		return false, callError
	}
	var valueAsIterHasNext Value
	valueAsIterHasNext, getError = valueAsIter.Get(HasNext)
	if getError != nil {
		return false, p.NewObjectWithNameNotFoundError(valueAsIter.GetClass(p), HasNext)
	}
	var valueAsIterNext Value
	valueAsIterNext, getError = valueAsIter.Get(Next)
	if getError != nil {
		return false, p.NewObjectWithNameNotFoundError(valueAsIter.GetClass(p), Next)
	}
	for index, receiver := range receivers {
		hasNextObject, callError = p.CallFunction(valueAsIterHasNext, valueAsIterHasNext.SymbolTable())
		if callError != nil {
			return false, callError
		}
		hasNextObjectBool, executionError = p.QuickGetBool(hasNextObject)
		if executionError != nil {
			return false, executionError
		}
		if !hasNextObjectBool {
			return false, p.NewInvalidNumberOfArgumentsError(numberOfReceivers, index+1)
		}
		value, callError = p.CallFunction(valueAsIterNext, valueAsIterNext.SymbolTable())
		if callError != nil {
			return false, callError
		}
		(*context)[receiver] = value
	}
	return true, nil
}

func (p *Plasma) setupForLoopOP(bytecode *Bytecode, code Code) *Object {
	source := p.PopObject()
	sourceIter, getError := source.Get(Iter)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(source.GetClass(p), Iter)
	}
	sourceAsIter, callError := p.CallFunction(sourceIter, sourceIter.SymbolTable())
	if callError != nil {
		return callError
	}
	loopSettings := code.Value.(ForLoopSettings)

	bodyBytecode := NewBytecodeFromArray(bytecode.NextN(loopSettings.BodyLength))
	receivers := loopSettings.Receivers
	forLoopEntry := &loopEntry{
		Action: NoAction,
	}
	p.LoopStack.Push(
		forLoopEntry,
	)
	defer p.LoopStack.Pop()
	context := map[string]Value{}
	var sourceHasNext Value
	sourceHasNext, getError = sourceAsIter.Get(HasNext)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(sourceAsIter.GetClass(p), HasNext)
	}
	var sourceNext Value
	sourceNext, getError = sourceAsIter.Get(Next)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(sourceAsIter.GetClass(p), Next)
	}
	receiversLength := len(receivers)
loop:
	for {
		// Update receivers
		// Check if the iteration can continue
		hasNext, loadSymbolsError := p.reloadForLoopContext(&context, receiversLength, receivers, sourceHasNext, sourceNext)
		if loadSymbolsError != nil {
			return loadSymbolsError
		}
		if !hasNext {
			break
		}
	redoLocation:
		// Load the receivers
		for receiver, object := range context {
			p.PeekSymbolTable().Set(receiver, object)
		}
		// Execute body
		_, bodyExecutionError := p.Execute(bodyBytecode)
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

func (p *Plasma) executeFinally(finally []Code) *Object {
	if finally != nil {
		_, executionError := p.Execute(NewBytecodeFromArray(finally))
		return executionError
	}
	return nil
}

func (p *Plasma) tryOP(code Code) *Object {
	tryInformation := code.Value.(*TryInformation)
	_, executionError := p.Execute(NewBytecodeFromArray(tryInformation.Body))
	if executionError == nil {
		return p.executeFinally(tryInformation.Finally)
	}
	var targetError Value
	var executionError2 *Object
	for _, exceptBlock := range tryInformation.ExceptBlocks {
		if exceptBlock.TargetErrors != nil {
			for _, targetErrorCode := range exceptBlock.TargetErrors {
				targetError, executionError2 = p.Execute(NewBytecodeFromArray(targetErrorCode))
				if executionError2 != nil {
					return executionError
				}
				if !targetError.Implements(p.ForceMasterGetAny(RuntimeError).(*Type)) {
					return p.NewInvalidTypeError(targetError.TypeName(), RuntimeError)
				}
				if executionError.Implements(targetError.(*Type)) {
					p.PeekSymbolTable().Set(exceptBlock.Receiver, executionError)
					_, executionError2 = p.Execute(NewBytecodeFromArray(exceptBlock.Body))
					if executionError2 == nil {
						return p.executeFinally(tryInformation.Finally)
					}
					return executionError2
				}
			}
		} else {
			p.PeekSymbolTable().Set(exceptBlock.Receiver, executionError)
			_, executionError2 = p.Execute(NewBytecodeFromArray(exceptBlock.Body))
			if executionError2 == nil {
				return p.executeFinally(tryInformation.Finally)
			}
			return executionError2
		}
	}
	if tryInformation.Else != nil {
		_, executionError2 = p.Execute(NewBytecodeFromArray(tryInformation.Else))
		if executionError2 != nil {
			return executionError2
		}
	}
	return p.executeFinally(tryInformation.Finally)
}

type ModuleInformation struct {
	Name       string
	CodeLength int
}

func (p *Plasma) newModuleOP(bytecode *Bytecode, code Code) *Object {
	moduleInformation := code.Value.(ModuleInformation)
	var moduleBody []Code
	for i := 0; i < moduleInformation.CodeLength; i++ {
		moduleBody = append(moduleBody, bytecode.Next())
	}
	module := p.NewModule(false, p.PeekSymbolTable())
	p.PushSymbolTable(module.SymbolTable())
	_, executionError := p.Execute(NewBytecodeFromArray(moduleBody))
	if executionError != nil {
		return executionError
	}
	p.PopSymbolTable()
	p.PeekSymbolTable().Set(moduleInformation.Name, module)
	return nil
}

type ClassInformation struct {
	Name       string
	BodyLength int
}

func (p *Plasma) newClassOP(bytecode *Bytecode, code Code) *Object {
	classInformation := code.Value.(ClassInformation)
	rawSubClasses := p.PopObject().GetContent()
	var subClasses []*Type
	for _, subClass := range rawSubClasses {
		if _, ok := subClass.(*Type); !ok {
			return p.NewInvalidTypeError(subClass.TypeName(), TypeName)
		}
		subClasses = append(subClasses, subClass.(*Type))
	}

	var classBody []Code
	for i := 0; i < classInformation.BodyLength; i++ {
		classBody = append(classBody, bytecode.Next())
	}
	class := p.NewType(
		false,
		classInformation.Name,
		p.PeekSymbolTable(),
		subClasses,
		NewPlasmaConstructor(classBody),
	)
	p.PeekSymbolTable().Set(classInformation.Name, class)
	return nil
}

func (p *Plasma) newClassFunctionOP(bytecode *Bytecode, code Code) *Object {
	functionInformation := code.Value.([2]int)
	codeLength := functionInformation[0]
	numberOfArguments := functionInformation[1]
	start := bytecode.index
	bytecode.index += codeLength
	end := bytecode.index
	functionCode := make([]Code, codeLength)
	copy(functionCode, bytecode.instructions[start:end])
	p.PushObject(p.NewFunction(false, p.PeekSymbolTable(), NewPlasmaClassFunction(p.PeekObject(), numberOfArguments, functionCode)))
	return nil
}

func (p *Plasma) raiseOP() *Object {
	if _, ok := p.PeekObject().(*Object); !ok {
		return p.NewInvalidTypeError(p.PeekObject().TypeName(), RuntimeError)
	}
	if !p.PeekObject().Implements(p.ForceMasterGetAny(RuntimeError).(*Type)) {
		return p.NewInvalidTypeError(p.PeekObject().TypeName(), RuntimeError)
	}
	return p.PeekObject().(*Object)
}

func (p *Plasma) caseOP(bytecode *Bytecode, code Code) *Object {
	references := p.PopObject()
	contains := p.ForceGetSelf(Contains, references)
	result, callError := p.CallFunction(contains, references.SymbolTable(), p.PeekObject())
	if callError != nil {
		return callError
	}
	boolResult, executionError := p.QuickGetBool(result)
	if executionError != nil {
		return executionError
	}
	if !boolResult {
		bytecode.index += code.Value.(int)
		return nil
	}
	p.PopObject()
	return nil
}
