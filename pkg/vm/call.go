package vm

func (p *Plasma) CallableInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(object Value) *Object {
		object.SetOnDemandSymbol(Call,
			func() Value {
				return p.NewFunction(isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(_ Value, _ ...Value) (Value, *Object) {
							return nil, p.NewNotImplementedCallableError(Call)
						},
					),
				)
			},
		)
		return nil
	}
}

func (p *Plasma) CallFunction(function Value, parent *SymbolTable, arguments ...Value) (Value, *Object) {
	var callFunction *Function
	if _, ok := function.(*Function); !ok {
		call, getError := function.Get(Call)
		if getError != nil {
			return nil, p.NewObjectNotCallable(function.GetClass(p))
		}
		if _, ok = call.(*Function); !ok {
			return nil, p.NewInvalidTypeError(function.TypeName(), CallableName)
		}
		callFunction = call.(*Function)
	} else {
		callFunction = function.(*Function)
	}
	if callFunction.Callable.NumberOfArguments() != len(arguments) {
		//  Return Here a error related to number of arguments
		return nil, p.NewInvalidNumberOfArgumentsError(len(arguments), callFunction.Callable.NumberOfArguments())
	}
	symbols := NewSymbolTable(parent)
	self, callback, code := callFunction.Callable.Call()
	if self != nil {
		symbols.Set(Self, self)
	} else {
		symbols.Set(Self, function)
	}
	p.PushSymbolTable(symbols)
	var result Value
	var callError *Object
	if callback != nil {
		result, callError = callback(self, arguments...)
	} else if code != nil {
		// Load the arguments
		for i := len(arguments) - 1; i > -1; i-- {
			p.PushObject(arguments[i])
		}
		result, callError = p.Execute(NewBytecodeFromArray(code))
	} else {
		panic("callback and code are nil")
	}
	p.PopSymbolTable()
	if callError != nil {
		return nil, callError
	}
	return result, nil
}

type FunctionCallback func(Value, ...Value) (Value, *Object)

type Callable interface {
	NumberOfArguments() int
	Call() (Value, FunctionCallback, []Code) // self should return directly the object or the code of the function
}

type PlasmaFunction struct {
	numberOfArguments int
	Code              []Code
}

func (p *PlasmaFunction) NumberOfArguments() int {
	return p.numberOfArguments
}

func (p *PlasmaFunction) Call() (Value, FunctionCallback, []Code) {
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
	Self              Value
}

func (p *PlasmaClassFunction) NumberOfArguments() int {
	return p.numberOfArguments
}

func (p *PlasmaClassFunction) Call() (Value, FunctionCallback, []Code) {
	return p.Self, nil, p.Code
}

func NewPlasmaClassFunction(self Value, numberOfArguments int, code []Code) *PlasmaClassFunction {
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

func (g *BuiltInFunction) Call() (Value, FunctionCallback, []Code) {
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
	Self              Value
}

func (g *BuiltInClassFunction) NumberOfArguments() int {
	return g.numberOfArguments
}

func (g *BuiltInClassFunction) Call() (Value, FunctionCallback, []Code) {
	return g.Self, g.callback, nil
}

func NewBuiltInClassFunction(self Value, numberOfArguments int, callback FunctionCallback) *BuiltInClassFunction {
	return &BuiltInClassFunction{
		numberOfArguments: numberOfArguments,
		callback:          callback,
		Self:              self,
	}
}
