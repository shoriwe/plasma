package vm

import (
	"bytes"
	"encoding/binary"
	"github.com/shoriwe/gplasma/pkg/tools"
	"math/big"
)

type Bytes struct {
	*Object
}

func (p *Plasma) NewBytes(isBuiltIn bool, parent *SymbolTable, content []uint8) Value {
	bytes_ := &Bytes{
		Object: p.NewObject(false, BytesName, nil, parent),
	}
	bytes_.SetBytes(content)
	bytes_.SetLength(len(content))
	p.BytesInitialize(isBuiltIn)(bytes_)
	bytes_.Set(Self, bytes_)
	return bytes_
}

func (p *Plasma) BytesInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(object Value) *Object {
		object.Set(Add,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						if _, ok := right.(*Bytes); !ok {
							return nil, p.NewInvalidTypeError(right.TypeName(), BytesName)
						}
						var newContent []uint8
						copy(newContent, self.GetBytes())
						newContent = append(newContent, right.GetBytes()...)
						return p.NewBytes(false, p.PeekSymbolTable(), newContent), nil
					},
				),
			),
		)
		object.Set(RightAdd,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						if _, ok := left.(*Bytes); !ok {
							return nil, p.NewInvalidTypeError(left.TypeName(), BytesName)
						}
						var newContent []uint8
						copy(newContent, left.GetBytes())
						newContent = append(newContent, self.GetBytes()...)
						return p.NewBytes(false, p.PeekSymbolTable(), newContent), nil
					},
				),
			),
		)
		object.Set(Mul,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						if _, ok := right.(*Integer); !ok {
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName)
						}
						return p.NewBytes(false, p.PeekSymbolTable(), bytes.Repeat(self.GetBytes(), int(right.GetInteger().Int64()))), nil
					},
				),
			),
		)
		object.Set(RightMul,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						if _, ok := left.(*Integer); !ok {
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName)
						}
						return p.NewBytes(false, p.PeekSymbolTable(), bytes.Repeat(left.GetBytes(), int(self.GetInteger().Int64()))), nil
					},
				),
			),
		)
		object.Set(Equals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						if _, ok := right.(*Bytes); !ok {
							return p.NewBool(false, p.PeekSymbolTable(), false), nil
						}
						if self.GetLength() != right.GetLength() {
							return p.NewBool(false, p.PeekSymbolTable(), false), nil
						}
						return p.NewBool(false, p.PeekSymbolTable(), bytes.Compare(self.GetBytes(), right.GetBytes()) == 0), nil
					},
				),
			),
		)
		object.Set(RightEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						if _, ok := left.(*Bytes); !ok {
							return p.NewBool(false, p.PeekSymbolTable(), false), nil
						}
						if left.GetLength() != self.GetLength() {
							return p.NewBool(false, p.PeekSymbolTable(), false), nil
						}
						return p.NewBool(false, p.PeekSymbolTable(), bytes.Compare(left.GetBytes(), self.GetBytes()) == 0), nil
					},
				),
			),
		)
		object.Set(NotEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						if _, ok := right.(*Bytes); !ok {
							return p.NewBool(false, p.PeekSymbolTable(), false), nil
						}
						if self.GetLength() != right.GetLength() {
							return p.NewBool(false, p.PeekSymbolTable(), false), nil
						}
						return p.NewBool(false, p.PeekSymbolTable(), bytes.Compare(self.GetBytes(), right.GetBytes()) == 0), nil
					},
				),
			),
		)
		object.Set(RightNotEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						if _, ok := left.(*Bytes); !ok {
							return p.NewBool(false, p.PeekSymbolTable(), false), nil
						}
						if left.GetLength() != self.GetLength() {
							return p.NewBool(false, p.PeekSymbolTable(), false), nil
						}
						return p.NewBool(false, p.PeekSymbolTable(), bytes.Compare(left.GetBytes(), self.GetBytes()) == 0), nil
					},
				),
			),
		)
		object.Set(Hash,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						selfHash := p.HashBytes(append(self.GetBytes(), []byte("Bytes")...))
						return p.NewInteger(false, p.PeekSymbolTable(), big.NewInt(selfHash)), nil
					},
				),
			),
		)
		object.Set(Copy,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						var newBytes []uint8
						copy(newBytes, self.GetBytes())
						return p.NewBytes(false, p.PeekSymbolTable(), newBytes), nil
					},
				),
			),
		)
		object.Set(Index,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						indexObject := arguments[0]
						var ok bool
						if _, ok = indexObject.(*Integer); ok {
							index, calcError := tools.CalcIndex(indexObject.GetInteger().Int64(), self.GetLength())
							if calcError != nil {
								return nil, p.NewIndexOutOfRange(self.GetLength(), indexObject.GetInteger().Int64())
							}
							return p.NewInteger(false, p.PeekSymbolTable(), big.NewInt(int64(self.GetBytes()[index]))), nil
						} else if _, ok = indexObject.(*Tuple); ok {
							if len(indexObject.GetContent()) != 2 {
								return nil, p.NewInvalidNumberOfArgumentsError(len(indexObject.GetContent()), 2)
							}
							startIndex, calcError := tools.CalcIndex(indexObject.GetContent()[0].GetInteger().Int64(), self.GetLength())
							if calcError != nil {
								return nil, p.NewIndexOutOfRange(self.GetLength(), indexObject.GetContent()[0].GetInteger().Int64())
							}
							var targetIndex int
							targetIndex, calcError = tools.CalcIndex(indexObject.GetContent()[1].GetInteger().Int64(), self.GetLength())
							if calcError != nil {
								return nil, p.NewIndexOutOfRange(self.GetLength(), indexObject.GetContent()[1].GetInteger().Int64())
							}
							return p.NewBytes(false, p.PeekSymbolTable(), self.GetBytes()[startIndex:targetIndex]), nil
						} else {
							return nil, p.NewInvalidTypeError(indexObject.TypeName(), IntegerName, TupleName)
						}
					},
				),
			),
		)
		object.Set(Iter,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						iterator := p.NewIterator(false, p.PeekSymbolTable())
						iterator.SetInteger(big.NewInt(0)) // This is the index
						iterator.SetBytes(self.GetBytes())
						iterator.SetLength(self.GetLength())
						iterator.Set(HasNext,
							p.NewFunction(isBuiltIn, iterator.SymbolTable(),
								NewBuiltInClassFunction(iterator,
									0,
									func(funcSelf Value, _ ...Value) (Value, *Object) {
										if int(funcSelf.GetInteger().Int64()) < funcSelf.GetLength() {
											return p.NewBool(false, p.PeekSymbolTable(), true), nil
										}
										return p.NewBool(false, p.PeekSymbolTable(), false), nil
									},
								),
							),
						)
						iterator.Set(Next,
							p.NewFunction(isBuiltIn, iterator.SymbolTable(),
								NewBuiltInClassFunction(iterator,
									0,
									func(funcSelf Value, _ ...Value) (Value, *Object) {
										char := funcSelf.GetBytes()[int(funcSelf.GetInteger().Int64())]
										funcSelf.SetInteger(new(big.Int).Add(funcSelf.GetInteger(), big.NewInt(1)))
										return p.NewInteger(false, p.PeekSymbolTable(), big.NewInt(int64(char))), nil
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
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewInteger(false, p.PeekSymbolTable(),
							big.NewInt(int64(binary.BigEndian.Uint32(self.GetBytes()))),
						), nil
					},
				),
			),
		)
		object.Set(ToString,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewString(false, p.PeekSymbolTable(), string(self.GetBytes())), nil
					},
				),
			),
		)
		object.Set(ToBool,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewBool(false, p.PeekSymbolTable(), self.GetLength() != 0), nil
					},
				),
			),
		)
		object.Set(ToArray,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						var newContent []Value
						for _, byte_ := range self.GetBytes() {
							newContent = append(newContent,
								p.NewInteger(false, p.PeekSymbolTable(),
									big.NewInt(int64(byte_)),
								),
							)
						}
						return p.NewArray(false, p.PeekSymbolTable(), newContent), nil
					},
				),
			),
		)
		object.Set(ToTuple,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						var newContent []Value
						for _, byte_ := range self.GetBytes() {
							newContent = append(newContent,
								p.NewInteger(false, p.PeekSymbolTable(),
									big.NewInt(int64(byte_)),
								),
							)
						}
						return p.NewTuple(false, p.PeekSymbolTable(), newContent), nil
					},
				),
			),
		)
		return nil
	}
}
