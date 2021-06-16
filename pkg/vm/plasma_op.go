package vm

import "fmt"

func (p *Plasma) Execute() (Value, *Object) {
	var executionError *Object
	for ; p.PeekBytecode().HasNext(); {
		code := p.PeekBytecode().Next()
		switch code.Instruction.OpCode {
		// Literals
		case NewStringOP:
			executionError = p.newStringOP(code)
		case NewBytesOP:
			executionError = p.newBytesOP(code)
		case NewIntegerOP:
			executionError = p.newIntegerOP(code)
		case NewFloatOP:
			executionError = p.newFloatOP(code)
		case NewTrueBoolOP:
			executionError = p.newTrueBoolOP()
		case NewFalseBoolOP:
			executionError = p.newFalseBoolOP()
		case GetNoneOP:
			executionError = p.getNoneOP()
		// Composite creation
		case NewTupleOP:
			executionError = p.newTupleOP(code)
		case NewArrayOP:
			executionError = p.newArrayOP(code)
		case NewHashOP:
			executionError = p.newHashOP(code)
		// Unary Expressions
		case NegateBitsOP:
			executionError = p.noArgsGetAndCall(NegBits)
		case BoolNegateOP:
			executionError = p.noArgsGetAndCall(Negate)
		case NegativeOP:
			executionError = p.noArgsGetAndCall(Negative)
		// Binary Expressions
		case AddOP:
			executionError = p.leftBinaryExpressionFuncCall(Add)
		case SubOP:
			executionError = p.leftBinaryExpressionFuncCall(Sub)
		case MulOP:
			executionError = p.leftBinaryExpressionFuncCall(Mul)
		case DivOP:
			executionError = p.leftBinaryExpressionFuncCall(Div)
		case FloorDivOP:
			executionError = p.leftBinaryExpressionFuncCall(FloorDiv)
		case ModOP:
			executionError = p.leftBinaryExpressionFuncCall(Mod)
		case PowOP:
			executionError = p.leftBinaryExpressionFuncCall(Pow)
		case BitXorOP:
			executionError = p.leftBinaryExpressionFuncCall(BitXor)
		case BitAndOP:
			executionError = p.leftBinaryExpressionFuncCall(BitAnd)
		case BitOrOP:
			executionError = p.leftBinaryExpressionFuncCall(BitOr)
		case BitLeftOP:
			executionError = p.leftBinaryExpressionFuncCall(BitLeft)
		case BitRightOP:
			executionError = p.leftBinaryExpressionFuncCall(BitRight)
		case AndOP:
			executionError = p.leftBinaryExpressionFuncCall(And)
		case OrOP:
			executionError = p.leftBinaryExpressionFuncCall(Or)
		case XorOP:
			executionError = p.leftBinaryExpressionFuncCall(Xor)
		case EqualsOP:
			executionError = p.leftBinaryExpressionFuncCall(Equals)
		case NotEqualsOP:
			executionError = p.leftBinaryExpressionFuncCall(NotEquals)
		case GreaterThanOP:
			executionError = p.leftBinaryExpressionFuncCall(GreaterThan)
		case LessThanOP:
			executionError = p.leftBinaryExpressionFuncCall(LessThan)
		case GreaterThanOrEqualOP:
			executionError = p.leftBinaryExpressionFuncCall(GreaterThanOrEqual)
		case LessThanOrEqualOP:
			executionError = p.leftBinaryExpressionFuncCall(LessThanOrEqual)
		case ContainsOP:
			leftHandSide := p.PopObject()
			rightHandSide := p.PopObject()
			p.PushObject(leftHandSide)
			p.PushObject(rightHandSide)
			executionError = p.leftBinaryExpressionFuncCall(Contains)
		// Other Expressions
		case GetIdentifierOP:
			executionError = p.getIdentifierOP(code)
		case IndexOP:
			executionError = p.indexOP()
		case SelectNameFromObjectOP:
			executionError = p.selectNameFromObjectOP(code)
		case MethodInvocationOP:
			executionError = p.methodInvocationOP(code)
		// Assign Statement
		case AssignIdentifierOP:
			executionError = p.assignIdentifierOP(code)
		case AssignSelectorOP:
			executionError = p.assignSelectorOP(code)
		case AssignIndexOP:
			executionError = p.assignIndexOP()
		case ReturnOP:
			executionError = p.returnOP(code)
			if executionError != nil {
				return nil, executionError
			}
			if p.MemoryStack.HasNext() {
				return p.PopObject(), nil
			}
			return p.PeekSymbolTable().Symbols[None], nil
		case IfJumpOP:
			executionError = p.ifJumpOP(code)
		case UnlessJumpOP:
			executionError = p.unlessJumpOP(code)
		// Special Instructions
		case LoadFunctionArgumentsOP:
			executionError = p.loadFunctionArgumentsOP(code)
		case NewFunctionOP:
			executionError = p.newFunctionOP(code)
		case JumpOP:
			executionError = p.jumpOP(code)
		case RedoOP:
			executionError = p.redoOP(code)
		case BreakOP:
			executionError = p.breakOP(code)
		case ContinueOP:
			executionError = p.continueOP(code)
		case PopOP:
			p.MemoryStack.Pop()
		case NOP:
			break
		case SetupForLoopOP:
			executionError = p.setupForLoopOP()
		case HasNextOP:
			executionError = p.hasNextOP(code)
		case UnpackReceiversPopOP:
			executionError = p.unpackReceiversPopOP()
		case UnpackReceiversPeekOP:
			executionError = p.unpackReceiversPeekOP(code)
		case PopIterOP:
			p.IterStack.Pop()
		case NewIteratorOP:
			executionError = p.newIteratorOP(code)
		case SetupTryBlockOP:
			executionError = p.setupTryBlockOP(code)
		case SetupTryExceptBlockOP:
			executionError = p.setupTryExceptBlockOP(code)
		case SetupTryElseBlockOP:
			executionError = p.setupTryElseBlockOP(code)
		case SetupTryFinallyBlockOP:
			executionError = p.setupTryFinallyBlockOP(code)
		case ExitTryBlockOP:
			executionError = p.exitTryBlockOP()
		case NewModuleOP:
			executionError = p.newModuleOP(code)
		case RaiseOP:
			executionError = p.raiseOP()
		case NewClassOP:
			executionError = p.newClassOP(code)
		case NewClassFunctionOP:
			executionError = p.newClassFunctionOP(code)
		case CaseOP:
			executionError = p.caseOP(code)
		default:
			panic(fmt.Sprintf("Unknown VM instruction %d", code.Instruction.OpCode))
		}
		if executionError != nil {
			// Here should be some of the code related to the try-except block
			if p.TryStack.HasNext() {
				executionError = p.handleTryExcepts(executionError)
				if executionError != nil && !p.TryStack.HasNext() {
					return nil, executionError
				}
				continue
			}
			return nil, executionError
		}
	}
	if p.MemoryStack.HasNext() {
		return p.PopObject(), nil
	}
	return p.NewNone(), nil
}

func (p *Plasma) newStringOP(code Code) *Object {
	value := code.Value.(string)
	stringObject := p.NewString(false, p.SymbolTableStack.Peek(), value)
	p.PushObject(stringObject)
	return nil
}

func (p *Plasma) newBytesOP(code Code) *Object {
	value := code.Value.([]byte)
	bytesObject := p.NewBytes(false, p.SymbolTableStack.Peek(), value)
	p.PushObject(bytesObject)
	return nil
}

func (p *Plasma) newIntegerOP(code Code) *Object {
	value := code.Value.(int64)
	integer := p.NewInteger(false, p.SymbolTableStack.Peek(), value)
	p.PushObject(integer)
	return nil
}

func (p *Plasma) newFloatOP(code Code) *Object {
	value := code.Value.(float64)
	float := p.NewFloat(false, p.SymbolTableStack.Peek(), value)
	p.PushObject(float)
	return nil
}

func (p *Plasma) newTrueBoolOP() *Object {
	p.PushObject(p.NewBool(false, p.PeekSymbolTable(), true))
	return nil
}

func (p *Plasma) newFalseBoolOP() *Object {
	p.PushObject(p.NewBool(false, p.PeekSymbolTable(), false))
	return nil
}

func (p *Plasma) getNoneOP() *Object {
	p.PushObject(p.NewNone())
	return nil
}

func (p *Plasma) newTupleOP(code Code) *Object {
	numberOfValues := code.Value.(int)
	var values []Value
	for i := 0; i < numberOfValues; i++ {
		values = append(values, p.PopObject())
	}
	p.PushObject(p.NewTuple(false, p.PeekSymbolTable(), values))
	return nil
}

func (p *Plasma) newArrayOP(code Code) *Object {
	numberOfValues := code.Value.(int)
	var values []Value
	for i := 0; i < numberOfValues; i++ {
		values = append(values, p.PopObject())
	}
	p.PushObject(p.NewArray(false, p.PeekSymbolTable(), values))
	return nil
}

func (p *Plasma) newHashOP(code Code) *Object {
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
		return p.NewObjectWithNameNotFoundError(Assign)
	}
	for _, keyValue := range keyValues {
		_, assignError := p.CallFunction(hashTableAssign.(*Function), hashTable.SymbolTable(), keyValue.Key, keyValue.Value)
		if assignError != nil {
			return assignError
		}
	}
	p.PushObject(hashTable)
	return nil
}

// Useful function to call those built ins that doesn't receive arguments of an object
func (p *Plasma) noArgsGetAndCall(operationName string) *Object {
	object := p.PopObject()
	operation, getError := object.Get(operationName)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(operationName)
	}
	if _, ok := operation.(*Function); !ok {
		return p.NewInvalidTypeError(operation.TypeName(), FunctionName)
	}
	result, callError := p.CallFunction(operation.(*Function), object.SymbolTable())
	if callError != nil {
		return callError
	}
	p.PushObject(result)
	return nil
}

// Function useful to cal object built-ins binary expression functions
func (p *Plasma) leftBinaryExpressionFuncCall(operationName string) *Object {
	leftHandSide := p.PopObject()
	rightHandSide := p.PopObject()
	operation, getError := leftHandSide.Get(operationName)
	if getError != nil {
		return p.rightBinaryExpressionFuncCall(leftHandSide, rightHandSide, operationName)
	}
	if _, ok := operation.(*Function); !ok {
		return p.rightBinaryExpressionFuncCall(leftHandSide, rightHandSide, operationName)
	}
	result, callError := p.CallFunction(operation.(*Function), leftHandSide.SymbolTable(), rightHandSide)
	if callError != nil {
		return p.rightBinaryExpressionFuncCall(leftHandSide, rightHandSide, operationName)
	}
	p.PushObject(result)
	return nil
}

func (p *Plasma) rightBinaryExpressionFuncCall(leftHandSide Value, rightHandSide Value, operationName string) *Object {
	operation, getError := rightHandSide.Get("Right" + operationName)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError("Right" + operationName)
	}
	if _, ok := operation.(*Function); !ok {
		return p.NewInvalidTypeError(operation.TypeName(), FunctionName)
	}
	result, callError := p.CallFunction(operation.(*Function), rightHandSide.SymbolTable(), leftHandSide)
	if callError != nil {
		return callError
	}
	p.PushObject(result)
	return nil
}

func (p *Plasma) indexOP() *Object {
	index := p.PopObject()
	source := p.PopObject()
	indexOperation, getError := source.Get(Index)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(Index)
	}
	if _, ok := indexOperation.(*Function); !ok {
		return p.NewInvalidTypeError(indexOperation.TypeName(), FunctionName)
	}
	result, callError := p.CallFunction(indexOperation.(*Function), source.SymbolTable(), index)
	if callError != nil {
		return callError
	}
	p.PushObject(result)
	return nil
}

func (p *Plasma) selectNameFromObjectOP(code Code) *Object {
	name := code.Value.(string)
	source := p.PopObject()
	value, getError := source.Get(name)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(name)
	}
	p.PushObject(value)
	return nil
}

func (p *Plasma) methodInvocationOP(code Code) *Object {
	numberOfArguments := code.Value.(int)
	function := p.PopObject()
	var arguments []Value
	for i := 0; i < numberOfArguments; i++ {
		if !p.MemoryStack.HasNext() {
			return p.NewInvalidNumberOfArgumentsError(i, numberOfArguments)
		}
		arguments = append(arguments, p.PopObject())
	}
	var result Value
	var callError *Object
	switch function.(type) {
	case *Function:
		result, callError = p.CallFunction(function.(*Function), NewSymbolTable(function.SymbolTable().Parent), arguments...)
	case *Type:
		result, callError = p.ConstructObject(function.(*Type), NewSymbolTable(function.SymbolTable().Parent))
		if callError != nil {
			return callError
		}
		resultInitialize, getError := result.Get(Initialize)
		if getError != nil {
			return p.NewObjectWithNameNotFoundError(Initialize)
		}
		if _, ok := resultInitialize.(*Function); !ok {
			return p.NewInvalidTypeError(resultInitialize.TypeName(), FunctionName)
		}
		_, callError = p.CallFunction(resultInitialize.(*Function), result.SymbolTable(), arguments...)
	default:
		// Try getting the method Call
		call, getError := function.Get(Call)
		if getError != nil {
			return p.NewObjectWithNameNotFoundError(Call)
		}
		if _, ok := call.(*Function); !ok {
			return p.NewInvalidTypeError(call.TypeName(), FunctionName)
		}
		result, callError = p.CallFunction(call.(*Function), NewSymbolTable(function.SymbolTable().Parent))
	}
	if callError != nil {
		return callError
	}
	p.PushObject(result)
	return nil
}

func (p *Plasma) getIdentifierOP(code Code) *Object {
	value, getError := p.PeekSymbolTable().GetAny(code.Value.(string))
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(code.Value.(string))
	}
	p.PushObject(value)
	return nil
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
		return p.NewObjectWithNameNotFoundError(Assign)
	}
	if _, ok := sourceAssign.(*Function); !ok {
		return p.NewInvalidTypeError(sourceAssign.TypeName(), FunctionName)
	}
	_, callError := p.CallFunction(sourceAssign.(*Function), p.PeekSymbolTable(), index, value)
	if callError != nil {
		return callError
	}
	return nil
}

func (p *Plasma) returnOP(code Code) *Object {
	numberOfReturnValues := code.Value.(int)
	if numberOfReturnValues == 0 {
		p.PushObject(p.NewNone())
		return nil
	}

	var values []Value
	for i := 0; i < numberOfReturnValues; i++ {
		if !p.MemoryStack.HasNext() {
			return p.NewInvalidNumberOfArgumentsError(i, numberOfReturnValues)
		}
		values = append(values, p.PopObject())
	}
	if len(values) == 1 {
		p.PushObject(values[0])
	} else {
		p.PushObject(p.NewTuple(false, p.PeekSymbolTable(), values))
	}
	return nil
}

func (p *Plasma) ifJumpOP(code Code) *Object {
	condition := p.PopObject()
	toBool, getError := condition.Get(ToBool)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(ToBool)
	}
	if _, ok := toBool.(*Function); !ok {
		return p.NewInvalidTypeError(toBool.TypeName(), FunctionName)
	}
	conditionBool, callError := p.CallFunction(toBool.(*Function), toBool.SymbolTable())
	if callError != nil {
		return callError
	}
	if !conditionBool.GetBool() {
		p.PeekBytecode().index += code.Value.(int)
	}
	return nil
}

func (p *Plasma) unlessJumpOP(code Code) *Object {
	condition := p.PopObject()
	toBool, getError := condition.Get(ToBool)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(ToBool)
	}
	if _, ok := toBool.(*Function); !ok {
		return p.NewInvalidTypeError(toBool.TypeName(), FunctionName)
	}
	conditionBool, callError := p.CallFunction(toBool.(*Function), toBool.SymbolTable())
	if callError != nil {
		return callError
	}
	if conditionBool.GetBool() {
		p.PeekBytecode().index += code.Value.(int)
	}
	return nil
}

// Special Instructions

func (p *Plasma) loadFunctionArgumentsOP(code Code) *Object {
	for _, argument := range code.Value.([]string) {
		value := p.PopObject()
		p.PeekSymbolTable().Set(argument, value)
	}
	return nil
}

func (p *Plasma) newFunctionOP(code Code) *Object {
	functionInformation := code.Value.([2]int)
	codeLength := functionInformation[0]
	numberOfArguments := functionInformation[1]
	start := p.PeekBytecode().index
	p.PeekBytecode().index += codeLength
	end := p.PeekBytecode().index
	functionCode := make([]Code, codeLength)
	copy(functionCode, p.PeekBytecode().instructions[start:end])
	p.PushObject(p.NewFunction(false, p.PeekSymbolTable(), NewPlasmaFunction(numberOfArguments, functionCode)))
	return nil
}

func (p *Plasma) jumpOP(code Code) *Object {
	p.PeekBytecode().index += code.Value.(int)
	return nil
}

func (p *Plasma) breakOP(code Code) *Object {
	p.PeekBytecode().index += code.Value.(int)
	return nil
}

func (p *Plasma) redoOP(code Code) *Object {
	p.PeekBytecode().index += code.Value.(int)
	return nil
}

func (p *Plasma) continueOP(code Code) *Object {
	p.PeekBytecode().index += code.Value.(int) // This could be positive or negative (positive only for do-while loop)
	return nil
}

func (p *Plasma) setupForLoopOP() *Object {
	// First check if it is a iterator
	value := p.PopObject()
	if _, ok := value.(*Iterator); ok {
		p.IterStack.Push(&IterEntry{
			Iterable:  value,
			LastValue: nil,
		})
		return nil
	}
	// Then check if the type implements iterator
	_, getError := value.Get(HasNext)
	if getError == nil {
		_, getError = value.Get(Next)
		if getError == nil {
			p.IterStack.Push(&IterEntry{
				Iterable:  value,
				LastValue: nil,
			})
			return nil
		}
	}
	// Finally transform it to iterable
	var valueIterFunc Value
	valueIterFunc, getError = value.Get(Iter)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(Iter)
	}
	if _, ok := valueIterFunc.(*Function); !ok {
		return p.NewInvalidTypeError(valueIterFunc.TypeName(), FunctionName)
	}
	valueIter, callError := p.CallFunction(valueIterFunc.(*Function), value.SymbolTable())
	if callError != nil {
		return callError
	}
	p.IterStack.Push(&IterEntry{
		Iterable:  valueIter,
		LastValue: nil,
	})
	return nil
}

func (p *Plasma) hasNextOP(code Code) *Object {
	hasNext, getError := p.IterStack.Peek().Iterable.Get(HasNext)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(HasNext)
	}
	if _, ok := hasNext.(*Function); !ok {
		return p.NewInvalidTypeError(hasNext.TypeName(), FunctionName)
	}
	result, callError := p.CallFunction(hasNext.(*Function), p.IterStack.Peek().Iterable.SymbolTable())
	if callError != nil {
		return callError
	}
	if _, ok := result.(*Bool); !ok {
		var resultToBool Value
		resultToBool, getError = hasNext.Get(ToBool)
		if getError != nil {
			return p.NewObjectWithNameNotFoundError(ToBool)
		}
		if _, ok = resultToBool.(*Function); !ok {
			return p.NewInvalidTypeError(resultToBool.TypeName(), FunctionName)
		}
		var resultBool Value
		resultBool, callError = p.CallFunction(resultToBool.(*Function), hasNext.SymbolTable())
		if callError != nil {
			return callError
		}
		if _, ok = resultBool.(*Bool); !ok {
			return p.NewInvalidTypeError(resultBool.TypeName(), BoolName)
		}
		result = resultBool
	}
	if !result.GetBool() {
		p.PeekBytecode().index += code.Value.(int)
	}
	return nil
}

func (p *Plasma) unpackReceiversPopOP() *Object {
	next, getError := p.IterStack.Peek().Iterable.Get(Next)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(Next)
	}
	if _, ok := next.(*Function); !ok {
		return p.NewInvalidTypeError(next.TypeName(), FunctionName)
	}
	nextValue, callError := p.CallFunction(next.(*Function), p.IterStack.Peek().Iterable.SymbolTable())
	if callError != nil {
		return callError
	}
	p.IterStack.Peek().LastValue = nextValue
	return nil
}

func (p *Plasma) unpackReceiversPeekOP(code Code) *Object {
	receivers := code.Value.([]string)
	if len(receivers) == 1 {
		p.PeekSymbolTable().Set(receivers[0], p.IterStack.Peek().LastValue)
		return nil
	}
	// First try to unpack iterators
	hasNext, getError := p.IterStack.Peek().LastValue.Get(HasNext)
	if _, ok := hasNext.(*Function); getError == nil && ok {
		var next Value
		next, getError = p.IterStack.Peek().LastValue.Get(Next)
		if _, ok = next.(*Function); getError == nil && ok {
			for _, receiver := range receivers {
				// First check if there is next value
				hasNextResult, callError := p.CallFunction(hasNext.(*Function), p.IterStack.Peek().LastValue.SymbolTable())
				if callError != nil {
					return callError
				}
				var hasNextResultBool Value
				if _, ok = hasNextResult.(*Bool); !ok {
					var hasNextResultToBool Value
					hasNextResultToBool, getError = hasNextResult.Get(ToBool)
					if getError != nil {
						return p.NewObjectWithNameNotFoundError(ToBool)
					}
					if _, ok = hasNextResultToBool.(*Function); !ok {
						return p.NewInvalidTypeError(hasNextResultToBool.TypeName(), FunctionName)
					}
					hasNextResultBool, callError = p.CallFunction(hasNextResultToBool.(*Function), hasNextResult.SymbolTable())
					if callError != nil {
						return callError
					}
					if _, ok = hasNextResultBool.(*Bool); !ok {
						return p.NewInvalidTypeError(hasNextResultBool.TypeName(), BoolName)
					}
					hasNextResult = hasNextResultBool
				}
				if hasNextResult.GetBool() {
					var value Value
					value, callError = p.CallFunction(next.(*Function), p.IterStack.Peek().LastValue.SymbolTable())
					if callError != nil {
						return callError
					}
					p.PeekSymbolTable().Set(receiver, value)
				}
			}
			return nil
		}
	}
	// Then try to unpack index-ables
	var lastValue Value
	if _, ok := p.IterStack.Peek().LastValue.(*Tuple); !ok {
		var toTuple Value
		toTuple, getError = p.IterStack.Peek().LastValue.Get(ToTuple)
		if getError != nil {
			return p.NewObjectWithNameNotFoundError(ToTuple)
		}
		if _, ok = toTuple.(*Function); !ok {
			return p.NewInvalidTypeError(toTuple.TypeName(), FunctionName)
		}
		var callError *Object
		lastValue, callError = p.CallFunction(toTuple.(*Function), p.IterStack.Peek().LastValue.SymbolTable())
		if callError != nil {
			return callError
		}
		if _, ok = lastValue.(*Tuple); !ok {
			return p.NewInvalidTypeError(lastValue.TypeName(), TupleName)
		}
	} else {
		lastValue = p.IterStack.Peek().LastValue
	}
	if len(receivers) != lastValue.GetLength() {
		return p.NewInvalidNumberOfArgumentsError(lastValue.GetLength(), len(receivers))
	}
	for i, receiver := range receivers {
		p.PeekSymbolTable().Set(receiver, lastValue.GetContent()[i])
	}
	return nil
}

func (p *Plasma) newIteratorOP(code Code) *Object {
	source := p.PopObject()
	var iterSource Value
	var callError *Object
	if _, ok := source.(*Iterator); ok {
		iterSource = source
	} else {
		iter, getError := source.Get(Iter)
		if getError != nil {
			return p.NewObjectWithNameNotFoundError(Iter)
		}
		if _, ok = iter.(*Function); !ok {
			return p.NewInvalidTypeError(iter.TypeName(), FunctionName)
		}
		iterSource, callError = p.CallFunction(iter.(*Function), source.SymbolTable())
		if callError != nil {
			return callError
		}
	}
	generatorIterator := p.NewIterator(false, p.PeekSymbolTable())
	generatorIterator.Set(Source, iterSource)

	hasNextCodeLength, nextCodeLength := code.Value.([2]int)[0], code.Value.([2]int)[1]
	var hasNextCode []Code
	for i := 0; i < hasNextCodeLength; i++ {
		hasNextCode = append(hasNextCode, p.PeekBytecode().Next())
	}
	var nextCode []Code
	for i := 0; i < nextCodeLength; i++ {
		nextCode = append(nextCode, p.PeekBytecode().Next())
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
	p.PushObject(generatorIterator)
	return nil
}

func (p *Plasma) setupTryBlockOP(code Code) *Object {
	p.TryStack.Push(
		&tryStackEntry{
			finalIndex:   p.PeekBytecode().index + code.Value.(int) - 1,
			exceptBlocks: nil,
			elseBlock:    nil,
			finallyBody:  nil,
		},
	)
	return nil
}

func (p *Plasma) setupTryExceptBlockOP(code Code) *Object {
	rawExcept := code.Value.(ExceptInformation)
	except := exceptBlock{
		targets:  nil,
		receiver: rawExcept.Receiver,
		body:     nil,
	}
	for ii := 0; ii < rawExcept.TargetsLength; ii++ {
		except.targets = append(except.targets, p.PeekBytecode().Next())
	}
	for iii := 0; iii < rawExcept.BodyLength; iii++ {
		except.body = append(except.body, p.PeekBytecode().Next())
	}
	p.TryStack.Peek().exceptBlocks = append(p.TryStack.Peek().exceptBlocks, except)
	return nil
}

func (p *Plasma) setupTryElseBlockOP(code Code) *Object {
	elseLength := code.Value.(int)
	for i := 0; i < elseLength; i++ {
		p.TryStack.Peek().elseBlock = append(p.TryStack.Peek().elseBlock, p.PeekBytecode().Next())
	}
	return nil
}

func (p *Plasma) setupTryFinallyBlockOP(code Code) *Object {
	finallyLength := code.Value.(int)
	for i := 0; i < finallyLength; i++ {
		p.TryStack.Peek().finallyBody = append(p.TryStack.Peek().finallyBody, p.PeekBytecode().Next())
	}
	return nil
}

func (p *Plasma) exitTryBlockOP() *Object {
	return p.executeFinally(p.TryStack.Pop().finallyBody)
}

func (p *Plasma) executeFinally(finallyBody []Code) *Object {
	if len(finallyBody) > 0 {
		p.PushBytecode(NewBytecodeFromArray(finallyBody))
		_, executionError := p.Execute()
		p.PopBytecode()
		if executionError != nil {
			return executionError
		}
	}
	return nil
}

func (p *Plasma) handleTryExcepts(exception *Object) *Object {
	entry := p.TryStack.Pop()
	for _, except := range entry.exceptBlocks {
		if len(except.targets) == 1 {
			p.PushBytecode(NewBytecodeFromArray(except.body))
			_, executionError := p.Execute()
			p.PopBytecode()
			if executionError != nil {
				return executionError
			}
			finallyExecutionError := p.executeFinally(entry.finallyBody)
			if finallyExecutionError != nil {
				return finallyExecutionError
			}
			p.PeekBytecode().index = entry.finalIndex
			return nil
		}
		p.PushBytecode(NewBytecodeFromArray(except.targets))
		targetsTuple, executionError := p.Execute()
		p.PopBytecode()
		if executionError != nil {
			return executionError
		}
		if _, ok := targetsTuple.(*Tuple); !ok {
			return p.NewInvalidTypeError(targetsTuple.TypeName(), TupleName)
		}
		runtimeError := p.ForceMasterGetAny(RuntimeError).(*Type)
		for _, target := range targetsTuple.GetContent() {
			if _, ok := target.(*Type); !ok {
				return p.NewInvalidTypeError(target.TypeName(), TypeName)
			}
			if !target.Implements(runtimeError) {
				return p.NewInvalidTypeError(target.TypeName(), RuntimeError)
			}
			if exception.class == target {
				p.PeekSymbolTable().Set(except.receiver, exception)
				p.PushBytecode(NewBytecodeFromArray(except.body))
				_, executionError = p.Execute()
				p.PopBytecode()
				if executionError != nil {
					return executionError
				}
				finallyExecutionError := p.executeFinally(entry.finallyBody)
				if finallyExecutionError != nil {
					return finallyExecutionError
				}
				p.PeekBytecode().index = entry.finalIndex
				return nil
			}
		}
	}
	if len(entry.elseBlock) > 0 {
		p.PushBytecode(NewBytecodeFromArray(entry.elseBlock))
		_, executionError := p.Execute()
		p.PopBytecode()
		if executionError != nil {
			return executionError
		}
		finallyExecutionError := p.executeFinally(entry.finallyBody)
		if finallyExecutionError != nil {
			return finallyExecutionError
		}
		p.PeekBytecode().index = entry.finalIndex
		return nil
	}
	return exception
}

type ModuleInformation struct {
	Name       string
	CodeLength int
}

func (p *Plasma) newModuleOP(code Code) *Object {
	moduleInformation := code.Value.(ModuleInformation)
	var moduleBody []Code
	for i := 0; i < moduleInformation.CodeLength; i++ {
		moduleBody = append(moduleBody, p.PeekBytecode().Next())
	}
	module := p.NewModule(false, p.PeekSymbolTable())
	p.PushSymbolTable(module.SymbolTable())
	p.PushBytecode(NewBytecodeFromArray(moduleBody))
	_, executionError := p.Execute()
	if executionError != nil {
		return executionError
	}
	p.PopBytecode()
	p.PopSymbolTable()
	p.PeekSymbolTable().Set(moduleInformation.Name, module)
	return nil
}

type ClassInformation struct {
	Name       string
	BodyLength int
}

func (p *Plasma) newClassOP(code Code) *Object {
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
		classBody = append(classBody, p.PeekBytecode().Next())
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

func (p *Plasma) newClassFunctionOP(code Code) *Object {
	functionInformation := code.Value.([2]int)
	codeLength := functionInformation[0]
	numberOfArguments := functionInformation[1]
	start := p.PeekBytecode().index
	p.PeekBytecode().index += codeLength
	end := p.PeekBytecode().index
	functionCode := make([]Code, codeLength)
	copy(functionCode, p.PeekBytecode().instructions[start:end])
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

func (p *Plasma) caseOP(code Code) *Object {
	references := p.PopObject()
	contains := p.ForceParentGetSelf(Contains, references.SymbolTable())
	result, callError := p.CallFunction(contains.(*Function), references.SymbolTable(), p.PeekObject())
	if callError != nil {
		return callError
	}
	var boolResult Value
	if _, ok := result.(*Bool); ok {
		boolResult = result
	} else {
		toBool, getError := result.Get(ToBool)
		if getError != nil {
			return p.NewObjectWithNameNotFoundError(ToBool)
		}
		if _, ok = toBool.(*Function); !ok {
			return p.NewInvalidTypeError(toBool.TypeName(), FunctionName)
		}
		boolResult, callError = p.CallFunction(toBool.(*Function), result.SymbolTable())
		if callError != nil {
			return callError
		}
	}
	if !boolResult.GetBool() {
		p.PeekBytecode().index += code.Value.(int)
		return nil
	}
	p.PopObject()
	return nil
}
