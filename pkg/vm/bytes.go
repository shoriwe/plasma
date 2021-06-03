package vm

import (
	"bytes"
	"encoding/binary"
	"github.com/shoriwe/gruby/pkg/errors"
)

type Bytes struct {
	*Object
}

func (p *Plasma) NewBytes(parent *SymbolTable, content []uint8) IObject {
	bytes_ := &Bytes{
		Object: p.NewObject(BytesName, nil, parent),
	}
	bytes_.SetBytes(content)
	bytes_.SetLength(len(content))
	p.BytesInitialize(bytes_)
	return bytes_
}

func (p *Plasma) BytesInitialize(object IObject) *errors.Error {
	object.Set(Add,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					if _, ok := right.(*Bytes); !ok {
						return nil, errors.NewTypeError(right.TypeName(), BytesName)
					}
					var newContent []uint8
					copy(newContent, self.GetBytes())
					newContent = append(newContent, right.GetBytes()...)
					return p.NewBytes(p.PeekSymbolTable(), newContent), nil
				},
			),
		),
	)
	object.Set(RightAdd,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					if _, ok := left.(*Bytes); !ok {
						return nil, errors.NewTypeError(left.TypeName(), BytesName)
					}
					var newContent []uint8
					copy(newContent, left.GetBytes())
					newContent = append(newContent, self.GetBytes()...)
					return p.NewBytes(p.PeekSymbolTable(), newContent), nil
				},
			),
		),
	)
	object.Set(Mul,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					if _, ok := right.(*Integer); !ok {
						return nil, errors.NewTypeError(right.TypeName(), IntegerName)
					}
					return p.NewBytes(p.PeekSymbolTable(), bytes.Repeat(self.GetBytes(), int(right.GetInteger64()))), nil
				},
			),
		),
	)
	object.Set(RightMul,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					if _, ok := left.(*Integer); !ok {
						return nil, errors.NewTypeError(left.TypeName(), IntegerName)
					}
					return p.NewBytes(p.PeekSymbolTable(), bytes.Repeat(left.GetBytes(), int(self.GetInteger64()))), nil
				},
			),
		),
	)
	object.Set(Equals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					if _, ok := right.(*Bytes); !ok {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					if self.GetLength() != right.GetLength() {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					return p.NewBool(p.PeekSymbolTable(), bytes.Compare(self.GetBytes(), right.GetBytes()) == 0), nil
				},
			),
		),
	)
	object.Set(RightEquals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					if _, ok := left.(*Bytes); !ok {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					if left.GetLength() != self.GetLength() {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					return p.NewBool(p.PeekSymbolTable(), bytes.Compare(left.GetBytes(), self.GetBytes()) == 0), nil
				},
			),
		),
	)
	object.Set(NotEquals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					if _, ok := right.(*Bytes); !ok {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					if self.GetLength() != right.GetLength() {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					return p.NewBool(p.PeekSymbolTable(), bytes.Compare(self.GetBytes(), right.GetBytes()) == 0), nil
				},
			),
		),
	)
	object.Set(RightNotEquals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					if _, ok := left.(*Bytes); !ok {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					if left.GetLength() != self.GetLength() {
						return p.NewBool(p.PeekSymbolTable(), false), nil
					}
					return p.NewBool(p.PeekSymbolTable(), bytes.Compare(left.GetBytes(), self.GetBytes()) == 0), nil
				},
			),
		),
	)
	object.Set(Hash,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					selfHash, hashingError := p.HashBytes(append(self.GetBytes(), []byte("Bytes")...))
					if hashingError != nil {
						return nil, hashingError
					}
					return p.NewInteger(p.PeekSymbolTable(), selfHash), nil
				},
			),
		),
	)
	object.Set(Copy,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					var newBytes []uint8
					copy(newBytes, self.GetBytes())
					return p.NewBytes(p.PeekSymbolTable(), newBytes), nil
				},
			),
		),
	)
	object.Set(Index,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					indexObject := arguments[0]
					var ok bool
					if _, ok = indexObject.(*Integer); ok {
						index, calcError := CalcIndex(indexObject, self.GetLength())
						if calcError != nil {
							return nil, calcError
						}
						return p.NewInteger(p.PeekSymbolTable(), int64(self.GetBytes()[index])), nil
					} else if _, ok = indexObject.(*Tuple); ok {
						if len(indexObject.GetContent()) != 2 {
							return nil, errors.NewInvalidNumberOfArguments(len(indexObject.GetContent()), 2)
						}
						startIndex, calcError := CalcIndex(indexObject.GetContent()[0], self.GetLength())
						if calcError != nil {
							return nil, calcError
						}
						targetIndex, calcError := CalcIndex(indexObject.GetContent()[1], self.GetLength())
						if calcError != nil {
							return nil, calcError
						}
						return p.NewBytes(p.PeekSymbolTable(), self.GetBytes()[startIndex:targetIndex]), nil
					} else {
						return nil, errors.NewTypeError(indexObject.TypeName(), IntegerName, TupleName)
					}
				},
			),
		),
	)
	object.Set(Iter,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					iterator := p.NewIterator(p.PeekSymbolTable())
					iterator.SetInteger64(0) // This is the index
					iterator.SetBytes(self.GetBytes())
					iterator.SetLength(self.GetLength())
					iterator.Set(HasNext,
						p.NewFunction(iterator.SymbolTable(),
							NewBuiltInClassFunction(iterator,
								0,
								func(self IObject, _ ...IObject) (IObject, *errors.Error) {
									funcSelf, funcGetError := p.PeekSymbolTable().GetSelf(Self)
									if funcGetError != nil {
										return nil, funcGetError
									}
									if int(funcSelf.GetInteger64()) < funcSelf.GetLength() {
										return p.NewBool(p.PeekSymbolTable(), true), nil
									}
									return p.NewBool(p.PeekSymbolTable(), false), nil
								},
							),
						),
					)
					iterator.Set(Next,
						p.NewFunction(iterator.SymbolTable(),
							NewBuiltInClassFunction(iterator,
								0,
								func(self IObject, _ ...IObject) (IObject, *errors.Error) {
									funcSelf, funcGetError := p.PeekSymbolTable().GetSelf(Self)
									if funcGetError != nil {
										return nil, funcGetError
									}
									char := funcSelf.GetBytes()[int(funcSelf.GetInteger64())]
									funcSelf.SetInteger64(funcSelf.GetInteger64() + 1)
									return p.NewInteger(p.PeekSymbolTable(), int64(char)), nil
								},
							),
						),
					)
					return iterator, nil
				},
			),
		),
	)
	object.Set(ToInteger,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewInteger(p.PeekSymbolTable(), int64(binary.BigEndian.Uint32(self.GetBytes()))), nil
				},
			),
		),
	)
	object.Set(ToString,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewString(p.PeekSymbolTable(), string(self.GetBytes())), nil
				},
			),
		),
	)
	object.Set(ToBool,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewBool(p.PeekSymbolTable(), self.GetLength() != 0), nil
				},
			),
		),
	)
	object.Set(ToArray,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					var newContent []IObject
					for _, byte_ := range self.GetBytes() {
						newContent = append(newContent, p.NewInteger(p.PeekSymbolTable(), int64(byte_)))
					}
					return p.NewArray(p.PeekSymbolTable(), newContent), nil
				},
			),
		),
	)
	object.Set(ToTuple,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					var newContent []IObject
					for _, byte_ := range self.GetBytes() {
						newContent = append(newContent, p.NewInteger(p.PeekSymbolTable(), int64(byte_)))
					}
					return p.NewTuple(p.PeekSymbolTable(), newContent), nil
				},
			),
		),
	)
	return nil
}
