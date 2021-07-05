package vm

import (
	"fmt"
)

type Function struct {
	*Object
	Callable Callable
}

func (p *Plasma) FunctionInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(object Value) *Object {
		object.Set(Hash, &Function{
			Object: &Object{
				id:         p.NextId(),
				typeName:   FunctionName,
				subClasses: nil,
				symbols:    NewSymbolTable(object.SymbolTable()),
				isBuiltIn:  isBuiltIn,
			},
			Callable: NewBuiltInClassFunction(object, 0,
				func(self Value, _ ...Value) (Value, *Object) {
					return p.NewInteger(false, p.PeekSymbolTable(), object.Id()), nil
				},
			),
		})
		object.Set(Equals,
			&Function{
				Object: &Object{
					id:         p.NextId(),
					typeName:   FunctionName,
					subClasses: nil,
					symbols:    NewSymbolTable(object.SymbolTable()),
					isBuiltIn:  isBuiltIn,
				},
				Callable: NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						if object.Id() == arguments[0].Id() {
							return p.GetTrue(), nil
						}
						return p.GetFalse(), nil
					},
				),
			},
		)
		object.Set(RightEquals,
			&Function{
				Object: &Object{
					id:         p.NextId(),
					typeName:   FunctionName,
					subClasses: nil,
					symbols:    NewSymbolTable(object.SymbolTable()),
					isBuiltIn:  isBuiltIn,
				},
				Callable: NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						if object.Id() == arguments[0].Id() {
							return p.GetTrue(), nil
						}
						return p.GetFalse(), nil
					},
				),
			},
		)
		object.Set(ToString,
			&Function{
				Object: &Object{
					id:         p.NextId(),
					typeName:   FunctionName,
					subClasses: nil,
					symbols:    NewSymbolTable(object.SymbolTable()),
					isBuiltIn:  isBuiltIn,
				},
				Callable: NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewString(false, p.PeekSymbolTable(), fmt.Sprintf("Function-%d", object.Id())), nil
					},
				),
			},
		)
		return nil
	}
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
	function.Set(Self, function)
	p.FunctionInitialize(isBuiltIn)(function)
	return function
}
