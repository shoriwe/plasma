package vm

import "github.com/shoriwe/gruby/pkg/errors"

type VirtualMachine interface {
	Initialize([]interface{}) *errors.Error
	Execute() (IObject, *errors.Error)
	GetStack() Stack
	MasterSymbolTable() *SymbolTable
	New(*SymbolTable) VirtualMachine
}
