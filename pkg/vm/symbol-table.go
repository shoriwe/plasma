package vm

import (
	"fmt"
	"sync"
)

var (
	SymbolNotFoundError = "symbol not found: %s"
)

type (
	Symbols struct {
		mutex  *sync.Mutex
		values map[string]*Value
		call   *Symbols
		Parent *Symbols
	}
)

/*
NewSymbols creates a new symbol table
*/
func NewSymbols(parent *Symbols) *Symbols {
	return &Symbols{
		mutex:  &sync.Mutex{},
		values: map[string]*Value{},
		call:   nil,
		Parent: parent,
	}
}

/*
Set Assigns a value to a symbol
*/
func (symbols *Symbols) Set(name string, value *Value) {
	symbols.mutex.Lock()
	defer symbols.mutex.Unlock()
	symbols.values[name] = value
}

/*
Get retrieves a value based on the symbol
*/
func (symbols *Symbols) Get(name string) (*Value, error) {
	symbols.mutex.Lock()
	defer symbols.mutex.Unlock()
	var (
		value *Value
		found bool
	)
	value, found = symbols.values[name]
	if found {
		return value, nil
	}
	for current := symbols.Parent; current != nil; current = current.Parent {
		value, found = current.values[name]
		if found {
			return value, nil
		}
	}
	return nil, fmt.Errorf(SymbolNotFoundError, name)
}

/*
Del deletes a symbol
*/
func (symbols *Symbols) Del(name string) error {
	symbols.mutex.Lock()
	defer symbols.mutex.Unlock()
	_, found := symbols.values[name]
	if !found {
		return fmt.Errorf(SymbolNotFoundError, name)
	}
	delete(symbols.values, name)
	return nil
}
