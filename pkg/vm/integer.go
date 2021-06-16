package vm

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/tools"
	"math"
)

type Integer struct {
	*Object
}

func (p *Plasma) NewInteger(isBuiltIn bool, parentSymbols *SymbolTable, value int64) *Integer {
	integer := &Integer{
		p.NewObject(isBuiltIn, IntegerName, nil, parentSymbols),
	}
	integer.SetInteger64(value)
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
						return p.NewInteger(false, p.PeekSymbolTable(), ^self.GetInteger64()), nil
					},
				),
			),
		)
		object.Set(Negative,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewInteger(false, p.PeekSymbolTable(), -self.GetInteger64()), nil
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
							return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger64()+right.GetInteger64()), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), float64(self.GetInteger64())+right.GetFloat64()), nil
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
							return p.NewInteger(false, p.PeekSymbolTable(), left.GetInteger64()+self.GetInteger64()), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), left.GetFloat64()+float64(self.GetInteger64())), nil
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
							return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger64()-right.GetInteger64()), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), float64(self.GetInteger64())-right.GetFloat64()), nil
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
							return p.NewInteger(false, p.PeekSymbolTable(), left.GetInteger64()-self.GetInteger64()), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), left.GetFloat64()-float64(self.GetInteger64())), nil
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
							return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger64()*right.GetInteger64()), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), float64(self.GetInteger64())*right.GetFloat64()), nil
						case *String:
							return p.NewString(false, p.PeekSymbolTable(), tools.Repeat(right.GetString(), self.GetInteger64())), nil
						case *Tuple:
							content, repetitionError := p.Repeat(right.GetContent(), int(self.GetInteger64()))
							if repetitionError != nil {
								return nil, repetitionError
							}
							return p.NewTuple(false, p.PeekSymbolTable(), content), nil
						case *Array:
							content, repetitionError := p.Repeat(right.GetContent(), int(self.GetInteger64()))
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
							return p.NewInteger(false, p.PeekSymbolTable(), left.GetInteger64()*self.GetInteger64()), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), left.GetFloat64()*float64(self.GetInteger64())), nil
						case *String:
							return p.NewString(false, p.PeekSymbolTable(), tools.Repeat(left.GetString(), self.GetInteger64())), nil
						case *Tuple:
							content, repetitionError := p.Repeat(left.GetContent(), int(self.GetInteger64()))
							if repetitionError != nil {
								return nil, repetitionError
							}
							return p.NewTuple(false, p.PeekSymbolTable(), content), nil
						case *Array:
							content, repetitionError := p.Repeat(left.GetContent(), int(self.GetInteger64()))
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
							return p.NewFloat(false, p.PeekSymbolTable(), float64(self.GetInteger64())/float64(right.GetInteger64())), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), float64(self.GetInteger64())/right.GetFloat64()), nil
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
							return p.NewFloat(false, p.PeekSymbolTable(), float64(left.GetInteger64())/float64(self.GetInteger64())), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), left.GetFloat64()/float64(self.GetInteger64())), nil
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
							return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger64()/right.GetInteger64()), nil
						case *Float:
							return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger64()/int64(right.GetFloat64())), nil
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
							return p.NewInteger(false, p.PeekSymbolTable(), left.GetInteger64()/self.GetInteger64()), nil
						case *Float:
							return p.NewInteger(false, p.PeekSymbolTable(), int64(left.GetFloat64())/self.GetInteger64()), nil
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
							return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger64()%right.GetInteger64()), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), math.Mod(float64(self.GetInteger64()), right.GetFloat64())), nil
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
							return p.NewInteger(false, p.PeekSymbolTable(), left.GetInteger64()%self.GetInteger64()), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), math.Mod(left.GetFloat64(), float64(self.GetInteger64()))), nil
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
							return p.NewFloat(false, p.PeekSymbolTable(), math.Pow(float64(self.GetInteger64()), float64(right.GetInteger64()))), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), math.Pow(float64(self.GetInteger64()), right.GetFloat64())), nil
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
							return p.NewFloat(false, p.PeekSymbolTable(), math.Pow(float64(left.GetInteger64()), float64(self.GetInteger64()))), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(), math.Pow(left.GetFloat64(), float64(self.GetInteger64()))), nil
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
						return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger64()^right.GetInteger64()), nil
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
						return p.NewInteger(false, p.PeekSymbolTable(), left.GetInteger64()^self.GetInteger64()), nil
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
						return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger64()&right.GetInteger64()), nil
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
						return p.NewInteger(false, p.PeekSymbolTable(), left.GetInteger64()&self.GetInteger64()), nil
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
						return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger64()|right.GetInteger64()), nil
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
						return p.NewInteger(false, p.PeekSymbolTable(), left.GetInteger64()|self.GetInteger64()), nil
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
						return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger64()<<right.GetInteger64()), nil
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
						return p.NewInteger(false, p.PeekSymbolTable(), left.GetInteger64()<<self.GetInteger64()), nil
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
						return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger64()>>right.GetInteger64()), nil
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
						return p.NewInteger(false, p.PeekSymbolTable(), left.GetInteger64()>>self.GetInteger64()), nil
					},
				),
			),
		)
		object.Set(Equals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
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
						return p.NewBool(false, p.PeekSymbolTable(), float64(self.GetInteger64()) == floatRight), nil
					},
				),
			),
		)
		object.Set(RightEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
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
						return p.NewBool(false, p.PeekSymbolTable(), floatLeft == float64(self.GetInteger64())), nil
					},
				),
			),
		)
		object.Set(NotEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
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
						return p.NewBool(false, p.PeekSymbolTable(), float64(self.GetInteger64()) != floatRight), nil
					},
				),
			),
		)
		object.Set(RightNotEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
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
						return p.NewBool(false, p.PeekSymbolTable(), floatLeft != float64(self.GetInteger64())), nil
					},
				),
			),
		)
		object.Set(GreaterThan,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
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
						return p.NewBool(false, p.PeekSymbolTable(), float64(self.GetInteger64()) > floatRight), nil
					},
				),
			),
		)
		object.Set(RightGreaterThan,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
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
						return p.NewBool(false, p.PeekSymbolTable(), floatLeft > float64(self.GetInteger64())), nil
					},
				),
			),
		)
		object.Set(LessThan,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
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
						return p.NewBool(false, p.PeekSymbolTable(), float64(self.GetInteger64()) < floatRight), nil
					},
				),
			),
		)
		object.Set(RightLessThan,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
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
						return p.NewBool(false, p.PeekSymbolTable(), floatLeft < float64(self.GetInteger64())), nil
					},
				),
			),
		)
		object.Set(GreaterThanOrEqual,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
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
						return p.NewBool(false, p.PeekSymbolTable(), float64(self.GetInteger64()) >= floatRight), nil
					},
				),
			),
		)
		object.Set(RightGreaterThanOrEqual,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
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
						return p.NewBool(false, p.PeekSymbolTable(), floatLeft >= float64(self.GetInteger64())), nil
					},
				),
			),
		)
		object.Set(LessThanOrEqual,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
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
						return p.NewBool(false, p.PeekSymbolTable(), float64(self.GetInteger64()) <= floatRight), nil
					},
				),
			),
		)
		object.Set(RightLessThanOrEqual,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self Value, arguments ...Value) (Value, *Object) {
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
						return p.NewBool(false, p.PeekSymbolTable(), floatLeft <= float64(self.GetInteger64())), nil
					},
				),
			),
		)
		object.Set(Hash,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger64()), nil
					},
				),
			),
		)
		object.Set(Copy,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger64()), nil
					},
				),
			),
		)
		object.Set(ToInteger,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger64()), nil
					},
				),
			),
		)
		object.Set(ToFloat,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewFloat(false, p.PeekSymbolTable(), float64(self.GetInteger64())), nil
					},
				),
			),
		)
		object.Set(ToString,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewString(false, p.PeekSymbolTable(), fmt.Sprint(self.GetInteger64())), nil
					},
				),
			),
		)
		object.Set(ToBool,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewBool(false, p.PeekSymbolTable(), self.GetInteger64() != 0), nil
					},
				),
			),
		)
		return nil
	}
}
