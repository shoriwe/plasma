package vm

import (
	"fmt"
)

/*
	Type         - (Done)
	Function     - (Done)
	Object       - (Done)
	Bool         - (Done)
	Bytes        - (Done)
	String       - (Done)
	HashTable    - (Done)
	Integer      - (Done)
	Array        - (Done)
	Tuple        - (Done)
	Hash         - (Done)
	Id           - (Done)
	Range        - (Done)
	Len          - (Done)
	DeleteFrom   - (Done)
	Dir          - ()
	Input        - (Done)
	ToString     - (Done)
	ToTuple      - (Done)
	ToArray      - (Done)
	ToInteger    - (Done)
	ToFloat      - (Done)
	ToBool       - (Done)
*/
func (p *Plasma) setBuiltInSymbols() {
	/*
		This is the master symbol table that is protected from writes
	*/
	p.builtInSymbolTable = NewSymbolTable(nil)

	// Types
	type_ := &Type{
		Object:      p.NewObject(ObjectName, nil, p.builtInSymbolTable),
		Constructor: NewBuiltInConstructor(p.ObjectInitialize),
		Name:        TypeName,
	}
	type_.Set(ToString,
		p.NewFunction(type_.symbols,
			NewBuiltInClassFunction(type_, 0,
				func(_ IObject, _ ...IObject) (IObject, *Object) {
					return p.NewString(p.PeekSymbolTable(), "Type@Object"), nil
				},
			),
		),
	)
	p.builtInSymbolTable.Set(TypeName, type_)
	//// Default Error Types
	exception := p.NewType(RuntimeError, p.builtInSymbolTable, []*Type{type_}, NewBuiltInConstructor(p.RuntimeErrorInitialize))
	p.builtInSymbolTable.Set(RuntimeError, exception)
	p.builtInSymbolTable.Set(InvalidTypeError,
		p.NewType(InvalidTypeError, p.builtInSymbolTable, []*Type{exception},
			NewBuiltInConstructor(
				func(object IObject) *Object {
					object.Set(Initialize,
						p.NewFunction(object.SymbolTable(),
							NewBuiltInClassFunction(object, 2,
								func(self IObject, arguments ...IObject) (IObject, *Object) {
									received := arguments[0]
									if _, ok := received.(*String); !ok {
										return nil, p.NewInvalidTypeError(received.TypeName(), StringName)
									}
									expecting := arguments[1]
									if _, ok := expecting.(*String); !ok {
										return nil, p.NewInvalidTypeError(expecting.TypeName(), StringName)
									}
									self.SetString(fmt.Sprintf("Expecting %s but received %s", expecting.GetString(), received.GetString()))
									return p.NewNone(), nil
								},
							),
						),
					)
					return nil
				},
			),
		),
	)
	p.builtInSymbolTable.Set(NotImplementedCallableError,
		p.NewType(NotImplementedCallableError, p.builtInSymbolTable, []*Type{exception},
			NewBuiltInConstructor(
				func(object IObject) *Object {
					object.Set(Initialize,
						p.NewFunction(object.SymbolTable(),
							NewBuiltInClassFunction(object, 1,
								func(self IObject, arguments ...IObject) (IObject, *Object) {
									methodNameObject := arguments[0]
									methodNameObjectToString, getError := methodNameObject.Get(ToString)
									if getError != nil {
										return nil, p.NewObjectWithNameNotFoundError(ToString)
									}
									if _, ok := methodNameObjectToString.(*Function); !ok {
										return nil, p.NewInvalidTypeError(methodNameObjectToString.TypeName(), FunctionName)
									}
									methodNameString, callError := p.CallFunction(methodNameObjectToString.(*Function), p.PeekSymbolTable())
									if callError != nil {
										return nil, callError
									}
									if _, ok := methodNameString.(*String); !ok {
										return nil, p.NewInvalidTypeError(methodNameString.TypeName(), StringName)
									}
									self.SetString(fmt.Sprintf("Method %s not implemented", methodNameString.GetString()))
									return p.NewNone(), nil
								},
							),
						),
					)
					return nil
				},
			),
		),
	)
	p.builtInSymbolTable.Set(ObjectConstructionError,
		p.NewType(ObjectConstructionError, p.builtInSymbolTable, []*Type{exception},
			NewBuiltInConstructor(
				func(object IObject) *Object {
					object.Set(Initialize,
						p.NewFunction(object.SymbolTable(),
							NewBuiltInClassFunction(object, 2,
								func(self IObject, arguments ...IObject) (IObject, *Object) {
									typeName := arguments[0]
									if _, ok := typeName.(*String); !ok {
										return nil, p.NewInvalidTypeError(typeName.TypeName(), StringName)
									}
									errorMessage := arguments[1]
									if _, ok := typeName.(*String); !ok {
										return nil, p.NewInvalidTypeError(errorMessage.TypeName(), StringName)
									}
									self.SetString(fmt.Sprintf("Could not construct object of Type: %s -> %s", typeName.GetString(), errorMessage.GetString()))
									return p.NewNone(), nil
								},
							),
						),
					)
					return nil
				},
			),
		),
	)
	p.builtInSymbolTable.Set(ObjectWithNameNotFoundError,
		p.NewType(ObjectWithNameNotFoundError, p.builtInSymbolTable, []*Type{exception},
			NewBuiltInConstructor(
				func(object IObject) *Object {
					object.Set(Initialize,
						p.NewFunction(object.SymbolTable(),
							NewBuiltInClassFunction(object, 1,
								func(self IObject, arguments ...IObject) (IObject, *Object) {
									name := arguments[0]
									if _, ok := name.(*String); !ok {
										return nil, p.NewInvalidTypeError(name.TypeName(), StringName)
									}
									self.SetString(fmt.Sprintf("Object with name %s not Found", name.GetString()))
									return p.NewNone(), nil
								},
							),
						),
					)
					return nil
				},
			),
		),
	)

	p.builtInSymbolTable.Set(InvalidNumberOfArgumentsError,
		p.NewType(InvalidNumberOfArgumentsError, p.builtInSymbolTable, []*Type{exception},
			NewBuiltInConstructor(
				func(object IObject) *Object {
					object.Set(Initialize,
						p.NewFunction(object.SymbolTable(),
							NewBuiltInClassFunction(object, 2,
								func(self IObject, arguments ...IObject) (IObject, *Object) {
									received := arguments[0]
									if _, ok := received.(*Integer); !ok {
										return nil, p.NewInvalidTypeError(received.TypeName(), IntegerName)
									}
									expecting := arguments[1]
									if _, ok := expecting.(*Integer); !ok {
										return nil, p.NewInvalidTypeError(expecting.TypeName(), IntegerName)
									}
									self.SetString(fmt.Sprintf("Received %d but expecting %d expecting", received.GetInteger64(), expecting.GetInteger64()))
									return p.NewNone(), nil
								},
							),
						),
					)
					return nil
				},
			),
		),
	)
	p.builtInSymbolTable.Set(GoRuntimeError,
		p.NewType(GoRuntimeError, p.builtInSymbolTable, []*Type{exception},
			NewBuiltInConstructor(
				func(object IObject) *Object {
					object.Set(Initialize,
						p.NewFunction(object.SymbolTable(),
							NewBuiltInClassFunction(object, 1,
								func(self IObject, arguments ...IObject) (IObject, *Object) {
									runtimeError := arguments[0]
									if _, ok := runtimeError.(*String); !ok {
										return nil, p.NewInvalidTypeError(runtimeError.TypeName(), StringName)
									}
									self.SetString(runtimeError.GetString())
									return p.NewNone(), nil
								},
							),
						),
					)
					return nil
				},
			),
		),
	)
	p.builtInSymbolTable.Set(UnhashableTypeError,
		p.NewType(UnhashableTypeError, p.builtInSymbolTable, []*Type{exception},
			NewBuiltInConstructor(
				func(object IObject) *Object {
					object.Set(Initialize,
						p.NewFunction(object.SymbolTable(),
							NewBuiltInClassFunction(object, 1,
								func(self IObject, arguments ...IObject) (IObject, *Object) {
									objectType := arguments[0]
									if _, ok := objectType.(*Type); !ok {
										return nil, p.NewInvalidTypeError(objectType.TypeName(), TypeName)
									}
									self.SetString(fmt.Sprintf("Object of type: %s is unhasable", objectType.(*Type).Name))
									return p.NewNone(), nil
								},
							),
						),
					)
					return nil
				},
			),
		),
	)
	p.builtInSymbolTable.Set(IndexOutOfRangeError,
		p.NewType(IndexOutOfRangeError, p.builtInSymbolTable, []*Type{exception},
			NewBuiltInConstructor(
				func(object IObject) *Object {
					object.Set(Initialize,
						p.NewFunction(object.SymbolTable(),
							NewBuiltInClassFunction(object, 2,
								func(self IObject, arguments ...IObject) (IObject, *Object) {
									length := arguments[0]
									if _, ok := length.(*Integer); !ok {
										return nil, p.NewInvalidTypeError(length.TypeName(), IntegerName)
									}
									index := arguments[1]
									if _, ok := length.(*Integer); !ok {
										return nil, p.NewInvalidTypeError(index.TypeName(), IntegerName)
									}
									self.SetString(fmt.Sprintf("Index: %d, out of range (Length=%d)", index.GetInteger64(), length.GetInteger64()))
									return p.NewNone(), nil
								},
							),
						),
					)
					return nil
				},
			),
		),
	)
	p.builtInSymbolTable.Set(KeyNotFoundError,
		p.NewType(KeyNotFoundError, p.builtInSymbolTable, []*Type{exception},
			NewBuiltInConstructor(
				func(object IObject) *Object {
					object.Set(Initialize,
						p.NewFunction(object.SymbolTable(),
							NewBuiltInClassFunction(object, 1,
								func(self IObject, arguments ...IObject) (IObject, *Object) {
									key := arguments[0]
									keyToString, getError := key.Get(ToString)
									if getError != nil {
										return nil, p.NewObjectWithNameNotFoundError(ToString)
									}
									if _, ok := keyToString.(*Function); !ok {
										return nil, p.NewInvalidTypeError(keyToString.TypeName(), FunctionName)
									}
									keyString, callError := p.CallFunction(keyToString.(*Function), p.PeekSymbolTable())
									if callError != nil {
										return nil, callError
									}
									if _, ok := keyString.(*String); !ok {
										return nil, p.NewInvalidTypeError(keyString.TypeName(), StringName)
									}
									self.SetString(fmt.Sprintf("Key %s not found", keyString.GetString()))
									return p.NewNone(), nil
								},
							),
						),
					)
					return nil
				},
			),
		),
	)
	p.builtInSymbolTable.Set(IntegerParsingError,
		p.NewType(IntegerParsingError, p.builtInSymbolTable, []*Type{exception},
			NewBuiltInConstructor(
				func(object IObject) *Object {
					object.Set(Initialize,
						p.NewFunction(object.SymbolTable(),
							NewBuiltInClassFunction(object, 0,
								func(self IObject, arguments ...IObject) (IObject, *Object) {
									self.SetString("Integer parsing error")
									return p.NewNone(), nil
								},
							),
						),
					)
					return nil
				},
			),
		),
	)
	p.builtInSymbolTable.Set(FloatParsingError,
		p.NewType(FloatParsingError, p.builtInSymbolTable, []*Type{exception},
			NewBuiltInConstructor(
				func(object IObject) *Object {
					object.Set(Initialize,
						p.NewFunction(object.SymbolTable(),
							NewBuiltInClassFunction(object, 0,
								func(self IObject, arguments ...IObject) (IObject, *Object) {
									self.SetString("Float parsing error")
									return p.NewNone(), nil
								},
							),
						),
					)
					return nil
				},
			),
		),
	)
	p.builtInSymbolTable.Set(BuiltInSymbolProtectionError,
		p.NewType(BuiltInSymbolProtectionError, p.builtInSymbolTable, []*Type{exception},
			NewBuiltInConstructor(
				func(object IObject) *Object {
					object.Set(Initialize,
						p.NewFunction(object.SymbolTable(),
							NewBuiltInClassFunction(object, 1,
								func(self IObject, arguments ...IObject) (IObject, *Object) {
									symbolName := arguments[0]
									if _, ok := symbolName.(*String); !ok {
										return nil, p.NewInvalidTypeError(symbolName.TypeName(), StringName)
									}
									self.SetString(fmt.Sprintf("cannot delete built-in symbol %s", symbolName.GetString()))
									return p.NewNone(), nil
								},
							),
						),
					)
					return nil
				},
			),
		),
	)
	//// Default Types
	p.builtInSymbolTable.Set(NoneName,
		p.NewType(NoneName, p.builtInSymbolTable, []*Type{type_},
			NewBuiltInConstructor(p.NoneInitialize),
		),
	)
	p.builtInSymbolTable.Set(BoolName,
		p.NewType(BoolName, p.builtInSymbolTable, []*Type{type_},
			NewBuiltInConstructor(p.BoolInitialize),
		),
	)
	p.builtInSymbolTable.Set(IteratorName,
		p.NewType(IteratorName, p.builtInSymbolTable, []*Type{type_},
			NewBuiltInConstructor(p.IteratorInitialize),
		),
	)
	p.builtInSymbolTable.Set(ObjectName,
		p.NewType(ObjectName, p.builtInSymbolTable, []*Type{type_},
			NewBuiltInConstructor(p.ObjectInitialize),
		),
	)
	p.builtInSymbolTable.Set(FunctionName,
		p.NewType(FunctionName, p.builtInSymbolTable, []*Type{type_},
			NewBuiltInConstructor(
				func(object IObject) *Object {
					return nil
				}),
		),
	)
	p.builtInSymbolTable.Set(IntegerName,
		p.NewType(IntegerName, p.builtInSymbolTable, []*Type{type_},
			NewBuiltInConstructor(p.IntegerInitialize),
		),
	)
	p.builtInSymbolTable.Set(StringName,
		p.NewType(StringName, p.builtInSymbolTable, []*Type{type_},
			NewBuiltInConstructor(p.StringInitialize),
		),
	)
	p.builtInSymbolTable.Set(BytesName,
		p.NewType(BytesName, p.builtInSymbolTable, []*Type{type_},
			NewBuiltInConstructor(p.BytesInitialize),
		),
	)
	p.builtInSymbolTable.Set(TupleName,
		p.NewType(TupleName, p.builtInSymbolTable, []*Type{type_},
			NewBuiltInConstructor(p.TupleInitialize),
		),
	)
	p.builtInSymbolTable.Set(ArrayName,
		p.NewType(ArrayName, p.builtInSymbolTable, []*Type{type_},
			NewBuiltInConstructor(p.ArrayInitialize),
		),
	)
	p.builtInSymbolTable.Set(HashName,
		p.NewType(HashName, p.builtInSymbolTable, []*Type{type_},
			NewBuiltInConstructor(p.HashTableInitialize),
		),
	)
	// Names

	// Functions
	p.builtInSymbolTable.Set("dir",
		p.NewFunction(p.builtInSymbolTable,
			NewBuiltInFunction(1,
				func(_ IObject, arguments ...IObject) (IObject, *Object) {
					object := arguments[0]
					var symbols []IObject
					for symbol := range object.SymbolTable().Symbols {
						symbols = append(symbols, p.NewString(p.PeekSymbolTable(), symbol))
					}
					return p.NewTuple(p.PeekSymbolTable(), symbols), nil
				},
			),
		),
	)
	p.builtInSymbolTable.Set("set",
		p.NewFunction(p.builtInSymbolTable,
			NewBuiltInFunction(3,
				func(_ IObject, arguments ...IObject) (IObject, *Object) {
					source := arguments[0]
					symbol := arguments[1]
					value := arguments[2]
					if _, ok := symbol.(*String); !ok {
						return nil, p.NewInvalidTypeError(symbol.TypeName(), StringName)
					}
					source.Set(symbol.GetString(), value)
					return p.NewNone(), nil
				},
			),
		),
	)
	p.builtInSymbolTable.Set("get_from",
		p.NewFunction(p.builtInSymbolTable,
			NewBuiltInFunction(2,
				func(_ IObject, arguments ...IObject) (IObject, *Object) {
					source := arguments[0]
					symbol := arguments[1]
					if _, ok := symbol.(*String); !ok {
						return nil, p.NewInvalidTypeError(symbol.TypeName(), StringName)
					}
					value, getError := source.Get(symbol.GetString())
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(symbol.GetString())
					}
					return value, nil
				},
			),
		),
	)
	p.builtInSymbolTable.Set("delete_from",
		p.NewFunction(p.builtInSymbolTable,
			NewBuiltInFunction(2,
				func(_ IObject, arguments ...IObject) (IObject, *Object) {
					source := arguments[0]
					symbol := arguments[1]
					if _, ok := symbol.(*String); !ok {
						return nil, p.NewInvalidTypeError(symbol.TypeName(), StringName)
					}
					_, getError := source.SymbolTable().GetSelf(symbol.GetString())
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(symbol.GetString())
					}
					delete(source.SymbolTable().Symbols, symbol.GetString())
					return p.NewNone(), nil
				},
			),
		),
	)
	p.builtInSymbolTable.Set("input",
		p.NewFunction(p.builtInSymbolTable,
			NewBuiltInFunction(1,
				func(_ IObject, arguments ...IObject) (IObject, *Object) {
					message := arguments[0]
					var messageString IObject
					if _, ok := message.(*String); !ok {
						messageToString, getError := message.Get(ToString)
						if getError != nil {
							return nil, p.NewObjectWithNameNotFoundError(ToString)
						}
						if _, ok = messageToString.(*Function); !ok {
							return nil, p.NewInvalidTypeError(messageToString.TypeName(), FunctionName)
						}
						toStringResult, callError := p.CallFunction(messageToString.(*Function), p.PeekSymbolTable())
						if callError != nil {
							return nil, callError
						}
						if _, ok = toStringResult.(*String); !ok {
							return nil, p.NewInvalidTypeError(toStringResult.TypeName(), StringName)
						}
						messageString = toStringResult
					} else {
						messageString = message
					}
					_, writingError := p.StdOut().Write([]byte(messageString.GetString()))
					if writingError != nil {
						return nil, p.NewGoRuntimeError(writingError)
					}
					if p.StdInScanner().Scan() {
						return p.NewString(p.PeekSymbolTable(), p.StdInScanner().Text()), nil
					}
					return nil, p.NewGoRuntimeError(p.StdInScanner().Err())
				},
			),
		),
	)
	p.builtInSymbolTable.Set("print",
		p.NewFunction(p.builtInSymbolTable,
			NewBuiltInFunction(1,
				func(_ IObject, arguments ...IObject) (IObject, *Object) {
					value := arguments[0]
					toString, getError := value.Get(ToString)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(ToString)
					}
					if _, ok := toString.(*Function); !ok {
						return nil, p.NewInvalidTypeError(toString.TypeName(), FunctionName)
					}
					stringValue, callError := p.CallFunction(toString.(*Function), value.SymbolTable())
					if callError != nil {
						return nil, callError
					}
					_, writeError := fmt.Fprintf(p.StdOut(), "%s", stringValue.GetString())
					if writeError != nil {
						return nil, p.NewGoRuntimeError(writeError)
					}
					return p.NewNone(), nil
				},
			),
		),
	)
	p.builtInSymbolTable.Set("println",
		p.NewFunction(p.builtInSymbolTable,
			NewBuiltInFunction(1,
				func(_ IObject, arguments ...IObject) (IObject, *Object) {
					value := arguments[0]
					toString, getError := value.Get(ToString)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(ToString)
					}
					if _, ok := toString.(*Function); !ok {
						return nil, p.NewInvalidTypeError(toString.TypeName(), FunctionName)
					}
					stringValue, callError := p.CallFunction(toString.(*Function), value.SymbolTable())
					if callError != nil {
						return nil, callError
					}
					_, writeError := fmt.Fprintf(p.StdOut(), "%s\n", stringValue.GetString())
					if writeError != nil {
						return nil, p.NewGoRuntimeError(writeError)
					}
					return p.NewNone(), nil
				},
			),
		),
	)
	p.builtInSymbolTable.Set("id",
		p.NewFunction(p.builtInSymbolTable,
			NewBuiltInFunction(1,
				func(_ IObject, arguments ...IObject) (IObject, *Object) {
					object := arguments[0]
					return p.NewInteger(p.PeekSymbolTable(), object.Id()), nil
				},
			),
		),
	)
	p.builtInSymbolTable.Set("hash",
		p.NewFunction(p.builtInSymbolTable,
			NewBuiltInFunction(1,
				func(_ IObject, arguments ...IObject) (IObject, *Object) {
					object := arguments[0]
					objectHashFunc, getError := object.Get(Hash)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(Hash)
					}
					if _, ok := objectHashFunc.(*Function); !ok {
						return nil, p.NewInvalidTypeError(objectHashFunc.TypeName(), FunctionName)
					}
					return p.CallFunction(objectHashFunc.(*Function), p.PeekSymbolTable())
				},
			),
		),
	)
	p.builtInSymbolTable.Set("range",
		p.NewFunction(p.builtInSymbolTable,
			NewBuiltInFunction(3,
				func(_ IObject, arguments ...IObject) (IObject, *Object) {
					start := arguments[0]
					if _, ok := start.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(start.TypeName(), IntegerName)
					}
					startValue := start.GetInteger64()

					end := arguments[1]
					if _, ok := end.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(end.TypeName(), IntegerName)
					}
					endValue := end.GetInteger64()

					step := arguments[2]
					if _, ok := step.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(step.TypeName(), IntegerName)
					}
					stepValue := step.GetInteger64()

					// This should return a iterator
					rangeIterator := p.NewIterator(p.PeekSymbolTable())
					rangeIterator.SetInteger64(startValue)

					rangeIterator.Set(HasNext,
						p.NewFunction(rangeIterator.SymbolTable(),
							NewBuiltInClassFunction(rangeIterator, 0,
								func(self IObject, _ ...IObject) (IObject, *Object) {
									return p.NewBool(p.PeekSymbolTable(), self.GetInteger64() < endValue), nil
								},
							),
						),
					)
					rangeIterator.Set(Next,
						p.NewFunction(rangeIterator.SymbolTable(),
							NewBuiltInClassFunction(rangeIterator, 0,
								func(self IObject, _ ...IObject) (IObject, *Object) {
									number := self.GetInteger64()
									self.SetInteger64(number + stepValue)
									return p.NewInteger(p.PeekSymbolTable(), number), nil
								},
							),
						),
					)

					return rangeIterator, nil
				},
			),
		),
	)
	p.builtInSymbolTable.Set("len",
		p.NewFunction(p.builtInSymbolTable,
			NewBuiltInFunction(1,
				func(_ IObject, arguments ...IObject) (IObject, *Object) {
					object := arguments[0]
					getLength, getError := object.Get(GetLength)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(GetLength)
					}
					if _, ok := getLength.(*Function); !ok {
						return nil, p.NewInvalidTypeError(getLength.TypeName(), FunctionName)
					}
					length, callError := p.CallFunction(getLength.(*Function), p.PeekSymbolTable())
					if callError != nil {
						return nil, callError
					}
					if _, ok := length.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(length.TypeName(), IntegerName)
					}
					return length, nil
				},
			),
		),
	)
	// To... (Transformations)
	p.builtInSymbolTable.Set(ToFloat,
		p.NewFunction(p.builtInSymbolTable,
			NewBuiltInFunction(1,
				func(_ IObject, arguments ...IObject) (IObject, *Object) {
					toFloat, getError := arguments[0].Get(ToFloat)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(ToFloat)
					}
					if _, ok := toFloat.(*Function); !ok {
						return nil, p.NewInvalidTypeError(toFloat.(IObject).TypeName(), FunctionName)
					}
					return p.CallFunction(toFloat.(*Function), arguments[0].SymbolTable().Parent)
				},
			),
		),
	)
	p.builtInSymbolTable.Set(ToString,
		p.NewFunction(p.builtInSymbolTable,
			NewBuiltInFunction(1,
				func(_ IObject, arguments ...IObject) (IObject, *Object) {
					toString, getError := arguments[0].Get(ToString)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(ToString)
					}
					if _, ok := toString.(*Function); !ok {
						return nil, p.NewInvalidTypeError(toString.(IObject).TypeName(), FunctionName)
					}
					return p.CallFunction(toString.(*Function), arguments[0].SymbolTable().Parent)
				},
			),
		),
	)
	p.builtInSymbolTable.Set(ToInteger,
		p.NewFunction(p.builtInSymbolTable,
			NewBuiltInFunction(1,
				func(_ IObject, arguments ...IObject) (IObject, *Object) {
					toInteger, getError := arguments[0].Get(ToInteger)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(ToInteger)
					}
					if _, ok := toInteger.(*Function); !ok {
						return nil, p.NewInvalidTypeError(toInteger.(IObject).TypeName(), FunctionName)
					}
					return p.CallFunction(toInteger.(*Function), arguments[0].SymbolTable().Parent)
				},
			),
		),
	)
	p.builtInSymbolTable.Set(ToArray,
		p.NewFunction(p.builtInSymbolTable,
			NewBuiltInFunction(1,
				func(_ IObject, arguments ...IObject) (IObject, *Object) {
					// First check if it is iterable
					// If not call its ToArray
					toArray, getError := arguments[0].Get(ToArray)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(ToArray)
					}
					if _, ok := toArray.(*Function); !ok {
						return nil, p.NewInvalidTypeError(toArray.(IObject).TypeName(), FunctionName)
					}
					return p.CallFunction(toArray.(*Function), arguments[0].SymbolTable().Parent)
				},
			),
		),
	)
	p.builtInSymbolTable.Set(ToTuple,
		p.NewFunction(p.builtInSymbolTable,
			NewBuiltInFunction(1,
				func(_ IObject, arguments ...IObject) (IObject, *Object) {
					// First check if it is iterable
					// If not call its ToTuple
					toTuple, getError := arguments[0].Get(ToTuple)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(ToTuple)
					}
					if _, ok := toTuple.(*Function); !ok {
						return nil, p.NewInvalidTypeError(toTuple.(IObject).TypeName(), FunctionName)
					}
					return p.CallFunction(toTuple.(*Function), arguments[0].SymbolTable().Parent)
				},
			),
		),
	)
	p.builtInSymbolTable.Set(ToBool,
		p.NewFunction(p.builtInSymbolTable,
			NewBuiltInFunction(1,
				func(_ IObject, arguments ...IObject) (IObject, *Object) {
					toBool, getError := arguments[0].Get(ToBool)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(ToBool)
					}
					if _, ok := toBool.(*Function); !ok {
						return nil, p.NewInvalidTypeError(toBool.(IObject).TypeName(), FunctionName)
					}
					return p.CallFunction(toBool.(*Function), arguments[0].SymbolTable().Parent)
				},
			),
		),
	)
	p.builtInSymbolTable = p.builtInSymbolTable
}
