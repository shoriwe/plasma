package vm

func (p *Plasma) constructSubClass(subClass *Type, object IObject) *Object {
	for _, subSubClass := range subClass.subClasses {
		object.SymbolTable().Parent = subSubClass.symbols.Parent
		subSubClassConstructionError := p.constructSubClass(subSubClass, object)
		if subSubClassConstructionError != nil {
			return subSubClassConstructionError
		}
	}
	object.SymbolTable().Parent = subClass.symbols.Parent
	baseInitializationError := subClass.Constructor.Construct(p, object)
	if baseInitializationError != nil {
		return baseInitializationError
	}
	return nil
}

func (p *Plasma) ConstructObject(type_ *Type, parent *SymbolTable) (IObject, *Object) {
	object := p.NewObject(type_.Name, type_.subClasses, parent)
	for _, subclass := range object.subClasses {
		subClassConstructionError := p.constructSubClass(subclass, object)
		if subClassConstructionError != nil {
			return nil, subClassConstructionError
		}
	}
	object.SymbolTable().Parent = parent
	object.class = type_
	baseInitializationError := type_.Constructor.Construct(p, object)
	if baseInitializationError != nil {
		return nil, baseInitializationError
	}
	return object, nil
}

type Constructor interface {
	Construct(*Plasma, IObject) *Object
}

type PlasmaConstructor struct {
	Constructor
	Code []Code
}

func (c *PlasmaConstructor) Construct(vm *Plasma, object IObject) *Object {
	vm.PushBytecode(NewBytecodeFromArray(c.Code))
	vm.PushSymbolTable(object.SymbolTable())
	vm.PushObject(object)
	_, executionError := vm.Execute()
	vm.PopBytecode()
	vm.PopSymbolTable()
	return executionError
}

func NewPlasmaConstructor(code []Code) *PlasmaConstructor {
	return &PlasmaConstructor{
		Code: code,
	}
}

type ConstructorCallBack func(IObject) *Object

type BuiltInConstructor struct {
	Constructor
	callback ConstructorCallBack
}

func (c *BuiltInConstructor) Construct(_ *Plasma, object IObject) *Object {
	return c.callback(object)
}

func NewBuiltInConstructor(callback ConstructorCallBack) *BuiltInConstructor {
	return &BuiltInConstructor{
		callback: callback,
	}
}
