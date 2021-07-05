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
		p.NewObject(isBuiltIn, FloatName, nil, parentSymbols),
	}
	float_.SetFloat(value)
	p.FloatInitialize(isBuiltIn)(float_)
	float_.Set(Self, float_)
	return float_
}

func (p *Plasma) FloatInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(object Value) *Object {
		object.SetOnDemandSymbol(Negative,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self Value, _ ...Value) (Value, *Object) {
							return p.NewFloat(false, p.PeekSymbolTable(),
								-self.GetFloat(),
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
								return p.NewFloat(false, p.PeekSymbolTable(),
									self.GetFloat()+float64(right.GetInteger()),
								), nil
							case *Float:
								return p.NewFloat(false, p.PeekSymbolTable(),
									self.GetFloat()+right.GetFloat(),
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
								return p.NewFloat(false, p.PeekSymbolTable(),
									float64(left.GetInteger())+self.GetFloat(),
								), nil
							case *Float:
								return p.NewFloat(false, p.PeekSymbolTable(),
									left.GetFloat()+self.GetFloat(),
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
								return p.NewFloat(false, p.PeekSymbolTable(),
									self.GetFloat()-float64(right.GetInteger()),
								), nil
							case *Float:
								return p.NewFloat(false, p.PeekSymbolTable(),
									self.GetFloat()-right.GetFloat(),
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
								return p.NewFloat(false, p.PeekSymbolTable(),
									float64(left.GetInteger())-self.GetFloat(),
								), nil
							case *Float:
								return p.NewFloat(false, p.PeekSymbolTable(),
									left.GetFloat()-self.GetFloat(),
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
								return p.NewFloat(false, p.PeekSymbolTable(),
									self.GetFloat()*float64(right.GetInteger()),
								), nil
							case *Float:
								return p.NewFloat(false, p.PeekSymbolTable(),
									self.GetFloat()*right.GetFloat(),
								), nil
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
								return p.NewFloat(false, p.PeekSymbolTable(),
									float64(left.GetInteger())*self.GetFloat(),
								), nil
							case *Float:
								return p.NewFloat(false, p.PeekSymbolTable(),
									left.GetFloat()*self.GetFloat(),
								), nil
							default:
								return nil, p.NewInvalidTypeError(left.TypeName(), IntegerName, FloatName)
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
									self.GetFloat()/float64(right.GetInteger()),
								), nil
							case *Float:
								return p.NewFloat(false, p.PeekSymbolTable(),
									self.GetFloat()/right.GetFloat(),
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
									float64(left.GetInteger())/self.GetFloat(),
								), nil
							case *Float:
								return p.NewFloat(false, p.PeekSymbolTable(),
									left.GetFloat()/self.GetFloat(),
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
									math.Mod(self.GetFloat(), float64(right.GetInteger())),
								), nil
							case *Float:
								return p.NewFloat(false, p.PeekSymbolTable(),
									math.Mod(self.GetFloat(), right.GetFloat()),
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
								return p.NewFloat(false, p.PeekSymbolTable(),
									math.Mod(float64(left.GetInteger()), self.GetFloat()),
								), nil
							case *Float:
								return p.NewFloat(false, p.PeekSymbolTable(),
									math.Mod(left.GetFloat(), self.GetFloat()),
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
									math.Pow(
										self.GetFloat(),
										float64(right.GetInteger()),
									),
								), nil
							case *Float:
								return p.NewFloat(false, p.PeekSymbolTable(),
									math.Pow(
										self.GetFloat(),
										right.GetFloat(),
									),
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
										self.GetFloat(),
									),
								), nil
							case *Float:
								return p.NewFloat(false, p.PeekSymbolTable(),
									math.Pow(
										float64(left.GetInteger()),
										self.GetFloat(),
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

		object.SetOnDemandSymbol(Equals,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self Value, arguments ...Value) (Value, *Object) {
							right := arguments[0]
							switch right.(type) {
							case *Integer:
								if self.GetFloat() == float64(right.GetInteger()) {
									return p.GetTrue(), nil
								}
								return p.GetFalse(), nil
							case *Float:
								if self.GetFloat() == right.GetFloat() {
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
								if float64(left.GetInteger()) == self.GetFloat() {
									return p.GetTrue(), nil
								}
								return p.GetFalse(), nil
							case *Float:
								if left.GetFloat() == self.GetFloat() {
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
								if self.GetFloat() != float64(right.GetInteger()) {
									return p.GetTrue(), nil
								}
								return p.GetFalse(), nil
							case *Float:
								if self.GetFloat() != right.GetFloat() {
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
								if float64(left.GetInteger()) != self.GetFloat() {
									return p.GetTrue(), nil
								}
								return p.GetFalse(), nil
							case *Float:
								if left.GetFloat() != self.GetFloat() {
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
								if self.GetFloat() > float64(right.GetInteger()) {
									return p.GetTrue(), nil
								}
								return p.GetFalse(), nil
							case *Float:
								if self.GetFloat() > right.GetFloat() {
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
								if float64(left.GetInteger()) > self.GetFloat() {
									return p.GetTrue(), nil
								}
								return p.GetFalse(), nil
							case *Float:
								if left.GetFloat() > self.GetFloat() {
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
								if self.GetFloat() < float64(right.GetInteger()) {
									return p.GetTrue(), nil
								}
								return p.GetFalse(), nil
							case *Float:
								if self.GetFloat() < right.GetFloat() {
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
								if float64(left.GetInteger()) < self.GetFloat() {
									return p.GetTrue(), nil
								}
								return p.GetFalse(), nil
							case *Float:
								if left.GetFloat() < self.GetFloat() {
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
								if self.GetFloat() >= float64(right.GetInteger()) {
									return p.GetTrue(), nil
								}
								return p.GetFalse(), nil
							case *Float:
								if self.GetFloat() >= right.GetFloat() {
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
								if float64(left.GetInteger()) >= self.GetFloat() {
									return p.GetTrue(), nil
								}
								return p.GetFalse(), nil
							case *Float:
								if left.GetFloat() >= self.GetFloat() {
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
								if self.GetFloat() <= float64(right.GetInteger()) {
									return p.GetTrue(), nil
								}
								return p.GetFalse(), nil
							case *Float:
								if self.GetFloat() <= right.GetFloat() {
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
								if float64(left.GetInteger()) <= self.GetFloat() {
									return p.GetTrue(), nil
								}
								return p.GetFalse(), nil
							case *Float:
								if left.GetFloat() <= self.GetFloat() {
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
							if self.GetHash() == 0 {
								floatHash := p.HashString(fmt.Sprintf("%f-%s", self.GetFloat(), FloatName))
								self.SetHash(floatHash)
							}
							return p.NewInteger(false, p.PeekSymbolTable(), self.GetHash()), nil
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
							return p.NewFloat(false, p.PeekSymbolTable(), self.GetFloat()), nil
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
							return p.NewInteger(false, p.PeekSymbolTable(), int64(self.GetFloat())), nil
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
							return p.NewFloat(false, p.PeekSymbolTable(), self.GetFloat()), nil
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
							return p.NewString(false, p.PeekSymbolTable(), fmt.Sprint(self.GetFloat())), nil
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
							if self.GetFloat() != 0 {
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
