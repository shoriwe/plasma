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
	for source := symbolTable; source != nil; source = source.Parent {
		result, found := source.Symbols[symbol]
		if found {
			return result, nil
		}
	}
	return nil, errors.NewNameNotFoundError()
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
