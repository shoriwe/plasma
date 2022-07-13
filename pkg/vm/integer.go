package vm

import (
	"bytes"
	"fmt"
	magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"
	special_symbols "github.com/shoriwe/gplasma/pkg/common/special-symbols"
	"math"
)

var (
	NotBitwiseOperable = fmt.Errorf("not bitwise operable")
	NotOperable        = fmt.Errorf("not operable")
)

/*
IntegerValue
	Class: Integer TODO
	Methods:
	- Positive
	- Negative
	- NegateBits
	- Equals: Integer == Any
	- NotEqual: Integer != Any
	- GreaterThan: Integer > Integer, Integer > Float
	- GreaterOrEqualThan: Integer >= Integer, Integer >= Float
	- LessThan: Integer < Integer, Integer < Float
	- LessOrEqualThan: Integer <= Integer, Integer <= Float
	- BitwiseOr: Integer | Integer
	- BitwiseXor: Integer ^ Integer
	- BitwiseAnd: Integer & Integer
	- BitwiseLeft: Integer << Integer
	- BitwiseRight: Integer >> Integer
	- Add: Integer + Integer, Integer + Float
	- Sub: Integer - Integer, Integer - Float
	- Mul: Integer * Integer, Integer * Float, Integer * String, Integer * Bytes, Integer * Array
	- Div: Integer / Integer, Integer / Float
	- FloorDiv: Integer // Integer, Integer // Float
	- Modulus: Integer % Integer
	- PowerOf: Integer ** Integer, Integer ** Float
	- Bool
	- Class
	- Copy
	- String
*/
func (ctx *Context) IntegerValue(i int64) *Value {
	value := ctx.NewValue()
	value.Int = i
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
				return ctx.IntegerValue(-self.GetInt()), nil
			},
		))
	}
	value.OnDemand[magic_functions.NegateBits] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				return ctx.IntegerValue(^self.GetInt()), nil
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
					if self.GetInt() == argument[0].GetInt() {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				case floatClass:
					if float64(self.GetInt()) == argument[0].GetFloat() {
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
					if self.GetInt() != argument[0].GetInt() {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				case floatClass:
					if float64(self.GetInt()) != argument[0].GetFloat() {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				default:
					return ctx.TrueValue(), nil
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
						if self.GetInt() > argument[0].GetInt() {
							return ctx.TrueValue(), nil
						}
						return ctx.FalseValue(), nil
					}
					if self.GetInt() < argument[0].GetInt() {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				case floatClass:
					if left {
						if float64(self.GetInt()) > argument[0].GetFloat() {
							return ctx.TrueValue(), nil
						}
						return ctx.FalseValue(), nil
					}
					if float64(self.GetInt()) < argument[0].GetFloat() {
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
						if self.GetInt() >= argument[0].GetInt() {
							return ctx.TrueValue(), nil
						}
						return ctx.FalseValue(), nil
					}
					if self.GetInt() <= argument[0].GetInt() {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				case floatClass:
					if left {
						if float64(self.GetInt()) >= argument[0].GetFloat() {
							return ctx.TrueValue(), nil
						}
						return ctx.FalseValue(), nil
					}
					if float64(self.GetInt()) <= argument[0].GetFloat() {
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
						if self.GetInt() < argument[0].GetInt() {
							return ctx.TrueValue(), nil
						}
						return ctx.FalseValue(), nil
					}
					if self.GetInt() > argument[0].GetInt() {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				case floatClass:
					if left {
						if float64(self.GetInt()) < argument[0].GetFloat() {
							return ctx.TrueValue(), nil
						}
						return ctx.FalseValue(), nil
					}
					if float64(self.GetInt()) > argument[0].GetFloat() {
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
						if self.GetInt() <= argument[0].GetInt() {
							return ctx.TrueValue(), nil
						}
						return ctx.FalseValue(), nil
					}
					if self.GetInt() >= argument[0].GetInt() {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				case floatClass:
					if left {
						if float64(self.GetInt()) <= argument[0].GetFloat() {
							return ctx.TrueValue(), nil
						}
						return ctx.FalseValue(), nil
					}
					if float64(self.GetInt()) >= argument[0].GetFloat() {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				default:
					return nil, NotComparable
				}
			},
		))
	}
	value.OnDemand[magic_functions.BitwiseOr] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				switch otherClass {
				case integerClass:
					return ctx.IntegerValue(self.GetInt() | argument[0].GetInt()), nil
				default:
					return nil, NotBitwiseOperable
				}
			},
		))
	}
	value.OnDemand[magic_functions.BitwiseXor] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				switch otherClass {
				case integerClass:
					if left {
						return ctx.IntegerValue(self.GetInt() ^ argument[0].GetInt()), nil
					}
					return ctx.IntegerValue(argument[0].GetInt() ^ self.GetInt()), nil
				default:
					return nil, NotBitwiseOperable
				}
			},
		))
	}
	value.OnDemand[magic_functions.BitwiseAnd] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				switch otherClass {
				case integerClass:
					return ctx.IntegerValue(self.GetInt() & argument[0].GetInt()), nil
				default:
					return nil, NotBitwiseOperable
				}
			},
		))
	}
	value.OnDemand[magic_functions.BitwiseLeft] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				switch otherClass {
				case integerClass:
					if left {
						return ctx.IntegerValue(self.GetInt() << argument[0].GetInt()), nil
					}
					return ctx.IntegerValue(argument[0].GetInt() << self.GetInt()), nil
				default:
					return nil, NotBitwiseOperable
				}
			},
		))
	}
	value.OnDemand[magic_functions.BitwiseRight] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				switch otherClass {
				case integerClass:
					if left {
						return ctx.IntegerValue(self.GetInt() >> argument[0].GetInt()), nil
					}
					return ctx.IntegerValue(argument[0].GetInt() >> self.GetInt()), nil
				default:
					return nil, NotBitwiseOperable
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
					return ctx.IntegerValue(self.GetInt() + argument[0].GetInt()), nil
				case floatClass:
					return ctx.FloatValue(float64(self.GetInt()) + argument[0].GetFloat()), nil
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
					return ctx.IntegerValue(self.GetInt() - argument[0].GetInt()), nil
				case floatClass:
					return ctx.FloatValue(float64(self.GetInt()) - argument[0].GetFloat()), nil
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
				stringClass, _ := ctx.VM.RootNamespace.Get(special_symbols.String)
				bytesClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Bytes)
				arrayClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Array)
				switch otherClass {
				case integerClass:
					return ctx.IntegerValue(self.GetInt() * argument[0].GetInt()), nil
				case floatClass:
					return ctx.FloatValue(float64(self.GetInt()) * argument[0].GetFloat()), nil
				case stringClass:
					return ctx.StringValue(bytes.Repeat(argument[0].GetContents(), int(self.GetInt()))), nil
				case bytesClass:
					return ctx.BytesValue(bytes.Repeat(argument[0].GetContents(), int(self.GetInt()))), nil
				case arrayClass:
					return ctx.ArrayValue(RepeatValues(argument[0].GetValues(), self.GetInt())), nil
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
						return ctx.FloatValue(float64(self.GetInt() / argument[0].GetInt())), nil
					}
					return ctx.FloatValue(float64(argument[0].GetInt() / self.GetInt())), nil
				case floatClass:
					if left {
						return ctx.FloatValue(float64(self.GetInt()) / argument[0].GetFloat()), nil
					}
					return ctx.FloatValue(argument[0].GetFloat() / float64(self.GetInt())), nil
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
						return ctx.IntegerValue(self.GetInt() / argument[0].GetInt()), nil
					}
					return ctx.IntegerValue(argument[0].GetInt() / self.GetInt()), nil
				case floatClass:
					if left {
						return ctx.IntegerValue(int64(float64(self.GetInt()) / argument[0].GetFloat())), nil
					}
					return ctx.IntegerValue(int64(argument[0].GetFloat() / float64(self.GetInt()))), nil
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
						return ctx.IntegerValue(self.GetInt() % argument[0].GetInt()), nil
					}
					return ctx.IntegerValue(argument[0].GetInt() % self.GetInt()), nil
				case floatClass:
					if left {
						return ctx.FloatValue(math.Mod(float64(self.GetInt()), argument[0].GetFloat())), nil
					}
					return ctx.FloatValue(math.Mod(argument[0].GetFloat(), float64(self.GetInt()))), nil
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
						return ctx.IntegerValue(int64(math.Pow(float64(self.GetInt()), float64(argument[0].GetInt())))), nil
					}
					return ctx.IntegerValue(int64(math.Pow(float64(argument[0].GetInt()), float64(self.GetInt())))), nil
				case floatClass:
					if left {
						return ctx.FloatValue(math.Pow(float64(self.GetInt()), argument[0].GetFloat())), nil
					}
					return ctx.FloatValue(math.Pow(argument[0].GetFloat(), float64(self.GetInt()))), nil
				default:
					return nil, NotOperable
				}
			},
		))
	}
	value.OnDemand[magic_functions.Bool] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				if self.GetInt() != 0 {
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
					self.Class, getError = ctx.VM.RootNamespace.Get(special_symbols.Integer)
					if getError != nil {
						panic("Integer class not implemented")
					}
				}
				return self.Class, nil
			},
		))
	}
	value.OnDemand[magic_functions.Copy] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				return ctx.IntegerValue(self.GetInt()), nil
			},
		))
	}
	value.OnDemand[magic_functions.String] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				return ctx.StringValue([]byte(fmt.Sprintf("%d", self.GetInt()))), nil
			},
		))
	}
	return value
}
