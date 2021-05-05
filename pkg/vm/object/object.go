package object

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/errors"
	"github.com/shoriwe/gruby/pkg/vm/utils"
	"sync"
)

const Empty = ""
const (
	// IObject Creation
	Initialize = "Initialize" // Executed just after New
	// Unary Operations
	NegBits = "NegBits"
	Negate  = "Negate"
	// Binary Operations
	//// Basic Binary
	Add      = "Add"
	RightAdd = "RightAdd"
	Sub      = "Sub"
	RightSub = "RightSub"
	Mul      = "Mul"
	RightMul = "RightMul"
	Div      = "Div"
	RightDiv = "RightDiv"
	Mod      = "Mod"
	RightMod = "RightMod"
	Pow      = "Pow"
	RightPow = "RightPow"
	//// Bitwise Binary
	BitXor        = "BitXor"
	RightBitXor   = "RightBitXor"
	BitAnd        = "BitAnd"
	RightBitAnd   = "RightBitAnd"
	BitOr         = "BitOr"
	RightBitOr    = "RightBitOr"
	BitLeft       = "BitLeft"
	RightBitLeft  = "RightBitLeft"
	BitRight      = "BitRight"
	RightBitRight = "RightBitRight"
	//// Logical Binary
	And      = "And"
	RightAnd = "RightAnd"
	Or       = "Or"
	RightOr  = "RightOr"
	Xor      = "Xor"
	RightXor = "RightXor"
	//// Comparison Binary
	Equals                  = "Equals"
	RightEquals             = "RightEquals"
	NotEquals               = "NotEquals"
	RightNotEquals          = "RightNotEquals"
	GreaterThan             = "GreaterThan"
	RightGreaterThan        = "RightGreaterThan"
	LessThan                = "LessThan"
	RightLessThan           = "RightLessThan"
	GreaterThanOrEqual      = "GreaterThanOrEqual"
	RightGreaterThanOrEqual = "RightGreaterThanOrEqual"
	LessThanOrEqual         = "LessThanOrEqual"
	RightLessThanOrEqual    = "RightLessThanOrEqual"
	// Behavior
	Assign     = "Assign"
	Copy       = "Copy"
	Dir        = "Dir"
	Index      = "Index"
	Call       = "Call"
	Iter       = "Iter"
	Class      = "Class"
	SubClasses = "SubClasses"
	// Transformation
	ToInteger = "ToInteger"
	ToFloat   = "ToFloat"
	ToString  = "ToString"
	ToBool    = "ToBool"
	ToArray   = "ToArray"
	ToTuple   = "ToTuple"
)

type ObjCounter struct {
	currentId uint
	mutex     *sync.Mutex
}

func (objCounter *ObjCounter) NextId() uint {
	objCounter.mutex.Lock()
	result := objCounter.currentId
	objCounter.currentId++
	objCounter.mutex.Unlock()
	return result
}

var counter = &ObjCounter{
	currentId: 1,
	mutex:     new(sync.Mutex),
}

type GoFunctionCallback func(...IObject) (IObject, *errors.Error)

type ConstructorFunction func(utils.VirtualMachine, *utils.SymbolTable, *utils.SymbolTable, []interface{}) (IObject, *errors.Error)

const TypeName = "Type"

// Meta-Class
type Type struct {
	*Object
	Constructor ConstructorFunction
}

func NewType(masterSymbols *utils.SymbolTable, parentSymbols *utils.SymbolTable, constructor ConstructorFunction) *Type {
	return &Type{
		Object: &Object{
			id:         counter.NextId(),
			typeName:   TypeName,
			subClasses: nil,
			symbols: &utils.SymbolTable{
				Master:  masterSymbols,
				Parent:  parentSymbols,
				Symbols: map[string]interface{}{},
			},
		},
		Constructor: constructor,
	}
}

const CallableName = "Callable"

type Callable interface {
	NumberOfArguments() int
	Call(...IObject) (IObject, *errors.Error) // This should return directly the object or the code of the function
}

type PlasmaFunctionType struct {
	numberOfArguments int
	Code              []interface{}
}

func (p *PlasmaFunctionType) NumberOfArguments() int {
	return p.numberOfArguments
}

func (p *PlasmaFunctionType) Call(arguments ...IObject) (IObject, *errors.Error) {
	if len(arguments) != p.numberOfArguments {
		return nil, errors.NewInvalidNumberOfArguments(len(arguments), p.numberOfArguments)
	}
	this := arguments[0]
	vmCopy := this.VirtualMachine().New()
	initializationError := vmCopy.Initialize(p.Code)
	if initializationError != nil {
		return nil, initializationError
	}
	result, executionError := vmCopy.Execute()
	if executionError != nil {
		return nil, executionError
	}
	return result.(IObject), nil
}

type GoFunctionType struct {
	numberOfArguments int
	callback          GoFunctionCallback
}

func (g *GoFunctionType) NumberOfArguments() int {
	return g.numberOfArguments
}

func (g *GoFunctionType) Call(arguments ...IObject) (IObject, *errors.Error) {
	if len(arguments) != g.numberOfArguments {
		return nil, errors.NewInvalidNumberOfArguments(len(arguments), g.numberOfArguments)
	}
	result, callError := g.callback(arguments...)
	if callError != nil {
		return nil, callError
	}
	return result, nil
}

func NewGoFunctionType(numberOfArguments int, callback GoFunctionCallback) *GoFunctionType {
	return &GoFunctionType{
		numberOfArguments: numberOfArguments,
		callback:          callback,
	}
}

func NewNotImplementedCallable(numberOfArguments int) *GoFunctionType {
	return NewGoFunctionType(numberOfArguments, func(_ ...IObject) (IObject, *errors.Error) {
		return nil, errors.NewNameNotFoundError()
	})
}

type IObject interface {
	Id() uint
	TypeName() string
	SymbolTable() *utils.SymbolTable
	SubClasses() []*Type
	Get(string) (interface{}, *errors.Error)
	Set(string, interface{})
	VirtualMachine() utils.VirtualMachine
}

// MetaClass for IObject
type Object struct {
	id             uint
	typeName       string
	subClasses     []*Type
	symbols        *utils.SymbolTable
	virtualMachine utils.VirtualMachine
}

func (o *Object) VirtualMachine() utils.VirtualMachine {
	return o.virtualMachine
}

func (o *Object) Id() uint {
	return o.id
}

func (o *Object) SubClasses() []*Type {
	return o.subClasses
}

func (o *Object) Get(symbol string) (interface{}, *errors.Error) {
	return o.symbols.GetSelf(symbol)
}

func (o *Object) Set(symbol string, object interface{}) {
	o.symbols.Set(symbol, object)
}

func (o *Object) TypeName() string {
	return o.typeName
}

func (o *Object) SymbolTable() *utils.SymbolTable {
	return o.symbols
}

func ObjInitialize(_ ...IObject) (IObject, *errors.Error) {
	return nil, nil
}

func ObjAnd(arguments ...IObject) (IObject, *errors.Error) {
	thisToBool, foundError := arguments[0].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var thisBool IObject
	var transformationError *errors.Error
	if _, ok := thisToBool.(*Function); !ok {
		if _, ok2 := thisToBool.(Callable); !ok2 {
			return nil, errors.NewTypeError([]string{FunctionName, CallableName}, thisToBool.(IObject).TypeName())
		} else {
			thisBool, transformationError = thisToBool.(Callable).Call(arguments[0])
		}
	} else {
		thisBool, transformationError = thisToBool.(Callable).Call(arguments[0])
	}
	if transformationError != nil {
		return nil, transformationError
	}

	var rightToBool interface{}
	rightToBool, foundError = arguments[1].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var rightBool IObject
	if _, ok := rightToBool.(*Function); !ok {
		if _, ok2 := rightToBool.(Callable); !ok2 {
			return nil, errors.NewTypeError([]string{FunctionName, CallableName}, rightToBool.(IObject).TypeName())
		} else {
			rightBool, transformationError = rightToBool.(Callable).Call(arguments[0])
		}
	} else {
		rightBool, transformationError = rightToBool.(Callable).Call(arguments[0])
	}
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(Empty, nil, thisBool.SymbolTable().Master, thisBool.SymbolTable().Parent, []interface{}{thisBool.(*Bool).Value && rightBool.(*Bool).Value}), nil
}

func ObjRightAnd(arguments ...IObject) (IObject, *errors.Error) {
	thisToBool, foundError := arguments[0].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var thisBool IObject
	var transformationError *errors.Error
	if _, ok := thisToBool.(*Function); !ok {
		if _, ok2 := thisToBool.(Callable); !ok2 {
			return nil, errors.NewTypeError([]string{FunctionName, CallableName}, thisToBool.(IObject).TypeName())
		} else {
			thisBool, transformationError = thisToBool.(Callable).Call(arguments[0])
		}
	} else {
		thisBool, transformationError = thisToBool.(Callable).Call(arguments[0])
	}
	if transformationError != nil {
		return nil, transformationError
	}

	var leftToBool interface{}
	leftToBool, foundError = arguments[1].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var leftBool IObject
	if _, ok := leftToBool.(*Function); !ok {
		if _, ok2 := leftToBool.(Callable); !ok2 {
			return nil, errors.NewTypeError([]string{FunctionName, CallableName}, leftToBool.(IObject).TypeName())
		} else {
			leftBool, transformationError = leftToBool.(Callable).Call(arguments[0])
		}
	} else {
		leftBool, transformationError = leftToBool.(Callable).Call(arguments[0])
	}
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(Empty, nil, thisBool.SymbolTable().Master, thisBool.SymbolTable().Parent, []interface{}{leftBool.(*Bool).Value && thisBool.(*Bool).Value}), nil
}

func ObjOr(arguments ...IObject) (IObject, *errors.Error) {
	thisToBool, foundError := arguments[0].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var thisBool IObject
	var transformationError *errors.Error
	if _, ok := thisToBool.(*Function); !ok {
		if _, ok2 := thisToBool.(Callable); !ok2 {
			return nil, errors.NewTypeError([]string{FunctionName, CallableName}, thisToBool.(IObject).TypeName())
		} else {
			thisBool, transformationError = thisToBool.(Callable).Call(arguments[0])
		}
	} else {
		thisBool, transformationError = thisToBool.(Callable).Call(arguments[0])
	}
	if transformationError != nil {
		return nil, transformationError
	}

	var rightToBool interface{}
	rightToBool, foundError = arguments[1].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var rightBool IObject
	if _, ok := rightToBool.(*Function); !ok {
		if _, ok2 := rightToBool.(Callable); !ok2 {
			return nil, errors.NewTypeError([]string{FunctionName, CallableName}, rightToBool.(IObject).TypeName())
		} else {
			rightBool, transformationError = rightToBool.(Callable).Call(arguments[0])
		}
	} else {
		rightBool, transformationError = rightToBool.(Callable).Call(arguments[0])
	}
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(Empty, nil, thisBool.SymbolTable().Master, thisBool.SymbolTable().Parent, []interface{}{thisBool.(*Bool).Value || rightBool.(*Bool).Value}), nil
}

func ObjRightOr(arguments ...IObject) (IObject, *errors.Error) {
	thisToBool, foundError := arguments[0].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var thisBool IObject
	var transformationError *errors.Error
	if _, ok := thisToBool.(*Function); !ok {
		if _, ok2 := thisToBool.(Callable); !ok2 {
			return nil, errors.NewTypeError([]string{FunctionName, CallableName}, thisToBool.(IObject).TypeName())
		} else {
			thisBool, transformationError = thisToBool.(Callable).Call(arguments[0])
		}
	} else {
		thisBool, transformationError = thisToBool.(Callable).Call(arguments[0])
	}
	if transformationError != nil {
		return nil, transformationError
	}

	var leftToBool interface{}
	leftToBool, foundError = arguments[1].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var leftBool IObject
	if _, ok := leftToBool.(*Function); !ok {
		if _, ok2 := leftToBool.(Callable); !ok2 {
			return nil, errors.NewTypeError([]string{FunctionName, CallableName}, leftToBool.(IObject).TypeName())
		} else {
			leftBool, transformationError = leftToBool.(Callable).Call(arguments[0])
		}
	} else {
		leftBool, transformationError = leftToBool.(Callable).Call(arguments[0])
	}
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(Empty, nil, thisBool.SymbolTable().Master, thisBool.SymbolTable().Parent, []interface{}{leftBool.(*Bool).Value || thisBool.(*Bool).Value}), nil
}

func ObjXor(arguments ...IObject) (IObject, *errors.Error) {
	thisToBool, foundError := arguments[0].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var thisBool IObject
	var transformationError *errors.Error
	if _, ok := thisToBool.(*Function); !ok {
		if _, ok2 := thisToBool.(Callable); !ok2 {
			return nil, errors.NewTypeError([]string{FunctionName, CallableName}, thisToBool.(IObject).TypeName())
		} else {
			thisBool, transformationError = thisToBool.(Callable).Call(arguments[0])
		}
	} else {
		thisBool, transformationError = thisToBool.(Callable).Call(arguments[0])
	}
	if transformationError != nil {
		return nil, transformationError
	}

	var rightToBool interface{}
	rightToBool, foundError = arguments[1].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var rightBool IObject
	if _, ok := rightToBool.(*Function); !ok {
		if _, ok2 := rightToBool.(Callable); !ok2 {
			return nil, errors.NewTypeError([]string{FunctionName, CallableName}, rightToBool.(IObject).TypeName())
		} else {
			rightBool, transformationError = rightToBool.(Callable).Call(arguments[0])
		}
	} else {
		rightBool, transformationError = rightToBool.(Callable).Call(arguments[0])
	}
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(Empty, nil, thisBool.SymbolTable().Master, thisBool.SymbolTable().Parent, []interface{}{thisBool.(*Bool).Value != rightBool.(*Bool).Value}), nil
}

func ObjRightXor(arguments ...IObject) (IObject, *errors.Error) {
	thisToBool, foundError := arguments[0].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var thisBool IObject
	var transformationError *errors.Error
	if _, ok := thisToBool.(*Function); !ok {
		if _, ok2 := thisToBool.(Callable); !ok2 {
			return nil, errors.NewTypeError([]string{FunctionName, CallableName}, thisToBool.(IObject).TypeName())
		} else {
			thisBool, transformationError = thisToBool.(Callable).Call(arguments[0])
		}
	} else {
		thisBool, transformationError = thisToBool.(Callable).Call(arguments[0])
	}
	if transformationError != nil {
		return nil, transformationError
	}

	var leftToBool interface{}
	leftToBool, foundError = arguments[1].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var leftBool IObject
	if _, ok := leftToBool.(*Function); !ok {
		if _, ok2 := leftToBool.(Callable); !ok2 {
			return nil, errors.NewTypeError([]string{FunctionName, CallableName}, leftToBool.(IObject).TypeName())
		} else {
			leftBool, transformationError = leftToBool.(Callable).Call(arguments[0])
		}
	} else {
		leftBool, transformationError = leftToBool.(Callable).Call(arguments[0])
	}
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(Empty, nil, thisBool.SymbolTable().Master, thisBool.SymbolTable().Parent, []interface{}{leftBool.(*Bool).Value != thisBool.(*Bool).Value}), nil
}

func ObjEquals(arguments ...IObject) (IObject, *errors.Error) {
	return NewBool(Empty, nil, arguments[0].SymbolTable().Master, arguments[0].SymbolTable().Parent, []interface{}{arguments[0].Id() == arguments[1].Id()}), nil
}

func ObjRightEquals(arguments ...IObject) (IObject, *errors.Error) {
	return NewBool(Empty, nil, arguments[0].SymbolTable().Master, arguments[0].SymbolTable().Parent, []interface{}{arguments[1].Id() == arguments[0].Id()}), nil
}

func ObjNotEquals(arguments ...IObject) (IObject, *errors.Error) {
	return NewBool(Empty, nil, arguments[0].SymbolTable().Master, arguments[0].SymbolTable().Parent, []interface{}{arguments[0].Id() != arguments[1].Id()}), nil
}

func ObjRightNotEquals(arguments ...IObject) (IObject, *errors.Error) {
	return NewBool(Empty, nil, arguments[0].SymbolTable().Master, arguments[0].SymbolTable().Parent, []interface{}{arguments[1].Id() != arguments[0].Id()}), nil
}

func ObjNegate(arguments ...IObject) (IObject, *errors.Error) {
	thisToBool, foundError := arguments[0].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var thisBool IObject
	var transformationError *errors.Error
	if _, ok := thisToBool.(*Function); !ok {
		if _, ok2 := thisToBool.(Callable); !ok2 {
			return nil, errors.NewTypeError([]string{FunctionName, CallableName}, thisToBool.(IObject).TypeName())
		} else {
			thisBool, transformationError = thisToBool.(Callable).Call(arguments[0])
		}
	} else {
		thisBool, transformationError = thisToBool.(Callable).Call(arguments[0])
	}
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(Empty, nil, arguments[0].SymbolTable().Master, arguments[0].SymbolTable().Parent, []interface{}{!thisBool.(*Bool).Value}), nil
}

func ObjToBool(arguments ...IObject) (IObject, *errors.Error) {
	return NewBool(Empty, nil, arguments[0].SymbolTable().Master, arguments[0].SymbolTable().Parent, []interface{}{true}), nil
}

func ObjToString(arguments ...IObject) (IObject, *errors.Error) {
	return NewString(Empty, nil, arguments[0].SymbolTable().Master, arguments[0].SymbolTable().Parent,
		[]interface{}{
			fmt.Sprintf("%s-%d", arguments[0].TypeName(), arguments[0].Id()),
		}), nil
}

const FunctionName = "Function"

/*
Should have only Call
*/
type Function struct {
	*Object
	Callable Callable
}

func NewFunction(_ string, _ []ConstructorFunction, masterSymbols *utils.SymbolTable, parentSymbols *utils.SymbolTable, values []interface{}) *Function {
	callable := values[0].(Callable)
	function := &Function{
		Object: &Object{
			id:         counter.NextId(),
			typeName:   FunctionName,
			subClasses: nil,
			symbols:    utils.NewSymbolTable(masterSymbols, parentSymbols),
		},
		Callable: callable,
	}
	function.symbols.Set(Call, callable.Call)
	return function
}

const ObjName = "Object"

/*
Classes will be handled as functions
*/
// Creates an Empty Object
/*
Pre-implemented methods
And - (Done)
Or - (Done)
Xor - (Done)
Equals  - (Done)
NotEquals - (Done)
Negate - (Done)
Initialize - (Done)
ToString - (Done)
ToBool - (Done)
Copy SymbolTable only
Dir SymbolTable to Special Hash
SubClasses
Class
*/
func NewObject(
	typeName string,
	subClasses []*Type,
	masterSymbols *utils.SymbolTable, parentSymbols *utils.SymbolTable,
	_ []interface{},
) *Object {
	symbols := utils.NewSymbolTable(masterSymbols, parentSymbols)
	symbols.Symbols = map[string]interface{}{
		// IObject Creation
		Initialize: NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewGoFunctionType(1, ObjInitialize)}),
		// Unary Operations
		NegBits: NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(3)}),
		Negate:  NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewGoFunctionType(1, ObjNegate)}),
		// Binary Operations
		//// Math binary
		Add:           NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		RightAdd:      NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		Sub:           NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		RightSub:      NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		Mul:           NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		RightMul:      NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		Div:           NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		RightDiv:      NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		Mod:           NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		RightMod:      NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		Pow:           NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		RightPow:      NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		BitXor:        NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		RightBitXor:   NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		BitAnd:        NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		RightBitAnd:   NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		BitOr:         NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		RightBitOr:    NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		BitLeft:       NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		RightBitLeft:  NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		BitRight:      NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		RightBitRight: NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		//// Logical binary
		And:      NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewGoFunctionType(2, ObjAnd)}),
		RightAnd: NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewGoFunctionType(2, ObjRightAnd)}),
		Or:       NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewGoFunctionType(2, ObjOr)}),
		RightOr:  NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewGoFunctionType(2, ObjRightOr)}),
		Xor:      NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewGoFunctionType(2, ObjXor)}),
		RightXor: NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewGoFunctionType(2, ObjRightXor)}),
		//// Comparison binary
		Equals:                  NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewGoFunctionType(2, ObjEquals)}),
		RightEquals:             NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewGoFunctionType(2, ObjRightEquals)}),
		NotEquals:               NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewGoFunctionType(2, ObjNotEquals)}),
		RightNotEquals:          NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewGoFunctionType(2, ObjRightNotEquals)}),
		GreaterThan:             NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		RightGreaterThan:        NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		LessThan:                NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		RightLessThan:           NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		GreaterThanOrEqual:      NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		RightGreaterThanOrEqual: NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		LessThanOrEqual:         NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		RightLessThanOrEqual:    NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(2)}),
		// Behavior
		Assign:     NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(3)}),
		Copy:       NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(1)}),
		Dir:        NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(1)}),
		Index:      NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(1)}),
		Call:       NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(1)}),
		Iter:       NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(1)}),
		Class:      NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(1)}),
		SubClasses: NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(1)}),
		// Transformation
		ToInteger: NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(3)}),
		ToFloat:   NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(3)}),
		ToString:  NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewGoFunctionType(1, ObjToString)}),
		ToBool:    NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewGoFunctionType(1, ObjToBool)}),
		ToArray:   NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(3)}),
		ToTuple:   NewFunction(FunctionName, nil, masterSymbols, symbols, []interface{}{NewNotImplementedCallable(3)}),
	}
	return &Object{
		id:         counter.NextId(),
		typeName:   typeName,
		subClasses: subClasses,
		symbols:    symbols,
	}
}

const StringName = "String"

type String struct {
	*Object
	Value string
}

func StringAdd(arguments ...IObject) (IObject, *errors.Error) {
	this := arguments[0]
	right := arguments[1]
	if _, ok := right.(*String); !ok {
		return nil, errors.NewTypeError([]string{StringName}, right.TypeName())
	}
	return NewString(
		Empty, nil,
		this.SymbolTable().Master, this.SymbolTable().Master,
		[]interface{}{
			this.(*String).Value + right.(*String).Value,
		},
	), nil
}

func StringRightAdd(arguments ...IObject) (IObject, *errors.Error) {
	this := arguments[0]
	left := arguments[1]
	if _, ok := left.(*String); !ok {
		return nil, errors.NewTypeError([]string{StringName}, left.TypeName())
	}
	return NewString(
		Empty, nil,
		left.SymbolTable().Master, left.SymbolTable().Master,
		[]interface{}{
			left.(*String).Value + this.(*String).Value,
		},
	), nil
}

func StringToString(arguments ...IObject) (IObject, *errors.Error) {
	return NewString(
		Empty, nil,
		arguments[0].SymbolTable().Master, arguments[0].SymbolTable().Parent,
		[]interface{}{
			arguments[0].(*String).Value,
		},
	), nil
}

func StringCopy(arguments ...IObject) (IObject, *errors.Error) {
	return NewString(
		Empty, nil,
		arguments[0].SymbolTable().Master, arguments[0].SymbolTable().Parent,
		[]interface{}{arguments[0].(*String).Value},
	), nil
}

func StringToBool(arguments ...IObject) (IObject, *errors.Error) {
	return NewBool(Empty, nil, arguments[0].SymbolTable().Master, arguments[0].SymbolTable().Parent, []interface{}{len(arguments[0].(*String).Value) != 0}), nil
}

/*
Methods for Strings
Supported Binary Operations
Add Strings only - (Done)
Mul String with integer
Class
Iter
Assign
Equals
NotEquals
Index - Integer or Tuple
Copy - (Done)
ToString - (Done)
ToInteger
ToFloat
ToBool - (Done)
*/
func NewString(
	_ string, _ []Type,
	masterSymbols *utils.SymbolTable, parentSymbols *utils.SymbolTable,
	values []interface{},
) *String {
	value := values[0].(string)
	string_ := &String{
		Value: value,
	}
	string_.Object = NewObject(StringName, nil, masterSymbols, parentSymbols, nil)
	string_.Set(Add, NewFunction(FunctionName, nil, masterSymbols, string_.symbols, []interface{}{NewGoFunctionType(2, StringAdd)}))
	string_.Set(RightAdd, NewFunction(FunctionName, nil, masterSymbols, string_.symbols, []interface{}{NewGoFunctionType(2, StringRightAdd)}))
	string_.Set(ToString, NewFunction(FunctionName, nil, masterSymbols, string_.symbols, []interface{}{NewGoFunctionType(1, StringToString)}))
	string_.Set(Copy, NewFunction(FunctionName, nil, masterSymbols, string_.symbols, []interface{}{NewGoFunctionType(1, StringCopy)}))
	string_.Set(ToBool, NewFunction(FunctionName, nil, masterSymbols, string_.symbols, []interface{}{NewGoFunctionType(1, StringToBool)}))

	return string_
}

const BoolName = "Bool"

/*
Methods for Strings
Supported Binary Operations
Class
Copy
ToString
ToInteger
ToFloat
ToBool
*/
type Bool struct {
	*Object
	Value bool
}

func NewBool(_ string, _ []ConstructorFunction, masterSymbols *utils.SymbolTable, parentSymbols *utils.SymbolTable, values []interface{}) *Bool {
	value := values[0].(bool)
	return &Bool{
		Object: NewObject(BoolName, nil, masterSymbols, parentSymbols, nil),
		Value:  value,
	}
}

func SetupDefaultTypes(vm utils.VirtualMachine) {
	// Object
	vm.MasterSymbolTable().Set(ObjName,
		NewType(vm.MasterSymbolTable(), vm.MasterSymbolTable(),
			func(vm utils.VirtualMachine, m *utils.SymbolTable, p *utils.SymbolTable, initArguments []interface{}) (IObject, *errors.Error) {
				return NewObject(Empty, nil, m, p, initArguments), nil
			},
		),
	)
	// Function
	vm.MasterSymbolTable().Set(FunctionName,
		NewType(vm.MasterSymbolTable(), vm.MasterSymbolTable(),
			func(vm utils.VirtualMachine, m *utils.SymbolTable, p *utils.SymbolTable, initArguments []interface{}) (IObject, *errors.Error) {
				return NewFunction(Empty, nil, m, p, initArguments), nil
			},
		),
	)
	// String
	vm.MasterSymbolTable().Set(StringName,
		NewType(vm.MasterSymbolTable(), vm.MasterSymbolTable(),
			func(vm utils.VirtualMachine, m *utils.SymbolTable, p *utils.SymbolTable, initArguments []interface{}) (IObject, *errors.Error) {
				return NewString(Empty, nil, m, p, initArguments), nil
			},
		),
	)
}
