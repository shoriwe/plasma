package vm

type (
	BytecodeCallable struct {
		Bytecode  []byte
		Arguments map[string]*Value
	}
	BuiltInCallable struct {
		left      bool
		arguments []*Value
		Callback  func(left bool, argument ...*Value) (*Value, error)
	}
)

func (b *BuiltInCallable) LoadArguments(left bool, argument ...*Value) {
	b.left = left
	b.arguments = argument
}

func (b *BuiltInCallable) Call() (*Value, error) {
	return b.Callback(b.left, b.arguments...)
}

func NewBuiltInCallable(callback func(left bool, argument ...*Value) (*Value, error)) *BuiltInCallable {
	return &BuiltInCallable{
		arguments: nil,
		Callback:  callback,
	}
}

func (ctx *Context) NewFunctionValue(callable Callable) (*Value, error) {
	value := ctx.NewValue()
	value.IsFunction = true
	value.Callable = callable
	return value, nil
}
