package vm

import "github.com/shoriwe/gruby/pkg/errors"

type Tuple struct {
	*Object
}

func (p *Plasma) NewTuple(parentSymbols *SymbolTable, content []IObject) *Tuple {
	tuple := &Tuple{
		Object: p.NewObject(TupleName, nil, parentSymbols),
	}
	tuple.SetContent(content)
	tuple.SetLength(len(content))
	p.TupleInitialize(tuple)
	return tuple
}

func (p *Plasma) TupleInitialize(object IObject) *errors.Error {
	object.Set(Equals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					if _, ok := right.(*Tuple); !ok {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					if self.GetLength() != right.GetLength() {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					var rightEquals IObject
					var comparisonResult IObject
					var callError *errors.Error
					var comparisonResultToBool IObject
					var comparisonBool IObject

					for i := 0; i < self.GetLength(); i++ {
						leftEquals, getError := self.GetContent()[i].Get(Equals)
						if getError != nil {
							rightEquals, getError = right.GetContent()[i].Get(RightEquals)
							if getError != nil {
								return nil, getError
							}
							if _, ok := rightEquals.(*Function); !ok {
								return nil, errors.NewTypeError(rightEquals.TypeName(), FunctionName)
							}
							comparisonResult, callError = p.CallFunction(rightEquals.(*Function), p.PeekSymbolTable(), self.GetContent()[i])
						} else {
							comparisonResult, callError = p.CallFunction(leftEquals.(*Function), p.PeekSymbolTable(), right.GetContent()[i])
						}
						if callError != nil {
							return nil, callError
						}
						comparisonResultToBool, getError = comparisonResult.Get(ToBool)
						if getError != nil {
							return nil, getError
						}
						if _, ok := comparisonResultToBool.(*Function); !ok {
							return nil, errors.NewTypeError(comparisonResultToBool.TypeName(), FunctionName)
						}
						comparisonBool, callError = p.CallFunction(comparisonResultToBool.(*Function), p.PeekSymbolTable())
						if callError != nil {
							return nil, callError
						}
						if _, ok := comparisonBool.(*Bool); !ok {
							return nil, errors.NewTypeError(comparisonBool.TypeName(), BoolName)
						}
						if !comparisonBool.GetBool() {
							return p.NewBool(p.PeekSymbolTable(), false), nil
						}
					}
					return p.NewBool(p.PeekSymbolTable(), true), nil
				},
			),
		),
	)
	object.Set(RightEquals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					if _, ok := left.(*Tuple); !ok {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					if self.GetLength() != left.GetLength() {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					var rightEquals IObject
					var comparisonResult IObject
					var callError *errors.Error
					var comparisonResultToBool IObject
					var comparisonBool IObject

					for i := 0; i < self.GetLength(); i++ {
						leftEquals, getError := left.GetContent()[i].Get(Equals)
						if getError != nil {
							rightEquals, getError = self.GetContent()[i].Get(RightEquals)
							if getError != nil {
								return nil, getError
							}
							if _, ok := rightEquals.(*Function); !ok {
								return nil, errors.NewTypeError(rightEquals.TypeName(), FunctionName)
							}
							comparisonResult, callError = p.CallFunction(rightEquals.(*Function), p.PeekSymbolTable(), left.GetContent()[i])
						} else {
							comparisonResult, callError = p.CallFunction(leftEquals.(*Function), p.PeekSymbolTable(), self.GetContent()[i])
						}
						if callError != nil {
							return nil, callError
						}
						comparisonResultToBool, getError = comparisonResult.Get(ToBool)
						if getError != nil {
							return nil, getError
						}
						if _, ok := comparisonResultToBool.(*Function); !ok {
							return nil, errors.NewTypeError(comparisonResultToBool.TypeName(), FunctionName)
						}
						comparisonBool, callError = p.CallFunction(comparisonResultToBool.(*Function), p.PeekSymbolTable())
						if callError != nil {
							return nil, callError
						}
						if _, ok := comparisonBool.(*Bool); !ok {
							return nil, errors.NewTypeError(comparisonBool.TypeName(), BoolName)
						}
						if !comparisonBool.GetBool() {
							return p.NewBool(p.PeekSymbolTable(), false), nil
						}
					}
					return p.NewBool(p.PeekSymbolTable(), true), nil
				},
			),
		),
	)
	object.Set(NotEquals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {

					right := arguments[0]
					if _, ok := right.(*Tuple); !ok {
						return p.NewBool(p.PeekSymbolTable(), true), nil
					}
					if self.GetLength() != right.GetLength() {
						return p.NewBool(p.PeekSymbolTable(), true), nil
					}
					var rightNotEquals IObject
					var comparisonResult IObject
					var callError *errors.Error
					var comparisonResultToBool IObject
					var comparisonBool IObject

					for i := 0; i < self.GetLength(); i++ {
						leftNotEquals, getError := self.GetContent()[i].Get(NotEquals)
						if getError != nil {
							rightNotEquals, getError = right.GetContent()[i].Get(RightNotEquals)
							if getError != nil {
								return nil, getError
							}
							if _, ok := rightNotEquals.(*Function); !ok {
								return nil, errors.NewTypeError(rightNotEquals.TypeName(), FunctionName)
							}
							comparisonResult, callError = p.CallFunction(rightNotEquals.(*Function), p.PeekSymbolTable(), self.GetContent()[i])
						} else {
							comparisonResult, callError = p.CallFunction(leftNotEquals.(*Function), p.PeekSymbolTable(), right.GetContent()[i])
						}
						if callError != nil {
							return nil, callError
						}
						comparisonResultToBool, getError = comparisonResult.Get(ToBool)
						if getError != nil {
							return nil, getError
						}
						if _, ok := comparisonResultToBool.(*Function); !ok {
							return nil, errors.NewTypeError(comparisonResultToBool.TypeName(), FunctionName)
						}
						comparisonBool, callError = p.CallFunction(comparisonResultToBool.(*Function), p.PeekSymbolTable())
						if callError != nil {
							return nil, callError
						}
						if _, ok := comparisonBool.(*Bool); !ok {
							return nil, errors.NewTypeError(comparisonBool.TypeName(), BoolName)
						}
						if !comparisonBool.GetBool() {
							return p.NewBool(p.PeekSymbolTable(), false), nil
						}
					}
					return p.NewBool(p.PeekSymbolTable(), true), nil
				},
			),
		),
	)
	object.Set(RightNotEquals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {

					left := arguments[0]
					if _, ok := left.(*Tuple); !ok {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					if self.GetLength() != left.GetLength() {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					var rightEquals IObject
					var comparisonResult IObject
					var callError *errors.Error
					var comparisonResultToBool IObject
					var comparisonBool IObject

					for i := 0; i < self.GetLength(); i++ {
						leftEquals, getError := left.GetContent()[i].Get(NotEquals)
						if getError != nil {
							rightEquals, getError = self.GetContent()[i].Get(RightNotEquals)
							if getError != nil {
								return nil, getError
							}
							if _, ok := rightEquals.(*Function); !ok {
								return nil, errors.NewTypeError(rightEquals.TypeName(), FunctionName)
							}
							comparisonResult, callError = p.CallFunction(rightEquals.(*Function), p.PeekSymbolTable(), left.GetContent()[i])
						} else {
							comparisonResult, callError = p.CallFunction(leftEquals.(*Function), p.PeekSymbolTable(), self.GetContent()[i])
						}
						if callError != nil {
							return nil, callError
						}
						comparisonResultToBool, getError = comparisonResult.Get(ToBool)
						if getError != nil {
							return nil, getError
						}
						if _, ok := comparisonResultToBool.(*Function); !ok {
							return nil, errors.NewTypeError(comparisonResultToBool.TypeName(), FunctionName)
						}
						comparisonBool, callError = p.CallFunction(comparisonResultToBool.(*Function), p.PeekSymbolTable())
						if callError != nil {
							return nil, callError
						}
						if _, ok := comparisonBool.(*Bool); !ok {
							return nil, errors.NewTypeError(comparisonBool.TypeName(), BoolName)
						}
						if !comparisonBool.GetBool() {
							return p.NewBool(p.PeekSymbolTable(), false), nil
						}
					}
					return p.NewBool(p.PeekSymbolTable(), true), nil
				},
			),
		),
	)
	object.Set(Hash,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					tupleHash := XXPrime5 ^ p.Seed()
					for _, contentObject := range self.GetContent() {
						objectHashFunc, getError := contentObject.Get(Hash)
						if getError != nil {
							return nil, getError
						}
						if _, ok := objectHashFunc.(*Function); !ok {
							return nil, errors.NewTypeError(objectHashFunc.TypeName(), FunctionName)
						}
						objectHash, callError := p.CallFunction(objectHashFunc.(*Function), self.SymbolTable())
						if callError != nil {
							return nil, callError
						}
						if _, ok := objectHash.(*Integer); !ok {
							return nil, errors.NewTypeError(objectHash.TypeName(), IntegerName)
						}
						tupleHash += uint64(objectHash.GetInteger64()) * XXPrime2
						tupleHash = (tupleHash << 31) | (tupleHash >> 33)
						tupleHash *= XXPrime1
						tupleHash &= (1 << 64) - 1
					}
					finalTupleHash := int64(tupleHash)
					if finalTupleHash < 0 {
						finalTupleHash = -finalTupleHash
					}
					return p.NewInteger(p.PeekSymbolTable(), finalTupleHash), nil
				},
			),
		),
	)
	object.Set(Copy,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {

					var copiedObjects []IObject
					for _, contentObject := range self.GetContent() {
						objectCopy, getError := contentObject.Get(Copy)
						if getError != nil {
							return nil, getError
						}
						if _, ok := objectCopy.(*Function); !ok {
							return nil, errors.NewTypeError(objectCopy.TypeName(), FunctionName)
						}
						copiedObject, copyError := p.CallFunction(objectCopy.(*Function), p.PeekSymbolTable())
						if copyError != nil {
							return nil, copyError
						}
						copiedObjects = append(copiedObjects, copiedObject)
					}
					return p.NewTuple(p.PeekSymbolTable(), copiedObjects), nil
				},
			),
		),
	)
	object.Set(Index,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					indexObject := arguments[0]
					var ok bool
					if _, ok = indexObject.(*Integer); ok {
						index, calcError := CalcIndex(indexObject, self.GetLength())
						if calcError != nil {
							return nil, calcError
						}
						return self.GetContent()[index], nil
					} else if _, ok = indexObject.(*Tuple); ok {
						if len(indexObject.GetContent()) != 2 {
							return nil, errors.NewInvalidNumberOfArguments(len(indexObject.GetContent()), 2)
						}
						startIndex, calcError := CalcIndex(indexObject.GetContent()[0], self.GetLength())
						if calcError != nil {
							return nil, calcError
						}
						targetIndex, calcError := CalcIndex(indexObject.GetContent()[1], self.GetLength())
						if calcError != nil {
							return nil, calcError
						}
						return p.NewTuple(p.PeekSymbolTable(), self.GetContent()[startIndex:targetIndex]), nil
					} else {
						return nil, errors.NewTypeError(indexObject.TypeName(), IntegerName, TupleName)
					}
				},
			),
		),
	)
	object.Set(Iter,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {

					iterator := p.NewIterator(p.PeekSymbolTable())
					iterator.SetInteger64(0)
					iterator.SetContent(self.GetContent())
					iterator.SetLength(self.GetLength())
					iterator.Set(HasNext,
						p.NewFunction(iterator.SymbolTable(),
							NewBuiltInClassFunction(iterator,
								0,
								func(funcSelf IObject, _ ...IObject) (IObject, *errors.Error) {
									return p.NewBool(p.PeekSymbolTable(), int(funcSelf.GetInteger64()) < funcSelf.GetLength()), nil
								},
							),
						),
					)
					iterator.Set(Next,
						p.NewFunction(iterator.SymbolTable(),
							NewBuiltInClassFunction(iterator,
								0,
								func(funcSelf IObject, _ ...IObject) (IObject, *errors.Error) {
									value := funcSelf.GetContent()[int(funcSelf.GetInteger64())]
									funcSelf.SetInteger64(funcSelf.GetInteger64() + 1)
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
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					result := "("
					var objectString IObject
					var callError *errors.Error
					for index, contentObject := range self.GetContent() {
						if index != 0 {
							result += ", "
						}
						objectToString, getError := contentObject.Get(ToString)
						if getError != nil {
							return nil, getError
						}
						if _, ok := objectToString.(*Function); !ok {
							return nil, errors.NewTypeError(objectToString.TypeName(), FunctionName)
						}
						objectString, callError = p.CallFunction(objectToString.(*Function), p.PeekSymbolTable())
						if callError != nil {
							return nil, callError
						}
						result += objectString.GetString()
					}
					return p.NewString(p.PeekSymbolTable(), result+")"), nil
				},
			),
		),
	)
	object.Set(ToBool,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewBool(p.PeekSymbolTable(), self.GetLength() != 0), nil
				},
			),
		),
	)
	object.Set(ToArray,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewArray(p.PeekSymbolTable(), append([]IObject{}, self.GetContent()...)), nil
				},
			),
		),
	)
	object.Set(ToTuple,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewTuple(p.PeekSymbolTable(), append([]IObject{}, self.GetContent()...)), nil
				},
			),
		),
	)
	// Tuple Special functions
	object.Set(Contains,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					value := arguments[0]
					valueRightEquals, getError := value.Get(RightEquals)
					if getError != nil {
						return nil, getError
					}
					if _, ok := valueRightEquals.(*Function); !ok {
						return nil, errors.NewTypeError(valueRightEquals.TypeName(), FunctionName)
					}
					for _, tupleValue := range self.GetContent() {
						callResult, callError := p.CallFunction(valueRightEquals.(*Function), value.SymbolTable(), tupleValue)
						if callError != nil {
							return nil, callError
						}
						var boolValue IObject
						if _, ok := callResult.(*Bool); ok {
							boolValue = callResult
						} else {
							var boolValueToBool IObject
							boolValueToBool, getError = callResult.Get(ToBool)
							if getError != nil {
								return nil, getError
							}
							if _, ok = boolValueToBool.(*Function); !ok {
								return nil, errors.NewTypeError(boolValueToBool.TypeName(), FunctionName)
							}
							callResult, callError = p.CallFunction(boolValueToBool.(*Function), callResult.SymbolTable())
							if callError != nil {
								return nil, callError
							}
							if _, ok = callResult.(*Bool); !ok {
								return nil, errors.NewTypeError(callResult.TypeName(), BoolName)
							}
							boolValue = callResult
						}
						if boolValue.(*Bool).GetBool() {
							return p.NewBool(p.PeekSymbolTable(), true), nil
						}
					}
					return p.NewBool(p.PeekSymbolTable(), false), nil
				},
			),
		),
	)
	return nil
}
