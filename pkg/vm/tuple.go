package vm

import (
	"github.com/shoriwe/gplasma/pkg/tools"
)

type Tuple struct {
	*Object
}

func (p *Plasma) NewTuple(context *Context, isBuiltIn bool, parentSymbols *SymbolTable, content []Value) *Tuple {
	tuple := &Tuple{
		Object: p.NewObject(context, isBuiltIn, TupleName, nil, parentSymbols),
	}
	tuple.SetContent(content)
	tuple.SetLength(len(content))
	p.TupleInitialize(isBuiltIn)(context, tuple)
	tuple.SetOnDemandSymbol(Self,
		func() Value {
			return tuple
		},
	)
	return tuple
}

func (p *Plasma) TupleInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object Value) *Object {
		object.SetOnDemandSymbol(Mul,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							switch right.(type) {
							case *Integer:
								content, repetitionError := p.Repeat(context, self.GetContent(), right.GetInteger())
								if repetitionError != nil {
									return nil, repetitionError
								}
								return p.NewTuple(context, false, context.PeekSymbolTable(), content), nil
							default:
								return nil, p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightMul,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							switch left.(type) {
							case *Integer:
								content, repetitionError := p.Repeat(context, self.GetContent(), left.GetInteger())
								if repetitionError != nil {
									return nil, repetitionError
								}
								return p.NewTuple(context, false, context.PeekSymbolTable(), content), nil
							default:
								return nil, p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Equals,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							if _, ok := right.(*Tuple); !ok {
								return p.GetFalse(), nil
							}
							if self.GetLength() != right.GetLength() {
								return p.GetFalse(), nil
							}
							var rightEquals Value
							var comparisonResult Value
							var callError *Object
							var comparisonBool bool

							for i := 0; i < self.GetLength(); i++ {
								leftEquals, getError := self.GetContent()[i].Get(Equals)
								if getError != nil {
									rightEquals, getError = right.GetContent()[i].Get(RightEquals)
									if getError != nil {
										return nil, p.NewObjectWithNameNotFoundError(context, right.GetContent()[i].GetClass(p), RightEquals)
									}
									comparisonResult, callError = p.CallFunction(context, rightEquals, context.PeekSymbolTable(), self.GetContent()[i])
								} else {
									comparisonResult, callError = p.CallFunction(context, leftEquals, context.PeekSymbolTable(), right.GetContent()[i])
								}
								if callError != nil {
									return nil, callError
								}
								comparisonBool, callError = p.QuickGetBool(context, comparisonResult)
								if !comparisonBool {
									return p.GetFalse(), nil
								}
							}
							return p.GetTrue(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightEquals,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							if _, ok := left.(*Tuple); !ok {
								return p.GetFalse(), nil
							}
							if self.GetLength() != left.GetLength() {
								return p.GetFalse(), nil
							}
							var rightEquals Value
							var comparisonResult Value
							var callError *Object
							var comparisonBool bool

							for i := 0; i < self.GetLength(); i++ {
								leftEquals, getError := left.GetContent()[i].Get(Equals)
								if getError != nil {
									rightEquals, getError = self.GetContent()[i].Get(RightEquals)
									if getError != nil {
										return nil, p.NewObjectWithNameNotFoundError(context, self.GetContent()[i].GetClass(p), RightEquals)
									}
									comparisonResult, callError = p.CallFunction(context, rightEquals, context.PeekSymbolTable(), left.GetContent()[i])
								} else {
									comparisonResult, callError = p.CallFunction(context, leftEquals, context.PeekSymbolTable(), self.GetContent()[i])
								}
								if callError != nil {
									return nil, callError
								}
								comparisonBool, callError = p.QuickGetBool(context, comparisonResult)
								if callError != nil {
									return nil, callError
								}
								if !comparisonBool {
									return p.GetFalse(), nil
								}
							}
							return p.GetTrue(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(NotEquals,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {

							right := arguments[0]
							if _, ok := right.(*Tuple); !ok {
								return p.GetTrue(), nil
							}
							if self.GetLength() != right.GetLength() {
								return p.GetTrue(), nil
							}
							var rightNotEquals Value
							var comparisonResult Value
							var callError *Object
							var comparisonBool bool

							for i := 0; i < self.GetLength(); i++ {
								leftNotEquals, getError := self.GetContent()[i].Get(NotEquals)
								if getError != nil {
									rightNotEquals, getError = right.GetContent()[i].Get(RightNotEquals)
									if getError != nil {
										return nil, p.NewObjectWithNameNotFoundError(context, right.GetContent()[i].GetClass(p), RightNotEquals)
									}
									comparisonResult, callError = p.CallFunction(context, rightNotEquals, context.PeekSymbolTable(), self.GetContent()[i])
								} else {
									comparisonResult, callError = p.CallFunction(context, leftNotEquals, context.PeekSymbolTable(), right.GetContent()[i])
								}
								if callError != nil {
									return nil, callError
								}
								comparisonBool, callError = p.QuickGetBool(context, comparisonResult)
								if !comparisonBool {
									return p.GetFalse(), nil
								}
							}
							return p.GetTrue(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightNotEquals,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {

							left := arguments[0]
							if _, ok := left.(*Tuple); !ok {
								return p.GetFalse(), nil
							}
							if self.GetLength() != left.GetLength() {
								return p.GetFalse(), nil
							}
							var rightEquals Value
							var comparisonResult Value
							var callError *Object
							var comparisonBool bool

							for i := 0; i < self.GetLength(); i++ {
								leftEquals, getError := left.GetContent()[i].Get(NotEquals)
								if getError != nil {
									rightEquals, getError = self.GetContent()[i].Get(RightNotEquals)
									if getError != nil {
										return nil, p.NewObjectWithNameNotFoundError(context, self.GetContent()[i].GetClass(p), RightNotEquals)
									}
									comparisonResult, callError = p.CallFunction(context, rightEquals, context.PeekSymbolTable(), left.GetContent()[i])
								} else {
									comparisonResult, callError = p.CallFunction(context, leftEquals, context.PeekSymbolTable(), self.GetContent()[i])
								}
								if callError != nil {
									return nil, callError
								}
								comparisonBool, callError = p.QuickGetBool(context, comparisonResult)
								if !comparisonBool {
									return p.GetFalse(), nil
								}
							}
							return p.GetTrue(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Contains,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							value := arguments[0]
							valueRightEquals, getError := value.Get(RightEquals)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(context, value.GetClass(p), RightEquals)
							}
							for _, tupleValue := range self.GetContent() {
								callResult, callError := p.CallFunction(context, valueRightEquals, value.SymbolTable(), tupleValue)
								if callError != nil {
									return nil, callError
								}
								var boolValue Value
								if _, ok := callResult.(*Bool); ok {
									boolValue = callResult
								} else {
									var boolValueToBool Value
									boolValueToBool, getError = callResult.Get(ToBool)
									if getError != nil {
										return nil, p.NewObjectWithNameNotFoundError(context, callResult.GetClass(p), ToBool)
									}
									callResult, callError = p.CallFunction(context, boolValueToBool, callResult.SymbolTable())
									if callError != nil {
										return nil, callError
									}
									if _, ok = callResult.(*Bool); !ok {
										return nil, p.NewInvalidTypeError(context, callResult.TypeName(), BoolName)
									}
									boolValue = callResult
								}
								if boolValue.(*Bool).GetBool() {
									return p.GetTrue(), nil
								}
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightContains,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							value := arguments[0]
							valueRightEquals, getError := value.Get(Equals)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(context, value.GetClass(p), Equals)
							}
							for _, tupleValue := range self.GetContent() {
								callResult, callError := p.CallFunction(context, valueRightEquals, value.SymbolTable(), tupleValue)
								if callError != nil {
									return nil, callError
								}
								var boolValue Value
								if _, ok := callResult.(*Bool); ok {
									boolValue = callResult
								} else {
									var boolValueToBool Value
									boolValueToBool, getError = callResult.Get(ToBool)
									if getError != nil {
										return nil, p.NewObjectWithNameNotFoundError(context, callResult.GetClass(p), ToBool)
									}
									callResult, callError = p.CallFunction(context, boolValueToBool, callResult.SymbolTable())
									if callError != nil {
										return nil, callError
									}
									if _, ok = callResult.(*Bool); !ok {
										return nil, p.NewInvalidTypeError(context, callResult.TypeName(), BoolName)
									}
									boolValue = callResult
								}
								if boolValue.(*Bool).GetBool() {
									return p.GetTrue(), nil
								}
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Hash,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							tupleHash := XXPrime5 ^ p.Seed()
							for _, contentObject := range self.GetContent() {
								objectHashFunc, getError := contentObject.Get(Hash)
								if getError != nil {
									return nil, p.NewObjectWithNameNotFoundError(context, contentObject.GetClass(p), Hash)
								}
								objectHash, callError := p.CallFunction(context, objectHashFunc, self.SymbolTable())
								if callError != nil {
									return nil, callError
								}
								if _, ok := objectHash.(*Integer); !ok {
									return nil, p.NewInvalidTypeError(context, objectHash.TypeName(), IntegerName)
								}
								tupleHash += uint64(objectHash.GetInteger()) * XXPrime2
								tupleHash = (tupleHash << 31) | (tupleHash >> 33)
								tupleHash *= XXPrime1
								tupleHash &= (1 << 64) - 1
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(), int64(tupleHash)), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Copy,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {

							var copiedObjects []Value
							for _, contentObject := range self.GetContent() {
								objectCopy, getError := contentObject.Get(Copy)
								if getError != nil {
									return nil, p.NewObjectWithNameNotFoundError(context, contentObject.GetClass(p), Copy)
								}
								copiedObject, copyError := p.CallFunction(context, objectCopy, context.PeekSymbolTable())
								if copyError != nil {
									return nil, copyError
								}
								copiedObjects = append(copiedObjects, copiedObject)
							}
							return p.NewTuple(context, false, context.PeekSymbolTable(), copiedObjects), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Index,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							indexObject := arguments[0]
							var ok bool
							if _, ok = indexObject.(*Integer); ok {
								index, calcError := tools.CalcIndex(indexObject.GetInteger(), self.GetLength())
								if calcError != nil {
									return nil, p.NewIndexOutOfRange(context, self.GetLength(), indexObject.GetInteger())
								}
								return self.GetContent()[index], nil
							} else if _, ok = indexObject.(*Tuple); ok {
								if len(indexObject.GetContent()) != 2 {
									return nil, p.NewInvalidNumberOfArgumentsError(context, len(indexObject.GetContent()), 2)
								}
								startIndex, calcError := tools.CalcIndex(indexObject.GetContent()[0].GetInteger(), self.GetLength())
								if calcError != nil {
									return nil, p.NewIndexOutOfRange(context, self.GetLength(), indexObject.GetContent()[0].GetInteger())
								}
								var targetIndex int
								targetIndex, calcError = tools.CalcIndex(indexObject.GetContent()[1].GetInteger(), self.GetLength())
								if calcError != nil {
									return nil, p.NewIndexOutOfRange(context, self.GetLength(), indexObject.GetContent()[1].GetInteger())
								}
								return p.NewTuple(context, false, context.PeekSymbolTable(), self.GetContent()[startIndex:targetIndex]), nil
							} else {
								return nil, p.NewInvalidTypeError(context, indexObject.TypeName(), IntegerName, TupleName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Iter,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							iterator := p.NewIterator(context, false, context.PeekSymbolTable())
							iterator.SetInteger(0)
							iterator.SetContent(self.GetContent())
							iterator.SetLength(self.GetLength())
							iterator.Set(HasNext,
								p.NewFunction(context, isBuiltIn, iterator.SymbolTable(),
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
								p.NewFunction(context, isBuiltIn, iterator.SymbolTable(),
									NewBuiltInClassFunction(iterator,
										0,
										func(funcSelf Value, _ ...Value) (Value, *Object) {
											value := funcSelf.GetContent()[int(funcSelf.GetInteger())]
											funcSelf.SetInteger(funcSelf.GetInteger() + 1)
											return value, nil
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
		object.SetOnDemandSymbol(ToString,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							result := "("
							var objectString Value
							var callError *Object
							for index, contentObject := range self.GetContent() {
								if index != 0 {
									result += ", "
								}
								objectToString, getError := contentObject.Get(ToString)
								if getError != nil {
									return nil, p.NewObjectWithNameNotFoundError(context, contentObject.GetClass(p), ToString)
								}
								objectString, callError = p.CallFunction(context, objectToString, context.PeekSymbolTable())
								if callError != nil {
									return nil, callError
								}
								result += objectString.GetString()
							}
							return p.NewString(context, false, context.PeekSymbolTable(), result+")"), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToBool,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							if self.GetLength() != 0 {
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
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewArray(context, false, context.PeekSymbolTable(), append([]Value{}, self.GetContent()...)), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToTuple,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewTuple(context, false, context.PeekSymbolTable(), append([]Value{}, self.GetContent()...)), nil
						},
					),
				)
			},
		)
		return nil
	}
}
