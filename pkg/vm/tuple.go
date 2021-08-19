package vm

import (
	"github.com/shoriwe/gplasma/pkg/tools"
)

func (p *Plasma) NewTuple(context *Context, isBuiltIn bool, parentSymbols *SymbolTable, content []*Value) *Value {
	tuple := p.NewValue(context, isBuiltIn, TupleName, nil, parentSymbols)
	tuple.BuiltInTypeId = TupleId
	tuple.SetContent(content)
	tuple.SetLength(len(content))
	p.TupleInitialize(isBuiltIn)(context, tuple)
	tuple.SetOnDemandSymbol(Self,
		func() *Value {
			return tuple
		},
	)
	return tuple
}

func (p *Plasma) TupleInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object *Value) *Value {
		object.SetOnDemandSymbol(Mul,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								content, repetitionError := p.Repeat(context, self.GetContent(), right.GetInteger())
								if repetitionError != nil {
									return repetitionError, false
								}
								return p.NewTuple(context, false, context.PeekSymbolTable(), content), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName), false
							}
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
							switch left.BuiltInTypeId {
							case IntegerId:
								content, repetitionError := p.Repeat(context, self.GetContent(), left.GetInteger())
								if repetitionError != nil {
									return repetitionError, false
								}
								return p.NewTuple(context, false, context.PeekSymbolTable(), content), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName), false
							}
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
							if !right.IsTypeById(TupleId) {
								return p.GetFalse(), true
							}
							if self.GetLength() != right.GetLength() {
								return p.GetFalse(), true
							}
							var rightEquals *Value
							var comparisonResult *Value
							var success bool
							var callError *Value
							var comparisonBool bool

							for i := 0; i < self.GetLength(); i++ {
								leftEquals, getError := self.GetContent()[i].Get(Equals)
								if getError != nil {
									rightEquals, getError = right.GetContent()[i].Get(RightEquals)
									if getError != nil {
										return p.NewObjectWithNameNotFoundError(context, right.GetContent()[i].GetClass(p), RightEquals), false
									}
									comparisonResult, success = p.CallFunction(context, rightEquals, self.GetContent()[i])
								} else {
									comparisonResult, success = p.CallFunction(context, leftEquals, right.GetContent()[i])
								}
								if !success {
									return comparisonResult, false
								}
								comparisonBool, callError = p.QuickGetBool(context, comparisonResult)
								if callError != nil {
									return callError, false
								}
								if !comparisonBool {
									return p.GetFalse(), true
								}
							}
							return p.GetTrue(), true
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
							if !left.IsTypeById(TupleId) {
								return p.GetFalse(), true
							}
							if self.GetLength() != left.GetLength() {
								return p.GetFalse(), true
							}
							var rightEquals *Value
							var comparisonResult *Value
							var success bool
							var comparisonBool bool

							for i := 0; i < self.GetLength(); i++ {
								leftEquals, getError := left.GetContent()[i].Get(Equals)
								if getError != nil {
									rightEquals, getError = self.GetContent()[i].Get(RightEquals)
									if getError != nil {
										return p.NewObjectWithNameNotFoundError(context, self.GetContent()[i].GetClass(p), RightEquals), false
									}
									comparisonResult, success = p.CallFunction(context, rightEquals, left.GetContent()[i])
								} else {
									comparisonResult, success = p.CallFunction(context, leftEquals, self.GetContent()[i])
								}
								if !success {
									return comparisonResult, false
								}
								var asBoolError *Value
								comparisonBool, asBoolError = p.QuickGetBool(context, comparisonResult)
								if asBoolError != nil {
									return asBoolError, false
								}
								if !comparisonBool {
									return p.GetFalse(), true
								}
							}
							return p.GetTrue(), true
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
							if !right.IsTypeById(TupleId) {
								return p.GetTrue(), true
							}
							if self.GetLength() != right.GetLength() {
								return p.GetTrue(), true
							}
							var rightNotEquals *Value
							var comparisonResult *Value
							var success bool
							var callError *Value
							var comparisonBool bool

							for i := 0; i < self.GetLength(); i++ {
								leftNotEquals, getError := self.GetContent()[i].Get(NotEquals)
								if getError != nil {
									rightNotEquals, getError = right.GetContent()[i].Get(RightNotEquals)
									if getError != nil {
										return p.NewObjectWithNameNotFoundError(context, right.GetContent()[i].GetClass(p), RightNotEquals), false
									}
									comparisonResult, success = p.CallFunction(context, rightNotEquals, self.GetContent()[i])
								} else {
									comparisonResult, success = p.CallFunction(context, leftNotEquals, right.GetContent()[i])
								}
								if !success {
									return comparisonResult, false
								}
								comparisonBool, callError = p.QuickGetBool(context, comparisonResult)
								if callError != nil {
									return callError, false
								}
								if !comparisonBool {
									return p.GetFalse(), true
								}
							}
							return p.GetTrue(), true
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
							if !left.IsTypeById(TupleId) {
								return p.GetFalse(), true
							}
							if self.GetLength() != left.GetLength() {
								return p.GetFalse(), true
							}
							var rightEquals *Value
							var comparisonResult *Value
							var success bool

							for i := 0; i < self.GetLength(); i++ {
								leftEquals, getError := left.GetContent()[i].Get(NotEquals)
								if getError != nil {
									rightEquals, getError = self.GetContent()[i].Get(RightNotEquals)
									if getError != nil {
										return p.NewObjectWithNameNotFoundError(context, self.GetContent()[i].GetClass(p), RightNotEquals), false
									}
									comparisonResult, success = p.CallFunction(context, rightEquals, left.GetContent()[i])
								} else {
									comparisonResult, success = p.CallFunction(context, leftEquals, self.GetContent()[i])
								}
								if !success {
									return comparisonResult, false
								}
								comparisonBool, callError := p.QuickGetBool(context, comparisonResult)
								if callError != nil {
									return callError, false
								}
								if !comparisonBool {
									return p.GetFalse(), true
								}
							}
							return p.GetTrue(), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Contains,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							value := arguments[0]
							valueRightEquals, getError := value.Get(RightEquals)
							if getError != nil {
								return p.NewObjectWithNameNotFoundError(context, value.GetClass(p), RightEquals), false
							}
							for _, tupleValue := range self.GetContent() {
								callResult, success := p.CallFunction(context, valueRightEquals, tupleValue)
								if !success {
									return callResult, false
								}
								var boolValue *Value
								if callResult.IsTypeById(BoolId) {
									boolValue = callResult
								} else {
									var boolValueToBool *Value
									boolValueToBool, getError = callResult.Get(ToBool)
									if getError != nil {
										return p.NewObjectWithNameNotFoundError(context, callResult.GetClass(p), ToBool), false
									}
									callResult, success = p.CallFunction(context, boolValueToBool)
									if !success {
										return callResult, false
									}
									if !callResult.IsTypeById(BoolId) {
										return p.NewInvalidTypeError(context, callResult.TypeName(), BoolName), false
									}
									boolValue = callResult
								}
								if boolValue.GetBool() {
									return p.GetTrue(), true
								}
							}
							return p.GetFalse(), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightContains,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							value := arguments[0]
							valueRightEquals, getError := value.Get(Equals)
							if getError != nil {
								return p.NewObjectWithNameNotFoundError(context, value.GetClass(p), Equals), false
							}
							for _, tupleValue := range self.GetContent() {
								callResult, success := p.CallFunction(context, valueRightEquals, tupleValue)
								if !success {
									return callResult, false
								}
								var boolValue *Value
								if callResult.IsTypeById(BoolId) {
									boolValue = callResult
								} else {
									var boolValueToBool *Value
									boolValueToBool, getError = callResult.Get(ToBool)
									if getError != nil {
										return p.NewObjectWithNameNotFoundError(context, callResult.GetClass(p), ToBool), false
									}
									callResult, success = p.CallFunction(context, boolValueToBool)
									if !success {
										return callResult, false
									}
									if !callResult.IsTypeById(BoolId) {
										return p.NewInvalidTypeError(context, callResult.TypeName(), BoolName), false
									}
									boolValue = callResult
								}
								if boolValue.GetBool() {
									return p.GetTrue(), true
								}
							}
							return p.GetFalse(), true
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
							tupleHash := XXPrime5 ^ p.Seed()
							for _, contentObject := range self.GetContent() {
								objectHashFunc, getError := contentObject.Get(Hash)
								if getError != nil {
									return p.NewObjectWithNameNotFoundError(context, contentObject.GetClass(p), Hash), false
								}
								objectHash, success := p.CallFunction(context, objectHashFunc)
								if !success {
									return objectHash, false
								}
								if !objectHash.IsTypeById(IntegerId) {
									return p.NewInvalidTypeError(context, objectHash.TypeName(), IntegerName), false
								}
								tupleHash += uint64(objectHash.GetInteger()) * XXPrime2
								tupleHash = (tupleHash << 31) | (tupleHash >> 33)
								tupleHash *= XXPrime1
								tupleHash &= (1 << 64) - 1
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(), int64(tupleHash)), true
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

							var copiedObjects []*Value
							for _, contentObject := range self.GetContent() {
								objectCopy, getError := contentObject.Get(Copy)
								if getError != nil {
									return p.NewObjectWithNameNotFoundError(context, contentObject.GetClass(p), Copy), false
								}
								copiedObject, success := p.CallFunction(context, objectCopy)
								if !success {
									return copiedObject, false
								}
								copiedObjects = append(copiedObjects, copiedObject)
							}
							return p.NewTuple(context, false, context.PeekSymbolTable(), copiedObjects), true
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
								return self.GetContent()[index], true
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
								return p.NewTuple(context, false, context.PeekSymbolTable(), self.GetContent()[startIndex:targetIndex]), true
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
							iterator.SetInteger(0)
							iterator.SetContent(self.GetContent())
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
											value := funcSelf.GetContent()[int(funcSelf.GetInteger())]
											funcSelf.SetInteger(funcSelf.GetInteger() + 1)
											return value, true
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
		object.SetOnDemandSymbol(ToString,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							result := "("
							var objectString *Value
							var success bool
							for index, contentObject := range self.GetContent() {
								if index != 0 {
									result += ", "
								}
								objectToString, getError := contentObject.Get(ToString)
								if getError != nil {
									return p.NewObjectWithNameNotFoundError(context, contentObject.GetClass(p), ToString), false
								}
								objectString, success = p.CallFunction(context, objectToString)
								if !success {
									return objectString, false
								}
								result += objectString.GetString()
							}
							return p.NewString(context, false, context.PeekSymbolTable(), result+")"), true
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
							return p.NewArray(context, false, context.PeekSymbolTable(), append([]*Value{}, self.GetContent()...)), true
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
							return p.NewTuple(context, false, context.PeekSymbolTable(), append([]*Value{}, self.GetContent()...)), true
						},
					),
				)
			},
		)
		return nil
	}
}
