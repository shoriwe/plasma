package runtime

import "github.com/shoriwe/gruby/pkg/errors"

type Function struct {
}

func (f Function) Initialize() (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) InitializeSubClass() (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) Iterator() (Iterable, *errors.Error) {
	panic("implement me")
}

func (f Function) AbsoluteValue() (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) NegateBits() (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) Negation(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) Addition(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) RightAddition(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) Subtraction(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) RightSubtraction(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) Modulus(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) RightModulus(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) Multiplication(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) RightMultiplication(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) Division(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) RightDivision(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) PowerOf(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) RightPowerOf(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) FloorDivision(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) RightFloorDivision(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) BitwiseRight(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) RightBitwiseRight(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) BitwiseLeft(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) RightBitwiseLeft(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) BitwiseAnd(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) RightBitwiseAnd(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) BitwiseOr(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) RightBitwiseOr(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) BitwiseXor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) RightBitwiseXor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) And(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) RightAnd(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) Or(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) RightOr(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) Xor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) RightXor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) Call(object ...Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) Index(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) Delete(object Object) *errors.Error {
	panic("implement me")
}

func (f Function) Equals(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) GreaterThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) GreaterOrEqualThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) LessThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) LessOrEqualThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) NotEqual(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) Integer() (*Integer, *errors.Error) {
	panic("implement me")
}

func (f Function) Float() (*Float, *errors.Error) {
	panic("implement me")
}

func (f Function) String() (*String, *errors.Error) {
	panic("implement me")
}

func (f Function) RawString() (string, *errors.Error) {
	panic("implement me")
}

func (f Function) Boolean() (Boolean, *errors.Error) {
	panic("implement me")
}

func (f Function) New() (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) Dir() (*Hash, *errors.Error) {
	panic("implement me")
}

func (f Function) GetAttribute(s *String) *errors.Error {
	panic("implement me")
}

func (f Function) SetAttribute(s *String, object Object) *errors.Error {
	panic("implement me")
}

func (f Function) DelAttribute(s *String) *errors.Error {
	panic("implement me")
}

func (f Function) Hash() (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) Class() (Object, *errors.Error) {
	panic("implement me")
}

func (f Function) SubClass() (*Array, *errors.Error) {
	panic("implement me")
}

func (f Function) Documentation() (*Hash, *errors.Error) {
	panic("implement me")
}

func (f Function) SymbolTable() *SymbolTable {
	panic("implement me")
}
