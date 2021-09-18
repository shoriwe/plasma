package json

import (
	"encoding/json"
	"github.com/shoriwe/gplasma/pkg/std/features/importlib"
	"github.com/shoriwe/gplasma/pkg/vm"
	"reflect"
)

func interpretJSON(context *vm.Context, p *vm.Plasma, i interface{}) *vm.Value {
	switch i.(type) {
	case bool:
		return p.InterpretAsBool(i.(bool))
	case nil:
		return p.GetNone()
	case string:
		return p.NewString(context, false, i.(string))
	case float64:
		return p.NewFloat(context, false, i.(float64))
	case map[string]interface{}:
		result := p.NewHashTable(context, false)
		for key, value := range i.(map[string]interface{}) {
			p.HashIndexAssign(context, result, p.NewString(context, false, key), interpretJSON(context, p, value))
		}
		return result
	case []interface{}:
		var elements []*vm.Value
		for _, element := range i.([]interface{}) {
			elements = append(elements, interpretJSON(context, p, element))
		}
		return p.NewArray(context, false, elements)
	default:
		panic(reflect.TypeOf(i))
	}
}

func dumpToJSON(context *vm.Context, p *vm.Plasma, value *vm.Value) (interface{}, *vm.Value) {
	switch value.BuiltInTypeId {
	case vm.IntegerId:
		return value.Integer, nil
	case vm.StringId:
		return value.String, nil
	case vm.FloatId:
		return value.Float, nil
	case vm.ArrayId:
		elements := make([]interface{}, len(value.Content))
		for index, plasmaElement := range value.Content {
			element, dumpError := dumpToJSON(context, p, plasmaElement)
			if dumpError != nil {
				return nil, dumpError
			}
			elements[index] = element
		}
		return elements, nil
	case vm.HashTableId:
		elements := map[string]interface{}{}
		for _, keyValues := range value.KeyValues {
			for _, keyValue := range keyValues {
				key, keyDumpError := dumpToJSON(context, p, keyValue.Key)
				if keyDumpError != nil {
					return nil, keyDumpError
				}
				val, valueDumpError := dumpToJSON(context, p, keyValue.Value)
				if valueDumpError != nil {
					return nil, valueDumpError
				}
				if _, ok := key.(string); !ok {
					return nil, p.NewInvalidTypeError(context, reflect.TypeOf(key).Name(), vm.StringName)
				}
				elements[key.(string)] = val
			}
		}
		return elements, nil
	}
	return nil, p.NewInvalidTypeError(context, value.GetClass(p).Name, vm.StringName, vm.HashName, vm.ArrayName, vm.IntegerName, vm.FloatName)
}

func jsonLoader(context *vm.Context, p *vm.Plasma) *vm.Value {
	result := p.NewModule(context, true)
	result.SetOnDemandSymbol(
		"loads",
		func() *vm.Value {
			return p.NewFunction(context, true, result.SymbolTable(),
				vm.NewBuiltInFunction(1,
					func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
						var target []byte
						if arguments[0].IsTypeById(vm.StringId) {
							target = []byte(arguments[0].String)
						} else if arguments[0].IsTypeById(vm.BytesId) {
							target = arguments[0].Bytes
						} else {
							return p.NewInvalidTypeError(context, arguments[0].GetClass(p).Name, vm.StringName, vm.BytesName), false
						}
						var receiver interface{}
						marshalError := json.Unmarshal(target, &receiver)
						if marshalError != nil {
							return p.NewGoRuntimeError(context, marshalError), false
						}
						return interpretJSON(context, p, receiver), true
					},
				),
			)
		},
	)
	result.SetOnDemandSymbol(
		"dumps",
		func() *vm.Value {
			return p.NewFunction(context, true, result.SymbolTable(),
				vm.NewBuiltInFunction(1,
					func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
						interfaceDump, rawDumpError := dumpToJSON(context, p, arguments[0])
						if rawDumpError != nil {
							return rawDumpError, false
						}
						bytesDump, dumpError := json.Marshal(interfaceDump)
						if dumpError != nil {
							return p.NewGoRuntimeError(context, dumpError), false
						}
						return p.NewString(context, false, string(bytesDump)), true
					},
				),
			)
		},
	)
	return result
}

var JSON = importlib.ModuleInformation{
	Name:   "json",
	Loader: jsonLoader,
}
