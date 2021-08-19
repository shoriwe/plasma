package vm

type KeyValue struct {
	Key   *Value
	Value *Value
}

func (p *Plasma) NewHashTable(context *Context, isBuiltIn bool, parent *SymbolTable, entries map[int64][]*KeyValue, entriesLength int) *Value {
	hashTable := p.NewValue(context, isBuiltIn, HashName, nil, parent)
	hashTable.BuiltInTypeId = HashTableId
	hashTable.SetKeyValues(entries)
	hashTable.SetLength(entriesLength)
	p.HashTableInitialize(isBuiltIn)(context, hashTable)
	hashTable.SetOnDemandSymbol(Self,
		func() *Value {
			return hashTable
		},
	)
	return hashTable
}

func (p *Plasma) HashTableInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object *Value) *Value {
		object.SetOnDemandSymbol(Equals,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {

							rawRight := arguments[0]
							if !rawRight.IsTypeById(HashTableId) {
								return p.GetFalse(), true
							}
							right := rawRight
							if self.GetLength() != right.Length {
								return p.GetFalse(), true
							}
							rightIndex, getError := right.Get(Index)
							if getError != nil {
								return p.NewObjectWithNameNotFoundError(context, right.GetClass(p), Index), false
							}
							for key, leftValue := range self.GetKeyValues() {
								// Check if other has the key
								rightValue, ok := right.KeyValues[key]
								if !ok {
									return p.GetFalse(), true
								}
								// Check if the each entry one has the same length
								if len(leftValue) != len(rightValue) {
									return p.GetFalse(), true
								}
								// Start comparing the entries
								for _, entry := range leftValue {
									_, success := p.CallFunction(context, rightIndex, entry.Key)
									if !success {
										return p.GetFalse(), true
									}
								}
							}
							return p.GetTrue(), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightEquals,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {

							rawLeft := arguments[0]
							if !rawLeft.IsTypeById(HashTableId) {
								return p.GetFalse(), true
							}
							left := rawLeft
							if self.GetLength() != left.Length {
								return p.GetFalse(), true
							}
							leftIndex, getError := left.Get(Index)
							if getError != nil {
								return p.NewObjectWithNameNotFoundError(context, left.GetClass(p), Index), false
							}
							for key, leftValue := range left.KeyValues {
								// Check if other has the key
								rightValue, ok := self.GetKeyValues()[key]
								if !ok {
									return p.GetFalse(), true
								}
								// Check if the each entry one has the same length
								if len(leftValue) != len(rightValue) {
									return p.GetFalse(), true
								}
								// Start comparing the entries
								for _, entry := range leftValue {
									_, success := p.CallFunction(context, leftIndex, entry.Key)
									if !success {
										return p.GetFalse(), true
									}
								}
							}
							return p.GetTrue(), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(NotEquals,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							rawRight := arguments[0]
							if !rawRight.IsTypeById(HashTableId) {
								return p.GetTrue(), true
							}
							right := rawRight
							if self.GetLength() != right.Length {
								return p.GetTrue(), true
							}
							rightIndex, getError := right.Get(Index)
							if getError != nil {
								return p.NewObjectWithNameNotFoundError(context, right.GetClass(p), Index), false
							}
							for key, leftValue := range self.GetKeyValues() {
								// Check if other has the key
								rightValue, ok := right.KeyValues[key]
								if !ok {
									return p.GetTrue(), true
								}
								// Check if the each entry one has the same length
								if len(leftValue) != len(rightValue) {
									return p.GetTrue(), true
								}
								// Start comparing the entries
								for _, entry := range leftValue {
									_, success := p.CallFunction(context, rightIndex, entry.Key)
									if !success {
										return p.GetTrue(), true
									}
								}
							}
							return p.GetFalse(), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightNotEquals,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							rawLeft := arguments[0]
							if !rawLeft.IsTypeById(HashTableId) {
								return p.GetTrue(), true
							}
							left := rawLeft
							if self.GetLength() != left.Length {
								return p.GetTrue(), true
							}
							leftIndex, getError := left.Get(Index)
							if getError != nil {
								return p.NewObjectWithNameNotFoundError(context, left.GetClass(p), Index), false
							}
							for key, leftValue := range left.KeyValues {
								// Check if other has the key
								rightValue, ok := self.GetKeyValues()[key]
								if !ok {
									return p.GetTrue(), true
								}
								// Check if the each entry one has the same length
								if len(leftValue) != len(rightValue) {
									return p.GetTrue(), true
								}
								// Start comparing the entries
								for _, entry := range leftValue {
									_, success := p.CallFunction(context, leftIndex, entry.Key)
									if !success {
										return p.GetTrue(), true
									}
								}
							}
							return p.GetFalse(), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Contains,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							value := arguments[0]
							valueHashFunc, getError := value.Get(Hash)
							if getError != nil {
								return p.NewObjectWithNameNotFoundError(context, value.GetClass(p), Hash), false
							}
							valueHashObject, success := p.CallFunction(context, valueHashFunc)
							if !success {
								return valueHashObject, false
							}
							if !valueHashObject.IsTypeById(IntegerId) {
								return p.NewInvalidTypeError(context, valueHashObject.TypeName(), IntegerName), false
							}
							valueHash := valueHashObject.GetInteger()
							entries, found := self.GetKeyValues()[valueHash]
							if !found {
								return p.GetFalse(), true
							}
							var valueEquals *Value
							valueEquals, getError = value.Get(RightEquals)
							if getError != nil {
								return p.NewObjectWithNameNotFoundError(context, value.GetClass(p), RightEquals), false
							}
							var comparisonResult *Value
							for _, entry := range entries {
								comparisonResult, success = p.CallFunction(context, valueEquals, entry.Key)
								if !success {
									return comparisonResult, false
								}
								comparisonResultBool, callError := p.QuickGetBool(context, comparisonResult)
								if callError != nil {
									return callError, false
								}
								if comparisonResultBool {
									return p.GetTrue(), true
								}
							}
							return p.GetFalse(), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightContains,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							value := arguments[0]
							valueHashFunc, getError := value.Get(Hash)
							if getError != nil {
								return p.NewObjectWithNameNotFoundError(context, value.GetClass(p), Hash), false
							}
							valueHashObject, success := p.CallFunction(context, valueHashFunc)
							if !success {
								return valueHashObject, false
							}
							if !valueHashObject.IsTypeById(IntegerId) {
								return p.NewInvalidTypeError(context, valueHashObject.TypeName(), IntegerName), false
							}
							valueHash := valueHashObject.GetInteger()
							entries, found := self.GetKeyValues()[valueHash]
							if !found {
								return p.GetFalse(), true
							}
							var valueEquals *Value
							valueEquals, getError = value.Get(Equals)
							if getError != nil {
								return p.NewObjectWithNameNotFoundError(context, value.GetClass(p), Equals), false
							}
							var comparisonResult *Value
							for _, entry := range entries {
								comparisonResult, success = p.CallFunction(context, valueEquals, entry.Key)
								if !success {
									return comparisonResult, false
								}
								comparisonResultBool, callError := p.QuickGetBool(context, comparisonResult)
								if callError != nil {
									return callError, false
								}
								if comparisonResultBool {
									return p.GetTrue(), true
								}
							}
							return p.GetFalse(), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Hash,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(_ *Value, _ ...*Value) (*Value, bool) {
							return p.NewUnhashableTypeError(context, p.ForceMasterGetAny(HashName)), false
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Copy,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(_ *Value, _ ...*Value) (*Value, bool) {
							return p.NewUnhashableTypeError(context, object.GetClass(p)), false
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Index,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							indexObject := arguments[0]
							indexObjectHash, getError := indexObject.Get(Hash)
							if getError != nil {
								return p.NewObjectWithNameNotFoundError(context, indexObject.GetClass(p), Hash), false
							}
							indexHash, success := p.CallFunction(context, indexObjectHash)
							if !success {
								return indexHash, false
							}
							if !indexHash.IsTypeById(IntegerId) {
								return p.NewInvalidTypeError(context, indexHash.TypeName(), IntegerName), false
							}
							keyValues, found := self.GetKeyValues()[indexHash.GetInteger()]
							if !found {
								return p.NewKeyNotFoundError(context, indexObject), false
							}
							var indexObjectEquals *Value
							indexObjectEquals, getError = indexObject.Get(Equals)
							var equals *Value
							for _, keyValue := range keyValues {
								equals, success = p.CallFunction(context, indexObjectEquals, keyValue.Key)
								if !success {
									return equals, false
								}
								equalsBool, callError := p.QuickGetBool(context, equals)
								if callError != nil {
									return callError, false
								}
								if equalsBool {
									return keyValue.Value, true
								}
							}
							return p.NewKeyNotFoundError(context, indexObject), false
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Assign,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 2,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							indexObject := arguments[0]
							newValue := arguments[1]
							indexObjectHash, getError := indexObject.Get(Hash)
							if getError != nil {
								return p.NewObjectWithNameNotFoundError(context, indexObject.GetClass(p), Hash), false
							}
							indexHash, success := p.CallFunction(context, indexObjectHash)
							if !success {
								return indexHash, false
							}
							if !indexHash.IsTypeById(IntegerId) {
								return p.NewInvalidTypeError(context, indexHash.TypeName(), IntegerName), false
							}
							keyValues, found := self.GetKeyValues()[indexHash.GetInteger()]
							if found {
								self.AddKeyValue(indexHash.GetInteger(), &KeyValue{
									Key:   indexObject,
									Value: newValue,
								})
								return p.GetNone(), true
							}
							var indexObjectEquals *Value
							indexObjectEquals, getError = indexObject.Get(Equals)
							var equals *Value
							for index, keyValue := range keyValues {
								equals, success = p.CallFunction(context, indexObjectEquals, keyValue.Key)
								if !success {
									return equals, false
								}
								equalsBool, callError := p.QuickGetBool(context, equals)
								if callError != nil {
									return callError, false
								}
								if equalsBool {
									self.GetKeyValues()[indexHash.GetInteger()][index].Value = newValue
									return p.GetNone(), true
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
							return p.GetNone(), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Iter,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							toTuple, getError := self.Get(ToTuple)
							if getError != nil {
								return p.NewObjectWithNameNotFoundError(context, self.GetClass(p), ToTuple), false
							}
							hashKeys, success := p.CallFunction(context, toTuple)
							if !success {
								return hashKeys, false
							}
							iterator := p.NewIterator(context, false, context.PeekSymbolTable())
							iterator.SetInteger(0) // This is the index
							iterator.SetContent(hashKeys.GetContent())
							iterator.SetLength(len(hashKeys.GetContent()))
							iterator.Set(HasNext,
								p.NewFunction(context, isBuiltIn, iterator.SymbolTable(),
									NewBuiltInClassFunction(iterator,
										0,
										func(funcSelf *Value, _ ...*Value) (*Value, bool) {
											return p.InterpretAsBool(funcSelf.GetInteger() < int64(funcSelf.GetLength())), true
										},
									),
								),
							)
							iterator.Set(Next,
								p.NewFunction(context, isBuiltIn, iterator.SymbolTable(),
									NewBuiltInClassFunction(iterator,
										0,
										func(funcSelf *Value, _ ...*Value) (*Value, bool) {
											value := funcSelf.GetContent()[int(funcSelf.GetInteger())]
											funcSelf.SetInteger(funcSelf.GetInteger() + 1)
											return value, true
										},
									),
								),
							)
							return iterator, true
						},
					),
				)
			},
		)

		object.SetOnDemandSymbol(ToString,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							result := "{"
							var (
								valueToString *Value
								valueString   *Value
							)
							for _, keyValues := range self.GetKeyValues() {
								for _, keyValue := range keyValues {
									keyToString, getError := keyValue.Key.Get(ToString)
									if getError != nil {
										return p.NewObjectWithNameNotFoundError(context, keyValue.Key.GetClass(p), ToString), false
									}
									keyString, success := p.CallFunction(context, keyToString)
									if !success {
										return keyString, false
									}
									result += keyString.GetString()
									valueToString, getError = keyValue.Value.Get(ToString)
									if getError != nil {
										return p.NewObjectWithNameNotFoundError(context, keyValue.Value.GetClass(p), ToString), false
									}
									valueString, success = p.CallFunction(context, valueToString)
									if !success {
										return valueString, false
									}
									result += ": " + valueString.GetString() + ", "
								}
							}
							if len(result) > 1 {
								result = result[:len(result)-2]
							}
							return p.NewString(context, false, context.PeekSymbolTable(), result+"}"), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToBool,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.InterpretAsBool(self.GetLength() > 0), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToArray,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							var keys []*Value
							for _, keyValues := range self.GetKeyValues() {
								for _, keyValue := range keyValues {
									keys = append(keys, keyValue.Key)
								}
							}
							return p.NewArray(context, false, context.PeekSymbolTable(), keys), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToTuple,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							var keys []*Value
							for _, keyValues := range self.GetKeyValues() {
								for _, keyValue := range keyValues {
									keys = append(keys, keyValue.Key)
								}
							}
							return p.NewTuple(context, false, context.PeekSymbolTable(), keys), true
						},
					),
				)
			},
		)
		return nil
	}
}
