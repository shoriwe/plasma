package vm

func (p *Plasma) GetNone() *Object {
	return p.ForceMasterGetAny(None).(*Object)
}

func (p *Plasma) NewNone(context *Context, isBuiltIn bool, parent *SymbolTable) *Object {
	result := p.NewObject(context, isBuiltIn, NoneName, nil, parent)
	p.NoneInitialize(isBuiltIn)(context, result)
	result.SetOnDemandSymbol(Self,
		func() Value {
			return result
		},
	)
	return result
}

func (p *Plasma) NoneInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object Value) *Object {
		object.SetOnDemandSymbol(Equals,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(_ Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							if right.GetClass(p) == p.ForceMasterGetAny(NoneName).(*Type) {
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
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(_ Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							if left.GetClass(p) == p.ForceMasterGetAny(NoneName).(*Type) {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToString,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(_ Value, _ ...Value) (Value, *Object) {
							return p.NewString(context, false, context.PeekSymbolTable(), "None"), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToBool,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(_ Value, _ ...Value) (Value, *Object) {
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		return nil
	}
}
