package vm

import magic_functions "github.com/shoriwe/plasma/pkg/common/magic-functions"

func (plasma *Plasma) noneClass() *Value {
	class := plasma.NewValue(plasma.rootSymbols, BuiltInClassId, plasma.class)
	class.SetAny(Callback(func(argument ...*Value) (*Value, error) {
		return plasma.NewNone(), nil
	}))
	return class
}

/*
NewNone Creates a new none Value
*/
func (plasma *Plasma) NewNone() *Value {
	if plasma.none != nil {
		return plasma.none
	}
	result := plasma.NewValue(plasma.rootSymbols, NoneId, plasma.noneType)
	result.Set(magic_functions.Bool, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.false, nil
		},
	))
	result.Set(magic_functions.String, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewString([]byte(result.String())), nil
		},
	))
	return result
}
