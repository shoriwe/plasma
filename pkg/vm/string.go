package vm

import (
	"github.com/shoriwe/gruby/pkg/errors"
)

type String struct {
	value string
}

func (string_ *String) Initialize() (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) InitializeSubClass() (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) Iterator() (Iterable, *errors.Error) {
	panic("implement me")
}

func (string_ *String) AbsoluteValue() (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) NegateBits() (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) Negation(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) Addition(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) RightAddition(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) Subtraction(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) RightSubtraction(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) Modulus(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) RightModulus(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) Multiplication(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) RightMultiplication(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) Division(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) RightDivision(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) PowerOf(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) RightPowerOf(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) FloorDivision(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) RightFloorDivision(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) BitwiseRight(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) RightBitwiseRight(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) BitwiseLeft(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) RightBitwiseLeft(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) BitwiseAnd(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) RightBitwiseAnd(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) BitwiseOr(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) RightBitwiseOr(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) BitwiseXor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) RightBitwiseXor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) And(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) RightAnd(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) Or(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) RightOr(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) Xor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) RightXor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) Call(object ...Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) Index(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) Delete(object Object) *errors.Error {
	panic("implement me")
}

func (string_ *String) Equals(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) GreaterThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) GreaterOrEqualThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) LessThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) LessOrEqualThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) NotEqual(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) Integer() (*Integer, *errors.Error) {
	panic("implement me")
}

func (string_ *String) Float() (*Float, *errors.Error) {
	panic("implement me")
}

func (string_ *String) String() (*String, *errors.Error) {
	panic("implement me")
}

func (string_ *String) Boolean() (Boolean, *errors.Error) {
	panic("implement me")
}

func (string_ *String) New() (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) Dir() (*Hash, *errors.Error) {
	panic("implement me")
}

func (string_ *String) GetAttribute(s *String) *errors.Error {
	panic("implement me")
}

func (string_ *String) SetAttribute(s *String, object Object) *errors.Error {
	panic("implement me")
}

func (string_ *String) DelAttribute(s *String) *errors.Error {
	panic("implement me")
}

func (string_ *String) Hash() (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) Class() (Object, *errors.Error) {
	panic("implement me")
}

func (string_ *String) SubClass() (*Array, *errors.Error) {
	panic("implement me")
}

func (string_ *String) Documentation() (*Hash, *errors.Error) {
	panic("implement me")
}

func (string_ *String) SymbolTable() *SymbolTable {
	panic("implement me")
}

func (string_ *String) RawString() string {
	return string_.value
}
