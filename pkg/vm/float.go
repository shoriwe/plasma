package vm

import (
	"fmt"
	"math"
)

type Float struct {
	*Object
}

func (p *Plasma) NewFloat(isBuiltIn bool, parentSymbols *SymbolTable, value float64) *Float {
	float_ := &Float{
		p.NewObject(isBuiltIn, IntegerName, nil, parentSymbols),
	}
	float_.SetFloat64(value)
	p.FloatInitialize(isBuiltIn)(float_)
	return float_
}

func (p *Plasma) FloatInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(object IObject) *Object {
		object.Set(Add,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							return p.NewFloat(false, p.PeekSymbolTable(), self.GetFloat64()+float64(right.GetInteger64())), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), self.GetFloat64()+right.GetFloat64()), nil
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
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						left := arguments[0]
						switch left.(type) {
						case *Integer:
							return p.NewFloat(false, p.PeekSymbolTable(), float64(left.GetInteger64())+self.GetFloat64()), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), left.GetFloat64()+self.GetFloat64()), nil
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
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							return p.NewFloat(false, p.PeekSymbolTable(), self.GetFloat64()-float64(right.GetInteger64())), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), self.GetFloat64()-right.GetFloat64()), nil
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
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						left := arguments[0]
						switch left.(type) {
						case *Integer:
							return p.NewFloat(false, p.PeekSymbolTable(), float64(left.GetInteger64())-self.GetFloat64()), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), left.GetFloat64()-self.GetFloat64()), nil
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
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							return p.NewFloat(false, p.PeekSymbolTable(), self.GetFloat64()*float64(right.GetInteger64())), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), self.GetFloat64()*right.GetFloat64()), nil
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
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						left := arguments[0]
						switch left.(type) {
						case *Integer:
							return p.NewFloat(false, p.PeekSymbolTable(), float64(left.GetInteger64())*self.GetFloat64()), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), left.GetFloat64()*self.GetFloat64()), nil
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
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							return p.NewFloat(false, p.PeekSymbolTable(), self.GetFloat64()/float64(right.GetInteger64())), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), self.GetFloat64()/right.GetFloat64()), nil
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
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						left := arguments[0]
						switch left.(type) {
						case *Integer:
							return p.NewFloat(false, p.PeekSymbolTable(), float64(left.GetInteger64())/self.GetFloat64()), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), left.GetFloat64()/self.GetFloat64()), nil
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
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							return p.NewFloat(false, p.PeekSymbolTable(), math.Mod(self.GetFloat64(), float64(right.GetInteger64()))), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), math.Mod(self.GetFloat64(), right.GetFloat64())), nil
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
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						left := arguments[0]
						switch left.(type) {
						case *Integer:
							return p.NewFloat(false, p.PeekSymbolTable(), math.Mod(float64(left.GetInteger64()), self.GetFloat64())), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), math.Mod(left.GetFloat64(), self.GetFloat64())), nil
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
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						right := arguments[0]
						switch right.(type) {
						case *Integer:
							return p.NewFloat(false, p.PeekSymbolTable(), math.Pow(self.GetFloat64(), float64(right.GetInteger64()))), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), math.Pow(self.GetFloat64(), right.GetFloat64())), nil
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
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						left := arguments[0]
						switch left.(type) {
						case *Integer:
							return p.NewFloat(false, p.PeekSymbolTable(), math.Pow(float64(left.GetInteger64()), self.GetFloat64())), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), math.Pow(left.GetFloat64(), self.GetFloat64())), nil
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
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						right := arguments[0]
						var floatRight float64
						switch right.(type) {
						case *Integer:
							floatRight = float64(right.GetInteger64())
						case *Float:
							floatRight = right.GetFloat64()
						default:
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
						}
						return p.NewBool(false, p.PeekSymbolTable(), self.GetFloat64() == floatRight), nil
					},
				),
			),
		)
		object.Set(RightEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						left := arguments[0]
						var floatLeft float64
						switch left.(type) {
						case *Integer:
							floatLeft = float64(left.GetInteger64())
						case *Float:
							floatLeft = left.GetFloat64()
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
						}
						return p.NewBool(false, p.PeekSymbolTable(), floatLeft == self.GetFloat64()), nil
					},
				),
			),
		)
		object.Set(NotEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						right := arguments[0]
						var floatRight float64
						switch right.(type) {
						case *Integer:
							floatRight = float64(right.GetInteger64())
						case *Float:
							floatRight = right.GetFloat64()
						default:
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
						}
						return p.NewBool(false, p.PeekSymbolTable(), self.GetFloat64() != floatRight), nil
					},
				),
			),
		)
		object.Set(RightNotEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						left := arguments[0]
						var floatLeft float64
						switch left.(type) {
						case *Integer:
							floatLeft = float64(left.GetInteger64())
						case *Float:
							floatLeft = left.GetFloat64()
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
						}
						return p.NewBool(false, p.PeekSymbolTable(), floatLeft != self.GetFloat64()), nil
					},
				),
			),
		)
		object.Set(GreaterThan,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						right := arguments[0]
						var floatRight float64
						switch right.(type) {
						case *Integer:
							floatRight = float64(right.GetInteger64())
						case *Float:
							floatRight = right.GetFloat64()
						default:
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
						}
						return p.NewBool(false, p.PeekSymbolTable(), self.GetFloat64() > floatRight), nil
					},
				),
			),
		)
		object.Set(RightGreaterThan,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						left := arguments[0]
						var floatLeft float64
						switch left.(type) {
						case *Integer:
							floatLeft = float64(left.GetInteger64())
						case *Float:
							floatLeft = left.GetFloat64()
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
						}
						return p.NewBool(false, p.PeekSymbolTable(), floatLeft > self.GetFloat64()), nil
					},
				),
			),
		)
		object.Set(LessThan,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						right := arguments[0]
						var floatRight float64
						switch right.(type) {
						case *Integer:
							floatRight = float64(right.GetInteger64())
						case *Float:
							floatRight = right.GetFloat64()
						default:
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
						}
						return p.NewBool(false, p.PeekSymbolTable(), self.GetFloat64() < floatRight), nil
					},
				),
			),
		)
		object.Set(RightLessThan,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						left := arguments[0]
						var floatLeft float64
						switch left.(type) {
						case *Integer:
							floatLeft = float64(left.GetInteger64())
						case *Float:
							floatLeft = left.GetFloat64()
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
						}
						return p.NewBool(false, p.PeekSymbolTable(), floatLeft < self.GetFloat64()), nil
					},
				),
			),
		)
		object.Set(GreaterThanOrEqual,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						right := arguments[0]
						var floatRight float64
						switch right.(type) {
						case *Integer:
							floatRight = float64(right.GetInteger64())
						case *Float:
							floatRight = right.GetFloat64()
						default:
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
						}
						return p.NewBool(false, p.PeekSymbolTable(), self.GetFloat64() >= floatRight), nil
					},
				),
			),
		)
		object.Set(RightGreaterThanOrEqual,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						left := arguments[0]
						var floatLeft float64
						switch left.(type) {
						case *Integer:
							floatLeft = float64(left.GetInteger64())
						case *Float:
							floatLeft = left.GetFloat64()
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
						}
						return p.NewBool(false, p.PeekSymbolTable(), floatLeft >= self.GetFloat64()), nil
					},
				),
			),
		)
		object.Set(LessThanOrEqual,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						right := arguments[0]
						var floatRight float64
						switch right.(type) {
						case *Integer:
							floatRight = float64(right.GetInteger64())
						case *Float:
							floatRight = right.GetFloat64()
						default:
							return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
						}
						return p.NewBool(false, p.PeekSymbolTable(), self.GetFloat64() <= floatRight), nil
					},
				),
			),
		)
		object.Set(RightLessThanOrEqual,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self IObject, arguments ...IObject) (IObject, *Object) {
						left := arguments[0]
						var floatLeft float64
						switch left.(type) {
						case *Integer:
							floatLeft = float64(left.GetInteger64())
						case *Float:
							floatLeft = left.GetFloat64()
						default:
							return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
						}
						return p.NewBool(false, p.PeekSymbolTable(), floatLeft <= self.GetFloat64()), nil
					},
				),
			),
		)

		object.Set(Hash,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self IObject, _ ...IObject) (IObject, *Object) {
						if self.GetHash() == 0 {
							floatHash := p.HashString(fmt.Sprintf("%f-%s", self.GetFloat64(), FloatName))
							self.SetHash(floatHash)
						}
						return p.NewInteger(false, p.PeekSymbolTable(), self.GetHash()), nil
					},
				),
			),
		)
		object.Set(Copy,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self IObject, _ ...IObject) (IObject, *Object) {
						return p.NewFloat(false, p.PeekSymbolTable(), self.GetFloat64()), nil
					},
				),
			),
		)

		object.Set(ToInteger,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self IObject, _ ...IObject) (IObject, *Object) {
						return p.NewInteger(false, p.PeekSymbolTable(), int64(self.GetFloat64())), nil
					},
				),
			),
		)
		object.Set(ToFloat,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self IObject, _ ...IObject) (IObject, *Object) {
						return p.NewFloat(false, p.PeekSymbolTable(), self.GetFloat64()), nil
					},
				),
			),
		)
		object.Set(ToString,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self IObject, _ ...IObject) (IObject, *Object) {
						return p.NewString(false, p.PeekSymbolTable(), fmt.Sprint(self.GetFloat64())), nil
					},
				),
			),
		)
		object.Set(ToBool,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self IObject, _ ...IObject) (IObject, *Object) {
						return p.NewBool(false, p.PeekSymbolTable(), self.GetFloat64() != 0), nil
					},
				),
			),
		)
		return nil
	}
}
