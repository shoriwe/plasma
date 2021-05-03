package runtime

import (
	"github.com/shoriwe/gruby/pkg/errors"
)

type Bool struct {
	symbolTable *SymbolTable
	value       bool
}

func (b Bool) Initialize() (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) InitializeSubClass() (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) Iterator() (Iterable, *errors.Error) {
	panic("implement me")
}

func (b Bool) AbsoluteValue() (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) NegateBits() (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) Negation(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) Addition(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) RightAddition(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) Subtraction(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) RightSubtraction(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) Modulus(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) RightModulus(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) Multiplication(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) RightMultiplication(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) Division(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) RightDivision(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) PowerOf(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) RightPowerOf(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) FloorDivision(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) RightFloorDivision(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) BitwiseRight(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) RightBitwiseRight(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) BitwiseLeft(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) RightBitwiseLeft(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) BitwiseAnd(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) RightBitwiseAnd(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) BitwiseOr(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) RightBitwiseOr(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) BitwiseXor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) RightBitwiseXor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) And(right Object) (Object, *errors.Error) {
	rightBoolean, ok := right.(*Bool)
	if !ok {
		var transformationError *errors.Error
		rightBoolean, transformationError = right.Boolean()
		if transformationError != nil {
			return nil, transformationError
		}
	} else {
		rightBoolean = right.(*Bool)
	}
	return &Bool{
		value: b.value && rightBoolean.value,
	}, nil
}

func (b Bool) RightAnd(left Object) (Object, *errors.Error) {
	leftBoolean, ok := left.(*Bool)
	if !ok {
		var transformationError *errors.Error
		leftBoolean, transformationError = left.Boolean()
		if transformationError != nil {
			return nil, transformationError
		}
	} else {
		leftBoolean = left.(*Bool)
	}
	return &Bool{
		value: leftBoolean.value && b.value,
	}, nil
}

func (b Bool) Or(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) RightOr(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) Xor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) RightXor(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) Call(object ...Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) Index(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) Delete(object Object) *errors.Error {
	panic("implement me")
}

func (b Bool) Equals(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) GreaterThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) GreaterOrEqualThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) LessThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) LessOrEqualThan(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) NotEqual(object Object) (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) Integer() (*Integer, *errors.Error) {
	panic("implement me")
}

func (b Bool) Float() (*Float, *errors.Error) {
	panic("implement me")
}

func (b Bool) String() (*String, *errors.Error) {
	if b.value {
		return &String{
			value: "True",
		}, nil
	}
	return &String{
		value: "False",
	}, nil
}

func (b Bool) RawString() (string, *errors.Error) {
	if b.value {
		return "True", nil
	}
	return "False", nil
}

func (b Bool) Boolean() (*Bool, *errors.Error) {
	return &Bool{value: b.value}, nil
}

func (b Bool) New() (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) Dir() (*Hash, *errors.Error) {
	panic("implement me")
}

func (b Bool) GetAttribute(s *String) *errors.Error {
	panic("implement me")
}

func (b Bool) SetAttribute(s *String, object Object) *errors.Error {
	panic("implement me")
}

func (b Bool) DelAttribute(s *String) *errors.Error {
	panic("implement me")
}

func (b Bool) Hash() (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) Class() (Object, *errors.Error) {
	panic("implement me")
}

func (b Bool) SubClass() (*Array, *errors.Error) {
	panic("implement me")
}

func (b Bool) Documentation() (*Hash, *errors.Error) {
	panic("implement me")
}

func (b Bool) SymbolTable() *SymbolTable {
	return b.symbolTable
}

func NewBool(parentSymbolTable *SymbolTable, value bool) *Bool {
	return &Bool{
		symbolTable: NewSymbolTable(parentSymbolTable),
		value:       value,
	}
}

func NewTrue(parentSymbolTable *SymbolTable) *Bool {
	return NewBool(parentSymbolTable, true)
}

func NewFalse(parentSymbolTable *SymbolTable) *Bool {
	return NewBool(parentSymbolTable, false)
}
