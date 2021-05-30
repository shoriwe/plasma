package vm

import "github.com/shoriwe/gruby/pkg/errors"

type VirtualMachine interface {
	Initialize([]Code) *errors.Error
	Execute() (IObject, *errors.Error)
	PushObject(IObject)
	PeekObject() IObject
	PopObject() IObject
	PushSymbolTable(*SymbolTable)
	PeekSymbolTable() *SymbolTable
	PopSymbolTable() *SymbolTable
	PushCode(*Bytecode)
	PeekCode() *Bytecode
	PopCode() *Bytecode
	HashString(string) (int64, *errors.Error)
	HashBytes([]byte) (int64, *errors.Error)
	InitializeByteCode(*Bytecode)
	Seed() uint64
}
