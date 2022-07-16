package vm

func (plasma *Plasma) floatClass() *Value {
	class := plasma.NewValue(plasma.rootSymbols, BuiltInClassId, plasma.class)
	class.SetAny(func(argument ...*Value) (*Value, error) {
		return plasma.NewFloat(argument[0].Float()), nil
	})
	return class
}

/*
NewFloat magic function:
TODO Positive:           __positive__
TODO Negative:           __negative__
TODO NegateBits:         __negate_its__
TODO Equals:             __equals__
TODO NotEqual:           __not_equal__
TODO GreaterThan:        __greater_than__
TODO GreaterOrEqualThan: __greater_or_equal_than__
TODO LessThan:           __less_than__
TODO LessOrEqualThan:    __less_or_equal_than__
TODO Add:                __add__
TODO Sub:                __sub__
TODO Mul:                __mul__
TODO Div:                __div__
TODO FloorDiv:           __floor_div__
TODO Modulus:            __mod__
TODO PowerOf:            __pow__
TODO Bool:               __bool__
TODO String             __string__
TODO Int                __int__
TODO Float              __float__
TODO Copy:               __copy__
*/
func (plasma *Plasma) NewFloat(f float64) *Value {
	result := plasma.NewValue(plasma.rootSymbols, FloatId, plasma.float)
	result.SetAny(f)
	// TODO: init magic functions
	return result
}
