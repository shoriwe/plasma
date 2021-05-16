package vm

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/errors"
	"math"
	"sync"
)

func Repeat(s string, times int64) string {
	result := ""
	var i int64
	for i = 0; i < times; i++ {
		result += s
	}
	return result
}

const (
	Self = "self"
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

type Callable interface {
	NumberOfArguments() int
	Call() (IObject, FunctionCallback, []Code) // self should return directly the object or the code of the function
}

type Constructor interface {
	Callable
	C()
}

type PlasmaFunction struct {
	numberOfArguments int
	Code              []Code
}

func (p *PlasmaFunction) NumberOfArguments() int {
	return p.numberOfArguments
}

func (p *PlasmaFunction) Call() (IObject, IObject, []Code) {
	return nil, nil, p.Code
}

func NewPlasmaFunction(numberOfArguments int, code []Code) *PlasmaFunction {
	return &PlasmaFunction{
		numberOfArguments: numberOfArguments,
		Code:              code,
	}
}

type PlasmaClassFunction struct {
	numberOfArguments int
	Code              []Code
	Self              IObject
}

func (p *PlasmaClassFunction) NumberOfArguments() int {
	return p.numberOfArguments
}

func (p *PlasmaClassFunction) Call() (IObject, IObject, []Code) {
	return p.Self, nil, p.Code
}

func NewPlasmaClassFunction(self IObject, numberOfArguments int, code []Code) *PlasmaClassFunction {
	return &PlasmaClassFunction{
		numberOfArguments: numberOfArguments,
		Code:              code,
		Self:              self,
	}
}

type BuiltInFunction struct {
	numberOfArguments int
	callback          FunctionCallback
}

func (g *BuiltInFunction) NumberOfArguments() int {
	return g.numberOfArguments
}

func (g *BuiltInFunction) Call() (IObject, FunctionCallback, []Code) {
	return nil, g.callback, nil
}

func NewBuiltInFunction(numberOfArguments int, callback FunctionCallback) *BuiltInFunction {
	return &BuiltInFunction{
		numberOfArguments: numberOfArguments,
		callback:          callback,
	}
}

type BuiltInClassFunction struct {
	numberOfArguments int
	callback          FunctionCallback
	Self              IObject
}

func (g *BuiltInClassFunction) NumberOfArguments() int {
	return g.numberOfArguments
}

func (g *BuiltInClassFunction) Call() (IObject, FunctionCallback, []Code) {
	return g.Self, g.callback, nil
}

func NewBuiltInClassFunction(self IObject, numberOfArguments int, callback FunctionCallback) *BuiltInClassFunction {
	return &BuiltInClassFunction{
		numberOfArguments: numberOfArguments,
		callback:          callback,
		Self:              self,
	}
}

type PlasmaConstructor struct {
	Constructor
	numberOfArguments int
	Code              []Code
}

func (c *PlasmaConstructor) NumberOfArguments() int {
	return c.numberOfArguments
}

func (c *PlasmaConstructor) Call() (IObject, FunctionCallback, []Code) {
	return nil, nil, c.Code
}

// self should  be used  at function definition time
func NewPlasmaConstructor(numberOfArguments int, code []Code) *PlasmaConstructor {
	return &PlasmaConstructor{
		numberOfArguments: numberOfArguments,
		Code:              code,
	}
}

type BuiltInConstructor struct {
	Constructor
	numberOfArguments int
	callback          FunctionCallback
}

func (c *BuiltInConstructor) NumberOfArguments() int {
	return c.numberOfArguments
}

func (c *BuiltInConstructor) Call() (IObject, FunctionCallback, []Code) {
	return nil, c.callback, nil
}

func NewBuiltInConstructor(numberOfArguments int, callback FunctionCallback) *BuiltInConstructor {
	return &BuiltInConstructor{
		numberOfArguments: numberOfArguments,
		callback:          callback,
	}
}

func NewNotImplementedCallable(numberOfArguments int) *BuiltInClassFunction {
	return NewBuiltInClassFunction(nil, numberOfArguments, func(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
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
}

// MetaClass for IObject
type Object struct {
	id             uint
	typeName       string
	subClasses     []*Function
	symbols        *SymbolTable
	virtualMachine VirtualMachine
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

func CallFunction(function *Function, vm VirtualMachine, parent *SymbolTable, arguments ...IObject) (IObject, *errors.Error) {
	symbols := NewSymbolTable(parent)
	self, callback, code := function.Callable.Call()
	if self != nil {
		symbols.Set(Self, self)
	}
	vm.PushSymbolTable(symbols)
	var result IObject
	var callError *errors.Error
	if callback != nil {
		result, callError = callback(vm, arguments...)
	} else if code != nil {
		vm.LoadCode(code)
		result, callError = vm.Execute()
	} else {
		panic("callback and code are nil")
	}
	if callError != nil {
		return nil, callError
	}
	return result, nil
}

func ObjAnd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	leftToBool, foundError := self.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := leftToBool.(*Function); !ok {
		return nil, errors.NewTypeError(leftToBool.(IObject).TypeName(), FunctionName)
	}
	leftBool, transformationError := CallFunction(leftToBool.(*Function), vm, self.SymbolTable(), self)
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
		return nil, errors.NewTypeError(rightToBool.(IObject).TypeName(), FunctionName)
	}
	var rightBool IObject
	rightBool, transformationError = CallFunction(rightToBool.(*Function), vm, other.SymbolTable(), other)
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(leftBool.SymbolTable().Parent, leftBool.(*Bool).Value && rightBool.(*Bool).Value), nil
}

func ObjRightAnd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	rightToBool, foundError := self.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := rightToBool.(*Function); !ok {
		return nil, errors.NewTypeError(rightToBool.(IObject).TypeName(), FunctionName)
	}
	rightBool, transformationError := CallFunction(rightToBool.(*Function), vm, self.SymbolTable(), self)
	if transformationError != nil {
		return nil, transformationError
	}
	other := arguments[1]
	var leftToBool interface{}
	leftToBool, foundError = other.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := leftToBool.(*Function); !ok {
		return nil, errors.NewTypeError(leftToBool.(IObject).TypeName(), FunctionName)
	}
	var leftBool IObject
	leftBool, transformationError = CallFunction(leftToBool.(*Function), vm, other.SymbolTable(), other)
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(rightBool.SymbolTable().Parent, leftBool.(*Bool).Value && rightBool.(*Bool).Value), nil
}

func ObjOr(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	leftToBool, foundError := self.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := leftToBool.(*Function); !ok {
		return nil, errors.NewTypeError(leftToBool.(IObject).TypeName(), FunctionName)
	}
	leftBool, transformationError := CallFunction(leftToBool.(*Function), vm, self.SymbolTable(), self)
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
		return nil, errors.NewTypeError(rightToBool.(IObject).TypeName(), FunctionName)
	}
	var rightBool IObject
	rightBool, transformationError = CallFunction(rightToBool.(*Function), vm, other.SymbolTable(), other)
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(leftBool.SymbolTable().Parent, leftBool.(*Bool).Value || rightBool.(*Bool).Value), nil
}

func ObjRightOr(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	rightToBool, foundError := self.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := rightToBool.(*Function); !ok {
		return nil, errors.NewTypeError(rightToBool.(IObject).TypeName(), FunctionName)
	}
	rightBool, transformationError := CallFunction(rightToBool.(*Function), vm, self.SymbolTable(), self)
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
		return nil, errors.NewTypeError(leftToBool.(IObject).TypeName(), FunctionName)
	}
	var leftBool IObject
	leftBool, transformationError = CallFunction(leftToBool.(*Function), vm, other.SymbolTable(), other)
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(rightBool.SymbolTable().Parent, leftBool.(*Bool).Value || rightBool.(*Bool).Value), nil
}

func ObjXor(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	leftToBool, foundError := self.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := leftToBool.(*Function); !ok {
		return nil, errors.NewTypeError(leftToBool.(IObject).TypeName(), FunctionName)
	}
	leftBool, transformationError := CallFunction(leftToBool.(*Function), vm, self.SymbolTable(), self)
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
		return nil, errors.NewTypeError(rightToBool.(IObject).TypeName(), FunctionName)
	}
	var rightBool IObject
	rightBool, transformationError = CallFunction(rightToBool.(*Function), vm, other.SymbolTable(), other)
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(leftBool.SymbolTable().Parent, leftBool.(*Bool).Value != rightBool.(*Bool).Value), nil
}

func ObjRightXor(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	leftToBool, foundError := self.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := leftToBool.(*Function); !ok {
		return nil, errors.NewTypeError(leftToBool.(IObject).TypeName(), FunctionName)
	}
	leftBool, transformationError := CallFunction(leftToBool.(*Function), vm, self.SymbolTable(), self)
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
		return nil, errors.NewTypeError(rightToBool.(IObject).TypeName(), FunctionName)
	}
	var rightBool IObject
	rightBool, transformationError = CallFunction(rightToBool.(*Function), vm, other.SymbolTable(), other)
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(leftBool.SymbolTable().Parent, rightBool.(*Bool).Value != leftBool.(*Bool).Value), nil
}

func ObjEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	other := arguments[1]
	return NewBool(self.SymbolTable().Parent, self.Id() == other.Id()), nil
}

func ObjRightEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	other := arguments[0]
	return NewBool(self.SymbolTable().Parent, other.Id() == self.Id()), nil
}

func ObjNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	other := arguments[1]
	return NewBool(self.SymbolTable().Parent, self.Id() != other.Id()), nil
}

func ObjRightNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	other := arguments[0]
	return NewBool(self.SymbolTable().Parent, other.Id() != self.Id()), nil
}

func ObjNegate(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	selfToBool, foundError := self.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := selfToBool.(*Function); !ok {
		return nil, errors.NewTypeError(selfToBool.(IObject).TypeName(), FunctionName)
	}
	var selfBool IObject
	var transformationError *errors.Error
	selfBool, transformationError = CallFunction(selfToBool.(*Function), vm, self.SymbolTable(), self)
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(self.SymbolTable().Parent, !selfBool.(*Bool).Value), nil
}

func ObjToBool(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewBool(self.SymbolTable().Parent, true), nil
}

func ObjToString(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
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
	result := &Object{
		id:         counter.NextId(),
		typeName:   typeName,
		subClasses: subClasses,
	}
	result.symbols = NewSymbolTable(parentSymbols)
	result.symbols.Symbols = map[string]IObject{
		// IObject Creation
		Initialize: NewFunction(result.symbols, NewBuiltInClassFunction(result, 0, ObjInitialize)),
		// Unary Operations
		NegBits: NewFunction(result.symbols, NewNotImplementedCallable(1)),
		Negate:  NewFunction(result.symbols, NewBuiltInClassFunction(result, 0, ObjNegate)),
		// Binary Operations
		//// Math binary
		Add:           NewFunction(result.symbols, NewNotImplementedCallable(2)),
		RightAdd:      NewFunction(result.symbols, NewNotImplementedCallable(2)),
		Sub:           NewFunction(result.symbols, NewNotImplementedCallable(2)),
		RightSub:      NewFunction(result.symbols, NewNotImplementedCallable(2)),
		Mul:           NewFunction(result.symbols, NewNotImplementedCallable(2)),
		RightMul:      NewFunction(result.symbols, NewNotImplementedCallable(2)),
		Div:           NewFunction(result.symbols, NewNotImplementedCallable(2)),
		RightDiv:      NewFunction(result.symbols, NewNotImplementedCallable(2)),
		Mod:           NewFunction(result.symbols, NewNotImplementedCallable(2)),
		RightMod:      NewFunction(result.symbols, NewNotImplementedCallable(2)),
		Pow:           NewFunction(result.symbols, NewNotImplementedCallable(2)),
		RightPow:      NewFunction(result.symbols, NewNotImplementedCallable(2)),
		BitXor:        NewFunction(result.symbols, NewNotImplementedCallable(2)),
		RightBitXor:   NewFunction(result.symbols, NewNotImplementedCallable(2)),
		BitAnd:        NewFunction(result.symbols, NewNotImplementedCallable(2)),
		RightBitAnd:   NewFunction(result.symbols, NewNotImplementedCallable(2)),
		BitOr:         NewFunction(result.symbols, NewNotImplementedCallable(2)),
		RightBitOr:    NewFunction(result.symbols, NewNotImplementedCallable(2)),
		BitLeft:       NewFunction(result.symbols, NewNotImplementedCallable(2)),
		RightBitLeft:  NewFunction(result.symbols, NewNotImplementedCallable(2)),
		BitRight:      NewFunction(result.symbols, NewNotImplementedCallable(2)),
		RightBitRight: NewFunction(result.symbols, NewNotImplementedCallable(2)),
		//// Logical binary
		And:      NewFunction(result.symbols, NewBuiltInClassFunction(result, 1, ObjAnd)),
		RightAnd: NewFunction(result.symbols, NewBuiltInClassFunction(result, 1, ObjRightAnd)),
		Or:       NewFunction(result.symbols, NewBuiltInClassFunction(result, 1, ObjOr)),
		RightOr:  NewFunction(result.symbols, NewBuiltInClassFunction(result, 1, ObjRightOr)),
		Xor:      NewFunction(result.symbols, NewBuiltInClassFunction(result, 1, ObjXor)),
		RightXor: NewFunction(result.symbols, NewBuiltInClassFunction(result, 1, ObjRightXor)),
		//// Comparison binary
		Equals:                  NewFunction(result.symbols, NewBuiltInClassFunction(result, 1, ObjEquals)),
		RightEquals:             NewFunction(result.symbols, NewBuiltInClassFunction(result, 1, ObjRightEquals)),
		NotEquals:               NewFunction(result.symbols, NewBuiltInClassFunction(result, 1, ObjNotEquals)),
		RightNotEquals:          NewFunction(result.symbols, NewBuiltInClassFunction(result, 1, ObjRightNotEquals)),
		GreaterThan:             NewFunction(result.symbols, NewNotImplementedCallable(2)),
		RightGreaterThan:        NewFunction(result.symbols, NewNotImplementedCallable(2)),
		LessThan:                NewFunction(result.symbols, NewNotImplementedCallable(2)),
		RightLessThan:           NewFunction(result.symbols, NewNotImplementedCallable(2)),
		GreaterThanOrEqual:      NewFunction(result.symbols, NewNotImplementedCallable(2)),
		RightGreaterThanOrEqual: NewFunction(result.symbols, NewNotImplementedCallable(2)),
		LessThanOrEqual:         NewFunction(result.symbols, NewNotImplementedCallable(2)),
		RightLessThanOrEqual:    NewFunction(result.symbols, NewNotImplementedCallable(2)),
		// Behavior
		Assign:     NewFunction(result.symbols, NewNotImplementedCallable(3)),
		Copy:       NewFunction(result.symbols, NewNotImplementedCallable(1)),
		Dir:        NewFunction(result.symbols, NewNotImplementedCallable(1)),
		Index:      NewFunction(result.symbols, NewNotImplementedCallable(1)),
		Call:       NewFunction(result.symbols, NewNotImplementedCallable(1)),
		Iter:       NewFunction(result.symbols, NewNotImplementedCallable(1)),
		Class:      NewFunction(result.symbols, NewNotImplementedCallable(1)),
		SubClasses: NewFunction(result.symbols, NewNotImplementedCallable(1)),
		// Transformation
		ToInteger: NewFunction(result.symbols, NewNotImplementedCallable(3)),
		ToFloat:   NewFunction(result.symbols, NewNotImplementedCallable(3)),
		ToString:  NewFunction(result.symbols, NewBuiltInClassFunction(result, 0, ObjToString)),
		ToBool:    NewFunction(result.symbols, NewBuiltInClassFunction(result, 0, ObjToBool)),
		ToArray:   NewFunction(result.symbols, NewNotImplementedCallable(3)),
		ToTuple:   NewFunction(result.symbols, NewNotImplementedCallable(3)),
	}
	return result
}

const StringName = "String"

type String struct {
	*Object
	Value string
}

func StringAdd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	other := arguments[1]
	if _, ok := other.(*String); !ok {
		return nil, errors.NewTypeError(other.TypeName(), StringName)
	}
	return NewString(
		self.SymbolTable().Parent,
		self.(*String).Value+other.(*String).Value,
	), nil
}

func StringRightAdd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	other := arguments[1]
	if _, ok := other.(*String); !ok {
		return nil, errors.NewTypeError(other.TypeName(), StringName)
	}
	return NewString(
		self.SymbolTable().Parent,
		other.(*String).Value+self.(*String).Value,
	), nil
}

func StringCopy(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewString(
		self.SymbolTable().Parent,
		self.(*String).Value,
	), nil
}

func StringToBool(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewBool(self.SymbolTable().Parent, len(self.(*String).Value) != 0), nil
}

func StringEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	if _, ok := right.(*String); !ok {
		return NewBool(self.SymbolTable().Parent, false), nil
	}
	return NewBool(self.SymbolTable().Parent, self.(*String).Value == right.(*String).Value), nil
}

func StringRightEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	if _, ok := left.(*String); !ok {
		return NewBool(self.SymbolTable().Parent, false), nil
	}
	return NewBool(self.SymbolTable().Parent, left.(*String).Value == self.(*String).Value), nil
}

func StringNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	if _, ok := right.(*String); !ok {
		return NewBool(self.SymbolTable().Parent, false), nil
	}
	return NewBool(self.SymbolTable().Parent, self.(*String).Value != right.(*String).Value), nil
}

func StringRightNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	if _, ok := left.(*String); !ok {
		return NewBool(self.SymbolTable().Parent, false), nil
	}
	return NewBool(self.SymbolTable().Parent, left.(*String).Value != self.(*String).Value), nil
}

/*
	Methods for Strings
	Supported Binary Operations
	Add Strings only - (Done)
	Mul String with integer
	Class
	SubClasses
	Iter
	Assign
	Equals - (Done)
	NotEquals - (Done)
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
		Object: NewObject(StringName, nil, parentSymbols),
		Value:  value,
	}
	string_.Set(Add, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 1, StringAdd)))
	string_.Set(RightAdd, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 1, StringRightAdd)))
	string_.Set(ToString, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringCopy)))
	string_.Set(Copy, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringCopy)))
	string_.Set(ToBool, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringToBool)))
	string_.Set(Equals, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringEquals)))
	string_.Set(RightEquals, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringRightEquals)))
	string_.Set(NotEquals, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringNotEquals)))
	string_.Set(RightNotEquals, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringRightNotEquals)))
	// string_.Set(ToInteger, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringToInteger)))
	// string_.Set(ToFloat, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringToFloat)))
	return string_
}

const IntegerName = "Integer"

type Integer struct {
	*Object
	Value int64
}

func IntegerToBool(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewBool(vm.PeekSymbolTable(), self.(*Integer).Value != 0), nil
}

func IntegerCopy(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewInteger(vm.PeekSymbolTable(), self.(*Integer).Value), nil
}

func IntegerToString(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewString(vm.PeekSymbolTable(), fmt.Sprint(self.(*Integer).Value)), nil
}

func IntegerToFloat(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewFloat(vm.PeekSymbolTable(), float64(self.(*Integer).Value)), nil
}

func IntegerNegBits(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewInteger(vm.PeekSymbolTable(), ^self.(Integer).Value), nil
}

func IntegerAdd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	switch right.(type) {
	case *Integer:
		return NewInteger(vm.PeekSymbolTable(), self.(*Integer).Value+right.(*Integer).Value), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), float64(self.(*Integer).Value)+right.(*Float).Value), nil
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
}

func IntegerRightAdd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	switch left.(type) {
	case *Integer:
		return NewInteger(vm.PeekSymbolTable(), left.(*Integer).Value+self.(*Integer).Value), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), left.(*Float).Value+float64(self.(*Integer).Value)), nil
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
}

func IntegerSub(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	switch right.(type) {
	case *Integer:
		return NewInteger(vm.PeekSymbolTable(), self.(*Integer).Value-right.(*Integer).Value), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), float64(self.(*Integer).Value)-right.(*Float).Value), nil
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
}

func IntegerRightSub(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	switch left.(type) {
	case *Integer:
		return NewInteger(vm.PeekSymbolTable(), left.(*Integer).Value-self.(*Integer).Value), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), left.(*Float).Value-float64(self.(*Integer).Value)), nil
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
}

func IntegerMul(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	switch right.(type) {
	case *Integer:
		return NewInteger(vm.PeekSymbolTable(), self.(*Integer).Value*right.(*Integer).Value), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), float64(self.(*Integer).Value)*right.(*Float).Value), nil
	case *String:
		return NewString(vm.PeekSymbolTable(), Repeat(right.(*String).Value, self.(*Integer).Value)), nil
	case *Tuple:
		panic(NewNotImplementedCallable(errors.UnknownLine))
	case *Array:
		panic(NewNotImplementedCallable(errors.UnknownLine))
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName)
	}
}

func IntegerRightMul(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	switch left.(type) {
	case *Integer:
		return NewInteger(vm.PeekSymbolTable(), left.(*Integer).Value*self.(*Integer).Value), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), left.(*Float).Value*float64(self.(*Integer).Value)), nil
	case *String:
		return NewString(vm.PeekSymbolTable(), Repeat(left.(*String).Value, self.(*Integer).Value)), nil
	case *Tuple:
		panic(NewNotImplementedCallable(errors.UnknownLine))
	case *Array:
		panic(NewNotImplementedCallable(errors.UnknownLine))
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName)
	}
}

func IntegerDiv(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	switch right.(type) {
	case *Integer:
		return NewInteger(vm.PeekSymbolTable(), self.(*Integer).Value/right.(*Integer).Value), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), float64(self.(*Integer).Value)/right.(*Float).Value), nil
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
}

func IntegerRightDiv(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	switch left.(type) {
	case *Integer:
		return NewInteger(vm.PeekSymbolTable(), left.(*Integer).Value/self.(*Integer).Value), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), left.(*Float).Value/float64(self.(*Integer).Value)), nil
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
}

func IntegerMod(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	switch right.(type) {
	case *Integer:
		return NewInteger(vm.PeekSymbolTable(), self.(*Integer).Value%right.(*Integer).Value), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), math.Mod(float64(self.(*Integer).Value), right.(*Float).Value)), nil
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
}

func IntegerRightMod(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	switch left.(type) {
	case *Integer:
		return NewInteger(vm.PeekSymbolTable(), left.(*Integer).Value%self.(*Integer).Value), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), math.Mod(left.(*Float).Value, float64(self.(*Integer).Value))), nil
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
}

func IntegerPow(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	switch right.(type) {
	case *Integer:
		return NewFloat(vm.PeekSymbolTable(), math.Pow(float64(self.(*Integer).Value), float64(right.(*Integer).Value))), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), math.Pow(float64(self.(*Integer).Value), right.(*Float).Value)), nil
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
}

func IntegerRightPow(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	switch left.(type) {
	case *Integer:
		return NewFloat(vm.PeekSymbolTable(), math.Pow(float64(left.(*Integer).Value), float64(self.(*Integer).Value))), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), math.Pow(left.(*Float).Value, float64(self.(*Integer).Value))), nil
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
}

func IntegerBitXor(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	if _, ok := right.(*Integer); !ok {
		return nil, errors.NewTypeError(right.TypeName(), IntegerName)
	}
	return NewInteger(vm.PeekSymbolTable(), self.(*Integer).Value^right.(*Integer).Value), nil
}

func IntegerRightBitXor(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	if _, ok := left.(*Integer); !ok {
		return nil, errors.NewTypeError(left.TypeName(), IntegerName)
	}
	return NewInteger(vm.PeekSymbolTable(), left.(*Integer).Value^self.(*Integer).Value), nil
}

func IntegerBitAnd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	if _, ok := right.(*Integer); !ok {
		return nil, errors.NewTypeError(right.TypeName(), IntegerName)
	}
	return NewInteger(vm.PeekSymbolTable(), self.(*Integer).Value&right.(*Integer).Value), nil
}

func IntegerRightBitAnd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	if _, ok := left.(*Integer); !ok {
		return nil, errors.NewTypeError(left.TypeName(), IntegerName)
	}
	return NewInteger(vm.PeekSymbolTable(), left.(*Integer).Value&self.(*Integer).Value), nil
}

func IntegerBitOr(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	if _, ok := right.(*Integer); !ok {
		return nil, errors.NewTypeError(right.TypeName(), IntegerName)
	}
	return NewInteger(vm.PeekSymbolTable(), self.(*Integer).Value|right.(*Integer).Value), nil
}

func IntegerRightBitOr(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	if _, ok := left.(*Integer); !ok {
		return nil, errors.NewTypeError(left.TypeName(), IntegerName)
	}
	return NewInteger(vm.PeekSymbolTable(), left.(*Integer).Value|self.(*Integer).Value), nil
}

func IntegerBitLeft(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	if _, ok := right.(*Integer); !ok {
		return nil, errors.NewTypeError(right.TypeName(), IntegerName)
	}
	return NewInteger(vm.PeekSymbolTable(), self.(*Integer).Value<<right.(*Integer).Value), nil
}

func IntegerRightBitLeft(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	if _, ok := left.(*Integer); !ok {
		return nil, errors.NewTypeError(left.TypeName(), IntegerName)
	}
	return NewInteger(vm.PeekSymbolTable(), left.(*Integer).Value<<self.(*Integer).Value), nil
}

func IntegerBitRight(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	if _, ok := right.(*Integer); !ok {
		return nil, errors.NewTypeError(right.TypeName(), IntegerName)
	}
	return NewInteger(vm.PeekSymbolTable(), self.(*Integer).Value>>right.(*Integer).Value), nil
}

func IntegerRightBitRight(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	if _, ok := left.(*Integer); !ok {
		return nil, errors.NewTypeError(left.TypeName(), IntegerName)
	}
	return NewInteger(vm.PeekSymbolTable(), left.(*Integer).Value>>self.(*Integer).Value), nil
}

func IntegerEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	var floatRight float64
	switch right.(type) {
	case *Integer:
		floatRight = float64(right.(*Integer).Value)
	case *Float:
		floatRight = right.(*Float).Value
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), float64(self.(*Integer).Value) == floatRight), nil
}
func IntegerRightEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	var floatLeft float64
	switch left.(type) {
	case *Integer:
		floatLeft = float64(left.(*Integer).Value)
	case *Float:
		floatLeft = left.(*Float).Value
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft == float64(self.(*Integer).Value)), nil
}
func IntegerNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	var floatRight float64
	switch right.(type) {
	case *Integer:
		floatRight = float64(right.(*Integer).Value)
	case *Float:
		floatRight = right.(*Float).Value
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), float64(self.(*Integer).Value) != floatRight), nil
}
func IntegerRightNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	var floatLeft float64
	switch left.(type) {
	case *Integer:
		floatLeft = float64(left.(*Integer).Value)
	case *Float:
		floatLeft = left.(*Float).Value
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft != float64(self.(*Integer).Value)), nil
}
func IntegerGreaterThan(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	var floatRight float64
	switch right.(type) {
	case *Integer:
		floatRight = float64(right.(*Integer).Value)
	case *Float:
		floatRight = right.(*Float).Value
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), float64(self.(*Integer).Value) > floatRight), nil
}
func IntegerRightGreaterThan(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	var floatLeft float64
	switch left.(type) {
	case *Integer:
		floatLeft = float64(left.(*Integer).Value)
	case *Float:
		floatLeft = left.(*Float).Value
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft > float64(self.(*Integer).Value)), nil
}
func IntegerLessThan(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	var floatRight float64
	switch right.(type) {
	case *Integer:
		floatRight = float64(right.(*Integer).Value)
	case *Float:
		floatRight = right.(*Float).Value
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), float64(self.(*Integer).Value) < floatRight), nil
}
func IntegerRightLessThan(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	var floatLeft float64
	switch left.(type) {
	case *Integer:
		floatLeft = float64(left.(*Integer).Value)
	case *Float:
		floatLeft = left.(*Float).Value
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft < float64(self.(*Integer).Value)), nil
}
func IntegerGreaterThanOrEqual(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	var floatRight float64
	switch right.(type) {
	case *Integer:
		floatRight = float64(right.(*Integer).Value)
	case *Float:
		floatRight = right.(*Float).Value
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), float64(self.(*Integer).Value) >= floatRight), nil
}
func IntegerRightGreaterThanOrEqual(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	var floatLeft float64
	switch left.(type) {
	case *Integer:
		floatLeft = float64(left.(*Integer).Value)
	case *Float:
		floatLeft = left.(*Float).Value
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft >= float64(self.(*Integer).Value)), nil
}
func IntegerLessThanOrEqual(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	var floatRight float64
	switch right.(type) {
	case *Integer:
		floatRight = float64(right.(*Integer).Value)
	case *Float:
		floatRight = right.(*Float).Value
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), float64(self.(*Integer).Value) <= floatRight), nil
}
func IntegerRightLessThanOrEqual(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	var floatLeft float64
	switch left.(type) {
	case *Integer:
		floatLeft = float64(left.(*Integer).Value)
	case *Float:
		floatLeft = left.(*Float).Value
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft <= float64(self.(*Integer).Value)), nil
}

/*
	Supported methods
	// Unary Operations
	NegBits - (Done)
	// Binary Operations
	//// Basic Binary
	Add          Integer or Float                          - (Done)
	RightAdd     Integer or Float                          - (Done)
	Sub          Integer or Float                          - (Done)
	RightSub     Integer or Float                          - (Done)
	Mul          String, Array, Tuple, Integer or Float
	RightMul     String, Array, Tuple, Integer or Float
	Div          Integer or Float                          - (Done)
	RightDiv     Integer or Float                          - (Done)
	Mod          Integer or Float                          - (Done)
	RightMod     Integer or Float                          - (Done)
	Pow          Integer or Float                          - (Done)
	RightPow     Integer or Float                          - (Done)
	//// Bitwise Binary
	BitXor          - (Done)
	RightBitXor     - (Done)
	BitAnd          - (Done)
	RightBitAnd     - (Done)
	BitOr           - (Done)
	RightBitOr      - (Done)
	BitLeft         - (Done)
	RightBitLeft    - (Done)
	BitRight        - (Done)
	RightBitRight   - (Done)
	//// Comparison Binary
	Equals                        Integer or Float - (Done)
	RightEquals                   Integer or Float - (Done)
	NotEquals                     Integer or Float - (Done)
	RightNotEquals                Integer or Float - (Done)
	GreaterThan                   Integer or Float - (Done)
	RightGreaterThan              Integer or Float - (Done)
	LessThan                      Integer or Float - (Done)
	RightLessThan                 Integer or Float - (Done)
	GreaterThanOrEqual            Integer or Float - (Done)
	RightGreaterThanOrEqual       Integer or Float - (Done)
	LessThanOrEqual               Integer or Float - (Done)
	RightLessThanOrEqual          Integer or Float - (Done)
	// Behavior
	Copy - (Done)
	// Transformation
	ToInteger - (Done)
	ToFloat - (Done)
	ToString - (Done)
	ToBool - (Done)
*/
func NewInteger(parentSymbols *SymbolTable, value int64) *Integer {
	integer := &Integer{
		NewObject(IntegerName, nil, parentSymbols),
		value,
	}

	integer.Set(NegBits, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerNegBits)))

	integer.Set(Add, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerAdd)))
	integer.Set(RightAdd, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerRightAdd)))
	integer.Set(Sub, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerSub)))
	integer.Set(RightSub, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerRightSub)))
	integer.Set(Mul, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerMul)))
	integer.Set(RightMul, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerRightMul)))
	integer.Set(Div, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerDiv)))
	integer.Set(RightDiv, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerRightDiv)))
	integer.Set(Mod, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerMod)))
	integer.Set(RightMod, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerRightMod)))
	integer.Set(Pow, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerPow)))
	integer.Set(RightPow, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerRightPow)))

	integer.Set(BitXor, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerBitXor)))
	integer.Set(RightBitXor, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerRightBitXor)))
	integer.Set(BitAnd, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerBitAnd)))
	integer.Set(RightBitAnd, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerRightBitAnd)))
	integer.Set(BitOr, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerBitOr)))
	integer.Set(RightBitOr, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerRightBitOr)))
	integer.Set(BitLeft, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerBitLeft)))
	integer.Set(RightBitLeft, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerRightBitLeft)))
	integer.Set(BitRight, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerBitRight)))
	integer.Set(RightBitRight, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerRightBitRight)))

	integer.Set(Equals, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerEquals)))
	integer.Set(RightEquals, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerRightEquals)))
	integer.Set(NotEquals, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerNotEquals)))
	integer.Set(RightNotEquals, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerRightNotEquals)))
	integer.Set(GreaterThan, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerGreaterThan)))
	integer.Set(RightGreaterThan, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerRightGreaterThan)))
	integer.Set(LessThan, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerLessThan)))
	integer.Set(RightLessThan, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerRightLessThan)))
	integer.Set(GreaterThanOrEqual, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerGreaterThanOrEqual)))
	integer.Set(RightGreaterThanOrEqual, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerRightGreaterThanOrEqual)))
	integer.Set(LessThanOrEqual, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerLessThanOrEqual)))
	integer.Set(RightLessThanOrEqual, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerRightLessThanOrEqual)))

	integer.Set(Copy, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerCopy)))

	integer.Set(ToInteger, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerCopy)))
	integer.Set(ToFloat, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerToFloat)))
	integer.Set(ToString, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerToString)))
	integer.Set(ToBool, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerToBool)))
	return integer
}

const FloatName = "Float"

type Float struct {
	*Object
	Value float64
}

func FloatToBool(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewBool(vm.PeekSymbolTable(), self.(*Float).Value != 0), nil
}

func FloatCopy(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewFloat(vm.PeekSymbolTable(), self.(*Float).Value), nil
}

func FloatToString(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewString(vm.PeekSymbolTable(), fmt.Sprint(self.(*Float).Value)), nil
}

func FloatToInteger(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewInteger(vm.PeekSymbolTable(), int64(self.(*Float).Value)), nil
}

func FloatAdd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	switch right.(type) {
	case *Integer:
		return NewFloat(vm.PeekSymbolTable(), self.(*Float).Value+float64(right.(*Integer).Value)), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), self.(*Float).Value+right.(*Float).Value), nil
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
}

func FloatRightAdd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	switch left.(type) {
	case *Integer:
		return NewFloat(vm.PeekSymbolTable(), float64(left.(*Integer).Value)+self.(*Float).Value), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), left.(*Float).Value+self.(*Float).Value), nil
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
}

func FloatSub(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	switch right.(type) {
	case *Integer:
		return NewFloat(vm.PeekSymbolTable(), self.(*Float).Value-float64(right.(*Integer).Value)), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), self.(*Float).Value-right.(*Float).Value), nil
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
}

func FloatRightSub(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	switch left.(type) {
	case *Integer:
		return NewFloat(vm.PeekSymbolTable(), float64(left.(*Integer).Value)-self.(*Float).Value), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), left.(*Float).Value-self.(*Float).Value), nil
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
}

func FloatMul(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	switch right.(type) {
	case *Integer:
		return NewFloat(vm.PeekSymbolTable(), self.(*Float).Value*float64(right.(*Integer).Value)), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), self.(*Float).Value*right.(*Float).Value), nil
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName)
	}
}

func FloatRightMul(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	switch left.(type) {
	case *Integer:
		return NewFloat(vm.PeekSymbolTable(), float64(left.(*Integer).Value)*self.(*Float).Value), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), left.(*Float).Value*self.(*Float).Value), nil
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
}

func FloatDiv(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	switch right.(type) {
	case *Integer:
		return NewFloat(vm.PeekSymbolTable(), self.(*Float).Value/float64(right.(*Integer).Value)), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), self.(*Float).Value/right.(*Float).Value), nil
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
}

func FloatRightDiv(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	switch left.(type) {
	case *Integer:
		return NewFloat(vm.PeekSymbolTable(), float64(left.(*Integer).Value)/self.(*Float).Value), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), left.(*Float).Value/self.(*Float).Value), nil
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
}

func FloatMod(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	switch right.(type) {
	case *Integer:
		return NewFloat(vm.PeekSymbolTable(), math.Mod(self.(*Float).Value, float64(right.(*Integer).Value))), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), math.Mod(self.(*Float).Value, right.(*Float).Value)), nil
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
}

func FloatRightMod(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	switch left.(type) {
	case *Integer:
		return NewFloat(vm.PeekSymbolTable(), math.Mod(float64(left.(*Integer).Value), self.(*Float).Value)), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), math.Mod(left.(*Float).Value, self.(*Float).Value)), nil
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
}

func FloatPow(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	switch right.(type) {
	case *Integer:
		return NewFloat(vm.PeekSymbolTable(), math.Pow(self.(*Float).Value, float64(right.(*Integer).Value))), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), math.Pow(self.(*Float).Value, right.(*Float).Value)), nil
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
}

func FloatRightPow(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	switch left.(type) {
	case *Integer:
		return NewFloat(vm.PeekSymbolTable(), math.Pow(float64(left.(*Integer).Value), self.(*Float).Value)), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), math.Pow(left.(*Float).Value, self.(*Float).Value)), nil
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
}

func FloatEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	var floatRight float64
	switch right.(type) {
	case *Integer:
		floatRight = float64(right.(*Integer).Value)
	case *Float:
		floatRight = right.(*Float).Value
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), self.(*Float).Value == floatRight), nil
}

func FloatRightEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	var floatLeft float64
	switch left.(type) {
	case *Integer:
		floatLeft = float64(left.(*Integer).Value)
	case *Float:
		floatLeft = left.(*Float).Value
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft == self.(*Float).Value), nil
}

func FloatNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	var floatRight float64
	switch right.(type) {
	case *Integer:
		floatRight = float64(right.(*Integer).Value)
	case *Float:
		floatRight = right.(*Float).Value
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), self.(*Float).Value != floatRight), nil
}

func FloatRightNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	var floatLeft float64
	switch left.(type) {
	case *Integer:
		floatLeft = float64(left.(*Integer).Value)
	case *Float:
		floatLeft = left.(*Float).Value
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft != self.(*Float).Value), nil
}

func FloatGreaterThan(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	var floatRight float64
	switch right.(type) {
	case *Integer:
		floatRight = float64(right.(*Integer).Value)
	case *Float:
		floatRight = right.(*Float).Value
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), self.(*Float).Value > floatRight), nil
}

func FloatRightGreaterThan(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	var floatLeft float64
	switch left.(type) {
	case *Integer:
		floatLeft = float64(left.(*Integer).Value)
	case *Float:
		floatLeft = left.(*Float).Value
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft > self.(*Float).Value), nil
}

func FloatLessThan(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	var floatRight float64
	switch right.(type) {
	case *Integer:
		floatRight = float64(right.(*Integer).Value)
	case *Float:
		floatRight = right.(*Float).Value
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), self.(*Float).Value < floatRight), nil
}
func FloatRightLessThan(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	var floatLeft float64
	switch left.(type) {
	case *Integer:
		floatLeft = float64(left.(*Integer).Value)
	case *Float:
		floatLeft = left.(*Float).Value
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft < self.(*Float).Value), nil
}
func FloatGreaterThanOrEqual(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	var floatRight float64
	switch right.(type) {
	case *Integer:
		floatRight = float64(right.(*Integer).Value)
	case *Float:
		floatRight = right.(*Float).Value
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), self.(*Float).Value >= floatRight), nil
}
func FloatRightGreaterThanOrEqual(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	var floatLeft float64
	switch left.(type) {
	case *Integer:
		floatLeft = float64(left.(*Integer).Value)
	case *Float:
		floatLeft = left.(*Float).Value
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft >= self.(*Float).Value), nil
}
func FloatLessThanOrEqual(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	var floatRight float64
	switch right.(type) {
	case *Integer:
		floatRight = float64(right.(*Integer).Value)
	case *Float:
		floatRight = right.(*Float).Value
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), self.(*Float).Value <= floatRight), nil
}
func FloatRightLessThanOrEqual(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	var floatLeft float64
	switch left.(type) {
	case *Integer:
		floatLeft = float64(left.(*Integer).Value)
	case *Float:
		floatLeft = left.(*Float).Value
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft <= self.(*Float).Value), nil
}

/*
	Supported methods
	// Binary Operations
	//// Basic Binary
	Add          Integer or Float  - (Done)
	RightAdd     Integer or Float  - (Done)
	Sub          Integer or Float  - (Done)
	RightSub     Integer or Float  - (Done)
	Mul          Integer or Float  - (Done)
	RightMul     Integer or Float  - (Done)
	Div          Integer or Float  - (Done)
	RightDiv     Integer or Float  - (Done)
	Mod          Integer or Float  - (Done)
	RightMod     Integer or Float  - (Done)
	Pow          Integer or Float  - (Done)
	RightPow     Integer or Float  - (Done)
	//// Comparison Binary
	Equals                        Integer or Float    - (Done)
	RightEquals                   Integer or Float    - (Done)
	NotEquals                     Integer or Float    - (Done)
	RightNotEquals                Integer or Float    - (Done)
	GreaterThan                   Integer or Float    - (Done)
	RightGreaterThan              Integer or Float    - (Done)
	LessThan                      Integer or Float    - (Done)
	RightLessThan                 Integer or Float    - (Done)
	GreaterThanOrEqual            Integer or Float    - (Done)
	RightGreaterThanOrEqual       Integer or Float    - (Done)
	LessThanOrEqual               Integer or Float    - (Done)
	RightLessThanOrEqual          Integer or Float    - (Done)
	// Behavior
	Copy          - (Done)
	// Transformation
	ToInteger    - (Done)
	ToFloat      - (Done)
	ToString     - (Done)
	ToBool       - (Done)
*/
func NewFloat(parentSymbols *SymbolTable, value float64) *Float {
	float_ := &Float{
		NewObject(IntegerName, nil, parentSymbols),
		value,
	}

	float_.Set(Add, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatAdd)))
	float_.Set(RightAdd, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatRightAdd)))
	float_.Set(Sub, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatSub)))
	float_.Set(RightSub, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatRightSub)))
	float_.Set(Mul, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatMul)))
	float_.Set(RightMul, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatRightMul)))
	float_.Set(Div, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatDiv)))
	float_.Set(RightDiv, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatRightDiv)))
	float_.Set(Mod, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatMod)))
	float_.Set(RightMod, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatRightMod)))
	float_.Set(Pow, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatPow)))
	float_.Set(RightPow, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatRightPow)))

	float_.Set(Equals, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatEquals)))
	float_.Set(RightEquals, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatRightEquals)))
	float_.Set(NotEquals, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatNotEquals)))
	float_.Set(RightNotEquals, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatRightNotEquals)))
	float_.Set(GreaterThan, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatGreaterThan)))
	float_.Set(RightGreaterThan, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatRightGreaterThan)))
	float_.Set(LessThan, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatLessThan)))
	float_.Set(RightLessThan, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatRightLessThan)))
	float_.Set(GreaterThanOrEqual, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatGreaterThanOrEqual)))
	float_.Set(RightGreaterThanOrEqual, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatRightGreaterThanOrEqual)))
	float_.Set(LessThanOrEqual, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatLessThanOrEqual)))
	float_.Set(RightLessThanOrEqual, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatRightLessThanOrEqual)))

	float_.Set(Copy, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatCopy)))

	float_.Set(ToInteger, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatToInteger)))
	float_.Set(ToFloat, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatCopy)))
	float_.Set(ToString, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatToString)))
	float_.Set(ToBool, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatToBool)))
	return float_
}

const ArrayName = "Array"

type Array struct {
	*Object
	Content []IObject
	Length  int
}

const TupleName = "Tuple"

type Tuple struct {
	*Object
	Content []IObject
	Length  int
}

const BoolName = "Bool"

const TrueName = "True"
const FalseName = "False"

type Bool struct {
	*Object
	Value bool
}

func BoolToString(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
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
	bool_.Set(ToString, NewFunction(bool_.symbols, NewBuiltInClassFunction(bool_, 0, BoolToString)))
	return bool_
}

func SetDefaultSymbolTable() *SymbolTable {
	symbolTable := NewSymbolTable(nil)
	// String
	symbolTable.Set(StringName,
		NewFunction(symbolTable,
			NewBuiltInFunction(1,
				func(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
					toString, getError := arguments[0].Get(ToString)
					if getError != nil {
						return nil, getError
					}
					if _, ok := toString.(*Function); !ok {
						return nil, errors.NewTypeError(toString.(IObject).TypeName(), FunctionName)
					}
					return CallFunction(toString.(*Function), vm, arguments[0].SymbolTable().Parent, arguments[0])
				},
			),
		),
	)
	return symbolTable
}
