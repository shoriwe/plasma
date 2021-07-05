package vm

import (
	"github.com/shoriwe/gplasma/pkg/errors"
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

	GetClass(*Plasma) *Type
	SetClass(*Type)

	GetBool() bool
	GetBytes() []uint8
	GetString() string
	GetInteger() int64
	GetFloat() float64
	GetContent() []Value
	GetKeyValues() map[int64][]*KeyValue
	GetLength() int

	SetBool(bool)
	SetBytes([]uint8)
	SetString(string)
	SetInteger(int64)
	SetFloat(float64)
	SetContent([]Value)
	SetKeyValues(map[int64][]*KeyValue)
	AddKeyValue(int64, *KeyValue)
	SetLength(int)
	IncreaseLength()
}

func (p *Plasma) QuickGetBool(value Value) (bool, *Object) {
	if _, ok := value.(*Bool); ok {
		return value.GetBool(), nil
	}
	valueToBool, getError := value.Get(ToBool)
	if getError != nil {
		return false, p.NewObjectWithNameNotFoundError(value.GetClass(p), ToBool)
	}
	valueBool, callError := p.CallFunction(valueToBool, valueToBool.SymbolTable().Parent)
	if callError != nil {
		return false, callError
	}
	if _, ok := valueBool.(*Bool); !ok {
		return false, p.NewInvalidTypeError(value.TypeName(), BoolName)
	}
	return valueBool.GetBool(), nil
}
