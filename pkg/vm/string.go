package vm

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/tools"
	"strings"
)

type String struct {
	*Object
}

func (p *Plasma) NewString(parentSymbols *SymbolTable, value string) *String {
	string_ := &String{
		Object: p.NewObject(StringName, nil, parentSymbols),
	}
	string_.SetString(value)
	string_.SetLength(len(value))
	p.StringInitialize(string_)
	return string_
}

func (p *Plasma) StringInitialize(object IObject) *Object {
	object.Set(Add,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					right := arguments[0]
					return p.NewString(
						p.PeekSymbolTable(),
						self.GetString()+right.GetString(),
					), nil
				},
			),
		),
	)
	object.Set(RightAdd,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					left := arguments[0]
					return p.NewString(
						p.PeekSymbolTable(),
						left.GetString()+self.GetString(),
					), nil
				},
			),
		),
	)
	object.Set(Mul,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					right := arguments[0]
					if _, ok := right.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName)
					}
					return p.NewString(
						p.PeekSymbolTable(),
						tools.Repeat(self.GetString(), right.GetInteger64()),
					), nil
				},
			),
		),
	)
	object.Set(RightMul,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					left := arguments[0]
					if _, ok := left.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName)
					}
					return p.NewString(
						p.PeekSymbolTable(),
						tools.Repeat(self.GetString(), left.GetInteger64()),
					), nil
				},
			),
		),
	)
	object.Set(Equals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					right := arguments[0]
					return p.NewBool(p.PeekSymbolTable(), self.GetString() == right.GetString()), nil
				},
			),
		),
	)
	object.Set(RightEquals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					left := arguments[0]
					return p.NewBool(p.PeekSymbolTable(), left.GetString() == self.GetString()), nil
				},
			),
		),
	)
	object.Set(NotEquals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					right := arguments[0]
					return p.NewBool(p.PeekSymbolTable(), self.GetString() != right.GetString()), nil
				},
			),
		),
	)
	object.Set(RightNotEquals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					left := arguments[0]
					return p.NewBool(p.PeekSymbolTable(), left.GetString() != self.GetString()), nil
				},
			),
		),
	)
	object.Set(Hash,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					if self.GetHash() == 0 {
						stringHash := p.HashString(fmt.Sprintf("%s-%s", self.GetString(), StringName))
						self.SetHash(stringHash)
					}
					return p.NewInteger(p.PeekSymbolTable(), self.GetHash()), nil
				},
			),
		),
	)
	object.Set(Copy,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewString(
						p.PeekSymbolTable(),
						self.GetString(),
					), nil
				},
			),
		),
	)
	object.Set(Index,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					indexObject := arguments[0]
					if _, ok := indexObject.(*Integer); ok {
						index, getIndexError := tools.CalcIndex(indexObject.GetInteger64(), self.GetLength())
						if getIndexError != nil {
							return nil, p.NewIndexOutOfRange(self.GetLength(), indexObject.GetInteger64())
						}
						return p.NewString(p.PeekSymbolTable(), string(self.GetString()[index])), nil
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
						return p.NewString(p.PeekSymbolTable(), self.GetString()[startIndex:targetIndex]), nil
					} else {
						return nil, p.NewInvalidTypeError(indexObject.TypeName(), IntegerName, TupleName)
					}
				},
			),
		),
	)
	object.Set(Iter,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					iterator := p.NewIterator(p.PeekSymbolTable())
					iterator.SetInteger64(0) // This is the index
					iterator.SetString(self.GetString())
					iterator.SetLength(self.GetLength())
					iterator.Set(HasNext,
						p.NewFunction(iterator.SymbolTable(),
							NewBuiltInClassFunction(iterator,
								0,
								func(funcSelf IObject, _ ...IObject) (IObject, *Object) {
									if int(funcSelf.GetInteger64()) < funcSelf.GetLength() {
										return p.NewBool(p.PeekSymbolTable(), true), nil
									}
									return p.NewBool(p.PeekSymbolTable(), false), nil
								},
							),
						),
					)
					iterator.Set(Next,
						p.NewFunction(iterator.SymbolTable(),
							NewBuiltInClassFunction(iterator,
								0,
								func(funcSelf IObject, _ ...IObject) (IObject, *Object) {
									char := string([]rune(funcSelf.GetString())[int(funcSelf.GetInteger64())])
									funcSelf.SetInteger64(funcSelf.GetInteger64() + 1)
									return p.NewString(p.PeekSymbolTable(), char), nil
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
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					number, parsingError := tools.ParseInteger(self.GetString())
					if parsingError != nil {
						return nil, p.NewIntegerParsingError()
					}
					return p.NewInteger(p.PeekSymbolTable(), number), nil
				},
			),
		),
	)
	object.Set(ToFloat,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					number, parsingError := tools.ParseFloat(strings.ReplaceAll(self.GetString(), "_", ""))
					if parsingError != nil {
						return nil, p.NewFloatParsingError()
					}
					return p.NewFloat(p.PeekSymbolTable(), number), nil
				},
			),
		),
	)
	object.Set(ToString,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewString(
						p.PeekSymbolTable(),
						self.GetString(),
					), nil
				},
			),
		),
	)
	object.Set(ToBool,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewBool(p.PeekSymbolTable(), self.GetLength() != 0), nil
				},
			),
		),
	)
	object.Set(ToArray,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					var content []IObject
					for _, char := range self.GetString() {
						content = append(content, p.NewString(
							p.PeekSymbolTable(), string(char),
						),
						)
					}
					return p.NewArray(p.PeekSymbolTable(), content), nil
				},
			),
		),
	)
	object.Set(ToTuple,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					var content []IObject
					for _, char := range self.GetString() {
						content = append(content, p.NewString(
							p.PeekSymbolTable(), string(char),
						),
						)
					}
					return p.NewTuple(p.PeekSymbolTable(), content), nil
				},
			),
		),
	)
	return nil
}
