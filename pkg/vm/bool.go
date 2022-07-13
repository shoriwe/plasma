package vm

import (
	magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

/*
NewBool
	Class: Bool TODO
	Methods:
	- Equals: Bool == Any
	- NotEqual: Bool != Any
	- Bool
	- Copy
	- String
*/
func (ctx *Context) NewBool() *Value {
	value := ctx.NewValue()
	value.OnDemand[magic_functions.Equals] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				if self.GetClass() == argument[0].GetClass() {
					return ctx.TrueValue(), nil
				}
				if self.Bool() == argument[0].Bool() {
					return ctx.TrueValue(), nil
				}
				return ctx.FalseValue(), nil
			},
		))
	}
	value.OnDemand[magic_functions.NotEqual] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				if self.GetClass() != argument[0].GetClass() {
					return ctx.TrueValue(), nil
				}
				if self.Bool() != argument[0].Bool() {
					return ctx.TrueValue(), nil
				}
				return ctx.FalseValue(), nil
			},
		))
	}
	value.OnDemand[magic_functions.Bool] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				if self.Bool() {
					return ctx.TrueValue(), nil
				}
				return ctx.FalseValue(), nil
			},
		))
	}
	value.OnDemand[magic_functions.Copy] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				if self.Bool() {
					return ctx.TrueValue(), nil
				}
				return ctx.FalseValue(), nil
			},
		))
	}
	value.OnDemand[magic_functions.String] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				if self.Bool() {
					return ctx.StringValue([]byte(lexer.TrueString)), nil
				}
				return ctx.StringValue([]byte(lexer.FalseString)), nil
			},
		))
	}
	return value
}

func (ctx *Context) TrueValue() *Value {
	if ctx.VM.TrueValue != nil {
		return ctx.VM.TrueValue
	}
	ctx.VM.mutex.Lock()
	defer ctx.VM.mutex.Unlock()
	ctx.VM.TrueValue = ctx.NewBool()
	ctx.VM.TrueValue.Int = 1
	return ctx.VM.TrueValue
}

func (ctx *Context) FalseValue() *Value {
	if ctx.VM.FalseValue == nil {
		return ctx.VM.FalseValue
	}
	ctx.VM.mutex.Lock()
	defer ctx.VM.mutex.Unlock()
	ctx.VM.FalseValue = ctx.NewBool()
	ctx.VM.FalseValue.Int = 0
	return ctx.VM.FalseValue
}
