package vm

import (
	"fmt"
)

type Bool struct {
	*Object
}

func (p *Plasma) NewBool(isBuiltIn bool, parentSymbols *SymbolTable, value bool) *Bool {
	bool_ := &Bool{
		Object: p.NewObject(false, BoolName, nil, parentSymbols),
	}
	bool_.SetBool(value)
	p.BoolInitialize(isBuiltIn)(bool_)
	return bool_
}

func (p *Plasma) BoolInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(object IObject) *Object {
		object.Set(Equals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						right := arguments[0]
						if _, ok := right.(*Bool); !ok {
							return p.NewBool(false, p.PeekSymbolTable(), false), nil
						}
						return p.NewBool(false, p.PeekSymbolTable(), self.GetBool() == right.GetBool()), nil
					},
				),
			),
		)
		object.Set(RightEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						left := arguments[0]
						if _, ok := left.(*Bool); !ok {
							return p.NewBool(false, p.PeekSymbolTable(), false), nil
						}
						return p.NewBool(false, p.PeekSymbolTable(), left.GetBool() == self.GetBool()), nil
					},
				),
			),
		)
		object.Set(NotEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						right := arguments[0]
						if _, ok := right.(*Bool); !ok {
							return p.NewBool(false, p.PeekSymbolTable(), false), nil
						}
						return p.NewBool(false, p.PeekSymbolTable(), self.GetBool() != right.GetBool()), nil
					},
				),
			),
		)
		object.Set(RightNotEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						left := arguments[0]
						if _, ok := left.(*Bool); !ok {
							return p.NewBool(false, p.PeekSymbolTable(), false), nil
						}
						return p.NewBool(false, p.PeekSymbolTable(), left.GetBool() != self.GetBool()), nil
					},
				),
			),
		)
		object.Set(Copy,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self IObject, _ ...IObject) (IObject, *Object) {

						if self.GetHash() == 0 {
							boolHash := p.HashString(fmt.Sprintf("%t-%s", self.GetBool(), BoolName))
							self.SetHash(boolHash)
						}
						return p.NewInteger(false, p.PeekSymbolTable(), self.GetHash()), nil
					},
				),
			),
		)
		object.Set(Copy,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self IObject, _ ...IObject) (IObject, *Object) {
						return p.NewBool(false, p.PeekSymbolTable(), self.GetBool()), nil
					},
				),
			),
		)
		object.Set(ToInteger,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self IObject, _ ...IObject) (IObject, *Object) {
						if self.GetBool() {
							return p.NewInteger(false, p.PeekSymbolTable(), 1), nil
						}
						return p.NewInteger(false, p.PeekSymbolTable(), 0), nil
					},
				),
			),
		)
		object.Set(ToFloat,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self IObject, _ ...IObject) (IObject, *Object) {
						if self.GetBool() {
							return p.NewFloat(false, p.PeekSymbolTable(), 1), nil
						}
						return p.NewFloat(false, p.PeekSymbolTable(), 0), nil
					},
				),
			),
		)
		object.Set(ToString,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self IObject, _ ...IObject) (IObject, *Object) {

						if self.GetBool() {
							return p.NewString(false, p.PeekSymbolTable(), TrueName), nil
						}
						return p.NewString(false, p.PeekSymbolTable(), FalseName), nil
					},
				),
			),
		)
		object.Set(ToBool,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self IObject, _ ...IObject) (IObject, *Object) {

						return p.NewBool(false, p.PeekSymbolTable(), self.GetBool()), nil
					},
				),
			),
		)
		return nil
	}
}
