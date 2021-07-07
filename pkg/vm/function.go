package vm

import (
	"fmt"
)

type Function struct {
	*Object
	Callable Callable
}

func (p *Plasma) FunctionInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object Value) *Object {
		object.SetOnDemandSymbol(Hash,
			func() Value {
				return &Function{
					Object: &Object{
						id:         p.NextId(),
						typeName:   FunctionName,
						subClasses: nil,
						symbols:    NewSymbolTable(object.SymbolTable()),
						isBuiltIn:  isBuiltIn,
					},
					Callable: NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewInteger(context, false, context.PeekSymbolTable(), object.Id()), nil
						},
					),
				}
			},
		)
		object.SetOnDemandSymbol(Equals,
			func() Value {
				return &Function{
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
				}
			},
		)
		object.SetOnDemandSymbol(RightEquals,
			func() Value {
				return &Function{
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
				}
			},
		)
		object.SetOnDemandSymbol(ToString,
			func() Value {
				return &Function{
					Object: &Object{
						id:         p.NextId(),
						typeName:   FunctionName,
						subClasses: nil,
						symbols:    NewSymbolTable(object.SymbolTable()),
						isBuiltIn:  isBuiltIn,
					},
					Callable: NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewString(context, false, context.PeekSymbolTable(), fmt.Sprintf("Function-%d", object.Id())), nil
						},
					),
				}
			},
		)
		return nil
	}
}

func (p *Plasma) NewFunction(context *Context, isBuiltIn bool, parentSymbols *SymbolTable, callable Callable) *Function {
	function := &Function{
		Object:   p.NewObject(context, isBuiltIn, FunctionName, nil, parentSymbols),
		Callable: callable,
	}
	function.SetOnDemandSymbol(Self,
		func() Value {
			return function
		},
	)
	p.FunctionInitialize(isBuiltIn)(context, function)
	return function
}
