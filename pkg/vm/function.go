package vm

import (
	"fmt"
)

type Function struct {
	*Object
	Callable Callable
}

func (p *Plasma) NewFunction(isBuiltIn bool, parentSymbols *SymbolTable, callable Callable) *Function {
	function := &Function{
		Object: &Object{
			id:         p.NextId(),
			typeName:   FunctionName,
			subClasses: nil,
			symbols:    NewSymbolTable(parentSymbols),
			isBuiltIn:  isBuiltIn,
		},
		Callable: callable,
	}
	function.Set(Hash, &Function{
		Object: nil,
		Callable: NewBuiltInClassFunction(function, 0,
			func(self IObject, _ ...IObject) (IObject, *Object) {
				return p.NewInteger(false, p.PeekSymbolTable(), int64(self.Id())), nil
			},
		),
	})
	function.Set(ToString, &Function{
		Object: nil,
		Callable: NewBuiltInClassFunction(function, 0,
			func(self IObject, _ ...IObject) (IObject, *Object) {
				return p.NewString(false, p.PeekSymbolTable(), fmt.Sprintf("Function-%d", function.Id())), nil
			},
		),
	})
	return function
}
