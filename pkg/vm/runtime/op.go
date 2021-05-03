package runtime

import (
	"github.com/shoriwe/gruby/pkg/errors"
)

type InstructionOP func(*Stack) *errors.Error

const (
	AdditionName            = "Addition"
	RightAdditionName       = "RightAddition"
	SubtractionName         = "Subtraction"
	RightSubtractionName    = "RightSubtraction"
	DivisionName            = "Division"
	RightDivisionName       = "RightDivision"
	ModulusName             = "Modulus"
	RightModulusName        = "RightModulus"
	MultiplicationName      = "Multiplication"
	RightMultiplicationName = "RightMultiplication"
	PowerOfName             = "PowerOf"
	RightPowerOfName        = "RightPowerOf"
	FloorDivisionName       = "FloorDivision"
	RightFloorDivisionName  = "RightFloorDivision"
	BitwiseLeftName         = "BitwiseLeftOP"
	RightBitwiseLeftName    = "RightBitwiseLeft"
	BitwiseRightName        = "BitwiseRightOP"
	RightBitwiseRightName   = "RightBitwiseRight"
	BitwiseAndName          = "BitwiseAndOP"
	RightBitwiseAndName     = "RightBitwiseAnd"
	BitwiseOrName           = "BitwiseOrOP"
	RightBitwiseOrName      = "RightBitwiseOr"
	BitwiseXorName          = "BitwiseXorOP"
	RightBitwiseXorName     = "RightBitwiseXor"
	BooleanName             = "Boolean"
	AndName                 = "And"
	RightAndName            = "RightAnd"
)

const (
	// Binary Operations
	AddOP uint = iota
	SubOP
	DivOP
	MulOP
	PowOP
	ModOP
	FloorDivOP
	BitwiseLeftOP
	BitwiseRightOP
	BitwiseAndOP
	BitwiseOrOP
	BitwiseXorOP
	// Binary logical operations
	AndOP
	OrOP

	NegateBitsOP

	// Memory Operations
	PushOP

	// Behavior
	ReturnOP
)

func getObjectBuiltInMethod(object Object, symbolName string) interface{} {
	switch symbolName {
	case AdditionName:
		return object.Addition
	case RightAdditionName:
		return object.RightAddition
	case SubtractionName:
		return object.Subtraction
	case RightSubtractionName:
		return object.RightSubtraction
	case DivisionName:
		return object.Division
	case RightDivisionName:
		return object.RightDivision
	case ModulusName:
		return object.Modulus
	case RightModulusName:
		return object.RightModulus
	case MultiplicationName:
		return object.Multiplication
	case RightMultiplicationName:
		return object.RightMultiplication
	case PowerOfName:
		return object.PowerOf
	case RightPowerOfName:
		return object.RightPowerOf
	case FloorDivisionName:
		return object.FloorDivision
	case RightFloorDivisionName:
		return object.RightFloorDivision
	case BitwiseLeftName:
		return object.BitwiseLeft
	case RightBitwiseLeftName:
		return object.RightBitwiseLeft
	case BitwiseRightName:
		return object.BitwiseRight
	case RightBitwiseRightName:
		return object.RightBitwiseRight
	case BitwiseAndName:
		return object.BitwiseAnd
	case RightBitwiseAndName:
		return object.RightBitwiseAnd
	case BitwiseOrName:
		return object.BitwiseOr
	case RightBitwiseOrName:
		return object.RightBitwiseOr
	case BitwiseXorName:
		return object.BitwiseXor
	case RightBitwiseXorName:
		return object.RightBitwiseXor
	case BooleanName:
		return object.Boolean
	case AndName:
		return object.And
	}
	return nil
}

func GetAttribute(object Object, symbolName string, useParent bool) (interface{}, *errors.Error) {
	var attribute interface{}
	var getError *errors.Error
	attribute, getError = object.SymbolTable().Get(symbolName, useParent)
	if getError != nil {
		attribute = getObjectBuiltInMethod(object, symbolName)
		if attribute == nil {
			return nil, getError
		}
	}
	return attribute, nil
}
