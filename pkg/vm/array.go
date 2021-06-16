package vm

import (
	"github.com/shoriwe/gplasma/pkg/tools"
)

type Array struct {
	*Object
}

func (p *Plasma) NewArray(isBuiltIn bool, parentSymbols *SymbolTable, content []Value) *Array {
	array := &Array{
		Object: p.NewObject(false, ArrayName, nil, parentSymbols),
	}
	array.SetContent(content)
	array.SetLength(len(content))
	p.ArrayInitialize(isBuiltIn)(array)
	array.Set(Self, array)
	return array
}

func (p *Plasma) ArrayInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(object Value) *Object {
		object.SymbolTable().Set(Mul,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							content, repetitionError := p.Repeat(self.GetContent(), int(right.GetInteger64()))
							if repetitionError != nil {
								return nil, repetitionError
							}
							return p.NewArray(false, p.PeekSymbolTable(), content), nil
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
							content, repetitionError := p.Repeat(self.GetContent(), int(left.GetInteger64()))
							if repetitionError != nil {
								return nil, repetitionError
							}
							return p.NewArray(false, p.PeekSymbolTable(), content), nil
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
						if _, ok := right.(*Array); !ok {
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
									return nil, p.NewObjectWithNameNotFoundError(RightEquals)
								}
								if _, ok := rightEquals.(*Function); !ok {
									return nil, p.NewInvalidTypeError(rightEquals.TypeName(), FunctionName)
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
								return nil, p.NewObjectWithNameNotFoundError(ToBool)
							}
							if _, ok := comparisonResultToBool.(*Function); !ok {
								return nil, p.NewInvalidTypeError(comparisonResultToBool.TypeName(), FunctionName)
							}
							comparisonBool, callError = p.CallFunction(comparisonResultToBool.(*Function), p.PeekSymbolTable())
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
						if _, ok := left.(*Array); !ok {
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
									return nil, p.NewObjectWithNameNotFoundError(RightEquals)
								}
								if _, ok := rightEquals.(*Function); !ok {
									return nil, p.NewInvalidTypeError(rightEquals.TypeName(), FunctionName)
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
								return nil, p.NewObjectWithNameNotFoundError(ToBool)
							}
							if _, ok := comparisonResultToBool.(*Function); !ok {
								return nil, p.NewInvalidTypeError(comparisonResultToBool.TypeName(), FunctionName)
							}
							comparisonBool, callError = p.CallFunction(comparisonResultToBool.(*Function), p.PeekSymbolTable())
							if callError != nil {
								return nil, callError
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
						if _, ok := right.(*Array); !ok {
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
									return nil, p.NewObjectWithNameNotFoundError(RightNotEquals)
								}
								if _, ok := rightNotEquals.(*Function); !ok {
									return nil, p.NewInvalidTypeError(rightNotEquals.TypeName(), FunctionName)
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
								return nil, p.NewObjectWithNameNotFoundError(ToBool)
							}
							if _, ok := comparisonResultToBool.(*Function); !ok {
								return nil, p.NewInvalidTypeError(comparisonResultToBool.TypeName(), FunctionName)
							}
							comparisonBool, callError = p.CallFunction(comparisonResultToBool.(*Function), p.PeekSymbolTable())
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
						if _, ok := left.(*Array); !ok {
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
									return nil, p.NewObjectWithNameNotFoundError(RightNotEquals)
								}
								if _, ok := rightEquals.(*Function); !ok {
									return nil, p.NewInvalidTypeError(rightEquals.TypeName(), FunctionName)
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
								return nil, p.NewObjectWithNameNotFoundError(ToBool)
							}
							if _, ok := comparisonResultToBool.(*Function); !ok {
								return nil, p.NewInvalidTypeError(comparisonResultToBool.TypeName(), FunctionName)
							}
							comparisonBool, callError = p.CallFunction(comparisonResultToBool.(*Function), p.PeekSymbolTable())
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
		object.Set(Hash,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(_ Value, _ ...Value) (Value, *Object) {
						return nil, p.NewUnhashableTypeError(object.GetClass())
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
								return nil, p.NewObjectWithNameNotFoundError(Copy)
							}
							if _, ok := objectCopy.(*Function); !ok {
								return nil, p.NewInvalidTypeError(objectCopy.TypeName(), FunctionName)
							}
							copiedObject, copyError := p.CallFunction(objectCopy.(*Function), p.PeekSymbolTable())
							if copyError != nil {
								return nil, copyError
							}
							copiedObjects = append(copiedObjects, copiedObject)
						}
						return p.NewArray(false, p.PeekSymbolTable(), copiedObjects), nil
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
							index, calcError := tools.CalcIndex(indexObject.GetInteger64(), self.GetLength())
							if calcError != nil {
								return nil, p.NewIndexOutOfRange(self.GetLength(), indexObject.GetInteger64())
							}
							return self.GetContent()[index], nil
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
							return p.NewArray(false, p.PeekSymbolTable(), self.GetContent()[startIndex:targetIndex]), nil
						} else {
							return nil, p.NewInvalidTypeError(indexObject.TypeName(), IntegerName, TupleName)
						}
					},
				),
			),
		)
		object.Set(Assign,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 2,
					func(self Value, arguments ...Value) (Value, *Object) {
						index, calcError := tools.CalcIndex(arguments[0].GetInteger64(), self.GetLength())
						if calcError != nil {
							return nil, p.NewIndexOutOfRange(self.GetLength(), arguments[0].GetInteger64())
						}
						self.GetContent()[index] = arguments[1]
						return p.NewNone(), nil
					},
				),
			),
		)
		object.Set(Iter,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {

						iterator := p.NewIterator(false, p.PeekSymbolTable())
						iterator.SetInteger64(0)
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
										return p.NewBool(false, p.PeekSymbolTable(), int(funcSelf.GetInteger64()) < funcSelf.GetLength()), nil
									},
								),
							),
						)
						iterator.Set(Next,
							p.NewFunction(isBuiltIn, iterator.SymbolTable(),
								NewBuiltInClassFunction(iterator,
									0,
									func(funcSelf Value, _ ...Value) (Value, *Object) {
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
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
								return nil, p.NewObjectWithNameNotFoundError(ToString)
							}
							if _, ok := objectToString.(*Function); !ok {
								return nil, p.NewInvalidTypeError(objectToString.TypeName(), FunctionName)
							}
							objectString, callError = p.CallFunction(objectToString.(*Function), p.PeekSymbolTable())
							if callError != nil {
								return nil, callError
							}
							result += objectString.GetString()
						}
						return p.NewString(false, p.PeekSymbolTable(), result+"]"), nil
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
