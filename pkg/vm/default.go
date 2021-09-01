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
	p.builtInContext.PeekSymbolTable().Set(TypeName, p.NewType(p.builtInContext, true, TypeName, nil, nil, NewBuiltInConstructor(p.ObjectInitialize(true))))
	//// Default Error Types
	exception := p.NewType(p.builtInContext, true, RuntimeError, p.builtInContext.PeekSymbolTable(), nil,
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
										self.String = fmt.Sprintf("Expecting %s but received %s", expecting.String, received.String)
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
										methodNameObjectToString, getError := methodNameObject.Get(p, context, ToString)
										if getError != nil {
											return getError, false
										}
										methodNameString, success := p.CallFunction(context, methodNameObjectToString)
										if !success {
											return methodNameString, false
										}
										if !methodNameString.IsTypeById(StringId) {
											return p.NewInvalidTypeError(context, methodNameString.TypeName(), StringName), false
										}
										self.String = fmt.Sprintf("Method %s not implemented", methodNameString.String)
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
										self.String = fmt.Sprintf("Could not construct object of Type: %s -> %s", typeName.String, errorMessage.String)
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
										self.String = fmt.Sprintf(" Value with name %s not Found inside object of type %s", name.String, objectType.Name)
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
										self.String = fmt.Sprintf("Received %d but expecting %d expecting", received.Integer, expecting.Integer)
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
										self.String = runtimeError.String
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
										self.String = fmt.Sprintf(" Value of type: %s is unhasable", objectType.Name)
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
										self.String = fmt.Sprintf("Index: %d, out of range (Length=%d)", index.Integer, length.Integer)
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
										keyToString, getError := key.Get(p, context, ToString)
										if getError != nil {
											return getError, false
										}
										keyString, success := p.CallFunction(context, keyToString)
										if !success {
											return keyString, false
										}
										if !keyString.IsTypeById(StringId) {
											return p.NewInvalidTypeError(context, keyString.TypeName(), StringName), false
										}
										self.String = fmt.Sprintf("Key %s not found", keyString.String)
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
										self.String = "Integer parsing error"
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
										self.String = "Float parsing error"
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
										self.String = fmt.Sprintf("cannot assign/delete built-in symbol %s", symbolName.String)
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
										self.String = fmt.Sprintf(" Value of type %s object is not callable", receivedType.Name)
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
		p.NewType(p.builtInContext, true, CallableName, p.builtInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.CallableInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(NoneName,
		p.NewType(p.builtInContext, true, NoneName, p.builtInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.NoneInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ModuleName,
		p.NewType(p.builtInContext, true, ModuleName, p.builtInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.ObjectInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(BoolName,
		p.NewType(p.builtInContext, true, BoolName, p.builtInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.BoolInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(IteratorName,
		p.NewType(p.builtInContext, true, IteratorName, p.builtInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.IteratorInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(FloatName,
		p.NewType(p.builtInContext, true, FloatName, p.builtInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.FloatInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ValueName,
		p.NewType(p.builtInContext, true, ValueName,
			nil, nil,
			NewBuiltInConstructor(p.ObjectInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(FunctionName,
		p.NewType(p.builtInContext, true, FunctionName, p.builtInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.ObjectInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(IntegerName,
		p.NewType(p.builtInContext, true, IntegerName, p.builtInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.IntegerInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(StringName,
		p.NewType(p.builtInContext, true, StringName, p.builtInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.StringInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(BytesName,
		p.NewType(p.builtInContext, true, BytesName, p.builtInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.BytesInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(TupleName,
		p.NewType(p.builtInContext, true, TupleName, p.builtInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.TupleInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(ArrayName,
		p.NewType(p.builtInContext, true, ArrayName, p.builtInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.ArrayInitialize(false)),
		),
	)
	p.builtInContext.PeekSymbolTable().Set(HashName,
		p.NewType(p.builtInContext, true, HashName, p.builtInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.HashTableInitialize(false)),
		),
	)
	// Names
	p.builtInContext.PeekSymbolTable().Set(TrueName, p.NewBool(p.builtInContext, true, true))
	p.builtInContext.PeekSymbolTable().Set(FalseName, p.NewBool(p.builtInContext, true, false))
	p.builtInContext.PeekSymbolTable().Set(None, p.NewNone(p.builtInContext, true, p.builtInContext.PeekSymbolTable()))
	// Functions
	p.builtInContext.PeekSymbolTable().Set("expand",
		p.NewFunction(p.builtInContext, true, p.builtInContext.PeekSymbolTable(),
			NewBuiltInFunction(2,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					receiver := arguments[0]
					for symbol, object := range arguments[1].SymbolTable().Symbols {
						receiver.Set(p, p.builtInContext, symbol, object)
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
						symbols = append(symbols, p.NewString(p.builtInContext, false, symbol))
					}
					return p.NewTuple(p.builtInContext, false, symbols), true
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
					if source.IsBuiltIn {
						return p.NewBuiltInSymbolProtectionError(p.builtInContext, symbol.String), false
					}
					source.Set(p, p.builtInContext, symbol.String, value)
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
					value, getError := source.Get(p, p.builtInContext, symbol.String)
					if getError != nil {
						return getError, false
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
					if source.IsBuiltIn {
						return p.NewBuiltInSymbolProtectionError(p.builtInContext, symbol.String), false
					}
					_, getError := source.SymbolTable().GetSelf(symbol.String)
					if getError != nil {
						return p.NewObjectWithNameNotFoundError(p.builtInContext, source.GetClass(p), symbol.String), false
					}
					delete(source.SymbolTable().Symbols, symbol.String)
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
						messageToString, getError := message.Get(p, p.builtInContext, ToString)
						if getError != nil {
							return getError, false
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
					_, writingError := p.StdOut().Write([]byte(messageString.String))
					if writingError != nil {
						return p.NewGoRuntimeError(p.builtInContext, writingError), false
					}
					if p.StdInScanner().Scan() {
						return p.NewString(p.builtInContext, false, p.StdInScanner().Text()), true
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
					toString, getError := value.Get(p, p.builtInContext, ToString)
					if getError != nil {
						return getError, false
					}
					stringValue, success := p.CallFunction(p.builtInContext, toString)
					if !success {
						return stringValue, false
					}
					_, writeError := fmt.Fprintf(p.StdOut(), "%s", stringValue.String)
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
					toString, getError := value.Get(p, p.builtInContext, ToString)
					if getError != nil {
						return getError, false
					}
					stringValue, success := p.CallFunction(p.builtInContext, toString)
					if !success {
						return stringValue, false
					}
					_, writeError := fmt.Fprintf(p.StdOut(), "%s\n", stringValue.String)
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
					return p.NewInteger(p.builtInContext, false, object.Id()), true
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

					end := arguments[1]
					if !end.IsTypeById(IntegerId) {
						return p.NewInvalidTypeError(p.builtInContext, end.TypeName(), IntegerName), false
					}

					step := arguments[2]
					if !step.IsTypeById(IntegerId) {
						return p.NewInvalidTypeError(p.builtInContext, step.TypeName(), IntegerName), false
					}

					rangeInformation := struct {
						current int64
						end     int64
						step    int64
					}{
						current: start.Integer,
						end:     end.Integer,
						step:    step.Integer,
					}

					// This should return a iterator
					rangeIterator := p.NewIterator(p.builtInContext, false)

					rangeIterator.Set(p, p.builtInContext, HasNext,
						p.NewFunction(p.builtInContext, true, rangeIterator.SymbolTable(),
							NewBuiltInClassFunction(rangeIterator, 0,
								func(self *Value, _ ...*Value) (*Value, bool) {
									if rangeInformation.current < rangeInformation.end {
										return p.GetTrue(), true
									}
									return p.GetFalse(), true
								},
							),
						),
					)
					rangeIterator.Set(p, p.builtInContext, Next,
						p.NewFunction(p.builtInContext, true, rangeIterator.SymbolTable(),
							NewBuiltInClassFunction(rangeIterator, 0,
								func(self *Value, _ ...*Value) (*Value, bool) {
									number := rangeInformation.current
									rangeInformation.current += rangeInformation.step
									return p.NewInteger(p.builtInContext, false, number), true
								},
							),
						),
					)
					return rangeIterator, true
				},
			),
		),
	)
}
