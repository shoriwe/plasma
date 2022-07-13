package vm

import (
	magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"
	special_symbols "github.com/shoriwe/gplasma/pkg/common/special-symbols"
)

func RepeatValues(values []*Value, times int64) []*Value {
	result := make([]*Value, 0, int64(len(values))*times)
	for i := int64(0); i < times; i++ {
		for _, value := range values {
			result = append(result, value.Copy())
		}
	}
	return result
}

/*
ArrayValue
	Class: Array TODO
	Methods:
	- Equals: Array == Any
	- NotEqual: Array != Any
	- Mul: Array * Integer
	- Length
	- Bool
	- Get
	- Set
	- Del
	- Class
	- Copy
	- String
	- Iter TODO
*/
func (ctx *Context) ArrayValue(values []*Value) *Value {
	value := ctx.NewValue()
	value.Values = values
	value.OnDemand[magic_functions.Equals] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				arrayClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Array)
				switch otherClass {
				case arrayClass:
					selfValues := self.GetValues()
					otherValues := argument[0].GetValues()
					if len(selfValues) != len(otherValues) {
						return ctx.FalseValue(), nil
					}
					for index, v := range selfValues {
						if !ctx.Equals(v, otherValues[index]) {
							return ctx.FalseValue(), nil
						}
					}
					return ctx.TrueValue(), nil
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
				arrayClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Array)
				switch otherClass {
				case arrayClass:
					selfValues := self.GetValues()
					otherValues := argument[0].GetValues()
					if len(selfValues) != len(otherValues) {
						return ctx.TrueValue(), nil
					}
					for index, v := range selfValues {
						if ctx.Equals(v, otherValues[index]) {
							return ctx.FalseValue(), nil
						}
					}
					return ctx.TrueValue(), nil
				default:
					return ctx.TrueValue(), nil
				}
			},
		))
	}
	value.OnDemand[magic_functions.Mul] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				switch otherClass {
				case integerClass:
					return ctx.ArrayValue(RepeatValues(self.GetValues(), argument[0].GetInt())), nil
				default:
					return nil, NotOperable
				}
			},
		))
	}
	value.OnDemand[magic_functions.Length] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				v := self.GetValues()
				return ctx.IntegerValue(int64(len(v))), nil
			},
		))
	}
	value.OnDemand[magic_functions.Bool] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				v := self.GetValues()
				if len(v) > 0 {
					return ctx.TrueValue(), nil
				}
				return ctx.FalseValue(), nil
			},
		))
	}
	value.OnDemand[magic_functions.Get] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				c := self.GetValues()
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				tupleClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Tuple)
				switch otherClass {
				case integerClass:
					return c[argument[0].GetInt()], nil
				case tupleClass:
					indexes := argument[0].GetValues()
					return ctx.ArrayValue(c[indexes[0].GetInt():indexes[1].GetInt()]), nil
				default:
					return nil, NotIndexable
				}
			},
		))
	}
	value.OnDemand[magic_functions.Set] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				switch otherClass {
				case integerClass:
					self.mutex.Lock()
					defer self.mutex.Unlock()
					self.Values[argument[0].GetInt()] = argument[1]
					return ctx.NoneValue(), nil
				default:
					return nil, NotIndexable
				}
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
					self.Class, getError = ctx.VM.RootNamespace.Get(special_symbols.Array)
					if getError != nil {
						panic("Array class not implemented")
					}
				}
				return self.Class, nil
			},
		))
	}
	value.OnDemand[magic_functions.Copy] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				c := self.GetValues()
				copyValues := make([]*Value, 0, len(c))
				for _, v := range c {
					copyValues = append(copyValues, v.Copy())
				}
				return ctx.ArrayValue(copyValues), nil
			},
		))
	}
	value.OnDemand[magic_functions.String] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				var contents []byte
				contents = append(contents, '[')
				for index, v := range self.GetValues() {
					if index != 0 {
						contents = append(contents, ',')
					}
					contents = append(contents, v.String()...)
				}
				contents = append(contents, ']')
				return ctx.StringValue(contents), nil
			},
		))
	}
	value.OnDemand[magic_functions.Iter] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				// TODO: Implement me!
				panic("implement me!")
			},
		))
	}
	return value
}
