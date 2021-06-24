package vm

func (p *Plasma) GetNone() *Object {
	return p.ForceMasterGetAny(None).(*Object)
}

func (p *Plasma) NewNone(isBuiltIn bool, parent *SymbolTable) *Object {
	result := p.NewObject(isBuiltIn, NoneName, nil, parent)
	p.NoneInitialize(isBuiltIn)(result)
	return result
}

func (p *Plasma) NoneInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(object Value) *Object {
		object.Set(Equals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(_ Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						return p.NewBool(false, p.PeekSymbolTable(), right.GetClass(p) == p.ForceMasterGetAny(NoneName).(*Type)), nil
					},
				),
			),
		)
		object.Set(NotEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(_ Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						return p.NewBool(false, p.PeekSymbolTable(), left.GetClass(p) == p.ForceMasterGetAny(NoneName).(*Type)), nil
					},
				),
			),
		)
		object.Set(ToString,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(_ Value, _ ...Value) (Value, *Object) {
						return p.NewString(false, p.PeekSymbolTable(), "None"), nil
					},
				),
			),
		)
		object.Set(ToBool,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(_ Value, _ ...Value) (Value, *Object) {
						return p.GetFalse(), nil
					},
				),
			),
		)
		return nil
	}
}
