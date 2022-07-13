package vm

import (
	"fmt"
	magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"
	special_symbols "github.com/shoriwe/gplasma/pkg/common/special-symbols"
	"sync"
)

var (
	NotCallable    = fmt.Errorf("value not callable")
	NotImplemented = fmt.Errorf("not implemented")
	NotComparable  = fmt.Errorf("not comparable")
)

var (
	NotImplementedCallable = NewBuiltInCallable(
		func(left bool, argument ...*Value) (*Value, error) {
			return nil, NotImplemented
		},
	)
)

type (
	OnDemand func(self *Value) (*Value, error)
	Callable interface {
		LoadArguments(left bool, argument ...*Value)
		Call() (*Value, error)
	}
	Value struct {
		IsFunction   bool
		mutex        *sync.Mutex
		Class        *Value
		Contents     []byte
		Int          int64
		Float        float64
		Values       []*Value
		VirtualTable *Symbols
		OnDemand     map[string]OnDemand
		Callable     Callable
	}
)

func (value *Value) GetClass() *Value {
	// TODO: implement me!
	panic("implement me!")
}

func (value *Value) GetIsFunction() bool {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	return value.IsFunction
}

func (value *Value) SetIsFunction(isFunction bool) {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	value.IsFunction = isFunction
}

func (value *Value) GetContents() []byte {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	return value.Contents
}

func (value *Value) SetContents(contents []byte) {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	value.Contents = contents
}

func (value *Value) GetInt() int64 {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	return value.Int
}

func (value *Value) SetInt(i int64) {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	value.Int = i
}

func (value *Value) GetFloat() float64 {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	return value.Float
}

func (value *Value) SetFloat(f float64) {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	value.Float = f
}

func (value *Value) GetValues() []*Value {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	return value.Values
}

func (value *Value) SetValues(values []*Value) {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	value.Values = values
}

func (ctx *Context) NewValue() *Value {
	onDemand := map[string]OnDemand{
		magic_functions.Init: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.HasNext: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.Next: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.Not: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NewBuiltInCallable(
				func(left bool, argument ...*Value) (*Value, error) {
					if self.Bool() {
						return ctx.FalseValue(), nil
					}
					return ctx.TrueValue(), nil
				},
			))
		},
		magic_functions.Positive: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.Negative: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.NegateBits: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.And: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NewBuiltInCallable(
				func(left bool, argument ...*Value) (*Value, error) {
					if self.Bool() && argument[0].Bool() {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				},
			))
		},
		magic_functions.Or: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NewBuiltInCallable(
				func(left bool, argument ...*Value) (*Value, error) {
					if self.Bool() || argument[0].Bool() {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				},
			))
		},
		magic_functions.Xor: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NewBuiltInCallable(
				func(left bool, argument ...*Value) (*Value, error) {
					if self.Bool() != argument[0].Bool() {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				},
			))
		},
		magic_functions.In: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.Is: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NewBuiltInCallable(
				func(left bool, argument ...*Value) (*Value, error) {
					self.mutex.Lock()
					defer self.mutex.Unlock()
					classMethod, getError := self.Get(magic_functions.Class)
					if getError != nil {
						panic(getError)
					}
					class, callError := classMethod.Call(false)
					if callError != nil {
						return nil, callError
					}
					if class == argument[0] {
						return ctx.TrueValue(), nil
					}
					return ctx.FalseValue(), nil
				},
			))
		},
		magic_functions.Implements: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NewBuiltInCallable(
				func(left bool, argument ...*Value) (*Value, error) {
					// TODO: Implement me!
					panic("implement me!")
				},
			))
		},
		magic_functions.Equals: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.NotEqual: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.GreaterThan: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.GreaterOrEqualThan: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.LessThan: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.LessOrEqualThan: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.BitwiseOr: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.BitwiseXor: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.BitwiseAnd: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.BitwiseLeft: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.BitwiseRight: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.Add: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.Sub: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.Mul: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.Div: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.FloorDiv: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.Modulus: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.PowerOf: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.Length: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.Bool: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NewBuiltInCallable(
				func(left bool, argument ...*Value) (*Value, error) {
					return ctx.TrueValue(), nil
				},
			))
		},
		magic_functions.Get: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.Set: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.Del: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.Call: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.Class: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NewBuiltInCallable(
				func(left bool, argument ...*Value) (*Value, error) {
					self.mutex.Lock()
					defer self.mutex.Unlock()
					if self.Class == nil {
						var getError error
						self.Class, getError = ctx.VM.RootNamespace.Get(special_symbols.Value)
						if getError != nil {
							panic("Value class not implemented")
						}
					}
					return self.Class, nil
				},
			))
		},
		magic_functions.Copy: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.String: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
		magic_functions.Iter: func(self *Value) (*Value, error) {
			return ctx.NewFunctionValue(NotImplementedCallable)
		},
	}
	return &Value{
		IsFunction:   false,
		mutex:        &sync.Mutex{},
		Class:        nil,
		Contents:     nil,
		Int:          0,
		Float:        0,
		Values:       nil,
		VirtualTable: NewSymbols(ctx.Namespace),
		OnDemand:     onDemand,
		Callable:     nil,
	}
}

func (value *Value) Get(symbol string) (*Value, error) {
	result, getError := value.VirtualTable.Get(symbol)
	if getError == nil {
		return result, nil
	}
	value.mutex.Lock()
	defer value.mutex.Unlock()
	onDemand, found := value.OnDemand[symbol]
	if !found {
		return nil, SymbolNotFoundError
	}
	var onDemandError error
	result, onDemandError = onDemand(value)
	if onDemandError != nil {
		return nil, onDemandError
	}
	value.VirtualTable.Set(symbol, result)
	return result, nil
}

func (value *Value) Call(left bool, argument ...*Value) (*Value, error) {
	if !value.IsFunction {
		call, getError := value.Get(magic_functions.Call)
		if getError != nil {
			return nil, NotCallable
		}
		return call.Call(left, argument...)
	}
	value.mutex.Lock()
	defer value.mutex.Unlock()
	value.Callable.LoadArguments(left, argument...)
	return value.Callable.Call()
}

func (value *Value) Bool() bool {
	boolMethod, getError := value.Get(magic_functions.Bool)
	if getError != nil {
		panic("Value doesn't implement __bool__")
	}
	result, callError := boolMethod.Call(false)
	if callError != nil {
		return false
	}
	return result.Int == 1
}

func (value *Value) String() string {
	stringMethod, getError := value.Get(magic_functions.String)
	if getError != nil {
		panic("Value doesn't implement __string__")
	}
	result, callError := stringMethod.Call(false)
	if callError != nil {
		return fmt.Sprintf("%#v", value)
	}
	return string(result.Contents)
}

func (value *Value) Copy() *Value {
	copyMethod, getError := value.Get(magic_functions.Copy)
	if getError != nil {
		panic("Value doesn't implement __copy__")
	}
	copied, callError := copyMethod.Call(false)
	if callError != nil {
		return value
	}
	return copied
}
