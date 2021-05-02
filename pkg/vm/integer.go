package vm

import (
	"github.com/shoriwe/gruby/pkg/errors"
	"math/big"
	"reflect"
)

type Integer struct {
	symbolTable *SymbolTable
	value       *big.Int
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

func (integer *Integer) RightAddition(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) Subtraction(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) RightSubtraction(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) Modulus(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) RightModulus(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) Multiplication(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) RightMultiplication(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) Division(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) RightDivision(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) PowerOf(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) RightPowerOf(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) FloorDivision(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) RightFloorDivision(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) BitwiseRight(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) RightBitwiseRight(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) BitwiseLeft(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) RightBitwiseLeft(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) BitwiseAnd(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) RightBitwiseAnd(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) BitwiseOr(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) RightBitwiseOr(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) BitwiseXor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (integer *Integer) RightBitwiseXor(object Object) (Object, *errors.Error) {
	panic("implement me")
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
	panic("implement me")
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

func (integer *Integer) Addition(other Object) (Object, *errors.Error) {
	switch other.(type) {
	case *Float:
		result := big.NewFloat(0)
		result.SetInt(integer.value)
		result.Add(result, other.(*Float).value)
		return &Float{
			value: result,
		}, nil

	case *Integer:
		result := big.NewInt(0)
		result.Add(result, integer.value)
		result.Add(result, other.(*Integer).value)
		return &Integer{
			value: result,
		}, nil
	default:
		operation, getError := GetAttribute(other, RightAddition)
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

func (integer *Integer) RawString() (string, *errors.Error) {
	return integer.value.String(), nil
}

func NewInteger(parentSymbolTable *SymbolTable, number string, base int) *Integer {
	n := big.NewInt(0)
	n.SetString(number, base)
	return &Integer{
		symbolTable: NewSymbolTable(parentSymbolTable),
		value:       n,
	}
}
