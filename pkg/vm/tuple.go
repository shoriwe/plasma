package vm

import (
	magic_functions "github.com/shoriwe/plasma/pkg/common/magic-functions"
)

func (plasma *Plasma) tupleClass() *Value {
	class := plasma.NewValue(plasma.rootSymbols, BuiltInClassId, plasma.class)
	class.SetAny(Callback(func(argument ...*Value) (*Value, error) {
		return plasma.NewTuple(argument[0].Values()), nil
	}))
	return class
}

/*
NewTuple Creates a new tuple Value
*/
func (plasma *Plasma) NewTuple(values []*Value) *Value {
	result := plasma.NewValue(plasma.rootSymbols, TupleId, plasma.tuple)
	result.SetAny(values)
	result.Set(magic_functions.In, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			for _, value := range result.GetValues() {
				if value.Equal(argument[0]) {
					return plasma.true, nil
				}
			}
			return plasma.false, nil
		}))
	result.Set(magic_functions.Equal, plasma.NewBuiltInFunction(
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
			rawString = append(rawString, '(')
			for index, value := range result.GetValues() {
				if index != 0 {
					rawString = append(rawString, ',', ' ')
				}
				rawString = append(rawString, value.String()...)
			}
			rawString = append(rawString, ')')
			return plasma.NewString(rawString), nil
		}))
	result.Set(magic_functions.Bytes, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			vs := result.GetValues()
			bytes := make([]byte, 0, len(vs))
			for _, v := range vs {
				bytes = append(bytes, Int[byte](v))
			}
			return plasma.NewBytes(bytes), nil
		}))
	result.Set(magic_functions.Array, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewArray(result.GetValues()), nil
		}))
	result.Set(magic_functions.Tuple, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return result, nil
		}))
	result.Set(magic_functions.Get, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId:
				s := result.Values()
				index := argument[0].GetInt64()
				if index < 0 {
					index += int64(len(s))
				}
				return s[index], nil
			case TupleId:
				s := result.Values()
				tupleIndex := argument[0].GetValues()
				var (
					startIndex int64
					endIndex   int64
				)
				if tupleIndex[0].TypeId() != NoneId {
					startIndex = tupleIndex[0].GetInt64()
					if startIndex < 0 {
						startIndex += int64(len(s))
					}
				} else {
					startIndex = 0
				}
				if len(tupleIndex) == 2 && tupleIndex[1].TypeId() != NoneId {
					endIndex = tupleIndex[1].GetInt64()
					if endIndex < 0 {
						endIndex += int64(len(s))
					}
				} else {
					endIndex = int64(len(s))
				}
				return plasma.NewTuple(s[startIndex:endIndex]), nil
			default:
				return nil, NotIndexable
			}
		}))
	result.Set(magic_functions.Iter, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			iter := plasma.NewValue(result.vtable, ValueId, plasma.value)
			iter.SetAny(int64(0))
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
