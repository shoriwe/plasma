package vm

import (
	"fmt"
)

/*
	Type         - (Done)
	Function     - (Done)
	 Value       - (Done)
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
	type_ := p.NewType(p.builtInContext, true, TypeName, nil, nil, NewBuiltInConstructor(p.ObjectInitialize(true)))

	type_.Set(ToString,
		p.NewFunction(p.builtInContext, true, type_.symbols,
			NewBuiltInClassFunction(type_, 0,
				func(_ *Value, _ ...*Value) (*Value, bool) {
					return p.NewString(p.builtInContext, false, p.builtInContext.PeekSymbolTable(), "Type@ Value"), true
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(TypeName, type_)
	//// Default Error Types
	exception := p.NewType(p.builtInContext, true, RuntimeError, p.builtInContext.PeekSymbolTable(), []*Value{type_},
		NewBuiltInConstructor(p.RuntimeErrorInitialize),
	)
	p.builtInContext.PeekSymbolTable().Set(RuntimeError, exception)
	p.builtInContext.PeekSymbolTable().Set(InvalidTypeError,
		p.NewType(p.builtInContext, true, InvalidTypeError, p.builtInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 2,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										received := arguments[0]
										if !received.IsTypeById(StringId) {
											return p.NewInvalidTypeError(context, received.TypeName(), StringName), false
										}
										expecting := arguments[1]
										if !expecting.IsTypeById(StringId) {
											return p.NewInvalidTypeError(context, expecting.TypeName(), StringName), false
										}
										self.SetString(fmt.Sprintf("Expecting %s but received %s", expecting.GetString(), received.GetString()))
										return p.GetNone(), true
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
		p.NewType(p.builtInContext, true, NotImplementedCallableError, p.builtInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 1,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										methodNameObject := arguments[0]
										methodNameObjectToString, getError := methodNameObject.Get(ToString)
										if getError != nil {
											return p.NewObjectWithNameNotFoundError(context, methodNameObject.GetClass(p), ToString), false
										}
										methodNameString, success := p.CallFunction(context, methodNameObjectToString)
										if !success {
											return methodNameString, false
										}
										if !methodNameString.IsTypeById(StringId) {
											return p.NewInvalidTypeError(context, methodNameString.TypeName(), StringName), false
										}
										self.SetString(fmt.Sprintf("Method %s not implemented", methodNameString.GetString()))
										return p.GetNone(), true
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
		p.NewType(p.builtInContext, true, ObjectConstructionError, p.builtInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 2,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										typeName := arguments[0]
										if !typeName.IsTypeById(StringId) {
											return p.NewInvalidTypeError(context, typeName.TypeName(), StringName), false
										}
										errorMessage := arguments[1]
										if !typeName.IsTypeById(StringId) {
											return p.NewInvalidTypeError(context, errorMessage.TypeName(), StringName), false
										}
										self.SetString(fmt.Sprintf("Could not construct object of Type: %s -> %s", typeName.GetString(), errorMessage.GetString()))
										return p.GetNone(), true
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
		p.NewType(p.builtInContext, true, ObjectWithNameNotFoundError, p.builtInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 2,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										objectType := arguments[0]
										if !objectType.IsTypeById(TypeId) {
											return p.NewInvalidTypeError(context, objectType.TypeName(), TypeName), false
										}
										name := arguments[1]
										if !name.IsTypeById(StringId) {
											return p.NewInvalidTypeError(context, name.TypeName(), StringName), false
										}
										self.SetString(fmt.Sprintf(" Value with name %s not Found inside object of type %s", name.GetString(), objectType.Name))
										return p.GetNone(), true
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
		p.NewType(p.builtInContext, true, InvalidNumberOfArgumentsError, p.builtInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 2,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										received := arguments[0]
										if !received.IsTypeById(IntegerId) {
											return p.NewInvalidTypeError(context, received.TypeName(), IntegerName), false
										}
										expecting := arguments[1]
										if !expecting.IsTypeById(IntegerId) {
											return p.NewInvalidTypeError(context, expecting.TypeName(), IntegerName), false
										}
										self.SetString(fmt.Sprintf("Received %d but expecting %d expecting", received.GetInteger(), expecting.GetInteger()))
										return p.GetNone(), true
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
		p.NewType(p.builtInContext, true, GoRuntimeError, p.builtInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 1,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										runtimeError := arguments[0]
										if !runtimeError.IsTypeById(StringId) {
											return p.NewInvalidTypeError(context, runtimeError.TypeName(), StringName), false
										}
										self.SetString(runtimeError.GetString())
										return p.GetNone(), true
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
		p.NewType(p.builtInContext, true, UnhashableTypeError, p.builtInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 1,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										objectType := arguments[0]
										if !objectType.IsTypeById(TypeId) {
											return p.NewInvalidTypeError(context, objectType.TypeName(), TypeName), false
										}
										self.SetString(fmt.Sprintf(" Value of type: %s is unhasable", objectType.Name))
										return p.GetNone(), true
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
		p.NewType(p.builtInContext, true, IndexOutOfRangeError, p.builtInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 2,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										length := arguments[0]
										if !length.IsTypeById(IntegerId) {
											return p.NewInvalidTypeError(context, length.TypeName(), IntegerName), false
										}
										index := arguments[1]
										if !index.IsTypeById(IntegerId) {
											return p.NewInvalidTypeError(context, index.TypeName(), IntegerName), false
										}
										self.SetString(fmt.Sprintf("Index: %d, out of range (Length=%d)", index.GetInteger(), length.GetInteger()))
										return p.GetNone(), true
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
		p.NewType(p.builtInContext, true, KeyNotFoundError, p.builtInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 1,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										key := arguments[0]
										keyToString, getError := key.Get(ToString)
										if getError != nil {
											return p.NewObjectWithNameNotFoundError(context, key.GetClass(p), ToString), false
										}
										keyString, success := p.CallFunction(context, keyToString)
										if !success {
											return keyString, false
										}
										if !keyString.IsTypeById(StringId) {
											return p.NewInvalidTypeError(context, keyString.TypeName(), StringName), false
										}
										self.SetString(fmt.Sprintf("Key %s not found", keyString.GetString()))
										return p.GetNone(), true
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
		p.NewType(p.builtInContext, true, IntegerParsingError, p.builtInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 0,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										self.SetString("Integer parsing error")
										return p.GetNone(), true
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
		p.NewType(p.builtInContext, true, FloatParsingError, p.builtInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 0,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										self.SetString("Float parsing error")
										return p.GetNone(), true
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
		p.NewType(p.builtInContext, true, BuiltInSymbolProtectionError, p.builtInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 1,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										symbolName := arguments[0]
										if !symbolName.IsTypeById(StringId) {
											return p.NewInvalidTypeError(context, symbolName.TypeName(), StringName), false
										}
										self.SetString(fmt.Sprintf("cannot assign/delete built-in symbol %s", symbolName.GetString()))
										return p.GetNone(), true
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
		p.NewType(p.builtInContext, true, ObjectNotCallableError, p.builtInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 1,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										receivedType := arguments[0]
										if !receivedType.IsTypeById(TypeId) {
											return p.NewInvalidTypeError(context, receivedType.TypeName(), TypeName), false
										}
										self.SetString(fmt.Sprintf(" Value of type %s object is not callable", receivedType.Name))
										return p.GetNone(), true
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
		p.NewType(p.builtInContext, true, CallableName, p.builtInContext.PeekSymbolTable(), []*Value{type_},
			NewBuiltInConstructor(p.CallableInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(NoneName,
		p.NewType(p.builtInContext, true, NoneName, p.builtInContext.PeekSymbolTable(), []*Value{type_},
			NewBuiltInConstructor(p.NoneInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ModuleName,
		p.NewType(p.builtInContext, true, ModuleName, p.builtInContext.PeekSymbolTable(), []*Value{type_},
			NewBuiltInConstructor(p.ObjectInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(BoolName,
		p.NewType(p.builtInContext, true, BoolName, p.builtInContext.PeekSymbolTable(), []*Value{type_},
			NewBuiltInConstructor(p.BoolInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(IteratorName,
		p.NewType(p.builtInContext, true, IteratorName, p.builtInContext.PeekSymbolTable(), []*Value{type_},
			NewBuiltInConstructor(p.IteratorInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(FloatName,
		p.NewType(p.builtInContext, true, FloatName, p.builtInContext.PeekSymbolTable(), []*Value{type_},
			NewBuiltInConstructor(p.FloatInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ValueName,
		p.NewType(p.builtInContext, true, ValueName, p.builtInContext.PeekSymbolTable(), []*Value{type_},
			NewBuiltInConstructor(p.ObjectInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(FunctionName,
		p.NewType(p.builtInContext, true, FunctionName, p.builtInContext.PeekSymbolTable(), []*Value{type_},
			NewBuiltInConstructor(p.ObjectInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(IntegerName,
		p.NewType(p.builtInContext, true, IntegerName, p.builtInContext.PeekSymbolTable(), []*Value{type_},
			NewBuiltInConstructor(p.IntegerInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(StringName,
		p.NewType(p.builtInContext, true, StringName, p.builtInContext.PeekSymbolTable(), []*Value{type_},
			NewBuiltInConstructor(p.StringInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(BytesName,
		p.NewType(p.builtInContext, true, BytesName, p.builtInContext.PeekSymbolTable(), []*Value{type_},
			NewBuiltInConstructor(p.BytesInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(TupleName,
		p.NewType(p.builtInContext, true, TupleName, p.builtInContext.PeekSymbolTable(), []*Value{type_},
			NewBuiltInConstructor(p.TupleInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ArrayName,
		p.NewType(p.builtInContext, true, ArrayName, p.builtInContext.PeekSymbolTable(), []*Value{type_},
			NewBuiltInConstructor(p.ArrayInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(HashName,
		p.NewType(p.builtInContext, true, HashName, p.builtInContext.PeekSymbolTable(), []*Value{type_},
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
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					receiver := arguments[0]
					for symbol, object := range arguments[1].SymbolTable().Symbols {
						receiver.Set(symbol, object)
					}
					return p.GetNone(), true
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("dir",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					object := arguments[0]
					var symbols []*Value
					for symbol := range object.Dir() {
						symbols = append(symbols, p.NewString(p.builtInContext, false, p.builtInContext.PeekSymbolTable(), symbol))
					}
					return p.NewTuple(p.builtInContext, false, p.builtInContext.PeekSymbolTable(), symbols), true
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("set",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(3,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					source := arguments[0]
					symbol := arguments[1]
					value := arguments[2]
					if !symbol.IsTypeById(StringId) {
						return p.NewInvalidTypeError(p.builtInContext, symbol.TypeName(), StringName), false
					}
					if source.IsBuiltIn() {
						return p.NewBuiltInSymbolProtectionError(p.builtInContext, symbol.GetString()), false
					}
					source.Set(symbol.GetString(), value)
					return p.GetNone(), true
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("get_from",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(2,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					source := arguments[0]
					symbol := arguments[1]
					if !symbol.IsTypeById(StringId) {
						return p.NewInvalidTypeError(p.builtInContext, symbol.TypeName(), StringName), false
					}
					value, getError := source.Get(symbol.GetString())
					if getError != nil {
						return p.NewObjectWithNameNotFoundError(p.builtInContext, source.GetClass(p), symbol.GetString()), false
					}
					return value, true
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("delete_from",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(2,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					source := arguments[0]
					symbol := arguments[1]
					if !symbol.IsTypeById(StringId) {
						return p.NewInvalidTypeError(p.builtInContext, symbol.TypeName(), StringName), false
					}
					if source.IsBuiltIn() {
						return p.NewBuiltInSymbolProtectionError(p.builtInContext, symbol.GetString()), false
					}
					_, getError := source.SymbolTable().GetSelf(symbol.GetString())
					if getError != nil {
						return p.NewObjectWithNameNotFoundError(p.builtInContext, source.GetClass(p), symbol.GetString()), false
					}
					delete(source.SymbolTable().Symbols, symbol.GetString())
					return p.GetNone(), true
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("input",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					message := arguments[0]
					var messageString *Value
					if !message.IsTypeById(StringId) {
						messageToString, getError := message.Get(ToString)
						if getError != nil {
							return p.NewObjectWithNameNotFoundError(p.builtInContext, message.GetClass(p), ToString), false
						}
						toStringResult, success := p.CallFunction(p.builtInContext, messageToString)
						if !success {
							return toStringResult, false
						}
						if !toStringResult.IsTypeById(StringId) {
							return p.NewInvalidTypeError(p.builtInContext, toStringResult.TypeName(), StringName), false
						}
						messageString = toStringResult
					} else {
						messageString = message
					}
					_, writingError := p.StdOut().Write([]byte(messageString.GetString()))
					if writingError != nil {
						return p.NewGoRuntimeError(p.builtInContext, writingError), false
					}
					if p.StdInScanner().Scan() {
						return p.NewString(p.builtInContext, false, p.builtInContext.PeekSymbolTable(), p.StdInScanner().Text()), true
					}
					return p.NewGoRuntimeError(p.builtInContext, p.StdInScanner().Err()), false
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("print",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					value := arguments[0]
					toString, getError := value.Get(ToString)
					if getError != nil {
						return p.NewObjectWithNameNotFoundError(p.builtInContext, value.GetClass(p), ToString), false
					}
					stringValue, success := p.CallFunction(p.builtInContext, toString)
					if !success {
						return stringValue, false
					}
					_, writeError := fmt.Fprintf(p.StdOut(), "%s", stringValue.GetString())
					if writeError != nil {
						return p.NewGoRuntimeError(p.builtInContext, writeError), false
					}
					return p.GetNone(), true
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("println",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					value := arguments[0]
					toString, getError := value.Get(ToString)
					if getError != nil {
						return p.NewObjectWithNameNotFoundError(p.builtInContext, value.GetClass(p), ToString), false
					}
					stringValue, success := p.CallFunction(p.builtInContext, toString)
					if !success {
						return stringValue, false
					}
					_, writeError := fmt.Fprintf(p.StdOut(), "%s\n", stringValue.GetString())
					if writeError != nil {
						return p.NewGoRuntimeError(p.builtInContext, writeError), false
					}
					return p.GetNone(), true
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("id",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					object := arguments[0]
					return p.NewInteger(p.builtInContext, false, p.builtInContext.PeekSymbolTable(), object.Id()), true
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("hash",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					object := arguments[0]
					objectHashFunc, getError := object.Get(Hash)
					if getError != nil {
						return p.NewObjectWithNameNotFoundError(p.builtInContext, object.GetClass(p), Hash), false
					}
					return p.CallFunction(p.builtInContext, objectHashFunc)
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("range",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(3,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					start := arguments[0]
					if !start.IsTypeById(IntegerId) {
						return p.NewInvalidTypeError(p.builtInContext, start.TypeName(), IntegerName), false
					}
					startValue := start.GetInteger()

					end := arguments[1]
					if !end.IsTypeById(IntegerId) {
						return p.NewInvalidTypeError(p.builtInContext, end.TypeName(), IntegerName), false
					}
					endValue := end.GetInteger()

					step := arguments[2]
					if !step.IsTypeById(IntegerId) {
						return p.NewInvalidTypeError(p.builtInContext, step.TypeName(), IntegerName), false
					}
					stepValue := step.GetInteger()

					// This should return a iterator
					rangeIterator := p.NewIterator(p.builtInContext, true, p.builtInContext.PeekSymbolTable())
					rangeIterator.SetInteger(startValue)

					rangeIterator.Set(HasNext,
						p.NewFunction(p.builtInContext, true, rangeIterator.SymbolTable(),
							NewBuiltInClassFunction(rangeIterator, 0,
								func(self *Value, _ ...*Value) (*Value, bool) {
									if self.GetInteger() < endValue {
										return p.GetTrue(), true
									}
									return p.GetFalse(), true
								},
							),
						),
					)
					rangeIterator.Set(Next,
						p.NewFunction(p.builtInContext, true, rangeIterator.SymbolTable(),
							NewBuiltInClassFunction(rangeIterator, 0,
								func(self *Value, _ ...*Value) (*Value, bool) {
									number := self.GetInteger()
									self.SetInteger(number + stepValue)
									return p.NewInteger(p.builtInContext, false, p.builtInContext.PeekSymbolTable(), number), true
								},
							),
						),
					)

					return rangeIterator, true
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set("len",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					object := arguments[0]
					getLength, getError := object.Get(GetLength)
					if getError != nil {
						return p.NewObjectWithNameNotFoundError(p.builtInContext, object.GetClass(p), GetLength), false
					}
					length, success := p.CallFunction(p.builtInContext, getLength)
					if !success {
						return length, false
					}
					if !length.IsTypeById(IntegerId) {
						return p.NewInvalidTypeError(p.builtInContext, length.TypeName(), IntegerName), false
					}
					return length, true
				},
			),
		),
	)
	// To... (Transformations)
	p.builtInContext.PeekSymbolTable().Set(ToFloat,
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					toFloat, getError := arguments[0].Get(ToFloat)
					if getError != nil {
						return p.NewObjectWithNameNotFoundError(p.builtInContext, arguments[0].GetClass(p), ToFloat), false
					}
					return p.CallFunction(p.builtInContext, toFloat)
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ToString,
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					toString, getError := arguments[0].Get(ToString)
					if getError != nil {
						return p.NewObjectWithNameNotFoundError(p.builtInContext, arguments[0].GetClass(p), ToString), false
					}
					return p.CallFunction(p.builtInContext, toString)
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ToInteger,
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					toInteger, getError := arguments[0].Get(ToInteger)
					if getError != nil {
						return p.NewObjectWithNameNotFoundError(p.builtInContext, arguments[0].GetClass(p), ToInteger), false
					}
					return p.CallFunction(p.builtInContext, toInteger)
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ToArray,
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					// First check if it is iterable
					// If not call its ToArray
					toArray, getError := arguments[0].Get(ToArray)
					if getError != nil {
						return p.NewObjectWithNameNotFoundError(p.builtInContext, arguments[0].GetClass(p), ToArray), false
					}
					return p.CallFunction(p.builtInContext, toArray)
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ToTuple,
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					// First check if it is iterable
					// If not call its ToTuple
					toTuple, getError := arguments[0].Get(ToTuple)
					if getError != nil {
						return p.NewObjectWithNameNotFoundError(p.builtInContext, arguments[0].GetClass(p), ToTuple), false
					}
					return p.CallFunction(p.builtInContext, toTuple)
				},
			),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ToBool,
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					toBool, getError := arguments[0].Get(ToBool)
					if getError != nil {
						return p.NewObjectWithNameNotFoundError(p.builtInContext, arguments[0].GetClass(p), ToBool), false
					}
					return p.CallFunction(p.builtInContext, toBool)
				},
			),
		),
	)
}
