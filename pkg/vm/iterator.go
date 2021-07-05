package vm

type Iterator struct {
	*Object
}

func (p *Plasma) NewIterator(isBuiltIn bool, parentSymbols *SymbolTable) *Iterator {
	iterator := &Iterator{
		Object: p.NewObject(isBuiltIn, IteratorName, nil, parentSymbols),
	}
	p.IteratorInitialize(isBuiltIn)(iterator)
	iterator.Set(Self, iterator)
	return iterator
}

func (p *Plasma) IteratorInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(object Value) *Object {
		object.SetOnDemandSymbol(HasNext,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
