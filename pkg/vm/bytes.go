package vm

import (
	"bytes"
	"encoding/binary"
	"github.com/shoriwe/gplasma/pkg/tools"
)

func (p *Plasma) NewBytes(context *Context, isBuiltIn bool, parent *SymbolTable, content []uint8) *Value {
	bytes_ := p.NewValue(context, isBuiltIn, BytesName, nil, parent)
	bytes_.BuiltInTypeId = BytesId
	bytes_.SetBytes(content)
	bytes_.SetLength(len(content))
	p.BytesInitialize(isBuiltIn)(context, bytes_)
	bytes_.SetOnDemandSymbol(Self,
		func() *Value {
			return bytes_
		},
	)
	return bytes_
}

func (p *Plasma) BytesInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object *Value) *Value {
		object.SetOnDemandSymbol(Add,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							if !right.IsTypeById(BytesId) {
								return p.NewInvalidTypeError(context, right.TypeName(), BytesName), false
							}
							var newContent []uint8
							copy(newContent, self.GetBytes())
							newContent = append(newContent, right.GetBytes()...)
							return p.NewBytes(context, false, context.PeekSymbolTable(), newContent), true
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
							if !left.IsTypeById(BytesId) {
								return p.NewInvalidTypeError(context, left.TypeName(), BytesName), false
							}
							var newContent []uint8
							copy(newContent, left.GetBytes())
							newContent = append(newContent, self.GetBytes()...)
							return p.NewBytes(context, false, context.PeekSymbolTable(), newContent), true
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
							if !right.IsTypeById(BytesId) {
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName), false
							}
							return p.NewBytes(context, false, context.PeekSymbolTable(), bytes.Repeat(self.GetBytes(), int(right.GetInteger()))), true
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
							return p.NewBytes(context, false, context.PeekSymbolTable(), bytes.Repeat(left.GetBytes(), int(self.GetInteger()))), true
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
							if !right.IsTypeById(BytesId) {
								return p.GetFalse(), true
							}
							if self.GetLength() != right.GetLength() {
								return p.GetFalse(), true
							}
							return p.InterpretAsBool(bytes.Compare(self.GetBytes(), right.GetBytes()) == 0), true
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
							if !left.IsTypeById(BytesId) {
								return p.GetFalse(), true
							}
							if left.GetLength() != self.GetLength() {
								return p.GetFalse(), true
							}
							return p.InterpretAsBool(bytes.Compare(left.GetBytes(), self.GetBytes()) == 0), true
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
							if !right.IsTypeById(BytesId) {
								return p.GetFalse(), true
							}
							if self.GetLength() != right.GetLength() {
								return p.GetFalse(), true
							}
							return p.InterpretAsBool(bytes.Compare(self.GetBytes(), right.GetBytes()) != 0), true
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
							if !left.IsTypeById(BytesId) {
								return p.GetFalse(), true
							}
							if left.GetLength() != self.GetLength() {
								return p.GetFalse(), true
							}
							return p.InterpretAsBool(bytes.Compare(left.GetBytes(), self.GetBytes()) != 0), true
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
							selfHash := p.HashBytes(append(self.GetBytes(), []byte("Bytes")...))
							return p.NewInteger(context, false, context.PeekSymbolTable(), selfHash), true
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
							var newBytes []uint8
							copy(newBytes, self.GetBytes())
							return p.NewBytes(context, false, context.PeekSymbolTable(), newBytes), true
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
								index, calcError := tools.CalcIndex(indexObject.GetInteger(), self.GetLength())
								if calcError != nil {
									return p.NewIndexOutOfRange(context, self.GetLength(), indexObject.GetInteger()), false
								}
								return p.NewInteger(context, false, context.PeekSymbolTable(), int64(self.GetBytes()[index])), true
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
								return p.NewBytes(context, false, context.PeekSymbolTable(), self.GetBytes()[startIndex:targetIndex]), true
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
							iterator.SetBytes(self.GetBytes())
							iterator.SetLength(self.GetLength())
							iterator.Set(HasNext,
								p.NewFunction(context, isBuiltIn, iterator.SymbolTable(),
									NewBuiltInClassFunction(iterator,
										0,
										func(funcSelf *Value, _ ...*Value) (*Value, bool) {
											return p.InterpretAsBool(int(funcSelf.GetInteger()) < funcSelf.GetLength()), true
										},
									),
								),
							)
							iterator.Set(Next,
								p.NewFunction(context, isBuiltIn, iterator.SymbolTable(),
									NewBuiltInClassFunction(iterator,
										0,
										func(funcSelf *Value, _ ...*Value) (*Value, bool) {
											char := funcSelf.GetBytes()[int(funcSelf.GetInteger())]
											funcSelf.SetInteger(funcSelf.GetInteger() + 1)
											return p.NewInteger(context, false, context.PeekSymbolTable(), int64(char)), true
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
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								int64(binary.BigEndian.Uint32(self.GetBytes())),
							), true
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
							return p.NewString(context, false, context.PeekSymbolTable(), string(self.GetBytes())), true
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
							return p.InterpretAsBool(self.GetLength() != 0), true
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
							var newContent []*Value
							for _, byte_ := range self.GetBytes() {
								newContent = append(newContent,
									p.NewInteger(context, false, context.PeekSymbolTable(),
										int64(byte_),
									),
								)
							}
							return p.NewArray(context, false, context.PeekSymbolTable(), newContent), true
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
							var newContent []*Value
							for _, byte_ := range self.GetBytes() {
								newContent = append(newContent,
									p.NewInteger(context, false, context.PeekSymbolTable(),
										int64(byte_),
									),
								)
							}
							return p.NewTuple(context, false, context.PeekSymbolTable(), newContent), true
						},
					),
				)
			},
		)
		return nil
	}
}
