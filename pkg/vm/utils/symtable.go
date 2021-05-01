package utils

import (
	"github.com/shoriwe/gruby/pkg/errors"
	"github.com/shoriwe/gruby/pkg/vm/object"
	"github.com/shoriwe/gruby/pkg/vm/vm-errors"
)

type SymbolTable struct {
	ownerName string
	symbols   map[string]object.Object
}

func (symbolTable *SymbolTable) Get(symbol string) (object.Object, *errors.Error) {
	value, ok := symbolTable.symbols[symbol]
	if !ok {
		return nil, vm_errors.NewAttributeNotFound(symbolTable.ownerName, symbol)
	}
	return value, nil
}

func (symbolTable *SymbolTable) Set(symbol string, obj object.Object) *errors.Error {
	if obj == nil {
		return errors.New(vm_errors.UnknownLine, "received a nil pointer as Object", vm_errors.NilObjectError)
	}
	symbolTable.symbols[symbol] = obj
	return nil
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		ownerName: "",
		symbols:   map[string]object.Object{},
	}
}
