package vm

import (
	"github.com/shoriwe/gplasma/pkg/errors"
	"math/big"
)

type Value interface {
	IsBuiltIn() bool
	Id() int64
	TypeName() string
	SymbolTable() *SymbolTable
	SubClasses() []*Type
	Get(string) (Value, *errors.Error)
	Set(string, Value)
	GetHash() int64
	SetHash(int64)

	Implements(*Type) bool // This should check if the object implements a class directly or indirectly

	GetClass() *Type
	SetClass(*Type)

	GetBool() bool
	GetBytes() []uint8
	GetString() string
	GetInteger() *big.Int
	GetFloat() *big.Float
	GetContent() []Value
	GetKeyValues() map[int64][]*KeyValue
	GetLength() int

	SetBool(bool)
	SetBytes([]uint8)
	SetString(string)
	SetInteger(*big.Int)
	SetFloat(*big.Float)
	SetContent([]Value)
	SetKeyValues(map[int64][]*KeyValue)
	AddKeyValue(int64, *KeyValue)
	SetLength(int)
	IncreaseLength()
}
