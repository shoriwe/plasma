package vm

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/errors"
)

type IObject interface {
	Id() int64
	TypeName() string
	SymbolTable() *SymbolTable
	SubClasses() []*Type
	Get(string) (IObject, *errors.Error)
	Set(string, IObject)
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
	GetContent() []IObject
	GetKeyValues() map[int64][]*KeyValue
	GetLength() int

	SetBool(bool)
	SetBytes([]uint8)
	SetString(string)
	SetInteger64(int64)
	SetFloat64(float64)
	SetContent([]IObject)
	SetKeyValues(map[int64][]*KeyValue)
	AddKeyValue(int64, *KeyValue)
	SetLength(int)
	IncreaseLength()
}

type Object struct {
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
	Content    []IObject
	KeyValues  map[int64][]*KeyValue
	Length     int
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

func (o *Object) GetContent() []IObject {
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

func (o *Object) SetContent(objects []IObject) {
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

func (o *Object) Get(symbol string) (IObject, *errors.Error) {
	value, getError := o.symbols.GetSelf(symbol)
	if getError != nil {
		return nil, getError
	}
	return value, nil
}

func (o *Object) Set(symbol string, object IObject) {
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

func (p *Plasma) ObjectInitialize(object IObject) *Object {
	object.SymbolTable().Update(map[string]IObject{
		Initialize: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(_ IObject, _ ...IObject) (IObject, *Object) {
					return nil, nil
				},
			),
		),
		NegBits: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(NegBits, 0)),
		Negate: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					selfToBool, foundError := self.Get(ToBool)
					if foundError != nil {
						return nil, p.NewObjectWithNameNotFoundError(ToBool)
					}
					if _, ok := selfToBool.(*Function); !ok {
						return nil, p.NewInvalidTypeError(selfToBool.(IObject).TypeName(), FunctionName)
					}
					selfBool, transformationError := p.CallFunction(selfToBool.(*Function), self.SymbolTable())
					if transformationError != nil {
						return nil, transformationError
					}
					return p.NewBool(p.PeekSymbolTable(), !selfBool.GetBool()), nil
				},
			),
		),
		Negative: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(Negative, 0)),
		Add: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(Add, 1)),
		RightAdd: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(RightAdd, 1)),
		Sub: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(Sub, 1)),
		RightSub: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(RightSub, 1)),
		Mul: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(Mul, 1)),
		RightMul: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(RightMul, 1)),
		Div: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(Div, 1)),
		RightDiv: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(RightDiv, 1)),
		Mod: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(Mod, 1)),
		RightMod: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(RightMod, 1)),
		Pow: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(Pow, 1)),
		RightPow: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(RightPow, 1)),
		BitXor: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(BitXor, 1)),
		RightBitXor: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(RightBitXor, 1)),
		BitAnd: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(BitAnd, 1)),
		RightBitAnd: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(RightBitAnd, 1)),
		BitOr: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(BitOr, 1)),
		RightBitOr: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(RightBitOr, 1)),
		BitLeft: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(BitLeft, 1)),
		RightBitLeft: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(RightBitLeft, 1)),
		BitRight: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(BitRight, 1)),
		RightBitRight: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(RightBitRight, 1)),
		And: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {

					leftToBool, foundError := self.Get(ToBool)
					if foundError != nil {
						return nil, p.NewObjectWithNameNotFoundError(ToBool)
					}
					if _, ok := leftToBool.(*Function); !ok {
						return nil, p.NewInvalidTypeError(leftToBool.(IObject).TypeName(), FunctionName)
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
						return nil, p.NewInvalidTypeError(rightToBool.(IObject).TypeName(), FunctionName)
					}
					var rightBool IObject
					rightBool, transformationError = p.CallFunction(rightToBool.(*Function), right.SymbolTable())
					if transformationError != nil {
						return nil, transformationError
					}
					return p.NewBool(p.PeekSymbolTable(), leftBool.GetBool() && rightBool.GetBool()), nil
				},
			),
		),
		RightAnd: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					rightToBool, foundError := self.Get(ToBool)
					if foundError != nil {
						return nil, p.NewObjectWithNameNotFoundError(ToBool)
					}
					if _, ok := rightToBool.(*Function); !ok {
						return nil, p.NewInvalidTypeError(rightToBool.(IObject).TypeName(), FunctionName)
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
						return nil, p.NewInvalidTypeError(leftToBool.(IObject).TypeName(), FunctionName)
					}
					var leftBool IObject
					leftBool, transformationError = p.CallFunction(leftToBool.(*Function), left.SymbolTable())
					if transformationError != nil {
						return nil, transformationError
					}
					return p.NewBool(p.PeekSymbolTable(), leftBool.GetBool() && rightBool.GetBool()), nil
				},
			),
		),
		Or: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					leftToBool, foundError := self.Get(ToBool)
					if foundError != nil {
						return nil, p.NewObjectWithNameNotFoundError(ToBool)
					}
					if _, ok := leftToBool.(*Function); !ok {
						return nil, p.NewInvalidTypeError(leftToBool.(IObject).TypeName(), FunctionName)
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
						return nil, p.NewInvalidTypeError(rightToBool.(IObject).TypeName(), FunctionName)
					}
					var rightBool IObject
					rightBool, transformationError = p.CallFunction(rightToBool.(*Function), right.SymbolTable())
					if transformationError != nil {
						return nil, transformationError
					}
					return p.NewBool(p.PeekSymbolTable(), leftBool.GetBool() || rightBool.GetBool()), nil
				},
			),
		),
		RightOr: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					rightToBool, foundError := self.Get(ToBool)
					if foundError != nil {
						return nil, p.NewObjectWithNameNotFoundError(ToBool)
					}
					if _, ok := rightToBool.(*Function); !ok {
						return nil, p.NewInvalidTypeError(rightToBool.(IObject).TypeName(), FunctionName)
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
						return nil, p.NewInvalidTypeError(leftToBool.(IObject).TypeName(), FunctionName)
					}
					var leftBool IObject
					leftBool, transformationError = p.CallFunction(leftToBool.(*Function), left.SymbolTable())
					if transformationError != nil {
						return nil, transformationError
					}
					return p.NewBool(p.PeekSymbolTable(), leftBool.GetBool() || rightBool.GetBool()), nil
				},
			),
		),
		Xor: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					leftToBool, foundError := self.Get(ToBool)
					if foundError != nil {
						return nil, p.NewObjectWithNameNotFoundError(ToBool)
					}
					if _, ok := leftToBool.(*Function); !ok {
						return nil, p.NewInvalidTypeError(leftToBool.(IObject).TypeName(), FunctionName)
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
						return nil, p.NewInvalidTypeError(rightToBool.(IObject).TypeName(), FunctionName)
					}
					var rightBool IObject
					rightBool, transformationError = p.CallFunction(rightToBool.(*Function), right.SymbolTable())
					if transformationError != nil {
						return nil, transformationError
					}
					return p.NewBool(p.PeekSymbolTable(), leftBool.GetBool() != rightBool.GetBool()), nil
				},
			),
		),
		RightXor: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					leftToBool, foundError := self.Get(ToBool)
					if foundError != nil {
						return nil, p.NewObjectWithNameNotFoundError(ToBool)
					}
					if _, ok := leftToBool.(*Function); !ok {
						return nil, p.NewInvalidTypeError(leftToBool.(IObject).TypeName(), FunctionName)
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
						return nil, p.NewInvalidTypeError(rightToBool.(IObject).TypeName(), FunctionName)
					}
					var rightBool IObject
					rightBool, transformationError = p.CallFunction(rightToBool.(*Function), left.SymbolTable())
					if transformationError != nil {
						return nil, transformationError
					}
					return p.NewBool(p.PeekSymbolTable(), rightBool.GetBool() != leftBool.GetBool()), nil
				},
			),
		),
		Equals: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1, func(self IObject, arguments ...IObject) (IObject, *Object) {
				right := arguments[0]
				return p.NewBool(p.PeekSymbolTable(), self.Id() == right.Id()), nil
			},
			),
		),
		RightEquals: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					left := arguments[0]
					return p.NewBool(p.PeekSymbolTable(), left.Id() == self.Id()), nil
				},
			),
		),
		NotEquals: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					right := arguments[0]
					return p.NewBool(p.PeekSymbolTable(), self.Id() != right.Id()), nil
				},
			),
		),
		RightNotEquals: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					left := arguments[0]
					return p.NewBool(p.PeekSymbolTable(), left.Id() != self.Id()), nil
				},
			),
		),
		GreaterThan: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(GreaterThan, 1)),
		RightGreaterThan: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(RightGreaterThan, 1)),
		LessThan: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(LessThan, 1)),
		RightLessThan: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(RightLessThan, 1)),
		GreaterThanOrEqual: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(GreaterThanOrEqual, 1)),
		RightGreaterThanOrEqual: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(RightGreaterThanOrEqual, 1)),
		LessThanOrEqual: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(LessThanOrEqual, 1)),
		RightLessThanOrEqual: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(RightLessThanOrEqual, 1)),
		Hash: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					if self.GetHash() == 0 {
						objectHash := p.HashString(fmt.Sprintf("%v-%s-%d", self, self.TypeName(), self.Id()))
						self.SetHash(objectHash)
					}
					return p.NewInteger(p.PeekSymbolTable(), self.GetHash()), nil
				},
			),
		),
		Index: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(Index, 1)),
		Assign: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(Assign, 2)),
		Call: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(Call, 0)),
		Iter: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(Iter, 0)),
		Class: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
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
		SubClasses: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					var subClassesCopy []IObject
					for _, class := range self.SubClasses() {
						subClassesCopy = append(subClassesCopy, class)
					}
					return p.NewTuple(p.PeekSymbolTable(), subClassesCopy), nil
				},
			),
		),
		ToInteger: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(ToInteger, 0)),
		ToFloat: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(ToFloat, 0)),
		ToString: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewString(p.PeekSymbolTable(),
						fmt.Sprintf("%s{%s}-%X", ObjectName, self.TypeName(), self.Id())), nil
				},
			),
		),
		ToBool: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(_ IObject, _ ...IObject) (IObject, *Object) {
					return p.NewBool(p.PeekSymbolTable(), true), nil
				},
			),
		),
		ToArray: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(ToArray, 0)),
		ToTuple: p.NewFunction(object.SymbolTable(),
			p.NewNotImplementedCallable(ToTuple, 0)),
		GetInteger64: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()), nil
				},
			),
		),
		GetBool: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewBool(p.PeekSymbolTable(), self.GetBool()), nil
				},
			),
		),
		GetBytes: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewBytes(p.PeekSymbolTable(), self.GetBytes()), nil
				},
			),
		),
		GetString: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewString(p.PeekSymbolTable(), self.GetString()), nil
				},
			),
		),
		GetFloat64: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewFloat(p.PeekSymbolTable(), self.GetFloat64()), nil
				},
			),
		),
		GetContent: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewArray(p.PeekSymbolTable(), self.GetContent()), nil
				},
			),
		),
		GetKeyValues: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewHashTable(p.PeekSymbolTable(), self.GetKeyValues(), self.GetLength()), nil
				},
			),
		),
		GetLength: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewInteger(p.PeekSymbolTable(), int64(self.GetLength())), nil
				},
			),
		),
		SetBool: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					self.SetBool(arguments[0].GetBool())
					return p.NewNone(), nil
				},
			),
		),
		SetBytes: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					self.SetBytes(arguments[0].GetBytes())
					self.SetLength(arguments[0].GetLength())
					return p.NewNone(), nil
				},
			),
		),
		SetString: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					self.SetString(arguments[0].GetString())
					self.SetLength(arguments[0].GetLength())
					return p.NewNone(), nil
				},
			),
		),
		SetInteger64: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					self.SetInteger64(arguments[0].GetInteger64())
					return p.NewNone(), nil
				},
			),
		),
		SetFloat64: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					self.SetFloat64(arguments[0].GetFloat64())
					return p.NewNone(), nil
				},
			),
		),
		SetContent: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					self.SetContent(arguments[0].GetContent())
					self.SetLength(arguments[0].GetLength())
					return p.NewNone(), nil
				},
			),
		),
		SetKeyValues: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					self.SetKeyValues(arguments[0].GetKeyValues())
					self.SetLength(arguments[0].GetLength())
					return p.NewNone(), nil
				},
			),
		),
		SetLength: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					self.SetLength(arguments[0].GetLength())
					return p.NewNone(), nil
				},
			),
		),
	})
	return nil
}

func (p *Plasma) NewObject(
	typeName string,
	subClasses []*Type,
	parentSymbols *SymbolTable,
) *Object {
	result := &Object{
		id:         p.NextId(),
		typeName:   typeName,
		subClasses: subClasses,
		symbols:    NewSymbolTable(parentSymbols),
	}
	result.Length = 0
	result.Bool = true
	result.String = ""
	result.Integer64 = 0
	result.Float64 = 0
	result.Content = []IObject{}
	result.KeyValues = map[int64][]*KeyValue{}
	result.Bytes = []uint8{}
	p.ObjectInitialize(result)
	return result
}
