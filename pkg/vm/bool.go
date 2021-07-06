package vm

type Bool struct {
	*Object
}

func (p *Plasma) QuickGetBool(value Value) (bool, *Object) {
	if _, ok := value.(*Bool); ok {
		return value.GetBool(), nil
	}
	valueToBool, getError := value.Get(ToBool)
	if getError != nil {
		return false, p.NewObjectWithNameNotFoundError(value.GetClass(p), ToBool)
	}
	valueBool, callError := p.CallFunction(valueToBool, valueToBool.SymbolTable().Parent)
	if callError != nil {
		return false, callError
	}
	if _, ok := valueBool.(*Bool); !ok {
		return false, p.NewInvalidTypeError(value.TypeName(), BoolName)
	}
	return valueBool.GetBool(), nil
}

func (p *Plasma) NewBool(isBuiltIn bool, parentSymbols *SymbolTable, value bool) *Bool {
	bool_ := &Bool{
		Object: p.NewObject(isBuiltIn, BoolName, nil, parentSymbols),
	}
	bool_.SetBool(value)
	p.BoolInitialize(isBuiltIn)(bool_)
	bool_.SetOnDemandSymbol(Self,
		func() Value {
			return bool_
		},
	)
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
		object.SetOnDemandSymbol(Equals,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(RightEquals,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(NotEquals,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(RightNotEquals,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(Copy,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							if self.GetBool() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Hash,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							if self.GetBool() {
								return p.NewInteger(false, p.PeekSymbolTable(), 1), nil
							}
							return p.NewInteger(false, p.PeekSymbolTable(), 0), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToInteger,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							if self.GetBool() {
								return p.NewInteger(false, p.PeekSymbolTable(), 1), nil
							}
							return p.NewInteger(false, p.PeekSymbolTable(), 0), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToFloat,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							if self.GetBool() {
								return p.NewFloat(false, p.PeekSymbolTable(), 1), nil
							}
							return p.NewFloat(false, p.PeekSymbolTable(), 0), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToString,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							if self.GetBool() {
								return p.NewString(false, p.PeekSymbolTable(), TrueName), nil
							}
							return p.NewString(false, p.PeekSymbolTable(), FalseName), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToBool,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							if self.GetBool() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		return nil
	}
}
