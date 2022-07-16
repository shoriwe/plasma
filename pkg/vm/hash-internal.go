package vm

import (
	"fmt"
	"sync"
)

var (
	NotHashable = fmt.Errorf("not hashable")
)

type Hash struct {
	mutex       *sync.Mutex
	internalMap map[comparable]*Value
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
		h.internalMap[string(key.GetBytes())] = value
	case BytesId:
		h.internalMap[key.GetBytes()] = value
	case BoolId:
		h.internalMap[key.GetBool()] = value
	case IntId:
		h.internalMap[key.GetInt64()] = value
	case FloatId:
		h.internalMap[key.GetFloat64()] = value
	}
	return NotHashable
}

func (h *Hash) Get(key *Value) (*Value, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	switch key.typeId {
	case StringId:
		return h.internalMap[string(key.GetBytes())], nil
	case BytesId:
		return h.internalMap[key.GetBytes()], nil
	case BoolId:
		return h.internalMap[key.GetBool()], nil
	case IntId:
		return h.internalMap[key.GetInt64()], nil
	case FloatId:
		return h.internalMap[key.GetFloat64()], nil
	}
	return nil, NotHashable
}

func (h *Hash) Del(key *Value) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	switch key.typeId {
	case StringId:
		delete(h.internalMap, string(key.GetBytes()))
	case BytesId:
		delete(h.internalMap, key.GetBytes())
	case BoolId:
		delete(h.internalMap, key.GetBool())
	case IntId:
		delete(h.internalMap, key.GetInt64())
	case FloatId:
		delete(h.internalMap, key.GetFloat64())
	}
	return NotHashable
}

func (plasma *Plasma) NewInternalHash() *Hash {
	return &Hash{
		mutex:       &sync.Mutex{},
		internalMap: map[comparable]*Value{},
	}
}
