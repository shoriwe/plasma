package runtime

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/errors"
	"math"
	"math/big"
	"reflect"
	"strconv"
)

type Float struct {
	symbolTable *SymbolTable
	value       float64
}

func (float *Float) Initialize() (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) InitializeSubClass() (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) Iterator() (Iterable, *errors.Error) {
	panic("implement me")
}

func (float *Float) AbsoluteValue() (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) NegateBits() (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) Negation(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) Addition(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return NewFloatFromNumber(float.symbolTable.parent, float.value+right.(*Float).value)
	case *Integer:
		return NewFloatFromNumber(float.symbolTable.parent, float.value+float64(right.(*Integer).value))
	default:
		operation, getError := GetAttribute(right, RightAdditionName, false)
		if getError != nil {
			return nil, getError
		}
		switch operation.(type) {
		case func(Object) (Object, *errors.Error):
			return operation.(func(Object) (Object, *errors.Error))(float)
		case *Function:
			return operation.(*Function).Call(float)
		default:
			return nil, NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
		}
	}
}

func (float *Float) RightAddition(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return NewFloatFromNumber(float.symbolTable.parent, left.(*Float).value+float.value)
	case *Integer:
		return NewFloatFromNumber(float.symbolTable.parent, float64(left.(*Integer).value)+float.value)
	default:
		return nil, NewMethodNotImplemented(RightAdditionName)
	}
}

func (float *Float) Subtraction(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return NewFloatFromNumber(float.symbolTable.parent, float.value-right.(*Float).value)
	case *Integer:
		return NewFloatFromNumber(float.symbolTable.parent, float.value-float64(right.(*Integer).value))
	default:
		operation, getError := GetAttribute(right, RightSubtractionName, false)
		if getError != nil {
			return nil, getError
		}
		switch operation.(type) {
		case func(Object) (Object, *errors.Error):
			return operation.(func(Object) (Object, *errors.Error))(float)
		case *Function:
			return operation.(*Function).Call(float)
		default:
			return nil, NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
		}
	}
}

func (float *Float) RightSubtraction(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return NewFloatFromNumber(float.symbolTable.parent, left.(*Float).value-float.value)
	case *Integer:
		return NewFloatFromNumber(float.symbolTable.parent, float64(left.(*Integer).value)-float.value)
	default:
		return nil, NewMethodNotImplemented(RightSubtractionName)
	}
}

func (float *Float) Modulus(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return NewFloatFromNumber(float.symbolTable.parent, math.Mod(float.value, right.(*Float).value))
	case *Integer:
		return NewFloatFromNumber(float.symbolTable.parent, math.Mod(float.value, float64(right.(*Integer).value)))
	default:
		operation, getError := GetAttribute(right, RightModulusName, false)
		if getError != nil {
			return nil, getError
		}
		switch operation.(type) {
		case func(Object) (Object, *errors.Error):
			return operation.(func(Object) (Object, *errors.Error))(float)
		case *Function:
			return operation.(*Function).Call(float)
		default:
			return nil, NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
		}
	}
}

func (float *Float) RightModulus(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return NewFloatFromNumber(float.symbolTable.parent, math.Mod(left.(*Float).value, float.value))
	case *Integer:
		return NewFloatFromNumber(float.symbolTable.parent, math.Mod(float64(left.(*Integer).value), float.value))
	default:
		return nil, NewMethodNotImplemented(RightModulusName)
	}
}

func (float *Float) Multiplication(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return NewFloatFromNumber(float.symbolTable.parent, float.value*right.(*Float).value)
	case *Integer:
		return NewFloatFromNumber(float.symbolTable.parent, float.value*float64(right.(*Integer).value))
	default:
		operation, getError := GetAttribute(right, RightMultiplicationName, false)
		if getError != nil {
			return nil, getError
		}
		switch operation.(type) {
		case func(Object) (Object, *errors.Error):
			return operation.(func(Object) (Object, *errors.Error))(float)
		case *Function:
			return operation.(*Function).Call(float)
		default:
			return nil, NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
		}
	}
}

func (float *Float) RightMultiplication(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return NewFloatFromNumber(float.symbolTable.parent, left.(*Float).value*float.value)
	case *Integer:
		return NewFloatFromNumber(float.symbolTable.parent, float64(left.(*Integer).value)*float.value)
	default:
		return nil, NewMethodNotImplemented(RightMultiplicationName)
	}
}

func (float *Float) Division(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return NewFloatFromNumber(float.symbolTable.parent, float.value/right.(*Float).value)
	case *Integer:
		return NewFloatFromNumber(float.symbolTable.parent, float.value/float64(right.(*Integer).value))
	default:
		operation, getError := GetAttribute(right, RightDivisionName, false)
		if getError != nil {
			return nil, getError
		}
		switch operation.(type) {
		case func(Object) (Object, *errors.Error):
			return operation.(func(Object) (Object, *errors.Error))(float)
		case *Function:
			return operation.(*Function).Call(float)
		default:
			return nil, NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
		}
	}
}

func (float *Float) RightDivision(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return NewFloatFromNumber(float.symbolTable.parent, left.(*Float).value/left.(*Float).value)
	case *Integer:
		return NewFloatFromNumber(float.symbolTable.parent, float64(left.(*Integer).value)/float.value)
	default:
		return nil, NewMethodNotImplemented(RightDivisionName)
	}
}

func (float *Float) PowerOf(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return NewFloatFromNumber(float.symbolTable.parent, math.Pow(float.value, right.(*Float).value))
	case *Integer:
		return NewFloatFromNumber(float.symbolTable.parent, math.Pow(float.value, float64(right.(*Integer).value)))
	default:
		operation, getError := GetAttribute(right, RightPowerOfName, false)
		if getError != nil {
			return nil, getError
		}
		switch operation.(type) {
		case func(Object) (Object, *errors.Error):
			return operation.(func(Object) (Object, *errors.Error))(float)
		case *Function:
			return operation.(*Function).Call(float)
		default:
			return nil, NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
		}
	}
}

func (float *Float) RightPowerOf(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return NewFloatFromNumber(float.symbolTable.parent, math.Pow(left.(*Float).value, float.value))
	case *Integer:
		return NewFloatFromNumber(float.symbolTable.parent, math.Pow(float64(left.(*Integer).value), float.value))
	default:
		return nil, NewMethodNotImplemented(RightPowerOfName)
	}
}

func (float *Float) FloorDivision(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return NewIntegerFromNumber(float.symbolTable.parent, int64(float.value/right.(*Float).value))
	case *Integer:
		return NewIntegerFromNumber(float.symbolTable.parent, int64(float.value/float64(right.(*Integer).value)))
	default:
		operation, getError := GetAttribute(right, RightFloorDivisionName, false)
		if getError != nil {
			return nil, getError
		}
		switch operation.(type) {
		case func(Object) (Object, *errors.Error):
			return operation.(func(Object) (Object, *errors.Error))(float)
		case *Function:
			return operation.(*Function).Call(float)
		default:
			return nil, NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
		}
	}
}

func (float *Float) RightFloorDivision(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return NewIntegerFromNumber(float.symbolTable.parent, int64(left.(*Float).value/float.value))
	case *Integer:
		return NewIntegerFromNumber(float.symbolTable.parent, int64(float64(left.(*Integer).value)/float.value))
	default:
		return nil, NewMethodNotImplemented(RightFloorDivisionName)
	}
}

func (float *Float) BitwiseRight(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) RightBitwiseRight(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) BitwiseLeft(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) RightBitwiseLeft(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) BitwiseAnd(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) RightBitwiseAnd(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) BitwiseOr(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) RightBitwiseOr(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) BitwiseXor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) RightBitwiseXor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) And(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) RightAnd(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) Or(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) RightOr(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) Xor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) RightXor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) Call(object ...Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) Index(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) Delete(object Object) *errors.Error {
	panic("implement me")
}

func (float *Float) Equals(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) GreaterThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) GreaterOrEqualThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) LessThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) LessOrEqualThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) NotEqual(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) Integer() (*Integer, *errors.Error) {
	panic("implement me")
}

func (float *Float) Float() (*Float, *errors.Error) {
	panic("implement me")
}

func (float *Float) RawString() (string, *errors.Error) {
	return fmt.Sprintf("%f", float.value), nil
}

func (float *Float) Boolean() (*Bool, *errors.Error) {
	panic("implement me")
}

func (float *Float) New() (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) Dir() (*Hash, *errors.Error) {
	panic("implement me")
}

func (float *Float) GetAttribute(v *String) *errors.Error {
	panic("implement me")
}

func (float *Float) SetAttribute(v *String, object Object) *errors.Error {
	panic("implement me")
}

func (float *Float) DelAttribute(v *String) *errors.Error {
	panic("implement me")
}

func (float *Float) Hash() (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) Class() (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) SubClass() (*Array, *errors.Error) {
	panic("implement me")
}

func (float *Float) Documentation() (*Hash, *errors.Error) {
	panic("implement me")
}

func (float *Float) SymbolTable() *SymbolTable {
	return float.symbolTable
}

func (float *Float) String() (*String, *errors.Error) {
	return NewString(float.symbolTable.parent, fmt.Sprintf("'%f'", float.value)), nil
}

func NewFloatFromNumber(parentSymbolTable *SymbolTable, number interface{}) (*Float, *errors.Error) {
	if _, ok := number.(int64); ok {
		return &Float{
			symbolTable: NewSymbolTable(parentSymbolTable),
			value:       float64(number.(int64)),
		}, nil
	} else if _, ok2 := number.(float64); ok2 {
		return &Float{
			symbolTable: NewSymbolTable(parentSymbolTable),
			value:       number.(float64),
		}, nil
	}
	return nil, NewTypeError(reflect.TypeOf(number).String(), "int64", "float64")
}

func NewFloatFromString(parentSymbolTable *SymbolTable, number string) *Float {
	value := big.NewFloat(0)
	value.SetString(number)
	float, _ := strconv.ParseFloat(number, 64)
	return &Float{
		symbolTable: NewSymbolTable(parentSymbolTable),
		value:       float,
	}
}
