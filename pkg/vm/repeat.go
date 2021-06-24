package vm

func (p *Plasma) Repeat(content []Value, times int64) ([]Value, *Object) {
	copyFunctions := map[int64]Value{}
	var result []Value
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
			objectCopy, callError := p.CallFunction(copyFunction, p.PeekSymbolTable())
			if callError != nil {
				return nil, callError
			}
			result = append(result, objectCopy)
		}
	}
	return result, nil
}
