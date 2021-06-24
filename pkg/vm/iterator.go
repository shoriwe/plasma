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
		object.Set(HasNext,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(_ Value, _ ...Value) (Value, *Object) {
						return p.NewBool(false, p.PeekSymbolTable(), false), nil
					},
				),
			),
		)
		object.Set(Next,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(_ Value, _ ...Value) (Value, *Object) {
						return p.NewNone(), nil
					},
				),
			),
		)
		object.Set(Iter,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return self, nil
					},
				),
			),
		)
		return nil
	}
}
