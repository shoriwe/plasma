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
		hash := value.GetHash()
		hash.mutex.Lock()
		defer hash.mutex.Unlock()
		result := make(map[any]any, len(hash.internalMap))
		for _, keyValue := range hash.internalMap {
			key, err := plasma.FromValue(keyValue.Key)
			if err != nil {
				return nil, err
			}
			v, err := plasma.FromValue(keyValue.Value)
			if err != nil {
				return nil, err
			}
			result[key] = v
		}
		return result, nil
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

func (plasma *Plasma) callGoFunc(symbols *Symbols, function reflect.Value, arguments ...reflect.Value) (*Value, error) {
	result := function.Call(arguments)
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

func reflectConvert(v reflect.Value, t reflect.Type) (reflect.Value, error) {
	var result reflect.Value
	nested := 0
	targetType := t
	for ; targetType.Kind() == reflect.Pointer; targetType = targetType.Elem() {
		nested++
	}
	vType := v.Type()
	switch {
	case v.Type() == targetType:
		result = v
	case v.CanConvert(targetType):
		result = v.Convert(targetType)
	case v.Kind() == reflect.Map && vType.Key().Kind() == reflect.String && vType.Elem().Kind() == reflect.Interface && targetType.Kind() == reflect.Struct:
		result = reflect.New(targetType).Elem()
		for fieldIndex := 0; fieldIndex < targetType.NumField(); fieldIndex++ {
			field := targetType.Field(fieldIndex)
			value := v.MapIndex(reflect.ValueOf(field.Name))
			if value.IsZero() {
				return reflect.Value{}, fmt.Errorf("field with name %s not found in plasma object", targetType.Field(fieldIndex).Name)
			}
			result.FieldByName(field.Name).Set(reflect.ValueOf(value.Interface()))
		}
	default:
		// TODO: Fix me this should work in any scenario
		return reflect.Value{}, fmt.Errorf("cannot convert %s to %s", v.Type(), targetType)
	}
	for doNested := 0; doNested < nested; doNested++ {
		result = result.Addr()
	}
	return result, nil
}

func (plasma *Plasma) toValueGoFunctionCall(symbols *Symbols, function reflect.Value) func(argument ...*Value) (*Value, error) {
	functionType := function.Type()
	numIn := functionType.NumIn()
	isVariadic := functionType.IsVariadic()
	if numIn == 0 {
		return func(argument ...*Value) (*Value, error) {
			return plasma.callGoFunc(symbols, function)
		}
	}
	return func(arguments ...*Value) (*Value, error) {
		if len(arguments) != numIn && !isVariadic {
			return nil, fmt.Errorf("expecting %d arguments but recieved %d", numIn, len(arguments))
		}
		callArguments := make([]reflect.Value, 0, numIn)
		var lastIndex int
		if isVariadic {
			lastIndex = numIn - 1
		} else {
			lastIndex = numIn
		}
		index := 0
		for ; index < lastIndex; index++ {
			argument := arguments[index]
			asGoValue, asGoValueError := plasma.FromValue(argument)
			if asGoValueError != nil {
				return nil, asGoValueError
			}
			argumentAsReflectValue, convertErr := reflectConvert(reflect.ValueOf(asGoValue), functionType.In(index))
			if convertErr != nil {
				return nil, convertErr
			}

			callArguments = append(callArguments, argumentAsReflectValue)
		}
		// Insert values to variadic
		if index < numIn {
			variadicArguments := make([]reflect.Value, 0, numIn-index)
			variadicArgumentType := functionType.In(index).Elem()
			for ; index < len(arguments); index++ {
				argument := arguments[index]
				asGoValue, asGoValueError := plasma.FromValue(argument)
				if asGoValueError != nil {
					return nil, asGoValueError
				}
				argumentAsReflectValue, convertErr := reflectConvert(reflect.ValueOf(asGoValue), variadicArgumentType)
				if convertErr != nil {
					return nil, convertErr
				}
				variadicArguments = append(variadicArguments, argumentAsReflectValue)
			}
			callArguments = append(callArguments, variadicArguments...)
		}
		return plasma.callGoFunc(symbols, function, callArguments...)
	}
}

func (plasma *Plasma) convertMethod(self reflect.Value, function reflect.Method) func(...any) any {
	return func(arguments ...any) any {
		in := make([]reflect.Value, 0, len(arguments)+1)
		in = append(in, self)
		for _, argument := range arguments {
			in = append(in, reflect.ValueOf(argument))
		}
		out := function.Func.Call(in)
		result := make([]any, 0, len(out))
		for _, o := range out {
			result = append(result, o.Interface())
		}
		switch len(result) {
		case 0:
			return nil
		case 1:
			return result[0]
		default:
			return result
		}
	}
}

func (plasma *Plasma) methodsToValue(symbols *Symbols, asReflectValue reflect.Value) (map[string]*Value, error) {
	asReflectValueType := asReflectValue.Type()
	numMethod := asReflectValueType.NumMethod()
	methods := make(map[string]*Value, numMethod)
	for index := 0; index < numMethod; index++ {
		method := asReflectValueType.Method(index)
		plasmaMethod, convertErr := plasma.ToValue(
			symbols,
			plasma.convertMethod(asReflectValue, method),
		)
		if convertErr != nil {
			return nil, convertErr
		}
		methods[method.Name] = plasmaMethod
	}
	return methods, nil
}

func (plasma *Plasma) ToValue(symbols *Symbols, v any) (*Value, error) {
	if v == nil {
		return plasma.None(), nil
	}
	var (
		obj     *Value
		methods map[string]*Value
	)
	if symbols == nil {
		symbols = plasma.Symbols()
	}
	asReflectValue := reflect.ValueOf(v)
	asReflectValueType := asReflectValue.Type()
	if asReflectValueType.NumMethod() > 0 {
		var transformError error
		methods, transformError = plasma.methodsToValue(symbols, asReflectValue)
		if transformError != nil {
			return nil, transformError
		}
	}
	switch asReflectValue.Kind() {
	case reflect.String:
		obj = plasma.NewString([]byte(asReflectValue.String()))
	case reflect.Bool:
		obj = plasma.NewBool(asReflectValue.Bool())
	case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		obj = plasma.NewInt(int64(asReflectValue.Uint()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		obj = plasma.NewInt(asReflectValue.Int())
	case reflect.Float32, reflect.Float64:
		obj = plasma.NewFloat(asReflectValue.Float())
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
		obj = plasma.NewArray(values)
	case reflect.Map:
		keys := asReflectValue.MapKeys()
		hash := plasma.NewInternalHash()
		hash.internalMap = make(map[string]HashKeyValue, len(keys))
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
		obj = plasma.NewHash(hash)
	case reflect.Struct:
		nFields := asReflectValue.NumField()
		obj = plasma.NewValue(symbols, ValueId, plasma.Value())
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
	case reflect.Func:
		obj = plasma.NewBuiltInFunction(symbols, plasma.toValueGoFunctionCall(symbols, asReflectValue))
	case reflect.Pointer:
		return plasma.ToValue(symbols, asReflectValue.Elem().Interface())
	case reflect.Chan:
		obj = plasma.NewValue(symbols, ValueId, plasma.Value())
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
	if methods == nil {
		return obj, nil
	}
	for methodName, method := range methods {
		obj.Set(methodName, method)
	}
	return obj, nil
}
