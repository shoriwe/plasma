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

func (p *Plasma) ForceGetSelf(name string, parent Value) Value {
	object, getError := parent.Get(name)
	if getError != nil {
		panic(getError.String())
	}
	return object
}

func (p *Plasma) ForceAnyGetAny(context *Context, name string) Value {
	object, getError := context.PeekSymbolTable().GetAny(name)
	if getError != nil {
		panic(getError.String())
	}
	return object
}

func (p *Plasma) ForceMasterGetAny(name string) Value {
	object, getError := p.builtInContext.PeekSymbolTable().GetAny(name)
	if getError != nil {
		panic(getError.String())
	}
	return object
}

func (p *Plasma) ForceConstruction(context *Context, type_ Value) Value {
	if _, ok := type_.(*Type); !ok {
		panic("don't modify built-ins, or fatal errors like this one will ocurr")
	}
	result, constructionError := p.ConstructObject(context, type_.(*Type), NewSymbolTable(context.PeekSymbolTable()))
	if constructionError != nil {
		panic(constructionError.typeName)
	}
	return result
}

func (p *Plasma) ForceInitialization(context *Context, object Value, arguments ...Value) {
	initialize, getError := object.Get(Initialize)
	if getError != nil {
		panic(getError.String())
	}
	_, callError := p.CallFunction(context,
		initialize, object.SymbolTable(),
		arguments...,
	)
	if callError != nil {
		panic(fmt.Sprintf("%s: %s", callError.TypeName(), callError.GetString()))
	}
}

func (p *Plasma) NewFloatParsingError(context *Context) *Object {
	errorType := p.ForceMasterGetAny(FloatParsingError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewIntegerParsingError(context *Context) *Object {
	errorType := p.ForceMasterGetAny(IntegerParsingError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewKeyNotFoundError(context *Context, key Value) *Object {
	errorType := p.ForceMasterGetAny(KeyNotFoundError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError,
		key,
	)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewIndexOutOfRange(context *Context, length int, index int64) *Object {
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
	return instantiatedError.(*Object)
}

func (p *Plasma) NewUnhashableTypeError(context *Context, objectType *Type) *Object {
	errorType := p.ForceMasterGetAny(UnhashableTypeError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError,
		objectType,
	)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewNotImplementedCallableError(context *Context, methodName string) *Object {
	errorType := p.ForceMasterGetAny(NotImplementedCallableError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError,
		p.NewString(context, false, context.PeekSymbolTable(), methodName),
	)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewGoRuntimeError(context *Context, er error) *Object {
	errorType := p.ForceMasterGetAny(GoRuntimeError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError,
		p.NewString(context, false, context.PeekSymbolTable(), er.Error()),
	)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewInvalidNumberOfArgumentsError(context *Context, received int, expecting int) *Object {
	errorType := p.ForceMasterGetAny(InvalidNumberOfArgumentsError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError,
		p.NewInteger(context, false, context.PeekSymbolTable(), int64(received)),
		p.NewInteger(context, false, context.PeekSymbolTable(), int64(expecting)),
	)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewObjectWithNameNotFoundError(context *Context, objectType *Type, name string) *Object {
	errorType := p.ForceMasterGetAny(ObjectWithNameNotFoundError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError,
		objectType, p.NewString(context, false, context.PeekSymbolTable(), name),
	)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewInvalidTypeError(context *Context, received string, expecting ...string) *Object {
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
		instantiatedErrorInitialize, instantiatedError.SymbolTable(),
		p.NewString(context, false, context.PeekSymbolTable(), received),
		p.NewString(context, false, context.PeekSymbolTable(), expectingSum),
	)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewObjectConstructionError(context *Context, typeName string, errorMessage string) *Object {
	errorType := p.ForceMasterGetAny(ObjectConstructionError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError,
		p.NewString(context, false, context.PeekSymbolTable(), typeName), p.NewString(context, false, context.PeekSymbolTable(), errorMessage),
	)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewBuiltInSymbolProtectionError(context *Context, symbolName string) *Object {
	errorType := p.ForceMasterGetAny(BuiltInSymbolProtectionError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError,
		p.NewString(context, false, context.PeekSymbolTable(), symbolName),
	)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewObjectNotCallable(context *Context, objectType *Type) *Object {
	errorType := p.ForceMasterGetAny(ObjectNotCallableError)
	instantiatedError := p.ForceConstruction(context, errorType)
	p.ForceInitialization(context, instantiatedError, objectType)
	return instantiatedError.(*Object)
}

func (p *Plasma) RuntimeErrorInitialize(context *Context, object Value) *Object {
	object.SetOnDemandSymbol(Initialize,
		func() Value {
			return p.NewFunction(context, false, object.SymbolTable(),
				NewBuiltInClassFunction(object, 2,
					func(self Value, arguments ...Value) (Value, *Object) {
						message := arguments[0]
						if _, ok := message.(*String); !ok {
							return nil, p.NewInvalidTypeError(context, message.TypeName(), StringName)
						}
						self.SetString(message.GetString())
						return p.GetNone(), nil
					},
				),
			)
		},
	)
	object.SetOnDemandSymbol(ToString,
		func() Value {
			return p.NewFunction(context, false, object.SymbolTable(),
				NewBuiltInClassFunction(object, 0,
					func(self Value, _ ...Value) (Value, *Object) {
						return p.NewString(context, false, context.PeekSymbolTable(), fmt.Sprintf("%s: %s", self.TypeName(), self.GetString())), nil
					},
				),
			)
		},
	)
	return nil
}
