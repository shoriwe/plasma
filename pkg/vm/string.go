package vm

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/tools"
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
		object.Set(Add,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)
		object.Set(RightAdd,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)
		object.Set(Mul,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						if _, ok := right.(*Integer); !ok {
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName)
						}
						return p.NewString(false,
							p.PeekSymbolTable(),
							tools.Repeat(self.GetString(), right.GetInteger64()),
						), nil
					},
				),
			),
		)
		object.Set(RightMul,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						if _, ok := left.(*Integer); !ok {
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName)
						}
						return p.NewString(false,
							p.PeekSymbolTable(),
							tools.Repeat(self.GetString(), left.GetInteger64()),
						), nil
					},
				),
			),
		)
		object.Set(Equals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						if _, ok := right.(*String); !ok {
							return p.NewBool(false, p.PeekSymbolTable(), false), nil
						}
						return p.NewBool(false, p.PeekSymbolTable(), self.GetString() == right.GetString()), nil
					},
				),
			),
		)
		object.Set(RightEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						if _, ok := left.(*String); !ok {
							return p.NewBool(false, p.PeekSymbolTable(), false), nil
						}
						return p.NewBool(false, p.PeekSymbolTable(), left.GetString() == self.GetString()), nil
					},
				),
			),
		)
		object.Set(NotEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						if _, ok := right.(*String); !ok {
							return p.NewBool(false, p.PeekSymbolTable(), true), nil
						}
						return p.NewBool(false, p.PeekSymbolTable(), self.GetString() != right.GetString()), nil
					},
				),
			),
		)
		object.Set(RightNotEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						if _, ok := left.(*String); !ok {
							return p.NewBool(false, p.PeekSymbolTable(), true), nil
						}
						return p.NewBool(false, p.PeekSymbolTable(), left.GetString() != self.GetString()), nil
					},
				),
			),
		)
		object.Set(Hash,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						if self.GetHash() == 0 {
							stringHash := p.HashString(fmt.Sprintf("%s-%s", self.GetString(), StringName))
							self.SetHash(stringHash)
						}
						return p.NewInteger(false, p.PeekSymbolTable(), self.GetHash()), nil
					},
				),
			),
		)
		object.Set(Copy,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewString(false,
							p.PeekSymbolTable(),
							self.GetString(),
						), nil
					},
				),
			),
		)
		object.Set(Index,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						indexObject := arguments[0]
						if _, ok := indexObject.(*Integer); ok {
							index, getIndexError := tools.CalcIndex(indexObject.GetInteger64(), self.GetLength())
							if getIndexError != nil {
								return nil, p.NewIndexOutOfRange(self.GetLength(), indexObject.GetInteger64())
							}
							return p.NewString(false, p.PeekSymbolTable(), string(self.GetString()[index])), nil
						} else if _, ok = indexObject.(*Tuple); ok {
							if len(indexObject.GetContent()) != 2 {
								return nil, p.NewInvalidNumberOfArgumentsError(len(indexObject.GetContent()), 2)
							}
							startIndex, calcError := tools.CalcIndex(indexObject.GetContent()[0].GetInteger64(), self.GetLength())
							if calcError != nil {
								return nil, p.NewIndexOutOfRange(self.GetLength(), indexObject.GetContent()[0].GetInteger64())
							}
							var targetIndex int
							targetIndex, calcError = tools.CalcIndex(indexObject.GetContent()[1].GetInteger64(), self.GetLength())
							if calcError != nil {
								return nil, p.NewIndexOutOfRange(self.GetLength(), indexObject.GetContent()[1].GetInteger64())
							}
							return p.NewString(false, p.PeekSymbolTable(), self.GetString()[startIndex:targetIndex]), nil
						} else {
							return nil, p.NewInvalidTypeError(indexObject.TypeName(), IntegerName, TupleName)
						}
					},
				),
			),
		)
		object.Set(Iter,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						iterator := p.NewIterator(false, p.PeekSymbolTable())
						iterator.SetInteger64(0) // This is the index
						iterator.SetString(self.GetString())
						iterator.SetLength(self.GetLength())
						iterator.Set(HasNext,
							p.NewFunction(isBuiltIn, iterator.SymbolTable(),
								NewBuiltInClassFunction(iterator,
									0,
									func(funcSelf Value, _ ...Value) (Value, *Object) {
										if int(funcSelf.GetInteger64()) < funcSelf.GetLength() {
											return p.NewBool(false, p.PeekSymbolTable(), true), nil
										}
										return p.NewBool(false, p.PeekSymbolTable(), false), nil
									},
								),
							),
						)
						iterator.Set(Next,
							p.NewFunction(isBuiltIn, iterator.SymbolTable(),
								NewBuiltInClassFunction(iterator,
									0,
									func(funcSelf Value, _ ...Value) (Value, *Object) {
										char := string([]rune(funcSelf.GetString())[int(funcSelf.GetInteger64())])
										funcSelf.SetInteger64(funcSelf.GetInteger64() + 1)
										return p.NewString(false, p.PeekSymbolTable(), char), nil
									},
								),
							),
						)
						return iterator, nil
					},
				),
			),
		)
		object.Set(ToInteger,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						number, parsingError := tools.ParseInteger(self.GetString())
						if parsingError != nil {
							return nil, p.NewIntegerParsingError()
						}
						return p.NewInteger(false, p.PeekSymbolTable(), number), nil
					},
				),
			),
		)
		object.Set(ToFloat,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						number, parsingError := tools.ParseFloat(strings.ReplaceAll(self.GetString(), "_", ""))
						if parsingError != nil {
							return nil, p.NewFloatParsingError()
						}
						return p.NewFloat(false, p.PeekSymbolTable(), number), nil
					},
				),
			),
		)
		object.Set(ToString,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewString(false,
							p.PeekSymbolTable(),
							self.GetString(),
						), nil
					},
				),
			),
		)
		object.Set(ToBool,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewBool(false, p.PeekSymbolTable(), self.GetLength() != 0), nil
					},
				),
			),
		)
		object.Set(ToArray,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)
		object.Set(ToTuple,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)
		return nil
	}
}
