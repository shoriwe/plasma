package vm

type KeyValue struct {
	Key   Value
	Value Value
}

type HashTable struct {
	*Object
}

func (p *Plasma) NewHashTable(context *Context, isBuiltIn bool, parent *SymbolTable, entries map[int64][]*KeyValue, entriesLength int) *HashTable {
	hashTable := &HashTable{
		Object: p.NewObject(context, isBuiltIn, HashName, nil, parent),
	}
	hashTable.SetKeyValues(entries)
	hashTable.SetLength(entriesLength)
	p.HashTableInitialize(isBuiltIn)(context, hashTable)
	hashTable.SetOnDemandSymbol(Self,
		func() Value {
			return hashTable
		},
	)
	return hashTable
}

func (p *Plasma) HashTableInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object Value) *Object {
		object.SetOnDemandSymbol(Equals,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {

							rawRight := arguments[0]
							if _, ok := rawRight.(*HashTable); !ok {
								return p.GetFalse(), nil
							}
							right := rawRight.(*HashTable)
							if self.GetLength() != right.Length {
								return p.GetFalse(), nil
							}
							rightIndex, getError := right.Get(Index)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(context, right.GetClass(p), Index)
							}
							for key, leftValue := range self.GetKeyValues() {
								// Check if other has the key
								rightValue, ok := right.KeyValues[key]
								if !ok {
									return p.GetFalse(), nil
								}
								// Check if the each entry one has the same length
								if len(leftValue) != len(rightValue) {
									return p.GetFalse(), nil
								}
								// Start comparing the entries
								for _, entry := range leftValue {
									_, indexingError := p.CallFunction(context, rightIndex, context.PeekSymbolTable(), entry.Key)
									if indexingError != nil {
										return p.GetFalse(), nil
									}
								}
							}
							return p.GetTrue(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightEquals,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {

							rawLeft := arguments[0]
							if _, ok := rawLeft.(*HashTable); !ok {
								return p.GetFalse(), nil
							}
							left := rawLeft.(*HashTable)
							if self.GetLength() != left.Length {
								return p.GetFalse(), nil
							}
							leftIndex, getError := left.Get(Index)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(context, left.GetClass(p), Index)
							}
							for key, leftValue := range left.KeyValues {
								// Check if other has the key
								rightValue, ok := self.GetKeyValues()[key]
								if !ok {
									return p.GetFalse(), nil
								}
								// Check if the each entry one has the same length
								if len(leftValue) != len(rightValue) {
									return p.GetFalse(), nil
								}
								// Start comparing the entries
								for _, entry := range leftValue {
									_, indexingError := p.CallFunction(context, leftIndex, context.PeekSymbolTable(), entry.Key)
									if indexingError != nil {
										return p.GetFalse(), nil
									}
								}
							}
							return p.GetTrue(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(NotEquals,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							rawRight := arguments[0]
							if _, ok := rawRight.(*HashTable); !ok {
								return p.GetTrue(), nil
							}
							right := rawRight.(*HashTable)
							if self.GetLength() != right.Length {
								return p.GetTrue(), nil
							}
							rightIndex, getError := right.Get(Index)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(context, right.GetClass(p), Index)
							}
							for key, leftValue := range self.GetKeyValues() {
								// Check if other has the key
								rightValue, ok := right.KeyValues[key]
								if !ok {
									return p.GetTrue(), nil
								}
								// Check if the each entry one has the same length
								if len(leftValue) != len(rightValue) {
									return p.GetTrue(), nil
								}
								// Start comparing the entries
								for _, entry := range leftValue {
									_, indexingError := p.CallFunction(context, rightIndex, context.PeekSymbolTable(), entry.Key)
									if indexingError != nil {
										return p.GetTrue(), nil
									}
								}
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightNotEquals,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							rawLeft := arguments[0]
							if _, ok := rawLeft.(*HashTable); !ok {
								return p.GetTrue(), nil
							}
							left := rawLeft.(*HashTable)
							if self.GetLength() != left.Length {
								return p.GetTrue(), nil
							}
							leftIndex, getError := left.Get(Index)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(context, left.GetClass(p), Index)
							}
							for key, leftValue := range left.KeyValues {
								// Check if other has the key
								rightValue, ok := self.GetKeyValues()[key]
								if !ok {
									return p.GetTrue(), nil
								}
								// Check if the each entry one has the same length
								if len(leftValue) != len(rightValue) {
									return p.GetTrue(), nil
								}
								// Start comparing the entries
								for _, entry := range leftValue {
									_, indexingError := p.CallFunction(context, leftIndex, context.PeekSymbolTable(), entry.Key)
									if indexingError != nil {
										return p.GetTrue(), nil
									}
								}
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Contains,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							value := arguments[0]
							valueHashFunc, getError := value.Get(Hash)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(context, value.GetClass(p), Hash)
							}
							valueHashObject, callError := p.CallFunction(context, valueHashFunc, context.PeekSymbolTable())
							if callError != nil {
								return nil, callError
							}
							if _, ok := valueHashObject.(*Integer); !ok {
								return nil, p.NewInvalidTypeError(context, valueHashObject.TypeName(), IntegerName)
							}
							valueHash := valueHashObject.GetInteger()
							entries, found := self.GetKeyValues()[valueHash]
							if !found {
								return p.GetFalse(), nil
							}
							var valueEquals Value
							valueEquals, getError = value.Get(RightEquals)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(context, value.GetClass(p), RightEquals)
							}
							var comparisonResult Value
							var comparisonResultBool bool
							for _, entry := range entries {
								comparisonResult, callError = p.CallFunction(context, valueEquals, context.PeekSymbolTable(), entry.Key)
								if callError != nil {
									return nil, callError
								}
								comparisonResultBool, callError = p.QuickGetBool(context, comparisonResult)
								if callError != nil {
									return nil, callError
								}
								if comparisonResultBool {
									return p.GetTrue(), nil
								}
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightContains,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							value := arguments[0]
							valueHashFunc, getError := value.Get(Hash)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(context, value.GetClass(p), Hash)
							}
							valueHashObject, callError := p.CallFunction(context, valueHashFunc, context.PeekSymbolTable())
							if callError != nil {
								return nil, callError
							}
							if _, ok := valueHashObject.(*Integer); !ok {
								return nil, p.NewInvalidTypeError(context, valueHashObject.TypeName(), IntegerName)
							}
							valueHash := valueHashObject.GetInteger()
							entries, found := self.GetKeyValues()[valueHash]
							if !found {
								return p.GetFalse(), nil
							}
							var valueEquals Value
							valueEquals, getError = value.Get(Equals)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(context, value.GetClass(p), Equals)
							}
							var comparisonResult Value
							var comparisonResultBool bool
							for _, entry := range entries {
								comparisonResult, callError = p.CallFunction(context, valueEquals, context.PeekSymbolTable(), entry.Key)
								if callError != nil {
									return nil, callError
								}
								comparisonResultBool, callError = p.QuickGetBool(context, comparisonResult)
								if callError != nil {
									return nil, callError
								}
								if comparisonResultBool {
									return p.GetTrue(), nil
								}
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Hash,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(_ Value, _ ...Value) (Value, *Object) {
							return nil, p.NewUnhashableTypeError(context, p.ForceMasterGetAny(HashName).(*Type))
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Copy,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(_ Value, _ ...Value) (Value, *Object) {
							return nil, p.NewUnhashableTypeError(context, object.GetClass(p))
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Index,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							indexObject := arguments[0]
							indexObjectHash, getError := indexObject.Get(Hash)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(context, indexObject.GetClass(p), Hash)
							}
							indexHash, callError := p.CallFunction(context, indexObjectHash, indexObject.SymbolTable())
							if callError != nil {
								return nil, callError
							}
							if _, ok := indexHash.(*Integer); !ok {
								return nil, p.NewInvalidTypeError(context, indexHash.TypeName(), IntegerName)
							}
							keyValues, found := self.GetKeyValues()[indexHash.GetInteger()]
							if !found {
								return nil, p.NewKeyNotFoundError(context, indexObject)
							}
							var indexObjectEquals Value
							indexObjectEquals, getError = indexObject.Get(Equals)
							var equals Value
							var equalsBool bool
							for _, keyValue := range keyValues {
								equals, callError = p.CallFunction(context, indexObjectEquals, indexObject.SymbolTable(), keyValue.Key)
								if callError != nil {
									return nil, callError
								}
								equalsBool, callError = p.QuickGetBool(context, equals)
								if callError != nil {
									return nil, callError
								}
								if equalsBool {
									return keyValue.Value, nil
								}
							}
							return nil, p.NewKeyNotFoundError(context, indexObject)
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Assign,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 2,
						func(self Value, arguments ...Value) (Value, *Object) {
							indexObject := arguments[0]
							newValue := arguments[1]
							indexObjectHash, getError := indexObject.Get(Hash)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(context, indexObject.GetClass(p), Hash)
							}
							indexHash, callError := p.CallFunction(context, indexObjectHash, indexObject.SymbolTable())
							if callError != nil {
								return nil, callError
							}
							if _, ok := indexHash.(*Integer); !ok {
								return nil, p.NewInvalidTypeError(context, indexHash.TypeName(), IntegerName)
							}
							keyValues, found := self.GetKeyValues()[indexHash.GetInteger()]
							if found {
								self.AddKeyValue(indexHash.GetInteger(), &KeyValue{
									Key:   indexObject,
									Value: newValue,
								})
								return p.GetNone(), nil
							}
							var indexObjectEquals Value
							indexObjectEquals, getError = indexObject.Get(Equals)
							var equals Value
							var equalsBool bool
							for index, keyValue := range keyValues {
								equals, callError = p.CallFunction(context, indexObjectEquals, indexObject.SymbolTable(), keyValue.Key)
								if callError != nil {
									return nil, callError
								}
								equalsBool, callError = p.QuickGetBool(context, equals)
								if callError != nil {
									return nil, callError
								}
								if equalsBool {
									self.GetKeyValues()[indexHash.GetInteger()][index].Value = newValue
									return p.GetNone(), nil
								}
							}
							self.IncreaseLength()
							self.GetKeyValues()[indexHash.GetInteger()] = append(
								self.GetKeyValues()[indexHash.GetInteger()],
								&KeyValue{
									Key:   indexObject,
									Value: newValue,
								},
							)
							return p.GetNone(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Iter,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							toTuple, getError := self.Get(ToTuple)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(context, self.GetClass(p), ToTuple)
							}
							hashKeys, callError := p.CallFunction(context, toTuple, self.SymbolTable())
							if callError != nil {
								return nil, callError
							}
							iterator := p.NewIterator(context, false, context.PeekSymbolTable())
							iterator.SetInteger(0) // This is the index
							iterator.SetContent(hashKeys.GetContent())
							iterator.SetLength(len(hashKeys.GetContent()))
							iterator.Set(HasNext,
								p.NewFunction(context, isBuiltIn, iterator.SymbolTable(),
									NewBuiltInClassFunction(iterator,
										0,
										func(funcSelf Value, _ ...Value) (Value, *Object) {
											if funcSelf.GetInteger() < int64(funcSelf.GetLength()) {
												return p.GetTrue(), nil
											}
											return p.GetFalse(), nil
										},
									),
								),
							)
							iterator.Set(Next,
								p.NewFunction(context, isBuiltIn, iterator.SymbolTable(),
									NewBuiltInClassFunction(iterator,
										0,
										func(funcSelf Value, _ ...Value) (Value, *Object) {
											value := funcSelf.GetContent()[int(funcSelf.GetInteger())]
											funcSelf.SetInteger(funcSelf.GetInteger() + 1)
											return value, nil
										},
									),
								),
							)
							return iterator, nil
						},
					),
				)
			},
		)

		object.SetOnDemandSymbol(ToString,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							result := "{"
							var (
								keyString     Value
								valueToString Value
								valueString   Value
								callError     *Object
							)
							for _, keyValues := range self.GetKeyValues() {
								for _, keyValue := range keyValues {
									keyToString, getError := keyValue.Key.Get(ToString)
									if getError != nil {
										return nil, p.NewObjectWithNameNotFoundError(context, keyValue.Key.GetClass(p), ToString)
									}
									keyString, callError = p.CallFunction(context, keyToString, keyValue.Key.SymbolTable())
									if callError != nil {
										return nil, callError
									}
									result += keyString.GetString()
									valueToString, getError = keyValue.Value.Get(ToString)
									if getError != nil {
										return nil, p.NewObjectWithNameNotFoundError(context, keyValue.Value.GetClass(p), ToString)
									}
									valueString, callError = p.CallFunction(context, valueToString, keyValue.Value.SymbolTable())
									if callError != nil {
										return nil, callError
									}
									result += ": " + valueString.GetString() + ", "
								}
							}
							if len(result) > 1 {
								result = result[:len(result)-2]
							}
							return p.NewString(context, false, context.PeekSymbolTable(), result+"}"), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToBool,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							if self.GetLength() > 0 {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToArray,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							var keys []Value
							for _, keyValues := range self.GetKeyValues() {
								for _, keyValue := range keyValues {
									keys = append(keys, keyValue.Key)
								}
							}
							return p.NewArray(context, false, context.PeekSymbolTable(), keys), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToTuple,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							var keys []Value
							for _, keyValues := range self.GetKeyValues() {
								for _, keyValue := range keyValues {
									keys = append(keys, keyValue.Key)
								}
							}
							return p.NewTuple(context, false, context.PeekSymbolTable(), keys), nil
						},
					),
				)
			},
		)
		return nil
	}
}
