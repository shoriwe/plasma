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
	integer.SetOnDemandSymbol(Self,
		func() Value {
			return integer
		},
	)
	return integer
}

func (p *Plasma) IntegerInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(object Value) *Object {
		object.SetOnDemandSymbol(NegBits,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewInteger(false, p.PeekSymbolTable(),
								^self.GetInteger(),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Negative,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewInteger(false, p.PeekSymbolTable(),
								-self.GetInteger(),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Add,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(RightAdd,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(Sub,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(RightSub,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(Mul,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(RightMul,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(Div,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(RightDiv,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(FloorDiv,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(RightFloorDiv,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(Mod,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(RightMod,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(Pow,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(RightPow,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(BitXor,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(RightBitXor,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(BitAnd,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(RightBitAnd,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(BitOr,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(RightBitOr,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(BitLeft,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(RightBitLeft,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(BitRight,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(RightBitRight,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(Equals,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(RightEquals,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(NotEquals,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(RightNotEquals,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(GreaterThan,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(RightGreaterThan,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(LessThan,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(RightLessThan,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(GreaterThanOrEqual,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(RightGreaterThanOrEqual,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(LessThanOrEqual,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(RightLessThanOrEqual,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
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
				)
			},
		)
		object.SetOnDemandSymbol(Hash,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger()), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Copy,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger()), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToInteger,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewInteger(false, p.PeekSymbolTable(), self.GetInteger()), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToFloat,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewFloat(false, p.PeekSymbolTable(),
								float64(self.GetInteger()),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToString,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewString(false, p.PeekSymbolTable(), fmt.Sprint(self.GetInteger())), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToBool,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							if self.GetInteger() != 0 {
								return p.GetTrue(), nil
							}
							return p.GetFalse(), nil
						},
					),
				)
			},
		)
		return nil
	}
}
