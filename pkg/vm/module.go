package vm

import (
	"fmt"
)

func (p *Plasma) ModuleInitialize(isBuiltIn bool) ConstructorCallBack {
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
						id:              p.NextId(),
						typeName:        FunctionName,
						subClasses:      nil,
						symbols:         NewSymbolTable(object.SymbolTable()),
						isBuiltIn:       isBuiltIn,
						onDemandSymbols: map[string]OnDemandLoader{},
					},
					Callable: NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewString(context, false, context.PeekSymbolTable(), fmt.Sprintf("Module-%d", object.Id())), nil
						},
					),
				}
			},
		)
		return nil
	}
}
func (p *Plasma) NewModule(context *Context, isBuiltIn bool, parent *SymbolTable) Value {
	module := p.NewObject(context, isBuiltIn, ModuleName, nil, parent)
	module.SetOnDemandSymbol(Self,
		func() Value {
			return module
		},
	)
	p.ModuleInitialize(isBuiltIn)(context, module)
	return module
}
