package vm

import (
	magic_functions "github.com/shoriwe/plasma/pkg/common/magic-functions"
)

func (plasma *Plasma) arrayClass() *Value {
	class := plasma.NewValue(plasma.rootSymbols, BuiltInClassId, plasma.class)
	class.SetAny(
		Callback(func(argument ...*Value) (*Value, error) {
			return plasma.NewArray(argument[0].Values()), nil
		}),
	)
	return class
}

/*
NewArray Creates a new array Value
*/
func (plasma *Plasma) NewArray(values []*Value) *Value {
	result := plasma.NewValue(plasma.rootSymbols, ArrayId, plasma.array)
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
			return plasma.NewString(result.Bytes()), nil
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
			return result, nil
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
				return plasma.NewArray(s[startIndex:endIndex]), nil
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
	result.Set(magic_functions.Append, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			result.SetAny(append(result.GetValues(), argument[0]))
			return plasma.none, nil
		},
	))
	result.Set(magic_functions.Clear, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			result.SetAny([]*Value{})
			return plasma.none, nil
		},
	))
	result.Set(magic_functions.Index, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			for index, value := range result.GetValues() {
				if value.Equal(argument[0]) {
					return plasma.NewInt(int64(index)), nil
				}
			}
			return plasma.NewInt(-1), nil
		},
	))
	result.Set(magic_functions.Pop, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			currentValues := result.GetValues()
			r := currentValues[len(currentValues)-1]
			currentValues = currentValues[:len(currentValues)-1]
			result.SetAny(currentValues)
			return r, nil
		},
	))
	result.Set(magic_functions.Insert, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			index := Int[int64](argument[0])
			value := argument[1]
			currentValues := result.GetValues()
			newValues := make([]*Value, 0, 1+int64(len(currentValues)))
			newValues = append(newValues, currentValues[:index]...)
			newValues = append(newValues, value)
			newValues = append(newValues, currentValues[index:]...)
			result.SetAny(newValues)
			return plasma.none, nil
		},
	))
	result.Set(magic_functions.Remove, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			index := Int[int64](argument[0])
			currentValues := result.GetValues()
			newValues := make([]*Value, 0, 1+int64(len(currentValues)))
			newValues = append(newValues, currentValues[:index]...)
			newValues = append(newValues, currentValues[index+1:]...)
			result.SetAny(newValues)
			return plasma.none, nil
		},
	))
	return result
}
