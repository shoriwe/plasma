package vm

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/errors"
	"math"
)

type Float struct {
	*Object
}

func (p *Plasma) NewFloat(parentSymbols *SymbolTable, value float64) *Float {
	float_ := &Float{
		p.NewObject(IntegerName, nil, parentSymbols),
	}
	float_.SetFloat64(value)
	p.FloatInitialize(float_)
	return float_
}

func (p *Plasma) FloatInitialize(object IObject) *errors.Error {
	object.Set(Add,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					right := arguments[0]
					switch right.(type) {
					case *Integer:
						return p.NewFloat(p.PeekSymbolTable(), self.GetFloat64()+float64(right.GetInteger64())), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), self.GetFloat64()+right.GetFloat64()), nil
					default:
						return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.Set(RightAdd,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					left := arguments[0]
					switch left.(type) {
					case *Integer:
						return p.NewFloat(p.PeekSymbolTable(), float64(left.GetInteger64())+self.GetFloat64()), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), left.GetFloat64()+self.GetFloat64()), nil
					default:
						return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.Set(Sub,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					right := arguments[0]
					switch right.(type) {
					case *Integer:
						return p.NewFloat(p.PeekSymbolTable(), self.GetFloat64()-float64(right.GetInteger64())), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), self.GetFloat64()-right.GetFloat64()), nil
					default:
						return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.Set(RightSub,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					left := arguments[0]
					switch left.(type) {
					case *Integer:
						return p.NewFloat(p.PeekSymbolTable(), float64(left.GetInteger64())-self.GetFloat64()), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), left.GetFloat64()-self.GetFloat64()), nil
					default:
						return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.Set(Mul,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					right := arguments[0]
					switch right.(type) {
					case *Integer:
						return p.NewFloat(p.PeekSymbolTable(), self.GetFloat64()*float64(right.GetInteger64())), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), self.GetFloat64()*right.GetFloat64()), nil
					default:
						return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName)
					}
				},
			),
		),
	)
	object.Set(RightMul,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					left := arguments[0]
					switch left.(type) {
					case *Integer:
						return p.NewFloat(p.PeekSymbolTable(), float64(left.GetInteger64())*self.GetFloat64()), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), left.GetFloat64()*self.GetFloat64()), nil
					default:
						return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.Set(Div,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					right := arguments[0]
					switch right.(type) {
					case *Integer:
						return p.NewFloat(p.PeekSymbolTable(), self.GetFloat64()/float64(right.GetInteger64())), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), self.GetFloat64()/right.GetFloat64()), nil
					default:
						return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.Set(RightDiv,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					left := arguments[0]
					switch left.(type) {
					case *Integer:
						return p.NewFloat(p.PeekSymbolTable(), float64(left.GetInteger64())/self.GetFloat64()), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), left.GetFloat64()/self.GetFloat64()), nil
					default:
						return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.Set(Mod,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					right := arguments[0]
					switch right.(type) {
					case *Integer:
						return p.NewFloat(p.PeekSymbolTable(), math.Mod(self.GetFloat64(), float64(right.GetInteger64()))), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), math.Mod(self.GetFloat64(), right.GetFloat64())), nil
					default:
						return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.Set(RightMod,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					left := arguments[0]
					switch left.(type) {
					case *Integer:
						return p.NewFloat(p.PeekSymbolTable(), math.Mod(float64(left.GetInteger64()), self.GetFloat64())), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), math.Mod(left.GetFloat64(), self.GetFloat64())), nil
					default:
						return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.Set(Pow,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					right := arguments[0]
					switch right.(type) {
					case *Integer:
						return p.NewFloat(p.PeekSymbolTable(), math.Pow(self.GetFloat64(), float64(right.GetInteger64()))), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), math.Pow(self.GetFloat64(), right.GetFloat64())), nil
					default:
						return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.Set(RightPow,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					left := arguments[0]
					switch left.(type) {
					case *Integer:
						return p.NewFloat(p.PeekSymbolTable(), math.Pow(float64(left.GetInteger64()), self.GetFloat64())), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), math.Pow(left.GetFloat64(), self.GetFloat64())), nil
					default:
						return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)

	object.Set(Equals,
		p.NewFunction(object.SymbolTable(),
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
					return p.NewBool(p.PeekSymbolTable(), self.GetFloat64() == floatRight), nil
				},
			),
		),
	)
	object.Set(RightEquals,
		p.NewFunction(object.SymbolTable(),
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
					return p.NewBool(p.PeekSymbolTable(), floatLeft == self.GetFloat64()), nil
				},
			),
		),
	)
	object.Set(NotEquals,
		p.NewFunction(object.SymbolTable(),
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
					return p.NewBool(p.PeekSymbolTable(), self.GetFloat64() != floatRight), nil
				},
			),
		),
	)
	object.Set(RightNotEquals,
		p.NewFunction(object.SymbolTable(),
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
					return p.NewBool(p.PeekSymbolTable(), floatLeft != self.GetFloat64()), nil
				},
			),
		),
	)
	object.Set(GreaterThan,
		p.NewFunction(object.SymbolTable(),
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
					return p.NewBool(p.PeekSymbolTable(), self.GetFloat64() > floatRight), nil
				},
			),
		),
	)
	object.Set(RightGreaterThan,
		p.NewFunction(object.SymbolTable(),
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
					return p.NewBool(p.PeekSymbolTable(), floatLeft > self.GetFloat64()), nil
				},
			),
		),
	)
	object.Set(LessThan,
		p.NewFunction(object.SymbolTable(),
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
					return p.NewBool(p.PeekSymbolTable(), self.GetFloat64() < floatRight), nil
				},
			),
		),
	)
	object.Set(RightLessThan,
		p.NewFunction(object.SymbolTable(),
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
					return p.NewBool(p.PeekSymbolTable(), floatLeft < self.GetFloat64()), nil
				},
			),
		),
	)
	object.Set(GreaterThanOrEqual,
		p.NewFunction(object.SymbolTable(),
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
					return p.NewBool(p.PeekSymbolTable(), self.GetFloat64() >= floatRight), nil
				},
			),
		),
	)
	object.Set(RightGreaterThanOrEqual,
		p.NewFunction(object.SymbolTable(),
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
					return p.NewBool(p.PeekSymbolTable(), floatLeft >= self.GetFloat64()), nil
				},
			),
		),
	)
	object.Set(LessThanOrEqual,
		p.NewFunction(object.SymbolTable(),
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
					return p.NewBool(p.PeekSymbolTable(), self.GetFloat64() <= floatRight), nil
				},
			),
		),
	)
	object.Set(RightLessThanOrEqual,
		p.NewFunction(object.SymbolTable(),
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
					return p.NewBool(p.PeekSymbolTable(), floatLeft <= self.GetFloat64()), nil
				},
			),
		),
	)

	object.Set(Hash,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					if self.GetHash() == 0 {
						floatHash := p.HashString(fmt.Sprintf("%f-%s", self.GetFloat64(), FloatName))
						self.SetHash(floatHash)
					}
					return p.NewInteger(p.PeekSymbolTable(), self.GetHash()), nil
				},
			),
		),
	)
	object.Set(Copy,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewFloat(p.PeekSymbolTable(), self.GetFloat64()), nil
				},
			),
		),
	)

	object.Set(ToInteger,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewInteger(p.PeekSymbolTable(), int64(self.GetFloat64())), nil
				},
			),
		),
	)
	object.Set(ToFloat,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewFloat(p.PeekSymbolTable(), self.GetFloat64()), nil
				},
			),
		),
	)
	object.Set(ToString,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewString(p.PeekSymbolTable(), fmt.Sprint(self.GetFloat64())), nil
				},
			),
		),
	)
	object.Set(ToBool,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewBool(p.PeekSymbolTable(), self.GetFloat64() != 0), nil
				},
			),
		),
	)
	return nil
}
