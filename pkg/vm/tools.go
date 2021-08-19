package vm

func (p *Plasma) Repeat(context *Context, content []*Value, times int64) ([]*Value, *Value) {
	copyFunctions := map[int64]*Value{}
	var result []*Value
	if times == 0 {
		for _, object := range content {
			copyObject, getError := object.Get(Copy)
			if getError != nil {
				copyFunctions[object.Id()] = nil
				continue
			}
			copyFunctions[object.Id()] = copyObject
		}
	}
	for i := int64(0); i < times; i++ {
		for _, object := range content {
			copyFunction := copyFunctions[object.Id()]
			if copyFunction == nil {
				result = append(result, object)
				continue
			}
			objectCopy, success := p.CallFunction(context, copyFunction)
			if !success {
				return nil, objectCopy
			}
			result = append(result, objectCopy)
		}
	}
	return result, nil
}

func (p *Plasma) Equals(context *Context, leftHandSide *Value, rightHandSide *Value) (bool, *Value) {
	equals, getError := leftHandSide.Get(Equals)
	if getError != nil {
		// Try with the rightHandSide
		var rightEquals *Value
		rightEquals, getError = rightHandSide.Get(RightEquals)
		if getError != nil {
			return false, p.NewObjectWithNameNotFoundError(p.builtInContext, rightHandSide.GetClass(p), Equals)
		}
		result, success := p.CallFunction(context, rightEquals, rightHandSide)
		if !success {
			return false, result
		}
		return p.QuickGetBool(context, result)
	}
	result, success := p.CallFunction(context, equals, rightHandSide)
	if !success {
		// Try with the rightHandSide
		var rightEquals *Value
		rightEquals, getError = rightHandSide.Get(RightEquals)
		if getError != nil {
			return false, p.NewObjectWithNameNotFoundError(p.builtInContext, rightHandSide.GetClass(p), Equals)
		}
		result, success = p.CallFunction(context, rightEquals, rightHandSide)
		if !success {
			return false, result
		}
	}
	return p.QuickGetBool(context, result)
}

func (p *Plasma) InterpretAsBool(expression bool) *Value {
	if expression {
		return p.GetTrue()
	}
	return p.GetFalse()
}
