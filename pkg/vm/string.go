package vm

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/tools"
	"strconv"
	"strings"
)

type String struct {
	*Object
}

func (p *Plasma) NewString(isBuiltIn bool, parentSymbols *SymbolTable, value string) *String {
	string_ := &String{
		Object: p.NewObject(isBuiltIn, StringName, nil, parentSymbols),
	}
	string_.SetString(value)
	string_.SetLength(len(value))
	p.StringInitialize(isBuiltIn)(string_)
	string_.Set(Self, string_)
	return string_
}

func (p *Plasma) StringInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(object Value) *Object {
		object.SetOnDemandSymbol(Add,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							if _, ok := right.(*String); !ok {
								return nil, p.NewInvalidTypeError(right.TypeName(), StringName)
							}
							return p.NewString(false,
								p.PeekSymbolTable(),
								self.GetString()+right.GetString(),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightAdd,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							if _, ok := left.(*String); !ok {
								return nil, p.NewInvalidTypeError(left.TypeName(), StringName)
							}
							return p.NewString(false,
								p.PeekSymbolTable(),
								left.GetString()+self.GetString(),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Mul,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							if _, ok := right.(*Integer); !ok {
								return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName)
							}
							return p.NewString(false,
								p.PeekSymbolTable(),
								tools.Repeat(self.GetString(), right.GetInteger()),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightMul,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							if _, ok := left.(*Integer); !ok {
								return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName)
							}
							return p.NewString(false,
								p.PeekSymbolTable(),
								tools.Repeat(self.GetString(), left.GetInteger()),
							), nil
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
							if _, ok := right.(*String); !ok {
								return p.GetFalse(), nil
							}
							if self.GetString() == right.GetString() {
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
							if _, ok := left.(*String); !ok {
								return p.GetFalse(), nil
							}
							if left.GetString() == self.GetString() {
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
							if _, ok := right.(*String); !ok {
								return p.GetTrue(), nil
							}
							if self.GetString() != right.GetString() {
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
							if _, ok := left.(*String); !ok {
								return p.GetTrue(), nil
							}
							if left.GetString() != self.GetString() {
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
								stringHash := p.HashString(fmt.Sprintf("%s-%s", self.GetString(), StringName))
								self.SetHash(stringHash)
							}
							return p.NewInteger(false, p.PeekSymbolTable(), self.GetHash()), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Copy,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewString(false,
								p.PeekSymbolTable(),
								self.GetString(),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Index,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							indexObject := arguments[0]
							if _, ok := indexObject.(*Integer); ok {
								index, getIndexError := tools.CalcIndex(indexObject.GetInteger(), self.GetLength())
								if getIndexError != nil {
									return nil, p.NewIndexOutOfRange(self.GetLength(), indexObject.GetInteger())
								}
								return p.NewString(false, p.PeekSymbolTable(), string(self.GetString()[index])), nil
							} else if _, ok = indexObject.(*Tuple); ok {
								if len(indexObject.GetContent()) != 2 {
									return nil, p.NewInvalidNumberOfArgumentsError(len(indexObject.GetContent()), 2)
								}
								startIndex, calcError := tools.CalcIndex(indexObject.GetContent()[0].GetInteger(), self.GetLength())
								if calcError != nil {
									return nil, p.NewIndexOutOfRange(self.GetLength(), indexObject.GetContent()[0].GetInteger())
								}
								var targetIndex int
								targetIndex, calcError = tools.CalcIndex(indexObject.GetContent()[1].GetInteger(), self.GetLength())
								if calcError != nil {
									return nil, p.NewIndexOutOfRange(self.GetLength(), indexObject.GetContent()[1].GetInteger())
								}
								return p.NewString(false, p.PeekSymbolTable(), self.GetString()[startIndex:targetIndex]), nil
							} else {
								return nil, p.NewInvalidTypeError(indexObject.TypeName(), IntegerName, TupleName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Iter,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							iterator := p.NewIterator(false, p.PeekSymbolTable())
							iterator.SetInteger(0) // This is the index
							iterator.SetString(self.GetString())
							iterator.SetLength(self.GetLength())
							iterator.Set(HasNext,
								p.NewFunction(isBuiltIn, iterator.SymbolTable(),
									NewBuiltInClassFunction(iterator,
										0,
										func(funcSelf Value, _ ...Value) (Value, *Object) {
											if funcSelf.GetInteger() < int64(funcSelf.GetLength()) {
												return p.GetTrue(), nil
											}
											return p.GetFalse(), nil
										},
									),
								),
							)
							iterator.Set(Next,
								p.NewFunction(isBuiltIn, iterator.SymbolTable(),
									NewBuiltInClassFunction(iterator,
										0,
										func(funcSelf Value, _ ...Value) (Value, *Object) {
											char := string([]rune(funcSelf.GetString())[int(funcSelf.GetInteger())])
											funcSelf.SetInteger(funcSelf.GetInteger() + 1)
											return p.NewString(false, p.PeekSymbolTable(), char), nil
										},
									),
								),
							)
							return iterator, nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToInteger,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							number, parsingError := strconv.ParseInt(strings.ReplaceAll(self.GetString(), "_", ""), 10, 64)
							if parsingError != nil {
								return nil, p.NewIntegerParsingError()
							}
							return p.NewInteger(false, p.PeekSymbolTable(), number), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToFloat,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							number, parsingError := strconv.ParseFloat(strings.ReplaceAll(self.GetString(), "_", ""), 64)
							if parsingError != nil {
								return nil, p.NewFloatParsingError()
							}
							return p.NewFloat(false, p.PeekSymbolTable(), number), nil
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
							return p.NewString(false,
								p.PeekSymbolTable(),
								self.GetString(),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToBool,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							if self.GetLength() > 0 {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToArray,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							var content []Value
							for _, char := range self.GetString() {
								content = append(content, p.NewString(false,
									p.PeekSymbolTable(), string(char),
								),
								)
							}
							return p.NewArray(false, p.PeekSymbolTable(), content), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToTuple,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							var content []Value
							for _, char := range self.GetString() {
								content = append(content, p.NewString(false,
									p.PeekSymbolTable(), string(char),
								),
								)
							}
							return p.NewTuple(false, p.PeekSymbolTable(), content), nil
						},
					),
				)
			},
		)
		return nil
	}
}
