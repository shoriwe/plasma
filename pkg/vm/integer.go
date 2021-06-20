package vm

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/tools"
	"math"
	"math/big"
)

type Integer struct {
	*Object
}

func (p *Plasma) NewInteger(isBuiltIn bool, parentSymbols *SymbolTable, value *big.Int) *Integer {
	integer := &Integer{
		p.NewObject(isBuiltIn, IntegerName, nil, parentSymbols),
	}
	integer.SetInteger(value)
	p.IntegerInitialize(isBuiltIn)(integer)
	integer.Set(Self, integer)
	return integer
}

func (p *Plasma) IntegerInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(object Value) *Object {
		object.Set(NegBits,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewInteger(false, p.PeekSymbolTable(), new(big.Int).Not(self.GetInteger())), nil
					},
				),
			),
		)
		object.Set(Negative,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewInteger(false, p.PeekSymbolTable(), new(big.Int).Neg(self.GetInteger())), nil
					},
				),
			),
		)
		object.Set(Add,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							return p.NewInteger(false, p.PeekSymbolTable(),
								new(big.Int).Add(self.GetInteger(), right.GetInteger()),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Add(new(big.Float).SetInt(self.GetInteger()), right.GetFloat()),
							), nil
						default:
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(RightAdd,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						switch left.(type) {
						case *Integer:
							return p.NewInteger(false, p.PeekSymbolTable(),
								new(big.Int).Add(left.GetInteger(), self.GetInteger()),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Add(left.GetFloat(), new(big.Float).SetInt(self.GetInteger())),
							), nil
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(Sub,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							return p.NewInteger(false, p.PeekSymbolTable(),
								new(big.Int).Sub(self.GetInteger(), right.GetInteger()),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Sub(new(big.Float).SetInt(self.GetInteger()), right.GetFloat()),
							), nil
						default:
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(RightSub,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						switch left.(type) {
						case *Integer:
							return p.NewInteger(false, p.PeekSymbolTable(),
								new(big.Int).Sub(left.GetInteger(), self.GetInteger()),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Sub(left.GetFloat(), new(big.Float).SetInt(self.GetInteger())),
							), nil
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(Mul,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							return p.NewInteger(false, p.PeekSymbolTable(),
								new(big.Int).Mul(self.GetInteger(), right.GetInteger()),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Mul(new(big.Float).SetInt(self.GetInteger()), right.GetFloat()),
							), nil
						case *String:
							return p.NewString(false, p.PeekSymbolTable(), tools.Repeat(right.GetString(), self.GetInteger())), nil
						case *Tuple:
							content, repetitionError := p.Repeat(right.GetContent(), self.GetInteger())
							if repetitionError != nil {
								return nil, repetitionError
							}
							return p.NewTuple(false, p.PeekSymbolTable(), content), nil
						case *Array:
							content, repetitionError := p.Repeat(right.GetContent(), self.GetInteger())
							if repetitionError != nil {
								return nil, repetitionError
							}
							return p.NewArray(false, p.PeekSymbolTable(), content), nil
						default:
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName)
						}
					},
				),
			),
		)
		object.Set(RightMul,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						switch left.(type) {
						case *Integer:
							return p.NewInteger(false, p.PeekSymbolTable(),
								new(big.Int).Mul(left.GetInteger(), self.GetInteger()),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Mul(left.GetFloat(), new(big.Float).SetInt(self.GetInteger())),
							), nil
						case *String:
							return p.NewString(false, p.PeekSymbolTable(), tools.Repeat(left.GetString(), self.GetInteger())), nil
						case *Tuple:
							content, repetitionError := p.Repeat(left.GetContent(), self.GetInteger())
							if repetitionError != nil {
								return nil, repetitionError
							}
							return p.NewTuple(false, p.PeekSymbolTable(), content), nil
						case *Array:
							content, repetitionError := p.Repeat(left.GetContent(), self.GetInteger())
							if repetitionError != nil {
								return nil, repetitionError
							}
							return p.NewArray(false, p.PeekSymbolTable(), content), nil
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName)
						}
					},
				),
			),
		)
		object.Set(Div,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Quo(new(big.Float).SetInt(self.GetInteger()), new(big.Float).SetInt(right.GetInteger())),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Quo(new(big.Float).SetInt(self.GetInteger()), right.GetFloat()),
							), nil
						default:
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(RightDiv,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						switch left.(type) {
						case *Integer:
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Quo(new(big.Float).SetInt(left.GetInteger()), new(big.Float).SetInt(self.GetInteger())),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Quo(left.GetFloat(), new(big.Float).SetInt(self.GetInteger())),
							), nil
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(FloorDiv,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							return p.NewInteger(false, p.PeekSymbolTable(),
								new(big.Int).Quo(self.GetInteger(), right.GetInteger()),
							), nil
						case *Float:
							rightHandSide, _ := right.GetFloat().Int(nil)
							return p.NewInteger(false, p.PeekSymbolTable(),
								new(big.Int).Quo(self.GetInteger(), rightHandSide),
							), nil
						default:
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(RightFloorDiv,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						switch left.(type) {
						case *Integer:
							return p.NewInteger(false, p.PeekSymbolTable(),
								new(big.Int).Quo(left.GetInteger(), self.GetInteger()),
							), nil
						case *Float:
							leftHandSide, _ := left.GetFloat().Int(nil)
							return p.NewInteger(false, p.PeekSymbolTable(),
								new(big.Int).Quo(leftHandSide, self.GetInteger()),
							), nil
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(Mod,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							return p.NewInteger(false, p.PeekSymbolTable(),
								new(big.Int).Mod(self.GetInteger(), right.GetInteger()),
							), nil
						case *Float:
							rightHandSide, _ := right.GetFloat().Int(nil)
							return p.NewInteger(false, p.PeekSymbolTable(),
								new(big.Int).Mod(self.GetInteger(), rightHandSide),
							), nil
						default:
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(RightMod,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						switch left.(type) {
						case *Integer:
							return p.NewInteger(false, p.PeekSymbolTable(),
								new(big.Int).Mod(left.GetInteger(), self.GetInteger()),
							), nil
						case *Float:
							leftHandSide, _ := left.GetFloat().Int(nil)
							return p.NewInteger(false, p.PeekSymbolTable(),
								new(big.Int).Mod(leftHandSide, self.GetInteger()),
							), nil
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(Pow,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							return p.NewFloat(false, p.PeekSymbolTable(),
								big.NewFloat(math.Pow(float64(self.GetInteger().Int64()), float64(right.GetInteger().Int64()))),
							), nil
						case *Float:
							rightHandSide, _ := right.GetFloat().Float64()
							return p.NewFloat(false, p.PeekSymbolTable(),
								big.NewFloat(math.Pow(float64(self.GetInteger().Int64()), rightHandSide)),
							), nil
						default:
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(RightPow,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						switch left.(type) {
						case *Integer:
							return p.NewFloat(false, p.PeekSymbolTable(),
								big.NewFloat(
									math.Pow(
										float64(left.GetInteger().Int64()),
										float64(self.GetInteger().Int64(),
										),
									),
								),
							), nil
						case *Float:
							leftHandSide, _ := left.GetFloat().Float64()
							return p.NewFloat(false, p.PeekSymbolTable(),
								big.NewFloat(
									math.Pow(
										leftHandSide,
										float64(self.GetInteger().Int64(),
										),
									),
								),
							), nil
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(BitXor,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						if _, ok := right.(*Integer); !ok {
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName)
						}
						return p.NewInteger(false, p.PeekSymbolTable(),
							new(big.Int).Xor(self.GetInteger(), right.GetInteger()),
						), nil
					},
				),
			),
		)
		object.Set(RightBitXor,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						if _, ok := left.(*Integer); !ok {
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName)
						}
						return p.NewInteger(false, p.PeekSymbolTable(),
							new(big.Int).Xor(left.GetInteger(), self.GetInteger()),
						), nil
					},
				),
			),
		)
		object.Set(BitAnd,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						if _, ok := right.(*Integer); !ok {
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName)
						}
						return p.NewInteger(false, p.PeekSymbolTable(),
							new(big.Int).And(self.GetInteger(), right.GetInteger()),
						), nil
					},
				),
			),
		)
		object.Set(RightBitAnd,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						if _, ok := left.(*Integer); !ok {
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName)
						}
						return p.NewInteger(false, p.PeekSymbolTable(),
							new(big.Int).And(left.GetInteger(), self.GetInteger()),
						), nil
					},
				),
			),
		)
		object.Set(BitOr,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						if _, ok := right.(*Integer); !ok {
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName)
						}
						return p.NewInteger(false, p.PeekSymbolTable(),
							new(big.Int).Or(self.GetInteger(), right.GetInteger()),
						), nil
					},
				),
			),
		)
		object.Set(RightBitOr,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						if _, ok := left.(*Integer); !ok {
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName)
						}
						return p.NewInteger(false, p.PeekSymbolTable(),
							new(big.Int).Or(left.GetInteger(), self.GetInteger()),
						), nil
					},
				),
			),
		)
		object.Set(BitLeft,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						if _, ok := right.(*Integer); !ok {
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName)
						}
						return p.NewInteger(false, p.PeekSymbolTable(),
							new(big.Int).Lsh(self.GetInteger(), uint(right.GetInteger().Uint64())),
						), nil
					},
				),
			),
		)
		object.Set(RightBitLeft,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						if _, ok := left.(*Integer); !ok {
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName)
						}
						return p.NewInteger(false, p.PeekSymbolTable(),
							new(big.Int).Lsh(left.GetInteger(), uint(self.GetInteger().Uint64())),
						), nil
					},
				),
			),
		)
		object.Set(BitRight,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						if _, ok := right.(*Integer); !ok {
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName)
						}
						return p.NewInteger(false, p.PeekSymbolTable(),
							new(big.Int).Rsh(self.GetInteger(), uint(right.GetInteger().Uint64())),
						), nil
					},
				),
			),
		)
		object.Set(RightBitRight,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						if _, ok := left.(*Integer); !ok {
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName)
						}
						return p.NewInteger(false, p.PeekSymbolTable(),
							new(big.Int).Rsh(left.GetInteger(), uint(self.GetInteger().Uint64())),
						), nil
					},
				),
			),
		)
		object.Set(Equals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							return p.NewBool(false, p.PeekSymbolTable(),
								self.GetInteger().Cmp(right.GetInteger()) == 0,
							), nil
						case *Float:
							return p.NewBool(false, p.PeekSymbolTable(),
								new(big.Float).SetInt(self.GetInteger()).Cmp(right.GetFloat()) == 0,
							), nil
						default:
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(RightEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						switch left.(type) {
						case *Integer:
							return p.NewBool(false, p.PeekSymbolTable(),
								left.GetInteger().Cmp(self.GetInteger()) == 0,
							), nil
						case *Float:
							return p.NewBool(false, p.PeekSymbolTable(),
								left.GetFloat().Cmp(new(big.Float).SetInt(self.GetInteger())) == 0,
							), nil
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(NotEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							return p.NewBool(false, p.PeekSymbolTable(),
								self.GetInteger().Cmp(right.GetInteger()) != 0,
							), nil
						case *Float:
							return p.NewBool(false, p.PeekSymbolTable(),
								new(big.Float).SetInt(self.GetInteger()).Cmp(right.GetFloat()) != 0,
							), nil
						default:
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(RightNotEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						switch left.(type) {
						case *Integer:
							return p.NewBool(false, p.PeekSymbolTable(),
								left.GetInteger().Cmp(self.GetInteger()) != 0,
							), nil
						case *Float:
							return p.NewBool(false, p.PeekSymbolTable(),
								left.GetFloat().Cmp(new(big.Float).SetInt(self.GetInteger())) != 0,
							), nil
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(GreaterThan,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							return p.NewBool(false, p.PeekSymbolTable(),
								self.GetInteger().Cmp(right.GetInteger()) == 1,
							), nil
						case *Float:
							return p.NewBool(false, p.PeekSymbolTable(),
								new(big.Float).SetInt(self.GetInteger()).Cmp(right.GetFloat()) == 1,
							), nil
						default:
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(RightGreaterThan,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						switch left.(type) {
						case *Integer:
							return p.NewBool(false, p.PeekSymbolTable(),
								left.GetInteger().Cmp(self.GetInteger()) == 1,
							), nil
						case *Float:
							return p.NewBool(false, p.PeekSymbolTable(),
								left.GetFloat().Cmp(new(big.Float).SetInt(self.GetInteger())) == 1,
							), nil
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(LessThan,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							return p.NewBool(false, p.PeekSymbolTable(),
								self.GetInteger().Cmp(right.GetInteger()) == -1,
							), nil
						case *Float:
							return p.NewBool(false, p.PeekSymbolTable(),
								new(big.Float).SetInt(self.GetInteger()).Cmp(right.GetFloat()) == -1,
							), nil
						default:
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(RightLessThan,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						switch left.(type) {
						case *Integer:
							return p.NewBool(false, p.PeekSymbolTable(),
								left.GetInteger().Cmp(self.GetInteger()) == -1,
							), nil
						case *Float:
							return p.NewBool(false, p.PeekSymbolTable(),
								left.GetFloat().Cmp(new(big.Float).SetInt(self.GetInteger())) == -1,
							), nil
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(GreaterThanOrEqual,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							result := self.GetInteger().Cmp(right.GetInteger())
							return p.NewBool(false, p.PeekSymbolTable(),
								result == 0 || result == 1,
							), nil
						case *Float:
							result := new(big.Float).SetInt(self.GetInteger()).Cmp(right.GetFloat())
							return p.NewBool(false, p.PeekSymbolTable(),
								result == 0 || result == 1,
							), nil
						default:
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(RightGreaterThanOrEqual,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						switch left.(type) {
						case *Integer:
							result := left.GetInteger().Cmp(self.GetInteger())
							return p.NewBool(false, p.PeekSymbolTable(),
								result == 0 || result == 1,
							), nil
						case *Float:
							result := left.GetFloat().Cmp(new(big.Float).SetInt(self.GetInteger()))
							return p.NewBool(false, p.PeekSymbolTable(),
								result == 0 || result == 1,
							), nil
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(LessThanOrEqual,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							result := self.GetInteger().Cmp(right.GetInteger())
							return p.NewBool(false, p.PeekSymbolTable(),
								result == 0 || result == -1,
							), nil
						case *Float:
							result := new(big.Float).SetInt(self.GetInteger()).Cmp(right.GetFloat())
							return p.NewBool(false, p.PeekSymbolTable(),
								result == 0 || result == -1,
							), nil
						default:
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(RightLessThanOrEqual,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						left := arguments[0]
						switch left.(type) {
						case *Integer:
							result := left.GetInteger().Cmp(self.GetInteger())
							return p.NewBool(false, p.PeekSymbolTable(),
								result == 0 || result == -1,
							), nil
						case *Float:
							result := left.GetFloat().Cmp(new(big.Float).SetInt(self.GetInteger()))
							return p.NewBool(false, p.PeekSymbolTable(),
								result == 0 || result == -1,
							), nil
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
						}
					},
				),
			),
		)
		object.Set(Hash,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger()), nil
					},
				),
			),
		)
		object.Set(Copy,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger()), nil
					},
				),
			),
		)
		object.Set(ToInteger,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger()), nil
					},
				),
			),
		)
		object.Set(ToFloat,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewFloat(false, p.PeekSymbolTable(), new(big.Float).SetInt(self.GetInteger())), nil
					},
				),
			),
		)
		object.Set(ToString,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewString(false, p.PeekSymbolTable(), fmt.Sprint(self.GetInteger())), nil
					},
				),
			),
		)
		object.Set(ToBool,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewBool(false, p.PeekSymbolTable(),
							self.GetInteger().Cmp(big.NewInt(0)) != 0,
						), nil
					},
				),
			),
		)
		return nil
	}
}
