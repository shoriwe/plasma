package vm

import (
	"fmt"
	"math/big"
)

func (p *Plasma) ModuleInitialize(isBuiltIn bool) ConstructorCallBack {
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
					return p.NewInteger(false, p.PeekSymbolTable(), big.NewInt(object.Id())), nil
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
						return p.NewBool(false, p.PeekSymbolTable(), object.Id() == arguments[0].Id()), nil
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
						return p.NewBool(false, p.PeekSymbolTable(), object.Id() == arguments[0].Id()), nil
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
						return p.NewString(false, p.PeekSymbolTable(), fmt.Sprintf("Module-%d", object.Id())), nil
					},
				),
			},
		)
		return nil
	}
}
func (p *Plasma) NewModule(isBuiltIn bool, parent *SymbolTable) Value {
	module := &Object{
		isBuiltIn:  isBuiltIn,
		id:         p.NextId(),
		typeName:   ModuleName,
		class:      nil,
		subClasses: nil,
		symbols:    NewSymbolTable(parent),
	}
	module.Set(Self, module)
	p.ModuleInitialize(isBuiltIn)(module)
	return module
}
