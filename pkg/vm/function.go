package vm

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/errors"
)

type Function struct {
	*Object
	Callable Callable
}

func (p *Plasma) NewFunction(parentSymbols *SymbolTable, callable Callable) *Function {
	function := &Function{
		Object: &Object{
			id:         p.NextId(),
			typeName:   FunctionName,
			subClasses: nil,
			symbols:    NewSymbolTable(parentSymbols),
		},
		Callable: callable,
	}
	function.Set(ToString, &Function{
		Object: nil,
		Callable: NewBuiltInClassFunction(function, 0,
			func(self IObject, _ ...IObject) (IObject, *errors.Error) {
				return p.NewString(p.PeekSymbolTable(), fmt.Sprintf("Function-%d", function.Id())), nil
			},
		),
	})
	return function
}
