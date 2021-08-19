package vm

func (p *Plasma) NewModule(context *Context, isBuiltIn bool, parent *SymbolTable) *Value {
	module := p.NewValue(context, isBuiltIn, ModuleName, nil, parent)
	module.BuiltInTypeId = ModuleId
	module.SetOnDemandSymbol(Self,
		func() *Value {
			return module
		},
	)
	return module
}
