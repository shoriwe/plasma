package vm

func (p *Plasma) NewNone() *Object {
	return p.ForceConstruction(p.ForceMasterGetAny(NoneName)).(*Object)
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
						return p.NewBool(false, p.PeekSymbolTable(), false), nil
					},
				),
			),
		)
		return nil
	}
}
