package vm

import (
	"github.com/shoriwe/gruby/pkg/errors"
)

const (
	Addition      = "Addition"
	RightAddition = "Addition"
)

type Object interface {
	Initialize() (Object, *errors.Error)
	InitializeSubClass() (Object, *errors.Error)

	// Iterator
	Iterator() (Iterable, *errors.Error)

	// Unary
	AbsoluteValue() (Object, *errors.Error)
	NegateBits() (Object, *errors.Error)     // Negate Bits
	Negation(Object) (Object, *errors.Error) // Boolean Negation
	// Binary
	Addition(Object) (Object, *errors.Error)
	RightAddition(Object) (Object, *errors.Error)
	Subtraction(Object) (Object, *errors.Error)
	RightSubtraction(Object) (Object, *errors.Error)
	Modulus(Object) (Object, *errors.Error)
	RightModulus(Object) (Object, *errors.Error)
	Multiplication(Object) (Object, *errors.Error)
	RightMultiplication(Object) (Object, *errors.Error)
	Division(Object) (Object, *errors.Error)
	RightDivision(Object) (Object, *errors.Error)
	PowerOf(Object) (Object, *errors.Error)
	RightPowerOf(Object) (Object, *errors.Error)
	FloorDivision(Object) (Object, *errors.Error)
	RightFloorDivision(Object) (Object, *errors.Error)
	BitwiseRight(Object) (Object, *errors.Error)
	RightBitwiseRight(Object) (Object, *errors.Error)
	BitwiseLeft(Object) (Object, *errors.Error)
	RightBitwiseLeft(Object) (Object, *errors.Error)
	BitwiseAnd(Object) (Object, *errors.Error)
	RightBitwiseAnd(Object) (Object, *errors.Error)
	BitwiseOr(Object) (Object, *errors.Error)
	RightBitwiseOr(Object) (Object, *errors.Error)
	BitwiseXor(Object) (Object, *errors.Error)
	RightBitwiseXor(Object) (Object, *errors.Error)
	// Logical Binary
	And(Object) (Object, *errors.Error)
	RightAnd(Object) (Object, *errors.Error)
	Or(Object) (Object, *errors.Error)
	RightOr(Object) (Object, *errors.Error)
	Xor(Object) (Object, *errors.Error)
	RightXor(Object) (Object, *errors.Error)
	// Special Behavior
	Call(...Object) (Object, *errors.Error)
	Index(Object) (Object, *errors.Error)
	Delete(Object) *errors.Error
	// Comparison Binary Operations
	Equals(Object) (Object, *errors.Error)
	GreaterThan(Object) (Object, *errors.Error)
	GreaterOrEqualThan(Object) (Object, *errors.Error)
	LessThan(Object) (Object, *errors.Error)
	LessOrEqualThan(Object) (Object, *errors.Error)
	NotEqual(Object) (Object, *errors.Error)

	// Type conversion
	Integer() (*Integer, *errors.Error)
	Float() (*Float, *errors.Error)
	String() (*String, *errors.Error)
	RawString() (string, *errors.Error)
	Boolean() (Boolean, *errors.Error)

	New() (Object, *errors.Error)
	Dir() (*Hash, *errors.Error)
	GetAttribute(*String) *errors.Error
	SetAttribute(*String, Object) *errors.Error
	DelAttribute(*String) *errors.Error
	Hash() (Object, *errors.Error)

	Class() (Object, *errors.Error)
	SubClass() (*Array, *errors.Error)
	Documentation() (*Hash, *errors.Error)

	SymbolTable() *SymbolTable
}

func getObjectBuiltInMethod(object Object, symbolName string) interface{} {
	switch symbolName {
	case Addition:
		return object.Addition
	}
	return nil
}

func GetAttribute(object Object, symbolName string) (interface{}, *errors.Error) {
	var attribute interface{}
	var getError *errors.Error
	attribute, getError = object.SymbolTable().Get(symbolName)
	if getError != nil {
		attribute = getObjectBuiltInMethod(object, symbolName)
		if attribute == nil {
			return nil, getError
		}
	}
	return attribute, nil
}
