package vm

type Iterator struct {
	*Object
}

func (p *Plasma) NewIterator(context *Context, isBuiltIn bool, parentSymbols *SymbolTable) *Iterator {
	iterator := &Iterator{
		Object: p.NewObject(context, isBuiltIn, IteratorName, nil, parentSymbols),
	}
	p.IteratorInitialize(isBuiltIn)(context, iterator)
	iterator.SetOnDemandSymbol(Self,
		func() Value {
			return iterator
		},
	)
	return iterator
}

func (p *Plasma) IteratorInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object Value) *Object {
		object.SetOnDemandSymbol(HasNext,
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
		object.SetOnDemandSymbol(Next,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(_ Value, _ ...Value) (Value, *Object) {
							return p.GetNone(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Iter,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return self, nil
						},
					),
				)
			},
		)
		return nil
	}
}
