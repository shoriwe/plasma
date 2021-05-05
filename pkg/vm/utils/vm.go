package utils

import "github.com/shoriwe/gruby/pkg/errors"

type VirtualMachine interface {
	Initialize([]interface{}) *errors.Error
	Execute() (interface{}, *errors.Error)
	GetStack() Stack
	MasterSymbolTable() *SymbolTable
	New() VirtualMachine
}
