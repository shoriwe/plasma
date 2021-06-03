package vm

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/errors"
)

type Bool struct {
	*Object
}

func (p *Plasma) NewBool(parentSymbols *SymbolTable, value bool) *Bool {
	bool_ := &Bool{
		Object: p.NewObject(BoolName, nil, parentSymbols),
	}
	bool_.SetBool(value)
	p.BoolInitialize(bool_)
	return bool_
}

func (p *Plasma) BoolInitialize(object IObject) *errors.Error {
	object.Set(Equals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					if _, ok := right.(*Bool); !ok {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					return p.NewBool(p.PeekSymbolTable(), self.GetBool() == right.GetBool()), nil
				},
			),
		),
	)
	object.Set(RightEquals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					if _, ok := left.(*Bool); !ok {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					return p.NewBool(p.PeekSymbolTable(), left.GetBool() == self.GetBool()), nil
				},
			),
		),
	)
	object.Set(NotEquals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					if _, ok := right.(*Bool); !ok {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					return p.NewBool(p.PeekSymbolTable(), self.GetBool() != right.GetBool()), nil
				},
			),
		),
	)
	object.Set(RightNotEquals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					if _, ok := left.(*Bool); !ok {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					return p.NewBool(p.PeekSymbolTable(), left.GetBool() != self.GetBool()), nil
				},
			),
		),
	)
	object.Set(Copy,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {

					if self.GetHash() == 0 {
						boolHash, hashingError := p.HashString(fmt.Sprintf("%t-%s", self.GetBool(), BoolName))
						if hashingError != nil {
							return nil, hashingError
						}
						self.SetHash(boolHash)
					}
					return p.NewInteger(p.PeekSymbolTable(), self.GetHash()), nil
				},
			),
		),
	)
	object.Set(Copy,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewBool(p.PeekSymbolTable(), self.GetBool()), nil
				},
			),
		),
	)
	object.Set(ToInteger,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					if self.GetBool() {
						return p.NewInteger(p.PeekSymbolTable(), 1), nil
					}
					return p.NewInteger(p.PeekSymbolTable(), 0), nil
				},
			),
		),
	)
	object.Set(ToFloat,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					if self.GetBool() {
						return p.NewFloat(p.PeekSymbolTable(), 1), nil
					}
					return p.NewFloat(p.PeekSymbolTable(), 0), nil
				},
			),
		),
	)
	object.Set(ToString,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {

					if self.GetBool() {
						return p.NewString(p.PeekSymbolTable(), TrueName), nil
					}
					return p.NewString(p.PeekSymbolTable(), FalseName), nil
				},
			),
		),
	)
	object.Set(ToBool,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {

					return p.NewBool(p.PeekSymbolTable(), self.GetBool()), nil
				},
			),
		),
	)
	return nil
}
