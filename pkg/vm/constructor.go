package vm

type Constructor interface {
	Construct(*Plasma, IObject) *Object
}

type PlasmaConstructor struct {
	Constructor
	Code []Code
}

func (c *PlasmaConstructor) Construct(vm *Plasma, object IObject) *Object {
	vm.PushCode(NewBytecodeFromArray(c.Code))
	vm.PushSymbolTable(object.SymbolTable())
	vm.PushObject(object)
	_, executionError := vm.Execute()
	vm.PopCode()
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
