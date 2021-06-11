package vm

type Iterator struct {
	*Object
}

func (p *Plasma) NewIterator(parentSymbols *SymbolTable) *Iterator {
	iterator := &Iterator{
		Object: p.NewObject(IteratorName, nil, parentSymbols),
	}
	return iterator
}

func (p *Plasma) IteratorInitialize(object IObject) *Object {
	object.Set(HasNext,
		p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(HasNext, 0),
		),
	)
	object.Set(Next,
		p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(Next, 0),
		),
	)
	return nil
}
