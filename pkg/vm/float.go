package vm

import (
	"fmt"
	"math"
)

func (p *Plasma) NewFloat(context *Context, isBuiltIn bool, parentSymbols *SymbolTable, value float64) *Value {
	float_ := p.NewValue(context, isBuiltIn, FloatName, nil, parentSymbols)
	float_.BuiltInTypeId = FloatId
	float_.SetFloat(value)
	p.FloatInitialize(isBuiltIn)(context, float_)
	float_.SetOnDemandSymbol(Self,
		func() *Value {
			return float_
		},
	)
	return float_
}

func (p *Plasma) FloatInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object *Value) *Value {
		object.SetOnDemandSymbol(Negative,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.NewFloat(context, false, context.PeekSymbolTable(),
								-self.GetFloat(),
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Add,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									self.GetFloat()+float64(right.GetInteger()),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									self.GetFloat()+right.GetFloat(),
								), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightAdd,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									float64(left.GetInteger())+self.GetFloat(),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									left.GetFloat()+self.GetFloat(),
								), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Sub,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									self.GetFloat()-float64(right.GetInteger()),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									self.GetFloat()-right.GetFloat(),
								), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightSub,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									float64(left.GetInteger())-self.GetFloat(),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									left.GetFloat()-self.GetFloat(),
								), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Mul,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									self.GetFloat()*float64(right.GetInteger()),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									self.GetFloat()*right.GetFloat(),
								), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightMul,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									float64(left.GetInteger())*self.GetFloat(),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									left.GetFloat()*self.GetFloat(),
								), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Div,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									self.GetFloat()/float64(right.GetInteger()),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									self.GetFloat()/right.GetFloat(),
								), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightDiv,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									float64(left.GetInteger())/self.GetFloat(),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									left.GetFloat()/self.GetFloat(),
								), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Mod,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									math.Mod(self.GetFloat(), float64(right.GetInteger())),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									math.Mod(self.GetFloat(), right.GetFloat()),
								), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightMod,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									math.Mod(float64(left.GetInteger()), self.GetFloat()),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									math.Mod(left.GetFloat(), self.GetFloat()),
								), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Pow,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									math.Pow(
										self.GetFloat(),
										float64(right.GetInteger()),
									),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									math.Pow(
										self.GetFloat(),
										right.GetFloat(),
									),
								), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightPow,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									math.Pow(
										float64(left.GetInteger()),
										self.GetFloat(),
									),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									math.Pow(
										float64(left.GetInteger()),
										self.GetFloat(),
									),
								), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)

		object.SetOnDemandSymbol(Equals,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(self.GetFloat() == float64(right.GetInteger())), true
							case FloatId:
								return p.InterpretAsBool(self.GetFloat() == right.GetFloat()), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightEquals,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(float64(left.GetInteger()) == self.GetFloat()), true
							case FloatId:
								return p.InterpretAsBool(left.GetFloat() == self.GetFloat()), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(NotEquals,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(self.GetFloat() != float64(right.GetInteger())), true
							case FloatId:
								return p.InterpretAsBool(self.GetFloat() != right.GetFloat()), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightNotEquals,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(float64(left.GetInteger()) != self.GetFloat()), true
							case FloatId:
								return p.InterpretAsBool(left.GetFloat() != self.GetFloat()), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(GreaterThan,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(self.GetFloat() > float64(right.GetInteger())), true
							case FloatId:
								return p.InterpretAsBool(self.GetFloat() > right.GetFloat()), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightGreaterThan,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(float64(left.GetInteger()) > self.GetFloat()), true
							case FloatId:
								return p.InterpretAsBool(left.GetFloat() > self.GetFloat()), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(LessThan,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(self.GetFloat() < float64(right.GetInteger())), true
							case FloatId:
								return p.InterpretAsBool(self.GetFloat() < right.GetFloat()), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightLessThan,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(float64(left.GetInteger()) < self.GetFloat()), true
							case FloatId:
								return p.InterpretAsBool(left.GetFloat() < self.GetFloat()), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(GreaterThanOrEqual,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(self.GetFloat() >= float64(right.GetInteger())), true
							case FloatId:
								return p.InterpretAsBool(self.GetFloat() >= right.GetFloat()), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightGreaterThanOrEqual,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(float64(left.GetInteger()) >= self.GetFloat()), true
							case FloatId:
								return p.InterpretAsBool(left.GetFloat() >= self.GetFloat()), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(LessThanOrEqual,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(self.GetFloat() <= float64(right.GetInteger())), true
							case FloatId:
								return p.InterpretAsBool(self.GetFloat() <= right.GetFloat()), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightLessThanOrEqual,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(float64(left.GetInteger()) <= self.GetFloat()), true
							case FloatId:
								return p.InterpretAsBool(left.GetFloat() <= self.GetFloat()), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)

		object.SetOnDemandSymbol(Hash,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							if self.GetHash() == 0 {
								floatHash := p.HashString(fmt.Sprintf("%f-%s", self.GetFloat(), FloatName))
								self.SetHash(floatHash)
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(), self.GetHash()), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Copy,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.NewFloat(context, false, context.PeekSymbolTable(), self.GetFloat()), true
						},
					),
				)
			},
		)

		object.SetOnDemandSymbol(ToInteger,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.NewInteger(context, false, context.PeekSymbolTable(), int64(self.GetFloat())), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToFloat,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.NewFloat(context, false, context.PeekSymbolTable(), self.GetFloat()), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToString,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.NewString(context, false, context.PeekSymbolTable(), fmt.Sprint(self.GetFloat())), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToBool,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.InterpretAsBool(self.GetFloat() != 0), true
						},
					),
				)
			},
		)
		return nil
	}
}
