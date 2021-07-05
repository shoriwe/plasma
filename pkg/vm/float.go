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
		object.Set(Negative,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewFloat(false, p.PeekSymbolTable(),
							-self.GetFloat(),
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
			),
		)
		object.Set(RightAdd,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)
		object.Set(Sub,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)
		object.Set(RightSub,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)
		object.Set(Mul,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)
		object.Set(RightMul,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)
		object.Set(RightMod,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)

		object.Set(Equals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)
		object.Set(RightEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)
		object.Set(NotEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)
		object.Set(RightNotEquals,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)
		object.Set(GreaterThan,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)
		object.Set(RightGreaterThan,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)
		object.Set(LessThan,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)
		object.Set(RightLessThan,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)
		object.Set(GreaterThanOrEqual,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)
		object.Set(RightGreaterThanOrEqual,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)
		object.Set(LessThanOrEqual,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)
		object.Set(RightLessThanOrEqual,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
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
			),
		)

		object.Set(Hash,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						if self.GetHash() == 0 {
							floatHash := p.HashString(fmt.Sprintf("%f-%s", self.GetFloat(), FloatName))
							self.SetHash(floatHash)
						}
						return p.NewInteger(false, p.PeekSymbolTable(), self.GetHash()), nil
					},
				),
			),
		)
		object.Set(Copy,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewFloat(false, p.PeekSymbolTable(), self.GetFloat()), nil
					},
				),
			),
		)

		object.Set(ToInteger,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewInteger(false, p.PeekSymbolTable(), int64(self.GetFloat())), nil
					},
				),
			),
		)
		object.Set(ToFloat,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewFloat(false, p.PeekSymbolTable(), self.GetFloat()), nil
					},
				),
			),
		)
		object.Set(ToString,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewString(false, p.PeekSymbolTable(), fmt.Sprint(self.GetFloat())), nil
					},
				),
			),
		)
		object.Set(ToBool,
			p.NewFunction(isBuiltIn, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						if self.GetFloat() != 0 {
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
