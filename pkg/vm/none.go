package vm

import "github.com/shoriwe/gruby/pkg/errors"

func (p *Plasma) NoneInitialize(object IObject) *errors.Error {
	object.Set(ToString,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(_ IObject, _ ...IObject) (IObject, *Object) {
					return p.NewString(p.PeekSymbolTable(), "None"), nil
				},
			),
		),
	)
	object.Set(ToBool,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(_ IObject, _ ...IObject) (IObject, *Object) {
					return p.NewBool(p.PeekSymbolTable(), false), nil
				},
			),
		),
	)
	return nil
}
