package vm

type Type struct {
	*Object
	Constructor Constructor
	Name        string
}

func (t *Type) Implements(class *Type) bool {
	if t == class {
		return true
	}
	for _, subClass := range t.subClasses {
		if subClass.Implements(class) {
			return true
		}
	}
	return false
}

func (p *Plasma) NewType(
	isBuiltIn bool,
	typeName string,
	parent *SymbolTable,
	subclasses []*Type,
	constructor Constructor,
) *Type {
	result := &Type{
		Object:      p.NewObject(isBuiltIn, TypeName, subclasses, parent),
		Constructor: constructor,
		Name:        typeName,
	}
	result.SetOnDemandSymbol(ToString,
		func() Value {
			return p.NewFunction(isBuiltIn, result.symbols,
				NewBuiltInClassFunction(result, 0,
					func(_ Value, _ ...Value) (Value, *Object) {
						return p.NewString(false, p.PeekSymbolTable(), "Type@"+typeName), nil
					},
				),
			)
		},
	)
	result.SetOnDemandSymbol(Self,
		func() Value {
			return result
		},
	)
	return result
}
