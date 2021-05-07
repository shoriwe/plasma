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
	for _, subClass := range obj.SubClasses() {
		_, subClassInitError := subClass.Callable.Call(obj.SymbolTable(), vm)
		if subClassInitError != nil {
			return nil, subClassInitError
		}
	}
	initialize, getError := obj.Get(Initialize)
	if getError != nil {
		return nil, getError
	}
	_, initializationError := initialize.(*Function).Callable.Call(obj.SymbolTable(), vm, append([]IObject{obj}, arguments...)...)
	if initializationError != nil {
		return nil, initializationError
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
	for _, subClass := range obj.SubClasses() {
		_, subClassInitError := subClass.Callable.Call(obj.SymbolTable(), vm)
		if subClassInitError != nil {
			return nil, subClassInitError
		}
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
	self := arguments[0]
	leftToBool, foundError := self.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var leftBool IObject
	var transformationError *errors.Error
	if _, ok := leftToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, leftToBool.(IObject).TypeName())
	}
	leftBool, transformationError = leftToBool.(*Function).Callable.Call(self.SymbolTable(), vm, self)
	if transformationError != nil {
		return nil, transformationError
	}

	other := arguments[1]
	var rightToBool interface{}
	rightToBool, foundError = other.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var rightBool IObject
	if _, ok := rightToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, rightToBool.(IObject).TypeName())
	}
	rightBool, transformationError = rightToBool.(*Function).Callable.Call(other.SymbolTable(), vm, other)
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(leftBool.SymbolTable().Parent, leftBool.(*Bool).Value && rightBool.(*Bool).Value), nil
}

func ObjRightAnd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self := arguments[1]
	rightToBool, foundError := self.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var rightBool IObject
	var transformationError *errors.Error
	if _, ok := rightToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, rightToBool.(IObject).TypeName())
	}
	rightBool, transformationError = rightToBool.(*Function).Callable.Call(self.SymbolTable(), vm, self)
	if transformationError != nil {
		return nil, transformationError
	}

	other := arguments[1]
	var leftToBool interface{}
	leftToBool, foundError = other.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var leftBool IObject
	if _, ok := leftToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, leftToBool.(IObject).TypeName())
	}
	leftBool, transformationError = leftToBool.(*Function).Callable.Call(other.SymbolTable(), vm, other)
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(rightBool.SymbolTable().Parent, leftBool.(*Bool).Value && rightBool.(*Bool).Value), nil
}

func ObjOr(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self := arguments[0]
	leftToBool, foundError := self.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var transformationError *errors.Error
	if _, ok := leftToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, leftToBool.(IObject).TypeName())
	}
	var leftBool IObject
	leftBool, transformationError = leftToBool.(*Function).Callable.Call(self.SymbolTable(), vm, self)
	if transformationError != nil {
		return nil, transformationError
	}

	other := arguments[1]
	var rightToBool interface{}
	rightToBool, foundError = other.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := rightToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, rightToBool.(IObject).TypeName())
	}
	var rightBool IObject
	rightBool, transformationError = rightToBool.(*Function).Callable.Call(other.SymbolTable(), vm, other)
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(leftBool.SymbolTable().Parent, leftBool.(*Bool).Value || rightBool.(*Bool).Value), nil
}

func ObjRightOr(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self := arguments[1]
	rightToBool, foundError := self.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var transformationError *errors.Error
	if _, ok := rightToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, rightToBool.(IObject).TypeName())
	}
	var rightBool IObject
	rightBool, transformationError = rightToBool.(*Function).Callable.Call(self.SymbolTable(), vm, self)
	if transformationError != nil {
		return nil, transformationError
	}

	other := arguments[0]
	var leftToBool interface{}
	leftToBool, foundError = other.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := leftToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, leftToBool.(IObject).TypeName())
	}
	var leftBool IObject
	leftBool, transformationError = leftToBool.(*Function).Callable.Call(other.SymbolTable(), vm, other)
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(rightBool.SymbolTable().Parent, leftBool.(*Bool).Value || rightBool.(*Bool).Value), nil
}

func ObjXor(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self := arguments[0]
	leftToBool, foundError := self.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var transformationError *errors.Error
	if _, ok := leftToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, leftToBool.(IObject).TypeName())
	}
	var leftBool IObject
	leftBool, transformationError = leftToBool.(*Function).Callable.Call(self.SymbolTable(), vm, self)
	if transformationError != nil {
		return nil, transformationError
	}

	other := arguments[1]
	var rightToBool interface{}
	rightToBool, foundError = arguments[1].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := rightToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, rightToBool.(IObject).TypeName())
	}
	var rightBool IObject
	rightBool, transformationError = rightToBool.(*Function).Callable.Call(other.SymbolTable(), vm, other)
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(leftBool.SymbolTable().Parent, leftBool.(*Bool).Value != rightBool.(*Bool).Value), nil
}

func ObjRightXor(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self := arguments[1]
	leftToBool, foundError := self.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var transformationError *errors.Error
	if _, ok := leftToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, leftToBool.(IObject).TypeName())
	}
	var leftBool IObject
	leftBool, transformationError = leftToBool.(*Function).Callable.Call(self.SymbolTable(), vm, self)
	if transformationError != nil {
		return nil, transformationError
	}

	other := arguments[0]
	var rightToBool interface{}
	rightToBool, foundError = arguments[1].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := rightToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, rightToBool.(IObject).TypeName())
	}
	var rightBool IObject
	rightBool, transformationError = rightToBool.(*Function).Callable.Call(other.SymbolTable(), vm, other)
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(leftBool.SymbolTable().Parent, rightBool.(*Bool).Value != leftBool.(*Bool).Value), nil
}

func ObjEquals(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self := arguments[0]
	other := arguments[1]
	return NewBool(self.SymbolTable().Parent, self.Id() == other.Id()), nil
}

func ObjRightEquals(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self := arguments[1]
	other := arguments[0]
	return NewBool(self.SymbolTable().Parent, other.Id() == self.Id()), nil
}

func ObjNotEquals(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self := arguments[0]
	other := arguments[1]
	return NewBool(self.SymbolTable().Parent, self.Id() != other.Id()), nil
}

func ObjRightNotEquals(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self := arguments[1]
	other := arguments[0]
	return NewBool(self.SymbolTable().Parent, other.Id() != self.Id()), nil
}

func ObjNegate(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self := arguments[0]
	thisToBool, foundError := self.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	var transformationError *errors.Error
	if _, ok := thisToBool.(*Function); !ok {
		return nil, errors.NewTypeError([]string{FunctionName, CallableName}, thisToBool.(IObject).TypeName())
	}
	var thisBool IObject
	thisBool, transformationError = thisToBool.(Callable).Call(self.SymbolTable(), vm, self)
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(self.SymbolTable().Parent, !thisBool.(*Bool).Value), nil
}

func ObjToBool(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self := arguments[0]
	return NewBool(self.SymbolTable().Parent, true), nil
}

func ObjToString(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self := arguments[0]
	return NewString(self.SymbolTable().Parent,
		fmt.Sprintf("%s-%d", self.TypeName(), self.Id())), nil
}

const FunctionName = "Function"

/*
   Should have only Call
*/
type Function struct {
	*Object
	Callable Callable
}

func NewFunction(parentSymbols *SymbolTable, callable Callable) *Function {
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
		Initialize: NewFunction(symbols, NewBuiltInFunction(1, ObjInitialize)),
		// Unary Operations
		NegBits: NewFunction(symbols, NewNotImplementedCallable(1)),
		Negate:  NewFunction(symbols, NewBuiltInFunction(1, ObjNegate)),
		// Binary Operations
		//// Math binary
		Add:           NewFunction(symbols, NewNotImplementedCallable(2)),
		RightAdd:      NewFunction(symbols, NewNotImplementedCallable(2)),
		Sub:           NewFunction(symbols, NewNotImplementedCallable(2)),
		RightSub:      NewFunction(symbols, NewNotImplementedCallable(2)),
		Mul:           NewFunction(symbols, NewNotImplementedCallable(2)),
		RightMul:      NewFunction(symbols, NewNotImplementedCallable(2)),
		Div:           NewFunction(symbols, NewNotImplementedCallable(2)),
		RightDiv:      NewFunction(symbols, NewNotImplementedCallable(2)),
		Mod:           NewFunction(symbols, NewNotImplementedCallable(2)),
		RightMod:      NewFunction(symbols, NewNotImplementedCallable(2)),
		Pow:           NewFunction(symbols, NewNotImplementedCallable(2)),
		RightPow:      NewFunction(symbols, NewNotImplementedCallable(2)),
		BitXor:        NewFunction(symbols, NewNotImplementedCallable(2)),
		RightBitXor:   NewFunction(symbols, NewNotImplementedCallable(2)),
		BitAnd:        NewFunction(symbols, NewNotImplementedCallable(2)),
		RightBitAnd:   NewFunction(symbols, NewNotImplementedCallable(2)),
		BitOr:         NewFunction(symbols, NewNotImplementedCallable(2)),
		RightBitOr:    NewFunction(symbols, NewNotImplementedCallable(2)),
		BitLeft:       NewFunction(symbols, NewNotImplementedCallable(2)),
		RightBitLeft:  NewFunction(symbols, NewNotImplementedCallable(2)),
		BitRight:      NewFunction(symbols, NewNotImplementedCallable(2)),
		RightBitRight: NewFunction(symbols, NewNotImplementedCallable(2)),
		//// Logical binary
		And:      NewFunction(symbols, NewBuiltInFunction(2, ObjAnd)),
		RightAnd: NewFunction(symbols, NewBuiltInFunction(2, ObjRightAnd)),
		Or:       NewFunction(symbols, NewBuiltInFunction(2, ObjOr)),
		RightOr:  NewFunction(symbols, NewBuiltInFunction(2, ObjRightOr)),
		Xor:      NewFunction(symbols, NewBuiltInFunction(2, ObjXor)),
		RightXor: NewFunction(symbols, NewBuiltInFunction(2, ObjRightXor)),
		//// Comparison binary
		Equals:                  NewFunction(symbols, NewBuiltInFunction(2, ObjEquals)),
		RightEquals:             NewFunction(symbols, NewBuiltInFunction(2, ObjRightEquals)),
		NotEquals:               NewFunction(symbols, NewBuiltInFunction(2, ObjNotEquals)),
		RightNotEquals:          NewFunction(symbols, NewBuiltInFunction(2, ObjRightNotEquals)),
		GreaterThan:             NewFunction(symbols, NewNotImplementedCallable(2)),
		RightGreaterThan:        NewFunction(symbols, NewNotImplementedCallable(2)),
		LessThan:                NewFunction(symbols, NewNotImplementedCallable(2)),
		RightLessThan:           NewFunction(symbols, NewNotImplementedCallable(2)),
		GreaterThanOrEqual:      NewFunction(symbols, NewNotImplementedCallable(2)),
		RightGreaterThanOrEqual: NewFunction(symbols, NewNotImplementedCallable(2)),
		LessThanOrEqual:         NewFunction(symbols, NewNotImplementedCallable(2)),
		RightLessThanOrEqual:    NewFunction(symbols, NewNotImplementedCallable(2)),
		// Behavior
		Assign:     NewFunction(symbols, NewNotImplementedCallable(3)),
		Copy:       NewFunction(symbols, NewNotImplementedCallable(1)),
		Dir:        NewFunction(symbols, NewNotImplementedCallable(1)),
		Index:      NewFunction(symbols, NewNotImplementedCallable(1)),
		Call:       NewFunction(symbols, NewNotImplementedCallable(1)),
		Iter:       NewFunction(symbols, NewNotImplementedCallable(1)),
		Class:      NewFunction(symbols, NewNotImplementedCallable(1)),
		SubClasses: NewFunction(symbols, NewNotImplementedCallable(1)),
		// Transformation
		ToInteger: NewFunction(symbols, NewNotImplementedCallable(3)),
		ToFloat:   NewFunction(symbols, NewNotImplementedCallable(3)),
		ToString:  NewFunction(symbols, NewBuiltInFunction(1, ObjToString)),
		ToBool:    NewFunction(symbols, NewBuiltInFunction(1, ObjToBool)),
		ToArray:   NewFunction(symbols, NewNotImplementedCallable(3)),
		ToTuple:   NewFunction(symbols, NewNotImplementedCallable(3)),
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
	self := arguments[0]
	other := arguments[1]
	if _, ok := other.(*String); !ok {
		return nil, errors.NewTypeError([]string{StringName}, other.TypeName())
	}
	return NewString(
		self.SymbolTable().Parent,
		self.(*String).Value+other.(*String).Value,
	), nil
}

func StringRightAdd(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self := arguments[0]
	other := arguments[1]
	if _, ok := other.(*String); !ok {
		return nil, errors.NewTypeError([]string{StringName}, other.TypeName())
	}
	return NewString(
		self.SymbolTable().Parent,
		other.(*String).Value+self.(*String).Value,
	), nil
}

func StringToString(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self := arguments[0]
	return NewString(
		self.SymbolTable().Parent,
		self.(*String).Value,
	), nil
}

func StringCopy(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self := arguments[0]
	return NewString(
		self.SymbolTable().Parent,
		self.(*String).Value,
	), nil
}

func StringToBool(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self := arguments[0]
	return NewBool(self.SymbolTable().Parent, len(self.(*String).Value) != 0), nil
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
	parentSymbols *SymbolTable,
	value string,
) *String {
	string_ := &String{
		Value: value,
	}
	string_.Object = NewObject(StringName, nil, parentSymbols)
	string_.Set(Add, NewFunction(string_.symbols, NewBuiltInFunction(2, StringAdd)))
	string_.Set(RightAdd, NewFunction(string_.symbols, NewBuiltInFunction(2, StringRightAdd)))
	string_.Set(ToString, NewFunction(string_.symbols, NewBuiltInFunction(1, StringToString)))
	string_.Set(Copy, NewFunction(string_.symbols, NewBuiltInFunction(1, StringCopy)))
	string_.Set(ToBool, NewFunction(string_.symbols, NewBuiltInFunction(1, StringToBool)))
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

const TrueName = "True"
const FalseName = "False"

type Bool struct {
	*Object
	Value bool
}

func BoolToString(_ VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self := arguments[0]
	if self.(*Bool).Value {
		return NewString(self.SymbolTable().Parent, TrueName), nil
	}
	return NewString(self.SymbolTable().Parent, FalseName), nil
}

func NewBool(parentSymbols *SymbolTable, value bool) *Bool {
	bool_ := &Bool{
		Object: NewObject(BoolName, nil, parentSymbols),
		Value:  value,
	}
	bool_.Set(ToString, NewFunction(bool_.symbols, NewBuiltInFunction(1, BoolToString)))
	return bool_
}

func SetupDefaultTypes(vm VirtualMachine) {
	// String
	vm.MasterSymbolTable().Set(StringName,
		NewFunction(vm.MasterSymbolTable(),
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
