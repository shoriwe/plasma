package vm

import (
	"bytes"
	"fmt"
	magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"
	special_symbols "github.com/shoriwe/gplasma/pkg/common/special-symbols"
)

var (
	NotIndexable = fmt.Errorf("not indexable")
)

/*
StringValue
Class: String TODO
Methods:
- Equals: String == Any
- NotEqual: String != Any
- Add: String + String
- Mul: String * Integer
- Length
- Bool
- Get: Integer, Tuple
- Class
- Copy
- String
- Iter TODO
*/
func (ctx *Context) StringValue(contents []byte) *Value {
	value := ctx.NewValue()
	value.Contents = contents
	value.OnDemand[magic_functions.Equals] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				stringClass, _ := ctx.VM.RootNamespace.Get(special_symbols.String)
				switch otherClass {
				case stringClass:
					if bytes.Equal(self.GetContents(), argument[0].GetContents()) {
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
				stringClass, _ := ctx.VM.RootNamespace.Get(special_symbols.String)
				switch otherClass {
				case stringClass:
					if !bytes.Equal(self.GetContents(), argument[0].GetContents()) {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				default:
					return ctx.FalseValue(), nil
				}
			},
		))
	}
	value.OnDemand[magic_functions.Add] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				otherClass := argument[0].GetClass()
				stringClass, _ := ctx.VM.RootNamespace.Get(special_symbols.String)
				switch otherClass {
				case stringClass:
					var newContents []byte
					if left {
						newContents = append(newContents, self.GetContents()...)
						newContents = append(newContents, argument[0].GetContents()...)
					} else {
						newContents = append(newContents, argument[0].GetContents()...)
						newContents = append(newContents, self.GetContents()...)
					}
					return ctx.StringValue(newContents), nil
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
				switch otherClass {
				case integerClass:
					return ctx.StringValue(bytes.Repeat(self.GetContents(), int(argument[0].GetInt()))), nil
				default:
					return nil, NotOperable
				}
			},
		))
	}
	value.OnDemand[magic_functions.Length] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				return ctx.IntegerValue(int64(len(self.GetContents()))), nil
			},
		))
	}
	value.OnDemand[magic_functions.Bool] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				if len(self.GetContents()) > 0 {
					return ctx.TrueValue(), nil
				}
				return ctx.FalseValue(), nil
			},
		))
	}
	value.OnDemand[magic_functions.Get] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				c := self.GetContents()
				otherClass := argument[0].GetClass()
				integerClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Integer)
				tupleClass, _ := ctx.VM.RootNamespace.Get(special_symbols.Tuple)
				switch otherClass {
				case integerClass:
					return ctx.StringValue([]byte{c[argument[0].GetInt()]}), nil
				case tupleClass:
					values := argument[0].GetValues()
					return ctx.StringValue(c[values[0].GetInt():values[1].GetInt()]), nil
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
					self.Class, getError = ctx.VM.RootNamespace.Get(special_symbols.String)
					if getError != nil {
						panic("String class not implemented")
					}
				}
				return self.Class, nil
			},
		))
	}
	value.OnDemand[magic_functions.Copy] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				c := self.GetContents()
				newChunk := make([]byte, len(c))
				copy(newChunk, c)
				return ctx.StringValue(newChunk), nil
			},
		))
	}
	value.OnDemand[magic_functions.String] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				return self.Copy(), nil
			},
		))
	}
	value.OnDemand[magic_functions.Iter] = func(self *Value) (*Value, error) {
		return ctx.NewFunctionValue(NewBuiltInCallable(
			func(left bool, argument ...*Value) (*Value, error) {
				// TODO: implement me!
				panic("implement me!")
			},
		))
	}
	return value
}
