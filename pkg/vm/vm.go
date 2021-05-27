package vm

import "github.com/shoriwe/gruby/pkg/errors"

type VirtualMachine interface {
	Initialize([]Code) *errors.Error
	Execute() (IObject, *errors.Error)
	LoadCode([]Code)
	PushSymbolTable(*SymbolTable)
	PeekSymbolTable() *SymbolTable
	PopSymbolTable() *SymbolTable
	HashString(string) (int64, *errors.Error)
	HashBytes([]byte) (int64, *errors.Error)
	InitializeByteCode(*Bytecode)
}
