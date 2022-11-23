package vm

import (
	"bytes"
	magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"
)

func (plasma *Plasma) bytesClass() *Value {
	class := plasma.NewValue(plasma.rootSymbols, BuiltInClassId, plasma.class)
	class.SetAny(Callback(func(argument ...*Value) (*Value, error) {
		return plasma.NewBytes(argument[0].Bytes()), nil
	}))
	return class
}

/*
NewBytes Creates a new bytes Value
*/
func (plasma *Plasma) NewBytes(contents []byte) *Value {
	result := plasma.NewValue(plasma.rootSymbols, BytesId, plasma.bytes)
	result.SetAny(contents)
	result.Set(magic_functions.In, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case BytesId:
				return plasma.NewBool(bytes.Contains(result.GetBytes(), argument[0].GetBytes())), nil
			case IntId:
				i := argument[0].GetInt64()
				for _, b := range result.GetBytes() {
					if int64(b) == i {
						return plasma.true, nil
					}
				}
				return plasma.false, nil
			}
			return plasma.false, nil
		},
	))
	result.Set(magic_functions.Equal, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewBool(result.Equal(argument[0])), nil
		},
	))
	result.Set(magic_functions.NotEqual, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewBool(!result.Equal(argument[0])), nil
		},
	))
	result.Set(magic_functions.Add, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case BytesId:
				s := result.GetBytes()
				otherS := argument[0].GetBytes()
				newString := make([]byte, 0, len(s)+len(otherS))
				newString = append(newString, s...)
				newString = append(newString, otherS...)
				return plasma.NewBytes(newString), nil
			}
			return nil, NotOperable
		},
	))
	result.Set(magic_functions.Mul, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId:
				s := result.GetBytes()
				times := argument[0].GetInt64()
				return plasma.NewBytes(bytes.Repeat(s, int(times))), nil
			}
			return nil, NotOperable
		},
	))
	result.Set(magic_functions.Length, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewInt(int64(len(result.GetBytes()))), nil
		},
	))
	result.Set(magic_functions.Bool, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewBool(len(result.GetBytes()) > 0), nil
		},
	))
	result.Set(magic_functions.String, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewString(result.GetBytes()), nil
		},
	))
	result.Set(magic_functions.Bytes, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return result, nil
		},
	))
	result.Set(magic_functions.Array, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			s := result.GetBytes()
			values := make([]*Value, 0, len(s))
			for _, b := range s {
				values = append(values, plasma.NewInt(int64(b)))
			}
			return plasma.NewArray(values), nil
		},
	))
	result.Set(magic_functions.Tuple, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			s := result.GetBytes()
			values := make([]*Value, 0, len(s))
			for _, b := range s {
				values = append(values, plasma.NewInt(int64(b)))
			}
			return plasma.NewTuple(values), nil
		},
	))
	result.Set(magic_functions.Get, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId:
				s := result.GetBytes()
				index := argument[0].GetInt64()
				if index < 0 {
					index += int64(len(s))
				}
				return plasma.NewInt(int64(s[index])), nil
			case TupleId:
				s := result.GetBytes()
				tupleIndex := argument[0].GetValues()
				startIndex := tupleIndex[0].GetInt64()
				if startIndex < 0 {
					startIndex += int64(len(s))
				}
				var endIndex int64
				if len(tupleIndex) == 2 {
					endIndex = tupleIndex[1].GetInt64()
				} else {
					endIndex = int64(len(s))
				}
				if endIndex < 0 {
					endIndex += int64(len(s))
				}
				return plasma.NewBytes(s[startIndex:endIndex]), nil
			default:
				return nil, NotIndexable
			}
		},
	))
	result.Set(magic_functions.Copy, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			s := result.GetBytes()
			newS := make([]byte, len(s))
			copy(newS, s)
			return plasma.NewBytes(newS), nil
		},
	))
	result.Set(magic_functions.Iter, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			iter := plasma.NewValue(result.vtable, ValueId, plasma.value)
			iter.SetAny(int64(0))
			iter.Set(magic_functions.HasNext, plasma.NewBuiltInFunction(iter.vtable,
				func(argument ...*Value) (*Value, error) {
					return plasma.NewBool(iter.GetInt64() < int64(len(result.GetBytes()))), nil
				},
			))
			iter.Set(magic_functions.Next, plasma.NewBuiltInFunction(iter.vtable,
				func(argument ...*Value) (*Value, error) {
					currentBytes := result.GetBytes()
					index := iter.GetInt64()
					iter.SetAny(index + 1)
					if index < int64(len(currentBytes)) {
						return plasma.NewBytes([]byte{currentBytes[index]}), nil
					}
					return plasma.none, nil
				},
			))
			return iter, nil
		},
	))
	result.Set(magic_functions.Join, plasma.NewBuiltInFunction(result.vtable,
		func(argument ...*Value) (*Value, error) {
			values := argument[0].Values()
			valuesBytes := make([][]byte, 0, len(values))
			for _, value := range values {
				valuesBytes = append(valuesBytes, []byte(value.String()))
			}
			return plasma.NewBytes(bytes.Join(valuesBytes, []byte(result.String()))), nil
		},
	))
	result.Set(magic_functions.Split, plasma.NewBuiltInFunction(result.vtable,
		func(argument ...*Value) (*Value, error) {
			sep := argument[0].String()
			splitted := bytes.Split(result.GetBytes(), []byte(sep))
			values := make([]*Value, 0, len(splitted))
			for _, b := range splitted {
				values = append(values, plasma.NewBytes(b))
			}
			return plasma.NewTuple(values), nil
		},
	))
	result.Set(magic_functions.Upper, plasma.NewBuiltInFunction(result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewBytes(bytes.ToUpper(result.GetBytes())), nil
		},
	))
	result.Set(magic_functions.Lower, plasma.NewBuiltInFunction(result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewBytes(bytes.ToLower(result.GetBytes())), nil
		},
	))
	result.Set(magic_functions.Count, plasma.NewBuiltInFunction(result.vtable,
		func(argument ...*Value) (*Value, error) {
			sep := argument[0].String()
			return plasma.NewInt(int64(bytes.Count(result.GetBytes(), []byte(sep)))), nil
		},
	))
	result.Set(magic_functions.Index, plasma.NewBuiltInFunction(result.vtable,
		func(argument ...*Value) (*Value, error) {
			sep := argument[0].String()
			return plasma.NewInt(int64(bytes.Index(result.GetBytes(), []byte(sep)))), nil
		},
	))
	return result
}
