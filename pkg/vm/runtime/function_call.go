package runtime

import (
	"github.com/shoriwe/gruby/pkg/errors"
	"reflect"
)

func NoArgumentsMethodCall(methodName string, object Object) (Object, *errors.Error) {
	method, getError := GetAttribute(object, methodName, false)
	if getError != nil {
		return nil, getError
	}
	switch method.(type) {
	case func() (*Bool, *errors.Error): // Boolean
		return method.(func() (*Bool, *errors.Error))()
	case func() (*String, *errors.Error): // String
		return method.(func() (*String, *errors.Error))()
	case func() (*Integer, *errors.Error): // Integer
		return method.(func() (*Integer, *errors.Error))()
	case func() (*Float, *errors.Error): // Float
		return method.(func() (*Float, *errors.Error))()
	case *Function:
		return method.(Function).Call()
	}
	return nil, NewTypeError(reflect.TypeOf(method).String(), FunctionTypeString)
}

func BasicBinaryOP(leftOP string, rightOP string, leftHandSide Object, rightHandSide Object) (Object, *errors.Error) {
	operation, getError := GetAttribute(leftHandSide, leftOP, false)
	var result Object
	var opError *errors.Error
	isRight := false
	if getError != nil {
		var getError2 *errors.Error
		operation, getError2 = GetAttribute(rightHandSide, rightOP, false)
		if getError2 != nil {
			return nil, getError
		}
		isRight = true
	}
	switch operation.(type) {
	case func(Object) (Object, *errors.Error):
		if isRight {
			result, opError = operation.(func(Object) (Object, *errors.Error))(leftHandSide.(Object))
		} else {
			result, opError = operation.(func(Object) (Object, *errors.Error))(rightHandSide.(Object))
		}
	case *Function:
		if isRight {
			result, opError = operation.(*Function).Call(rightHandSide.(Object))
		} else {
			result, opError = operation.(*Function).Call(leftHandSide.(Object))
		}
	default:
		return nil, NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
	}
	if opError != nil {
		return nil, opError
	}
	return result, nil
}
