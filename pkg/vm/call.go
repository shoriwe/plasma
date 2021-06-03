package vm

import "github.com/shoriwe/gruby/pkg/errors"

type FunctionCallback func(IObject, ...IObject) (IObject, *errors.Error)

func NewNotImplementedCallable(numberOfArguments int) *BuiltInClassFunction {
	return NewBuiltInClassFunction(nil, numberOfArguments,
		func(self IObject, _ ...IObject) (IObject, *errors.Error) {
			return nil, errors.NewNameNotFoundError()
		},
	)
}

type Callable interface {
	NumberOfArguments() int
	Call() (IObject, FunctionCallback, []Code) // self should return directly the object or the code of the function
}

type PlasmaFunction struct {
	numberOfArguments int
	Code              []Code
}

func (p *PlasmaFunction) NumberOfArguments() int {
	return p.numberOfArguments
}

func (p *PlasmaFunction) Call() (IObject, FunctionCallback, []Code) {
	return nil, nil, p.Code
}

func NewPlasmaFunction(numberOfArguments int, code []Code) *PlasmaFunction {
	return &PlasmaFunction{
		numberOfArguments: numberOfArguments,
		Code:              code,
	}
}

type PlasmaClassFunction struct {
	numberOfArguments int
	Code              []Code
	Self              IObject
}

func (p *PlasmaClassFunction) NumberOfArguments() int {
	return p.numberOfArguments
}

func (p *PlasmaClassFunction) Call() (IObject, FunctionCallback, []Code) {
	return p.Self, nil, p.Code
}

func NewPlasmaClassFunction(self IObject, numberOfArguments int, code []Code) *PlasmaClassFunction {
	return &PlasmaClassFunction{
		numberOfArguments: numberOfArguments,
		Code:              code,
		Self:              self,
	}
}

type BuiltInFunction struct {
	numberOfArguments int
	callback          FunctionCallback
}

func (g *BuiltInFunction) NumberOfArguments() int {
	return g.numberOfArguments
}

func (g *BuiltInFunction) Call() (IObject, FunctionCallback, []Code) {
	return nil, g.callback, nil
}

func NewBuiltInFunction(numberOfArguments int, callback FunctionCallback) *BuiltInFunction {
	return &BuiltInFunction{
		numberOfArguments: numberOfArguments,
		callback:          callback,
	}
}

type BuiltInClassFunction struct {
	numberOfArguments int
	callback          FunctionCallback
	Self              IObject
}

func (g *BuiltInClassFunction) NumberOfArguments() int {
	return g.numberOfArguments
}

func (g *BuiltInClassFunction) Call() (IObject, FunctionCallback, []Code) {
	return g.Self, g.callback, nil
}

func NewBuiltInClassFunction(self IObject, numberOfArguments int, callback FunctionCallback) *BuiltInClassFunction {
	return &BuiltInClassFunction{
		numberOfArguments: numberOfArguments,
		callback:          callback,
		Self:              self,
	}
}
