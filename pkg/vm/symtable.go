package vm

import (
	"github.com/shoriwe/gruby/pkg/errors"
)

type SymbolTable struct {
	parent  *SymbolTable
	symbols map[string]Object
}

func (symbolTable *SymbolTable) Get(symbol string, useParent bool) (Object, *errors.Error) {
	value, ok := symbolTable.symbols[symbol]
	if !ok {
		if symbolTable.parent == nil || !useParent {
			return nil, NewAttributeNotFound("", symbol)
		}
		var parentFoundError *errors.Error
		value, parentFoundError = symbolTable.parent.Get(symbol, true)
		if parentFoundError != nil {
			return nil, parentFoundError
		}
	}
	return value, nil
}

func (symbolTable *SymbolTable) Set(symbol string, obj Object) *errors.Error {
	if obj == nil {
		return errors.New(UnknownLine, "received a nil pointer as Object", NilObjectError)
	}
	symbolTable.symbols[symbol] = obj
	return nil
}

func NewSymbolTable(parentSymbolTable *SymbolTable) *SymbolTable {
	return &SymbolTable{
		parent:  parentSymbolTable,
		symbols: map[string]Object{},
	}
}
