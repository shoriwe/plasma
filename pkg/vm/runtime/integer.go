package runtime

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/errors"
	"math"
	"reflect"
	"strconv"
)

type Integer struct {
	symbolTable *SymbolTable
	value       int64
}

func (integer *Integer) Initialize() (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) InitializeSubClass() (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) Iterator() (Iterable, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) AbsoluteValue() (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) NegateBits() (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) Negation(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) Addition(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return &Float{
			value: float64(integer.value) + right.(*Float).value,
		}, nil
	case *Integer:
		return &Integer{
			value: integer.value + right.(*Integer).value,
		}, nil
	default:
		operation, getError := GetAttribute(right, RightAdditionName, false)
		if getError != nil {
			return nil, getError
		}
		switch operation.(type) {
		case func(Object) (Object, *errors.Error):
			return operation.(func(Object) (Object, *errors.Error))(integer)
		case *Function:
			return operation.(*Function).Call(integer)
		default:
			return nil, NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
		}
	}
}

func (integer *Integer) RightAddition(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return &Float{
			value: left.(*Float).value + float64(integer.value),
		}, nil
	case *Integer:
		return &Integer{
			value: left.(*Integer).value + integer.value,
		}, nil
	default:
		return nil, NewMethodNotImplemented(RightAdditionName)
	}
}

func (integer *Integer) Subtraction(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return &Float{
			value: float64(integer.value) - right.(*Float).value,
		}, nil

	case *Integer:
		return &Integer{
			value: integer.value - right.(*Integer).value,
		}, nil
	default:
		operation, getError := GetAttribute(right, RightSubtractionName, false)
		if getError != nil {
			return nil, getError
		}
		switch operation.(type) {
		case func(Object) (Object, *errors.Error):
			return operation.(func(Object) (Object, *errors.Error))(integer)
		case *Function:
			return operation.(*Function).Call(integer)
		default:
			return nil, NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
		}
	}
}

func (integer *Integer) RightSubtraction(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return &Float{
			value: left.(*Float).value - float64(integer.value),
		}, nil

	case *Integer:
		return &Integer{
			value: left.(*Integer).value - integer.value,
		}, nil
	default:
		return nil, NewMethodNotImplemented(RightSubtractionName)
	}
}

func (integer *Integer) Modulus(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return &Float{
			value: math.Mod(float64(integer.value), right.(*Float).value),
		}, nil

	case *Integer:
		return &Integer{
			value: integer.value % right.(*Integer).value,
		}, nil
	default:
		operation, getError := GetAttribute(right, RightModulusName, false)
		if getError != nil {
			return nil, getError
		}
		switch operation.(type) {
		case func(Object) (Object, *errors.Error):
			return operation.(func(Object) (Object, *errors.Error))(integer)
		case *Function:
			return operation.(*Function).Call(integer)
		default:
			return nil, NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
		}
	}
}

func (integer *Integer) RightModulus(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return &Float{
			value: math.Mod(left.(*Float).value, float64(integer.value)),
		}, nil
	case *Integer:
		return &Integer{
			value: left.(*Integer).value % integer.value,
		}, nil
	default:
		return nil, NewMethodNotImplemented(RightModulusName)
	}
}

func (integer *Integer) Multiplication(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return &Float{
			value: float64(integer.value) * right.(*Float).value,
		}, nil
	case *Integer:
		return &Integer{
			value: integer.value * right.(*Integer).value,
		}, nil
	default:
		operation, getError := GetAttribute(right, RightMultiplicationName, false)
		if getError != nil {
			return nil, getError
		}
		switch operation.(type) {
		case func(Object) (Object, *errors.Error):
			return operation.(func(Object) (Object, *errors.Error))(integer)
		case *Function:
			return operation.(*Function).Call(integer)
		default:
			return nil, NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
		}
	}
}

func (integer *Integer) RightMultiplication(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return &Float{
			value: left.(*Float).value * float64(integer.value),
		}, nil

	case *Integer:
		return &Integer{
			value: left.(*Integer).value * integer.value,
		}, nil
	default:
		return nil, NewMethodNotImplemented(RightMultiplicationName)
	}
}

func (integer *Integer) Division(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return &Float{
			value: float64(integer.value) / right.(*Float).value,
		}, nil

	case *Integer:
		return &Float{
			value: float64(integer.value) / float64(right.(*Integer).value),
		}, nil
	default:
		operation, getError := GetAttribute(right, RightDivisionName, false)
		if getError != nil {
			return nil, getError
		}
		switch operation.(type) {
		case func(Object) (Object, *errors.Error):
			return operation.(func(Object) (Object, *errors.Error))(integer)
		case *Function:
			return operation.(*Function).Call(integer)
		default:
			return nil, NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
		}
	}
}

func (integer *Integer) RightDivision(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return &Float{
			value: left.(*Float).value / float64(integer.value),
		}, nil

	case *Integer:
		return &Float{
			value: left.(*Float).value / float64(integer.value),
		}, nil
	default:
		return nil, NewMethodNotImplemented(RightDivisionName)
	}
}

func (integer *Integer) PowerOf(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return &Float{
			value: math.Pow(float64(integer.value), right.(*Float).value),
		}, nil
	case *Integer:
		return &Integer{
			value: int64(math.Pow(float64(integer.value), float64(right.(*Integer).value))),
		}, nil
	default:
		operation, getError := GetAttribute(right, RightPowerOfName, false)
		if getError != nil {
			return nil, getError
		}
		switch operation.(type) {
		case func(Object) (Object, *errors.Error):
			return operation.(func(Object) (Object, *errors.Error))(integer)
		case *Function:
			return operation.(*Function).Call(integer)
		default:
			return nil, NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
		}
	}
}

func (integer *Integer) RightPowerOf(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return &Float{
			value: math.Pow(left.(*Float).value, float64(integer.value)),
		}, nil
	case *Integer:
		return &Integer{
			value: int64(math.Pow(float64(left.(*Integer).value), float64(integer.value))),
		}, nil
	default:
		return nil, NewMethodNotImplemented(RightPowerOfName)
	}
}

func (integer *Integer) FloorDivision(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return &Integer{
			value: int64(float64(integer.value) / right.(*Float).value),
		}, nil

	case *Integer:
		return &Integer{
			value: integer.value / right.(*Integer).value,
		}, nil
	default:
		operation, getError := GetAttribute(right, RightFloorDivisionName, false)
		if getError != nil {
			return nil, getError
		}
		switch operation.(type) {
		case func(Object) (Object, *errors.Error):
			return operation.(func(Object) (Object, *errors.Error))(integer)
		case *Function:
			return operation.(*Function).Call(integer)
		default:
			return nil, NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
		}
	}
}

func (integer *Integer) RightFloorDivision(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return &Integer{
			value: int64(left.(*Float).value / float64(integer.value)),
		}, nil
	case *Integer:
		return &Integer{
			value: left.(*Integer).value / integer.value,
		}, nil
	default:
		return nil, NewMethodNotImplemented(RightFloorDivisionName)
	}
}

func (integer *Integer) BitwiseRight(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Integer:
		return &Integer{
			value: integer.value >> right.(*Integer).value,
		}, nil
	default:
		operation, getError := GetAttribute(right, RightBitwiseRightName, false)
		if getError != nil {
			return nil, getError
		}
		switch operation.(type) {
		case func(Object) (Object, *errors.Error):
			return operation.(func(Object) (Object, *errors.Error))(integer)
		case *Function:
			return operation.(*Function).Call(integer)
		default:
			return nil, NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
		}
	}
}

func (integer *Integer) RightBitwiseRight(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Integer:
		return &Integer{
			value: left.(*Integer).value >> integer.value,
		}, nil
	default:
		return nil, NewMethodNotImplemented(RightBitwiseRightName)
	}
}

func (integer *Integer) BitwiseLeft(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Integer:
		return &Integer{
			value: integer.value << right.(*Integer).value,
		}, nil
	default:
		operation, getError := GetAttribute(right, RightBitwiseLeftName, false)
		if getError != nil {
			return nil, getError
		}
		switch operation.(type) {
		case func(Object) (Object, *errors.Error):
			return operation.(func(Object) (Object, *errors.Error))(integer)
		case *Function:
			return operation.(*Function).Call(integer)
		default:
			return nil, NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
		}
	}
}

func (integer *Integer) RightBitwiseLeft(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Integer:
		return &Integer{
			value: left.(*Integer).value << integer.value,
		}, nil
	default:
		return nil, NewMethodNotImplemented(RightBitwiseLeftName)
	}
}

func (integer *Integer) BitwiseAnd(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Integer:
		return &Integer{
			value: integer.value & right.(*Integer).value,
		}, nil
	default:
		operation, getError := GetAttribute(right, RightBitwiseAndName, false)
		if getError != nil {
			return nil, getError
		}
		switch operation.(type) {
		case func(Object) (Object, *errors.Error):
			return operation.(func(Object) (Object, *errors.Error))(integer)
		case *Function:
			return operation.(*Function).Call(integer)
		default:
			return nil, NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
		}
	}
}

func (integer *Integer) RightBitwiseAnd(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Integer:
		return &Integer{
			value: left.(*Integer).value & integer.value,
		}, nil
	default:
		return nil, NewMethodNotImplemented(RightBitwiseAndName)
	}
}

func (integer *Integer) BitwiseOr(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Integer:
		return &Integer{
			value: integer.value | right.(*Integer).value,
		}, nil
	default:
		operation, getError := GetAttribute(right, RightBitwiseOrName, false)
		if getError != nil {
			return nil, getError
		}
		switch operation.(type) {
		case func(Object) (Object, *errors.Error):
			return operation.(func(Object) (Object, *errors.Error))(integer)
		case *Function:
			return operation.(*Function).Call(integer)
		default:
			return nil, NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
		}
	}
}

func (integer *Integer) RightBitwiseOr(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Integer:
		return &Integer{
			value: left.(*Integer).value | integer.value,
		}, nil
	default:
		return nil, NewMethodNotImplemented(RightBitwiseOrName)
	}
}

func (integer *Integer) BitwiseXor(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Integer:
		return &Integer{
			value: integer.value ^ right.(*Integer).value,
		}, nil
	default:
		operation, getError := GetAttribute(right, RightBitwiseXorName, false)
		if getError != nil {
			return nil, getError
		}
		switch operation.(type) {
		case func(Object) (Object, *errors.Error):
			return operation.(func(Object) (Object, *errors.Error))(integer)
		case *Function:
			return operation.(*Function).Call(integer)
		default:
			return nil, NewTypeError(FunctionTypeString, reflect.TypeOf(operation).String())
		}
	}
}

func (integer *Integer) RightBitwiseXor(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Integer:
		return &Integer{
			value: left.(*Integer).value ^ integer.value,
		}, nil
	default:
		return nil, NewMethodNotImplemented(RightBitwiseXorName)
	}
}

func (integer *Integer) And(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) RightAnd(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) Or(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) RightOr(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) Xor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) RightXor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) Call(object ...Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) Index(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) Delete(object Object) *errors.Error {
	panic("implement me")
}

func (integer *Integer) Equals(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) GreaterThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) GreaterOrEqualThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) LessThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) LessOrEqualThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) NotEqual(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) Integer() (*Integer, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) Float() (*Float, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) String() (*String, *errors.Error) {
	return &String{
		value: fmt.Sprintf("%d", integer.value),
	}, nil
}

func (integer *Integer) Boolean() (Boolean, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) New() (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) Dir() (*Hash, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) GetAttribute(s *String) *errors.Error {
	panic("implement me")
}

func (integer *Integer) SetAttribute(s *String, object Object) *errors.Error {
	panic("implement me")
}

func (integer *Integer) DelAttribute(s *String) *errors.Error {
	panic("implement me")
}

func (integer *Integer) Hash() (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) Class() (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) SubClass() (*Array, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) Documentation() (*Hash, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) SymbolTable() *SymbolTable {
	return integer.symbolTable
}

func (integer *Integer) RawString() (string, *errors.Error) {
	return fmt.Sprintf("%d", integer.value), nil
}

func NewInteger(parentSymbolTable *SymbolTable, number string, base int) *Integer {
	value, _ := strconv.ParseInt(number, base, 64)
	return &Integer{
		symbolTable: NewSymbolTable(parentSymbolTable),
		value:       value,
	}
}
