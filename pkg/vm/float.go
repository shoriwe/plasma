package vm

import (
	"encoding/binary"
	magic_functions "github.com/shoriwe/plasma/pkg/common/magic-functions"
	"math"
)

func (plasma *Plasma) floatClass() *Value {
	class := plasma.NewValue(plasma.rootSymbols, BuiltInClassId, plasma.class)
	class.SetAny(Callback(func(argument ...*Value) (*Value, error) {
		return plasma.NewFloat(Float[float64](argument[0])), nil
	}))
	return class
}

/*
NewFloat creates a new float Value
*/
func (plasma *Plasma) NewFloat(f float64) *Value {
	result := plasma.NewValue(plasma.rootSymbols, FloatId, plasma.float)
	result.SetAny(f)
	result.Set(magic_functions.Positive, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return result, nil
		},
	))
	result.Set(magic_functions.Negative, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewFloat(-Float[float64](result)), nil
		},
	))
	result.Set(magic_functions.NegateBits, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewFloat(math.Float64frombits(^math.Float64bits(Float[float64](result)))), nil
		},
	))
	result.Set(magic_functions.Equal, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId, FloatId:
				return plasma.NewBool(result.Equal(argument[0])), nil
			}
			return plasma.false, nil
		},
	))
	result.Set(magic_functions.NotEqual, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId, FloatId:
				return plasma.NewBool(!result.Equal(argument[0])), nil
			}
			return plasma.true, nil
		},
	))
	result.Set(magic_functions.GreaterThan, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId, FloatId:
				return plasma.NewBool(Float[float64](result) > Float[float64](argument[0])), nil
			}
			return nil, NotComparable
		},
	))
	result.Set(magic_functions.GreaterOrEqualThan, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId, FloatId:
				return plasma.NewBool(Float[float64](result) >= Float[float64](argument[0])), nil
			}
			return nil, NotComparable
		},
	))
	result.Set(magic_functions.LessThan, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId, FloatId:
				return plasma.NewBool(Float[float64](result) < Float[float64](argument[0])), nil
			}
			return nil, NotComparable
		},
	))
	result.Set(magic_functions.LessOrEqualThan, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId, FloatId:
				return plasma.NewBool(Float[float64](result) <= Float[float64](argument[0])), nil
			}
			return nil, NotComparable
		},
	))
	result.Set(magic_functions.BitwiseOr, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId, FloatId:
				return plasma.NewFloat(
					math.Float64frombits(
						math.Float64bits(Float[float64](result)) | math.Float64bits(Float[float64](argument[0])),
					),
				), nil
			}
			return nil, NotOperable
		},
	))
	result.Set(magic_functions.BitwiseXor, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId, FloatId:
				return plasma.NewFloat(
					math.Float64frombits(
						math.Float64bits(Float[float64](result)) ^ math.Float64bits(Float[float64](argument[0])),
					),
				), nil
			}
			return nil, NotOperable
		},
	))
	result.Set(magic_functions.BitwiseAnd, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId, FloatId:
				return plasma.NewFloat(
					math.Float64frombits(
						math.Float64bits(Float[float64](result)) & math.Float64bits(Float[float64](argument[0])),
					),
				), nil
			}
			return nil, NotOperable
		},
	))
	result.Set(magic_functions.BitwiseLeft, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId, FloatId:
				return plasma.NewFloat(
					math.Float64frombits(
						math.Float64bits(Float[float64](result)) << math.Float64bits(Float[float64](argument[0])),
					),
				), nil
			}
			return nil, NotOperable
		},
	))
	result.Set(magic_functions.BitwiseRight, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId, FloatId:
				return plasma.NewFloat(
					math.Float64frombits(
						math.Float64bits(Float[float64](result)) >> math.Float64bits(Float[float64](argument[0])),
					),
				), nil
			}
			return nil, NotOperable
		},
	))
	result.Set(magic_functions.Add, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId, FloatId:
				return plasma.NewFloat(Float[float64](result) + Float[float64](argument[0])), nil
			}
			return nil, NotOperable
		},
	))
	result.Set(magic_functions.Sub, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId, FloatId:
				return plasma.NewFloat(Float[float64](result) - Float[float64](argument[0])), nil
			}
			return nil, NotOperable
		},
	))
	result.Set(magic_functions.Mul, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId, FloatId:
				return plasma.NewFloat(Float[float64](result) * Float[float64](argument[0])), nil
			}
			return nil, NotOperable
		},
	))
	result.Set(magic_functions.Div, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId, FloatId:
				return plasma.NewFloat(Float[float64](result) / Float[float64](argument[0])), nil
			}
			return nil, NotOperable
		},
	))
	result.Set(magic_functions.FloorDiv, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId, FloatId:
				return plasma.NewInt(int64(Float[float64](result) / Float[float64](argument[0]))), nil
			}
			return nil, NotOperable
		},
	))
	result.Set(magic_functions.Modulus, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId, FloatId:
				return plasma.NewFloat(math.Mod(Float[float64](result), Float[float64](argument[0]))), nil
			}
			return nil, NotOperable
		},
	))
	result.Set(magic_functions.PowerOf, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			switch argument[0].TypeId() {
			case IntId, FloatId:
				return plasma.NewFloat(math.Pow(Float[float64](result), Float[float64](argument[0]))), nil
			}
			return nil, NotOperable
		},
	))
	result.Set(magic_functions.Bool, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewBool(result.Bool()), nil
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
			return plasma.NewInt(Int[int64](result)), nil
		},
	))
	result.Set(magic_functions.Float, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return result, nil
		},
	))
	result.Set(magic_functions.Copy, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewFloat(Float[float64](result)), nil
		},
	))
	result.Set(magic_functions.BigEndian, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			b := make([]byte, 8)
			binary.BigEndian.PutUint64(b, math.Float64bits(Float[float64](result)))
			return plasma.NewBytes(b), nil
		},
	))
	result.Set(magic_functions.LittleEndian, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			b := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, math.Float64bits(Float[float64](result)))
			return plasma.NewBytes(b), nil
		},
	))
	result.Set(magic_functions.FromBig, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewFloat(math.Float64frombits(binary.BigEndian.Uint64(argument[0].GetBytes()))), nil
		},
	))
	result.Set(magic_functions.FromLittle, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewFloat(math.Float64frombits(binary.LittleEndian.Uint64(argument[0].GetBytes()))), nil
		},
	))
	return result
}
