package runtime

import "github.com/shoriwe/gruby/pkg/errors"

type None struct {
}

func (n None) Initialize() (Object, *errors.Error) {
	panic("implement me")
}

func (n None) InitializeSubClass() (Object, *errors.Error) {
	panic("implement me")
}

func (n None) Iterator() (Iterable, *errors.Error) {
	panic("implement me")
}

func (n None) AbsoluteValue() (Object, *errors.Error) {
	panic("implement me")
}

func (n None) NegateBits() (Object, *errors.Error) {
	panic("implement me")
}

func (n None) Negation() (Object, *errors.Error) {
	panic("implement me")
}

func (n None) Addition(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) RightAddition(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) Subtraction(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) RightSubtraction(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) Modulus(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) RightModulus(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) Multiplication(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) RightMultiplication(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) Division(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) RightDivision(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) PowerOf(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) RightPowerOf(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) FloorDivision(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) RightFloorDivision(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) BitwiseRight(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) RightBitwiseRight(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) BitwiseLeft(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) RightBitwiseLeft(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) BitwiseAnd(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) RightBitwiseAnd(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) BitwiseOr(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) RightBitwiseOr(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) BitwiseXor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) RightBitwiseXor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) And(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) RightAnd(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) Or(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) RightOr(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) Xor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) RightXor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) Call(object ...Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) Index(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) Delete(object Object) *errors.Error {
	panic("implement me")
}

func (n None) Equals(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) GreaterThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) GreaterOrEqualThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) LessThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) LessOrEqualThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) NotEqual(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (n None) Integer() (*Integer, *errors.Error) {
	panic("implement me")
}

func (n None) Float() (*Float, *errors.Error) {
	panic("implement me")
}

func (n None) String() (*String, *errors.Error) {
	panic("implement me")
}

func (n None) RawString() (string, *errors.Error) {
	panic("implement me")
}

func (n None) Boolean() (*Bool, *errors.Error) {
	panic("implement me")
}

func (n None) New() (Object, *errors.Error) {
	panic("implement me")
}

func (n None) Dir() (*Hash, *errors.Error) {
	panic("implement me")
}

func (n None) GetAttribute(s *String) *errors.Error {
	panic("implement me")
}

func (n None) SetAttribute(s *String, object Object) *errors.Error {
	panic("implement me")
}

func (n None) DelAttribute(s *String) *errors.Error {
	panic("implement me")
}

func (n None) Hash() (Object, *errors.Error) {
	panic("implement me")
}

func (n None) Class() (Object, *errors.Error) {
	panic("implement me")
}

func (n None) SubClass() (*Array, *errors.Error) {
	panic("implement me")
}

func (n None) Documentation() (*Hash, *errors.Error) {
	panic("implement me")
}

func (n None) SymbolTable() *SymbolTable {
	panic("implement me")
}

func NewNone() *None {
	return &None{}
}
