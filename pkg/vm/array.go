package vm

import magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"

func (plasma *Plasma) arrayClass() *Value {
	class := plasma.NewValue(plasma.rootSymbols, BuiltInClassId, plasma.class)
	class.SetAny(func(argument ...*Value) (*Value, error) {
		return plasma.NewArray(argument[0].Values()), nil
	})
	return class
}

/*
NewArray magic function:
In                  __in__
Equals              __equals__
NotEqual            __not_equal__
Mul                 __mul__
Length              __len__
Bool                __bool__
String              __string__
Bytes               __bytes__
Array               __array__
Tuple               __tuple__
Get                 __get__
Set                 __set__
Iter                __iter__
*/
func (plasma *Plasma) NewArray(values []*Value) *Value {
	result := plasma.NewValue(plasma.rootSymbols, ArrayId, plasma.array)
	result.SetAny(values)
	result.Set(magic_functions.In, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			for _, value := range result.GetValues() {
				if !value.Equal(argument[0]) {
					return plasma.false, nil
				}
			}
			return plasma.true, nil
		}))
	result.Set(magic_functions.Equals, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case ArrayId:
				otherValues := argument[0].GetValues()
				for index, value := range result.GetValues() {
					if !value.Equal(otherValues[index]) {
						return plasma.false, nil
					}
				}
			}
			return plasma.true, nil
		}))
	result.Set(magic_functions.NotEqual, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case ArrayId:
				otherValues := argument[0].GetValues()
				for index, value := range result.GetValues() {
					if value.Equal(otherValues[index]) {
						return plasma.false, nil
					}
				}
			}
			return plasma.true, nil
		}))
	result.Set(magic_functions.Mul, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId:
				times := argument[0].GetInt64()
				currentValues := result.GetValues()
				newValues := make([]*Value, 0, int64(len(currentValues))*times)
				for i := int64(0); i < times; i++ {
					for _, value := range currentValues {
						newValues = append(newValues, value)
					}
				}
				return plasma.NewArray(newValues), nil
			default:
				return nil, NotOperable
			}
		}))
	result.Set(magic_functions.Length, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewInt(int64(len(result.GetValues()))), nil
		}))
	result.Set(magic_functions.Bool, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewBool(len(result.GetValues()) > 0), nil
		}))
	result.Set(magic_functions.String, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			var rawString []byte
			rawString = append(rawString, '[')
			for index, value := range result.GetValues() {
				if index != 0 {
					rawString = append(rawString, ',', ' ')
				}
				rawString = append(rawString, value.String()...)
			}
			rawString = append(rawString, ']')
			return plasma.NewString(rawString), nil
		}))
	result.Set(magic_functions.Bytes, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			var rawString []byte
			rawString = append(rawString, '[')
			for index, value := range result.GetValues() {
				if index != 0 {
					rawString = append(rawString, ',', ' ')
				}
				rawString = append(rawString, value.String()...)
			}
			rawString = append(rawString, ']')
			return plasma.NewBytes(rawString), nil
		}))
	result.Set(magic_functions.Tuple, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewTuple(result.GetValues()), nil
		}))
	result.Set(magic_functions.Get, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId:
				return result.GetValues()[argument[0].GetInt64()], nil
			default:
				return nil, NotIndexable
			}
		}))
	result.Set(magic_functions.Set, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId:
				result.GetValues()[argument[0].GetInt64()] = argument[1]
				return plasma.none, nil
			default:
				return nil, NotIndexable
			}
		}))
	result.Set(magic_functions.Iter, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			iter := plasma.NewValue(result.vtable, ValueId, plasma.value)
			iter.SetAny(0)
			iter.Set(magic_functions.HasNext, plasma.NewBuiltInFunction(iter.vtable,
				func(argument ...*Value) (*Value, error) {
					return plasma.NewBool(iter.GetInt64() < int64(len(result.GetValues()))), nil
				},
			))
			iter.Set(magic_functions.Next, plasma.NewBuiltInFunction(iter.vtable,
				func(argument ...*Value) (*Value, error) {
					currentValues := result.GetValues()
					index := iter.GetInt64()
					iter.SetAny(index + 1)
					if index < int64(len(currentValues)) {
						return currentValues[index], nil
					}
					return plasma.none, nil
				},
			))
			return iter, nil
		}))
	return result
}
