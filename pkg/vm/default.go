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
	Expand       - (Done)
	Id           - (Done)
	Range        - (Done)
	Len          - (Done)
	DeleteFrom   - (Done)
	Dir          - (Done)
	Input        - (Done)
	ToString     - (Done)
	ToTuple      - (Done)
	ToArray      - (Done)
	ToInteger    - (Done)
	ToFloat      - (Done)
	ToBool       - (Done)
*/
func (p *Plasma) InitializeBuiltIn() {
	if !p.builtInContext.SymbolStack.HasNext() {

	}
	/*
		This is the master symbol table that is protected from writes
	*/

	// Types
	type_ := &Type{
		Object:      p.NewObject(p.builtInContext, true, ObjectName, nil, p.builtInContext.PeekSymbolTable()),
		Constructor: NewBuiltInConstructor(p.ObjectInitialize(true)),
		Name:        TypeName,
	}
	type_.Set(ToString,
		p.NewFunction(p.builtInContext, true, type_.symbols,
			NewBuiltInClassFunction(type_, 0,
				func(_ Value, _ ...Value) (Value, *Object) {
					return p.NewString(p.builtInContext, false, p.builtInContext.PeekSymbolTable(), "Type@Object"), nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(TypeName, type_)
	//// Default Error Types
	exception := p.NewType(p.builtInContext, true, RuntimeError, p.builtInContext.PeekSymbolTable(), []*Type{type_},
		NewBuiltInConstructor(p.RuntimeErrorInitialize),
	)
	p.builtInContext.PeekSymbolTable().Set(RuntimeError, exception)
	p.builtInContext.PeekSymbolTable().Set(InvalidTypeError,
		p.NewType(p.builtInContext, true, InvalidTypeError, p.builtInContext.PeekSymbolTable(), []*Type{exception},
			NewBuiltInConstructor(
				func(context *Context, object Value) *Object {
					object.SetOnDemandSymbol(Initialize,
						func() Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 2,
									func(self Value, arguments ...Value) (Value, *Object) {
										received := arguments[0]
										if _, ok := received.(*String); !ok {
											return nil, p.NewInvalidTypeError(context, received.TypeName(), StringName)
										}
										expecting := arguments[1]
										if _, ok := expecting.(*String); !ok {
											return nil, p.NewInvalidTypeError(context, expecting.TypeName(), StringName)
										}
										self.SetString(fmt.Sprintf("Expecting %s but received %s", expecting.GetString(), received.GetString()))
										return p.GetNone(), nil
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(NotImplementedCallableError,
		p.NewType(p.builtInContext, true, NotImplementedCallableError, p.builtInContext.PeekSymbolTable(), []*Type{exception},
			NewBuiltInConstructor(
				func(context *Context, object Value) *Object {
					object.SetOnDemandSymbol(Initialize,
						func() Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 1,
									func(self Value, arguments ...Value) (Value, *Object) {
										methodNameObject := arguments[0]
										methodNameObjectToString, getError := methodNameObject.Get(ToString)
										if getError != nil {
											return nil, p.NewObjectWithNameNotFoundError(context, methodNameObject.GetClass(p), ToString)
										}
										methodNameString, callError := p.CallFunction(context, methodNameObjectToString, context.PeekSymbolTable())
										if callError != nil {
											return nil, callError
										}
										if _, ok := methodNameString.(*String); !ok {
											return nil, p.NewInvalidTypeError(context, methodNameString.TypeName(), StringName)
										}
										self.SetString(fmt.Sprintf("Method %s not implemented", methodNameString.GetString()))
										return p.GetNone(), nil
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ObjectConstructionError,
		p.NewType(p.builtInContext, true, ObjectConstructionError, p.builtInContext.PeekSymbolTable(), []*Type{exception},
			NewBuiltInConstructor(
				func(context *Context, object Value) *Object {
					object.SetOnDemandSymbol(Initialize,
						func() Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 2,
									func(self Value, arguments ...Value) (Value, *Object) {
										typeName := arguments[0]
										if _, ok := typeName.(*String); !ok {
											return nil, p.NewInvalidTypeError(context, typeName.TypeName(), StringName)
										}
										errorMessage := arguments[1]
										if _, ok := typeName.(*String); !ok {
											return nil, p.NewInvalidTypeError(context, errorMessage.TypeName(), StringName)
										}
										self.SetString(fmt.Sprintf("Could not construct object of Type: %s -> %s", typeName.GetString(), errorMessage.GetString()))
										return p.GetNone(), nil
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ObjectWithNameNotFoundError,
		p.NewType(p.builtInContext, true, ObjectWithNameNotFoundError, p.builtInContext.PeekSymbolTable(), []*Type{exception},
			NewBuiltInConstructor(
				func(context *Context, object Value) *Object {
					object.SetOnDemandSymbol(Initialize,
						func() Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 2,
									func(self Value, arguments ...Value) (Value, *Object) {
										objectType := arguments[0]
										if _, ok := objectType.(*Type); !ok {
											return nil, p.NewInvalidTypeError(context, objectType.TypeName(), TypeName)
										}
										name := arguments[1]
										if _, ok := name.(*String); !ok {
											return nil, p.NewInvalidTypeError(context, name.TypeName(), StringName)
										}
										self.SetString(fmt.Sprintf("Object with name %s not Found inside object of type %s", name.GetString(), objectType.(*Type).Name))
										return p.GetNone(), nil
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)

	p.builtInContext.PeekSymbolTable().Set(InvalidNumberOfArgumentsError,
		p.NewType(p.builtInContext, true, InvalidNumberOfArgumentsError, p.builtInContext.PeekSymbolTable(), []*Type{exception},
			NewBuiltInConstructor(
				func(context *Context, object Value) *Object {
					object.SetOnDemandSymbol(Initialize,
						func() Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 2,
									func(self Value, arguments ...Value) (Value, *Object) {
										received := arguments[0]
										if _, ok := received.(*Integer); !ok {
											return nil, p.NewInvalidTypeError(context, received.TypeName(), IntegerName)
										}
										expecting := arguments[1]
										if _, ok := expecting.(*Integer); !ok {
											return nil, p.NewInvalidTypeError(context, expecting.TypeName(), IntegerName)
										}
										self.SetString(fmt.Sprintf("Received %d but expecting %d expecting", received.GetInteger(), expecting.GetInteger()))
										return p.GetNone(), nil
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(GoRuntimeError,
		p.NewType(p.builtInContext, true, GoRuntimeError, p.builtInContext.PeekSymbolTable(), []*Type{exception},
			NewBuiltInConstructor(
				func(context *Context, object Value) *Object {
					object.SetOnDemandSymbol(Initialize,
						func() Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 1,
									func(self Value, arguments ...Value) (Value, *Object) {
										runtimeError := arguments[0]
										if _, ok := runtimeError.(*String); !ok {
											return nil, p.NewInvalidTypeError(context, runtimeError.TypeName(), StringName)
										}
										self.SetString(runtimeError.GetString())
										return p.GetNone(), nil
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(UnhashableTypeError,
		p.NewType(p.builtInContext, true, UnhashableTypeError, p.builtInContext.PeekSymbolTable(), []*Type{exception},
			NewBuiltInConstructor(
				func(context *Context, object Value) *Object {
					object.SetOnDemandSymbol(Initialize,
						func() Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 1,
									func(self Value, arguments ...Value) (Value, *Object) {
										objectType := arguments[0]
										if _, ok := objectType.(*Type); !ok {
											return nil, p.NewInvalidTypeError(context, objectType.TypeName(), TypeName)
										}
										self.SetString(fmt.Sprintf("Object of type: %s is unhasable", objectType.(*Type).Name))
										return p.GetNone(), nil
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(IndexOutOfRangeError,
		p.NewType(p.builtInContext, true, IndexOutOfRangeError, p.builtInContext.PeekSymbolTable(), []*Type{exception},
			NewBuiltInConstructor(
				func(context *Context, object Value) *Object {
					object.SetOnDemandSymbol(Initialize,
						func() Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 2,
									func(self Value, arguments ...Value) (Value, *Object) {
										length := arguments[0]
										if _, ok := length.(*Integer); !ok {
											return nil, p.NewInvalidTypeError(context, length.TypeName(), IntegerName)
										}
										index := arguments[1]
										if _, ok := length.(*Integer); !ok {
											return nil, p.NewInvalidTypeError(context, index.TypeName(), IntegerName)
										}
										self.SetString(fmt.Sprintf("Index: %d, out of range (Length=%d)", index.GetInteger(), length.GetInteger()))
										return p.GetNone(), nil
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(KeyNotFoundError,
		p.NewType(p.builtInContext, true, KeyNotFoundError, p.builtInContext.PeekSymbolTable(), []*Type{exception},
			NewBuiltInConstructor(
				func(context *Context, object Value) *Object {
					object.SetOnDemandSymbol(Initialize,
						func() Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 1,
									func(self Value, arguments ...Value) (Value, *Object) {
										key := arguments[0]
										keyToString, getError := key.Get(ToString)
										if getError != nil {
											return nil, p.NewObjectWithNameNotFoundError(context, key.GetClass(p), ToString)
										}
										keyString, callError := p.CallFunction(context, keyToString, context.PeekSymbolTable())
										if callError != nil {
											return nil, callError
										}
										if _, ok := keyString.(*String); !ok {
											return nil, p.NewInvalidTypeError(context, keyString.TypeName(), StringName)
										}
										self.SetString(fmt.Sprintf("Key %s not found", keyString.GetString()))
										return p.GetNone(), nil
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(IntegerParsingError,
		p.NewType(p.builtInContext, true, IntegerParsingError, p.builtInContext.PeekSymbolTable(), []*Type{exception},
			NewBuiltInConstructor(
				func(context *Context, object Value) *Object {
					object.SetOnDemandSymbol(Initialize,
						func() Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 0,
									func(self Value, arguments ...Value) (Value, *Object) {
										self.SetString("Integer parsing error")
										return p.GetNone(), nil
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(FloatParsingError,
		p.NewType(p.builtInContext, true, FloatParsingError, p.builtInContext.PeekSymbolTable(), []*Type{exception},
			NewBuiltInConstructor(
				func(context *Context, object Value) *Object {
					object.SetOnDemandSymbol(Initialize,
						func() Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 0,
									func(self Value, arguments ...Value) (Value, *Object) {
										self.SetString("Float parsing error")
										return p.GetNone(), nil
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(BuiltInSymbolProtectionError,
		p.NewType(p.builtInContext, true, BuiltInSymbolProtectionError, p.builtInContext.PeekSymbolTable(), []*Type{exception},
			NewBuiltInConstructor(
				func(context *Context, object Value) *Object {
					object.SetOnDemandSymbol(Initialize,
						func() Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 1,
									func(self Value, arguments ...Value) (Value, *Object) {
										symbolName := arguments[0]
										if _, ok := symbolName.(*String); !ok {
											return nil, p.NewInvalidTypeError(context, symbolName.TypeName(), StringName)
										}
										self.SetString(fmt.Sprintf("cannot assign/delete built-in symbol %s", symbolName.GetString()))
										return p.GetNone(), nil
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ObjectNotCallableError,
		p.NewType(p.builtInContext, true, ObjectNotCallableError, p.builtInContext.PeekSymbolTable(), []*Type{exception},
			NewBuiltInConstructor(
				func(context *Context, object Value) *Object {
					object.SetOnDemandSymbol(Initialize,
						func() Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 1,
									func(self Value, arguments ...Value) (Value, *Object) {
										receivedType := arguments[0]
										if _, ok := receivedType.(*Type); !ok {
											return nil, p.NewInvalidTypeError(context, receivedType.TypeName(), TypeName)
										}
										self.SetString(fmt.Sprintf("Object of type %s object is not callable", receivedType.(*Type).Name))
										return p.GetNone(), nil
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	//// Default Types
	p.builtInContext.PeekSymbolTable().Set(CallableName,
		p.NewType(p.builtInContext, true, CallableName, p.builtInContext.PeekSymbolTable(), []*Type{type_},
			NewBuiltInConstructor(p.CallableInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(NoneName,
		p.NewType(p.builtInContext, true, NoneName, p.builtInContext.PeekSymbolTable(), []*Type{type_},
			NewBuiltInConstructor(p.NoneInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ModuleName,
		p.NewType(p.builtInContext, true, ModuleName, p.builtInContext.PeekSymbolTable(), []*Type{type_},
			NewBuiltInConstructor(p.ObjectInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(BoolName,
		p.NewType(p.builtInContext, true, BoolName, p.builtInContext.PeekSymbolTable(), []*Type{type_},
			NewBuiltInConstructor(p.BoolInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(IteratorName,
		p.NewType(p.builtInContext, true, IteratorName, p.builtInContext.PeekSymbolTable(), []*Type{type_},
			NewBuiltInConstructor(p.IteratorInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(FloatName,
		p.NewType(p.builtInContext, true, FloatName, p.builtInContext.PeekSymbolTable(), []*Type{type_},
			NewBuiltInConstructor(p.FloatInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ObjectName,
		p.NewType(p.builtInContext, true, ObjectName, p.builtInContext.PeekSymbolTable(), []*Type{type_},
			NewBuiltInConstructor(p.ObjectInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(FunctionName,
		p.NewType(p.builtInContext, true, FunctionName, p.builtInContext.PeekSymbolTable(), []*Type{type_},
			NewBuiltInConstructor(p.ObjectInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(IntegerName,
		p.NewType(p.builtInContext, true, IntegerName, p.builtInContext.PeekSymbolTable(), []*Type{type_},
			NewBuiltInConstructor(p.IntegerInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(StringName,
		p.NewType(p.builtInContext, true, StringName, p.builtInContext.PeekSymbolTable(), []*Type{type_},
			NewBuiltInConstructor(p.StringInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(BytesName,
		p.NewType(p.builtInContext, true, BytesName, p.builtInContext.PeekSymbolTable(), []*Type{type_},
			NewBuiltInConstructor(p.BytesInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(TupleName,
		p.NewType(p.builtInContext, true, TupleName, p.builtInContext.PeekSymbolTable(), []*Type{type_},
			NewBuiltInConstructor(p.TupleInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ArrayName,
		p.NewType(p.builtInContext, true, ArrayName, p.builtInContext.PeekSymbolTable(), []*Type{type_},
			NewBuiltInConstructor(p.ArrayInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(HashName,
		p.NewType(p.builtInContext, true, HashName, p.builtInContext.PeekSymbolTable(), []*Type{type_},
			NewBuiltInConstructor(p.HashTableInitialize(false)),
		),
	)
	// Names
	p.builtInContext.PeekSymbolTable().Set(TrueName, p.NewBool(p.builtInContext, true, p.builtInContext.PeekSymbolTable(), true))
	p.builtInContext.PeekSymbolTable().Set(FalseName, p.NewBool(p.builtInContext, true, p.builtInContext.PeekSymbolTable(), false))
	p.builtInContext.PeekSymbolTable().Set(None, p.NewNone(p.builtInContext, true, p.builtInContext.PeekSymbolTable()))
	// Functions
	p.builtInContext.PeekSymbolTable().Set("expand",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(2,
				func(_ Value, arguments ...Value) (Value, *Object) {
					receiver := arguments[0]
					for symbol, object := range arguments[1].SymbolTable().Symbols {
						receiver.Set(symbol, object)
					}
					return p.GetNone(), nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("dir",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ Value, arguments ...Value) (Value, *Object) {
					object := arguments[0]
					var symbols []Value
					for symbol := range object.Dir() {
						symbols = append(symbols, p.NewString(p.builtInContext, false, p.builtInContext.PeekSymbolTable(), symbol))
					}
					return p.NewTuple(p.builtInContext, false, p.builtInContext.PeekSymbolTable(), symbols), nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("set",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(3,
				func(_ Value, arguments ...Value) (Value, *Object) {
					source := arguments[0]
					symbol := arguments[1]
					value := arguments[2]
					if _, ok := symbol.(*String); !ok {
						return nil, p.NewInvalidTypeError(p.builtInContext, symbol.TypeName(), StringName)
					}
					if source.IsBuiltIn() {
						return nil, p.NewBuiltInSymbolProtectionError(p.builtInContext, symbol.GetString())
					}
					source.Set(symbol.GetString(), value)
					return p.GetNone(), nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("get_from",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(2,
				func(_ Value, arguments ...Value) (Value, *Object) {
					source := arguments[0]
					symbol := arguments[1]
					if _, ok := symbol.(*String); !ok {
						return nil, p.NewInvalidTypeError(p.builtInContext, symbol.TypeName(), StringName)
					}
					value, getError := source.Get(symbol.GetString())
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(p.builtInContext, source.GetClass(p), symbol.GetString())
					}
					return value, nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("delete_from",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(2,
				func(_ Value, arguments ...Value) (Value, *Object) {
					source := arguments[0]
					symbol := arguments[1]
					if _, ok := symbol.(*String); !ok {
						return nil, p.NewInvalidTypeError(p.builtInContext, symbol.TypeName(), StringName)
					}
					if source.IsBuiltIn() {
						return nil, p.NewBuiltInSymbolProtectionError(p.builtInContext, symbol.GetString())
					}
					_, getError := source.SymbolTable().GetSelf(symbol.GetString())
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(p.builtInContext, source.GetClass(p), symbol.GetString())
					}
					delete(source.SymbolTable().Symbols, symbol.GetString())
					return p.GetNone(), nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("input",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ Value, arguments ...Value) (Value, *Object) {
					message := arguments[0]
					var messageString Value
					if _, ok := message.(*String); !ok {
						messageToString, getError := message.Get(ToString)
						if getError != nil {
							return nil, p.NewObjectWithNameNotFoundError(p.builtInContext, message.GetClass(p), ToString)
						}
						toStringResult, callError := p.CallFunction(p.builtInContext, messageToString, p.builtInContext.PeekSymbolTable())
						if callError != nil {
							return nil, callError
						}
						if _, ok = toStringResult.(*String); !ok {
							return nil, p.NewInvalidTypeError(p.builtInContext, toStringResult.TypeName(), StringName)
						}
						messageString = toStringResult
					} else {
						messageString = message
					}
					_, writingError := p.StdOut().Write([]byte(messageString.GetString()))
					if writingError != nil {
						return nil, p.NewGoRuntimeError(p.builtInContext, writingError)
					}
					if p.StdInScanner().Scan() {
						return p.NewString(p.builtInContext, false, p.builtInContext.PeekSymbolTable(), p.StdInScanner().Text()), nil
					}
					return nil, p.NewGoRuntimeError(p.builtInContext, p.StdInScanner().Err())
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("print",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ Value, arguments ...Value) (Value, *Object) {
					value := arguments[0]
					toString, getError := value.Get(ToString)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(p.builtInContext, value.GetClass(p), ToString)
					}
					stringValue, callError := p.CallFunction(p.builtInContext, toString, value.SymbolTable())
					if callError != nil {
						return nil, callError
					}
					_, writeError := fmt.Fprintf(p.StdOut(), "%s", stringValue.GetString())
					if writeError != nil {
						return nil, p.NewGoRuntimeError(p.builtInContext, writeError)
					}
					return p.GetNone(), nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("println",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ Value, arguments ...Value) (Value, *Object) {
					value := arguments[0]
					toString, getError := value.Get(ToString)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(p.builtInContext, value.GetClass(p), ToString)
					}
					stringValue, callError := p.CallFunction(p.builtInContext, toString, value.SymbolTable())
					if callError != nil {
						return nil, callError
					}
					_, writeError := fmt.Fprintf(p.StdOut(), "%s\n", stringValue.GetString())
					if writeError != nil {
						return nil, p.NewGoRuntimeError(p.builtInContext, writeError)
					}
					return p.GetNone(), nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("id",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ Value, arguments ...Value) (Value, *Object) {
					object := arguments[0]
					return p.NewInteger(p.builtInContext, false, p.builtInContext.PeekSymbolTable(), object.Id()), nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("hash",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ Value, arguments ...Value) (Value, *Object) {
					object := arguments[0]
					objectHashFunc, getError := object.Get(Hash)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(p.builtInContext, object.GetClass(p), Hash)
					}
					return p.CallFunction(p.builtInContext, objectHashFunc, p.builtInContext.PeekSymbolTable())
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("range",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(3,
				func(_ Value, arguments ...Value) (Value, *Object) {
					start := arguments[0]
					if _, ok := start.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(p.builtInContext, start.TypeName(), IntegerName)
					}
					startValue := start.GetInteger()

					end := arguments[1]
					if _, ok := end.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(p.builtInContext, end.TypeName(), IntegerName)
					}
					endValue := end.GetInteger()

					step := arguments[2]
					if _, ok := step.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(p.builtInContext, step.TypeName(), IntegerName)
					}
					stepValue := step.GetInteger()

					// This should return a iterator
					rangeIterator := p.NewIterator(p.builtInContext, true, p.builtInContext.PeekSymbolTable())
					rangeIterator.SetInteger(startValue)

					rangeIterator.Set(HasNext,
						p.NewFunction(p.builtInContext, true, rangeIterator.SymbolTable(),
							NewBuiltInClassFunction(rangeIterator, 0,
								func(self Value, _ ...Value) (Value, *Object) {
									if self.GetInteger() < endValue {
										return p.GetTrue(), nil
									}
									return p.GetFalse(), nil
								},
							),
						),
					)
					rangeIterator.Set(Next,
						p.NewFunction(p.builtInContext, true, rangeIterator.SymbolTable(),
							NewBuiltInClassFunction(rangeIterator, 0,
								func(self Value, _ ...Value) (Value, *Object) {
									number := self.GetInteger()
									self.SetInteger(number + stepValue)
									return p.NewInteger(p.builtInContext, false, p.builtInContext.PeekSymbolTable(), number), nil
								},
							),
						),
					)

					return rangeIterator, nil
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("len",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ Value, arguments ...Value) (Value, *Object) {
					object := arguments[0]
					getLength, getError := object.Get(GetLength)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(p.builtInContext, object.GetClass(p), GetLength)
					}
					length, callError := p.CallFunction(p.builtInContext, getLength, p.builtInContext.PeekSymbolTable())
					if callError != nil {
						return nil, callError
					}
					if _, ok := length.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(p.builtInContext, length.TypeName(), IntegerName)
					}
					return length, nil
				},
			),
		),
	)
	// To... (Transformations)
	p.builtInContext.PeekSymbolTable().Set(ToFloat,
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ Value, arguments ...Value) (Value, *Object) {
					toFloat, getError := arguments[0].Get(ToFloat)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(p.builtInContext, arguments[0].GetClass(p), ToFloat)
					}
					return p.CallFunction(p.builtInContext, toFloat, arguments[0].SymbolTable().Parent)
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ToString,
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ Value, arguments ...Value) (Value, *Object) {
					toString, getError := arguments[0].Get(ToString)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(p.builtInContext, arguments[0].GetClass(p), ToString)
					}
					return p.CallFunction(p.builtInContext, toString, arguments[0].SymbolTable().Parent)
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ToInteger,
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ Value, arguments ...Value) (Value, *Object) {
					toInteger, getError := arguments[0].Get(ToInteger)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(p.builtInContext, arguments[0].GetClass(p), ToInteger)
					}
					return p.CallFunction(p.builtInContext, toInteger, arguments[0].SymbolTable().Parent)
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ToArray,
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ Value, arguments ...Value) (Value, *Object) {
					// First check if it is iterable
					// If not call its ToArray
					toArray, getError := arguments[0].Get(ToArray)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(p.builtInContext, arguments[0].GetClass(p), ToArray)
					}
					return p.CallFunction(p.builtInContext, toArray, arguments[0].SymbolTable().Parent)
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ToTuple,
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ Value, arguments ...Value) (Value, *Object) {
					// First check if it is iterable
					// If not call its ToTuple
					toTuple, getError := arguments[0].Get(ToTuple)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(p.builtInContext, arguments[0].GetClass(p), ToTuple)
					}
					return p.CallFunction(p.builtInContext, toTuple, arguments[0].SymbolTable().Parent)
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ToBool,
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ Value, arguments ...Value) (Value, *Object) {
					toBool, getError := arguments[0].Get(ToBool)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(p.builtInContext, arguments[0].GetClass(p), ToBool)
					}
					return p.CallFunction(p.builtInContext, toBool, arguments[0].SymbolTable().Parent)
				},
			),
		),
	)
}
