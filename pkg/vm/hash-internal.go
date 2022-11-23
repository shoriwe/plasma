package vm

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	NotHashable = fmt.Errorf("not hashable")
)

func createHashString(a any) string {
	return fmt.Sprintf("%s -- %v", reflect.TypeOf(a), a)
}

type (
	HashKeyValue struct {
		Key   *Value
		Value *Value
	}
	Hash struct {
		mutex       *sync.Mutex
		internalMap map[string]HashKeyValue
	}
)

/*
Size Returns the internal size of the hash
*/
func (h *Hash) Size() int64 {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	return int64(len(h.internalMap))
}

/*
Set sets a key value pair
*/
func (h *Hash) Set(key, value *Value) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	var keyString string
	switch key.TypeId() {
	case StringId, BytesId:
		keyString = createHashString(string(key.GetBytes()))
	case BoolId:
		keyString = createHashString(key.GetBool())
	case IntId:
		keyString = createHashString(key.GetInt64())
	case FloatId:
		keyString = createHashString(key.GetFloat64())
	default:
		return NotHashable
	}
	h.internalMap[keyString] = HashKeyValue{key, value}
	return nil
}

/*
Get retrieves a value based on the key
*/
func (h *Hash) Get(key *Value) (*Value, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	var keyString string
	switch key.TypeId() {
	case StringId, BytesId:
		keyString = createHashString(string(key.GetBytes()))
	case BoolId:
		keyString = createHashString(key.GetBool())
	case IntId:
		keyString = createHashString(key.GetInt64())
	case FloatId:
		keyString = createHashString(key.GetFloat64())
	default:
		return nil, NotHashable
	}
	return h.internalMap[keyString].Value, nil
}

/*
Del deletes a key from the hash
*/
func (h *Hash) Del(key *Value) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	var keyString string
	switch key.TypeId() {
	case StringId, BytesId:
		keyString = createHashString(string(key.GetBytes()))
	case BoolId:
		keyString = createHashString(key.GetBool())
	case IntId:
		keyString = createHashString(key.GetInt64())
	case FloatId:
		keyString = createHashString(key.GetFloat64())
	default:
		return NotHashable
	}
	delete(h.internalMap, keyString)
	return nil
}

/*
Copy creates a copy of the hash
*/
func (h *Hash) Copy() *Hash {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	result := &Hash{
		mutex:       &sync.Mutex{},
		internalMap: make(map[string]HashKeyValue, len(h.internalMap)),
	}
	for key, value := range h.internalMap {
		result.internalMap[key] = value
	}
	return result
}

/*
In verifies the key is inside the hash
*/
func (h *Hash) In(key *Value) (bool, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	var found bool
	switch key.TypeId() {
	case StringId:
		_, found = h.internalMap[createHashString(string(key.GetBytes()))]
	case BytesId:
		_, found = h.internalMap[createHashString(key.GetBytes())]
	case BoolId:
		_, found = h.internalMap[createHashString(key.GetBool())]
	case IntId:
		_, found = h.internalMap[createHashString(key.GetInt64())]
	case FloatId:
		_, found = h.internalMap[createHashString(key.GetFloat64())]
	default:
		return false, NotHashable
	}
	return found, nil
}

/*
NewInternalHash Creates a new Hash object to handle hashing operations
*/
func (plasma *Plasma) NewInternalHash() *Hash {
	return &Hash{
		mutex:       &sync.Mutex{},
		internalMap: map[string]HashKeyValue{},
	}
}
