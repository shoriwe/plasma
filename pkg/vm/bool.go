package vm

func (plasma *Plasma) boolClass() *Value {
	class := plasma.NewValue(plasma.rootSymbols, BuiltInClassId, plasma.class)
	class.SetAny(func(argument ...*Value) (*Value, error) {
		return plasma.NewBool(argument[0].Bool()), nil
	})
	return class
}

/*
NewBool magic function:
TODO Not                 __not__
TODO And                 __and__
TODO Or                  __or__
TODO Xor                 __xor__
TODO Equals              __equals__
TODO NotEqual            __not_equal__
TODO Bool                __bool__
TODO String              __string__
TODO Int                 __int__
TODO Float               __float__
TODO Bytes               __bytes__
TODO Call                __call__
TODO Copy                __copy__
*/
func (plasma *Plasma) NewBool(b bool) *Value {
	if b && plasma.true != nil {
		return plasma.true
	} else if !b && plasma.false != nil {
		return plasma.false
	}
	result := plasma.NewValue(plasma.rootSymbols, BoolId, plasma.bool)
	result.SetAny(b)
	// TODO: init magic functions
	return result
}
