package vm

import (
	"github.com/shoriwe/gplasma/pkg/tools"
)

type Array struct {
	*Object
}

func (p *Plasma) NewArray(isBuiltIn bool, parentSymbols *SymbolTable, content []Value) *Array {
	array := &Array{
		Object: p.NewObject(isBuiltIn, ArrayName, nil, parentSymbols),
	}
	array.SetContent(content)
	array.SetLength(len(content))
	p.ArrayInitialize(isBuiltIn)(array)
	array.Set(Self, array)
	return array
}

func (p *Plasma) ArrayInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(object Value) *Object {
		object.SetOnDemandSymbol(Mul,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							switch right.(type) {
							case *Integer:
								content, repetitionError := p.Repeat(self.GetContent(),
									right.GetInteger())
								if repetitionError != nil {
									return nil, repetitionError
								}
								return p.NewArray(false, p.PeekSymbolTable(), content), nil
							default:
								return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName)
							}
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
							switch left.(type) {
							case *Integer:
								content, repetitionError := p.Repeat(self.GetContent(),
									left.GetInteger(),
								)
								if repetitionError != nil {
									return nil, repetitionError
								}
								return p.NewArray(false, p.PeekSymbolTable(), content), nil
							default:
								return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName)
							}
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
							if _, ok := right.(*Array); !ok {
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
										return nil, p.NewObjectWithNameNotFoundError(right.GetClass(p), RightEquals)
									}
									comparisonResult, callError = p.CallFunction(rightEquals, p.PeekSymbolTable(), self.GetContent()[i])
								} else {
									comparisonResult, callError = p.CallFunction(leftEquals, p.PeekSymbolTable(), right.GetContent()[i])
								}
								if callError != nil {
									return nil, callError
								}
								comparisonBool, callError = p.QuickGetBool(comparisonResult)
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
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							if _, ok := left.(*Array); !ok {
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
										return nil, p.NewObjectWithNameNotFoundError(self.GetContent()[i].GetClass(p), RightEquals)
									}
									comparisonResult, callError = p.CallFunction(rightEquals, p.PeekSymbolTable(), left.GetContent()[i])
								} else {
									comparisonResult, callError = p.CallFunction(leftEquals, p.PeekSymbolTable(), self.GetContent()[i])
								}
								if callError != nil {
									return nil, callError
								}
								comparisonBool, callError = p.QuickGetBool(comparisonResult)
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
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							if _, ok := right.(*Array); !ok {
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
										return nil, p.NewObjectWithNameNotFoundError(right.GetContent()[i].GetClass(p), RightNotEquals)
									}
									comparisonResult, callError = p.CallFunction(rightNotEquals, p.PeekSymbolTable(), self.GetContent()[i])
								} else {
									comparisonResult, callError = p.CallFunction(leftNotEquals, p.PeekSymbolTable(), right.GetContent()[i])
								}
								if callError != nil {
									return nil, callError
								}
								comparisonBool, callError = p.QuickGetBool(comparisonResult)
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
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							if _, ok := left.(*Array); !ok {
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
										return nil, p.NewObjectWithNameNotFoundError(self.GetContent()[i].GetClass(p), RightNotEquals)
									}
									comparisonResult, callError = p.CallFunction(rightEquals, p.PeekSymbolTable(), left.GetContent()[i])
								} else {
									comparisonResult, callError = p.CallFunction(leftEquals, p.PeekSymbolTable(), self.GetContent()[i])
								}
								if callError != nil {
									return nil, callError
								}
								comparisonBool, callError = p.QuickGetBool(comparisonResult)
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
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							value := arguments[0]
							valueRightEquals, getError := value.Get(RightEquals)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(value.GetClass(p), RightEquals)
							}
							for _, tupleValue := range self.GetContent() {
								callResult, callError := p.CallFunction(valueRightEquals, value.SymbolTable(), tupleValue)
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
										return nil, p.NewObjectWithNameNotFoundError(callError.GetClass(p), ToBool)
									}
									callResult, callError = p.CallFunction(boolValueToBool, callResult.SymbolTable())
									if callError != nil {
										return nil, callError
									}
									if _, ok = callResult.(*Bool); !ok {
										return nil, p.NewInvalidTypeError(callResult.TypeName(), BoolName)
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
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							value := arguments[0]
							valueRightEquals, getError := value.Get(RightEquals)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(value.GetClass(p), RightEquals)
							}
							for _, tupleValue := range self.GetContent() {
								callResult, callError := p.CallFunction(valueRightEquals, value.SymbolTable(), tupleValue)
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
										return nil, p.NewObjectWithNameNotFoundError(callResult.GetClass(p), ToBool)
									}
									callResult, callError = p.CallFunction(boolValueToBool, callResult.SymbolTable())
									if callError != nil {
										return nil, callError
									}
									if _, ok = callResult.(*Bool); !ok {
										return nil, p.NewInvalidTypeError(callResult.TypeName(), BoolName)
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
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(_ Value, _ ...Value) (Value, *Object) {
							return nil, p.NewUnhashableTypeError(object.GetClass(p))
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
							var copiedObjects []Value
							for _, contentObject := range self.GetContent() {
								objectCopy, getError := contentObject.Get(Copy)
								if getError != nil {
									return nil, p.NewObjectWithNameNotFoundError(contentObject.GetClass(p), Copy)
								}
								copiedObject, copyError := p.CallFunction(objectCopy, p.PeekSymbolTable())
								if copyError != nil {
									return nil, copyError
								}
								copiedObjects = append(copiedObjects, copiedObject)
							}
							return p.NewArray(false, p.PeekSymbolTable(), copiedObjects), nil
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
							var ok bool
							if _, ok = indexObject.(*Integer); ok {
								index, calcError := tools.CalcIndex(indexObject.GetInteger(), self.GetLength())
								if calcError != nil {
									return nil, p.NewIndexOutOfRange(self.GetLength(), indexObject.GetInteger())
								}
								return self.GetContent()[index], nil
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
								return p.NewArray(false, p.PeekSymbolTable(), self.GetContent()[startIndex:targetIndex]), nil
							} else {
								return nil, p.NewInvalidTypeError(indexObject.TypeName(), IntegerName, TupleName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Assign,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 2,
						func(self Value, arguments ...Value) (Value, *Object) {
							index, calcError := tools.CalcIndex(arguments[0].GetInteger(), self.GetLength())
							if calcError != nil {
								return nil, p.NewIndexOutOfRange(self.GetLength(), arguments[0].GetInteger())
							}
							self.GetContent()[index] = arguments[1]
							return p.GetNone(), nil
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
							iterator.SetInteger(0)
							iterator.SetContent(self.GetContent())
							iterator.SetLength(self.GetLength())
							iterator.Set(HasNext,
								p.NewFunction(isBuiltIn, iterator.SymbolTable(),
									NewBuiltInClassFunction(iterator,
										0,
										func(funcSelf Value, _ ...Value) (Value, *Object) {
											if funcSelf.GetLength() != self.GetLength() {
												funcSelf.SetLength(self.GetLength())
											}
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
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							result := "["
							var objectString Value
							var callError *Object
							for index, contentObject := range self.GetContent() {
								if index != 0 {
									result += ", "
								}
								objectToString, getError := contentObject.Get(ToString)
								if getError != nil {
									return nil, p.NewObjectWithNameNotFoundError(contentObject.GetClass(p), ToString)
								}
								objectString, callError = p.CallFunction(objectToString, p.PeekSymbolTable())
								if callError != nil {
									return nil, callError
								}
								result += objectString.GetString()
							}
							return p.NewString(false, p.PeekSymbolTable(), result+"]"), nil
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
							if self.GetLength() != 0 {
								return p.GetTrue(), nil
							}
							return p.GetTrue(), nil
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
							return p.NewArray(false, p.PeekSymbolTable(), append([]Value{}, self.GetContent()...)), nil
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
							return p.NewTuple(false, p.PeekSymbolTable(), append([]Value{}, self.GetContent()...)), nil
						},
					),
				)
			},
		)
		return nil
	}
}
