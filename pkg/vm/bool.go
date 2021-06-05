package vm

import (
	"fmt"
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

func (p *Plasma) BoolInitialize(object IObject) *Object {
	object.Set(Equals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
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
				func(self IObject, arguments ...IObject) (IObject, *Object) {
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
				func(self IObject, arguments ...IObject) (IObject, *Object) {
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
				func(self IObject, arguments ...IObject) (IObject, *Object) {
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
				func(self IObject, _ ...IObject) (IObject, *Object) {

					if self.GetHash() == 0 {
						boolHash := p.HashString(fmt.Sprintf("%t-%s", self.GetBool(), BoolName))
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
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewBool(p.PeekSymbolTable(), self.GetBool()), nil
				},
			),
		),
	)
	object.Set(ToInteger,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
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
				func(self IObject, _ ...IObject) (IObject, *Object) {
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
				func(self IObject, _ ...IObject) (IObject, *Object) {

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
				func(self IObject, _ ...IObject) (IObject, *Object) {

					return p.NewBool(p.PeekSymbolTable(), self.GetBool()), nil
				},
			),
		),
	)
	return nil
}
