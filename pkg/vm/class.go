package vm

func (plasma *Plasma) metaClass() *Value {
	plasma.class = plasma.NewValue(plasma.rootSymbols, BuiltInClassId, plasma.class)
	plasma.class.class = plasma.class
	plasma.class.SetAny(func(argument ...*Value) (*Value, error) {
		return plasma.NewClass(), nil
	})
	return plasma.class
}

/*
NewClass magic function:
TODO Init                __init__
TODO HasNext             __has_next__
TODO Next                __next__
TODO Not                 __not__
TODO Positive            __positive__
TODO Negative            __negative__
TODO NegateBits          __negate_its__
TODO And                 __and__
TODO Or                  __or__
TODO Xor                 __xor__
TODO In                  __in__
TODO Is                  __is__
TODO Implements          __implements__
TODO Equals              __equals__
TODO NotEqual            __not_equal__
TODO GreaterThan         __greater_than__
TODO GreaterOrEqualThan  __greater_or_equal_than__
TODO LessThan            __less_than__
TODO LessOrEqualThan     __less_or_equal_than__
TODO BitwiseOr           __bitwise_or__
TODO BitwiseXor          __bitwise_xor__
TODO BitwiseAnd          __bitwise_and__
TODO BitwiseLeft         __bitwise_left__
TODO BitwiseRight        __bitwise_right__
TODO Add                 __add__
TODO Sub                 __sub__
TODO Mul                 __mul__
TODO Div                 __div__
TODO FloorDiv            __floor_div__
TODO Modulus             __mod__
TODO PowerOf             __pow__
TODO Length              __len__
TODO Bool                __bool__
TODO String              __string__
TODO Int                 __int__
TODO Float               __float__
TODO Bytes               __bytes__
TODO Array               __array__
TODO Tuple               __tuple__
TODO Get                 __get__
TODO Set                 __set__
TODO Del                 __del__
TODO Call                __call__
TODO Class               __class__
TODO Copy                __copy__
TODO Iter                __iter__
*/
func (plasma *Plasma) NewClass() *Value {
	result := plasma.NewValue(plasma.rootSymbols, BuiltInClassId, plasma.class)
	// TODO: init magic functions
	return result
}
