package vm

import (
	"github.com/shoriwe/gplasma/pkg/tools"
	"math/big"
)

type Tuple struct {
	*Object
}

func (p *Plasma) NewTuple(isBuiltIn bool, parentSymbols *SymbolTable, content []Value) *Tuple {
	tuple := &Tuple{
		Object: p.NewObject(false, TupleName, nil, parentSymbols),
	}
	tuple.SetContent(content)
	tuple.SetLength(len(content))
	p.TupleInitialize(isBuiltIn)(tuple)
	tuple.Set(Self, tuple)
	return tuple
}

func (p *Plasma) TupleInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(object Value) *Object {
		object.SymbolTable().Set(Mul,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							content, repetitionError := p.Repeat(self.GetContent(), right.GetInteger())
							if repetitionError != nil {
								return nil, repetitionError
							}
							return p.NewTuple(false, p.PeekSymbolTable(), content), nil
						default:
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName)
						}
					},
				),
			),
		)
		object.SymbolTable().Set(RightMul,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						switch left.(type) {
						case *Integer:
							content, repetitionError := p.Repeat(self.GetContent(), left.GetInteger())
							if repetitionError != nil {
								return nil, repetitionError
							}
							return p.NewTuple(false, p.PeekSymbolTable(), content), nil
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName)
						}
					},
				),
			),
		)
		object.Set(Equals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						if _, ok := right.(*Tuple); !ok {
							return p.NewBool(false, p.PeekSymbolTable(), false), nil
						}
						if self.GetLength() != right.GetLength() {
							return p.NewBool(false, p.PeekSymbolTable(), false), nil
						}
						var rightEquals Value
						var comparisonResult Value
						var callError *Object
						var comparisonResultToBool Value
						var comparisonBool Value

						for i := 0; i < self.GetLength(); i++ {
							leftEquals, getError := self.GetContent()[i].Get(Equals)
							if getError != nil {
								rightEquals, getError = right.GetContent()[i].Get(RightEquals)
								if getError != nil {
									return nil, p.NewObjectWithNameNotFoundError(right.GetContent()[i].GetClass(p), RightEquals)
								}
								comparisonResult, callError = p.CallFunction(rightEquals, p.PeekSymbolTable(), self.GetContent()[i])
							} else {
								comparisonResult, callError = p.CallFunction(leftEquals, p.PeekSymbolTable(), right.GetContent()[i])
							}
							if callError != nil {
								return nil, callError
							}
							comparisonResultToBool, getError = comparisonResult.Get(ToBool)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(comparisonResult.GetClass(p), ToBool)
							}
							comparisonBool, callError = p.CallFunction(comparisonResultToBool, p.PeekSymbolTable())
							if callError != nil {
								return nil, callError
							}
							if _, ok := comparisonBool.(*Bool); !ok {
								return nil, p.NewInvalidTypeError(comparisonBool.TypeName(), BoolName)
							}
							if !comparisonBool.GetBool() {
								return p.NewBool(false, p.PeekSymbolTable(), false), nil
							}
						}
						return p.NewBool(false, p.PeekSymbolTable(), true), nil
					},
				),
			),
		)
		object.Set(RightEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						if _, ok := left.(*Tuple); !ok {
							return p.NewBool(false, p.PeekSymbolTable(), false), nil
						}
						if self.GetLength() != left.GetLength() {
							return p.NewBool(false, p.PeekSymbolTable(), false), nil
						}
						var rightEquals Value
						var comparisonResult Value
						var callError *Object
						var comparisonResultToBool Value
						var comparisonBool Value

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
							comparisonResultToBool, getError = comparisonResult.Get(ToBool)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(comparisonResult.GetClass(p), ToBool)
							}
							comparisonBool, callError = p.CallFunction(comparisonResultToBool, p.PeekSymbolTable())
							if callError != nil {
								return nil, callError
							}
							if _, ok := comparisonBool.(*Bool); !ok {
								return nil, p.NewInvalidTypeError(comparisonBool.TypeName(), BoolName)
							}
							if !comparisonBool.GetBool() {
								return p.NewBool(false, p.PeekSymbolTable(), false), nil
							}
						}
						return p.NewBool(false, p.PeekSymbolTable(), true), nil
					},
				),
			),
		)
		object.Set(NotEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {

						right := arguments[0]
						if _, ok := right.(*Tuple); !ok {
							return p.NewBool(false, p.PeekSymbolTable(), true), nil
						}
						if self.GetLength() != right.GetLength() {
							return p.NewBool(false, p.PeekSymbolTable(), true), nil
						}
						var rightNotEquals Value
						var comparisonResult Value
						var callError *Object
						var comparisonResultToBool Value
						var comparisonBool Value

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
							comparisonResultToBool, getError = comparisonResult.Get(ToBool)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(comparisonResult.GetClass(p), ToBool)
							}
							comparisonBool, callError = p.CallFunction(comparisonResultToBool, p.PeekSymbolTable())
							if callError != nil {
								return nil, callError
							}
							if _, ok := comparisonBool.(*Bool); !ok {
								return nil, p.NewInvalidTypeError(comparisonBool.TypeName(), BoolName)
							}
							if !comparisonBool.GetBool() {
								return p.NewBool(false, p.PeekSymbolTable(), false), nil
							}
						}
						return p.NewBool(false, p.PeekSymbolTable(), true), nil
					},
				),
			),
		)
		object.Set(RightNotEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {

						left := arguments[0]
						if _, ok := left.(*Tuple); !ok {
							return p.NewBool(false, p.PeekSymbolTable(), false), nil
						}
						if self.GetLength() != left.GetLength() {
							return p.NewBool(false, p.PeekSymbolTable(), false), nil
						}
						var rightEquals Value
						var comparisonResult Value
						var callError *Object
						var comparisonResultToBool Value
						var comparisonBool Value

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
							comparisonResultToBool, getError = comparisonResult.Get(ToBool)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(comparisonResult.GetClass(p), ToBool)
							}
							comparisonBool, callError = p.CallFunction(comparisonResultToBool, p.PeekSymbolTable())
							if callError != nil {
								return nil, callError
							}
							if _, ok := comparisonBool.(*Bool); !ok {
								return nil, p.NewInvalidTypeError(comparisonBool.TypeName(), BoolName)
							}
							if !comparisonBool.GetBool() {
								return p.NewBool(false, p.PeekSymbolTable(), false), nil
							}
						}
						return p.NewBool(false, p.PeekSymbolTable(), true), nil
					},
				),
			),
		)
		object.Set(Contains,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
								return p.NewBool(false, p.PeekSymbolTable(), true), nil
							}
						}
						return p.NewBool(false, p.PeekSymbolTable(), false), nil
					},
				),
			),
		)
		object.Set(RightContains,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						value := arguments[0]
						valueRightEquals, getError := value.Get(Equals)
						if getError != nil {
							return nil, p.NewObjectWithNameNotFoundError(value.GetClass(p), Equals)
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
								return p.NewBool(false, p.PeekSymbolTable(), true), nil
							}
						}
						return p.NewBool(false, p.PeekSymbolTable(), false), nil
					},
				),
			),
		)
		object.Set(Hash,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						tupleHash := XXPrime5 ^ p.Seed()
						for _, contentObject := range self.GetContent() {
							objectHashFunc, getError := contentObject.Get(Hash)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(contentObject.GetClass(p), Hash)
							}
							objectHash, callError := p.CallFunction(objectHashFunc, self.SymbolTable())
							if callError != nil {
								return nil, callError
							}
							if _, ok := objectHash.(*Integer); !ok {
								return nil, p.NewInvalidTypeError(objectHash.TypeName(), IntegerName)
							}
							tupleHash += objectHash.GetInteger().Uint64() * XXPrime2
							tupleHash = (tupleHash << 31) | (tupleHash >> 33)
							tupleHash *= XXPrime1
							tupleHash &= (1 << 64) - 1
						}
						return p.NewInteger(false, p.PeekSymbolTable(), new(big.Int).SetUint64(tupleHash)), nil
					},
				),
			),
		)
		object.Set(Copy,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
						return p.NewTuple(false, p.PeekSymbolTable(), copiedObjects), nil
					},
				),
			),
		)
		object.Set(Index,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						indexObject := arguments[0]
						var ok bool
						if _, ok = indexObject.(*Integer); ok {
							index, calcError := tools.CalcIndex(indexObject.GetInteger().Int64(), self.GetLength())
							if calcError != nil {
								return nil, p.NewIndexOutOfRange(self.GetLength(), indexObject.GetInteger().Int64())
							}
							return self.GetContent()[index], nil
						} else if _, ok = indexObject.(*Tuple); ok {
							if len(indexObject.GetContent()) != 2 {
								return nil, p.NewInvalidNumberOfArgumentsError(len(indexObject.GetContent()), 2)
							}
							startIndex, calcError := tools.CalcIndex(indexObject.GetContent()[0].GetInteger().Int64(), self.GetLength())
							if calcError != nil {
								return nil, p.NewIndexOutOfRange(self.GetLength(), indexObject.GetContent()[0].GetInteger().Int64())
							}
							var targetIndex int
							targetIndex, calcError = tools.CalcIndex(indexObject.GetContent()[1].GetInteger().Int64(), self.GetLength())
							if calcError != nil {
								return nil, p.NewIndexOutOfRange(self.GetLength(), indexObject.GetContent()[1].GetInteger().Int64())
							}
							return p.NewTuple(false, p.PeekSymbolTable(), self.GetContent()[startIndex:targetIndex]), nil
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
						iterator.SetInteger(big.NewInt(0))
						iterator.SetContent(self.GetContent())
						iterator.SetLength(self.GetLength())
						iterator.Set(HasNext,
							p.NewFunction(isBuiltIn, iterator.SymbolTable(),
								NewBuiltInClassFunction(iterator,
									0,
									func(funcSelf Value, _ ...Value) (Value, *Object) {
										return p.NewBool(false, p.PeekSymbolTable(), funcSelf.GetInteger().Cmp(big.NewInt(int64(funcSelf.GetLength()))) == -1), nil
									},
								),
							),
						)
						iterator.Set(Next,
							p.NewFunction(isBuiltIn, iterator.SymbolTable(),
								NewBuiltInClassFunction(iterator,
									0,
									func(funcSelf Value, _ ...Value) (Value, *Object) {
										value := funcSelf.GetContent()[int(funcSelf.GetInteger().Int64())]
										funcSelf.SetInteger(new(big.Int).Add(funcSelf.GetInteger(), big.NewInt(1)))
										return value, nil
									},
								),
							),
						)
						return iterator, nil
					},
				),
			),
		)
		object.Set(ToString,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
								return nil, p.NewObjectWithNameNotFoundError(contentObject.GetClass(p), ToString)
							}
							objectString, callError = p.CallFunction(objectToString, p.PeekSymbolTable())
							if callError != nil {
								return nil, callError
							}
							result += objectString.GetString()
						}
						return p.NewString(false, p.PeekSymbolTable(), result+")"), nil
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
						return p.NewArray(false, p.PeekSymbolTable(), append([]Value{}, self.GetContent()...)), nil
					},
				),
			),
		)
		object.Set(ToTuple,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewTuple(false, p.PeekSymbolTable(), append([]Value{}, self.GetContent()...)), nil
					},
				),
			),
		)
		return nil
	}
}
