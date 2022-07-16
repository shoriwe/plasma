package vm

func (plasma *Plasma) stringClass() *Value {
	class := plasma.NewValue(plasma.rootSymbols, BuiltInClassId, plasma.class)
	class.SetAny(func(argument ...*Value) (*Value, error) {
		return plasma.NewString(argument[0].Contents()), nil
	})
	return class
}

/*
NewString magic function:
TODO Equals              __equals__
TODO NotEqual            __not_equal__
TODO Add                 __add__
TODO Mul                 __mul__
TODO Length              __len__
TODO Bool                __bool__
TODO String              __string__
TODO Bytes               __bytes__
TODO Array               __array__
TODO Tuple               __tuple__
TODO Get                 __get__
TODO Copy                __copy__
TODO Iter                __iter__
*/
func (plasma *Plasma) NewString(contents []byte) *Value {
	result := plasma.NewValue(plasma.rootSymbols, StringId, plasma.string)
	result.SetAny(contents)
	// TODO: init magic functions
	return result
}
