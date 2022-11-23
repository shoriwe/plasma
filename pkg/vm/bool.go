package vm

import magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"

func (plasma *Plasma) boolClass() *Value {
	class := plasma.NewValue(plasma.rootSymbols, BuiltInClassId, plasma.class)
	class.SetAny(Callback(func(argument ...*Value) (*Value, error) {
		return plasma.NewBool(argument[0].Bool()), nil
	}))
	return class
}

/*
NewBool Creates a new bool Value
*/
func (plasma *Plasma) NewBool(b bool) *Value {
	if b && plasma.true != nil {
		return plasma.true
	} else if !b && plasma.false != nil {
		return plasma.false
	}
	result := plasma.NewValue(plasma.rootSymbols, BoolId, plasma.bool)
	result.SetAny(b)
	result.Set(magic_functions.Not, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewBool(!result.GetBool()), nil
		},
	))
	result.Set(magic_functions.Equal, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case BoolId:
				return plasma.NewBool(result.GetBool() == argument[0].GetBool()), nil
			}
			return plasma.false, nil
		},
	))
	result.Set(magic_functions.NotEqual, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case BoolId:
				return plasma.NewBool(result.GetBool() != argument[0].GetBool()), nil
			}
			return plasma.true, nil
		},
	))
	result.Set(magic_functions.Bool, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return result, nil
		},
	))
	result.Set(magic_functions.String, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewString([]byte(result.String())), nil
		},
	))
	result.Set(magic_functions.Int, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			if result.GetBool() {
				return plasma.NewInt(1), nil
			}
			return plasma.NewInt(0), nil
		},
	))
	result.Set(magic_functions.Float, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			if result.GetBool() {
				return plasma.NewFloat(1), nil
			}
			return plasma.NewFloat(0), nil
		},
	))
	result.Set(magic_functions.Bytes, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewBytes([]byte(result.String())), nil
		},
	))
	result.Set(magic_functions.Copy, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return result, nil
		},
	))
	return result
}
