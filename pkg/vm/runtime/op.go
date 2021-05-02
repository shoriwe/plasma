package runtime

import (
	"github.com/shoriwe/gruby/pkg/errors"
)

type InstructionOP func(*Stack) *errors.Error

const (
	Addition            = "Addition"
	RightAddition       = "RightAddition"
	Subtraction         = "Subtraction"
	RightSubtraction    = "RightSubtraction"
	Division            = "Division"
	RightDivision       = "RightDivision"
	Modulus             = "Modulus"
	RightModulus        = "RightModulus"
	Multiplication      = "Multiplication"
	RightMultiplication = "RightMultiplication"
	PowerOf             = "PowerOf"
	RightPowerOf        = "RightPowerOf"
	FloorDivision       = "FloorDivision"
	RightFloorDivision  = "RightFloorDivision"
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
	NegateBitsOP
	BitAndOP
	BitOrOP
	BitXorOP
	BitLeftOP
	BitRightOP

	// Memory Operations
	PushOP

	// Behavior
	ReturnOP
)

func getObjectBuiltInMethod(object Object, symbolName string) interface{} {
	switch symbolName {
	case Addition:
		return object.Addition
	case RightAddition:
		return object.RightAddition
	case Subtraction:
		return object.Subtraction
	case RightSubtraction:
		return object.RightSubtraction
	case Division:
		return object.Division
	case RightDivision:
		return object.RightDivision
	case Modulus:
		return object.Modulus
	case RightModulus:
		return object.RightModulus
	case Multiplication:
		return object.Multiplication
	case RightMultiplication:
		return object.RightMultiplication
	case PowerOf:
		return object.PowerOf
	case RightPowerOf:
		return object.RightPowerOf
	case FloorDivision:
		return object.FloorDivision
	case RightFloorDivision:
		return object.RightFloorDivision
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
