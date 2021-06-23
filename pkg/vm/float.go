package vm

import (
	"fmt"
	"math"
	"math/big"
)

type Float struct {
	*Object
}

func (p *Plasma) NewFloat(isBuiltIn bool, parentSymbols *SymbolTable, value *big.Float) *Float {
	float_ := &Float{
		p.NewObject(isBuiltIn, FloatName, nil, parentSymbols),
	}
	float_.SetFloat(value)
	p.FloatInitialize(isBuiltIn)(float_)
	float_.Set(Self, float_)
	return float_
}

func (p *Plasma) FloatInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(object Value) *Object {
		object.Set(Add,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Add(self.GetFloat(), new(big.Float).SetInt(right.GetInteger())),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Add(self.GetFloat(), right.GetFloat()),
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
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Add(new(big.Float).SetInt(left.GetInteger()), self.GetFloat()),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Add(left.GetFloat(), self.GetFloat()),
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
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Sub(self.GetFloat(), new(big.Float).SetInt(right.GetInteger())),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Sub(self.GetFloat(), right.GetFloat()),
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
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Sub(new(big.Float).SetInt(left.GetInteger()), self.GetFloat()),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Sub(left.GetFloat(), self.GetFloat()),
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
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Mul(self.GetFloat(), new(big.Float).SetInt(right.GetInteger())),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Mul(self.GetFloat(), right.GetFloat()),
							), nil
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
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Mul(new(big.Float).SetInt(left.GetInteger()), self.GetFloat()),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Mul(left.GetFloat(), self.GetFloat()),
							), nil
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
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
								new(big.Float).Quo(self.GetFloat(), new(big.Float).SetInt(right.GetInteger())),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Quo(self.GetFloat(), right.GetFloat()),
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
								new(big.Float).Quo(new(big.Float).SetInt(left.GetInteger()), self.GetFloat()),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								new(big.Float).Quo(left.GetFloat(), self.GetFloat()),
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
							leftIntPart, _ := self.GetFloat().Int(nil)
							return p.NewInteger(false, p.PeekSymbolTable(),
								new(big.Int).Mod(leftIntPart, right.GetInteger()),
							), nil
						case *Float:
							leftInt, _ := self.GetFloat().Int(nil)
							rightInt, _ := right.GetFloat().Int(nil)
							return p.NewInteger(false, p.PeekSymbolTable(),
								new(big.Int).Mod(leftInt, rightInt),
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
							rightInt, _ := self.GetFloat().Int(nil)
							return p.NewInteger(false, p.PeekSymbolTable(),
								new(big.Int).Mod(left.GetInteger(), rightInt),
							), nil
						case *Float:
							leftInt, _ := left.GetFloat().Int(nil)
							rightInt, _ := self.GetFloat().Int(nil)
							return p.NewInteger(false, p.PeekSymbolTable(),
								new(big.Int).Mod(leftInt, rightInt),
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
							leftHandSide, _ := self.GetFloat().Float64()
							return p.NewFloat(false, p.PeekSymbolTable(),
								big.NewFloat(
									math.Pow(
										leftHandSide,
										float64(right.GetInteger().Int64(),
										),
									),
								),
							), nil
						case *Float:
							leftHandSide, _ := self.GetFloat().Float64()
							rightHandSide, _ := right.GetFloat().Float64()
							return p.NewFloat(false, p.PeekSymbolTable(),
								big.NewFloat(
									math.Pow(
										leftHandSide,
										rightHandSide,
									),
								),
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
							rightHandSide, _ := self.GetFloat().Float64()
							return p.NewFloat(false, p.PeekSymbolTable(),
								big.NewFloat(
									math.Pow(
										float64(left.GetInteger().Int64()),
										rightHandSide,
									),
								),
							), nil
						case *Float:
							leftHandSide, _ := left.GetFloat().Float64()
							rightHandSide, _ := self.GetFloat().Float64()
							return p.NewFloat(false, p.PeekSymbolTable(),
								big.NewFloat(
									math.Pow(
										leftHandSide,
										rightHandSide,
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

		object.Set(Equals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							return p.NewBool(false, p.PeekSymbolTable(),
								self.GetFloat().Cmp(new(big.Float).SetInt(right.GetInteger())) == 0,
							), nil
						case *Float:
							return p.NewBool(false, p.PeekSymbolTable(),
								self.GetFloat().Cmp(right.GetFloat()) == 0,
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
								new(big.Float).SetInt(left.GetInteger()).Cmp(self.GetFloat()) == 0,
							), nil
						case *Float:
							return p.NewBool(false, p.PeekSymbolTable(),
								left.GetFloat().Cmp(self.GetFloat()) == 0,
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
								self.GetFloat().Cmp(new(big.Float).SetInt(right.GetInteger())) != 0,
							), nil
						case *Float:
							return p.NewBool(false, p.PeekSymbolTable(),
								self.GetFloat().Cmp(right.GetFloat()) != 0,
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
								new(big.Float).SetInt(left.GetInteger()).Cmp(self.GetFloat()) != 0,
							), nil
						case *Float:
							return p.NewBool(false, p.PeekSymbolTable(),
								left.GetFloat().Cmp(self.GetFloat()) != 0,
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
								self.GetFloat().Cmp(new(big.Float).SetInt(right.GetInteger())) == 1,
							), nil
						case *Float:
							return p.NewBool(false, p.PeekSymbolTable(),
								self.GetFloat().Cmp(right.GetFloat()) == 1,
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
								new(big.Float).SetInt(left.GetInteger()).Cmp(self.GetFloat()) == 1,
							), nil
						case *Float:
							return p.NewBool(false, p.PeekSymbolTable(),
								left.GetFloat().Cmp(self.GetFloat()) == 1,
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
								self.GetFloat().Cmp(new(big.Float).SetInt(right.GetInteger())) == -1,
							), nil
						case *Float:
							return p.NewBool(false, p.PeekSymbolTable(),
								self.GetFloat().Cmp(right.GetFloat()) == -1,
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
								new(big.Float).SetInt(left.GetInteger()).Cmp(self.GetFloat()) == -1,
							), nil
						case *Float:
							return p.NewBool(false, p.PeekSymbolTable(),
								left.GetFloat().Cmp(self.GetFloat()) == -1,
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
							result := self.GetFloat().Cmp(new(big.Float).SetInt(right.GetInteger()))
							return p.NewBool(false, p.PeekSymbolTable(),
								result == 0 || result == 1,
							), nil
						case *Float:
							result := self.GetFloat().Cmp(right.GetFloat())
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
							result := new(big.Float).SetInt(left.GetInteger()).Cmp(self.GetFloat())
							return p.NewBool(false, p.PeekSymbolTable(),
								result == 0 || result == 1,
							), nil
						case *Float:
							result := left.GetFloat().Cmp(self.GetFloat())
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
							result := self.GetFloat().Cmp(new(big.Float).SetInt(right.GetInteger()))
							return p.NewBool(false, p.PeekSymbolTable(),
								result == 0 || result == -1,
							), nil
						case *Float:
							result := self.GetFloat().Cmp(right.GetFloat())
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
							result := new(big.Float).SetInt(left.GetInteger()).Cmp(self.GetFloat())
							return p.NewBool(false, p.PeekSymbolTable(),
								result == 0 || result == -1,
							), nil
						case *Float:
							result := left.GetFloat().Cmp(self.GetFloat())
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
						if self.GetHash() == 0 {
							floatHash := p.HashString(fmt.Sprintf("%f-%s", self.GetFloat(), FloatName))
							self.SetHash(floatHash)
						}
						return p.NewInteger(false, p.PeekSymbolTable(), big.NewInt(self.GetHash())), nil
					},
				),
			),
		)
		object.Set(Copy,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewFloat(false, p.PeekSymbolTable(), self.GetFloat()), nil
					},
				),
			),
		)

		object.Set(ToInteger,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						result, _ := self.GetFloat().Int(nil)
						return p.NewInteger(false, p.PeekSymbolTable(), result), nil
					},
				),
			),
		)
		object.Set(ToFloat,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewFloat(false, p.PeekSymbolTable(), self.GetFloat()), nil
					},
				),
			),
		)
		object.Set(ToString,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewString(false, p.PeekSymbolTable(), fmt.Sprint(self.GetFloat())), nil
					},
				),
			),
		)
		object.Set(ToBool,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewBool(false, p.PeekSymbolTable(), self.GetFloat().Cmp(big.NewFloat(0)) != 0), nil
					},
				),
			),
		)
		return nil
	}
}
