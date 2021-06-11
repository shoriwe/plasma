package vm

import (
	"fmt"
)

/*
	SetDefaultSymbolTable
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
	Delete		 - ()
	ToString     - (Done)
	ToTuple      - (Done)
	ToArray      - (Done)
	ToInteger    - (Done)
	ToFloat      - (Done)
	ToBool       - (Done)
	ToHashTable  - ()
*/
func (p *Plasma) SetDefaultSymbolTable() {
	symbolTable := NewSymbolTable(nil)

	// Types
	type_ := &Type{
		Object:      p.NewObject(ObjectName, nil, symbolTable),
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
	symbolTable.Set(TypeName, type_)
	//// Default Error Types
	exception := p.NewType(RuntimeError, symbolTable, []*Type{type_}, NewBuiltInConstructor(p.RuntimeErrorInitialize))
	symbolTable.Set(RuntimeError, exception)
	symbolTable.Set(InvalidTypeError,
		p.NewType(InvalidTypeError, symbolTable, []*Type{exception},
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
	symbolTable.Set(NotImplementedCallableError,
		p.NewType(NotImplementedCallableError, symbolTable, []*Type{exception},
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
	symbolTable.Set(ObjectConstructionError,
		p.NewType(ObjectConstructionError, symbolTable, []*Type{exception},
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
	symbolTable.Set(ObjectWithNameNotFoundError,
		p.NewType(ObjectWithNameNotFoundError, symbolTable, []*Type{exception},
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

	symbolTable.Set(InvalidNumberOfArgumentsError,
		p.NewType(InvalidNumberOfArgumentsError, symbolTable, []*Type{exception},
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
	symbolTable.Set(GoRuntimeError,
		p.NewType(GoRuntimeError, symbolTable, []*Type{exception},
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
	symbolTable.Set(UnhashableTypeError,
		p.NewType(UnhashableTypeError, symbolTable, []*Type{exception},
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
	symbolTable.Set(IndexOutOfRangeError,
		p.NewType(IndexOutOfRangeError, symbolTable, []*Type{exception},
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
	symbolTable.Set(KeyNotFoundError,
		p.NewType(KeyNotFoundError, symbolTable, []*Type{exception},
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
	symbolTable.Set(IntegerParsingError,
		p.NewType(IntegerParsingError, symbolTable, []*Type{exception},
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
	symbolTable.Set(FloatParsingError,
		p.NewType(FloatParsingError, symbolTable, []*Type{exception},
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
	//// Default Types
	symbolTable.Set(NoneName,
		p.NewType(NoneName, symbolTable, []*Type{type_},
			NewBuiltInConstructor(p.NoneInitialize),
		),
	)
	symbolTable.Set(BoolName,
		p.NewType(BoolName, symbolTable, []*Type{type_},
			NewBuiltInConstructor(p.BoolInitialize),
		),
	)
	symbolTable.Set(IteratorName,
		p.NewType(IteratorName, symbolTable, []*Type{type_},
			NewBuiltInConstructor(p.IteratorInitialize),
		),
	)
	symbolTable.Set(ObjectName,
		p.NewType(ObjectName, symbolTable, []*Type{type_},
			NewBuiltInConstructor(p.ObjectInitialize),
		),
	)
	symbolTable.Set(FunctionName,
		p.NewType(FunctionName, symbolTable, []*Type{type_},
			NewBuiltInConstructor(
				func(object IObject) *Object {
					return nil
				}),
		),
	)
	symbolTable.Set(IntegerName,
		p.NewType(IntegerName, symbolTable, []*Type{type_},
			NewBuiltInConstructor(p.IntegerInitialize),
		),
	)
	symbolTable.Set(StringName,
		p.NewType(StringName, symbolTable, []*Type{type_},
			NewBuiltInConstructor(p.StringInitialize),
		),
	)
	symbolTable.Set(BytesName,
		p.NewType(BytesName, symbolTable, []*Type{type_},
			NewBuiltInConstructor(p.BytesInitialize),
		),
	)
	symbolTable.Set(TupleName,
		p.NewType(TupleName, symbolTable, []*Type{type_},
			NewBuiltInConstructor(p.TupleInitialize),
		),
	)
	symbolTable.Set(ArrayName,
		p.NewType(ArrayName, symbolTable, []*Type{type_},
			NewBuiltInConstructor(p.ArrayInitialize),
		),
	)
	symbolTable.Set(HashName,
		p.NewType(HashName, symbolTable, []*Type{type_},
			NewBuiltInConstructor(p.HashTableInitialize),
		),
	)
	// Names

	// Functions
	symbolTable.Set("print",
		p.NewFunction(symbolTable,
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
	symbolTable.Set("println",
		p.NewFunction(symbolTable,
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
	symbolTable.Set("id",
		p.NewFunction(symbolTable,
			NewBuiltInFunction(1,
				func(_ IObject, arguments ...IObject) (IObject, *Object) {
					object := arguments[0]
					// ToDo: Fix me, what about those times the id is greater than int64 max value
					return p.NewInteger(p.PeekSymbolTable(), int64(object.Id())), nil
				},
			),
		),
	)
	symbolTable.Set("hash",
		p.NewFunction(symbolTable,
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
	symbolTable.Set("range",
		p.NewFunction(symbolTable,
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
	symbolTable.Set("len",
		p.NewFunction(symbolTable,
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
	symbolTable.Set(ToFloat,
		p.NewFunction(symbolTable,
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
	symbolTable.Set(ToString,
		p.NewFunction(symbolTable,
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
	symbolTable.Set(ToInteger,
		p.NewFunction(symbolTable,
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
	symbolTable.Set(ToArray,
		p.NewFunction(symbolTable,
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
	symbolTable.Set(ToTuple,
		p.NewFunction(symbolTable,
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
	symbolTable.Set(ToBool,
		p.NewFunction(symbolTable,
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
	p.programMasterSymbolTable = symbolTable
}
