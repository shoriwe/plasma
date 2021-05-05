package utils

import "github.com/shoriwe/gruby/pkg/errors"

type SymbolTable struct {
	Master  *SymbolTable
	Parent  *SymbolTable
	Symbols map[string]interface{}
}

func (symbolTable *SymbolTable) Set(s string, object interface{}) {
	symbolTable.Symbols[s] = object
}

func (symbolTable *SymbolTable) GetSelf(symbol string) (interface{}, *errors.Error) {
	result, found := symbolTable.Symbols[symbol]
	if !found {
		return nil, errors.NewNameNotFoundError()
	}
	return result, nil
}

func (symbolTable *SymbolTable) GetAny(symbol string) (interface{}, *errors.Error) {
	result, found := symbolTable.Symbols[symbol]
	if !found {
		var foundError *errors.Error
		if symbolTable.Master != nil {
			result, foundError = symbolTable.Master.GetSelf(symbol)
			if foundError != nil {
				if symbolTable.Parent != nil && symbolTable.Parent != symbolTable.Master {
					result, foundError = symbolTable.Parent.GetAny(symbol)
					if foundError != nil {
						return nil, foundError
					}
				}
				return nil, foundError
			}
		} else if symbolTable.Parent != nil {
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

func NewSymbolTable(masterSymbol *SymbolTable, parentSymbols *SymbolTable) *SymbolTable {
	return &SymbolTable{
		Master:  masterSymbol,
		Parent:  parentSymbols,
		Symbols: map[string]interface{}{},
	}
}
