package vm

import "github.com/shoriwe/gruby/pkg/errors"

type Iterator struct {
	*Object
}

func (p *Plasma) NewIterator(parentSymbols *SymbolTable) *Iterator {
	iterator := &Iterator{
		Object: p.NewObject(IteratorName, nil, parentSymbols),
	}
	return iterator
}

func (p *Plasma) IteratorInitialize(object IObject) *errors.Error {
	object.Set(HasNext,
		p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(0),
		),
	)
	object.Set(Next,
		p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(0),
		),
	)
	return nil
}
