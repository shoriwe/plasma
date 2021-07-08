package vm

func (p *Plasma) NewModule(context *Context, isBuiltIn bool, parent *SymbolTable) Value {
	module := p.NewObject(context, isBuiltIn, ModuleName, nil, parent)
	module.SetOnDemandSymbol(Self,
		func() Value {
			return module
		},
	)
	return module
}
