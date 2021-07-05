package vm

type Bool struct {
	*Object
}

func (p *Plasma) NewBool(isBuiltIn bool, parentSymbols *SymbolTable, value bool) *Bool {
	bool_ := &Bool{
		Object: p.NewObject(isBuiltIn, BoolName, nil, parentSymbols),
	}
	bool_.SetBool(value)
	p.BoolInitialize(isBuiltIn)(bool_)
	bool_.Set(Self, bool_)
	return bool_
}

func (p *Plasma) GetFalse() *Bool {
	return p.ForceMasterGetAny(FalseName).(*Bool)
}

func (p *Plasma) GetTrue() *Bool {
	return p.ForceMasterGetAny(TrueName).(*Bool)
}

func (p *Plasma) BoolInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(object Value) *Object {
		object.Set(Equals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						if _, ok := right.(*Bool); !ok {
							return p.GetFalse(), nil
						}
						if self.GetBool() == right.GetBool() {
							return p.GetTrue(), nil
						}
						return p.GetFalse(), nil
					},
				),
			),
		)
		object.Set(RightEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						if _, ok := left.(*Bool); !ok {
							return p.GetFalse(), nil
						}
						if left.GetBool() == self.GetBool() {
							return p.GetTrue(), nil
						}
						return p.GetFalse(), nil
					},
				),
			),
		)
		object.Set(NotEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						if _, ok := right.(*Bool); !ok {
							return p.GetFalse(), nil
						}
						if self.GetBool() != right.GetBool() {
							return p.GetTrue(), nil
						}
						return p.GetFalse(), nil
					},
				),
			),
		)
		object.Set(RightNotEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						if _, ok := left.(*Bool); !ok {
							return p.GetFalse(), nil
						}
						if left.GetBool() != self.GetBool() {
							return p.GetTrue(), nil
						}
						return p.GetFalse(), nil
					},
				),
			),
		)
		object.Set(Copy,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						if self.GetBool() {
							return p.GetTrue(), nil
						}
						return p.GetFalse(), nil
					},
				),
			),
		)
		object.Set(Hash,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						if self.GetBool() {
							return p.NewInteger(false, p.PeekSymbolTable(), 1), nil
						}
						return p.NewInteger(false, p.PeekSymbolTable(), 0), nil
					},
				),
			),
		)
		object.Set(ToInteger,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
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
					func(self Value, _ ...Value) (Value, *Object) {
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
					func(self Value, _ ...Value) (Value, *Object) {
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
					func(self Value, _ ...Value) (Value, *Object) {
						if self.GetBool() {
							return p.GetTrue(), nil
						}
						return p.GetFalse(), nil
					},
				),
			),
		)
		return nil
	}
}
