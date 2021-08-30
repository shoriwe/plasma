package vm

func (p *Plasma) newStringOP(context *Context, s string) *Value {
	context.LastObject = p.NewString(context, false, s)
	return nil
}

func (p *Plasma) newBytesOP(context *Context, bytes []uint8) *Value {
	context.LastObject = p.NewBytes(context, false, bytes)
	return nil
}

func (p *Plasma) newIntegerOP(context *Context, i int64) *Value {
	context.LastObject = p.NewInteger(context, false, i)
	return nil
}

func (p *Plasma) newFloatOP(context *Context, f float64) *Value {
	context.LastObject = p.NewFloat(context, false, f)
	return nil
}

func (p *Plasma) newArrayOP(context *Context, length int) *Value {
	content := make([]*Value, length)
	for index := 0; index < length; index++ {
		content[index] = context.PopObject()
	}
	context.LastObject = p.NewArray(context, false, content)
	return nil
}

func (p *Plasma) newTupleOP(context *Context, length int) *Value {
	content := make([]*Value, length)
	for index := 0; index < length; index++ {
		content[index] = context.PopObject()
	}
	context.LastObject = p.NewTuple(context, false, content)
	return nil
}

func (p *Plasma) newHashTableOP(context *Context, length int) *Value {
	result := p.NewHashTable(context, false)
	for index := 0; index < length; index++ {
		key := context.PopObject()
		value := context.PopObject()
		assignResult, success := p.HashIndexAssign(context, result, key, value)
		if !success {
			return assignResult
		}
	}
	context.LastObject = result
	return nil
}

func (p *Plasma) unaryOP(context *Context, unaryOperation uint8) *Value {
	unaryName := unaryInstructionsFunctions[unaryOperation]
	target := context.PopObject()
	operation, getError := target.Get(p, context, unaryName)
	if getError != nil {
		return getError
	}
	result, success := p.CallFunction(context, operation)
	if !success {
		return result
	}
	context.LastObject = result
	return nil
}

func (p *Plasma) binaryOPRightHandSide(context *Context, leftHandSide *Value, rightHandSide *Value, rightOperationName string) *Value {
	rightHandSideOperation, getError := rightHandSide.Get(p, context, rightOperationName)
	if getError != nil {
		return getError
	}
	result, success := p.CallFunction(context, rightHandSideOperation, leftHandSide)
	if !success {
		return result
	}
	context.LastObject = result
	return nil
}

func (p *Plasma) binaryOP(context *Context, binaryOperation uint8) *Value {
	binaryNames := binaryInstructionsFunctions[binaryOperation]
	leftHandSide := context.PopObject()
	rightHandSide := context.PopObject()
	leftHandSideOperation, getError := leftHandSide.Get(p, context, binaryNames[0])
	if getError != nil {
		return p.binaryOPRightHandSide(context, leftHandSide, rightHandSide, binaryNames[1])
	}
	result, success := p.CallFunction(context, leftHandSideOperation, rightHandSide)
	if !success {
		return p.binaryOPRightHandSide(context, leftHandSide, rightHandSide, binaryNames[1])
	}
	context.LastObject = result
	return nil
}

func (p *Plasma) methodInvocationOP(context *Context, numberOfArguments int) *Value {
	method := context.PopObject()
	var arguments []*Value
	for index := 0; index < numberOfArguments; index++ {
		arguments = append(arguments, context.PopObject())
	}
	result, success := p.CallFunction(context, method, arguments...)
	if !success {
		return result
	}
	context.LastObject = result
	return nil
}

func (p *Plasma) getIdentifierOP(context *Context, symbol string) *Value {
	result, getError := context.PeekSymbolTable().GetAny(symbol)
	if getError != nil {
		return p.NewObjectWithNameNotFoundError(context, p.ForceMasterGetAny(ValueName), symbol)
	}
	context.LastObject = result
	return nil
}

func (p *Plasma) selectNameFromObjectOP(context *Context, symbol string) *Value {
	source := context.PopObject()
	result, getError := source.Get(p, context, symbol)
	if getError != nil {
		return getError
	}
	context.LastObject = result
	return nil
}

func (p *Plasma) indexOP(context *Context) *Value {
	index := context.PopObject()
	source := context.PopObject()
	result, success := p.IndexCall(context, source, index)
	if !success {
		return result
	}
	context.LastObject = result
	return nil
}

func (p *Plasma) pushOP(context *Context) *Value {
	if context.LastObject != nil {
		context.PushObject(context.LastObject)
		context.LastObject = nil
	}
	return nil
}

func (p *Plasma) assignIdentifierOP(context *Context, symbol string) *Value {
	value := context.PopObject()
	context.PeekSymbolTable().Set(symbol, value)
	return nil
}

func (p *Plasma) newClassOP(context *Context, bytecode *Bytecode, classInformation ClassInformation) *Value {
	bases := context.PopObject()
	body := bytecode.NextN(classInformation.BodyLength)
	result := p.NewType(context, false, classInformation.Name, context.PeekSymbolTable(), bases.Content,
		NewPlasmaConstructor(body),
	)
	context.LastObject = result
	return nil
}

func (p *Plasma) newClassFunctionOP(context *Context, bytecode *Bytecode, functionInformation FunctionInformation) *Value {
	body := bytecode.NextN(functionInformation.BodyLength)
	context.LastObject = p.NewFunction(
		context,
		false,
		context.PeekObject().SymbolTable().Parent,
		NewPlasmaClassFunction(
			context.PeekObject(),
			functionInformation.NumberOfArguments,
			body,
		),
	)
	return nil
}

func (p *Plasma) newFunctionOP(context *Context, bytecode *Bytecode, functionInformation FunctionInformation) *Value {
	body := bytecode.NextN(functionInformation.BodyLength)
	context.LastObject = p.NewFunction(
		context,
		false,
		context.PeekSymbolTable(),
		NewPlasmaFunction(
			functionInformation.NumberOfArguments,
			body,
		),
	)
	return nil
}

func (p *Plasma) loadFunctionArgumentsOP(context *Context, receivers []string) *Value {
	for _, receiver := range receivers {
		context.PeekSymbolTable().Set(receiver, context.PopObject())
	}
	return nil
}

func (p *Plasma) returnOP(context *Context, numberOfResults int) *Value {
	if numberOfResults == 0 {
		return p.GetNone()
	} else if numberOfResults == 1 {
		return context.PopObject()
	}
	resultValues := make([]*Value, numberOfResults)
	for i := 0; i < numberOfResults; i++ {
		resultValues[i] = context.PopObject()
	}
	return p.NewTuple(context, false, resultValues)
}

func (p *Plasma) ifOneLinerOP(context *Context, bytecode *Bytecode, information ConditionInformation) *Value {
	condition := context.PopObject()
	ifBody := bytecode.NextN(information.BodyLength)
	elseBody := bytecode.NextN(information.ElseBodyLength)
	asBool, asBoolError := p.QuickGetBool(context, condition)
	if asBoolError != nil {
		return asBoolError
	}
	var codeToExecute []Code
	if asBool {
		codeToExecute = ifBody
	} else {
		codeToExecute = elseBody
	}
	result, success := p.Execute(context, NewBytecodeFromArray(codeToExecute))
	if !success {
		return result
	}
	context.LastObject = result
	return nil
}

func (p *Plasma) unlessOneLinerOP(context *Context, bytecode *Bytecode, information ConditionInformation) *Value {
	condition := context.PopObject()
	unlessBody := bytecode.NextN(information.BodyLength)
	elseBody := bytecode.NextN(information.ElseBodyLength)
	asBool, asBoolError := p.QuickGetBool(context, condition)
	if asBoolError != nil {
		return asBoolError
	}
	var codeToExecute []Code
	if !asBool {
		codeToExecute = unlessBody
	} else {
		codeToExecute = elseBody
	}
	result, success := p.Execute(context, NewBytecodeFromArray(codeToExecute))
	if !success {
		return result
	}
	context.LastObject = result
	return nil
}

func (p *Plasma) assignSelectorOP(context *Context, symbol string) *Value {
	source := context.PopObject()
	value := context.PopObject()
	source.Set(p, context, symbol, value)
	return nil
}

func (p *Plasma) assignIndexOP(context *Context) *Value {
	index := context.PopObject()
	source := context.PopObject()
	value := context.PopObject()
	assign, getError := source.Get(p, context, Assign)
	if getError != nil {
		return getError
	}
	result, success := p.CallFunction(context, assign, index, value)
	if !success {
		return result
	}
	return nil
}

func (p *Plasma) ifOP(context *Context, bytecode *Bytecode, information ConditionInformation) *Value {
	condition := context.PopObject()
	ifBody := bytecode.NextN(information.BodyLength)
	elseBody := bytecode.NextN(information.ElseBodyLength)
	conditionAsBool, transformationError := p.QuickGetBool(context, condition)
	if transformationError != nil {
		return transformationError
	}
	var codeToExecute []Code
	if conditionAsBool {
		codeToExecute = ifBody
	} else {
		codeToExecute = elseBody
	}
	result, success := p.Execute(context, NewBytecodeFromArray(codeToExecute))
	if !success {
		return result
	}
	context.LastObject = result
	return nil
}

func (p *Plasma) unlessOP(context *Context, bytecode *Bytecode, information ConditionInformation) *Value {
	condition := context.PopObject()
	ifBody := bytecode.NextN(information.BodyLength)
	elseBody := bytecode.NextN(information.ElseBodyLength)
	conditionAsBool, transformationError := p.QuickGetBool(context, condition)
	if transformationError != nil {
		return transformationError
	}
	var codeToExecute []Code
	if !conditionAsBool {
		codeToExecute = ifBody
	} else {
		codeToExecute = elseBody
	}
	result, success := p.Execute(context, NewBytecodeFromArray(codeToExecute))
	if !success {
		return result
	}
	context.LastObject = result
	return nil
}

func (p *Plasma) forLoopOP(context *Context, bytecode *Bytecode, information LoopInformation) *Value {
	source, success := p.InterpretAsIterator(context, context.PopObject())
	if !success {
		return source
	}
	next, nextGetError := source.Get(p, context, Next)
	if nextGetError != nil {
		return nextGetError
	}
	hasNext, hasNextGetError := source.Get(p, context, HasNext)
	if hasNextGetError != nil {
		return hasNextGetError
	}
	body := bytecode.NextN(information.BodyLength)
	numberOfReceivers := len(information.Receivers)
	bodyBytecode := NewBytecodeFromArray(body)
	var (
		doesHasNext *Value
		nextValue   *Value
		result      *Value
	)
loop:
	for {
	continueState:
		// Check  if the iter has a next value
		doesHasNext, success = p.CallFunction(context, hasNext)
		if !success {
			return doesHasNext
		}
		doesHasNextAsBool, boolInterpretationError := p.QuickGetBool(context, doesHasNext)
		if boolInterpretationError != nil {
			return boolInterpretationError
		}
		if !doesHasNextAsBool {
			break
		}

		// Get the value
		nextValue, success = p.CallFunction(context, next)
		if !success {
			return nextValue
		}
	redoState:
		// Unpack the value
		unpackedValues, unpackError := p.UnpackValues(context, nextValue, numberOfReceivers)
		if unpackError != nil {
			return unpackError
		}
		numberOfUnpackedValues := len(unpackedValues)
		if numberOfUnpackedValues != numberOfReceivers {
			return p.NewIndexOutOfRange(context, numberOfUnpackedValues, int64(numberOfReceivers))
		}
		for index, symbol := range information.Receivers {
			context.PeekSymbolTable().Set(symbol, unpackedValues[index])
		}
		// Reset  the bytecode
		bodyBytecode.index = 0
		// Execute the body
		result, success = p.Execute(context, bodyBytecode)
		if !success {
			return result
		}
		switch context.LastState {
		case ReturnState:
			context.LastObject = result
			return nil
		case BreakState:
			break loop
		case ContinueState:
			goto continueState
		case RedoState:
			goto redoState
		}
	}
	context.LastState = NoState
	return nil
}

func (p *Plasma) newGeneratorOP(context *Context, numberOfReceivers int) *Value {
	operation := context.PopObject()
	source := context.PopObject()
	sourceAsIter, interpretationSuccess := p.InterpretAsIterator(context, source)
	if !interpretationSuccess {
		return sourceAsIter
	}
	sourceAsIterHasNext, hasNextError := sourceAsIter.Get(p, context, HasNext)
	if hasNextError != nil {
		return hasNextError
	}
	sourceAsIterNext, nextGetError := sourceAsIter.Get(p, context, Next)
	if nextGetError != nil {
		return nextGetError
	}
	result := p.NewIterator(context, false)
	result.SetOnDemandSymbol(
		HasNext,
		func() *Value {
			return p.NewFunction(context, false, result.SymbolTable(),
				NewBuiltInFunction(0,
					func(self *Value, _ ...*Value) (*Value, bool) {
						return p.CallFunction(context, sourceAsIterHasNext)
					},
				),
			)
		},
	)
	result.SetOnDemandSymbol(
		Next,
		func() *Value {
			return p.NewFunction(context, false, result.SymbolTable(),
				NewBuiltInFunction(0,
					func(self *Value, _ ...*Value) (*Value, bool) {
						nextValues, success := p.CallFunction(context, sourceAsIterNext)
						if !success {
							return nextValues, false
						}
						unpackedValues, unpackError := p.UnpackValues(context, nextValues, numberOfReceivers)
						if unpackError != nil {
							return unpackError, false
						}
						return p.CallFunction(context, operation, unpackedValues...)
					},
				),
			)
		},
	)
	context.LastObject = result
	return nil
}

func (p *Plasma) whileLoopOP(context *Context, bytecode *Bytecode, information LoopInformation) *Value {
	conditionCode := bytecode.NextN(information.ConditionLength)
	conditionBytecode := NewBytecodeFromArray(conditionCode)

	body := bytecode.NextN(information.BodyLength)
	bodyBytecode := NewBytecodeFromArray(body)

	var (
		result *Value
	)
loop:
	for {
	continueState:
		// Reset condition bytecode
		conditionBytecode.index = 0
		// Check the condition
		condition, success := p.Execute(context, conditionBytecode)
		if !success {
			return condition
		}
		// Interpret as boolean
		conditionAsBool, interpretationError := p.QuickGetBool(context, condition)
		if interpretationError != nil {
			return interpretationError
		}
		if !conditionAsBool {
			break
		}
	redoState:
		// Reset  the bytecode
		bodyBytecode.index = 0
		// Execute the body
		result, success = p.Execute(context, bodyBytecode)
		if !success {
			return result
		}
		switch context.LastState {
		case ReturnState:
			context.LastObject = result
			return nil
		case BreakState:
			break loop
		case ContinueState:
			goto continueState
		case RedoState:
			goto redoState
		}
	}
	context.LastState = NoState
	return nil
}

func (p *Plasma) doWhileLoopOP(context *Context, bytecode *Bytecode, information LoopInformation) *Value {
	conditionCode := bytecode.NextN(information.ConditionLength)
	conditionBytecode := NewBytecodeFromArray(conditionCode)

	body := bytecode.NextN(information.BodyLength)
	bodyBytecode := NewBytecodeFromArray(body)

	var (
		condition *Value
	)
loop:
	for {
	redoState:
		// Execute code
		// Reset  the bytecode
		bodyBytecode.index = 0
		// Execute the body
		result, success := p.Execute(context, bodyBytecode)
		if !success {
			return result
		}
		// Check state
		switch context.LastState {
		case ReturnState:
			context.LastObject = result
			return nil
		case BreakState:
			break loop
		case RedoState:
			goto redoState
		}
		// Reset condition bytecode
		conditionBytecode.index = 0
		// Check the condition
		condition, success = p.Execute(context, conditionBytecode)
		if !success {
			return condition
		}
		// Interpret as boolean
		conditionAsBool, interpretationError := p.QuickGetBool(context, condition)
		if interpretationError != nil {
			return interpretationError
		}
		if !conditionAsBool {
			break
		}
	}
	context.LastState = NoState
	return nil
}

func (p *Plasma) untilLoopOP(context *Context, bytecode *Bytecode, information LoopInformation) *Value {
	conditionCode := bytecode.NextN(information.ConditionLength)
	conditionBytecode := NewBytecodeFromArray(conditionCode)

	body := bytecode.NextN(information.BodyLength)
	bodyBytecode := NewBytecodeFromArray(body)

	var (
		result *Value
	)
loop:
	for {
	continueState:
		// Reset condition bytecode
		conditionBytecode.index = 0
		// Check the condition
		condition, success := p.Execute(context, conditionBytecode)
		if !success {
			return condition
		}
		// Interpret as boolean
		conditionAsBool, interpretationError := p.QuickGetBool(context, condition)
		if interpretationError != nil {
			return interpretationError
		}
		if conditionAsBool {
			break
		}
	redoState:
		// Reset  the bytecode
		bodyBytecode.index = 0
		// Execute the body
		result, success = p.Execute(context, bodyBytecode)
		if !success {
			return result
		}
		switch context.LastState {
		case ReturnState:
			context.LastObject = result
			return nil
		case BreakState:
			break loop
		case ContinueState:
			goto continueState
		case RedoState:
			goto redoState
		}
	}
	context.LastState = NoState
	return nil
}
