package vm

import (
	"github.com/shoriwe/gruby/pkg/errors"
)

type SymbolTable struct {
	Parent  *SymbolTable
	Symbols map[string]IObject
}

func (symbolTable *SymbolTable) Set(s string, object IObject) {
	symbolTable.Symbols[s] = object
}

func (symbolTable *SymbolTable) GetSelf(symbol string) (IObject, *errors.Error) {
	result, found := symbolTable.Symbols[symbol]
	if !found {
		return nil, errors.NewNameNotFoundError()
	}
	return result, nil
}

func (symbolTable *SymbolTable) GetAny(symbol string) (IObject, *errors.Error) {
	result, found := symbolTable.Symbols[symbol]
	if !found {
		var foundError *errors.Error
		if symbolTable.Parent != nil {
			result, foundError = symbolTable.Parent.GetAny(symbol)
			if foundError != nil {
				return nil, foundError
			}
		} else {
			return nil, errors.NewNameNotFoundError()
		}
	}
	return result, nil
}

func (symbolTable *SymbolTable) Update(newEntries map[string]IObject) {
	for key, value := range newEntries {
		symbolTable.Symbols[key] = value
	}
}

func NewSymbolTable(parentSymbols *SymbolTable) *SymbolTable {
	return &SymbolTable{
		Parent:  parentSymbols,
		Symbols: map[string]IObject{},
	}
}
