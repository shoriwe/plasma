package vm

func (p *Plasma) constructSubClass(context *Context, subClass *Type, object Value) *Object {
	for _, subSubClass := range subClass.subClasses {
		object.SymbolTable().Parent = subSubClass.symbols.Parent
		subSubClassConstructionError := p.constructSubClass(context, subSubClass, object)
		if subSubClassConstructionError != nil {
			return subSubClassConstructionError
		}
	}
	object.SymbolTable().Parent = subClass.symbols.Parent
	baseInitializationError := subClass.Constructor.Construct(context, p, object)
	if baseInitializationError != nil {
		return baseInitializationError
	}
	return nil
}

func (p *Plasma) ConstructObject(context *Context, type_ *Type, parent *SymbolTable) (Value, *Object) {
	object := p.NewObject(context, false, type_.Name, type_.subClasses, parent)
	for _, subclass := range object.subClasses {
		subClassConstructionError := p.constructSubClass(context, subclass, object)
		if subClassConstructionError != nil {
			return nil, subClassConstructionError
		}
	}
	object.SymbolTable().Parent = parent
	object.class = type_
	baseInitializationError := type_.Constructor.Construct(context, p, object)
	if baseInitializationError != nil {
		return nil, baseInitializationError
	}
	return object, nil
}

type Constructor interface {
	Construct(*Context, *Plasma, Value) *Object
}

type PlasmaConstructor struct {
	Constructor
	Code []Code
}

func (c *PlasmaConstructor) Construct(context *Context, vm *Plasma, object Value) *Object {
	context.PushSymbolTable(object.SymbolTable())
	context.PushObject(object)
	_, executionError := vm.Execute(context, NewBytecodeFromArray(c.Code))
	context.PopSymbolTable()
	return executionError
}

func NewPlasmaConstructor(code []Code) *PlasmaConstructor {
	return &PlasmaConstructor{
		Code: code,
	}
}

type ConstructorCallBack func(*Context, Value) *Object

type BuiltInConstructor struct {
	Constructor
	callback ConstructorCallBack
}

func (c *BuiltInConstructor) Construct(context *Context, _ *Plasma, object Value) *Object {
	return c.callback(context, object)
}

func NewBuiltInConstructor(callback ConstructorCallBack) *BuiltInConstructor {
	return &BuiltInConstructor{
		callback: callback,
	}
}
