package vm

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/tools"
	"math"
)

func (p *Plasma) NewInteger(context *Context, isBuiltIn bool, parentSymbols *SymbolTable, value int64) *Value {
	integer := p.NewValue(context, isBuiltIn, IntegerName, nil, parentSymbols)
	integer.BuiltInTypeId = IntegerId
	integer.SetInteger(value)
	p.IntegerInitialize(isBuiltIn)(context, integer)
	integer.SetOnDemandSymbol(Self,
		func() *Value {
			return integer
		},
	)
	return integer
}

func (p *Plasma) IntegerInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object *Value) *Value {
		object.SetOnDemandSymbol(NegBits,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								^self.GetInteger(),
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Negative,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								-self.GetInteger(),
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
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									self.GetInteger()+right.GetInteger(),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									float64(self.GetInteger())+right.GetFloat(),
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
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									left.GetInteger()+self.GetInteger(),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									left.GetFloat()+float64(self.GetInteger()),
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
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									self.GetInteger()-right.GetInteger(),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									float64(self.GetInteger())-right.GetFloat(),
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
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									left.GetInteger()-self.GetInteger(),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									left.GetFloat()-float64(self.GetInteger()),
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
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									self.GetInteger()*right.GetInteger(),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									float64(self.GetInteger())*right.GetFloat(),
								), true
							case StringId:
								return p.NewString(context, false, context.PeekSymbolTable(), tools.Repeat(right.GetString(), self.GetInteger())), true
							case BytesId:
								content, repetitionError := p.Repeat(context, right.GetContent(), self.GetInteger())
								if repetitionError != nil {
									return repetitionError, false
								}
								return p.NewTuple(context, false, context.PeekSymbolTable(), content), true
							case ArrayId:
								content, repetitionError := p.Repeat(context, right.GetContent(), self.GetInteger())
								if repetitionError != nil {
									return repetitionError, false
								}
								return p.NewArray(context, false, context.PeekSymbolTable(), content), true
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
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									left.GetInteger()*self.GetInteger(),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									left.GetFloat()*float64(self.GetInteger()),
								), true
							case StringId:
								return p.NewString(context, false, context.PeekSymbolTable(), tools.Repeat(left.GetString(), self.GetInteger())), true
							case BytesId:
								content, repetitionError := p.Repeat(context, left.GetContent(), self.GetInteger())
								if repetitionError != nil {
									return repetitionError, false
								}
								return p.NewTuple(context, false, context.PeekSymbolTable(), content), true
							case ArrayId:
								content, repetitionError := p.Repeat(context, left.GetContent(), self.GetInteger())
								if repetitionError != nil {
									return repetitionError, false
								}
								return p.NewArray(context, false, context.PeekSymbolTable(), content), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName), false
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
									float64(self.GetInteger())/float64(right.GetInteger()),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									float64(self.GetInteger())/right.GetFloat(),
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
									float64(left.GetInteger())/float64(self.GetInteger()),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									left.GetFloat()/float64(self.GetInteger()),
								), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(FloorDiv,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									self.GetInteger()/right.GetInteger(),
								), true
							case FloatId:
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									self.GetInteger()/int64(right.GetFloat()),
								), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightFloorDiv,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									left.GetInteger()/self.GetInteger(),
								), true
							case FloatId:
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									int64(left.GetFloat())/self.GetInteger(),
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
									math.Mod(float64(self.GetInteger()), float64(right.GetInteger())),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									math.Mod(float64(self.GetInteger()), right.GetFloat()),
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
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									left.GetInteger()%self.GetInteger(),
								), true
							case FloatId:
								return p.NewInteger(context, false, context.PeekSymbolTable(),
									int64(left.GetFloat())%self.GetInteger(),
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
									math.Pow(float64(self.GetInteger()), float64(right.GetInteger())),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									math.Pow(float64(self.GetInteger()), right.GetFloat()),
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
										float64(self.GetInteger()),
									),
								), true
							case FloatId:
								return p.NewFloat(context, false, context.PeekSymbolTable(),
									math.Pow(
										left.GetFloat(),
										float64(self.GetInteger()),
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
		object.SetOnDemandSymbol(BitXor,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							if !right.IsTypeById(BytesId) {
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName), false
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								self.GetInteger()^right.GetInteger(),
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightBitXor,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							if !left.IsTypeById(IntegerId) {
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName), false
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								left.GetInteger()^self.GetInteger(),
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(BitAnd,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							if !right.IsTypeById(BytesId) {
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName), false
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								self.GetInteger()&right.GetInteger(),
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightBitAnd,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							if !left.IsTypeById(IntegerId) {
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName), false
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								left.GetInteger()&self.GetInteger(),
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(BitOr,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							if !right.IsTypeById(BytesId) {
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName), false
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								self.GetInteger()|right.GetInteger(),
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightBitOr,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							if !left.IsTypeById(IntegerId) {
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName), false
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								left.GetInteger()|self.GetInteger(),
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(BitLeft,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							if !right.IsTypeById(BytesId) {
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName), false
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								self.GetInteger()<<uint(right.GetInteger()),
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightBitLeft,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							if !left.IsTypeById(IntegerId) {
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName), false
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								left.GetInteger()<<uint(self.GetInteger()),
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(BitRight,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							if !right.IsTypeById(BytesId) {
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName), false
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								self.GetInteger()>>uint(right.GetInteger()),
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightBitRight,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							if !left.IsTypeById(IntegerId) {
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName), false
							}
							return p.NewInteger(context, false, context.PeekSymbolTable(),
								left.GetInteger()>>uint(self.GetInteger()),
							), true
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
								return p.InterpretAsBool(self.GetInteger() == right.GetInteger()), true
							case FloatId:
								return p.InterpretAsBool(float64(self.GetInteger()) == right.GetFloat()), true
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
								return p.InterpretAsBool(left.GetInteger() == self.GetInteger()), true
							case FloatId:
								return p.InterpretAsBool(left.GetFloat() == float64(self.GetInteger())), true
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
								return p.InterpretAsBool(self.GetInteger() != right.GetInteger()), true
							case FloatId:
								return p.InterpretAsBool(float64(self.GetInteger()) != (right.GetFloat())), true
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
								return p.InterpretAsBool(left.GetInteger() != self.GetInteger()), true
							case FloatId:
								return p.InterpretAsBool(left.GetFloat() != float64(self.GetInteger())), true
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
								return p.InterpretAsBool(self.GetInteger() > right.GetInteger()), true
							case FloatId:
								return p.InterpretAsBool(float64(self.GetInteger()) > right.GetFloat()), true
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
								return p.InterpretAsBool(left.GetInteger() > self.GetInteger()), true
							case FloatId:
								return p.InterpretAsBool(left.GetFloat() > float64(self.GetInteger())), true
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
								return p.InterpretAsBool(self.GetInteger() < right.GetInteger()), true
							case FloatId:
								return p.InterpretAsBool(float64(self.GetInteger()) < right.GetFloat()), true
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
								return p.InterpretAsBool(left.GetInteger() < self.GetInteger()), true
							case FloatId:
								return p.InterpretAsBool(left.GetFloat() < float64(self.GetInteger())), true
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
								return p.InterpretAsBool(self.GetInteger() >= right.GetInteger()), true
							case FloatId:
								return p.InterpretAsBool(float64(self.GetInteger()) >= right.GetFloat()), true
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
								return p.InterpretAsBool(left.GetInteger() >= self.GetInteger()), true
							case FloatId:
								return p.InterpretAsBool(left.GetFloat() >= float64(self.GetInteger())), true
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
								return p.InterpretAsBool(self.GetInteger() <= right.GetInteger()), true
							case FloatId:
								return p.InterpretAsBool(float64(self.GetInteger()) <= right.GetFloat()), true
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
								return p.InterpretAsBool(left.GetInteger() <= self.GetInteger()), true
							case FloatId:
								return p.InterpretAsBool(left.GetFloat() <= float64(self.GetInteger())), true
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
							return p.NewInteger(context, false, context.PeekSymbolTable(), self.GetInteger()), true
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
							return p.NewInteger(context, false, context.PeekSymbolTable(), self.GetInteger()), true
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
							return p.NewInteger(context, false, context.PeekSymbolTable(), self.GetInteger()), true
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
							return p.NewFloat(context, false, context.PeekSymbolTable(),
								float64(self.GetInteger()),
							), true
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
							return p.NewString(context, false, context.PeekSymbolTable(), fmt.Sprint(self.GetInteger())), true
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
							return p.InterpretAsBool(self.GetInteger() != 0), true
						},
					),
				)
			},
		)
		return nil
	}
}
