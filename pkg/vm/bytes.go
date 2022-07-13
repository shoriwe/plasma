package vm

import (
	magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"
	special_symbols "github.com/shoriwe/gplasma/pkg/common/special-symbols"
)

/*
BytesValue
Class: Bytes TODO
Methods:
- *All String methods*
- Class
*/
func (ctx *Context) BytesValue(contents []byte) *Value {
	value := ctx.StringValue(contents)
	value.OnDemand[magic_functions.Class] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				self.mutex.Lock()
				defer self.mutex.Unlock()
				if self.Class == nil {
					var getError error
					self.Class, getError = ctx.VM.RootNamespace.Get(special_symbols.Bytes)
					if getError != nil {
						panic("Bytes class not implemented")
					}
				}
				return self.Class, nil
			},
		))
	}
	return value
}
