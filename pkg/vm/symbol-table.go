package vm

import (
	"fmt"
	"sync"
)

var (
	SymbolNotFoundError = fmt.Errorf("symbol not found")
)

type (
	Symbols struct {
		mutex  *sync.Mutex
		Values map[string]*Value
		Parent *Symbols
	}
)

func NewSymbols(parent *Symbols) *Symbols {
	return &Symbols{
		mutex:  &sync.Mutex{},
		Values: map[string]*Value{},
		Parent: parent,
	}
}

func (symbols *Symbols) Set(name string, value *Value) {
	symbols.mutex.Lock()
	defer symbols.mutex.Unlock()
	symbols.Values[name] = value
}

func (symbols *Symbols) Get(name string) (*Value, error) {
	symbols.mutex.Lock()
	defer symbols.mutex.Unlock()
	var (
		value *Value
		found bool
	)
	value, found = symbols.Values[name]
	if found {
		return value, nil
	}
	for current := symbols.Parent; current != nil; current = current.Parent {
		value, found = current.Values[name]
		if found {
			return value, nil
		}
	}
	return nil, SymbolNotFoundError
}

func (symbols *Symbols) Del(name string) error {
	symbols.mutex.Lock()
	defer symbols.mutex.Unlock()
	_, found := symbols.Values[name]
	if !found {
		return SymbolNotFoundError
	}
	delete(symbols.Values, name)
	return nil
}
