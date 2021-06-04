package vm

import "github.com/shoriwe/gruby/pkg/errors"

type KeyValue struct {
	Key   IObject
	Value IObject
}

type HashTable struct {
	*Object
}

func (p *Plasma) NewHashTable(parent *SymbolTable, entries map[int64][]*KeyValue, entriesLength int) *HashTable {
	hashTable := &HashTable{
		Object: p.NewObject(HashName, nil, parent),
	}
	hashTable.SetKeyValues(entries)
	hashTable.SetLength(entriesLength)
	p.HashTableInitialize(hashTable)
	return hashTable
}

func (p *Plasma) HashTableInitialize(object IObject) *errors.Error {
	object.Set(Equals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {

					rawRight := arguments[0]
					if _, ok := rawRight.(*HashTable); !ok {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					right := rawRight.(*HashTable)
					if self.GetLength() != right.Length {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					rightIndex, getError := right.Get(Index)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(Index)
					}
					if _, ok := rightIndex.(*Function); !ok {
						return nil, p.NewInvalidTypeError(rightIndex.TypeName(), FunctionName)
					}
					for key, leftValue := range self.GetKeyValues() {
						// Check if other has the key
						rightValue, ok := right.KeyValues[key]
						if !ok {
							return p.NewBool(p.PeekSymbolTable(), false), nil
						}
						// Check if the each entry one has the same length
						if len(leftValue) != len(rightValue) {
							return p.NewBool(p.PeekSymbolTable(), false), nil
						}
						// Start comparing the entries
						for _, entry := range leftValue {
							_, indexingError := p.CallFunction(rightIndex.(*Function), p.PeekSymbolTable(), entry.Key)
							if indexingError != nil {
								return p.NewBool(p.PeekSymbolTable(), false), nil
							}
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
				func(self IObject, arguments ...IObject) (IObject, *Object) {

					rawLeft := arguments[0]
					if _, ok := rawLeft.(*HashTable); !ok {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					left := rawLeft.(*HashTable)
					if self.GetLength() != left.Length {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					leftIndex, getError := left.Get(Index)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(Index)
					}
					if _, ok := leftIndex.(*Function); !ok {
						return nil, p.NewInvalidTypeError(leftIndex.TypeName(), FunctionName)
					}
					for key, leftValue := range left.KeyValues {
						// Check if other has the key
						rightValue, ok := self.GetKeyValues()[key]
						if !ok {
							return p.NewBool(p.PeekSymbolTable(), false), nil
						}
						// Check if the each entry one has the same length
						if len(leftValue) != len(rightValue) {
							return p.NewBool(p.PeekSymbolTable(), false), nil
						}
						// Start comparing the entries
						for _, entry := range leftValue {
							_, indexingError := p.CallFunction(leftIndex.(*Function), p.PeekSymbolTable(), entry.Key)
							if indexingError != nil {
								return p.NewBool(p.PeekSymbolTable(), false), nil
							}
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
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					rawRight := arguments[0]
					if _, ok := rawRight.(*HashTable); !ok {
						return p.NewBool(p.PeekSymbolTable(), true), nil
					}
					right := rawRight.(*HashTable)
					if self.GetLength() != right.Length {
						return p.NewBool(p.PeekSymbolTable(), true), nil
					}
					rightIndex, getError := right.Get(Index)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(Index)
					}
					if _, ok := rightIndex.(*Function); !ok {
						return nil, p.NewInvalidTypeError(rightIndex.TypeName(), FunctionName)
					}
					for key, leftValue := range self.GetKeyValues() {
						// Check if other has the key
						rightValue, ok := right.KeyValues[key]
						if !ok {
							return p.NewBool(p.PeekSymbolTable(), true), nil
						}
						// Check if the each entry one has the same length
						if len(leftValue) != len(rightValue) {
							return p.NewBool(p.PeekSymbolTable(), true), nil
						}
						// Start comparing the entries
						for _, entry := range leftValue {
							_, indexingError := p.CallFunction(rightIndex.(*Function), p.PeekSymbolTable(), entry.Key)
							if indexingError != nil {
								return p.NewBool(p.PeekSymbolTable(), true), nil
							}
						}
					}
					return p.NewBool(p.PeekSymbolTable(), false), nil
				},
			),
		),
	)
	object.Set(RightNotEquals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					rawLeft := arguments[0]
					if _, ok := rawLeft.(*HashTable); !ok {
						return p.NewBool(p.PeekSymbolTable(), true), nil
					}
					left := rawLeft.(*HashTable)
					if self.GetLength() != left.Length {
						return p.NewBool(p.PeekSymbolTable(), true), nil
					}
					leftIndex, getError := left.Get(Index)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(Index)
					}
					if _, ok := leftIndex.(*Function); !ok {
						return nil, p.NewInvalidTypeError(leftIndex.TypeName(), FunctionName)
					}
					for key, leftValue := range left.KeyValues {
						// Check if other has the key
						rightValue, ok := self.GetKeyValues()[key]
						if !ok {
							return p.NewBool(p.PeekSymbolTable(), true), nil
						}
						// Check if the each entry one has the same length
						if len(leftValue) != len(rightValue) {
							return p.NewBool(p.PeekSymbolTable(), true), nil
						}
						// Start comparing the entries
						for _, entry := range leftValue {
							_, indexingError := p.CallFunction(leftIndex.(*Function), p.PeekSymbolTable(), entry.Key)
							if indexingError != nil {
								return p.NewBool(p.PeekSymbolTable(), true), nil
							}
						}
					}
					return p.NewBool(p.PeekSymbolTable(), false), nil
				},
			),
		),
	)

	object.Set(Hash,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(_ IObject, _ ...IObject) (IObject, *Object) {
					panic("Implement me!!!")
				},
			),
		),
	)
	object.Set(Copy,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(_ IObject, _ ...IObject) (IObject, *Object) {
					return nil, p.NewUnhashableTypeError()
				},
			),
		),
	)
	object.Set(Index,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					indexObject := arguments[0]
					indexObjectHash, getError := indexObject.Get(Hash)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(Hash)
					}
					if _, ok := indexObjectHash.(*Function); !ok {
						return nil, p.NewInvalidTypeError(indexObjectHash.TypeName(), FunctionName)
					}
					indexHash, callError := p.CallFunction(indexObjectHash.(*Function), indexObject.SymbolTable())
					if callError != nil {
						return nil, callError
					}
					if _, ok := indexHash.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(indexHash.TypeName(), IntegerName)
					}
					keyValues, found := self.GetKeyValues()[indexHash.GetInteger64()]
					if !found {
						return nil, p.NewKeyNotFoundError(indexObject)
					}
					var indexObjectEquals IObject
					indexObjectEquals, getError = indexObject.Get(Equals)
					if _, ok := indexObjectEquals.(*Function); !ok {
						return nil, p.NewInvalidTypeError(indexObjectEquals.TypeName(), FunctionName)
					}
					var equals IObject
					for _, keyValue := range keyValues {
						equals, callError = p.CallFunction(indexObjectEquals.(*Function), indexObject.SymbolTable(), keyValue.Key)
						if callError != nil {
							return nil, callError
						}
						if _, ok := equals.(*Bool); !ok {
							return nil, p.NewInvalidTypeError(equals.TypeName(), BoolName)
						}
						if equals.GetBool() {
							return keyValue.Value, nil
						}
					}
					return nil, p.NewKeyNotFoundError(indexObject)
				},
			),
		),
	)
	object.Set(Assign,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 2,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					indexObject := arguments[0]
					newValue := arguments[1]
					indexObjectHash, getError := indexObject.Get(Hash)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(Hash)
					}
					if _, ok := indexObjectHash.(*Function); !ok {
						return nil, p.NewInvalidTypeError(indexObjectHash.TypeName(), FunctionName)
					}
					indexHash, callError := p.CallFunction(indexObjectHash.(*Function), indexObject.SymbolTable())
					if callError != nil {
						return nil, callError
					}
					if _, ok := indexHash.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(indexHash.TypeName(), IntegerName)
					}
					keyValues, found := self.GetKeyValues()[indexHash.GetInteger64()]
					if found {
						self.AddKeyValue(indexHash.GetInteger64(), &KeyValue{
							Key:   indexObject,
							Value: newValue,
						})
						return p.GetNone()
					}
					var indexObjectEquals IObject
					indexObjectEquals, getError = indexObject.Get(Equals)
					if _, ok := indexObjectEquals.(*Function); !ok {
						return nil, p.NewInvalidTypeError(indexObjectEquals.TypeName(), FunctionName)
					}
					var equals IObject
					for index, keyValue := range keyValues {
						equals, callError = p.CallFunction(indexObjectEquals.(*Function), indexObject.SymbolTable(), keyValue.Key)
						if callError != nil {
							return nil, callError
						}
						if _, ok := equals.(*Bool); !ok {
							return nil, p.NewInvalidTypeError(equals.TypeName(), BoolName)
						}
						if equals.GetBool() {
							self.GetKeyValues()[indexHash.GetInteger64()][index].Value = newValue
							return p.GetNone()
						}
					}
					self.IncreaseLength()
					self.GetKeyValues()[indexHash.GetInteger64()] = append(
						self.GetKeyValues()[indexHash.GetInteger64()],
						&KeyValue{
							Key:   indexObject,
							Value: newValue,
						},
					)
					return p.GetNone()
				},
			),
		),
	)
	object.Set(Iter,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					toTuple, getError := self.Get(ToTuple)
					if getError != nil {
						return nil, p.NewObjectWithNameNotFoundError(ToTuple)
					}
					if _, ok := toTuple.(*Function); !ok {
						return nil, p.NewInvalidTypeError(toTuple.TypeName(), FunctionName)
					}
					hashKeys, callError := p.CallFunction(toTuple.(*Function), self.SymbolTable())
					if callError != nil {
						return nil, callError
					}
					iterator := p.NewIterator(p.PeekSymbolTable())
					iterator.SetInteger64(0) // This is the index
					iterator.SetContent(hashKeys.GetContent())
					iterator.SetLength(len(hashKeys.GetContent()))
					iterator.Set(HasNext,
						p.NewFunction(iterator.SymbolTable(),
							NewBuiltInClassFunction(iterator,
								0,
								func(funcSelf IObject, _ ...IObject) (IObject, *Object) {
									return p.NewBool(p.PeekSymbolTable(), int(funcSelf.GetInteger64()) < funcSelf.GetLength()), nil
								},
							),
						),
					)
					iterator.Set(Next,
						p.NewFunction(iterator.SymbolTable(),
							NewBuiltInClassFunction(iterator,
								0,
								func(funcSelf IObject, _ ...IObject) (IObject, *Object) {
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
				func(self IObject, _ ...IObject) (IObject, *Object) {
					result := "{"
					var (
						keyString     IObject
						valueToString IObject
						valueString   IObject
						callError     *Object
					)
					for _, keyValues := range self.GetKeyValues() {
						for _, keyValue := range keyValues {
							keyToString, getError := keyValue.Key.Get(ToString)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(ToString)
							}
							keyString, callError = p.CallFunction(keyToString.(*Function), keyValue.Key.SymbolTable())
							if callError != nil {
								return nil, callError
							}
							result += keyString.GetString()
							valueToString, getError = keyValue.Value.Get(ToString)
							if getError != nil {
								return nil, p.NewObjectWithNameNotFoundError(ToString)
							}
							valueString, callError = p.CallFunction(valueToString.(*Function), keyValue.Value.SymbolTable())
							if callError != nil {
								return nil, callError
							}
							result += ": " + valueString.GetString() + ", "
						}
					}
					if len(result) > 1 {
						result = result[:len(result)-2]
					}
					return p.NewString(p.PeekSymbolTable(), result+"}"), nil
				},
			),
		),
	)
	object.Set(ToBool,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					if self.GetLength() > 0 {
						return p.NewBool(p.PeekSymbolTable(), true), nil
					}
					return p.NewBool(p.PeekSymbolTable(), false), nil
				},
			),
		),
	)
	object.Set(ToArray,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					var keys []IObject
					for _, keyValues := range self.GetKeyValues() {
						for _, keyValue := range keyValues {
							keys = append(keys, keyValue.Key)
						}
					}
					return p.NewArray(p.PeekSymbolTable(), keys), nil
				},
			),
		),
	)
	object.Set(ToTuple,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					var keys []IObject
					for _, keyValues := range self.GetKeyValues() {
						for _, keyValue := range keyValues {
							keys = append(keys, keyValue.Key)
						}
					}
					return p.NewTuple(p.PeekSymbolTable(), keys), nil
				},
			),
		),
	)
	return nil
}
