package vm

import (
	"fmt"
	magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"
	"reflect"
)

type PlasmaCallback func(arg ...any) (any, error)

func (plasma *Plasma) ZeroCopyArray(vt *Symbols, v any) *Value {
	r := plasma.NewValue(vt, ValueId, plasma.Value())
	var get, length *Value
	switch array := v.(type) {
	case []int:
		get = plasma.NewBuiltInFunction(
			r.VirtualTable(),
			func(argument ...*Value) (*Value, error) {
				value := array[argument[0].Int()]
				return plasma.NewInt(int64(value)), nil
			})
		length = plasma.NewBuiltInFunction(
			r.VirtualTable(),
			func(argument ...*Value) (*Value, error) {
				return plasma.NewInt(int64(len(array))), nil
			})
	case []int32:
		get = plasma.NewBuiltInFunction(
			r.VirtualTable(),
			func(argument ...*Value) (*Value, error) {
				value := array[argument[0].Int()]
				return plasma.NewInt(int64(value)), nil
			})
		length = plasma.NewBuiltInFunction(
			r.VirtualTable(),
			func(argument ...*Value) (*Value, error) {
				return plasma.NewInt(int64(len(array))), nil
			})
	case []int64:
		get = plasma.NewBuiltInFunction(
			r.VirtualTable(),
			func(argument ...*Value) (*Value, error) {
				value := array[argument[0].Int()]
				return plasma.NewInt(value), nil
			})
		length = plasma.NewBuiltInFunction(
			r.VirtualTable(),
			func(argument ...*Value) (*Value, error) {
				return plasma.NewInt(int64(len(array))), nil
			})
	case []uint:
		get = plasma.NewBuiltInFunction(
			r.VirtualTable(),
			func(argument ...*Value) (*Value, error) {
				value := array[argument[0].Int()]
				return plasma.NewInt(int64(value)), nil
			})
		length = plasma.NewBuiltInFunction(
			r.VirtualTable(),
			func(argument ...*Value) (*Value, error) {
				return plasma.NewInt(int64(len(array))), nil
			})
	case []uint32:
		get = plasma.NewBuiltInFunction(
			r.VirtualTable(),
			func(argument ...*Value) (*Value, error) {
				value := array[argument[0].Int()]
				return plasma.NewInt(int64(value)), nil
			})
		length = plasma.NewBuiltInFunction(
			r.VirtualTable(),
			func(argument ...*Value) (*Value, error) {
				return plasma.NewInt(int64(len(array))), nil
			})
	case []uint64:
		get = plasma.NewBuiltInFunction(
			r.VirtualTable(),
			func(argument ...*Value) (*Value, error) {
				value := array[argument[0].Int()]
				return plasma.NewInt(int64(value)), nil
			})
		length = plasma.NewBuiltInFunction(
			r.VirtualTable(),
			func(argument ...*Value) (*Value, error) {
				return plasma.NewInt(int64(len(array))), nil
			})
	case []float32:
		get = plasma.NewBuiltInFunction(
			r.VirtualTable(),
			func(argument ...*Value) (*Value, error) {
				value := array[argument[0].Int()]
				return plasma.NewFloat(float64(value)), nil
			})
		length = plasma.NewBuiltInFunction(
			r.VirtualTable(),
			func(argument ...*Value) (*Value, error) {
				return plasma.NewInt(int64(len(array))), nil
			})
	case []float64:
		get = plasma.NewBuiltInFunction(
			r.VirtualTable(),
			func(argument ...*Value) (*Value, error) {
				value := array[argument[0].Int()]
				return plasma.NewFloat(value), nil
			})
		length = plasma.NewBuiltInFunction(
			r.VirtualTable(),
			func(argument ...*Value) (*Value, error) {
				return plasma.NewInt(int64(len(array))), nil
			})
	default:
		panic("unknown type")
	}
	r.Set(magic_functions.Get, get)
	r.Set(magic_functions.Length, length)
	return r
}

func (plasma *Plasma) FromValue(value *Value) (any, error) {
	switch id := value.TypeId(); id {
	case ValueId:
		value.mutex.Lock()
		defer value.mutex.Unlock()
		r := make(map[string]any, len(value.vtable.values))
		for key, objValue := range value.vtable.values {
			v, err := plasma.FromValue(objValue)
			if err != nil {
				return nil, err
			}
			r[key] = v
		}
		return r, nil
	case StringId:
		return value.String(), nil
	case BytesId:
		return value.GetBytes(), nil
	case BoolId:
		return value.Bool(), nil
	case NoneId:
		return nil, nil
	case IntId:
		return value.Int(), nil
	case FloatId:
		return value.Float(), nil
	case ArrayId, TupleId:
		values := value.GetValues()
		r := make([]any, 0, len(values))
		for _, arrayValue := range values {
			v, err := plasma.FromValue(arrayValue)
			if err != nil {
				return nil, err
			}
			r = append(r, v)
		}
		return r, nil
	case HashId:
		return nil, fmt.Errorf("built-in functions cannot cannot be converted to Go object")
	case BuiltInFunctionId:
		return nil, fmt.Errorf("built-in functions cannot cannot be converted to Go object")
	case FunctionId:
		return nil, fmt.Errorf("plasma functions cannot cannot be converted to Go object")
	case BuiltInClassId:
		return nil, fmt.Errorf("built-in class cannot cannot be converted to Go object")
	case ClassId:
		return nil, fmt.Errorf("plasma class cannot cannot be converted to Go object")
	default:
		return nil, fmt.Errorf("unknown value with type id %d", id)
	}
}

func (plasma *Plasma) toValueGoFunctionCall(symbols *Symbols, function reflect.Value) func(argument ...*Value) (*Value, error) {
	return func(argument ...*Value) (*Value, error) {
		functionType := function.Type()
		if len(argument) != functionType.NumIn() {
			return nil, fmt.Errorf("expecting %d arguments", functionType.NumIn())
		}
		callArguments := make([]reflect.Value, 0, len(argument))
		for argIndex, arg := range argument {
			asGoValue, err := plasma.FromValue(arg)
			if err != nil {
				return nil, err
			}
			targetType := functionType.In(argIndex)
			argumentAsReflectValue := reflect.ValueOf(asGoValue)
			asMapStringAny, asMapStringAnyOk := asGoValue.(map[string]any)
			switch {
			case argumentAsReflectValue.Type() == targetType:
				break
			case argumentAsReflectValue.CanConvert(targetType):
				argumentAsReflectValue = argumentAsReflectValue.Convert(targetType)
			case asMapStringAnyOk:
				nested := 0
				for ; targetType.Kind() == reflect.Pointer; targetType = targetType.Elem() {
					nested++
				}
				argumentAsReflectValue = reflect.New(targetType).Elem()
				for fieldIndex := 0; fieldIndex < targetType.NumField(); fieldIndex++ {
					asMapStringAnyValue, found := asMapStringAny[targetType.Field(fieldIndex).Name]
					if !found {
						return nil, fmt.Errorf("field with name %s not found in plasma object", targetType.Field(fieldIndex).Name)
					}
					argumentAsReflectValue.FieldByName(targetType.Field(fieldIndex).Name).Set(reflect.ValueOf(asMapStringAnyValue))
				}
				for doNested := 0; doNested < nested; doNested++ {
					argumentAsReflectValue = argumentAsReflectValue.Addr()
				}
			default:
				// TODO: Fix me this should work in any scenario
				return nil, fmt.Errorf("cannot convert unknown type of Go Obj")
			}
			callArguments = append(callArguments, argumentAsReflectValue)
		}
		result := function.Call(callArguments)
		if len(result) == 0 {
			return plasma.None(), nil
		}
		plasmaResult := make([]*Value, 0, len(result))
		for _, r := range result {
			asPlasma, err := plasma.ToValue(symbols, r.Interface())
			if err != nil {
				return nil, err
			}
			plasmaResult = append(plasmaResult, asPlasma)
		}
		if len(plasmaResult) == 1 {
			return plasmaResult[0], nil
		}
		return plasma.NewArray(plasmaResult), nil
	}
}

func (plasma *Plasma) ToValue(symbols *Symbols, v any) (*Value, error) {
	if v == nil {
		return plasma.None(), nil
	}
	if symbols == nil {
		symbols = plasma.Symbols()
	}
	asReflectValue := reflect.ValueOf(v)
	switch asReflectValue.Kind() {
	case reflect.String:
		return plasma.NewString([]byte(v.(string))), nil
	case reflect.Bool:
		return plasma.NewBool(v.(bool)), nil
	case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return plasma.NewInt(int64(asReflectValue.Uint())), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return plasma.NewInt(asReflectValue.Int()), nil
	case reflect.Float32, reflect.Float64:
		return plasma.NewFloat(asReflectValue.Float()), nil
	case reflect.Complex64, reflect.Complex128:
		panic("Complex not supported by the VM")
	case reflect.Slice, reflect.Array:
		if asReflectValue.Type() == reflect.TypeOf([]byte{}) {
			return plasma.NewBytes(v.([]byte)), nil
		}
		values := make([]*Value, 0, asReflectValue.Len())
		for index := 0; index < asReflectValue.Len(); index++ {
			value, err := plasma.ToValue(symbols, asReflectValue.Index(index).Interface())
			if err != nil {
				return nil, fmt.Errorf("transform error at index %d: %w", index, err)
			}
			values = append(values, value)
		}
		return plasma.NewArray(values), nil
	case reflect.Map:
		keys := asReflectValue.MapKeys()
		hash := plasma.NewInternalHash()
		hash.internalMap = make(map[string]*Value, len(keys))
		for _, key := range keys {
			keyV, keyErr := plasma.ToValue(symbols, key.Interface())
			if keyErr != nil {
				return nil, fmt.Errorf("transform key %v error: %w", key, keyErr)
			}
			value := asReflectValue.MapIndex(key)
			valueV, valueErr := plasma.ToValue(symbols, value.Interface())
			if valueErr != nil {
				return nil, fmt.Errorf("transform key %v error: %w", key, valueErr)
			}
			setErr := hash.Set(keyV, valueV)
			if setErr != nil {
				return nil, fmt.Errorf("transform error: %w", setErr)
			}
		}
		return plasma.NewHash(hash), nil
	case reflect.Struct:
		nFields := asReflectValue.NumField()
		obj := plasma.NewValue(symbols, ValueId, plasma.Value())
		for i := 0; i < nFields; i++ {
			if !asReflectValue.Field(i).CanInterface() {
				continue
			}
			fieldValue, err := plasma.ToValue(symbols, asReflectValue.Field(i).Interface())
			if err != nil {
				return nil, fmt.Errorf("transform struct error at field %s: %w", asReflectValue.Type().Field(i).Name, err)
			}
			obj.Set(asReflectValue.Type().Field(i).Name, fieldValue)
		}
		return obj, nil
	case reflect.Func:
		return plasma.NewBuiltInFunction(symbols, plasma.toValueGoFunctionCall(symbols, asReflectValue)), nil
	case reflect.Pointer:
		return plasma.ToValue(symbols, asReflectValue.Elem().Interface())
	case reflect.Chan:
		obj := plasma.NewValue(symbols, ValueId, plasma.Value())
		obj.Set("recv", plasma.NewBuiltInFunction(obj.VirtualTable(),
			func(argument ...*Value) (*Value, error) {
				vv, ok := asReflectValue.Recv()
				if !ok {
					return nil, fmt.Errorf("channel is closed")
				}
				return plasma.ToValue(obj.VirtualTable(), vv.Interface())
			},
		))
		obj.Set("send", plasma.NewBuiltInFunction(obj.VirtualTable(),
			func(argument ...*Value) (*Value, error) {
				vv, err := plasma.FromValue(argument[0])
				if err != nil {
					return nil, err
				}
				vvAsReflectValue := reflect.ValueOf(vv)
				if vvAsReflectValue.Type() == asReflectValue.Type().Elem() {

				} else if vvAsReflectValue.Type() != asReflectValue.Type().Elem() && vvAsReflectValue.CanConvert(asReflectValue.Type().Elem()) {
					vvAsReflectValue = vvAsReflectValue.Convert(asReflectValue.Type().Elem())
				} else {
					return nil, fmt.Errorf("cannot convert type inside channel")
				}
				asReflectValue.Send(vvAsReflectValue)
				return plasma.None(), nil
			},
		))
		return obj, nil
	case reflect.UnsafePointer:
		return plasma.ToValue(symbols, uintptr(asReflectValue.UnsafePointer()))
	case reflect.Interface:
		return nil, fmt.Errorf("cannot convert to plasma object")

	default:
		return nil, fmt.Errorf("cannot convert to plasma object")
	}
}
