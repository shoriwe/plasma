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
				p.NewNotImplementedCallable(HasNext, 0),
			),
		)
		object.Set(Next,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				p.NewNotImplementedCallable(Next, 0),
			),
		)
		return nil
	}
}
