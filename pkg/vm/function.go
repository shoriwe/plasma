package vm

type Function struct {
	*Object
	Callable Callable
}

func (p *Plasma) NewFunction(context *Context, isBuiltIn bool, parentSymbols *SymbolTable, callable Callable) *Function {
	function := &Function{
		Object:   p.NewObject(context, isBuiltIn, FunctionName, nil, parentSymbols),
		Callable: callable,
	}
	function.SetOnDemandSymbol(Self,
		func() Value {
			return function
		},
	)
	return function
}
