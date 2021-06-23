package vm

import "math/big"

func (p *Plasma) Repeat(content []Value, times *big.Int) ([]Value, *Object) {
	copyFunctions := map[int64]Value{}
	var result []Value
	if times.Cmp(big.NewInt(0)) == 1 {
		for _, object := range content {
			copyObject, getError := object.Get(Copy)
			if getError != nil {
				copyFunctions[object.Id()] = nil
				continue
			}
			copyFunctions[object.Id()] = copyObject
		}
	}
	for i := big.NewInt(0); i.Cmp(times) == -1; i.Add(i, big.NewInt(1)) {
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
