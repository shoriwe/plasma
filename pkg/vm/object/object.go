package object

import (
	"github.com/shoriwe/gruby/pkg/errors"
	"github.com/shoriwe/gruby/pkg/vm/utils"
	"sync"
)

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

type Callable interface {
	NumberOfArguments() int
	Call(...IObject) (IObject, []interface{}, *errors.Error) // This should return directly the object or the code of the function
}

type PlasmaFunctionType struct {
	numberOfArguments int
	Code              []interface{}
}

func (p *PlasmaFunctionType) NumberOfArguments() int {
	return p.numberOfArguments
}

func (p *PlasmaFunctionType) Call(arguments ...IObject) (IObject, []interface{}, *errors.Error) {
	if len(arguments) != p.numberOfArguments {
		return nil, nil, errors.NewInvalidNumberOfArguments(len(arguments), p.numberOfArguments)
	}
	return nil, p.Code, nil
}

type GoFunctionType struct {
	numberOfArguments int
	callback          GoFunctionCallback
}

func (g *GoFunctionType) NumberOfArguments() int {
	return g.numberOfArguments
}

func (g *GoFunctionType) Call(arguments ...IObject) (IObject, []interface{}, *errors.Error) {
	if len(arguments) != g.numberOfArguments {
		return nil, nil, errors.NewInvalidNumberOfArguments(len(arguments), g.numberOfArguments)
	}
	result, callError := g.callback(arguments...)
	if callError != nil {
		return nil, nil, callError
	}
	return result, nil, nil
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
	SubClasses() []IObject
	Get(string) (interface{}, *errors.Error)
	Set(string, interface{})
}

// MetaClass for IObject
/*
Pre-implemented methods
Copy SymbolTable only
Dir SymbolTable to Special Hash
And
Or
Xor
Equals
NotEquals
SubClasses
Initialize - (Done)
ToString
ToBool
*/
type Object struct {
	id         uint
	typeName   string
	subClasses []IObject
	symbols    *utils.SymbolTable
}

func (o *Object) Id() uint {
	return o.id
}

func (o *Object) SubClasses() []IObject {
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

// Creates an Empty Object Class
func NewObject(
	typeName string,
	subClasses []IObject,
	masterSymbols *utils.SymbolTable, parentSymbols *utils.SymbolTable,
) *Object {
	symbols := utils.NewSymbolTable(masterSymbols, parentSymbols)
	symbols.Symbols = map[string]interface{}{
		// IObject Creation
		Initialize: NewGoFunctionType(1, ObjInitialize),
		// Unary Operations
		NegBits: NewNotImplementedCallable(1),
		Negate:  NewNotImplementedCallable(1),
		// Binary Operations
		Add:                     NewNotImplementedCallable(2),
		RightAdd:                NewNotImplementedCallable(2),
		Sub:                     NewNotImplementedCallable(2),
		RightSub:                NewNotImplementedCallable(2),
		Mul:                     NewNotImplementedCallable(2),
		RightMul:                NewNotImplementedCallable(2),
		Div:                     NewNotImplementedCallable(2),
		RightDiv:                NewNotImplementedCallable(2),
		Mod:                     NewNotImplementedCallable(2),
		RightMod:                NewNotImplementedCallable(2),
		Pow:                     NewNotImplementedCallable(2),
		RightPow:                NewNotImplementedCallable(2),
		BitXor:                  NewNotImplementedCallable(2),
		RightBitXor:             NewNotImplementedCallable(2),
		BitAnd:                  NewNotImplementedCallable(2),
		RightBitAnd:             NewNotImplementedCallable(2),
		BitOr:                   NewNotImplementedCallable(2),
		RightBitOr:              NewNotImplementedCallable(2),
		BitLeft:                 NewNotImplementedCallable(2),
		RightBitLeft:            NewNotImplementedCallable(2),
		BitRight:                NewNotImplementedCallable(2),
		RightBitRight:           NewNotImplementedCallable(2),
		And:                     NewNotImplementedCallable(2),
		RightAnd:                NewNotImplementedCallable(2),
		Or:                      NewNotImplementedCallable(2),
		RightOr:                 NewNotImplementedCallable(2),
		Xor:                     NewNotImplementedCallable(2),
		RightXor:                NewNotImplementedCallable(2),
		Equals:                  NewNotImplementedCallable(2),
		RightEquals:             NewNotImplementedCallable(2),
		NotEquals:               NewNotImplementedCallable(2),
		RightNotEquals:          NewNotImplementedCallable(2),
		GreaterThan:             NewNotImplementedCallable(2),
		RightGreaterThan:        NewNotImplementedCallable(2),
		LessThan:                NewNotImplementedCallable(2),
		RightLessThan:           NewNotImplementedCallable(2),
		GreaterThanOrEqual:      NewNotImplementedCallable(2),
		RightGreaterThanOrEqual: NewNotImplementedCallable(2),
		LessThanOrEqual:         NewNotImplementedCallable(2),
		RightLessThanOrEqual:    NewNotImplementedCallable(2),
		// Behavior
		Copy:       NewNotImplementedCallable(1),
		Dir:        NewNotImplementedCallable(1),
		Index:      NewNotImplementedCallable(2),
		Call:       NewNotImplementedCallable(1),
		Iter:       NewNotImplementedCallable(1),
		Class:      NewNotImplementedCallable(1),
		SubClasses: NewNotImplementedCallable(1),
		// Transformation
		ToInteger: NewNotImplementedCallable(1),
		ToFloat:   NewNotImplementedCallable(1),
		ToString:  NewNotImplementedCallable(1),
		ToBool:    NewNotImplementedCallable(1),
		ToArray:   NewNotImplementedCallable(1),
		ToTuple:   NewNotImplementedCallable(1),
	}
	return &Object{
		id:         counter.NextId(),
		typeName:   typeName,
		subClasses: subClasses,
		symbols:    symbols,
	}
}

const StringName = "String"

/*
Methods for Strings
Supported Binary Operations
Add Strings only - (Done)
Mul String with integer
And
Or
Xor
Construct
Initialize
Class
Iter Dir
Index - Integer or Tuple
Copy - (Done)
ToString - (Done)
ToInteger
ToFloat
ToBool
*/
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
		this.SymbolTable().Master, this.SymbolTable().Master,
		this.(*String).Value+right.(*String).Value,
	), nil
}

func StringRightAdd(arguments ...IObject) (IObject, *errors.Error) {
	this := arguments[0]
	left := arguments[1]
	if _, ok := left.(*String); !ok {
		return nil, errors.NewTypeError([]string{StringName}, left.TypeName())
	}
	return NewString(
		left.SymbolTable().Master, left.SymbolTable().Master,
		left.(*String).Value+this.(*String).Value,
	), nil
}

func StringToString(arguments ...IObject) (IObject, *errors.Error) {
	this := arguments[0]
	return NewString(
		this.SymbolTable().Master, this.SymbolTable().Parent,
		this.(*String).Value,
	), nil
}

func StringCopy(arguments ...IObject) (IObject, *errors.Error) {
	return StringToString(arguments...)
}

func NewString(
	masterSymbols *utils.SymbolTable, parentSymbols *utils.SymbolTable,
	value string,
) *String {
	result := &String{
		Value: value,
	}
	result.Object = NewObject(StringName, nil, masterSymbols, parentSymbols)
	result.Set(Add, NewGoFunctionType(2, StringAdd))
	result.Set(RightAdd, NewGoFunctionType(2, StringRightAdd))
	result.Set(ToString, NewGoFunctionType(1, StringToString))
	result.Set(Copy, NewGoFunctionType(1, StringCopy))
	return result
}
