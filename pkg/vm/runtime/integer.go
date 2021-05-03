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

func (integer *Integer) Negation() (Object, *errors.Error) {
	rawIntegerAsBoolean, callError := NoArgumentsMethodCall(BooleanName, integer)
	if callError != nil {
		return nil, callError
	}
	if _, ok := rawIntegerAsBoolean.(*Bool); !ok {
		return nil, NewTypeError(reflect.TypeOf(rawIntegerAsBoolean).String(), BoolName)
	}
	rawIntegerAsBoolean.(*Bool).value = !rawIntegerAsBoolean.(*Bool).value
	return rawIntegerAsBoolean, nil
}

func (integer *Integer) Addition(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return NewFloatFromNumber(integer.symbolTable.parent, float64(integer.value)+right.(*Float).value)
	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, integer.value+right.(*Integer).value)
	default:
		return nil, NewTypeError(reflect.TypeOf(right).String(), IntegerName, FloatName)
	}
}

func (integer *Integer) RightAddition(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return NewFloatFromNumber(integer.symbolTable.parent, left.(*Float).value+float64(integer.value))
	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, left.(*Integer).value+integer.value)
	default:
		return nil, NewTypeError(reflect.TypeOf(left).String(), IntegerName, FloatName)
	}
}

func (integer *Integer) Subtraction(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return NewFloatFromNumber(integer.symbolTable.parent, float64(integer.value)-right.(*Float).value)

	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, integer.value-right.(*Integer).value)
	default:
		return nil, NewTypeError(reflect.TypeOf(right).String(), IntegerName, FloatName)
	}
}

func (integer *Integer) RightSubtraction(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return NewFloatFromNumber(integer.symbolTable.parent, left.(*Float).value-float64(integer.value))

	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, left.(*Integer).value-integer.value)
	default:
		return nil, NewTypeError(reflect.TypeOf(left).String(), IntegerName, FloatName)
	}
}

func (integer *Integer) Modulus(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return NewFloatFromNumber(integer.symbolTable.parent, math.Mod(float64(integer.value), right.(*Float).value))

	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, integer.value%right.(*Integer).value)
	default:
		return nil, NewTypeError(reflect.TypeOf(right).String(), IntegerName, FloatName)
	}
}

func (integer *Integer) RightModulus(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return NewFloatFromNumber(integer.symbolTable.parent, math.Mod(left.(*Float).value, float64(integer.value)))
	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, left.(*Integer).value%integer.value)
	default:
		return nil, NewTypeError(reflect.TypeOf(left).String(), IntegerName, FloatName)
	}
}

func (integer *Integer) Multiplication(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return NewFloatFromNumber(integer.symbolTable.parent, float64(integer.value)*right.(*Float).value)
	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, integer.value*right.(*Integer).value)
	case *String:
		panic("")
	default:
		return nil, NewTypeError(reflect.TypeOf(right).String(), IntegerName, FloatName, StringName, TupleName, ArrayName)
	}
}

func (integer *Integer) RightMultiplication(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return NewFloatFromNumber(integer.symbolTable.parent, left.(*Float).value*float64(integer.value))

	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, left.(*Integer).value*integer.value)
	default:
		return nil, NewTypeError(reflect.TypeOf(left).String(), IntegerName, FloatName, StringName, TupleName, ArrayName)
	}
}

func (integer *Integer) Division(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return NewFloatFromNumber(integer.symbolTable.parent, float64(integer.value)/right.(*Float).value)

	case *Integer:
		return NewFloatFromNumber(integer.symbolTable.parent, float64(integer.value)/float64(right.(*Integer).value))
	default:
		return nil, NewTypeError(reflect.TypeOf(right).String(), IntegerName, FloatName)
	}
}

func (integer *Integer) RightDivision(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return NewFloatFromNumber(integer.symbolTable.parent, left.(*Float).value/float64(integer.value))

	case *Integer:
		return NewFloatFromNumber(integer.symbolTable.parent, left.(*Float).value/float64(integer.value))
	default:
		return nil, NewTypeError(reflect.TypeOf(left).String(), IntegerName, FloatName)
	}
}

func (integer *Integer) PowerOf(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return NewFloatFromNumber(integer.symbolTable.parent, math.Pow(float64(integer.value), right.(*Float).value))
	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, int64(math.Pow(float64(integer.value), float64(right.(*Integer).value))))
	default:
		return nil, NewTypeError(reflect.TypeOf(right).String(), IntegerName, FloatName)
	}
}

func (integer *Integer) RightPowerOf(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return NewFloatFromNumber(integer.symbolTable.parent, math.Pow(left.(*Float).value, float64(integer.value)))
	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, int64(math.Pow(float64(left.(*Integer).value), float64(integer.value))))
	default:
		return nil, NewTypeError(reflect.TypeOf(left).String(), IntegerName, FloatName)
	}
}

func (integer *Integer) FloorDivision(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Float:
		return NewIntegerFromNumber(integer.symbolTable.parent, int64(float64(integer.value)/right.(*Float).value))

	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, integer.value/right.(*Integer).value)
	default:
		return nil, NewTypeError(reflect.TypeOf(right).String(), IntegerName, FloatName)
	}
}

func (integer *Integer) RightFloorDivision(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Float:
		return NewIntegerFromNumber(integer.symbolTable.parent, int64(left.(*Float).value/float64(integer.value)))
	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, left.(*Integer).value/integer.value)
	default:
		return nil, NewTypeError(reflect.TypeOf(left).String(), IntegerName, FloatName)
	}
}

func (integer *Integer) BitwiseRight(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, integer.value>>right.(*Integer).value)
	default:
		return nil, NewTypeError(reflect.TypeOf(right).String(), IntegerName)
	}
}

func (integer *Integer) RightBitwiseRight(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, left.(*Integer).value>>integer.value)
	default:
		return nil, NewTypeError(reflect.TypeOf(left).String(), IntegerName)
	}
}

func (integer *Integer) BitwiseLeft(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, integer.value<<right.(*Integer).value)
	default:
		return nil, NewTypeError(reflect.TypeOf(right).String(), IntegerName)
	}
}

func (integer *Integer) RightBitwiseLeft(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, left.(*Integer).value<<integer.value)
	default:
		return nil, NewTypeError(reflect.TypeOf(left).String(), IntegerName)
	}
}

func (integer *Integer) BitwiseAnd(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, integer.value&right.(*Integer).value)
	default:
		return nil, NewTypeError(reflect.TypeOf(right).String(), IntegerName)
	}
}

func (integer *Integer) RightBitwiseAnd(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, left.(*Integer).value&integer.value)
	default:
		return nil, NewTypeError(reflect.TypeOf(left).String(), IntegerName)
	}
}

func (integer *Integer) BitwiseOr(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, integer.value|right.(*Integer).value)
	default:
		return nil, NewTypeError(reflect.TypeOf(right).String(), IntegerName)
	}
}

func (integer *Integer) RightBitwiseOr(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, left.(*Integer).value|integer.value)
	default:
		return nil, NewTypeError(reflect.TypeOf(left).String(), IntegerName)
	}
}

func (integer *Integer) BitwiseXor(right Object) (Object, *errors.Error) {
	switch right.(type) {
	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, integer.value^right.(*Integer).value)
	default:
		return nil, NewTypeError(reflect.TypeOf(right).String(), IntegerName)
	}
}

func (integer *Integer) RightBitwiseXor(left Object) (Object, *errors.Error) {
	switch left.(type) {
	case *Integer:
		return NewIntegerFromNumber(integer.symbolTable.parent, left.(*Integer).value^integer.value)
	default:
		return nil, NewTypeError(reflect.TypeOf(left).String(), IntegerName)
	}
}

func (integer *Integer) And(right Object) (Object, *errors.Error) {
	rightBoolean, ok := right.(*Bool)
	if !ok {
		rawRightBoolean, callError := NoArgumentsMethodCall(BooleanName, right)
		if callError != nil {
			return nil, callError
		}
		if _, ok2 := rawRightBoolean.(*Bool); !ok2 {
			return nil, NewTypeError(reflect.TypeOf(rawRightBoolean).String(), BoolName)
		}
		rightBoolean = rawRightBoolean.(*Bool)
	} else {
		rightBoolean = right.(*Bool)
	}
	leftBoolean, _ := integer.Boolean()
	return NewBool(integer.symbolTable.parent, leftBoolean.value && rightBoolean.value), nil
}

func (integer *Integer) RightAnd(left Object) (Object, *errors.Error) {
	leftBoolean, ok := left.(*Bool)
	if !ok {
		rawLeftBoolean, callError := NoArgumentsMethodCall(BooleanName, left)
		if callError != nil {
			return nil, callError
		}
		if _, ok2 := rawLeftBoolean.(*Bool); !ok2 {
			return nil, NewTypeError(reflect.TypeOf(rawLeftBoolean).String(), BoolName)
		}
		leftBoolean = rawLeftBoolean.(*Bool)
	} else {
		leftBoolean = left.(*Bool)
	}
	rightBoolean, _ := integer.Boolean()
	return NewBool(integer.symbolTable.parent, leftBoolean.value && rightBoolean.value), nil
}

func (integer *Integer) Or(right Object) (Object, *errors.Error) {
	rightBoolean, ok := right.(*Bool)
	if !ok {
		rawRightBoolean, callError := NoArgumentsMethodCall(BooleanName, right)
		if callError != nil {
			return nil, callError
		}
		if _, ok2 := rawRightBoolean.(*Bool); !ok2 {
			return nil, NewTypeError(reflect.TypeOf(rawRightBoolean).String(), BoolName)
		}
		rightBoolean = rawRightBoolean.(*Bool)
	} else {
		rightBoolean = right.(*Bool)
	}
	leftBoolean, _ := integer.Boolean()
	return NewBool(integer.symbolTable.parent, leftBoolean.value || rightBoolean.value), nil
}

func (integer *Integer) RightOr(left Object) (Object, *errors.Error) {
	leftBoolean, ok := left.(*Bool)
	if !ok {
		rawLeftBoolean, callError := NoArgumentsMethodCall(BooleanName, left)
		if callError != nil {
			return nil, callError
		}
		if _, ok2 := rawLeftBoolean.(*Bool); !ok2 {
			return nil, NewTypeError(reflect.TypeOf(rawLeftBoolean).String(), BoolName)
		}
		leftBoolean = rawLeftBoolean.(*Bool)
	} else {
		leftBoolean = left.(*Bool)
	}
	rightBoolean, _ := integer.Boolean()
	return NewBool(integer.symbolTable.parent, leftBoolean.value || rightBoolean.value), nil
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
	fmt.Println(integer.symbolTable)
	return NewString(integer.symbolTable.parent, fmt.Sprintf("'%d'", integer.value)), nil
}

func (integer *Integer) Boolean() (*Bool, *errors.Error) {
	return NewBool(NewSymbolTable(integer.symbolTable.parent), integer.value != 0), nil
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

func NewIntegerFromNumber(parentSymbolTable *SymbolTable, number interface{}) (*Integer, *errors.Error) {
	if _, ok := number.(int64); ok {
		return &Integer{
			symbolTable: NewSymbolTable(parentSymbolTable),
			value:       number.(int64),
		}, nil
	} else if _, ok2 := number.(float64); ok2 {
		return &Integer{
			symbolTable: NewSymbolTable(parentSymbolTable),
			value:       int64(number.(float64)),
		}, nil
	}
	return nil, NewTypeError(reflect.TypeOf(number).String(), "int64", "float64")
}

func NewIntegerFromString(parentSymbolTable *SymbolTable, number string, base int) *Integer {
	value, _ := strconv.ParseInt(number, base, 64)
	return &Integer{
		symbolTable: NewSymbolTable(parentSymbolTable),
		value:       value,
	}
}
