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
	var content []*Value
	for index := 0; index < length; index++ {
		content = append(content, context.PopObject())
	}
	context.LastObject = p.NewArray(context, false, content)
	return nil
}

func (p *Plasma) newTupleOP(context *Context, length int) *Value {
	var content []*Value
	for index := 0; index < length; index++ {
		content = append(content, context.PopObject())
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
		return getError
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
