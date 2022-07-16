package vm

func (plasma *Plasma) arrayClass() *Value {
	class := plasma.NewValue(plasma.rootSymbols, BuiltInClassId, plasma.class)
	class.SetAny(func(argument ...*Value) (*Value, error) {
		return plasma.NewArray(argument[0].Values()), nil
	})
	return class
}

/*
NewArray magic function:
TODO In                  __in__
TODO Equals              __equals__
TODO NotEqual            __not_equal__
TODO Mul                 __mul__
TODO Length              __len__
TODO Bool                __bool__
TODO String              __string__
TODO Bytes               __bytes__
TODO Array               __array__
TODO Tuple               __tuple__
TODO Get                 __get__
TODO Set                 __set__
TODO Copy                __copy__
TODO Iter                __iter__
*/
func (plasma *Plasma) NewArray(values []*Value) *Value {
	result := plasma.NewValue(plasma.rootSymbols, ArrayId, plasma.array)
	result.SetAny(values)
	// TODO: init magic functions
	return result
}
