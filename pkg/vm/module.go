package vm

import (
	"fmt"
)

func (p *Plasma) ModuleInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(object Value) *Object {
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
							return p.NewInteger(false, p.PeekSymbolTable(), object.Id()), nil
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
							return p.NewString(false, p.PeekSymbolTable(), fmt.Sprintf("Module-%d", object.Id())), nil
						},
					),
				}
			},
		)
		return nil
	}
}
func (p *Plasma) NewModule(isBuiltIn bool, parent *SymbolTable) Value {
	module := &Object{
		isBuiltIn:       isBuiltIn,
		id:              p.NextId(),
		typeName:        ModuleName,
		class:           nil,
		subClasses:      nil,
		symbols:         NewSymbolTable(parent),
		onDemandSymbols: map[string]OnDemandLoader{},
	}
	module.Set(Self, module)
	p.ModuleInitialize(isBuiltIn)(module)
	return module
}
