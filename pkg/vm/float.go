package vm

import (
	"fmt"
	magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"
	special_symbols "github.com/shoriwe/gplasma/pkg/common/special-symbols"
	"math"
)

/*
FloatValue
	Class: Float TODO
	Methods:
	- Positive
	- Negative
	- Equals: Float == Any
	- NotEqual: Float != Any
	- GreaterThan: Float > Float, Float > Integer
	- GreaterOrEqualThan: Float >= Float, Float >= Integer
	- LessThan: Float < Float, Float < Integer
	- LessOrEqualThan: Float <= Float, Float <= Integer
	- Add: Float + Float, Float + Integer
	- Sub: Float - Float, Float - Integer
	- Mul: Float * Float, Float * Integer
	- Div: Float / Float, Float / Integer
	- FloorDiv: Float // Float, Float // Integer
	- PowerOf: Float ** Float, Float ** Integer
	- Bool
	- Class
	- Copy
	- String
*/
func (ctx *Context) FloatValue(i float64) *Value {
	value := ctx.NewValue()
	value.Float = i
	value.OnDemand[magic_functions.Positive] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				return self.Copy(), nil
			},
		))
	}
	value.OnDemand[magic_functions.Negative] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				return ctx.FloatValue(-self.GetFloat()), nil
			},
		))
	}
	value.OnDemand[magic_functions.Equals] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				floatClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Float)
				switch otherClass {
				case integerClass:
					if self.GetFloat() == float64(argument[0].GetInt()) {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				case floatClass:
					if self.GetFloat() == argument[0].GetFloat() {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				default:
					return ctx.FalseValue(), nil
				}
			},
		))
	}
	value.OnDemand[magic_functions.NotEqual] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				floatClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Float)
				switch otherClass {
				case integerClass:
					if self.GetFloat() != float64(argument[0].GetInt()) {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				case floatClass:
					if self.GetFloat() != argument[0].GetFloat() {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				default:
					return ctx.FalseValue(), nil
				}
			},
		))
	}
	value.OnDemand[magic_functions.GreaterThan] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				floatClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Float)
				switch otherClass {
				case integerClass:
					if left {
						if self.GetFloat() > float64(argument[0].GetInt()) {
							return ctx.TrueValue(), nil
						}
						return ctx.FalseValue(), nil
					}
					if self.GetFloat() < float64(argument[0].GetInt()) {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				case floatClass:
					if left {
						if self.GetFloat() > argument[0].GetFloat() {
							return ctx.TrueValue(), nil
						}
						return ctx.FalseValue(), nil
					}
					if self.GetFloat() < argument[0].GetFloat() {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				default:
					return nil, NotComparable
				}
			},
		))
	}
	value.OnDemand[magic_functions.GreaterOrEqualThan] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				floatClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Float)
				switch otherClass {
				case integerClass:
					if left {
						if self.GetFloat() >= float64(argument[0].GetInt()) {
							return ctx.TrueValue(), nil
						}
						return ctx.FalseValue(), nil
					}
					if self.GetFloat() <= float64(argument[0].GetInt()) {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				case floatClass:
					if left {
						if self.GetFloat() >= argument[0].GetFloat() {
							return ctx.TrueValue(), nil
						}
						return ctx.FalseValue(), nil
					}
					if self.GetFloat() <= argument[0].GetFloat() {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				default:
					return nil, NotComparable
				}
			},
		))
	}
	value.OnDemand[magic_functions.LessThan] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				floatClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Float)
				switch otherClass {
				case integerClass:
					if left {
						if self.GetFloat() < float64(argument[0].GetInt()) {
							return ctx.TrueValue(), nil
						}
						return ctx.FalseValue(), nil
					}
					if self.GetFloat() > float64(argument[0].GetInt()) {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				case floatClass:
					if left {
						if self.GetFloat() < argument[0].GetFloat() {
							return ctx.TrueValue(), nil
						}
						return ctx.FalseValue(), nil
					}
					if self.GetFloat() > argument[0].GetFloat() {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				default:
					return nil, NotComparable
				}
			},
		))
	}
	value.OnDemand[magic_functions.LessOrEqualThan] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				floatClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Float)
				switch otherClass {
				case integerClass:
					if left {
						if self.GetFloat() <= float64(argument[0].GetInt()) {
							return ctx.TrueValue(), nil
						}
						return ctx.FalseValue(), nil
					}
					if self.GetFloat() >= float64(argument[0].GetInt()) {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				case floatClass:
					if left {
						if self.GetFloat() <= argument[0].GetFloat() {
							return ctx.TrueValue(), nil
						}
						return ctx.FalseValue(), nil
					}
					if self.GetFloat() >= argument[0].GetFloat() {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				default:
					return nil, NotComparable
				}
			},
		))
	}
	value.OnDemand[magic_functions.Add] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				floatClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Float)
				switch otherClass {
				case integerClass:
					return ctx.FloatValue(self.GetFloat() + float64(argument[0].GetInt())), nil
				case floatClass:
					return ctx.FloatValue(self.GetFloat() + argument[0].GetFloat()), nil
				default:
					return nil, NotOperable
				}
			},
		))
	}
	value.OnDemand[magic_functions.Sub] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				floatClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Float)
				switch otherClass {
				case integerClass:
					return ctx.FloatValue(self.GetFloat() - float64(argument[0].GetInt())), nil
				case floatClass:
					return ctx.FloatValue(self.GetFloat() - argument[0].GetFloat()), nil
				default:
					return nil, NotOperable
				}
			},
		))
	}
	value.OnDemand[magic_functions.Mul] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				floatClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Float)
				switch otherClass {
				case integerClass:
					return ctx.FloatValue(self.GetFloat() * float64(argument[0].GetInt())), nil
				case floatClass:
					return ctx.FloatValue(self.GetFloat() * argument[0].GetFloat()), nil
				default:
					return nil, NotOperable
				}
			},
		))
	}
	value.OnDemand[magic_functions.Div] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				floatClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Float)
				switch otherClass {
				case integerClass:
					if left {
						return ctx.FloatValue(self.GetFloat() / float64(argument[0].GetInt())), nil
					}
					return ctx.FloatValue(float64(argument[0].GetInt()) / self.GetFloat()), nil
				case floatClass:
					if left {
						return ctx.FloatValue(self.GetFloat() / argument[0].GetFloat()), nil
					}
					return ctx.FloatValue(argument[0].GetFloat() / self.GetFloat()), nil
				default:
					return nil, NotOperable
				}
			},
		))
	}
	value.OnDemand[magic_functions.FloorDiv] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				floatClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Float)
				switch otherClass {
				case integerClass:
					if left {
						return ctx.IntegerValue(int64(self.GetFloat() / float64(argument[0].GetInt()))), nil
					}
					return ctx.IntegerValue(int64(float64(argument[0].GetInt()) / self.GetFloat())), nil
				case floatClass:
					if left {
						return ctx.IntegerValue(int64(self.GetFloat() / argument[0].GetFloat())), nil
					}
					return ctx.IntegerValue(int64(argument[0].GetFloat() / self.GetFloat())), nil
				default:
					return nil, NotOperable
				}
			},
		))
	}
	value.OnDemand[magic_functions.Modulus] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				floatClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Float)
				switch otherClass {
				case integerClass:
					if left {
						return ctx.FloatValue(math.Mod(self.GetFloat(), float64(argument[0].GetInt()))), nil
					}
					return ctx.FloatValue(math.Mod(float64(argument[0].GetInt()), self.GetFloat())), nil
				case floatClass:
					if left {
						return ctx.FloatValue(math.Mod(self.GetFloat(), argument[0].GetFloat())), nil
					}
					return ctx.FloatValue(math.Mod(argument[0].GetFloat(), self.GetFloat())), nil
				default:
					return nil, NotOperable
				}
			},
		))
	}
	value.OnDemand[magic_functions.PowerOf] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				floatClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Float)
				switch otherClass {
				case integerClass:
					if left {
						return ctx.FloatValue(math.Pow(self.GetFloat(), float64(argument[0].GetInt()))), nil
					}
					return ctx.FloatValue(math.Pow(float64(argument[0].GetInt()), self.GetFloat())), nil
				case floatClass:
					if left {
						return ctx.FloatValue(math.Pow(self.GetFloat(), argument[0].GetFloat())), nil
					}
					return ctx.FloatValue(math.Pow(argument[0].GetFloat(), self.GetFloat())), nil
				default:
					return nil, NotOperable
				}
			},
		))
	}
	value.OnDemand[magic_functions.Bool] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				if self.GetFloat() != 0 {
					return ctx.TrueValue(), nil
				}
				return ctx.FalseValue(), nil
			},
		))
	}
	value.OnDemand[magic_functions.Class] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				self.mutex.Lock()
				defer self.mutex.Unlock()
				if self.Class == nil {
					var getError error
					self.Class, getError = ctx.VM.RootNamespace.Get(special_symbols.Float)
					if getError != nil {
						panic("Float class not implemented")
					}
				}
				return self.Class, nil
			},
		))
	}
	value.OnDemand[magic_functions.Copy] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				return ctx.FloatValue(self.GetFloat()), nil
			},
		))
	}
	value.OnDemand[magic_functions.String] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				return ctx.StringValue([]byte(fmt.Sprintf("%f", self.GetFloat()))), nil
			},
		))
	}
	return value
}
