package vm

import (
	magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

/*
NoneValue
	Class: NoneType TODO
	Methods:
	- Bool
	- Copy
	- String
*/
func (ctx *Context) NoneValue() *Value {
	if ctx.VM.NoneValue != nil {
		return ctx.VM.NoneValue
	}
	ctx.VM.mutex.Lock()
	defer ctx.VM.mutex.Unlock()
	ctx.VM.NoneValue = ctx.NewValue()
	ctx.VM.NoneValue.OnDemand[magic_functions.Bool] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				return ctx.FalseValue(), nil
			},
		))
	}
	ctx.VM.NoneValue.OnDemand[magic_functions.Copy] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				return self, nil
			},
		))
	}
	ctx.VM.NoneValue.OnDemand[magic_functions.String] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				return ctx.StringValue([]byte(lexer.NoneString)), nil
			},
		))
	}
	return ctx.VM.NoneValue
}
