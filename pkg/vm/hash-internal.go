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

type Hash struct {
	mutex       *sync.Mutex
	internalMap map[string]*Value
}

func (h *Hash) Size() int64 {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	return int64(len(h.internalMap))
}

func (h *Hash) Set(key, value *Value) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	switch key.typeId {
	case StringId:
		h.internalMap[createHashString(createHashString(key.GetBytes()))] = value
	case BytesId:
		h.internalMap[createHashString(key.GetBytes())] = value
	case BoolId:
		h.internalMap[createHashString(key.GetBool())] = value
	case IntId:
		h.internalMap[createHashString(key.GetInt64())] = value
	case FloatId:
		h.internalMap[createHashString(key.GetFloat64())] = value
	}
	return NotHashable
}

func (h *Hash) Get(key *Value) (*Value, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	switch key.typeId {
	case StringId:
		return h.internalMap[createHashString(string(key.GetBytes()))], nil
	case BytesId:
		return h.internalMap[createHashString(key.GetBytes())], nil
	case BoolId:
		return h.internalMap[createHashString(key.GetBool())], nil
	case IntId:
		return h.internalMap[createHashString(key.GetInt64())], nil
	case FloatId:
		return h.internalMap[createHashString(key.GetFloat64())], nil
	}
	return nil, NotHashable
}

func (h *Hash) Del(key *Value) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	switch key.typeId {
	case StringId:
		delete(h.internalMap, createHashString(key.GetBytes()))
	case BytesId:
		delete(h.internalMap, createHashString(key.GetBytes()))
	case BoolId:
		delete(h.internalMap, createHashString(key.GetBool()))
	case IntId:
		delete(h.internalMap, createHashString(key.GetInt64()))
	case FloatId:
		delete(h.internalMap, createHashString(key.GetFloat64()))
	}
	return NotHashable
}

func (plasma *Plasma) NewInternalHash() *Hash {
	return &Hash{
		mutex:       &sync.Mutex{},
		internalMap: map[string]*Value{},
	}
}
