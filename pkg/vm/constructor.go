package vm

import "github.com/shoriwe/gruby/pkg/errors"

type Constructor interface {
	Construct(*Plasma, IObject) *errors.Error
}

type PlasmaConstructor struct {
	Constructor
	Code []Code
}

func (c *PlasmaConstructor) Construct(vm *Plasma, object IObject) *Object {
	vm.PushCode(NewBytecodeFromArray(c.Code))
	vm.PushSymbolTable(object.SymbolTable())
	_, executionError := vm.Execute()
	return executionError
}

func NewPlasmaConstructor(code []Code) *PlasmaConstructor {
	return &PlasmaConstructor{
		Code: code,
	}
}

type ConstructorCallBack func(IObject) *errors.Error

type BuiltInConstructor struct {
	Constructor
	callback ConstructorCallBack
}

func (c *BuiltInConstructor) Construct(_ *Plasma, object IObject) *errors.Error {
	return c.callback(object)
}

func NewBuiltInConstructor(callback ConstructorCallBack) *BuiltInConstructor {
	return &BuiltInConstructor{
		callback: callback,
	}
}
