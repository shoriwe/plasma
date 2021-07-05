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
						return p.NewInteger(false, p.PeekSymbolTable(),
							^self.GetInteger(),
						), nil
					},
				),
			),
		)
		object.Set(Negative,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewInteger(false, p.PeekSymbolTable(),
							-self.GetInteger(),
						), nil
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
								self.GetInteger()+right.GetInteger(),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								float64(self.GetInteger())+right.GetFloat(),
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
								left.GetInteger()+self.GetInteger(),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								left.GetFloat()+float64(self.GetInteger()),
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
								self.GetInteger()-right.GetInteger(),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								float64(self.GetInteger())-right.GetFloat(),
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
								left.GetInteger()-self.GetInteger(),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								left.GetFloat()-float64(self.GetInteger()),
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
								self.GetInteger()*right.GetInteger(),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								float64(self.GetInteger())*right.GetFloat(),
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
								left.GetInteger()*self.GetInteger(),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								left.GetFloat()*float64(self.GetInteger()),
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
								float64(self.GetInteger())/float64(right.GetInteger()),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								float64(self.GetInteger())/right.GetFloat(),
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
								float64(left.GetInteger())/float64(self.GetInteger()),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								left.GetFloat()/float64(self.GetInteger()),
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
								self.GetInteger()/right.GetInteger(),
							), nil
						case *Float:
							return p.NewInteger(false, p.PeekSymbolTable(),
								self.GetInteger()/int64(right.GetFloat()),
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
								left.GetInteger()/self.GetInteger(),
							), nil
						case *Float:
							return p.NewInteger(false, p.PeekSymbolTable(),
								int64(left.GetFloat())/self.GetInteger(),
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
							return p.NewFloat(false, p.PeekSymbolTable(),
								math.Mod(float64(self.GetInteger()), float64(right.GetInteger())),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								math.Mod(float64(self.GetInteger()), right.GetFloat()),
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
								left.GetInteger()%self.GetInteger(),
							), nil
						case *Float:
							return p.NewInteger(false, p.PeekSymbolTable(),
								int64(left.GetFloat())%self.GetInteger(),
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
								math.Pow(float64(self.GetInteger()), float64(right.GetInteger())),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								math.Pow(float64(self.GetInteger()), right.GetFloat()),
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
								math.Pow(
									float64(left.GetInteger()),
									float64(self.GetInteger(),
									),
								),
							), nil
						case *Float:
							return p.NewFloat(false, p.PeekSymbolTable(),
								math.Pow(
									left.GetFloat(),
									float64(self.GetInteger(),
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
							self.GetInteger()^right.GetInteger(),
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
							left.GetInteger()^self.GetInteger(),
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
							self.GetInteger()&right.GetInteger(),
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
							left.GetInteger()&self.GetInteger(),
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
							self.GetInteger()|right.GetInteger(),
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
							left.GetInteger()|self.GetInteger(),
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
							self.GetInteger()<<uint(right.GetInteger()),
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
							left.GetInteger()<<uint(self.GetInteger()),
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
							self.GetInteger()>>uint(right.GetInteger()),
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
							left.GetInteger()>>uint(self.GetInteger()),
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
							if self.GetInteger() == right.GetInteger() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						case *Float:
							if float64(self.GetInteger()) == right.GetFloat() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
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
							if left.GetInteger() == self.GetInteger() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						case *Float:
							if left.GetFloat() == float64(self.GetInteger()) {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
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
							if self.GetInteger() != right.GetInteger() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						case *Float:
							if float64(self.GetInteger()) != (right.GetFloat()) {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
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
							if left.GetInteger() != self.GetInteger() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						case *Float:
							if left.GetFloat() != float64(self.GetInteger()) {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
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
							if self.GetInteger() > right.GetInteger() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						case *Float:
							if float64(self.GetInteger()) > right.GetFloat() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
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
							if left.GetInteger() > self.GetInteger() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						case *Float:
							if left.GetFloat() > float64(self.GetInteger()) {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
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
							if self.GetInteger() < right.GetInteger() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						case *Float:
							if float64(self.GetInteger()) < right.GetFloat() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
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
							if left.GetInteger() < self.GetInteger() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						case *Float:
							if left.GetFloat() < float64(self.GetInteger()) {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
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
							if self.GetInteger() >= right.GetInteger() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						case *Float:
							if float64(self.GetInteger()) >= right.GetFloat() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
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
							if left.GetInteger() >= self.GetInteger() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						case *Float:
							if left.GetFloat() >= float64(self.GetInteger()) {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
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
							if self.GetInteger() <= right.GetInteger() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						case *Float:
							if float64(self.GetInteger()) <= right.GetFloat() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
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
							if left.GetInteger() <= self.GetInteger() {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						case *Float:
							if left.GetFloat() <= float64(self.GetInteger()) {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
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
						return p.NewFloat(false, p.PeekSymbolTable(),
							float64(self.GetInteger()),
						), nil
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
						if self.GetInteger() != 0 {
							return p.GetTrue(), nil
						}
						return p.GetFalse(), nil
					},
				),
			),
		)
		return nil
	}
}
