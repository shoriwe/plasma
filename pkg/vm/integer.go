package vm

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/tools"
	"math"
)

type Integer struct {
	*Object
}

func (p *Plasma) NewInteger(context *Context, isBuiltIn bool, parentSymbols *SymbolTable, value int64) *Integer {
	integer := &Integer{
		p.NewObject(context, isBuiltIn, IntegerName, nil, parentSymbols),
	}
	integer.SetInteger(value)
	p.IntegerInitialize(isBuiltIn)(context, integer)
	integer.SetOnDemandSymbol(Self,
		func() Value {
			return integer
		},
	)
	return integer
}

func (p *Plasma) IntegerInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object Value) *Object {
		object.SetOnDemandSymbol(NegBits,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								^self.GetInteger(),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Negative,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								-self.GetInteger(),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Add,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							switch right.(type) {
							case *Integer:
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									self.GetInteger()+right.GetInteger(),
								), nil
							case *Float:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									float64(self.GetInteger())+right.GetFloat(),
								), nil
							default:
								return nil, p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightAdd,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							switch left.(type) {
							case *Integer:
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									left.GetInteger()+self.GetInteger(),
								), nil
							case *Float:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									left.GetFloat()+float64(self.GetInteger()),
								), nil
							default:
								return nil, p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Sub,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							switch right.(type) {
							case *Integer:
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									self.GetInteger()-right.GetInteger(),
								), nil
							case *Float:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									float64(self.GetInteger())-right.GetFloat(),
								), nil
							default:
								return nil, p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightSub,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							switch left.(type) {
							case *Integer:
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									left.GetInteger()-self.GetInteger(),
								), nil
							case *Float:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									left.GetFloat()-float64(self.GetInteger()),
								), nil
							default:
								return nil, p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Mul,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							switch right.(type) {
							case *Integer:
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									self.GetInteger()*right.GetInteger(),
								), nil
							case *Float:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									float64(self.GetInteger())*right.GetFloat(),
								), nil
							case *String:
								return p.NewString(context, false, context.PeekSymbolTable(), tools.Repeat(right.GetString(), self.GetInteger())), nil
							case *Tuple:
								content, repetitionError := p.Repeat(context, right.GetContent(), self.GetInteger())
								if repetitionError != nil {
									return nil, repetitionError
								}
								return p.NewTuple(context, false, context.PeekSymbolTable(), content), nil
							case *Array:
								content, repetitionError := p.Repeat(context, right.GetContent(), self.GetInteger())
								if repetitionError != nil {
									return nil, repetitionError
								}
								return p.NewArray(context, false, context.PeekSymbolTable(), content), nil
							default:
								return nil, p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightMul,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							switch left.(type) {
							case *Integer:
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									left.GetInteger()*self.GetInteger(),
								), nil
							case *Float:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									left.GetFloat()*float64(self.GetInteger()),
								), nil
							case *String:
								return p.NewString(context, false, context.PeekSymbolTable(), tools.Repeat(left.GetString(), self.GetInteger())), nil
							case *Tuple:
								content, repetitionError := p.Repeat(context, left.GetContent(), self.GetInteger())
								if repetitionError != nil {
									return nil, repetitionError
								}
								return p.NewTuple(context, false, context.PeekSymbolTable(), content), nil
							case *Array:
								content, repetitionError := p.Repeat(context, left.GetContent(), self.GetInteger())
								if repetitionError != nil {
									return nil, repetitionError
								}
								return p.NewArray(context, false, context.PeekSymbolTable(), content), nil
							default:
								return nil, p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Div,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							switch right.(type) {
							case *Integer:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									float64(self.GetInteger())/float64(right.GetInteger()),
								), nil
							case *Float:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									float64(self.GetInteger())/right.GetFloat(),
								), nil
							default:
								return nil, p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightDiv,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							switch left.(type) {
							case *Integer:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									float64(left.GetInteger())/float64(self.GetInteger()),
								), nil
							case *Float:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									left.GetFloat()/float64(self.GetInteger()),
								), nil
							default:
								return nil, p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(FloorDiv,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							switch right.(type) {
							case *Integer:
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									self.GetInteger()/right.GetInteger(),
								), nil
							case *Float:
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									self.GetInteger()/int64(right.GetFloat()),
								), nil
							default:
								return nil, p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightFloorDiv,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							switch left.(type) {
							case *Integer:
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									left.GetInteger()/self.GetInteger(),
								), nil
							case *Float:
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									int64(left.GetFloat())/self.GetInteger(),
								), nil
							default:
								return nil, p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Mod,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							switch right.(type) {
							case *Integer:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									math.Mod(float64(self.GetInteger()), float64(right.GetInteger())),
								), nil
							case *Float:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									math.Mod(float64(self.GetInteger()), right.GetFloat()),
								), nil
							default:
								return nil, p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightMod,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							switch left.(type) {
							case *Integer:
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									left.GetInteger()%self.GetInteger(),
								), nil
							case *Float:
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									int64(left.GetFloat())%self.GetInteger(),
								), nil
							default:
								return nil, p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Pow,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							switch right.(type) {
							case *Integer:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									math.Pow(float64(self.GetInteger()), float64(right.GetInteger())),
								), nil
							case *Float:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									math.Pow(float64(self.GetInteger()), right.GetFloat()),
								), nil
							default:
								return nil, p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightPow,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							switch left.(type) {
							case *Integer:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									math.Pow(
										float64(left.GetInteger()),
										float64(self.GetInteger(),
										),
									),
								), nil
							case *Float:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									math.Pow(
										left.GetFloat(),
										float64(self.GetInteger(),
										),
									),
								), nil
							default:
								return nil, p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(BitXor,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							if _, ok := right.(*Integer); !ok {
								return nil, p.NewInvalidTypeError(context, right.TypeName(), IntegerName)
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								self.GetInteger()^right.GetInteger(),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightBitXor,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							if _, ok := left.(*Integer); !ok {
								return nil, p.NewInvalidTypeError(context, left.TypeName(), IntegerName)
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								left.GetInteger()^self.GetInteger(),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(BitAnd,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							if _, ok := right.(*Integer); !ok {
								return nil, p.NewInvalidTypeError(context, right.TypeName(), IntegerName)
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								self.GetInteger()&right.GetInteger(),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightBitAnd,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							if _, ok := left.(*Integer); !ok {
								return nil, p.NewInvalidTypeError(context, left.TypeName(), IntegerName)
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								left.GetInteger()&self.GetInteger(),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(BitOr,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							if _, ok := right.(*Integer); !ok {
								return nil, p.NewInvalidTypeError(context, right.TypeName(), IntegerName)
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								self.GetInteger()|right.GetInteger(),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightBitOr,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							if _, ok := left.(*Integer); !ok {
								return nil, p.NewInvalidTypeError(context, left.TypeName(), IntegerName)
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								left.GetInteger()|self.GetInteger(),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(BitLeft,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							if _, ok := right.(*Integer); !ok {
								return nil, p.NewInvalidTypeError(context, right.TypeName(), IntegerName)
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								self.GetInteger()<<uint(right.GetInteger()),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightBitLeft,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							if _, ok := left.(*Integer); !ok {
								return nil, p.NewInvalidTypeError(context, left.TypeName(), IntegerName)
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								left.GetInteger()<<uint(self.GetInteger()),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(BitRight,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							if _, ok := right.(*Integer); !ok {
								return nil, p.NewInvalidTypeError(context, right.TypeName(), IntegerName)
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								self.GetInteger()>>uint(right.GetInteger()),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightBitRight,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							left := arguments[0]
							if _, ok := left.(*Integer); !ok {
								return nil, p.NewInvalidTypeError(context, left.TypeName(), IntegerName)
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								left.GetInteger()>>uint(self.GetInteger()),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Equals,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
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
								return nil, p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightEquals,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
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
								return nil, p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(NotEquals,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
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
								return nil, p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightNotEquals,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
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
								return nil, p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(GreaterThan,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
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
								return nil, p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightGreaterThan,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
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
								return nil, p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(LessThan,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
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
								return nil, p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightLessThan,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
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
								return nil, p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(GreaterThanOrEqual,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
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
								return nil, p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightGreaterThanOrEqual,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
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
								return nil, p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(LessThanOrEqual,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
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
								return nil, p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightLessThanOrEqual,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
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
								return nil, p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName)
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Hash,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewInteger(context, false, context.PeekSymbolTable(), self.GetInteger()), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Copy,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewInteger(context, false, context.PeekSymbolTable(), self.GetInteger()), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToInteger,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewInteger(context, false, context.PeekSymbolTable(), self.GetInteger()), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToFloat,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewFloat(context, false, context.PeekSymbolTable(),
								float64(self.GetInteger()),
							), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToString,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewString(context, false, context.PeekSymbolTable(), fmt.Sprint(self.GetInteger())), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToBool,
			func() Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
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
