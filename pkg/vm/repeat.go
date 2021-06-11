package vm

func (p *Plasma) Repeat(content []IObject, times int) ([]IObject, *Object) {
	copyFunctions := map[int64]*Function{}
	var result []IObject
	if times > 0 {
		for _, object := range content {
			copyObject, getError := object.Get(Copy)
			if getError != nil {
				copyFunctions[object.Id()] = nil
				continue
			}
			if _, ok := copyObject.(*Function); !ok {
				return nil, p.NewInvalidTypeError(copyObject.TypeName(), FunctionName)
			}
			copyFunctions[object.Id()] = copyObject.(*Function)
		}
	}
	for i := 0; i < times; i++ {
		for _, object := range content {
			copyFunction := copyFunctions[object.Id()]
			if copyFunction == nil {
				result = append(result, object)
				continue
			}
			objectCopy, callError := p.CallFunction(copyFunction, p.PeekSymbolTable())
			if callError != nil {
				return nil, callError
			}
			result = append(result, objectCopy)
		}
	}
	return result, nil
}
