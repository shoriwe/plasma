package vm

import (
	"bytes"
	"fmt"
	"github.com/shoriwe/plasma/pkg/lexer"
	"golang.org/x/exp/constraints"
	"sync"
)

const (
	ValueId TypeId = iota
	StringId
	BytesId
	BoolId
	NoneId
	IntId
	FloatId
	ArrayId
	TupleId
	HashId
	BuiltInFunctionId
	FunctionId
	BuiltInClassId
	ClassId
)

type (
	TypeId   int
	Callback func(argument ...*Value) (*Value, error)
	FuncInfo struct {
		Arguments []string
		Bytecode  []byte
	}
	ClassInfo struct {
		prepared bool
		Bases    []*Value
		Bytecode []byte
	}
	Value struct {
		onDemand map[string]func(self *Value) *Value
		class    *Value
		typeId   TypeId
		mutex    *sync.Mutex
		v        any
		vtable   *Symbols
	}
)

func (plasma *Plasma) valueClass() *Value {
	class := plasma.NewValue(plasma.rootSymbols, BuiltInClassId, plasma.class)
	class.SetAny(Callback(func(argument ...*Value) (*Value, error) {
		return plasma.NewValue(plasma.rootSymbols, ValueId, plasma.value), nil
	}))
	return class
}

func (value *Value) GetClass() *Value {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	return value.class
}

func (value *Value) TypeId() TypeId {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	return value.typeId
}

func (value *Value) VirtualTable() *Symbols {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	return value.vtable
}

func (value *Value) SetAny(v any) {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	value.v = v
}

/*
GetHash cast the internal value to *Hash
*/
func (value *Value) GetHash() *Hash {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	return value.v.(*Hash)
}

/*
GetCallback cast the internal value to Callback
*/
func (value *Value) GetCallback() Callback {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	return value.v.(Callback)
}

/*
GetValues cast the internal value to []*Value
*/
func (value *Value) GetValues() []*Value {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	return value.v.([]*Value)
}

/*
GetFuncInfo cast the internal value to FuncInfo
*/
func (value *Value) GetFuncInfo() FuncInfo {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	return value.v.(FuncInfo)
}

/*
GetClassInfo cast the internal value to *ClassInfo
*/
func (value *Value) GetClassInfo() *ClassInfo {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	return value.v.(*ClassInfo)
}

/*
GetBytes cast the internal value to []byte
*/
func (value *Value) GetBytes() []byte {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	return value.v.([]byte)
}

/*
GetBool cast the internal value to bool
*/
func (value *Value) GetBool() bool {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	return value.v.(bool)
}

/*
GetInt64 cast the internal value to int64
*/
func (value *Value) GetInt64() int64 {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	return value.v.(int64)
}

/*
GetFloat64 cast the internal value to float64
*/
func (value *Value) GetFloat64() float64 {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	return value.v.(float64)
}

/*
GetAny returns whatever the internal value contains
*/
func (value *Value) GetAny() any {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	return value.v
}

/*
Set a plasma *Value to a symbol inside the Object
*/
func (value *Value) Set(symbol string, v *Value) {
	value.vtable.Set(symbol, v)
}

/*
Get Retrieves the value named as the symbol
*/
func (value *Value) Get(symbol string) (*Value, error) {
	result, getError := value.vtable.Get(symbol)
	if getError == nil {
		return result, nil
	}
	value.mutex.Lock()
	defer value.mutex.Unlock()
	onDemand, found := value.onDemand[symbol]
	if !found {
		return nil, fmt.Errorf(SymbolNotFoundError, symbol)
	}
	result = onDemand(value)
	value.vtable.Set(symbol, result)
	return result, nil
}

/*
Del deletes the reference of the value named by the symbol
*/
func (value *Value) Del(symbol string) error {
	return value.vtable.Del(symbol)
}

/*
Bool interpret the Value as a go bool
*/
func (value *Value) Bool() bool {
	switch value.TypeId() {
	case ValueId:
		return true
	case StringId, BytesId:
		return len(value.GetBytes()) > 0
	case BoolId:
		return value.GetBool()
	case NoneId:
		return false
	case IntId:
		return value.GetInt64() != 0
	case FloatId:
		return value.GetFloat64() != 0
	case ArrayId, TupleId:
		return len(value.GetValues()) > 0
	case HashId:
		return value.GetHash().Size() > 0
	case BuiltInFunctionId:
		return true
	case FunctionId:
		return true
	case BuiltInClassId:
		return true
	case ClassId:
		return true
	}
	return false
}

/*
Bytes interpret the Value as a go []byte
*/
func (value *Value) Bytes() []byte {
	switch value.TypeId() {
	case ValueId:
		return []byte("?Value")
	case StringId, BytesId:
		return value.GetBytes()
	case BoolId:
		if value.GetBool() {
			return []byte(lexer.TrueString)
		}
		return []byte(lexer.FalseString)
	case NoneId:
		return []byte(lexer.NoneString)
	case IntId:
		return []byte(fmt.Sprintf("%d", value.GetInt64()))
	case FloatId:
		return []byte(fmt.Sprintf("%f", value.GetFloat64()))
	case ArrayId:
		value.mutex.Lock()
		defer value.mutex.Unlock()
		builder := &bytes.Buffer{}
		builder.WriteByte('[')
		for index, v := range value.v.([]*Value) {
			if index != 0 {
				builder.Write([]byte{',', ' '})
			}
			if v.TypeId() == StringId {
				builder.WriteByte('"')
				builder.Write(v.Bytes())
				builder.WriteByte('"')
			} else {
				builder.Write(v.Bytes())
			}
		}
		builder.WriteByte(']')
		return builder.Bytes()
	case TupleId:
		value.mutex.Lock()
		defer value.mutex.Unlock()
		builder := &bytes.Buffer{}
		builder.WriteByte('(')
		for index, v := range value.v.([]*Value) {
			if index != 0 {
				builder.Write([]byte{',', ' '})
			}
			if v.TypeId() == StringId {
				builder.WriteByte('"')
				builder.Write(v.Bytes())
				builder.WriteByte('"')
			} else {
				builder.Write(v.Bytes())
			}
		}
		builder.WriteByte(')')
		return builder.Bytes()
	case HashId:
		hash := value.GetHash()
		hash.mutex.Lock()
		defer hash.mutex.Unlock()
		builder := &bytes.Buffer{}
		builder.WriteByte('{')
		index := 0
		for _, keyValue := range hash.internalMap {
			if index != 0 {
				builder.Write([]byte{',', ' '})
			}
			if keyValue.Key.TypeId() == StringId {
				builder.WriteByte('"')
				builder.Write(keyValue.Key.Bytes())
				builder.WriteByte('"')
			} else {
				builder.Write(keyValue.Key.Bytes())
			}
			builder.Write([]byte{':', ' '})
			if keyValue.Value.TypeId() == StringId {
				builder.WriteByte('"')
				builder.Write(keyValue.Value.Bytes())
				builder.WriteByte('"')
			} else {
				builder.Write(keyValue.Value.Bytes())
			}
			index++
		}
		builder.WriteByte('}')
		return builder.Bytes()
	case BuiltInFunctionId:
		return []byte("?BuiltInFunction")
	case FunctionId:
		return []byte("?Function")
	case BuiltInClassId:
		return []byte("?BuiltInClass")
	case ClassId:
		return []byte("?Class")
	}
	return []byte{}
}

/*
String interpret the Value as a go string
*/
func (value *Value) String() string {
	return string(value.Bytes())
}

/*
Int interpret the Value as a go int64
*/
func Int[T constraints.Integer](value *Value) T {
	switch value.TypeId() {
	case ValueId:
		return 0
	case StringId:
		return 0
	case BytesId:
		return 0
	case BoolId:
		if value.GetBool() {
			return 1
		}
		return 0
	case NoneId:
		return 0
	case IntId:
		return T(value.GetInt64())
	case FloatId:
		return T(value.GetFloat64())
	case ArrayId:
		return 0
	case TupleId:
		return 0
	case HashId:
		return 0
	case BuiltInFunctionId:
		return 0
	case FunctionId:
		return 0
	case BuiltInClassId:
		return 0
	case ClassId:
		return 0
	}
	return 0
}

func Float[T constraints.Float](value *Value) T {
	switch value.TypeId() {
	case ValueId:
		return 0
	case StringId:
		return 0
	case BytesId:
		return 0
	case BoolId:
		if value.GetBool() {
			return 1
		}
		return 0
	case NoneId:
		return 0
	case IntId:
		return T(value.GetInt64())
	case FloatId:
		return T(value.GetFloat64())
	case ArrayId:
		return 0
	case TupleId:
		return 0
	case HashId:
		return 0
	case BuiltInFunctionId:
		return 0
	case FunctionId:
		return 0
	case BuiltInClassId:
		return 0
	case ClassId:
		return 0
	}
	return 0
}

func (value *Value) Values() []*Value {
	switch value.TypeId() {
	case ValueId:
		return nil
	case StringId:
		return nil
	case BytesId:
		return nil
	case BoolId:
		return nil
	case NoneId:
		return nil
	case IntId:
		return nil
	case FloatId:
		return nil
	case ArrayId, TupleId:
		return value.GetValues()
	case HashId:
		return nil
	case BuiltInFunctionId:
		return nil
	case FunctionId:
		return nil
	case BuiltInClassId:
		return nil
	case ClassId:
		return nil
	}
	return nil
}

func (value *Value) Call(argument ...*Value) (*Value, error) {
	return value.GetCallback()(argument...)
}

func (value *Value) Implements(class *Value) bool {
	if value == class {
		return true
	}
	for _, base := range value.GetClassInfo().Bases {
		if base.Implements(class) {
			return true
		}
	}
	return false
}

func (value *Value) Equal(other *Value) bool {
	switch value.TypeId() {
	case ValueId:
		return value.ValueEqual(other)
	case StringId:
		return value.StringEqual(other)
	case BytesId:
		return value.BytesEqual(other)
	case BoolId:
		return value.BoolEqual(other)
	case NoneId:
		return value.NoneEqual(other)
	case IntId:
		return value.IntEqual(other)
	case FloatId:
		return value.FloatEqual(other)
	case ArrayId:
		return value.ArrayEqual(other)
	case TupleId:
		return value.TupleEqual(other)
	case HashId:
		return value.HashEqual(other)
	case BuiltInFunctionId:
		return value.BuiltInFunctionEqual(other)
	case FunctionId:
		return value.FunctionEqual(other)
	case BuiltInClassId:
		return value.BuiltInClassEqual(other)
	case ClassId:
		return value.ClassEqual(other)
	}
	return false
}

func (value *Value) ValueEqual(other *Value) bool {
	return value == other
}

func (value *Value) StringEqual(other *Value) bool {
	return value.String() == other.String()
}

func (value *Value) BytesEqual(other *Value) bool {
	return bytes.Equal(value.GetBytes(), other.GetBytes())
}

func (value *Value) BoolEqual(other *Value) bool {
	return value.Bool() == other.Bool()
}

func (value *Value) NoneEqual(other *Value) bool {
	return value == other
}

func (value *Value) IntEqual(other *Value) bool {
	switch other.TypeId() {
	case IntId:
		return Int[int64](value) == Int[int64](other)
	case FloatId:
		return Float[float64](value) == Float[float64](other)
	}
	return false
}

func (value *Value) FloatEqual(other *Value) bool {
	switch other.TypeId() {
	case IntId:
		return Float[float64](value) == Float[float64](other)
	case FloatId:
		return Float[float64](value) == Float[float64](other)
	}
	return false
}

func (value *Value) ArrayEqual(other *Value) bool {
	switch other.TypeId() {
	case ArrayId:
		values := value.GetValues()
		otherValues := other.GetValues()
		if len(values) != len(otherValues) {
			return false
		}
		for index, internalValue := range values {
			if !internalValue.Equal(otherValues[index]) {
				return false
			}
		}
		return true
	}
	return false
}

func (value *Value) TupleEqual(other *Value) bool {
	switch other.TypeId() {
	case TupleId:
		values := value.GetValues()
		otherValues := other.GetValues()
		if len(values) != len(otherValues) {
			return false
		}
		for index, internalValue := range values {
			if !internalValue.Equal(otherValues[index]) {
				return false
			}
		}
		return true
	}
	return false
}

func (value *Value) HashEqual(other *Value) bool {
	// TODO: implement me!
	return value == other
}

func (value *Value) BuiltInFunctionEqual(other *Value) bool {
	return value == other
}

func (value *Value) FunctionEqual(other *Value) bool {
	return value == other
}

func (value *Value) BuiltInClassEqual(other *Value) bool {
	return value == other
}

func (value *Value) ClassEqual(other *Value) bool {
	return value == other
}

/*
NewValue Creates a new Value
*/
func (plasma *Plasma) NewValue(parent *Symbols, typeId TypeId, class *Value) *Value {
	return &Value{
		onDemand: plasma.onDemand,
		class:    class,
		typeId:   typeId,
		mutex:    &sync.Mutex{},
		v:        nil,
		vtable:   NewSymbols(parent),
	}
}
