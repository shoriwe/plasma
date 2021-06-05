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
	Hash         - ()
	Id           - ()
	Range        - ()
	Len          - ()
	ToString     - (Done)
	ToTuple      - (Done)
	ToArray      - (Done)
	ToInteger    - (Done)
	ToFloat      - (Done)
	ToBool       - (Done)
	ToHashTable  - ()
	ToIter       - ()
	ToBytes      - ()
	ToObject     - ()
*/
func (p *Plasma) SetDefaultSymbolTable() {
	symbolTable := NewSymbolTable(nil)
	noneObject := p.NewObject(NoneName, nil, symbolTable)

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
									return p.GetNone()
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
							NewBuiltInClassFunction(object, 0,
								func(self IObject, arguments ...IObject) (IObject, *Object) {
									self.SetString("Method not implemented")
									return p.GetNone()
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
									return p.GetNone()
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
									return p.GetNone()
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
									return p.GetNone()
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
	symbolTable.Set(None,
		noneObject,
	)
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
					return p.GetNone()
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
					return p.GetNone()
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