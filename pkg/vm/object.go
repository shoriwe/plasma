package vm

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/errors"
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

type FunctionCallback func(VirtualMachine, ...IObject) (IObject, *errors.Error)

type ConstructorFunction func(VirtualMachine, *SymbolTable, *SymbolTable, []interface{}) (IObject, *errors.Error)

const CallableName = "Callable"

type Callable interface {
	NumberOfArguments() int
	Call(*SymbolTable, VirtualMachine, ...IObject) (IObject, *errors.Error) // This should return directly the object or the code of the function
}

type PlasmaFunction struct {
	numberOfArguments int
	Code              []interface{}
}

func (p *PlasmaFunction) NumberOfArguments() int {
	return p.numberOfArguments
}

func (p *PlasmaFunction) Call(parent *SymbolTable, vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	if len(arguments) != p.numberOfArguments {
		return nil, errors.NewInvalidNumberOfArguments(len(arguments), p.numberOfArguments)
	}
	vmCopy := vm.New(parent)
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

type BuiltInFunction struct {
	numberOfArguments int
	callback          FunctionCallback
}

func (g *BuiltInFunction) NumberOfArguments() int {
	return g.numberOfArguments
}

func (g *BuiltInFunction) Call(_ *SymbolTable, vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	if len(arguments) != g.numberOfArguments {
		return nil, errors.NewInvalidNumberOfArguments(len(arguments), g.numberOfArguments)
	}
	result, callError := g.callback(vm, arguments...)
	if callError != nil {
		return nil, callError
	}
	return result, nil
}

func NewBuiltInFunction(numberOfArguments int, callback FunctionCallback) *BuiltInFunction {
	return &BuiltInFunction{
		numberOfArguments: numberOfArguments,
		callback:          callback,
	}
}

type PlasmaConstructor struct {
	numberOfArguments int
	Code              []interface{}
}

func (c *PlasmaConstructor) NumberOfArguments() int {
	return c.numberOfArguments
}

func (c *PlasmaConstructor) Call(parent *SymbolTable, vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	if len(arguments) != c.numberOfArguments {
		return nil, errors.NewInvalidNumberOfArguments(len(arguments), c.numberOfArguments)
	}
	obj := NewObject(Empty, nil, parent)
	vmCopy := vm.New(obj.SymbolTable())
	vmInitializationError := vmCopy.Initialize(c.Code)
	if vmInitializationError != nil {
		return nil, vmInitializationError
	}
	_, executionError := vmCopy.Execute()
	if executionError != nil {
		return nil, executionError
	}
	initialize, getError := obj.Get(Initialize)
	if getError != nil {
		return nil, getError
	}
	_, objectInitializationError := initialize.(*Function).Callable.Call(obj.SymbolTable(), vm, obj)
	if objectInitializationError != nil {
		return nil, objectInitializationError
	}
	return obj, nil
}

// This should  be used  at function definition time
func NewPlasmaConstructor(numberOfArguments int, code []interface{}) *PlasmaConstructor {
	return &PlasmaConstructor{
		numberOfArguments: numberOfArguments,
		Code:              code,
	}
}

type BuiltInConstructor struct {
	numberOfArguments int
	callback          FunctionCallback
}

func (c *BuiltInConstructor) NumberOfArguments() int {
	return c.numberOfArguments
}

func (c *BuiltInConstructor) Call(_ *SymbolTable, vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	if len(arguments) != c.numberOfArguments {
		return nil, errors.NewInvalidNumberOfArguments(len(arguments), c.numberOfArguments)
	}
	obj, callError := c.callback(vm, arguments...)
	if callError != nil {
		return nil, callError
	}
	initialize, getError := obj.Get(Initialize)
	if getError != nil {
		return nil, getError
	}
	_, initializationError := initialize.(*Function).Callable.Call(obj.SymbolTable(), vm, obj)
	if initializationError != nil {
		return nil, initializationError
	}
	return obj, nil
}

func NewBuiltInConstructor(numberOfArguments int, callback FunctionCallback) *BuiltInConstructor {
	return &BuiltInConstructor{
		numberOfArguments: numberOfArguments,
		callback:          callback,
	}
}

func NewNotImplementedCallable(numberOfArguments int) *BuiltInFunction {
	return NewBuiltInFunction(numberOfArguments, func(_ VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
		return nil, errors.NewNameNotFoundError()
	})
}

type IObject interface {
	Id() uint
	TypeName() string
	SymbolTable() *SymbolTable
	SubClasses() []*Function
	Get(string) (IObject, *errors.Error)
	Set(string, IObject)
	VirtualMachine() VirtualMachine
}

// MetaClass for IObject
type Object struct {
	id             uint
	typeName       string
	subClasses     []*Function
	symbols        *SymbolTable
	virtualMachine VirtualMachine
}

func (o *Object) VirtualMachine() VirtualMachine {
	return o.virtualMachine
}

func (o *Object) Id() uint {
	return o.id
}

func (o *Object) SubClasses() []*Function {
	return o.subClasses
}

func (o *Object) Get(symbol string) (IObject, *errors.Error) {
	return o.symbols.GetSelf(symbol)
}

func (o *Object) Set(symbol string, object IObject) {
	o.symbols.Set(symbol, object)
}

func (o *Object) TypeName() string {
	return o.typeName
}

func (o *Object) SymbolTable() *SymbolTable {
	return o.symbols
}

func ObjInitialize(_ VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	return nil, nil
}

func ObjAnd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	thisToBool, foundError := arguments[0].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var thisBool IObject
	var transformationError *errors.Error
	if _, ok := thisToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, thisToBool.(IObject).TypeName())
	}
	thisBool, transformationError = thisToBool.(*Function).Callable.Call(arguments[0].SymbolTable(), vm, arguments[0])
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
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, rightToBool.(IObject).TypeName())
	}
	rightBool, transformationError = rightToBool.(*Function).Callable.Call(arguments[1].SymbolTable(), vm, arguments[1])
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(Empty, nil, thisBool.SymbolTable().Parent, thisBool.(*Bool).Value && rightBool.(*Bool).Value), nil
}

func ObjRightAnd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	thisToBool, foundError := arguments[0].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var thisBool IObject
	var transformationError *errors.Error
	if _, ok := thisToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, thisToBool.(IObject).TypeName())
	}
	thisBool, transformationError = thisToBool.(*Function).Callable.Call(arguments[0].SymbolTable(), vm, arguments[0])
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
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, leftToBool.(IObject).TypeName())
	}
	leftBool, transformationError = leftToBool.(*Function).Callable.Call(arguments[1].SymbolTable(), vm, arguments[1])
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(Empty, nil, thisBool.SymbolTable().Parent, leftBool.(*Bool).Value && thisBool.(*Bool).Value), nil
}

func ObjOr(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	thisToBool, foundError := arguments[0].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var transformationError *errors.Error
	if _, ok := thisToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, thisToBool.(IObject).TypeName())
	}
	var thisBool IObject
	thisBool, transformationError = thisToBool.(*Function).Callable.Call(arguments[0].SymbolTable(), vm, arguments[0])
	if transformationError != nil {
		return nil, transformationError
	}

	var rightToBool interface{}
	rightToBool, foundError = arguments[1].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := rightToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, rightToBool.(IObject).TypeName())
	}
	var rightBool IObject
	rightBool, transformationError = rightToBool.(*Function).Callable.Call(arguments[1].SymbolTable(), vm, arguments[1])
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(Empty, nil, thisBool.SymbolTable().Parent, thisBool.(*Bool).Value || rightBool.(*Bool).Value), nil
}

func ObjRightOr(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	leftToBool, foundError := arguments[0].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var transformationError *errors.Error
	if _, ok := leftToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, leftToBool.(IObject).TypeName())
	}
	var leftBool IObject
	leftBool, transformationError = leftToBool.(*Function).Callable.Call(arguments[0].SymbolTable(), vm, arguments[0])
	if transformationError != nil {
		return nil, transformationError
	}

	var rightToBool interface{}
	rightToBool, foundError = arguments[1].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := rightToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, rightToBool.(IObject).TypeName())
	}
	var rightBool IObject
	rightBool, transformationError = rightToBool.(*Function).Callable.Call(arguments[0].SymbolTable(), vm, arguments[0])
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(Empty, nil, leftBool.SymbolTable().Parent, rightBool.(*Bool).Value || leftBool.(*Bool).Value), nil
}

func ObjXor(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	leftToBool, foundError := arguments[0].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var transformationError *errors.Error
	if _, ok := leftToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, leftToBool.(IObject).TypeName())
	}
	var leftBool IObject
	leftBool, transformationError = leftToBool.(*Function).Callable.Call(arguments[0].SymbolTable(), vm, arguments[0])
	if transformationError != nil {
		return nil, transformationError
	}

	var rightToBool interface{}
	rightToBool, foundError = arguments[1].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := rightToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, rightToBool.(IObject).TypeName())
	}
	var rightBool IObject
	rightBool, transformationError = rightToBool.(*Function).Callable.Call(arguments[0].SymbolTable(), vm, arguments[0])
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(Empty, nil, leftBool.SymbolTable().Parent, leftBool.(*Bool).Value != rightBool.(*Bool).Value), nil
}

func ObjRightXor(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	leftToBool, foundError := arguments[0].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var transformationError *errors.Error
	if _, ok := leftToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, leftToBool.(IObject).TypeName())
	}
	var leftBool IObject
	leftBool, transformationError = leftToBool.(*Function).Callable.Call(arguments[0].SymbolTable(), vm, arguments[0])
	if transformationError != nil {
		return nil, transformationError
	}

	var rightToBool interface{}
	rightToBool, foundError = arguments[1].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := rightToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, rightToBool.(IObject).TypeName())
	}
	var rightBool IObject
	rightBool, transformationError = rightToBool.(*Function).Callable.Call(arguments[0].SymbolTable(), vm, arguments[0])
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(Empty, nil, leftBool.SymbolTable().Parent, rightBool.(*Bool).Value != leftBool.(*Bool).Value), nil
}

func ObjEquals(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	return NewBool(Empty, nil, arguments[0].SymbolTable().Parent, arguments[0].Id() == arguments[1].Id()), nil
}

func ObjRightEquals(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	return NewBool(Empty, nil, arguments[0].SymbolTable().Parent, arguments[1].Id() == arguments[0].Id()), nil
}

func ObjNotEquals(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	return NewBool(Empty, nil, arguments[0].SymbolTable().Parent, arguments[0].Id() != arguments[1].Id()), nil
}

func ObjRightNotEquals(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	return NewBool(Empty, nil, arguments[0].SymbolTable().Parent, arguments[1].Id() != arguments[0].Id()), nil
}

func ObjNegate(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	thisToBool, foundError := arguments[0].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var transformationError *errors.Error
	if _, ok := thisToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, thisToBool.(IObject).TypeName())
	}
	var thisBool IObject
	thisBool, transformationError = thisToBool.(Callable).Call(arguments[0].SymbolTable(), vm, arguments[0])
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(Empty, nil, arguments[0].SymbolTable().Parent, !thisBool.(*Bool).Value), nil
}

func ObjToBool(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	return NewBool(Empty, nil, arguments[0].SymbolTable().Parent, true), nil
}

func ObjToString(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	return NewString(Empty, nil, arguments[0].SymbolTable().Parent,
		fmt.Sprintf("%s-%d", arguments[0].TypeName(), arguments[0].Id())), nil
}

const FunctionName = "Function"

/*
   Should have only Call
*/
type Function struct {
	*Object
	Callable Callable
}

func NewFunction(_ string, _ []ConstructorFunction, parentSymbols *SymbolTable, callable Callable) *Function {
	function := &Function{
		Object: &Object{
			id:         counter.NextId(),
			typeName:   FunctionName,
			subClasses: nil,
			symbols:    NewSymbolTable(parentSymbols),
		},
		Callable: callable,
	}
	return function
}

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
	subClasses []*Function,
	parentSymbols *SymbolTable,
) *Object {
	symbols := NewSymbolTable(parentSymbols)
	symbols.Symbols = map[string]IObject{
		// IObject Creation
		Initialize: NewFunction(FunctionName, nil, symbols, NewBuiltInFunction(1, ObjInitialize)),
		// Unary Operations
		NegBits: NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(3)),
		Negate:  NewFunction(FunctionName, nil, symbols, NewBuiltInFunction(1, ObjNegate)),
		// Binary Operations
		//// Math binary
		Add:           NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		RightAdd:      NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		Sub:           NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		RightSub:      NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		Mul:           NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		RightMul:      NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		Div:           NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		RightDiv:      NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		Mod:           NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		RightMod:      NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		Pow:           NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		RightPow:      NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		BitXor:        NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		RightBitXor:   NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		BitAnd:        NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		RightBitAnd:   NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		BitOr:         NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		RightBitOr:    NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		BitLeft:       NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		RightBitLeft:  NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		BitRight:      NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		RightBitRight: NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		//// Logical binary
		And:      NewFunction(FunctionName, nil, symbols, NewBuiltInFunction(2, ObjAnd)),
		RightAnd: NewFunction(FunctionName, nil, symbols, NewBuiltInFunction(2, ObjRightAnd)),
		Or:       NewFunction(FunctionName, nil, symbols, NewBuiltInFunction(2, ObjOr)),
		RightOr:  NewFunction(FunctionName, nil, symbols, NewBuiltInFunction(2, ObjRightOr)),
		Xor:      NewFunction(FunctionName, nil, symbols, NewBuiltInFunction(2, ObjXor)),
		RightXor: NewFunction(FunctionName, nil, symbols, NewBuiltInFunction(2, ObjRightXor)),
		//// Comparison binary
		Equals:                  NewFunction(FunctionName, nil, symbols, NewBuiltInFunction(2, ObjEquals)),
		RightEquals:             NewFunction(FunctionName, nil, symbols, NewBuiltInFunction(2, ObjRightEquals)),
		NotEquals:               NewFunction(FunctionName, nil, symbols, NewBuiltInFunction(2, ObjNotEquals)),
		RightNotEquals:          NewFunction(FunctionName, nil, symbols, NewBuiltInFunction(2, ObjRightNotEquals)),
		GreaterThan:             NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		RightGreaterThan:        NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		LessThan:                NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		RightLessThan:           NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		GreaterThanOrEqual:      NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		RightGreaterThanOrEqual: NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		LessThanOrEqual:         NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		RightLessThanOrEqual:    NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(2)),
		// Behavior
		Assign:     NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(3)),
		Copy:       NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(1)),
		Dir:        NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(1)),
		Index:      NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(1)),
		Call:       NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(1)),
		Iter:       NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(1)),
		Class:      NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(1)),
		SubClasses: NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(1)),
		// Transformation
		ToInteger: NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(3)),
		ToFloat:   NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(3)),
		ToString:  NewFunction(FunctionName, nil, symbols, NewBuiltInFunction(1, ObjToString)),
		ToBool:    NewFunction(FunctionName, nil, symbols, NewBuiltInFunction(1, ObjToBool)),
		ToArray:   NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(3)),
		ToTuple:   NewFunction(FunctionName, nil, symbols, NewNotImplementedCallable(3)),
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

func StringAdd(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	this := arguments[0]
	right := arguments[1]
	if _, ok := right.(*String); !ok {
		return nil, errors.NewTypeError([]string{StringName}, right.TypeName())
	}
	return NewString(
		Empty, nil,
		this.SymbolTable().Parent,
		this.(*String).Value+right.(*String).Value,
	), nil
}

func StringRightAdd(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	this := arguments[0]
	left := arguments[1]
	if _, ok := left.(*String); !ok {
		return nil, errors.NewTypeError([]string{StringName}, left.TypeName())
	}
	return NewString(
		Empty, nil,
		this.SymbolTable().Parent,
		left.(*String).Value+this.(*String).Value,
	), nil
}

func StringToString(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	return NewString(
		Empty, nil,
		arguments[0].SymbolTable().Parent,
		arguments[0].(*String).Value,
	), nil
}

func StringCopy(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	return NewString(
		Empty, nil,
		arguments[0].SymbolTable().Parent,
		arguments[0].(*String).Value,
	), nil
}

func StringToBool(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	return NewBool(Empty, nil, arguments[0].SymbolTable().Parent, len(arguments[0].(*String).Value) != 0), nil
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
	_ string, _ []*Function,
	parentSymbols *SymbolTable,
	value string,
) *String {
	string_ := &String{
		Value: value,
	}
	string_.Object = NewObject(StringName, nil, parentSymbols)
	string_.Set(Add, NewFunction(FunctionName, nil, string_.symbols, NewBuiltInFunction(2, StringAdd)))
	string_.Set(RightAdd, NewFunction(FunctionName, nil, string_.symbols, NewBuiltInFunction(2, StringRightAdd)))
	string_.Set(ToString, NewFunction(FunctionName, nil, string_.symbols, NewBuiltInFunction(1, StringToString)))
	string_.Set(Copy, NewFunction(FunctionName, nil, string_.symbols, NewBuiltInFunction(1, StringCopy)))
	string_.Set(ToBool, NewFunction(FunctionName, nil, string_.symbols, NewBuiltInFunction(1, StringToBool)))
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

func NewBool(_ string, _ []ConstructorFunction, parentSymbols *SymbolTable, value bool) *Bool {
	return &Bool{
		Object: NewObject(BoolName, nil, parentSymbols),
		Value:  value,
	}
}

func SetupDefaultTypes(vm VirtualMachine) {
	// String
	vm.MasterSymbolTable().Set(StringName,
		NewFunction(Empty, nil, vm.MasterSymbolTable(),
			NewBuiltInConstructor(1,
				func(vm2 VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
					toString, getError := arguments[0].Get(ToString)
					if getError != nil {
						return nil, getError
					}
					if _, ok := toString.(*Function); !ok {
						return nil, errors.NewTypeError([]string{FunctionName}, toString.(IObject).TypeName())
					}
					return toString.(*Function).Callable.Call(arguments[0].SymbolTable(), vm2, arguments[0])
				},
			),
		),
	)
}
