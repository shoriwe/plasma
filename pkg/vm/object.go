package vm

import (
	"fmt"
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

	GetClass() *Type
	SetClass(*Type)

	GetBool() bool
	GetBytes() []uint8
	GetString() string
	GetInteger64() int64
	GetFloat64() float64
	GetContent() []Value
	GetKeyValues() map[int64][]*KeyValue
	GetLength() int

	SetBool(bool)
	SetBytes([]uint8)
	SetString(string)
	SetInteger64(int64)
	SetFloat64(float64)
	SetContent([]Value)
	SetKeyValues(map[int64][]*KeyValue)
	AddKeyValue(int64, *KeyValue)
	SetLength(int)
	IncreaseLength()
}

type Object struct {
	isBuiltIn  bool
	id         int64
	typeName   string
	class      *Type
	subClasses []*Type
	symbols    *SymbolTable
	hash       int64
	Bool       bool
	String     string
	Bytes      []uint8
	Integer64  int64
	Float64    float64
	Content    []Value
	KeyValues  map[int64][]*KeyValue
	Length     int
}

func (o *Object) IsBuiltIn() bool {
	return o.isBuiltIn
}

func (o *Object) IncreaseLength() {
	o.Length++
}

func (o *Object) GetBool() bool {
	return o.Bool
}

func (o *Object) GetBytes() []uint8 {
	return o.Bytes
}

func (o *Object) GetString() string {
	return o.String
}

func (o *Object) GetInteger64() int64 {
	return o.Integer64
}

func (o *Object) GetFloat64() float64 {
	return o.Float64
}

func (o *Object) GetContent() []Value {
	return o.Content
}

func (o *Object) GetKeyValues() map[int64][]*KeyValue {
	return o.KeyValues
}

func (o *Object) GetLength() int {
	return o.Length
}

func (o *Object) SetString(s string) {
	o.String = s
}

func (o *Object) SetInteger64(i int64) {
	o.Integer64 = i
}

func (o *Object) SetFloat64(f float64) {
	o.Float64 = f
}

func (o *Object) SetContent(objects []Value) {
	o.Content = objects
}

func (o *Object) SetKeyValues(m map[int64][]*KeyValue) {
	o.KeyValues = m
}

func (o *Object) AddKeyValue(hash int64, keyValue *KeyValue) {
	o.KeyValues[hash] = append(o.KeyValues[hash], keyValue)
}

func (o *Object) SetLength(i int) {
	o.Length = i
}

func (o *Object) SetBool(b bool) {
	o.Bool = b
}

func (o *Object) SetBytes(b []uint8) {
	o.Bytes = b
}

func (o *Object) Id() int64 {
	return o.id
}

func (o *Object) SubClasses() []*Type {
	return o.subClasses
}

func (o *Object) Get(symbol string) (Value, *errors.Error) {
	value, getError := o.symbols.GetSelf(symbol)
	if getError != nil {
		return nil, getError
	}
	return value, nil
}

func (o *Object) Set(symbol string, object Value) {
	o.symbols.Set(symbol, object)
}

func (o *Object) TypeName() string {
	return o.typeName
}

func (o *Object) SymbolTable() *SymbolTable {
	return o.symbols
}

func (o *Object) GetHash() int64 {
	return o.hash
}

func (o *Object) SetHash(newHash int64) {
	o.hash = newHash
}

func (o *Object) GetClass() *Type {
	return o.class
}

func (o *Object) SetClass(class *Type) {
	o.class = class
}

func (o *Object) Implements(class *Type) bool {
	if o.class == class {
		return true
	}
	for _, subClass := range o.subClasses {
		if subClass.Implements(class) {
			return true
		}
	}
	return false
}

func (p *Plasma) ObjectInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(object Value) *Object {
		object.SymbolTable().Update(map[string]Value{
			Initialize: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(_ Value, _ ...Value) (Value, *Object) {
						return p.NewNone(), nil
					},
				),
			),
			Negate: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						selfToBool, foundError := self.Get(ToBool)
						if foundError != nil {
							return nil, p.NewObjectWithNameNotFoundError(ToBool)
						}
						if _, ok := selfToBool.(*Function); !ok {
							return nil, p.NewInvalidTypeError(selfToBool.(Value).TypeName(), FunctionName)
						}
						selfBool, transformationError := p.CallFunction(selfToBool.(*Function), self.SymbolTable())
						if transformationError != nil {
							return nil, transformationError
						}
						return p.NewBool(false, p.PeekSymbolTable(), !selfBool.GetBool()), nil
					},
				),
			),
			And: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {

						leftToBool, foundError := self.Get(ToBool)
						if foundError != nil {
							return nil, p.NewObjectWithNameNotFoundError(ToBool)
						}
						if _, ok := leftToBool.(*Function); !ok {
							return nil, p.NewInvalidTypeError(leftToBool.(Value).TypeName(), FunctionName)
						}
						leftBool, transformationError := p.CallFunction(leftToBool.(*Function), self.SymbolTable())
						if transformationError != nil {
							return nil, transformationError
						}
						right := arguments[0]
						var rightToBool interface{}
						rightToBool, foundError = right.Get(ToBool)
						if foundError != nil {
							return nil, p.NewObjectWithNameNotFoundError(ToBool)
						}
						if _, ok := rightToBool.(*Function); !ok {
							return nil, p.NewInvalidTypeError(rightToBool.(Value).TypeName(), FunctionName)
						}
						var rightBool Value
						rightBool, transformationError = p.CallFunction(rightToBool.(*Function), right.SymbolTable())
						if transformationError != nil {
							return nil, transformationError
						}
						return p.NewBool(false, p.PeekSymbolTable(), leftBool.GetBool() && rightBool.GetBool()), nil
					},
				),
			),
			RightAnd: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						rightToBool, foundError := self.Get(ToBool)
						if foundError != nil {
							return nil, p.NewObjectWithNameNotFoundError(ToBool)
						}
						if _, ok := rightToBool.(*Function); !ok {
							return nil, p.NewInvalidTypeError(rightToBool.(Value).TypeName(), FunctionName)
						}
						rightBool, transformationError := p.CallFunction(rightToBool.(*Function), self.SymbolTable())
						if transformationError != nil {
							return nil, transformationError
						}
						left := arguments[0]
						var leftToBool interface{}
						leftToBool, foundError = left.Get(ToBool)
						if foundError != nil {
							return nil, p.NewObjectWithNameNotFoundError(ToBool)
						}
						if _, ok := leftToBool.(*Function); !ok {
							return nil, p.NewInvalidTypeError(leftToBool.(Value).TypeName(), FunctionName)
						}
						var leftBool Value
						leftBool, transformationError = p.CallFunction(leftToBool.(*Function), left.SymbolTable())
						if transformationError != nil {
							return nil, transformationError
						}
						return p.NewBool(false, p.PeekSymbolTable(), leftBool.GetBool() && rightBool.GetBool()), nil
					},
				),
			),
			Or: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						leftToBool, foundError := self.Get(ToBool)
						if foundError != nil {
							return nil, p.NewObjectWithNameNotFoundError(ToBool)
						}
						if _, ok := leftToBool.(*Function); !ok {
							return nil, p.NewInvalidTypeError(leftToBool.(Value).TypeName(), FunctionName)
						}
						leftBool, transformationError := p.CallFunction(leftToBool.(*Function), self.SymbolTable())
						if transformationError != nil {
							return nil, transformationError
						}

						right := arguments[0]
						var rightToBool interface{}
						rightToBool, foundError = right.Get(ToBool)
						if foundError != nil {
							return nil, p.NewObjectWithNameNotFoundError(ToBool)
						}
						if _, ok := rightToBool.(*Function); !ok {
							return nil, p.NewInvalidTypeError(rightToBool.(Value).TypeName(), FunctionName)
						}
						var rightBool Value
						rightBool, transformationError = p.CallFunction(rightToBool.(*Function), right.SymbolTable())
						if transformationError != nil {
							return nil, transformationError
						}
						return p.NewBool(false, p.PeekSymbolTable(), leftBool.GetBool() || rightBool.GetBool()), nil
					},
				),
			),
			RightOr: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						rightToBool, foundError := self.Get(ToBool)
						if foundError != nil {
							return nil, p.NewObjectWithNameNotFoundError(ToBool)
						}
						if _, ok := rightToBool.(*Function); !ok {
							return nil, p.NewInvalidTypeError(rightToBool.(Value).TypeName(), FunctionName)
						}
						rightBool, transformationError := p.CallFunction(rightToBool.(*Function), self.SymbolTable())
						if transformationError != nil {
							return nil, transformationError
						}
						left := arguments[0]
						var leftToBool interface{}
						leftToBool, foundError = left.Get(ToBool)
						if foundError != nil {
							return nil, p.NewObjectWithNameNotFoundError(ToBool)
						}
						if _, ok := leftToBool.(*Function); !ok {
							return nil, p.NewInvalidTypeError(leftToBool.(Value).TypeName(), FunctionName)
						}
						var leftBool Value
						leftBool, transformationError = p.CallFunction(leftToBool.(*Function), left.SymbolTable())
						if transformationError != nil {
							return nil, transformationError
						}
						return p.NewBool(false, p.PeekSymbolTable(), leftBool.GetBool() || rightBool.GetBool()), nil
					},
				),
			),
			Xor: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						leftToBool, foundError := self.Get(ToBool)
						if foundError != nil {
							return nil, p.NewObjectWithNameNotFoundError(ToBool)
						}
						if _, ok := leftToBool.(*Function); !ok {
							return nil, p.NewInvalidTypeError(leftToBool.(Value).TypeName(), FunctionName)
						}
						leftBool, transformationError := p.CallFunction(leftToBool.(*Function), self.SymbolTable())
						if transformationError != nil {
							return nil, transformationError
						}

						right := arguments[0]
						var rightToBool interface{}
						rightToBool, foundError = right.Get(ToBool)
						if foundError != nil {
							return nil, p.NewObjectWithNameNotFoundError(ToBool)
						}
						if _, ok := rightToBool.(*Function); !ok {
							return nil, p.NewInvalidTypeError(rightToBool.(Value).TypeName(), FunctionName)
						}
						var rightBool Value
						rightBool, transformationError = p.CallFunction(rightToBool.(*Function), right.SymbolTable())
						if transformationError != nil {
							return nil, transformationError
						}
						return p.NewBool(false, p.PeekSymbolTable(), leftBool.GetBool() != rightBool.GetBool()), nil
					},
				),
			),
			RightXor: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						leftToBool, foundError := self.Get(ToBool)
						if foundError != nil {
							return nil, p.NewObjectWithNameNotFoundError(ToBool)
						}
						if _, ok := leftToBool.(*Function); !ok {
							return nil, p.NewInvalidTypeError(leftToBool.(Value).TypeName(), FunctionName)
						}
						leftBool, transformationError := p.CallFunction(leftToBool.(*Function), self.SymbolTable())
						if transformationError != nil {
							return nil, transformationError
						}

						left := arguments[0]
						var rightToBool interface{}
						rightToBool, foundError = left.Get(ToBool)
						if foundError != nil {
							return nil, p.NewObjectWithNameNotFoundError(ToBool)
						}
						if _, ok := rightToBool.(*Function); !ok {
							return nil, p.NewInvalidTypeError(rightToBool.(Value).TypeName(), FunctionName)
						}
						var rightBool Value
						rightBool, transformationError = p.CallFunction(rightToBool.(*Function), left.SymbolTable())
						if transformationError != nil {
							return nil, transformationError
						}
						return p.NewBool(false, p.PeekSymbolTable(), rightBool.GetBool() != leftBool.GetBool()), nil
					},
				),
			),
			Equals: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1, func(self Value, arguments ...Value) (Value, *Object) {
					right := arguments[0]
					return p.NewBool(false, p.PeekSymbolTable(), self.Id() == right.Id()), nil
				},
				),
			),
			RightEquals: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						return p.NewBool(false, p.PeekSymbolTable(), left.Id() == self.Id()), nil
					},
				),
			),
			NotEquals: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						return p.NewBool(false, p.PeekSymbolTable(), self.Id() != right.Id()), nil
					},
				),
			),
			RightNotEquals: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						return p.NewBool(false, p.PeekSymbolTable(), left.Id() != self.Id()), nil
					},
				),
			),
			Hash: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						if self.GetHash() == 0 {
							objectHash := p.HashString(fmt.Sprintf("%v-%s-%d", self, self.TypeName(), self.Id()))
							self.SetHash(objectHash)
						}
						return p.NewInteger(false, p.PeekSymbolTable(), self.GetHash()), nil
					},
				),
			),
			Class: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						if self.GetClass() == nil { // This should only happen with built-ins
							class, getError := p.BuiltInSymbols().GetAny(self.TypeName())
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(self.TypeName())
							}
							if _, ok := class.(*Type); !ok {
								return nil, p.NewInvalidTypeError(class.TypeName(), TypeName)
							}
							self.SetClass(class.(*Type))
						}
						return self.GetClass(), nil
					},
				),
			),
			SubClasses: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						var subClassesCopy []Value
						for _, class := range self.SubClasses() {
							subClassesCopy = append(subClassesCopy, class)
						}
						return p.NewTuple(false, p.PeekSymbolTable(), subClassesCopy), nil
					},
				),
			),
			ToString: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewString(false, p.PeekSymbolTable(),
							fmt.Sprintf("%s{%s}-%X", ObjectName, self.TypeName(), self.Id())), nil
					},
				),
			),
			ToBool: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(_ Value, _ ...Value) (Value, *Object) {
						return p.NewBool(false, p.PeekSymbolTable(), true), nil
					},
				),
			),
			GetInteger64: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger64()), nil
					},
				),
			),
			GetBool: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewBool(false, p.PeekSymbolTable(), self.GetBool()), nil
					},
				),
			),
			GetBytes: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewBytes(false, p.PeekSymbolTable(), self.GetBytes()), nil
					},
				),
			),
			GetString: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewString(false, p.PeekSymbolTable(), self.GetString()), nil
					},
				),
			),
			GetFloat64: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewFloat(false, p.PeekSymbolTable(), self.GetFloat64()), nil
					},
				),
			),
			GetContent: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewArray(false, p.PeekSymbolTable(), self.GetContent()), nil
					},
				),
			),
			GetKeyValues: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewHashTable(false, p.PeekSymbolTable(), self.GetKeyValues(), self.GetLength()), nil
					},
				),
			),
			GetLength: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewInteger(false, p.PeekSymbolTable(), int64(self.GetLength())), nil
					},
				),
			),
			SetBool: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						self.SetBool(arguments[0].GetBool())
						return p.NewNone(), nil
					},
				),
			),
			SetBytes: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						self.SetBytes(arguments[0].GetBytes())
						self.SetLength(arguments[0].GetLength())
						return p.NewNone(), nil
					},
				),
			),
			SetString: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						self.SetString(arguments[0].GetString())
						self.SetLength(arguments[0].GetLength())
						return p.NewNone(), nil
					},
				),
			),
			SetInteger64: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						self.SetInteger64(arguments[0].GetInteger64())
						return p.NewNone(), nil
					},
				),
			),
			SetFloat64: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						self.SetFloat64(arguments[0].GetFloat64())
						return p.NewNone(), nil
					},
				),
			),
			SetContent: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						self.SetContent(arguments[0].GetContent())
						self.SetLength(arguments[0].GetLength())
						return p.NewNone(), nil
					},
				),
			),
			SetKeyValues: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						self.SetKeyValues(arguments[0].GetKeyValues())
						self.SetLength(arguments[0].GetLength())
						return p.NewNone(), nil
					},
				),
			),
			SetLength: p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						self.SetLength(arguments[0].GetLength())
						return p.NewNone(), nil
					},
				),
			),
		})
		return nil
	}
}
func (p *Plasma) NewObject(
	isBuiltIn bool,
	typeName string,
	subClasses []*Type,
	parentSymbols *SymbolTable,
) *Object {
	result := &Object{
		id:         p.NextId(),
		typeName:   typeName,
		subClasses: subClasses,
		symbols:    NewSymbolTable(parentSymbols),
		isBuiltIn:  isBuiltIn,
	}
	result.Length = 0
	result.Bool = true
	result.String = ""
	result.Integer64 = 0
	result.Float64 = 0
	result.Content = []Value{}
	result.KeyValues = map[int64][]*KeyValue{}
	result.Bytes = []uint8{}
	result.Set(Self, result)
	p.ObjectInitialize(isBuiltIn)(result)
	return result
}
