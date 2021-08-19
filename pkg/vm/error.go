package vm

import (
	"fmt"
)

const (
	RuntimeError                  = "RuntimeError"                  // Done
	InvalidTypeError              = "InvalidTypeError"              // Done
	NotImplementedCallableError   = "NotImplementedCallableError"   // Done
	ObjectConstructionError       = "ObjectConstructionError"       // Done
	ObjectWithNameNotFoundError   = "ObjectWithNameNotFoundError"   // Done
	InvalidNumberOfArgumentsError = "InvalidNumberOfArgumentsError" // Done
	GoRuntimeError                = "GoRuntimeError"                // Done
	UnhashableTypeError           = "UnhashableTypeError"           // Done
	IndexOutOfRangeError          = "IndexOutOfRangeError"          // Done
	KeyNotFoundError              = "KeyNotFoundError"              // Done
	IntegerParsingError           = "IntegerParsingError"           // Done
	FloatParsingError             = "FloatParsingError"             // Done
	BuiltInSymbolProtectionError  = "BuiltInSymbolProtectionError"  // Done
	ObjectNotCallableError        = "ObjectNotCallableError"        // Done
)

func (p *Plasma) ForceGetSelf(name string, parent *Value) *Value {
	object, getError := parent.Get(name)
	if getError != nil {
		panic(getError.String())
	}
	return object
}

func (p *Plasma) ForceMasterGetAny(name string) *Value {
	object, getError := p.builtInContext.PeekSymbolTable().GetAny(name)
	if getError != nil {
		panic(getError.String())
	}
	return object
}

func (p *Plasma) ForceConstruction(context *Context, type_ *Value) *Value {
	if !type_.IsTypeById(TypeId) {
		panic("don't modify built-ins, or fatal errors like this one will ocurr")
	}
	result, success := p.ConstructObject(context, type_, NewSymbolTable(context.PeekSymbolTable()))
	if !success {
		panic(result.typeName)
	}
	return result
}

func (p *Plasma) ForceInitialization(context *Context, object *Value, arguments ...*Value) {
	initialize, getError := object.Get(Initialize)
	if getError != nil {
		panic(getError.String())
	}
	callError, success := p.CallFunction(context,
		initialize,
		arguments...,
	)
	if !success {
		panic(fmt.Sprintf("%s: %s", callError.TypeName(), callError.GetString()))
	}
}

func (p *Plasma) NewFloatParsingError(context *Context) *Value {
	errorType := p.ForceMasterGetAny(FloatParsingError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError)
	return instantiatedError
}

func (p *Plasma) NewIntegerParsingError(context *Context) *Value {
	errorType := p.ForceMasterGetAny(IntegerParsingError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError)
	return instantiatedError
}

func (p *Plasma) NewKeyNotFoundError(context *Context, key *Value) *Value {
	errorType := p.ForceMasterGetAny(KeyNotFoundError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError,
		key,
	)
	return instantiatedError
}

func (p *Plasma) NewIndexOutOfRange(context *Context, length int, index int64) *Value {
	errorType := p.ForceMasterGetAny(IndexOutOfRangeError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError,
		p.NewInteger(context,
			false,
			context.PeekSymbolTable(),
			int64(length),
		),
		p.NewInteger(context, false, context.PeekSymbolTable(), index),
	)
	return instantiatedError
}

func (p *Plasma) NewUnhashableTypeError(context *Context, objectType *Value) *Value {
	errorType := p.ForceMasterGetAny(UnhashableTypeError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError,
		objectType,
	)
	return instantiatedError
}

func (p *Plasma) NewNotImplementedCallableError(context *Context, methodName string) *Value {
	errorType := p.ForceMasterGetAny(NotImplementedCallableError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError,
		p.NewString(context, false, context.PeekSymbolTable(), methodName),
	)
	return instantiatedError
}

func (p *Plasma) NewGoRuntimeError(context *Context, er error) *Value {
	errorType := p.ForceMasterGetAny(GoRuntimeError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError,
		p.NewString(context, false, context.PeekSymbolTable(), er.Error()),
	)
	return instantiatedError
}

func (p *Plasma) NewInvalidNumberOfArgumentsError(context *Context, received int, expecting int) *Value {
	errorType := p.ForceMasterGetAny(InvalidNumberOfArgumentsError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError,
		p.NewInteger(context, false, context.PeekSymbolTable(), int64(received)),
		p.NewInteger(context, false, context.PeekSymbolTable(), int64(expecting)),
	)
	return instantiatedError
}

func (p *Plasma) NewObjectWithNameNotFoundError(context *Context, objectType *Value, name string) *Value {
	errorType := p.ForceMasterGetAny(ObjectWithNameNotFoundError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError,
		objectType, p.NewString(context, false, context.PeekSymbolTable(), name),
	)
	return instantiatedError
}

func (p *Plasma) NewInvalidTypeError(context *Context, received string, expecting ...string) *Value {
	errorType := p.ForceMasterGetAny(InvalidTypeError)
	instantiatedError := p.ForceConstruction(context, errorType)
	instantiatedErrorInitialize, _ := instantiatedError.Get(Initialize)
	expectingSum := ""
	for index, s := range expecting {
		if index != 0 {
			expectingSum += ", "
		}
		expectingSum += s
	}
	_, _ = p.CallFunction(context,
		instantiatedErrorInitialize,
		p.NewString(context, false, context.PeekSymbolTable(), received),
		p.NewString(context, false, context.PeekSymbolTable(), expectingSum),
	)
	return instantiatedError
}

func (p *Plasma) NewObjectConstructionError(context *Context, typeName string, errorMessage string) *Value {
	errorType := p.ForceMasterGetAny(ObjectConstructionError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError,
		p.NewString(context, false, context.PeekSymbolTable(), typeName), p.NewString(context, false, context.PeekSymbolTable(), errorMessage),
	)
	return instantiatedError
}

func (p *Plasma) NewBuiltInSymbolProtectionError(context *Context, symbolName string) *Value {
	errorType := p.ForceMasterGetAny(BuiltInSymbolProtectionError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError,
		p.NewString(context, false, context.PeekSymbolTable(), symbolName),
	)
	return instantiatedError
}

func (p *Plasma) NewObjectNotCallable(context *Context, objectType *Value) *Value {
	errorType := p.ForceMasterGetAny(ObjectNotCallableError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError, objectType)
	return instantiatedError
}

func (p *Plasma) RuntimeErrorInitialize(context *Context, object *Value) *Value {
	object.SetOnDemandSymbol(Initialize,
		func() *Value {
			return p.NewFunction(context, false, object.SymbolTable(),
				NewBuiltInClassFunction(object, 1,
					func(self *Value, arguments ...*Value) (*Value, bool) {
						message := arguments[0]
						if !message.IsTypeById(StringId) {
							return p.NewInvalidTypeError(context, message.TypeName(), StringName), false
						}
						self.SetString(message.GetString())
						return p.GetNone(), true
					},
				),
			)
		},
	)
	object.SetOnDemandSymbol(ToString,
		func() *Value {
			return p.NewFunction(context, false, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self *Value, _ ...*Value) (*Value, bool) {
						return p.NewString(context, false, context.PeekSymbolTable(), fmt.Sprintf("%s: %s", self.TypeName(), self.GetString())), true
					},
				),
			)
		},
	)
	return nil
}
