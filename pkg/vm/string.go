package vm

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/tools"
	"strconv"
	"strings"
)

func (p *Plasma) NewString(context *Context, isBuiltIn bool, parentSymbols *SymbolTable, value string) *Value {
	string_ := p.NewValue(context, isBuiltIn, StringName, nil, parentSymbols)
	string_.BuiltInTypeId = StringId
	string_.SetString(value)
	string_.SetLength(len(value))
	p.StringInitialize(isBuiltIn)(context, string_)
	string_.SetOnDemandSymbol(Self,
		func() *Value {
			return string_
		},
	)
	return string_
}

func (p *Plasma) StringInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object *Value) *Value {
		object.SetOnDemandSymbol(Add,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							if !right.IsTypeById(StringId) {
								return p.NewInvalidTypeError(context, right.TypeName(), StringName), false
							}
							return p.NewString(context, false,
								context.PeekSymbolTable(),
								self.GetString()+right.GetString(),
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightAdd,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							if !left.IsTypeById(StringId) {
								return p.NewInvalidTypeError(context, left.TypeName(), StringName), false
							}
							return p.NewString(context, false,
								context.PeekSymbolTable(),
								left.GetString()+self.GetString(),
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Mul,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							if !right.IsTypeById(IntegerId) {
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName), false
							}
							return p.NewString(context, false,
								context.PeekSymbolTable(),
								tools.Repeat(self.GetString(), right.GetInteger()),
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightMul,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							if !left.IsTypeById(IntegerId) {
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName), false
							}
							return p.NewString(context, false,
								context.PeekSymbolTable(),
								tools.Repeat(self.GetString(), left.GetInteger()),
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Equals,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							if !right.IsTypeById(StringId) {
								return p.GetFalse(), true
							}
							return p.InterpretAsBool(self.GetString() == right.GetString()), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightEquals,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							if !left.IsTypeById(StringId) {
								return p.GetFalse(), true
							}
							return p.InterpretAsBool(left.GetString() == self.GetString()), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(NotEquals,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							if !right.IsTypeById(StringId) {
								return p.GetTrue(), true
							}
							return p.InterpretAsBool(self.GetString() != right.GetString()), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightNotEquals,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							if !left.IsTypeById(StringId) {
								return p.GetTrue(), true
							}
							return p.InterpretAsBool(left.GetString() != self.GetString()), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Hash,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							if self.GetHash() == 0 {
								stringHash := p.HashString(fmt.Sprintf("%s-%s", self.GetString(), StringName))
								self.SetHash(stringHash)
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(), self.GetHash()), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Copy,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.NewString(context, false,
								context.PeekSymbolTable(),
								self.GetString(),
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Index,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							indexObject := arguments[0]
							if indexObject.IsTypeById(IntegerId) {
								index, getIndexError := tools.CalcIndex(indexObject.GetInteger(), self.GetLength())
								if getIndexError != nil {
									return p.NewIndexOutOfRange(context, self.GetLength(), indexObject.GetInteger()), false
								}
								return p.NewString(context, false, context.PeekSymbolTable(), string(self.GetString()[index])), true
							} else if indexObject.IsTypeById(TupleId) {
								if len(indexObject.GetContent()) != 2 {
									return p.NewInvalidNumberOfArgumentsError(context, len(indexObject.GetContent()), 2), false
								}
								startIndex, calcError := tools.CalcIndex(indexObject.GetContent()[0].GetInteger(), self.GetLength())
								if calcError != nil {
									return p.NewIndexOutOfRange(context, self.GetLength(), indexObject.GetContent()[0].GetInteger()), false
								}
								var targetIndex int
								targetIndex, calcError = tools.CalcIndex(indexObject.GetContent()[1].GetInteger(), self.GetLength())
								if calcError != nil {
									return p.NewIndexOutOfRange(context, self.GetLength(), indexObject.GetContent()[1].GetInteger()), false
								}
								return p.NewString(context, false, context.PeekSymbolTable(), self.GetString()[startIndex:targetIndex]), true
							} else {
								return p.NewInvalidTypeError(context, indexObject.TypeName(), IntegerName, TupleName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Iter,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							iterator := p.NewIterator(context, false, context.PeekSymbolTable())
							iterator.SetInteger(0) // This is the index
							iterator.SetString(self.GetString())
							iterator.SetLength(self.GetLength())
							iterator.Set(HasNext,
								p.NewFunction(context, isBuiltIn, iterator.SymbolTable(),
									NewBuiltInClassFunction(iterator,
										0,
										func(funcSelf *Value, _ ...*Value) (*Value, bool) {
											return p.InterpretAsBool(funcSelf.GetInteger() < int64(funcSelf.GetLength())), true
										},
									),
								),
							)
							iterator.Set(Next,
								p.NewFunction(context, isBuiltIn, iterator.SymbolTable(),
									NewBuiltInClassFunction(iterator,
										0,
										func(funcSelf *Value, _ ...*Value) (*Value, bool) {
											char := string([]rune(funcSelf.GetString())[int(funcSelf.GetInteger())])
											funcSelf.SetInteger(funcSelf.GetInteger() + 1)
											return p.NewString(context, false, context.PeekSymbolTable(), char), true
										},
									),
								),
							)
							return iterator, true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToInteger,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							number, parsingError := strconv.ParseInt(strings.ReplaceAll(self.GetString(), "_", ""), 10, 64)
							if parsingError != nil {
								return p.NewIntegerParsingError(context), false
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(), number), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToFloat,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							number, parsingError := strconv.ParseFloat(strings.ReplaceAll(self.GetString(), "_", ""), 64)
							if parsingError != nil {
								return p.NewFloatParsingError(context), false
							}
							return p.NewFloat(context, false, context.PeekSymbolTable(), number), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToString,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.NewString(context, false,
								context.PeekSymbolTable(),
								self.GetString(),
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToBool,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.InterpretAsBool(self.GetLength() > 0), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToArray,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							var content []*Value
							for _, char := range self.GetString() {
								content = append(content, p.NewString(context, false,
									context.PeekSymbolTable(), string(char),
								),
								)
							}
							return p.NewArray(context, false, context.PeekSymbolTable(), content), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToTuple,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							var content []*Value
							for _, char := range self.GetString() {
								content = append(content, p.NewString(context, false,
									context.PeekSymbolTable(), string(char),
								),
								)
							}
							return p.NewTuple(context, false, context.PeekSymbolTable(), content), true
						},
					),
				)
			},
		)
		return nil
	}
}
