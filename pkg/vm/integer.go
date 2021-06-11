package vm

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/tools"
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

func (p *Plasma) IntegerInitialize(object IObject) *Object {
	object.Set(NegBits,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewInteger(p.PeekSymbolTable(), ^self.GetInteger64()), nil
				},
			),
		),
	)
	object.Set(Negative,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewInteger(p.PeekSymbolTable(), -self.GetInteger64()), nil
				},
			),
		),
	)
	object.Set(Add,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					right := arguments[0]
					switch right.(type) {
					case *Integer:
						return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()+right.GetInteger64()), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), float64(self.GetInteger64())+right.GetFloat64()), nil
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
						return p.NewInteger(p.PeekSymbolTable(), left.GetInteger64()+self.GetInteger64()), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), left.GetFloat64()+float64(self.GetInteger64())), nil
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
						return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()-right.GetInteger64()), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), float64(self.GetInteger64())-right.GetFloat64()), nil
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
						return p.NewInteger(p.PeekSymbolTable(), left.GetInteger64()-self.GetInteger64()), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), left.GetFloat64()-float64(self.GetInteger64())), nil
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
						return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()*right.GetInteger64()), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), float64(self.GetInteger64())*right.GetFloat64()), nil
					case *String:
						return p.NewString(p.PeekSymbolTable(), tools.Repeat(right.GetString(), self.GetInteger64())), nil
					case *Tuple:
						content, repetitionError := p.Repeat(right.GetContent(), int(self.GetInteger64()))
						if repetitionError != nil {
							return nil, repetitionError
						}
						return p.NewTuple(p.PeekSymbolTable(), content), nil
					case *Array:
						content, repetitionError := p.Repeat(right.GetContent(), int(self.GetInteger64()))
						if repetitionError != nil {
							return nil, repetitionError
						}
						return p.NewArray(p.PeekSymbolTable(), content), nil
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
						return p.NewInteger(p.PeekSymbolTable(), left.GetInteger64()*self.GetInteger64()), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), left.GetFloat64()*float64(self.GetInteger64())), nil
					case *String:
						return p.NewString(p.PeekSymbolTable(), tools.Repeat(left.GetString(), self.GetInteger64())), nil
					case *Tuple:
						content, repetitionError := p.Repeat(left.GetContent(), int(self.GetInteger64()))
						if repetitionError != nil {
							return nil, repetitionError
						}
						return p.NewTuple(p.PeekSymbolTable(), content), nil
					case *Array:
						content, repetitionError := p.Repeat(left.GetContent(), int(self.GetInteger64()))
						if repetitionError != nil {
							return nil, repetitionError
						}
						return p.NewArray(p.PeekSymbolTable(), content), nil
					default:
						return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName)
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
						return p.NewFloat(p.PeekSymbolTable(), float64(self.GetInteger64())/float64(right.GetInteger64())), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), float64(self.GetInteger64())/right.GetFloat64()), nil
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
						return p.NewFloat(p.PeekSymbolTable(), float64(left.GetInteger64())/float64(self.GetInteger64())), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), left.GetFloat64()/float64(self.GetInteger64())), nil
					default:
						return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.Set(FloorDiv,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					right := arguments[0]
					switch right.(type) {
					case *Integer:
						return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()/right.GetInteger64()), nil
					case *Float:
						return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()/int64(right.GetFloat64())), nil
					default:
						return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.Set(RightFloorDiv,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					left := arguments[0]
					switch left.(type) {
					case *Integer:
						return p.NewInteger(p.PeekSymbolTable(), left.GetInteger64()/self.GetInteger64()), nil
					case *Float:
						return p.NewInteger(p.PeekSymbolTable(), int64(left.GetFloat64())/self.GetInteger64()), nil
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
						return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()%right.GetInteger64()), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), math.Mod(float64(self.GetInteger64()), right.GetFloat64())), nil
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
						return p.NewInteger(p.PeekSymbolTable(), left.GetInteger64()%self.GetInteger64()), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), math.Mod(left.GetFloat64(), float64(self.GetInteger64()))), nil
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
						return p.NewFloat(p.PeekSymbolTable(), math.Pow(float64(self.GetInteger64()), float64(right.GetInteger64()))), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), math.Pow(float64(self.GetInteger64()), right.GetFloat64())), nil
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
						return p.NewFloat(p.PeekSymbolTable(), math.Pow(float64(left.GetInteger64()), float64(self.GetInteger64()))), nil
					case *Float:
						return p.NewFloat(p.PeekSymbolTable(), math.Pow(left.GetFloat64(), float64(self.GetInteger64()))), nil
					default:
						return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
					}
				},
			),
		),
	)
	object.Set(BitXor,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					right := arguments[0]
					if _, ok := right.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName)
					}
					return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()^right.GetInteger64()), nil
				},
			),
		),
	)
	object.Set(RightBitXor,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					left := arguments[0]
					if _, ok := left.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName)
					}
					return p.NewInteger(p.PeekSymbolTable(), left.GetInteger64()^self.GetInteger64()), nil
				},
			),
		),
	)
	object.Set(BitAnd,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					right := arguments[0]
					if _, ok := right.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName)
					}
					return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()&right.GetInteger64()), nil
				},
			),
		),
	)
	object.Set(RightBitAnd,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					left := arguments[0]
					if _, ok := left.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName)
					}
					return p.NewInteger(p.PeekSymbolTable(), left.GetInteger64()&self.GetInteger64()), nil
				},
			),
		),
	)
	object.Set(BitOr,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					right := arguments[0]
					if _, ok := right.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName)
					}
					return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()|right.GetInteger64()), nil
				},
			),
		),
	)
	object.Set(RightBitOr,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					left := arguments[0]
					if _, ok := left.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName)
					}
					return p.NewInteger(p.PeekSymbolTable(), left.GetInteger64()|self.GetInteger64()), nil
				},
			),
		),
	)
	object.Set(BitLeft,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					right := arguments[0]
					if _, ok := right.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName)
					}
					return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()<<right.GetInteger64()), nil
				},
			),
		),
	)
	object.Set(RightBitLeft,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					left := arguments[0]
					if _, ok := left.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName)
					}
					return p.NewInteger(p.PeekSymbolTable(), left.GetInteger64()<<self.GetInteger64()), nil
				},
			),
		),
	)
	object.Set(BitRight,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					right := arguments[0]
					if _, ok := right.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(right.TypeName(), IntegerName)
					}
					return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()>>right.GetInteger64()), nil
				},
			),
		),
	)
	object.Set(RightBitRight,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 1,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					left := arguments[0]
					if _, ok := left.(*Integer); !ok {
						return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName)
					}
					return p.NewInteger(p.PeekSymbolTable(), left.GetInteger64()>>self.GetInteger64()), nil
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
					return p.NewBool(p.PeekSymbolTable(), float64(self.GetInteger64()) == floatRight), nil
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
					return p.NewBool(p.PeekSymbolTable(), floatLeft == float64(self.GetInteger64())), nil
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
					return p.NewBool(p.PeekSymbolTable(), float64(self.GetInteger64()) != floatRight), nil
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
					return p.NewBool(p.PeekSymbolTable(), floatLeft != float64(self.GetInteger64())), nil
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
					return p.NewBool(p.PeekSymbolTable(), float64(self.GetInteger64()) > floatRight), nil
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
					return p.NewBool(p.PeekSymbolTable(), floatLeft > float64(self.GetInteger64())), nil
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
					return p.NewBool(p.PeekSymbolTable(), float64(self.GetInteger64()) < floatRight), nil
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
					return p.NewBool(p.PeekSymbolTable(), floatLeft < float64(self.GetInteger64())), nil
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
					return p.NewBool(p.PeekSymbolTable(), float64(self.GetInteger64()) >= floatRight), nil
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
					return p.NewBool(p.PeekSymbolTable(), floatLeft >= float64(self.GetInteger64())), nil
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
					return p.NewBool(p.PeekSymbolTable(), float64(self.GetInteger64()) <= floatRight), nil
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
					return p.NewBool(p.PeekSymbolTable(), floatLeft <= float64(self.GetInteger64())), nil
				},
			),
		),
	)
	object.Set(Hash,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()), nil
				},
			),
		),
	)
	object.Set(Copy,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()), nil
				},
			),
		),
	)
	object.Set(ToInteger,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewInteger(p.PeekSymbolTable(), self.GetInteger64()), nil
				},
			),
		),
	)
	object.Set(ToFloat,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewFloat(p.PeekSymbolTable(), float64(self.GetInteger64())), nil
				},
			),
		),
	)
	object.Set(ToString,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewString(p.PeekSymbolTable(), fmt.Sprint(self.GetInteger64())), nil
				},
			),
		),
	)
	object.Set(ToBool,
		p.NewFunction(object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewBool(p.PeekSymbolTable(), self.GetInteger64() != 0), nil
				},
			),
		),
	)
	return nil
}
