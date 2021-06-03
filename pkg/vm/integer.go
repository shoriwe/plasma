package vm

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/errors"
	"github.com/shoriwe/gruby/pkg/tools"
	"math"
)

type Integer struct {
	*Object
}

func (p *Plasma) NewInteger(parentSymbols *SymbolTable, value int64) *Integer {
	integer := &Integer{
		p.NewObject(IntegerName, nil, parentSymbols),
	}
	integer.SetInteger64(value)
	p.IntegerInitialize(integer)
	return integer
}

func (p *Plasma) IntegerInitialize(object IObject) *errors.Error {
	object.SymbolTable().Set(NegBits,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewInteger(p.PeekSymbolTable(), ^self.GetInteger64()), nil
				},
			),
		),
	)
	object.SymbolTable().Set(Negative,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewInteger(p.PeekSymbolTable(), -self.GetInteger64()), nil
				},
			),
		),
	)
	object.SymbolTable().Set(Add,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					switch right.(type) {
					case *Integer:
						return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()+right.GetInteger64()), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), float64(self.GetInteger64())+right.GetFloat64()), nil
					default:
						return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.SymbolTable().Set(RightAdd,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					switch left.(type) {
					case *Integer:
						return p.NewInteger(p.PeekSymbolTable(), left.GetInteger64()+self.GetInteger64()), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), left.GetFloat64()+float64(self.GetInteger64())), nil
					default:
						return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.SymbolTable().Set(Sub,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					switch right.(type) {
					case *Integer:
						return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()-right.GetInteger64()), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), float64(self.GetInteger64())-right.GetFloat64()), nil
					default:
						return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.SymbolTable().Set(RightSub,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					switch left.(type) {
					case *Integer:
						return p.NewInteger(p.PeekSymbolTable(), left.GetInteger64()-self.GetInteger64()), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), left.GetFloat64()-float64(self.GetInteger64())), nil
					default:
						return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.SymbolTable().Set(Mul,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					switch right.(type) {
					case *Integer:
						return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()*right.GetInteger64()), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), float64(self.GetInteger64())*right.GetFloat64()), nil
					case *String:
						return p.NewString(p.PeekSymbolTable(), tools.Repeat(right.GetString(), self.GetInteger64())), nil
					case *Tuple:
						panic(NewNotImplementedCallable(errors.UnknownLine))
					case *Array:
						panic(NewNotImplementedCallable(errors.UnknownLine))
					default:
						return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName)
					}
				},
			),
		),
	)
	object.SymbolTable().Set(RightMul,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					switch left.(type) {
					case *Integer:
						return p.NewInteger(p.PeekSymbolTable(), left.GetInteger64()*self.GetInteger64()), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), left.GetFloat64()*float64(self.GetInteger64())), nil
					case *String:
						return p.NewString(p.PeekSymbolTable(), tools.Repeat(left.GetString(), self.GetInteger64())), nil
					case *Tuple:
						panic(NewNotImplementedCallable(errors.UnknownLine))
					case *Array:
						panic(NewNotImplementedCallable(errors.UnknownLine))
					default:
						return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName)
					}
				},
			),
		),
	)
	object.SymbolTable().Set(Div,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					switch right.(type) {
					case *Integer:
						return p.NewFloat(p.PeekSymbolTable(), float64(self.GetInteger64())/float64(right.GetInteger64())), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), float64(self.GetInteger64())/right.GetFloat64()), nil
					default:
						return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.SymbolTable().Set(RightDiv,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					switch left.(type) {
					case *Integer:
						return p.NewFloat(p.PeekSymbolTable(), float64(left.GetInteger64())/float64(self.GetInteger64())), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), left.GetFloat64()/float64(self.GetInteger64())), nil
					default:
						return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.SymbolTable().Set(Mod,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					switch right.(type) {
					case *Integer:
						return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()%right.GetInteger64()), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), math.Mod(float64(self.GetInteger64()), right.GetFloat64())), nil
					default:
						return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.SymbolTable().Set(RightMod,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					switch left.(type) {
					case *Integer:
						return p.NewInteger(p.PeekSymbolTable(), left.GetInteger64()%self.GetInteger64()), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), math.Mod(left.GetFloat64(), float64(self.GetInteger64()))), nil
					default:
						return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.SymbolTable().Set(Pow,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					switch right.(type) {
					case *Integer:
						return p.NewFloat(p.PeekSymbolTable(), math.Pow(float64(self.GetInteger64()), float64(right.GetInteger64()))), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), math.Pow(float64(self.GetInteger64()), right.GetFloat64())), nil
					default:
						return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.SymbolTable().Set(RightPow,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					switch left.(type) {
					case *Integer:
						return p.NewFloat(p.PeekSymbolTable(), math.Pow(float64(left.GetInteger64()), float64(self.GetInteger64()))), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), math.Pow(left.GetFloat64(), float64(self.GetInteger64()))), nil
					default:
						return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.SymbolTable().Set(BitXor,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					if _, ok := right.(*Integer); !ok {
						return nil, errors.NewTypeError(right.TypeName(), IntegerName)
					}
					return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()^right.GetInteger64()), nil
				},
			),
		),
	)
	object.SymbolTable().Set(RightBitXor,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					if _, ok := left.(*Integer); !ok {
						return nil, errors.NewTypeError(left.TypeName(), IntegerName)
					}
					return p.NewInteger(p.PeekSymbolTable(), left.GetInteger64()^self.GetInteger64()), nil
				},
			),
		),
	)
	object.SymbolTable().Set(BitAnd,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					if _, ok := right.(*Integer); !ok {
						return nil, errors.NewTypeError(right.TypeName(), IntegerName)
					}
					return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()&right.GetInteger64()), nil
				},
			),
		),
	)
	object.SymbolTable().Set(RightBitAnd,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					if _, ok := left.(*Integer); !ok {
						return nil, errors.NewTypeError(left.TypeName(), IntegerName)
					}
					return p.NewInteger(p.PeekSymbolTable(), left.GetInteger64()&self.GetInteger64()), nil
				},
			),
		),
	)
	object.SymbolTable().Set(BitOr,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					if _, ok := right.(*Integer); !ok {
						return nil, errors.NewTypeError(right.TypeName(), IntegerName)
					}
					return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()|right.GetInteger64()), nil
				},
			),
		),
	)
	object.SymbolTable().Set(RightBitOr,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					if _, ok := left.(*Integer); !ok {
						return nil, errors.NewTypeError(left.TypeName(), IntegerName)
					}
					return p.NewInteger(p.PeekSymbolTable(), left.GetInteger64()|self.GetInteger64()), nil
				},
			),
		),
	)
	object.SymbolTable().Set(BitLeft,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					if _, ok := right.(*Integer); !ok {
						return nil, errors.NewTypeError(right.TypeName(), IntegerName)
					}
					return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()<<right.GetInteger64()), nil
				},
			),
		),
	)
	object.SymbolTable().Set(RightBitLeft,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					if _, ok := left.(*Integer); !ok {
						return nil, errors.NewTypeError(left.TypeName(), IntegerName)
					}
					return p.NewInteger(p.PeekSymbolTable(), left.GetInteger64()<<self.GetInteger64()), nil
				},
			),
		),
	)
	object.SymbolTable().Set(BitRight,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					if _, ok := right.(*Integer); !ok {
						return nil, errors.NewTypeError(right.TypeName(), IntegerName)
					}
					return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()>>right.GetInteger64()), nil
				},
			),
		),
	)
	object.SymbolTable().Set(RightBitRight,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					if _, ok := left.(*Integer); !ok {
						return nil, errors.NewTypeError(left.TypeName(), IntegerName)
					}
					return p.NewInteger(p.PeekSymbolTable(), left.GetInteger64()>>self.GetInteger64()), nil
				},
			),
		),
	)
	object.SymbolTable().Set(Equals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					var floatRight float64
					switch right.(type) {
					case *Integer:
						floatRight = float64(right.GetInteger64())
					case *Float:
						floatRight = right.GetFloat64()
					default:
						return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
					}
					return p.NewBool(p.PeekSymbolTable(), float64(self.GetInteger64()) == floatRight), nil
				},
			),
		),
	)
	object.SymbolTable().Set(RightEquals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					var floatLeft float64
					switch left.(type) {
					case *Integer:
						floatLeft = float64(left.GetInteger64())
					case *Float:
						floatLeft = left.GetFloat64()
					default:
						return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
					}
					return p.NewBool(p.PeekSymbolTable(), floatLeft == float64(self.GetInteger64())), nil
				},
			),
		),
	)
	object.SymbolTable().Set(NotEquals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					var floatRight float64
					switch right.(type) {
					case *Integer:
						floatRight = float64(right.GetInteger64())
					case *Float:
						floatRight = right.GetFloat64()
					default:
						return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
					}
					return p.NewBool(p.PeekSymbolTable(), float64(self.GetInteger64()) != floatRight), nil
				},
			),
		),
	)
	object.SymbolTable().Set(RightNotEquals,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					var floatLeft float64
					switch left.(type) {
					case *Integer:
						floatLeft = float64(left.GetInteger64())
					case *Float:
						floatLeft = left.GetFloat64()
					default:
						return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
					}
					return p.NewBool(p.PeekSymbolTable(), floatLeft != float64(self.GetInteger64())), nil
				},
			),
		),
	)
	object.SymbolTable().Set(GreaterThan,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					var floatRight float64
					switch right.(type) {
					case *Integer:
						floatRight = float64(right.GetInteger64())
					case *Float:
						floatRight = right.GetFloat64()
					default:
						return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
					}
					return p.NewBool(p.PeekSymbolTable(), float64(self.GetInteger64()) > floatRight), nil
				},
			),
		),
	)
	object.SymbolTable().Set(RightGreaterThan,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					var floatLeft float64
					switch left.(type) {
					case *Integer:
						floatLeft = float64(left.GetInteger64())
					case *Float:
						floatLeft = left.GetFloat64()
					default:
						return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
					}
					return p.NewBool(p.PeekSymbolTable(), floatLeft > float64(self.GetInteger64())), nil
				},
			),
		),
	)
	object.SymbolTable().Set(LessThan,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					var floatRight float64
					switch right.(type) {
					case *Integer:
						floatRight = float64(right.GetInteger64())
					case *Float:
						floatRight = right.GetFloat64()
					default:
						return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
					}
					return p.NewBool(p.PeekSymbolTable(), float64(self.GetInteger64()) < floatRight), nil
				},
			),
		),
	)
	object.SymbolTable().Set(RightLessThan,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					var floatLeft float64
					switch left.(type) {
					case *Integer:
						floatLeft = float64(left.GetInteger64())
					case *Float:
						floatLeft = left.GetFloat64()
					default:
						return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
					}
					return p.NewBool(p.PeekSymbolTable(), floatLeft < float64(self.GetInteger64())), nil
				},
			),
		),
	)
	object.SymbolTable().Set(GreaterThanOrEqual,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					var floatRight float64
					switch right.(type) {
					case *Integer:
						floatRight = float64(right.GetInteger64())
					case *Float:
						floatRight = right.GetFloat64()
					default:
						return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
					}
					return p.NewBool(p.PeekSymbolTable(), float64(self.GetInteger64()) >= floatRight), nil
				},
			),
		),
	)
	object.SymbolTable().Set(RightGreaterThanOrEqual,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					var floatLeft float64
					switch left.(type) {
					case *Integer:
						floatLeft = float64(left.GetInteger64())
					case *Float:
						floatLeft = left.GetFloat64()
					default:
						return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
					}
					return p.NewBool(p.PeekSymbolTable(), floatLeft >= float64(self.GetInteger64())), nil
				},
			),
		),
	)
	object.SymbolTable().Set(LessThanOrEqual,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					right := arguments[0]
					var floatRight float64
					switch right.(type) {
					case *Integer:
						floatRight = float64(right.GetInteger64())
					case *Float:
						floatRight = right.GetFloat64()
					default:
						return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
					}
					return p.NewBool(p.PeekSymbolTable(), float64(self.GetInteger64()) <= floatRight), nil
				},
			),
		),
	)
	object.SymbolTable().Set(RightLessThanOrEqual,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *errors.Error) {
					left := arguments[0]
					var floatLeft float64
					switch left.(type) {
					case *Integer:
						floatLeft = float64(left.GetInteger64())
					case *Float:
						floatLeft = left.GetFloat64()
					default:
						return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
					}
					return p.NewBool(p.PeekSymbolTable(), floatLeft <= float64(self.GetInteger64())), nil
				},
			),
		),
	)
	object.SymbolTable().Set(Hash,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()), nil
				},
			),
		),
	)
	object.SymbolTable().Set(Copy,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()), nil
				},
			),
		),
	)
	object.SymbolTable().Set(ToInteger,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()), nil
				},
			),
		),
	)
	object.SymbolTable().Set(ToFloat,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewFloat(p.PeekSymbolTable(), float64(self.GetInteger64())), nil
				},
			),
		),
	)
	object.SymbolTable().Set(ToString,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewString(p.PeekSymbolTable(), fmt.Sprint(self.GetInteger64())), nil
				},
			),
		),
	)
	object.SymbolTable().Set(ToBool,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *errors.Error) {
					return p.NewBool(p.PeekSymbolTable(), self.GetInteger64() != 0), nil
				},
			),
		),
	)
	return nil
}
