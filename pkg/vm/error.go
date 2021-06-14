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
)

func (p *Plasma) ForceParentGetSelf(name string, parent *SymbolTable) IObject {
	object, getError := parent.GetSelf(name)
	if getError != nil {
		panic(getError.String())
	}
	return object
}

func (p *Plasma) ForceAnyGetAny(name string) IObject {
	object, getError := p.PeekSymbolTable().GetAny(name)
	if getError != nil {
		panic(getError.String())
	}
	return object
}

func (p *Plasma) ForceMasterGetAny(name string) IObject {
	object, getError := p.builtInSymbolTable.GetAny(name)
	if getError != nil {
		panic(getError.String())
	}
	return object
}

func (p *Plasma) ForceConstruction(type_ IObject) IObject {
	if _, ok := type_.(*Type); !ok {
		panic("don't modify built-ins, or fatal errors like this one will ocurr")
	}
	result, constructionError := p.ConstructObject(type_.(*Type), NewSymbolTable(p.PeekSymbolTable()))
	if constructionError != nil {
		panic(constructionError.typeName)
	}
	return result
}

func (p *Plasma) ForceInitialization(object IObject, arguments ...IObject) {
	initialize, getError := object.Get(Initialize)
	if getError != nil {
		panic(getError.String())
	}
	_, callError := p.CallFunction(
		initialize.(*Function), object.SymbolTable(),
		arguments...,
	)
	if callError != nil {
		panic(fmt.Sprintf("%s: %s", callError.TypeName(), callError.GetString()))
	}
}

func (p *Plasma) NewFloatParsingError() *Object {
	errorType := p.ForceMasterGetAny(FloatParsingError)
	instantiatedError := p.ForceConstruction(errorType)
	p.ForceInitialization(instantiatedError)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewIntegerParsingError() *Object {
	errorType := p.ForceMasterGetAny(IntegerParsingError)
	instantiatedError := p.ForceConstruction(errorType)
	p.ForceInitialization(instantiatedError)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewKeyNotFoundError(key IObject) *Object {
	errorType := p.ForceMasterGetAny(KeyNotFoundError)
	instantiatedError := p.ForceConstruction(errorType)
	p.ForceInitialization(instantiatedError,
		key,
	)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewIndexOutOfRange(length int, index int64) *Object {
	errorType := p.ForceMasterGetAny(IndexOutOfRangeError)
	instantiatedError := p.ForceConstruction(errorType)
	p.ForceInitialization(instantiatedError,
		p.NewInteger(
			false,
			p.PeekSymbolTable(),
			int64(length),
		),
		p.NewInteger(false, p.PeekSymbolTable(), index),
	)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewUnhashableTypeError(objectType *Type) *Object {
	errorType := p.ForceMasterGetAny(UnhashableTypeError)
	instantiatedError := p.ForceConstruction(errorType)
	p.ForceInitialization(instantiatedError,
		objectType,
	)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewNotImplementedCallableError(methodName string) *Object {
	errorType := p.ForceMasterGetAny(NotImplementedCallableError)
	instantiatedError := p.ForceConstruction(errorType)
	p.ForceInitialization(instantiatedError,
		p.NewString(false, p.PeekSymbolTable(), methodName),
	)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewGoRuntimeError(er error) *Object {
	errorType := p.ForceMasterGetAny(GoRuntimeError)
	instantiatedError := p.ForceConstruction(errorType)
	p.ForceInitialization(instantiatedError,
		p.NewString(false, p.PeekSymbolTable(), er.Error()),
	)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewInvalidNumberOfArgumentsError(received int, expecting int) *Object {
	errorType := p.ForceMasterGetAny(InvalidNumberOfArgumentsError)
	instantiatedError := p.ForceConstruction(errorType)
	p.ForceInitialization(instantiatedError,
		p.NewInteger(false, p.PeekSymbolTable(), int64(received)),
		p.NewInteger(false, p.PeekSymbolTable(), int64(expecting)),
	)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewObjectWithNameNotFoundError(name string) *Object {
	errorType := p.ForceMasterGetAny(ObjectWithNameNotFoundError)
	instantiatedError := p.ForceConstruction(errorType)
	p.ForceInitialization(instantiatedError,
		p.NewString(false, p.PeekSymbolTable(), name),
	)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewInvalidTypeError(received string, expecting ...string) *Object {
	errorType := p.ForceMasterGetAny(InvalidTypeError)
	instantiatedError := p.ForceConstruction(errorType)
	instantiatedErrorInitialize, _ := instantiatedError.Get(Initialize)
	expectingSum := ""
	for index, s := range expecting {
		if index != 0 {
			expectingSum += ", "
		}
		expectingSum += s
	}
	_, _ = p.CallFunction(
		instantiatedErrorInitialize.(*Function), instantiatedError.SymbolTable(),
		p.NewString(false, p.PeekSymbolTable(), received),
		p.NewString(false, p.PeekSymbolTable(), expectingSum),
	)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewObjectConstructionError(typeName string, errorMessage string) *Object {
	errorType := p.ForceMasterGetAny(ObjectConstructionError)
	instantiatedError := p.ForceConstruction(errorType)
	p.ForceInitialization(instantiatedError,
		p.NewString(false, p.PeekSymbolTable(), typeName), p.NewString(false, p.PeekSymbolTable(), errorMessage),
	)
	return instantiatedError.(*Object)
}

func (p *Plasma) NewBuiltInSymbolProtectionError(symbolName string) *Object {
	errorType := p.ForceMasterGetAny(BuiltInSymbolProtectionError)
	instantiatedError := p.ForceConstruction(errorType)
	p.ForceInitialization(instantiatedError,
		p.NewString(false, p.PeekSymbolTable(), symbolName),
	)
	return instantiatedError.(*Object)
}

func (p *Plasma) RuntimeErrorInitialize(object IObject) *Object {
	object.Set(Initialize,
		p.NewFunction(false, object.SymbolTable(),
			NewBuiltInClassFunction(object, 2,
				func(self IObject, arguments ...IObject) (IObject, *Object) {
					message := arguments[0]
					if _, ok := message.(*String); !ok {
						return nil, p.NewInvalidTypeError(message.TypeName(), StringName)
					}
					self.SetString(message.GetString())
					return p.NewNone(), nil
				},
			),
		),
	)
	object.Set(ToString,
		p.NewFunction(false, object.SymbolTable(),
			NewBuiltInClassFunction(object, 0,
				func(self IObject, _ ...IObject) (IObject, *Object) {
					return p.NewString(false, p.PeekSymbolTable(), fmt.Sprintf("%s: %s", self.TypeName(), self.GetString())), nil
				},
			),
		),
	)
	return nil
}
