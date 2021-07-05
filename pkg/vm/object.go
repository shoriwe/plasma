package vm

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/errors"
)

type Object struct {
	isBuiltIn       bool
	id              int64
	typeName        string
	class           *Type
	subClasses      []*Type
	symbols         *SymbolTable
	hash            int64
	Bool            bool
	String          string
	Bytes           []uint8
	Integer         int64
	Float           float64
	Content         []Value
	KeyValues       map[int64][]*KeyValue
	Length          int
	onDemandSymbols map[string]OnDemandLoader
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

func (o *Object) GetInteger() int64 {
	return o.Integer
}

func (o *Object) GetFloat() float64 {
	return o.Float
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

func (o *Object) SetInteger(i int64) {
	o.Integer = i
}

func (o *Object) SetFloat(f float64) {
	o.Float = f
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
	result, getError := o.symbols.GetSelf(symbol)
	if getError != nil {
		loader, found := o.onDemandSymbols[symbol]
		if !found {
			return nil, getError
		}
		result = loader()
		o.Set(symbol, result)
	}
	return result, nil
}

func (o *Object) SetOnDemandSymbol(symbol string, loader OnDemandLoader) {
	o.onDemandSymbols[symbol] = loader
}

func (o *Object) GetOnDemandSymbolLoader(symbol string) OnDemandLoader {
	return o.onDemandSymbols[symbol]
}
func (o *Object) GetOnDemandSymbols() map[string]OnDemandLoader {
	return o.onDemandSymbols
}

func (o *Object) Dir() map[string]byte {
	result := map[string]byte{}
	for symbol := range o.symbols.Symbols {
		result[symbol] = 0
	}
	for symbol := range o.onDemandSymbols {
		result[symbol] = 0
	}
	return result
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

func (o *Object) GetClass(p *Plasma) *Type {
	if o.class == nil { // This should only happen with built-ins
		o.class = p.ForceMasterGetAny(o.typeName).(*Type)
	}
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

func (p *Plasma) NewObject(
	isBuiltIn bool,
	typeName string,
	subClasses []*Type,
	parentSymbols *SymbolTable,
) *Object {
	result := &Object{
		id:              p.NextId(),
		typeName:        typeName,
		subClasses:      subClasses,
		symbols:         NewSymbolTable(parentSymbols),
		isBuiltIn:       isBuiltIn,
		onDemandSymbols: map[string]OnDemandLoader{},
	}
	result.Length = 0
	result.Bool = true
	result.String = ""
	result.Integer = 0
	result.Float = 0
	result.Content = []Value{}
	result.KeyValues = map[int64][]*KeyValue{}
	result.Bytes = []uint8{}
	result.Set(Self, result)
	p.ObjectInitialize(isBuiltIn)(result)
	return result
}

func (p *Plasma) ObjectInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(object Value) *Object {
		object.SetOnDemandSymbol(Initialize,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(_ Value, _ ...Value) (Value, *Object) {
							return p.GetNone(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Negate,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							selfBool, callError := p.QuickGetBool(self)
							if callError != nil {
								return nil, callError
							}
							if selfBool {
								return p.GetFalse(), nil
							}
							return p.GetTrue(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(And,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							leftBool, callError := p.QuickGetBool(self)
							if callError != nil {
								return nil, callError
							}
							var rightBool bool
							rightBool, callError = p.QuickGetBool(arguments[0])
							if callError != nil {
								return nil, callError
							}
							if leftBool && rightBool {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightAnd,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							leftBool, callError := p.QuickGetBool(arguments[0])
							if callError != nil {
								return nil, callError
							}
							var rightBool bool
							rightBool, callError = p.QuickGetBool(self)
							if callError != nil {
								return nil, callError
							}
							if leftBool && rightBool {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Or,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							leftBool, callError := p.QuickGetBool(self)
							if callError != nil {
								return nil, callError
							}
							var rightBool bool
							rightBool, callError = p.QuickGetBool(arguments[0])
							if callError != nil {
								return nil, callError
							}
							if leftBool || rightBool {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightOr,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							leftBool, callError := p.QuickGetBool(arguments[0])
							if callError != nil {
								return nil, callError
							}
							var rightBool bool
							rightBool, callError = p.QuickGetBool(self)
							if callError != nil {
								return nil, callError
							}
							if leftBool || rightBool {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Xor,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							leftBool, callError := p.QuickGetBool(self)
							if callError != nil {
								return nil, callError
							}
							var rightBool bool
							rightBool, callError = p.QuickGetBool(arguments[0])
							if callError != nil {
								return nil, callError
							}
							if leftBool != rightBool {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightXor,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							leftBool, callError := p.QuickGetBool(arguments[0])
							if callError != nil {
								return nil, callError
							}
							var rightBool bool
							rightBool, callError = p.QuickGetBool(self)
							if callError != nil {
								return nil, callError
							}
							if leftBool != rightBool {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Equals,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							if self.Id() == right.Id() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Equals,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							if self.Id() == right.Id() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightEquals,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							if left.Id() == self.Id() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(NotEquals,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							if self.Id() != right.Id() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightNotEquals,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							if left.Id() != self.Id() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Hash,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							if self.GetHash() == 0 {
								objectHash := p.HashString(fmt.Sprintf("%v-%s-%d", self, self.TypeName(), self.Id()))
								self.SetHash(objectHash)
							}
							return p.NewInteger(false, p.PeekSymbolTable(), self.GetHash()), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Class,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return self.GetClass(p), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(SubClasses,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							var subClassesCopy []Value
							for _, class := range self.SubClasses() {
								subClassesCopy = append(subClassesCopy, class)
							}
							return p.NewTuple(false, p.PeekSymbolTable(), subClassesCopy), nil
						},
					),
				)
			},
		)

		object.SetOnDemandSymbol(ToString,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewString(false, p.PeekSymbolTable(),
								fmt.Sprintf("%s{%s}-%X", ObjectName, self.TypeName(), self.Id())), nil
						},
					),
				)
			})
		object.SetOnDemandSymbol(ToBool,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(_ Value, _ ...Value) (Value, *Object) {
							return p.GetTrue(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(GetInteger,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger()), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(GetBool,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							if self.GetBool() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(GetBytes,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewBytes(false, p.PeekSymbolTable(), self.GetBytes()), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(GetString,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewString(false, p.PeekSymbolTable(), self.GetString()), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(GetFloat,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewFloat(false, p.PeekSymbolTable(), self.GetFloat()), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(GetContent,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewArray(false, p.PeekSymbolTable(), self.GetContent()), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(GetKeyValues,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewHashTable(false, p.PeekSymbolTable(), self.GetKeyValues(), self.GetLength()), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(GetLength,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewInteger(false, p.PeekSymbolTable(), int64(self.GetLength())), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(SetBool,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							self.SetBool(arguments[0].GetBool())
							return p.GetNone(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(SetBytes,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							self.SetBytes(arguments[0].GetBytes())
							self.SetLength(arguments[0].GetLength())
							return p.GetNone(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(SetString,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							self.SetString(arguments[0].GetString())
							self.SetLength(arguments[0].GetLength())
							return p.GetNone(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(SetInteger,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							self.SetInteger(arguments[0].GetInteger())
							return p.GetNone(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(SetFloat,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							self.SetFloat(arguments[0].GetFloat())
							return p.GetNone(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(SetContent,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							self.SetContent(arguments[0].GetContent())
							self.SetLength(arguments[0].GetLength())
							return p.GetNone(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(SetKeyValues,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							self.SetKeyValues(arguments[0].GetKeyValues())
							self.SetLength(arguments[0].GetLength())
							return p.GetNone(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(SetLength,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							self.SetLength(arguments[0].GetLength())
							return p.GetNone(), nil
						},
					),
				)
			},
		)
		return nil
	}
}
