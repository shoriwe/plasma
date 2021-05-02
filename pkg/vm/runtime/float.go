package runtime

import (
	"github.com/shoriwe/gruby/pkg/errors"
	"math"
	"math/big"
	"reflect"
)

type Float struct {
	symbolTable *SymbolTable
	value       *big.Float
	accuracy    big.Accuracy
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
		result := big.NewFloat(0)
		result.Add(result, float.value)
		result.Add(result, right.(*Float).value)
		return &Float{
			value: result,
		}, nil

	case *Integer:
		result := big.NewFloat(0)
		result.SetInt(right.(*Integer).value)
		result.Add(result, float.value)
		return &Float{
			value: result,
		}, nil
	default:
		operation, getError := GetAttribute(right, RightAddition, false)
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
		result := big.NewFloat(0)
		result.Add(result, float.value)
		result.Add(result, left.(*Float).value)
		return &Float{
			value: result,
		}, nil

	case *Integer:
		result := big.NewFloat(0)
		result.SetInt(left.(*Integer).value)
		result.Add(result, float.value)
		return &Float{
			value: result,
		}, nil
	default:
		return nil, NewMethodNotImplemented(RightAddition)
	}
}

func (float *Float) Subtraction(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		result := big.NewFloat(0)
		result.Set(float.value)
		result.Sub(result, right.(*Float).value)
		return &Float{
			value: result,
		}, nil

	case *Integer:
		result := big.NewFloat(0)
		result.Add(result, float.value)
		result.Sub(result, new(big.Float).SetInt(right.(*Integer).value))
		return &Float{
			value: result,
		}, nil
	default:
		operation, getError := GetAttribute(right, RightAddition, false)
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
		result := big.NewFloat(0)
		result.Set(float.value)
		result.Mul(result, big.NewFloat(-1))
		result.Add(result, left.(*Float).value)
		return &Float{
			value: result,
		}, nil

	case *Integer:
		result := big.NewFloat(0)
		result.Add(result, new(big.Float).SetInt(left.(*Integer).value))
		result.Sub(result, float.value)
		return &Float{
			value: result,
		}, nil
	default:
		return nil, NewMethodNotImplemented(RightAddition)
	}
}

func (float *Float) Modulus(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		leftFloat, accuracy := float.value.Float64()
		rightFloat, _ := right.(*Float).value.Float64()
		result := big.NewFloat(math.Mod(leftFloat, rightFloat))
		return &Float{
			accuracy: accuracy,
			value:    result,
		}, nil
	case *Integer:
		leftFloat, accuracy := float.value.Float64()
		result := big.NewFloat(math.Mod(leftFloat, float64(right.(*Integer).value.Int64())))
		return &Float{
			accuracy: accuracy,
			value:    result,
		}, nil
	default:
		operation, getError := GetAttribute(right, RightAddition, false)
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
		rightFloat, accuracy := float.value.Float64()
		leftFloat, _ := left.(*Float).value.Float64()
		result := big.NewFloat(math.Mod(leftFloat, rightFloat))
		return &Float{
			accuracy: accuracy,
			value:    result,
		}, nil
	case *Integer:
		rightFloat, accuracy := float.value.Float64()
		result := big.NewFloat(math.Mod(float64(left.(*Integer).value.Int64()), rightFloat))
		return &Float{
			accuracy: accuracy,
			value:    result,
		}, nil
	default:
		return nil, NewMethodNotImplemented(RightAddition)
	}
}

func (float *Float) Multiplication(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		result := big.NewFloat(0)
		result.Set(float.value)
		result.Mul(result, right.(*Float).value)
		return &Float{
			value: result,
		}, nil
	case *Integer:
		result := big.NewFloat(0)
		result.Set(float.value)
		result.Mul(result, new(big.Float).SetInt(right.(*Integer).value))
		return &Float{
			value: result,
		}, nil
	default:
		operation, getError := GetAttribute(right, RightAddition, false)
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
		result := big.NewFloat(0)
		result.Add(result, left.(*Float).value)
		result.Mul(result, float.value)
		return &Float{
			value: result,
		}, nil

	case *Integer:
		result := big.NewFloat(0)
		result.SetInt(left.(*Integer).value)
		result.Mul(result, float.value)
		return &Float{
			value: result,
		}, nil
	default:
		return nil, NewMethodNotImplemented(RightAddition)
	}
}

func (float *Float) Division(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		result := big.NewFloat(0)
		result.Set(float.value)
		result.Quo(result, right.(*Float).value)
		return &Float{
			value: result,
		}, nil

	case *Integer:
		result := big.NewFloat(0)
		result.Set(float.value)
		result.Quo(result, new(big.Float).SetInt(right.(*Integer).value))
		return &Float{
			value: result,
		}, nil
	default:
		operation, getError := GetAttribute(right, RightAddition, false)
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
		result := big.NewFloat(0)
		result.Add(result, left.(*Float).value)
		result.Quo(result, float.value)
		return &Float{
			value: result,
		}, nil

	case *Integer:
		result := big.NewFloat(0)
		result.SetInt(left.(*Integer).value)
		result.Quo(result, float.value)
		return &Float{
			value: result,
		}, nil
	default:
		return nil, NewMethodNotImplemented(RightAddition)
	}
}

func (float *Float) PowerOf(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		leftFloat, accuracy := float.value.Float64()
		rightFloat, _ := right.(*Float).value.Float64()
		result := big.NewFloat(math.Pow(leftFloat, rightFloat))
		return &Float{
			accuracy: accuracy,
			value:    result,
		}, nil
	case *Integer:
		leftFloat, accuracy := float.value.Float64()
		result := big.NewFloat(math.Pow(leftFloat, float64(right.(*Integer).value.Int64())))
		return &Float{
			accuracy: accuracy,
			value:    result,
		}, nil
	default:
		operation, getError := GetAttribute(right, RightAddition, false)
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
		leftFloat, accuracy := left.(*Float).value.Float64()
		rightFloat, _ := float.value.Float64()
		result := big.NewFloat(math.Pow(leftFloat, rightFloat))
		return &Float{
			accuracy: accuracy,
			value:    result,
		}, nil
	case *Integer:
		rightFloat, accuracy := float.value.Float64()
		result := big.NewFloat(math.Pow(float64(left.(*Integer).value.Int64()), rightFloat))
		return &Float{
			accuracy: accuracy,
			value:    result,
		}, nil
	default:
		return nil, NewMethodNotImplemented(RightAddition)
	}
}

func (float *Float) FloorDivision(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (float *Float) RightFloorDivision(object Object) (Object, *errors.Error) {
	panic("implement me")
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
	panic("implement me")
}

func (float *Float) Boolean() (Boolean, *errors.Error) {
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
	return &String{
		value: float.value.String(),
	}, nil
}

func NewFloat(parentSymbolTable *SymbolTable, number string) *Float {
	value := big.NewFloat(0)
	value.SetString(number)
	return &Float{
		symbolTable: NewSymbolTable(parentSymbolTable),
		value:       value,
	}
}
