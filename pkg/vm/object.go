package vm

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/errors"
)

type IObject interface {
	Id() uint64
	TypeName() string
	SymbolTable() *SymbolTable
	SubClasses() []*Type
	Get(string) (IObject, *errors.Error)
	Set(string, IObject)
	GetHash() int64
	SetHash(int64)

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
	id             uint64
	typeName       string
	class          *Type
	subClasses     []*Type
	symbols        *SymbolTable
	virtualMachine VirtualMachine
	hash           int64
	Bool           bool
	String         string
	Bytes          []uint8
	Integer64      int64
	Float64        float64
	Content        []IObject
	KeyValues      map[int64][]*KeyValue
	Length         int
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

func (o *Object) Id() uint64 {
	return o.id
}

func (o *Object) SubClasses() []*Type {
	return o.subClasses
}

func (o *Object) Get(symbol string) (IObject, *errors.Error) {
	return o.symbols.GetSelf(symbol)
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

func (p *Plasma) ConstructObject(type_ *Type, vm VirtualMachine, parent *SymbolTable) (IObject, *errors.Error) {
	object := p.NewObject(type_.typeName, type_.subClasses, parent)
	for _, subclass := range object.subClasses {
		initializationError := subclass.Constructor.Initialize(vm, object)
		if initializationError != nil {
			return nil, initializationError
		}
	}
	return object, nil
}

func (p *Plasma) ObjectInitialize(object IObject) *errors.Error {
	object.SymbolTable().Update(map[string]IObject{
		Initialize: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(_ IObject, _ ...IObject) (IObject, *errors.Error) {
					return nil, nil
				},
			),
		),
		NegBits: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(0)),
		Negate: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					selfToBool, foundError := self.Get(ToBool)
					if foundError != nil {
						return nil, foundError
					}
					if _, ok := selfToBool.(*Function); !ok {
						return nil, errors.NewTypeError(selfToBool.(IObject).TypeName(), FunctionName)
					}
					var selfBool IObject
					var transformationError *errors.Error
					selfBool, transformationError = p.CallFunction(selfToBool.(*Function), self.SymbolTable())
					if transformationError != nil {
						return nil, transformationError
					}
					return p.NewBool(p.PeekSymbolTable(), !selfBool.GetBool()), nil
				},
			),
		),
		Negative: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(0)),
		Add: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		RightAdd: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		Sub: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		RightSub: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		Mul: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		RightMul: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		Div: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		RightDiv: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		Mod: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		RightMod: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		Pow: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		RightPow: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		BitXor: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		RightBitXor: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		BitAnd: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		RightBitAnd: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		BitOr: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		RightBitOr: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		BitLeft: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		RightBitLeft: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		BitRight: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		RightBitRight: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		And: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {

					leftToBool, foundError := self.Get(ToBool)
					if foundError != nil {
						return nil, foundError
					}
					if _, ok := leftToBool.(*Function); !ok {
						return nil, errors.NewTypeError(leftToBool.(IObject).TypeName(), FunctionName)
					}
					leftBool, transformationError := p.CallFunction(leftToBool.(*Function), self.SymbolTable())
					if transformationError != nil {
						return nil, transformationError
					}
					right := arguments[0]
					var rightToBool interface{}
					rightToBool, foundError = right.Get(ToBool)
					if foundError != nil {
						return nil, foundError
					}
					if _, ok := rightToBool.(*Function); !ok {
						return nil, errors.NewTypeError(rightToBool.(IObject).TypeName(), FunctionName)
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
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					rightToBool, foundError := self.Get(ToBool)
					if foundError != nil {
						return nil, foundError
					}
					if _, ok := rightToBool.(*Function); !ok {
						return nil, errors.NewTypeError(rightToBool.(IObject).TypeName(), FunctionName)
					}
					rightBool, transformationError := p.CallFunction(rightToBool.(*Function), self.SymbolTable())
					if transformationError != nil {
						return nil, transformationError
					}
					left := arguments[0]
					var leftToBool interface{}
					leftToBool, foundError = left.Get(ToBool)
					if foundError != nil {
						return nil, foundError
					}
					if _, ok := leftToBool.(*Function); !ok {
						return nil, errors.NewTypeError(leftToBool.(IObject).TypeName(), FunctionName)
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
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					leftToBool, foundError := self.Get(ToBool)
					if foundError != nil {
						return nil, foundError
					}
					if _, ok := leftToBool.(*Function); !ok {
						return nil, errors.NewTypeError(leftToBool.(IObject).TypeName(), FunctionName)
					}
					leftBool, transformationError := p.CallFunction(leftToBool.(*Function), self.SymbolTable())
					if transformationError != nil {
						return nil, transformationError
					}

					right := arguments[0]
					var rightToBool interface{}
					rightToBool, foundError = right.Get(ToBool)
					if foundError != nil {
						return nil, foundError
					}
					if _, ok := rightToBool.(*Function); !ok {
						return nil, errors.NewTypeError(rightToBool.(IObject).TypeName(), FunctionName)
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
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					rightToBool, foundError := self.Get(ToBool)
					if foundError != nil {
						return nil, foundError
					}
					if _, ok := rightToBool.(*Function); !ok {
						return nil, errors.NewTypeError(rightToBool.(IObject).TypeName(), FunctionName)
					}
					rightBool, transformationError := p.CallFunction(rightToBool.(*Function), self.SymbolTable(), self)
					if transformationError != nil {
						return nil, transformationError
					}
					left := arguments[0]
					var leftToBool interface{}
					leftToBool, foundError = left.Get(ToBool)
					if foundError != nil {
						return nil, foundError
					}
					if _, ok := leftToBool.(*Function); !ok {
						return nil, errors.NewTypeError(leftToBool.(IObject).TypeName(), FunctionName)
					}
					var leftBool IObject
					leftBool, transformationError = p.CallFunction(leftToBool.(*Function), left.SymbolTable(), left)
					if transformationError != nil {
						return nil, transformationError
					}
					return p.NewBool(p.PeekSymbolTable(), leftBool.GetBool() || rightBool.GetBool()), nil
				},
			),
		),
		Xor: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					leftToBool, foundError := self.Get(ToBool)
					if foundError != nil {
						return nil, foundError
					}
					if _, ok := leftToBool.(*Function); !ok {
						return nil, errors.NewTypeError(leftToBool.(IObject).TypeName(), FunctionName)
					}
					leftBool, transformationError := p.CallFunction(leftToBool.(*Function), self.SymbolTable(), self)
					if transformationError != nil {
						return nil, transformationError
					}

					right := arguments[0]
					var rightToBool interface{}
					rightToBool, foundError = right.Get(ToBool)
					if foundError != nil {
						return nil, foundError
					}
					if _, ok := rightToBool.(*Function); !ok {
						return nil, errors.NewTypeError(rightToBool.(IObject).TypeName(), FunctionName)
					}
					var rightBool IObject
					rightBool, transformationError = p.CallFunction(rightToBool.(*Function), right.SymbolTable(), right)
					if transformationError != nil {
						return nil, transformationError
					}
					return p.NewBool(p.PeekSymbolTable(), leftBool.GetBool() != rightBool.GetBool()), nil
				},
			),
		),
		RightXor: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					leftToBool, foundError := self.Get(ToBool)
					if foundError != nil {
						return nil, foundError
					}
					if _, ok := leftToBool.(*Function); !ok {
						return nil, errors.NewTypeError(leftToBool.(IObject).TypeName(), FunctionName)
					}
					leftBool, transformationError := p.CallFunction(leftToBool.(*Function), self.SymbolTable(), self)
					if transformationError != nil {
						return nil, transformationError
					}

					left := arguments[0]
					var rightToBool interface{}
					rightToBool, foundError = left.Get(ToBool)
					if foundError != nil {
						return nil, foundError
					}
					if _, ok := rightToBool.(*Function); !ok {
						return nil, errors.NewTypeError(rightToBool.(IObject).TypeName(), FunctionName)
					}
					var rightBool IObject
					rightBool, transformationError = p.CallFunction(rightToBool.(*Function), left.SymbolTable(), left)
					if transformationError != nil {
						return nil, transformationError
					}
					return p.NewBool(p.PeekSymbolTable(), rightBool.GetBool() != leftBool.GetBool()), nil
				},
			),
		),
		Equals: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1, func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
				right := arguments[0]
				return p.NewBool(p.PeekSymbolTable(), self.Id() == right.Id()), nil
			},
			),
		),
		RightEquals: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					return p.NewBool(p.PeekSymbolTable(), left.Id() == self.Id()), nil
				},
			),
		),
		NotEquals: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					return p.NewBool(p.PeekSymbolTable(), self.Id() != right.Id()), nil
				},
			),
		),
		RightNotEquals: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					return p.NewBool(p.PeekSymbolTable(), left.Id() != self.Id()), nil
				},
			),
		),
		GreaterThan: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		RightGreaterThan: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		LessThan: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		RightLessThan: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		GreaterThanOrEqual: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		RightGreaterThanOrEqual: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		LessThanOrEqual: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		RightLessThanOrEqual: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		Hash: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					if self.GetHash() == 0 {
						objectHash, hashingError := p.HashString(fmt.Sprintf("%v-%s-%d", self, self.TypeName(), self.Id()))
						if hashingError != nil {
							return nil, hashingError
						}
						self.SetHash(objectHash)
					}
					return p.NewInteger(p.PeekSymbolTable(), self.GetHash()), nil
				},
			),
		),
		Copy: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(0)),
		Index: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(1)),
		Assign: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(2)),
		Call: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(0)),
		Iter: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(0)),
		Class: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					if self.GetClass() == nil { // This should only happen with built-ins
						class, getError := p.MasterSymbolTable().GetAny(self.TypeName())
						if getError != nil {
							return nil, getError
						}
						if _, ok := class.(*Type); !ok {
							return nil, errors.NewTypeError(class.TypeName(), TypeName)
						}
						self.SetClass(class.(*Type))
					}
					return self.GetClass(), nil
				}),
		),
		SubClasses: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(0)),
		ToInteger: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(0)),
		ToFloat: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(0)),
		ToString: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewString(p.PeekSymbolTable(),
						fmt.Sprintf("%s-%d", self.TypeName(), self.Id())), nil
				},
			),
		),
		ToBool: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(_ IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewBool(p.PeekSymbolTable(), true), nil
				},
			),
		),
		ToArray: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(0)),
		ToTuple: p.NewFunction(object.SymbolTable(),
			NewNotImplementedCallable(0)),
		GetInteger64: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()), nil
				},
			),
		),
		GetBool: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewBool(p.PeekSymbolTable(), self.GetBool()), nil
				},
			),
		),
		GetBytes: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewBytes(p.PeekSymbolTable(), self.GetBytes()), nil
				},
			),
		),
		GetString: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewString(p.PeekSymbolTable(), self.GetString()), nil
				},
			),
		),
		GetFloat64: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewFloat(p.PeekSymbolTable(), self.GetFloat64()), nil
				},
			),
		),
		GetContent: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewArray(p.PeekSymbolTable(), self.GetContent()), nil
				},
			),
		),
		GetKeyValues: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewHashTable(p.PeekSymbolTable(), self.GetKeyValues(), self.GetLength()), nil
				},
			),
		),
		GetLength: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewInteger(p.PeekSymbolTable(), int64(self.GetLength())), nil
				},
			),
		),
		SetBool: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					self.SetBool(arguments[0].GetBool())
					return p.PeekSymbolTable().GetAny(None)
				},
			),
		),
		SetBytes: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					self.SetBytes(arguments[0].GetBytes())
					self.SetLength(arguments[0].GetLength())
					return p.PeekSymbolTable().GetAny(None)
				},
			),
		),
		SetString: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					self.SetString(arguments[0].GetString())
					self.SetLength(arguments[0].GetLength())
					return p.PeekSymbolTable().GetAny(None)
				},
			),
		),
		SetInteger64: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					self.SetInteger64(arguments[0].GetInteger64())
					return p.PeekSymbolTable().GetAny(None)
				},
			),
		),
		SetFloat64: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					self.SetFloat64(arguments[0].GetFloat64())
					return p.PeekSymbolTable().GetAny(None)
				},
			),
		),
		SetContent: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					self.SetContent(arguments[0].GetContent())
					self.SetLength(arguments[0].GetLength())
					return p.PeekSymbolTable().GetAny(None)
				},
			),
		),
		SetKeyValues: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					self.SetKeyValues(arguments[0].GetKeyValues())
					self.SetLength(arguments[0].GetLength())
					return p.PeekSymbolTable().GetAny(None)
				},
			),
		),
		SetLength: p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					self.SetLength(arguments[0].GetLength())
					return p.PeekSymbolTable().GetAny(None)
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
