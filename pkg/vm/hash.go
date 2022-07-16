package vm

func (plasma *Plasma) hashClass() *Value {
	class := plasma.NewValue(plasma.rootSymbols, BuiltInClassId, plasma.class)
	class.SetAny(func(argument ...*Value) (*Value, error) {
		return plasma.NewHash(argument[0].GetHash()), nil
	})
	return class
}

/*
NewHash magic function:
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
TODO Del				 __del__
TODO Copy                __copy__
TODO Iter                __iter__
*/
func (plasma *Plasma) NewHash(hash *Hash) *Value {
	result := plasma.NewValue(plasma.rootSymbols, HashId, plasma.hash)
	result.SetAny(hash)
	// TODO: init magic functions
	return result
}
