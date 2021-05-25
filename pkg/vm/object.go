package vm

import (
	"encoding/binary"
	"fmt"
	"github.com/shoriwe/gruby/pkg/errors"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const (
	TypeName     = "Type"
	ObjectName   = "Object"
	FunctionName = "Function"
	StringName   = "String"
	BoolName     = "Bool"
	TrueName     = "True"
	FalseName    = "False"
	TupleName    = "Tuple"
	IntegerName  = "Integer"
	FloatName    = "Float"
	ArrayName    = "Array"
	NoneName     = "NoneType"
	BytesName    = "Bytes"
	HashName     = "Hash"
	IterName     = "Iterator"
	ModuleName   = "Module"
	None         = "None"
)

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
	Hash       = "Hash"
	Copy       = "Copy"
	Index      = "Index"
	Assign     = "Assign"
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

var (
	integerPattern = regexp.MustCompile("(?m)^-?\\d+(_|\\d)*$")
	hexPattern     = regexp.MustCompile("(?m)^0[xX][a-zA-Z0-9]+(_|[a-zA-Z0-9])*$")
	octalPattern   = regexp.MustCompile("(?m)^0[oO][0-7]+(_|[0-7])*$")
	binaryPattern  = regexp.MustCompile("(?m)^0[bB][01]+(_|[01])*$")
)

var (
	floatPattern          = regexp.MustCompile("(?m)^-?\\d+(_|\\d)*\\.\\d+(_|\\d)*$")
	scientificPattern     = regexp.MustCompile("(?m)^-?((\\d+(_|\\d)*)|(\\d+(_|\\d)*\\.\\d+(_|\\d)*))[eE][-+]\\d+(_|\\d)*$")
	numberPartPattern     = regexp.MustCompile("(?m)^-?((\\d+(_|\\d)*)|(\\d+(_|\\d)*\\.\\d+(_|\\d)*))[eE][-+]")
	scientificPartPattern = regexp.MustCompile("(?m)[eE][-+]\\d+(_|\\d)*$")
)

func ParseInteger(s string) (string, int, bool) {
	if integerPattern.MatchString(s) {
		return strings.ReplaceAll(s, "_", ""), 10, true
	} else if hexPattern.MatchString(s) {
		return strings.ReplaceAll(s[2:], "_", ""), 8, true
	} else if octalPattern.MatchString(s) {
		return strings.ReplaceAll(s[2:], "_", ""), 16, true
	} else if binaryPattern.MatchString(s) {
		return strings.ReplaceAll(s[2:], "_", ""), 2, true
	}
	return "", 0, false
}

func ParseFloat(s string) (string, bool) {
	if floatPattern.MatchString(s) {
		return strings.ReplaceAll(s, "_", ""), true
	} else if scientificPattern.MatchString(s) {
		// ToDo: Support scientific numbers
	}
	return "", false
}

func Repeat(s string, times int64) string {
	result := ""
	var i int64
	for i = 0; i < times; i++ {
		result += s
	}
	return result
}

func CalcIndex(object IObject, length int) (int, *errors.Error) {
	index := int(object.(*Integer).Value)
	if length <= index {
		return 0, errors.NewIndexOutOfRange(errors.UnknownLine, length, index)
	}
	if index < 0 {
		index = length + index
		if index < 0 {
			return 0, errors.NewIndexOutOfRange(errors.UnknownLine, length, index)
		}
	}
	return index, nil
}

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

type Constructor interface {
	Initialize(VirtualMachine, IObject) *errors.Error
}

type PlasmaConstructor struct {
	Constructor
	Code []Code
}

func (c *PlasmaConstructor) Initialize(vm VirtualMachine, object IObject) *errors.Error {
	vm.LoadCode(c.Code)
	vm.PushSymbolTable(object.SymbolTable())
	_, executionError := vm.Execute()
	return executionError
}

func NewPlasmaConstructor(code []Code) *PlasmaConstructor {
	return &PlasmaConstructor{
		Code: code,
	}
}

type ConstructorCallBack func(VirtualMachine, IObject) *errors.Error

type BuiltInConstructor struct {
	Constructor
	callback ConstructorCallBack
}

func (c *BuiltInConstructor) Initialize(vm VirtualMachine, object IObject) *errors.Error {
	return c.callback(vm, object)
}

func NewBuiltInConstructor(callback ConstructorCallBack) *BuiltInConstructor {
	return &BuiltInConstructor{
		callback: callback,
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
	SubClasses() []*Type
	Get(string) (IObject, *errors.Error)
	Set(string, IObject)
	GetHash() int64
	SetHash(int64)
}

// MetaClass for IObject
type Object struct {
	id             uint
	typeName       string
	subClasses     []*Type
	symbols        *SymbolTable
	virtualMachine VirtualMachine
	hash           int64
}

func (o *Object) Id() uint {
	return o.id
}

func (o *Object) SubClasses() []*Type {
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

func (o *Object) GetHash() int64 {
	return o.hash
}

func (o *Object) SetHash(newHash int64) {
	o.hash = newHash
}

func ObjInitialize(_ VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	return nil, nil
}

func CallFunction(function *Function, vm VirtualMachine, parent *SymbolTable, arguments ...IObject) (IObject, *errors.Error) {
	if function.Callable.NumberOfArguments() != len(arguments) {
		return nil, errors.NewInvalidNumberOfArguments(len(arguments), function.Callable.NumberOfArguments())
	}
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

func ConstructObject(type_ *Type, vm VirtualMachine, parent *SymbolTable) (IObject, *errors.Error) {
	object := NewObject(type_.typeName, type_.subClasses, parent)
	for _, subclass := range object.subClasses {
		initializationError := subclass.Constructor.Initialize(vm, object)
		if initializationError != nil {
			return nil, initializationError
		}
	}
	return object, nil
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
	right := arguments[1]
	var rightToBool interface{}
	rightToBool, foundError = right.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := rightToBool.(*Function); !ok {
		return nil, errors.NewTypeError(rightToBool.(IObject).TypeName(), FunctionName)
	}
	var rightBool IObject
	rightBool, transformationError = CallFunction(rightToBool.(*Function), vm, right.SymbolTable(), right)
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(vm.PeekSymbolTable(), leftBool.(*Bool).Value && rightBool.(*Bool).Value), nil
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
	left := arguments[1]
	var leftToBool interface{}
	leftToBool, foundError = left.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := leftToBool.(*Function); !ok {
		return nil, errors.NewTypeError(leftToBool.(IObject).TypeName(), FunctionName)
	}
	var leftBool IObject
	leftBool, transformationError = CallFunction(leftToBool.(*Function), vm, left.SymbolTable(), left)
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(vm.PeekSymbolTable(), leftBool.(*Bool).Value && rightBool.(*Bool).Value), nil
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

	right := arguments[1]
	var rightToBool interface{}
	rightToBool, foundError = right.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := rightToBool.(*Function); !ok {
		return nil, errors.NewTypeError(rightToBool.(IObject).TypeName(), FunctionName)
	}
	var rightBool IObject
	rightBool, transformationError = CallFunction(rightToBool.(*Function), vm, right.SymbolTable(), right)
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(vm.PeekSymbolTable(), leftBool.(*Bool).Value || rightBool.(*Bool).Value), nil
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
	left := arguments[0]
	var leftToBool interface{}
	leftToBool, foundError = left.Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := leftToBool.(*Function); !ok {
		return nil, errors.NewTypeError(leftToBool.(IObject).TypeName(), FunctionName)
	}
	var leftBool IObject
	leftBool, transformationError = CallFunction(leftToBool.(*Function), vm, left.SymbolTable(), left)
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(vm.PeekSymbolTable(), leftBool.(*Bool).Value || rightBool.(*Bool).Value), nil
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

	right := arguments[1]
	var rightToBool interface{}
	rightToBool, foundError = arguments[1].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := rightToBool.(*Function); !ok {
		return nil, errors.NewTypeError(rightToBool.(IObject).TypeName(), FunctionName)
	}
	var rightBool IObject
	rightBool, transformationError = CallFunction(rightToBool.(*Function), vm, right.SymbolTable(), right)
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(vm.PeekSymbolTable(), leftBool.(*Bool).Value != rightBool.(*Bool).Value), nil
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

	left := arguments[0]
	var rightToBool interface{}
	rightToBool, foundError = arguments[1].Get(ToBool)
	if foundError != nil {
		return nil, foundError
	}
	if _, ok := rightToBool.(*Function); !ok {
		return nil, errors.NewTypeError(rightToBool.(IObject).TypeName(), FunctionName)
	}
	var rightBool IObject
	rightBool, transformationError = CallFunction(rightToBool.(*Function), vm, left.SymbolTable(), left)
	if transformationError != nil {
		return nil, transformationError
	}
	return NewBool(vm.PeekSymbolTable(), rightBool.(*Bool).Value != leftBool.(*Bool).Value), nil
}

func ObjEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[1]
	return NewBool(vm.PeekSymbolTable(), self.Id() == right.Id()), nil
}

func ObjRightEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	return NewBool(vm.PeekSymbolTable(), left.Id() == self.Id()), nil
}

func ObjNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[1]
	return NewBool(vm.PeekSymbolTable(), self.Id() != right.Id()), nil
}

func ObjRightNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	return NewBool(vm.PeekSymbolTable(), left.Id() != self.Id()), nil
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
	return NewBool(vm.PeekSymbolTable(), !selfBool.(*Bool).Value), nil
}

func ObjToBool(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	return NewBool(vm.PeekSymbolTable(), true), nil
}

func ObjToString(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewString(vm.PeekSymbolTable(),
		fmt.Sprintf("%s-%d", self.TypeName(), self.Id())), nil
}

func ObjHash(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	if self.GetHash() == 0 {
		objectHash, hashingError := vm.HashString(fmt.Sprintf("%v-%s-%d", self, self.TypeName(), self.Id()))
		if hashingError != nil {
			return nil, hashingError
		}
		self.SetHash(objectHash)

	}
	return NewInteger(vm.PeekSymbolTable(), self.GetHash()), nil
}

/*
	Negate
	//// Logical Binary
	And           - (Done)
	RightAnd      - (Done)
	Or            - (Done)
	RightOr       - (Done)
	Xor           - (Done)
	RightXor      - (Done)
	//// Comparison Binary
	Equals            - (Done)
	RightEquals       - (Done)
	NotEquals         - (Done)
	RightNotEquals    - (Done)
	// Behavior
	Hash - (Done)
	Copy - (Done)
	Dir
	Index - (Done)
	Call - (Done)
	Iter - (Done)
	Class
	SubClasses
	// Transformation
	ToString - (Done)
	ToBool - (Done)
*/
func ObjectInitialize(_ VirtualMachine, object IObject) *errors.Error {
	object.SymbolTable().Update(map[string]IObject{
		// IObject Creation
		Initialize: NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, ObjInitialize)),
		// Unary Operations
		NegBits: NewFunction(object.SymbolTable(), NewNotImplementedCallable(0)),
		Negate:  NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, ObjNegate)),
		// Binary Operations
		//// Math binary
		Add:           NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		RightAdd:      NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		Sub:           NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		RightSub:      NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		Mul:           NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		RightMul:      NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		Div:           NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		RightDiv:      NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		Mod:           NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		RightMod:      NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		Pow:           NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		RightPow:      NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		BitXor:        NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		RightBitXor:   NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		BitAnd:        NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		RightBitAnd:   NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		BitOr:         NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		RightBitOr:    NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		BitLeft:       NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		RightBitLeft:  NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		BitRight:      NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		RightBitRight: NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		//// Logical binary
		And:      NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ObjAnd)),
		RightAnd: NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ObjRightAnd)),
		Or:       NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ObjOr)),
		RightOr:  NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ObjRightOr)),
		Xor:      NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ObjXor)),
		RightXor: NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ObjRightXor)),
		//// Comparison binary
		Equals:                  NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ObjEquals)),
		RightEquals:             NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ObjRightEquals)),
		NotEquals:               NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ObjNotEquals)),
		RightNotEquals:          NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ObjRightNotEquals)),
		GreaterThan:             NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		RightGreaterThan:        NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		LessThan:                NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		RightLessThan:           NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		GreaterThanOrEqual:      NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		RightGreaterThanOrEqual: NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		LessThanOrEqual:         NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		RightLessThanOrEqual:    NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		// Behavior
		Hash:       NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, ObjHash)),
		Copy:       NewFunction(object.SymbolTable(), NewNotImplementedCallable(0)),
		Index:      NewFunction(object.SymbolTable(), NewNotImplementedCallable(1)),
		Assign:     NewFunction(object.SymbolTable(), NewNotImplementedCallable(2)),
		Call:       NewFunction(object.SymbolTable(), NewNotImplementedCallable(0)),
		Iter:       NewFunction(object.SymbolTable(), NewNotImplementedCallable(0)),
		Class:      NewFunction(object.SymbolTable(), NewNotImplementedCallable(0)),
		SubClasses: NewFunction(object.SymbolTable(), NewNotImplementedCallable(0)),
		// Transformation
		ToInteger: NewFunction(object.SymbolTable(), NewNotImplementedCallable(0)),
		ToFloat:   NewFunction(object.SymbolTable(), NewNotImplementedCallable(0)),
		ToString:  NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, ObjToString)),
		ToBool:    NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, ObjToBool)),
		ToArray:   NewFunction(object.SymbolTable(), NewNotImplementedCallable(0)),
		ToTuple:   NewFunction(object.SymbolTable(), NewNotImplementedCallable(0)),
	})
	return nil
}

// NewObject Creates an Empty Object
func NewObject(
	typeName string,
	subClasses []*Type,
	parentSymbols *SymbolTable,
) *Object {
	result := &Object{
		id:         counter.NextId(),
		typeName:   typeName,
		subClasses: subClasses,
		symbols:    NewSymbolTable(parentSymbols),
	}
	ObjectInitialize(nil, result)
	return result
}

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

type String struct {
	*Object
	Value  string
	Length int
}

func StringAdd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[1]
	if _, ok := right.(*String); !ok {
		return nil, errors.NewTypeError(right.TypeName(), StringName)
	}
	return NewString(
		vm.PeekSymbolTable(),
		self.(*String).Value+right.(*String).Value,
	), nil
}

func StringRightAdd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[1]
	if _, ok := left.(*String); !ok {
		return nil, errors.NewTypeError(left.TypeName(), StringName)
	}
	return NewString(
		vm.PeekSymbolTable(),
		left.(*String).Value+self.(*String).Value,
	), nil
}

func StringMul(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[1]
	if _, ok := right.(*Integer); !ok {
		return nil, errors.NewTypeError(right.TypeName(), IntegerName)
	}
	return NewString(
		vm.PeekSymbolTable(),
		Repeat(self.(*String).Value, right.(*Integer).Value),
	), nil
}

func StringRightMul(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[1]
	if _, ok := left.(*Integer); !ok {
		return nil, errors.NewTypeError(left.TypeName(), IntegerName)
	}
	return NewString(
		vm.PeekSymbolTable(),
		Repeat(self.(*String).Value, left.(*Integer).Value),
	), nil
}

func StringCopy(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewString(
		vm.PeekSymbolTable(),
		self.(*String).Value,
	), nil
}

func StringToBool(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewBool(vm.PeekSymbolTable(), self.(*String).Length != 0), nil
}

func StringEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	if _, ok := right.(*String); !ok {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	return NewBool(vm.PeekSymbolTable(), self.(*String).Value == right.(*String).Value), nil
}

func StringRightEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	if _, ok := left.(*String); !ok {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	return NewBool(vm.PeekSymbolTable(), left.(*String).Value == self.(*String).Value), nil
}

func StringNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	if _, ok := right.(*String); !ok {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	return NewBool(vm.PeekSymbolTable(), self.(*String).Value != right.(*String).Value), nil
}

func StringRightNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	if _, ok := left.(*String); !ok {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	return NewBool(vm.PeekSymbolTable(), left.(*String).Value != self.(*String).Value), nil
}

func StringIndex(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	indexObject := arguments[0]
	if _, ok := indexObject.(*Integer); ok {
		index, getIndexError := CalcIndex(indexObject, self.(*String).Length)
		if getIndexError != nil {
			return nil, getIndexError
		}
		return NewString(vm.PeekSymbolTable(), string(self.(*String).Value[index])), nil
	} else if _, ok = indexObject.(*Tuple); ok {
		if len(indexObject.(*Tuple).Content) != 2 {
			return nil, errors.NewInvalidNumberOfArguments(len(indexObject.(*Tuple).Content), 2)
		}
		startIndex, calcError := CalcIndex(indexObject.(*Tuple).Content[0], self.(*Array).Length)
		if calcError != nil {
			return nil, calcError
		}
		targetIndex, calcError := CalcIndex(indexObject.(*Tuple).Content[1], self.(*Array).Length)
		if calcError != nil {
			return nil, calcError
		}
		return NewString(vm.PeekSymbolTable(), self.(*String).Value[startIndex:targetIndex]), nil
	} else {
		return nil, errors.NewTypeError(indexObject.TypeName(), IntegerName, TupleName)
	}
}

func StringToInteger(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	content, base, valid := ParseInteger(self.(*String).Value)
	if !valid {
		return nil, errors.NewInvalidIntegerDefinition(errors.UnknownLine, self.(*String).Value)
	}
	number, parsingError := strconv.ParseInt(content, base, 64)
	if parsingError != nil {
		return nil, errors.NewInvalidIntegerDefinition(errors.UnknownLine, self.(*String).Value)
	}
	return NewInteger(vm.PeekSymbolTable(), number), nil
}

func StringToFloat(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	content, valid := ParseFloat(self.(*String).Value)
	if !valid {
		return nil, errors.NewInvalidFloatDefinition(errors.UnknownLine, self.(*String).Value)
	}
	number, parsingError := strconv.ParseFloat(content, 64)
	if parsingError != nil {
		return nil, errors.NewInvalidFloatDefinition(errors.UnknownLine, self.(*String).Value)
	}
	return NewFloat(vm.PeekSymbolTable(), number), nil
}

func StringToArray(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	var content []IObject
	for _, char := range self.(*String).Value {
		content = append(content, NewString(
			vm.PeekSymbolTable(), string(char),
		),
		)
	}
	return NewArray(vm.PeekSymbolTable(), content), nil
}

func StringToTuple(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	var content []IObject
	for _, char := range self.(*String).Value {
		content = append(content, NewString(
			vm.PeekSymbolTable(), string(char),
		),
		)
	}
	return NewTuple(vm.PeekSymbolTable(), content), nil
}

func StringHash(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	if self.GetHash() == 0 {
		stringHash, hashingError := vm.HashString(fmt.Sprintf("%s-%s", self.(*String).Value, StringName))
		if hashingError != nil {
			return nil, hashingError
		}
		self.SetHash(stringHash)
	}
	return NewInteger(vm.PeekSymbolTable(), self.GetHash()), nil
}

/*
	// Binary Operations
	//// Basic Binary
	Add         String only - (Done)
	RightAdd    String only - (Done)
	Mul         String with Integer only      - (Done)
	RightMul    String with Integer only      - (Done)
	//// Comparison Binary
	Equals              String only - (Done)
	RightEquals         String only - (Done)
	NotEquals           String only - (Done)
	RightNotEquals      String only - (Done)
	// Behavior
	Hash			   - (Done)
	Copy               - (Done)
	Index    Integer or Tuple
	Iter
	// Transformation
	ToInteger
	ToFloat
	ToString       - (Done)
	ToBool         - (Done)
	ToArray        - (Done)
	ToTuple        - (Done)
*/
func StringInitialize(_ VirtualMachine, object IObject) *errors.Error {
	object.Set(Add, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, StringAdd)))
	object.Set(RightAdd, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, StringRightAdd)))
	object.Set(Mul, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, StringMul)))
	object.Set(RightMul, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, StringRightMul)))
	object.Set(Equals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, StringEquals)))
	object.Set(RightEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, StringRightEquals)))
	object.Set(NotEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, StringNotEquals)))
	object.Set(RightNotEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, StringRightNotEquals)))
	object.Set(Hash, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, StringHash)))
	object.Set(Copy, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, StringCopy)))
	object.Set(Index, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, StringIndex)))
	// object.Set(Iter, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, StringIter)))
	object.Set(ToInteger, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, StringToInteger)))
	object.Set(ToFloat, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, StringToFloat)))
	object.Set(ToString, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, StringCopy)))
	object.Set(ToBool, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, StringToBool)))
	object.Set(ToArray, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, StringToArray)))
	object.Set(ToTuple, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, StringToTuple)))
	return nil
}

func NewString(
	parentSymbols *SymbolTable,
	value string,
) *String {
	string_ := &String{
		Object: NewObject(StringName, nil, parentSymbols),
		Value:  value,
		Length: len(value),
	}
	StringInitialize(nil, string_)
	return string_
}

func BytesAdd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	if _, ok := right.(*Bytes); !ok {
		return nil, errors.NewTypeError(right.TypeName(), BytesName)
	}
	var newContent []*Integer
	var byteCopyFunc IObject
	for _, byte_ := range self.(*Bytes).Content {
		byteCopyFunc, getError = byte_.Get(Copy)
		if getError != nil {
			return nil, getError
		}
		if _, ok := byteCopyFunc.(*Function); !ok {
			return nil, errors.NewTypeError(byteCopyFunc.TypeName(), FunctionName)
		}
		byteCopy, callError := CallFunction(byteCopyFunc.(*Function), vm, vm.PeekSymbolTable())
		if callError != nil {
			return nil, callError
		}
		if _, ok := byteCopy.(*Integer); !ok {
			return nil, errors.NewTypeError(byteCopy.TypeName(), IntegerName)
		}
		newContent = append(newContent, byteCopy.(*Integer))
	}
	for _, byte_ := range right.(*Bytes).Content {
		byteCopyFunc, getError = byte_.Get(Copy)
		if getError != nil {
			return nil, getError
		}
		if _, ok := byteCopyFunc.(*Function); !ok {
			return nil, errors.NewTypeError(byteCopyFunc.TypeName(), FunctionName)
		}
		byteCopy, callError := CallFunction(byteCopyFunc.(*Function), vm, vm.PeekSymbolTable())
		if callError != nil {
			return nil, callError
		}
		if _, ok := byteCopy.(*Integer); !ok {
			return nil, errors.NewTypeError(byteCopy.TypeName(), IntegerName)
		}
		newContent = append(newContent, byteCopy.(*Integer))
	}
	return NewBytes(vm.PeekSymbolTable(), newContent), nil
}

func BytesRightAdd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	if _, ok := left.(*Bytes); !ok {
		return nil, errors.NewTypeError(left.TypeName(), BytesName)
	}
	var newContent []*Integer
	var byteCopyFunc IObject
	for _, byte_ := range left.(*Bytes).Content {
		byteCopyFunc, getError = byte_.Get(Copy)
		if getError != nil {
			return nil, getError
		}
		if _, ok := byteCopyFunc.(*Function); !ok {
			return nil, errors.NewTypeError(byteCopyFunc.TypeName(), FunctionName)
		}
		byteCopy, callError := CallFunction(byteCopyFunc.(*Function), vm, vm.PeekSymbolTable())
		if callError != nil {
			return nil, callError
		}
		if _, ok := byteCopy.(*Integer); !ok {
			return nil, errors.NewTypeError(byteCopy.TypeName(), IntegerName)
		}
		newContent = append(newContent, byteCopy.(*Integer))
	}
	for _, byte_ := range self.(*Bytes).Content {
		byteCopyFunc, getError = byte_.Get(Copy)
		if getError != nil {
			return nil, getError
		}
		if _, ok := byteCopyFunc.(*Function); !ok {
			return nil, errors.NewTypeError(byteCopyFunc.TypeName(), FunctionName)
		}
		byteCopy, callError := CallFunction(byteCopyFunc.(*Function), vm, vm.PeekSymbolTable())
		if callError != nil {
			return nil, callError
		}
		if _, ok := byteCopy.(*Integer); !ok {
			return nil, errors.NewTypeError(byteCopy.TypeName(), IntegerName)
		}
		newContent = append(newContent, byteCopy.(*Integer))
	}
	return NewBytes(vm.PeekSymbolTable(), newContent), nil
}

func BytesMul(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	if _, ok := right.(*Integer); !ok {
		return nil, errors.NewTypeError(right.TypeName(), IntegerName)
	}
	var i int64
	var newContent []*Integer
	var byteCopyFunc IObject
	for i = 0; i < right.(*Integer).Value; i++ {
		for _, byte_ := range self.(*Bytes).Content {
			byteCopyFunc, getError = byte_.Get(Copy)
			if getError != nil {
				return nil, getError
			}
			if _, ok := byteCopyFunc.(*Function); !ok {
				return nil, errors.NewTypeError(byteCopyFunc.TypeName(), FunctionName)
			}
			byteCopy, callError := CallFunction(byteCopyFunc.(*Function), vm, vm.PeekSymbolTable())
			if callError != nil {
				return nil, callError
			}
			if _, ok := byteCopy.(*Integer); !ok {
				return nil, errors.NewTypeError(byteCopy.TypeName(), IntegerName)
			}
			newContent = append(newContent, byteCopy.(*Integer))
		}
	}
	return NewBytes(vm.PeekSymbolTable(), newContent), nil
}

func BytesRightMul(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	if _, ok := left.(*Integer); !ok {
		return nil, errors.NewTypeError(left.TypeName(), IntegerName)
	}
	var i int64
	var newContent []*Integer
	var byteCopyFunc IObject
	for i = 0; i < left.(*Integer).Value; i++ {
		for _, byte_ := range self.(*Bytes).Content {
			byteCopyFunc, getError = byte_.Get(Copy)
			if getError != nil {
				return nil, getError
			}
			if _, ok := byteCopyFunc.(*Function); !ok {
				return nil, errors.NewTypeError(byteCopyFunc.TypeName(), FunctionName)
			}
			byteCopy, callError := CallFunction(byteCopyFunc.(*Function), vm, vm.PeekSymbolTable())
			if callError != nil {
				return nil, callError
			}
			if _, ok := byteCopy.(*Integer); !ok {
				return nil, errors.NewTypeError(byteCopy.TypeName(), IntegerName)
			}
			newContent = append(newContent, byteCopy.(*Integer))
		}
	}
	return NewBytes(vm.PeekSymbolTable(), newContent), nil
}

func BytesEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	if _, ok := right.(*Bytes); !ok {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	if self.(*Bytes).Length != right.(*Bytes).Length {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftEquals IObject
	for index, byte_ := range self.(*Bytes).Content {
		leftEquals, getError = byte_.Get(Equals)
		if getError != nil {
			return nil, getError
		}
		if _, ok := leftEquals.(*Function); !ok {
			return nil, errors.NewTypeError(leftEquals.TypeName(), FunctionName)
		}
		result, callError := CallFunction(leftEquals.(*Function), vm, vm.PeekSymbolTable(), right.(*Bytes).Content[index])
		if callError != nil {
			return nil, callError
		}
		if _, ok := result.(*Bool); !ok {
			return nil, errors.NewTypeError(result.TypeName(), BoolName)
		}
		if !result.(*Bool).Value {
			return NewBool(vm.PeekSymbolTable(), false), nil
		}
	}
	return NewBool(vm.PeekSymbolTable(), true), nil
}

func BytesRightEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	if _, ok := left.(*Bytes); !ok {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	if left.(*Bytes).Length != self.(*Bytes).Length {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftEquals IObject
	for index, byte_ := range left.(*Bytes).Content {
		leftEquals, getError = byte_.Get(RightEquals)
		if getError != nil {
			return nil, getError
		}
		if _, ok := leftEquals.(*Function); !ok {
			return nil, errors.NewTypeError(leftEquals.TypeName(), FunctionName)
		}
		result, callError := CallFunction(leftEquals.(*Function), vm, vm.PeekSymbolTable(), self.(*Bytes).Content[index])
		if callError != nil {
			return nil, callError
		}
		if _, ok := result.(*Bool); !ok {
			return nil, errors.NewTypeError(result.TypeName(), BoolName)
		}
		if !result.(*Bool).Value {
			return NewBool(vm.PeekSymbolTable(), false), nil
		}
	}
	return NewBool(vm.PeekSymbolTable(), true), nil
}

func BytesNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	if _, ok := right.(*Bytes); !ok {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	if self.(*Bytes).Length != right.(*Bytes).Length {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftEquals IObject
	for index, byte_ := range self.(*Bytes).Content {
		leftEquals, getError = byte_.Get(NotEquals)
		if getError != nil {
			return nil, getError
		}
		if _, ok := leftEquals.(*Function); !ok {
			return nil, errors.NewTypeError(leftEquals.TypeName(), FunctionName)
		}
		result, callError := CallFunction(leftEquals.(*Function), vm, vm.PeekSymbolTable(), right.(*Bytes).Content[index])
		if callError != nil {
			return nil, callError
		}
		if _, ok := result.(*Bool); !ok {
			return nil, errors.NewTypeError(result.TypeName(), BoolName)
		}
		if !result.(*Bool).Value {
			return NewBool(vm.PeekSymbolTable(), false), nil
		}
	}
	return NewBool(vm.PeekSymbolTable(), true), nil
}

func BytesRightNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	if _, ok := left.(*Bytes); !ok {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	if left.(*Bytes).Length != self.(*Bytes).Length {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftEquals IObject
	for index, byte_ := range left.(*Bytes).Content {
		leftEquals, getError = byte_.Get(RightNotEquals)
		if getError != nil {
			return nil, getError
		}
		if _, ok := leftEquals.(*Function); !ok {
			return nil, errors.NewTypeError(leftEquals.TypeName(), FunctionName)
		}
		result, callError := CallFunction(leftEquals.(*Function), vm, vm.PeekSymbolTable(), self.(*Bytes).Content[index])
		if callError != nil {
			return nil, callError
		}
		if _, ok := result.(*Bool); !ok {
			return nil, errors.NewTypeError(result.TypeName(), BoolName)
		}
		if !result.(*Bool).Value {
			return NewBool(vm.PeekSymbolTable(), false), nil
		}
	}
	return NewBool(vm.PeekSymbolTable(), true), nil
}

func BytesHash(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	var hashNumber int64
	for index, byte_ := range self.(*Bytes).Content {
		if index == 0 {
			hashNumber = byte_.hash
		} else {
			hashNumber <<= byte_.hash
		}
	}
	return NewInteger(vm.PeekSymbolTable(), hashNumber), nil
}

func BytesCopy(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	var newBytes []*Integer
	for _, byte_ := range self.(*Bytes).Content {
		newBytes = append(newBytes, NewInteger(vm.PeekSymbolTable(), byte_.Value))
	}
	return NewBytes(vm.PeekSymbolTable(), newBytes), nil
}

func BytesIndex(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	indexObject := arguments[0]
	var ok bool
	if _, ok = indexObject.(*Integer); ok {
		index, calcError := CalcIndex(indexObject, self.(*Bytes).Length)
		if calcError != nil {
			return nil, calcError
		}
		return self.(*Bytes).Content[index], nil
	} else if _, ok = indexObject.(*Tuple); ok {
		if len(indexObject.(*Tuple).Content) != 2 {
			return nil, errors.NewInvalidNumberOfArguments(len(indexObject.(*Tuple).Content), 2)
		}
		startIndex, calcError := CalcIndex(indexObject.(*Tuple).Content[0], self.(*Array).Length)
		if calcError != nil {
			return nil, calcError
		}
		targetIndex, calcError := CalcIndex(indexObject.(*Tuple).Content[1], self.(*Array).Length)
		if calcError != nil {
			return nil, calcError
		}
		return NewBytes(vm.PeekSymbolTable(), self.(*Bytes).Content[startIndex:targetIndex]), nil
	} else {
		return nil, errors.NewTypeError(indexObject.TypeName(), IntegerName, TupleName)
	}
}

func BytesToInteger(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	var numbers []uint8
	for _, byte_ := range self.(*Bytes).Content {
		numbers = append(numbers, uint8(byte_.Value))
	}
	return NewInteger(vm.PeekSymbolTable(), int64(binary.BigEndian.Uint32(numbers))), nil
}

func BytesToString(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	var numbers []uint8
	for _, byte_ := range self.(*Bytes).Content {
		numbers = append(numbers, uint8(byte_.Value))
	}
	return NewString(vm.PeekSymbolTable(), string(numbers)), nil
}

func BytesToBool(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewBool(vm.PeekSymbolTable(), self.(*Bytes).Length != 0), nil
}

func BytesToArray(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	var newContent []IObject
	var byteCopyFunc IObject
	for _, byte_ := range self.(*Bytes).Content {
		byteCopyFunc, getError = byte_.Get(Copy)
		if getError != nil {
			return nil, getError
		}
		if _, ok := byteCopyFunc.(*Function); !ok {
			return nil, errors.NewTypeError(byteCopyFunc.TypeName(), FunctionName)
		}
		byteCopy, callError := CallFunction(byteCopyFunc.(*Function), vm, vm.PeekSymbolTable())
		if callError != nil {
			return nil, callError
		}
		if _, ok := byteCopy.(*Integer); !ok {
			return nil, errors.NewTypeError(byteCopy.TypeName(), IntegerName)
		}
		newContent = append(newContent, byteCopy.(*Integer))
	}
	return NewArray(vm.PeekSymbolTable(), newContent), nil
}

func BytesToTuple(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	var newContent []IObject
	var byteCopyFunc IObject
	for _, byte_ := range self.(*Bytes).Content {
		byteCopyFunc, getError = byte_.Get(Copy)
		if getError != nil {
			return nil, getError
		}
		if _, ok := byteCopyFunc.(*Function); !ok {
			return nil, errors.NewTypeError(byteCopyFunc.TypeName(), FunctionName)
		}
		byteCopy, callError := CallFunction(byteCopyFunc.(*Function), vm, vm.PeekSymbolTable())
		if callError != nil {
			return nil, callError
		}
		if _, ok := byteCopy.(*Integer); !ok {
			return nil, errors.NewTypeError(byteCopy.TypeName(), IntegerName)
		}
		newContent = append(newContent, byteCopy.(*Integer))
	}
	return NewTuple(vm.PeekSymbolTable(), newContent), nil
}

type Bytes struct {
	*Object
	Content []*Integer
	Length  int
}

/*
	// Binary Operations
	//// Basic Binary
	Add           - (Done)
	RightAdd      - (Done)
	Mul           - (Done)
	RightMul      - (Done)
	//// Comparison Binary
	Equals            - (Done)
	RightEquals       - (Done)
	NotEquals         - (Done)
	RightNotEquals    - (Done)
	// Behavior
	Hash            - (Done)
	Copy            - (Done)
	Index           - (Done)
	Iter            - ()
	// Transformation
	ToInteger       - (Done)
	ToString        - (Done)
	ToBool          - (Done)
	ToArray         - (Done)
	ToTuple         - (Done)
*/
func BytesInitialize(_ VirtualMachine, object IObject) *errors.Error {
	object.Set(Add, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, BytesAdd)))
	object.Set(RightAdd, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, BytesRightAdd)))
	object.Set(Mul, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, BytesMul)))
	object.Set(RightMul, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, BytesRightMul)))

	object.Set(Equals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, BytesEquals)))
	object.Set(RightEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, BytesRightEquals)))
	object.Set(NotEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, BytesNotEquals)))
	object.Set(RightNotEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, BytesRightNotEquals)))

	object.Set(Hash, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, BytesHash)))
	object.Set(Copy, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, BytesCopy)))
	object.Set(Index, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, BytesIndex)))

	object.Set(ToInteger, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, BytesToInteger)))
	object.Set(ToString, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, BytesToString)))
	object.Set(ToBool, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, BytesToBool)))
	object.Set(ToArray, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, BytesToArray)))
	object.Set(ToTuple, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, BytesToTuple)))
	return nil
}

func NewBytes(parent *SymbolTable, content []*Integer) IObject {
	bytes_ := &Bytes{
		Object:  NewObject(BytesName, nil, parent),
		Content: content,
		Length:  len(content),
	}
	BytesInitialize(nil, bytes_)
	return bytes_
}

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

func IntegerHash(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	if self.GetHash() == 0 {
		integerHash, hashingError := vm.HashString(fmt.Sprintf("%d-%s", self.(*Integer).hash, IntegerName))
		if hashingError != nil {
			return nil, hashingError
		}
		self.SetHash(integerHash)
	}
	return NewInteger(vm.PeekSymbolTable(), self.GetHash()), nil
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
	Hash - (Done)
	Copy - (Done)
	// Transformation
	ToInteger - (Done)
	ToFloat - (Done)
	ToString - (Done)
	ToBool - (Done)
*/
func IntegerInitialize(_ VirtualMachine, object IObject) *errors.Error {
	object.SymbolTable().Set(NegBits, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerNegBits)))

	object.SymbolTable().Set(Add, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerAdd)))
	object.SymbolTable().Set(RightAdd, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerRightAdd)))
	object.SymbolTable().Set(Sub, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerSub)))
	object.SymbolTable().Set(RightSub, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerRightSub)))
	object.SymbolTable().Set(Mul, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerMul)))
	object.SymbolTable().Set(RightMul, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerRightMul)))
	object.SymbolTable().Set(Div, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerDiv)))
	object.SymbolTable().Set(RightDiv, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerRightDiv)))
	object.SymbolTable().Set(Mod, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerMod)))
	object.SymbolTable().Set(RightMod, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerRightMod)))
	object.SymbolTable().Set(Pow, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerPow)))
	object.SymbolTable().Set(RightPow, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerRightPow)))

	object.SymbolTable().Set(BitXor, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerBitXor)))
	object.SymbolTable().Set(RightBitXor, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerRightBitXor)))
	object.SymbolTable().Set(BitAnd, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerBitAnd)))
	object.SymbolTable().Set(RightBitAnd, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerRightBitAnd)))
	object.SymbolTable().Set(BitOr, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerBitOr)))
	object.SymbolTable().Set(RightBitOr, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerRightBitOr)))
	object.SymbolTable().Set(BitLeft, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerBitLeft)))
	object.SymbolTable().Set(RightBitLeft, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerRightBitLeft)))
	object.SymbolTable().Set(BitRight, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerBitRight)))
	object.SymbolTable().Set(RightBitRight, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerRightBitRight)))

	object.SymbolTable().Set(Equals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerEquals)))
	object.SymbolTable().Set(RightEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerRightEquals)))
	object.SymbolTable().Set(NotEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerNotEquals)))
	object.SymbolTable().Set(RightNotEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerRightNotEquals)))
	object.SymbolTable().Set(GreaterThan, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerGreaterThan)))
	object.SymbolTable().Set(RightGreaterThan, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerRightGreaterThan)))
	object.SymbolTable().Set(LessThan, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerLessThan)))
	object.SymbolTable().Set(RightLessThan, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerRightLessThan)))
	object.SymbolTable().Set(GreaterThanOrEqual, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerGreaterThanOrEqual)))
	object.SymbolTable().Set(RightGreaterThanOrEqual, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerRightGreaterThanOrEqual)))
	object.SymbolTable().Set(LessThanOrEqual, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerLessThanOrEqual)))
	object.SymbolTable().Set(RightLessThanOrEqual, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerRightLessThanOrEqual)))

	object.SymbolTable().Set(Hash, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerHash)))
	object.SymbolTable().Set(Copy, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerCopy)))

	object.SymbolTable().Set(ToInteger, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerCopy)))
	object.SymbolTable().Set(ToFloat, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerToFloat)))
	object.SymbolTable().Set(ToString, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerToString)))
	object.SymbolTable().Set(ToBool, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, IntegerToBool)))
	return nil
}

func NewInteger(parentSymbols *SymbolTable, value int64) *Integer {
	integer := &Integer{
		NewObject(IntegerName, nil, parentSymbols),
		value,
	}
	IntegerInitialize(nil, integer)
	return integer
}

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

func FloatHash(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	if self.GetHash() == 0 {
		floatHash, hashingError := vm.HashString(fmt.Sprintf("%f-%s", self.(*Float).Value, FloatName))
		if hashingError != nil {
			return nil, hashingError
		}
		self.SetHash(floatHash)
	}
	return NewInteger(vm.PeekSymbolTable(), self.GetHash()), nil
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
	Hash          - (Done)
	Copy          - (Done)
	// Transformation
	ToInteger    - (Done)
	ToFloat      - (Done)
	ToString     - (Done)
	ToBool       - (Done)
*/
func FloatInitialize(_ VirtualMachine, object IObject) *errors.Error {
	object.Set(Add, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatAdd)))
	object.Set(RightAdd, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatRightAdd)))
	object.Set(Sub, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatSub)))
	object.Set(RightSub, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatRightSub)))
	object.Set(Mul, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatMul)))
	object.Set(RightMul, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatRightMul)))
	object.Set(Div, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatDiv)))
	object.Set(RightDiv, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatRightDiv)))
	object.Set(Mod, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatMod)))
	object.Set(RightMod, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatRightMod)))
	object.Set(Pow, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatPow)))
	object.Set(RightPow, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatRightPow)))

	object.Set(Equals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatEquals)))
	object.Set(RightEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatRightEquals)))
	object.Set(NotEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatNotEquals)))
	object.Set(RightNotEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatRightNotEquals)))
	object.Set(GreaterThan, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatGreaterThan)))
	object.Set(RightGreaterThan, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatRightGreaterThan)))
	object.Set(LessThan, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatLessThan)))
	object.Set(RightLessThan, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatRightLessThan)))
	object.Set(GreaterThanOrEqual, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatGreaterThanOrEqual)))
	object.Set(RightGreaterThanOrEqual, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatRightGreaterThanOrEqual)))
	object.Set(LessThanOrEqual, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatLessThanOrEqual)))
	object.Set(RightLessThanOrEqual, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatRightLessThanOrEqual)))

	object.Set(Hash, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatHash)))
	object.Set(Copy, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatCopy)))

	object.Set(ToInteger, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatToInteger)))
	object.Set(ToFloat, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatCopy)))
	object.Set(ToString, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatToString)))
	object.Set(ToBool, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, FloatToBool)))
	return nil
}

func NewFloat(parentSymbols *SymbolTable, value float64) *Float {
	float_ := &Float{
		NewObject(IntegerName, nil, parentSymbols),
		value,
	}
	FloatInitialize(nil, float_)
	return float_
}

type Array struct {
	*Object
	Content []IObject
	Length  int
}

func ArrayEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	if _, ok := right.(*Array); !ok {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	if self.(*Array).Length != right.(*Array).Length {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftEquals IObject
	var rightEquals IObject
	var comparisonResult IObject
	var callError *errors.Error
	var comparisonResultToBool IObject
	var comparisonBool IObject

	for i := 0; i < self.(*Array).Length; i++ {
		leftEquals, getError = self.(*Array).Content[i].Get(Equals)
		if getError != nil {
			rightEquals, getError = right.(*Array).Content[i].Get(RightEquals)
			if getError != nil {
				return nil, getError
			}
			if _, ok := rightEquals.(*Function); !ok {
				return nil, errors.NewTypeError(rightEquals.TypeName(), FunctionName)
			}
			comparisonResult, callError = CallFunction(rightEquals.(*Function), vm, vm.PeekSymbolTable(), self.(*Array).Content[i])
		} else {
			comparisonResult, callError = CallFunction(leftEquals.(*Function), vm, vm.PeekSymbolTable(), right.(*Array).Content[i])
		}
		if callError != nil {
			return nil, callError
		}
		comparisonResultToBool, getError = comparisonResult.Get(ToBool)
		if getError != nil {
			return nil, getError
		}
		if _, ok := comparisonResultToBool.(*Function); !ok {
			return nil, errors.NewTypeError(comparisonResultToBool.TypeName(), FunctionName)
		}
		comparisonBool, callError = CallFunction(comparisonResultToBool.(*Function), vm, vm.PeekSymbolTable())
		if callError != nil {
			return nil, callError
		}
		if _, ok := comparisonBool.(*Bool); !ok {
			return nil, errors.NewTypeError(comparisonBool.TypeName(), BoolName)
		}
		if !comparisonBool.(*Bool).Value {
			return NewBool(vm.PeekSymbolTable(), false), nil
		}
	}
	return NewBool(vm.PeekSymbolTable(), true), nil
}

func ArrayRightEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	if _, ok := left.(*Array); !ok {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	if self.(*Array).Length != left.(*Array).Length {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftEquals IObject
	var rightEquals IObject
	var comparisonResult IObject
	var callError *errors.Error
	var comparisonResultToBool IObject
	var comparisonBool IObject

	for i := 0; i < self.(*Array).Length; i++ {
		leftEquals, getError = left.(*Array).Content[i].Get(Equals)
		if getError != nil {
			rightEquals, getError = self.(*Array).Content[i].Get(RightEquals)
			if getError != nil {
				return nil, getError
			}
			if _, ok := rightEquals.(*Function); !ok {
				return nil, errors.NewTypeError(rightEquals.TypeName(), FunctionName)
			}
			comparisonResult, callError = CallFunction(rightEquals.(*Function), vm, vm.PeekSymbolTable(), left.(*Array).Content[i])
		} else {
			comparisonResult, callError = CallFunction(leftEquals.(*Function), vm, vm.PeekSymbolTable(), self.(*Array).Content[i])
		}
		if callError != nil {
			return nil, callError
		}
		comparisonResultToBool, getError = comparisonResult.Get(ToBool)
		if getError != nil {
			return nil, getError
		}
		if _, ok := comparisonResultToBool.(*Function); !ok {
			return nil, errors.NewTypeError(comparisonResultToBool.TypeName(), FunctionName)
		}
		comparisonBool, callError = CallFunction(comparisonResultToBool.(*Function), vm, vm.PeekSymbolTable())
		if callError != nil {
			return nil, callError
		}
		if _, ok := comparisonBool.(*Bool); !ok {
			return nil, errors.NewTypeError(comparisonBool.TypeName(), BoolName)
		}
		if !comparisonBool.(*Bool).Value {
			return NewBool(vm.PeekSymbolTable(), false), nil
		}
	}
	return NewBool(vm.PeekSymbolTable(), true), nil
}

func ArrayNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	if _, ok := right.(*Array); !ok {
		return NewBool(vm.PeekSymbolTable(), true), nil
	}
	if self.(*Array).Length != right.(*Array).Length {
		return NewBool(vm.PeekSymbolTable(), true), nil
	}
	var leftNotEquals IObject
	var rightNotEquals IObject
	var comparisonResult IObject
	var callError *errors.Error
	var comparisonResultToBool IObject
	var comparisonBool IObject

	for i := 0; i < self.(*Array).Length; i++ {
		leftNotEquals, getError = self.(*Array).Content[i].Get(NotEquals)
		if getError != nil {
			rightNotEquals, getError = right.(*Array).Content[i].Get(RightNotEquals)
			if getError != nil {
				return nil, getError
			}
			if _, ok := rightNotEquals.(*Function); !ok {
				return nil, errors.NewTypeError(rightNotEquals.TypeName(), FunctionName)
			}
			comparisonResult, callError = CallFunction(rightNotEquals.(*Function), vm, vm.PeekSymbolTable(), self.(*Array).Content[i])
		} else {
			comparisonResult, callError = CallFunction(leftNotEquals.(*Function), vm, vm.PeekSymbolTable(), right.(*Array).Content[i])
		}
		if callError != nil {
			return nil, callError
		}
		comparisonResultToBool, getError = comparisonResult.Get(ToBool)
		if getError != nil {
			return nil, getError
		}
		if _, ok := comparisonResultToBool.(*Function); !ok {
			return nil, errors.NewTypeError(comparisonResultToBool.TypeName(), FunctionName)
		}
		comparisonBool, callError = CallFunction(comparisonResultToBool.(*Function), vm, vm.PeekSymbolTable())
		if callError != nil {
			return nil, callError
		}
		if _, ok := comparisonBool.(*Bool); !ok {
			return nil, errors.NewTypeError(comparisonBool.TypeName(), BoolName)
		}
		if !comparisonBool.(*Bool).Value {
			return NewBool(vm.PeekSymbolTable(), false), nil
		}
	}
	return NewBool(vm.PeekSymbolTable(), true), nil
}

func ArrayRightNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	if _, ok := left.(*Array); !ok {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	if self.(*Array).Length != left.(*Array).Length {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftEquals IObject
	var rightEquals IObject
	var comparisonResult IObject
	var callError *errors.Error
	var comparisonResultToBool IObject
	var comparisonBool IObject

	for i := 0; i < self.(*Array).Length; i++ {
		leftEquals, getError = left.(*Array).Content[i].Get(NotEquals)
		if getError != nil {
			rightEquals, getError = self.(*Array).Content[i].Get(RightNotEquals)
			if getError != nil {
				return nil, getError
			}
			if _, ok := rightEquals.(*Function); !ok {
				return nil, errors.NewTypeError(rightEquals.TypeName(), FunctionName)
			}
			comparisonResult, callError = CallFunction(rightEquals.(*Function), vm, vm.PeekSymbolTable(), left.(*Array).Content[i])
		} else {
			comparisonResult, callError = CallFunction(leftEquals.(*Function), vm, vm.PeekSymbolTable(), self.(*Array).Content[i])
		}
		if callError != nil {
			return nil, callError
		}
		comparisonResultToBool, getError = comparisonResult.Get(ToBool)
		if getError != nil {
			return nil, getError
		}
		if _, ok := comparisonResultToBool.(*Function); !ok {
			return nil, errors.NewTypeError(comparisonResultToBool.TypeName(), FunctionName)
		}
		comparisonBool, callError = CallFunction(comparisonResultToBool.(*Function), vm, vm.PeekSymbolTable())
		if callError != nil {
			return nil, callError
		}
		if _, ok := comparisonBool.(*Bool); !ok {
			return nil, errors.NewTypeError(comparisonBool.TypeName(), BoolName)
		}
		if !comparisonBool.(*Bool).Value {
			return NewBool(vm.PeekSymbolTable(), false), nil
		}
	}
	return NewBool(vm.PeekSymbolTable(), true), nil
}

func ArrayCopy(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	var copiedObjects []IObject
	var objectCopy IObject
	for _, object := range self.(*Array).Content {
		objectCopy, getError = object.Get(Copy)
		if _, ok := objectCopy.(*Function); !ok {
			return nil, errors.NewTypeError(objectCopy.TypeName(), FunctionName)
		}
		copiedObject, copyError := CallFunction(objectCopy.(*Function), vm, vm.PeekSymbolTable())
		if copyError != nil {
			return nil, copyError
		}
		copiedObjects = append(copiedObjects, copiedObject)
	}
	return NewArray(vm.PeekSymbolTable(), copiedObjects), nil
}

func ArrayIndex(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	indexObject := arguments[0]
	var ok bool
	if _, ok = indexObject.(*Integer); ok {
		index, calcError := CalcIndex(indexObject, self.(*Array).Length)
		if calcError != nil {
			return nil, calcError
		}
		return self.(*Array).Content[index], nil
	} else if _, ok = indexObject.(*Tuple); ok {
		if len(indexObject.(*Tuple).Content) != 2 {
			return nil, errors.NewInvalidNumberOfArguments(len(indexObject.(*Tuple).Content), 2)
		}
		startIndex, calcError := CalcIndex(indexObject.(*Tuple).Content[0], self.(*Array).Length)
		if calcError != nil {
			return nil, calcError
		}
		targetIndex, calcError := CalcIndex(indexObject.(*Tuple).Content[1], self.(*Array).Length)
		if calcError != nil {
			return nil, calcError
		}
		return NewArray(vm.PeekSymbolTable(), self.(*Array).Content[startIndex:targetIndex]), nil
	} else {
		return nil, errors.NewTypeError(indexObject.TypeName(), IntegerName, TupleName)
	}
}

func ArrayAssign(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	index, calcError := CalcIndex(arguments[0], self.(*Array).Length)
	if calcError != nil {
		return nil, calcError
	}
	self.(*Array).Content[index] = arguments[1]
	var none IObject
	none, getError = vm.PeekSymbolTable().GetAny(None)
	if getError != nil {
		return nil, getError
	}
	return none, nil
}

func ArrayToString(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	result := "("
	var objectToString IObject
	var objectString IObject
	var callError *errors.Error
	for index, object := range self.(*Array).Content {
		if index != 0 {
			result += ", "
		}
		objectToString, getError = object.Get(ToString)
		if getError != nil {
			return nil, getError
		}
		if _, ok := objectToString.(*Function); !ok {
			return nil, errors.NewTypeError(objectToString.TypeName(), FunctionName)
		}
		objectString, callError = CallFunction(objectToString.(*Function), vm, vm.PeekSymbolTable())
		if callError != nil {
			return nil, callError
		}
		if _, ok := objectString.(*String); !ok {
			return nil, errors.NewTypeError(objectString.TypeName(), StringName)
		}
		result += objectString.(*String).Value
	}
	return NewString(vm.PeekSymbolTable(), result+")"), nil
}

func ArrayToBool(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewBool(vm.PeekSymbolTable(), self.(*Array).Length != 0), nil
}

func ArrayToArray(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewArray(vm.PeekSymbolTable(), append([]IObject{}, self.(*Array).Content...)), nil
}

func ArrayToTuple(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewTuple(vm.PeekSymbolTable(), append([]IObject{}, self.(*Array).Content...)), nil
}

func ArrayHash(_ VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	return nil, errors.NewUnhashableTypeError(errors.UnknownLine)
}

/*
	// Binary Operations
	//// Comparison Binary
	Equals            - (Done)
	RightEquals       - (Done)
	NotEquals         - (Done)
	RightNotEquals    - (Done)
	// Behavior
	Hash  	  - (Done)
	Copy      - (Done)
	Index     Integer or Tuple
	Assign	  - (Done)
	Iter
	// Transformation
	ToString      - (Done)
	ToBool        - (Done)
	ToArray       - (Done)
	ToTuple       - (Done)
*/
func ArrayInitialize(_ VirtualMachine, object IObject) *errors.Error {
	object.Set(Equals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ArrayEquals)))
	object.Set(RightEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ArrayRightEquals)))
	object.Set(NotEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ArrayNotEquals)))
	object.Set(RightNotEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ArrayRightNotEquals)))
	object.Set(Hash, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, ArrayHash)))
	object.Set(Copy, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, ArrayCopy)))
	object.Set(Index, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ArrayIndex)))
	object.Set(Assign, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 2, ArrayAssign)))
	// object.Set(Iter, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, ArrayIter)))
	object.Set(ToString, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, ArrayToString)))
	object.Set(ToBool, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, ArrayToBool)))
	object.Set(ToArray, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, ArrayToArray)))
	object.Set(ToTuple, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, ArrayToTuple)))
	return nil
}

func NewArray(parentSymbols *SymbolTable, content []IObject) *Array {
	array := &Array{
		Object:  NewObject(ArrayName, nil, parentSymbols),
		Content: content,
		Length:  len(content),
	}
	ArrayInitialize(nil, array)
	return array
}

type Tuple struct {
	*Object
	Content []IObject
	Length  int
}

func TupleEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	if _, ok := right.(*Tuple); !ok {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	if self.(*Tuple).Length != right.(*Tuple).Length {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftEquals IObject
	var rightEquals IObject
	var comparisonResult IObject
	var callError *errors.Error
	var comparisonResultToBool IObject
	var comparisonBool IObject

	for i := 0; i < self.(*Tuple).Length; i++ {
		leftEquals, getError = self.(*Tuple).Content[i].Get(Equals)
		if getError != nil {
			rightEquals, getError = right.(*Tuple).Content[i].Get(RightEquals)
			if getError != nil {
				return nil, getError
			}
			if _, ok := rightEquals.(*Function); !ok {
				return nil, errors.NewTypeError(rightEquals.TypeName(), FunctionName)
			}
			comparisonResult, callError = CallFunction(rightEquals.(*Function), vm, vm.PeekSymbolTable(), self.(*Tuple).Content[i])
		} else {
			comparisonResult, callError = CallFunction(leftEquals.(*Function), vm, vm.PeekSymbolTable(), right.(*Tuple).Content[i])
		}
		if callError != nil {
			return nil, callError
		}
		comparisonResultToBool, getError = comparisonResult.Get(ToBool)
		if getError != nil {
			return nil, getError
		}
		if _, ok := comparisonResultToBool.(*Function); !ok {
			return nil, errors.NewTypeError(comparisonResultToBool.TypeName(), FunctionName)
		}
		comparisonBool, callError = CallFunction(comparisonResultToBool.(*Function), vm, vm.PeekSymbolTable())
		if callError != nil {
			return nil, callError
		}
		if _, ok := comparisonBool.(*Bool); !ok {
			return nil, errors.NewTypeError(comparisonBool.TypeName(), BoolName)
		}
		if !comparisonBool.(*Bool).Value {
			return NewBool(vm.PeekSymbolTable(), false), nil
		}
	}
	return NewBool(vm.PeekSymbolTable(), true), nil
}

func TupleRightEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	if _, ok := left.(*Tuple); !ok {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	if self.(*Tuple).Length != left.(*Tuple).Length {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftEquals IObject
	var rightEquals IObject
	var comparisonResult IObject
	var callError *errors.Error
	var comparisonResultToBool IObject
	var comparisonBool IObject

	for i := 0; i < self.(*Tuple).Length; i++ {
		leftEquals, getError = left.(*Tuple).Content[i].Get(Equals)
		if getError != nil {
			rightEquals, getError = self.(*Tuple).Content[i].Get(RightEquals)
			if getError != nil {
				return nil, getError
			}
			if _, ok := rightEquals.(*Function); !ok {
				return nil, errors.NewTypeError(rightEquals.TypeName(), FunctionName)
			}
			comparisonResult, callError = CallFunction(rightEquals.(*Function), vm, vm.PeekSymbolTable(), left.(*Tuple).Content[i])
		} else {
			comparisonResult, callError = CallFunction(leftEquals.(*Function), vm, vm.PeekSymbolTable(), self.(*Tuple).Content[i])
		}
		if callError != nil {
			return nil, callError
		}
		comparisonResultToBool, getError = comparisonResult.Get(ToBool)
		if getError != nil {
			return nil, getError
		}
		if _, ok := comparisonResultToBool.(*Function); !ok {
			return nil, errors.NewTypeError(comparisonResultToBool.TypeName(), FunctionName)
		}
		comparisonBool, callError = CallFunction(comparisonResultToBool.(*Function), vm, vm.PeekSymbolTable())
		if callError != nil {
			return nil, callError
		}
		if _, ok := comparisonBool.(*Bool); !ok {
			return nil, errors.NewTypeError(comparisonBool.TypeName(), BoolName)
		}
		if !comparisonBool.(*Bool).Value {
			return NewBool(vm.PeekSymbolTable(), false), nil
		}
	}
	return NewBool(vm.PeekSymbolTable(), true), nil
}

func TupleNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	if _, ok := right.(*Tuple); !ok {
		return NewBool(vm.PeekSymbolTable(), true), nil
	}
	if self.(*Tuple).Length != right.(*Tuple).Length {
		return NewBool(vm.PeekSymbolTable(), true), nil
	}
	var leftNotEquals IObject
	var rightNotEquals IObject
	var comparisonResult IObject
	var callError *errors.Error
	var comparisonResultToBool IObject
	var comparisonBool IObject

	for i := 0; i < self.(*Tuple).Length; i++ {
		leftNotEquals, getError = self.(*Tuple).Content[i].Get(NotEquals)
		if getError != nil {
			rightNotEquals, getError = right.(*Tuple).Content[i].Get(RightNotEquals)
			if getError != nil {
				return nil, getError
			}
			if _, ok := rightNotEquals.(*Function); !ok {
				return nil, errors.NewTypeError(rightNotEquals.TypeName(), FunctionName)
			}
			comparisonResult, callError = CallFunction(rightNotEquals.(*Function), vm, vm.PeekSymbolTable(), self.(*Tuple).Content[i])
		} else {
			comparisonResult, callError = CallFunction(leftNotEquals.(*Function), vm, vm.PeekSymbolTable(), right.(*Tuple).Content[i])
		}
		if callError != nil {
			return nil, callError
		}
		comparisonResultToBool, getError = comparisonResult.Get(ToBool)
		if getError != nil {
			return nil, getError
		}
		if _, ok := comparisonResultToBool.(*Function); !ok {
			return nil, errors.NewTypeError(comparisonResultToBool.TypeName(), FunctionName)
		}
		comparisonBool, callError = CallFunction(comparisonResultToBool.(*Function), vm, vm.PeekSymbolTable())
		if callError != nil {
			return nil, callError
		}
		if _, ok := comparisonBool.(*Bool); !ok {
			return nil, errors.NewTypeError(comparisonBool.TypeName(), BoolName)
		}
		if !comparisonBool.(*Bool).Value {
			return NewBool(vm.PeekSymbolTable(), false), nil
		}
	}
	return NewBool(vm.PeekSymbolTable(), true), nil
}

func TupleRightNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	if _, ok := left.(*Tuple); !ok {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	if self.(*Tuple).Length != left.(*Tuple).Length {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftEquals IObject
	var rightEquals IObject
	var comparisonResult IObject
	var callError *errors.Error
	var comparisonResultToBool IObject
	var comparisonBool IObject

	for i := 0; i < self.(*Tuple).Length; i++ {
		leftEquals, getError = left.(*Tuple).Content[i].Get(NotEquals)
		if getError != nil {
			rightEquals, getError = self.(*Tuple).Content[i].Get(RightNotEquals)
			if getError != nil {
				return nil, getError
			}
			if _, ok := rightEquals.(*Function); !ok {
				return nil, errors.NewTypeError(rightEquals.TypeName(), FunctionName)
			}
			comparisonResult, callError = CallFunction(rightEquals.(*Function), vm, vm.PeekSymbolTable(), left.(*Tuple).Content[i])
		} else {
			comparisonResult, callError = CallFunction(leftEquals.(*Function), vm, vm.PeekSymbolTable(), self.(*Tuple).Content[i])
		}
		if callError != nil {
			return nil, callError
		}
		comparisonResultToBool, getError = comparisonResult.Get(ToBool)
		if getError != nil {
			return nil, getError
		}
		if _, ok := comparisonResultToBool.(*Function); !ok {
			return nil, errors.NewTypeError(comparisonResultToBool.TypeName(), FunctionName)
		}
		comparisonBool, callError = CallFunction(comparisonResultToBool.(*Function), vm, vm.PeekSymbolTable())
		if callError != nil {
			return nil, callError
		}
		if _, ok := comparisonBool.(*Bool); !ok {
			return nil, errors.NewTypeError(comparisonBool.TypeName(), BoolName)
		}
		if !comparisonBool.(*Bool).Value {
			return NewBool(vm.PeekSymbolTable(), false), nil
		}
	}
	return NewBool(vm.PeekSymbolTable(), true), nil
}

func TupleCopy(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	var copiedObjects []IObject
	var objectCopy IObject
	for _, object := range self.(*Tuple).Content {
		objectCopy, getError = object.Get(Copy)
		if _, ok := objectCopy.(*Function); !ok {
			return nil, errors.NewTypeError(objectCopy.TypeName(), FunctionName)
		}
		copiedObject, copyError := CallFunction(objectCopy.(*Function), vm, vm.PeekSymbolTable())
		if copyError != nil {
			return nil, copyError
		}
		copiedObjects = append(copiedObjects, copiedObject)
	}
	return NewTuple(vm.PeekSymbolTable(), copiedObjects), nil
}

func TupleIndex(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	indexObject := arguments[0]
	var ok bool
	if _, ok = indexObject.(*Integer); ok {
		index, calcError := CalcIndex(indexObject, self.(*Tuple).Length)
		if calcError != nil {
			return nil, calcError
		}
		return self.(*Tuple).Content[index], nil
	} else if _, ok = indexObject.(*Tuple); ok {
		if len(indexObject.(*Tuple).Content) != 2 {
			return nil, errors.NewInvalidNumberOfArguments(len(indexObject.(*Tuple).Content), 2)
		}
		startIndex, calcError := CalcIndex(indexObject.(*Tuple).Content[0], self.(*Array).Length)
		if calcError != nil {
			return nil, calcError
		}
		targetIndex, calcError := CalcIndex(indexObject.(*Tuple).Content[1], self.(*Array).Length)
		if calcError != nil {
			return nil, calcError
		}
		return NewTuple(vm.PeekSymbolTable(), self.(*Tuple).Content[startIndex:targetIndex]), nil
	} else {
		return nil, errors.NewTypeError(indexObject.TypeName(), IntegerName, TupleName)
	}
}

func TupleToString(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	result := "("
	var objectToString IObject
	var objectString IObject
	var callError *errors.Error
	for index, object := range self.(*Tuple).Content {
		if index != 0 {
			result += ", "
		}
		objectToString, getError = object.Get(ToString)
		if getError != nil {
			return nil, getError
		}
		if _, ok := objectToString.(*Function); !ok {
			return nil, errors.NewTypeError(objectToString.TypeName(), FunctionName)
		}
		objectString, callError = CallFunction(objectToString.(*Function), vm, vm.PeekSymbolTable())
		if callError != nil {
			return nil, callError
		}
		if _, ok := objectString.(*String); !ok {
			return nil, errors.NewTypeError(objectString.TypeName(), StringName)
		}
		result += objectString.(*String).Value
	}
	return NewString(vm.PeekSymbolTable(), result+")"), nil
}

func TupleToBool(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewBool(vm.PeekSymbolTable(), self.(*Tuple).Length != 0), nil
}

func TupleToArray(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewArray(vm.PeekSymbolTable(), append([]IObject{}, self.(*Tuple).Content...)), nil
}

func TupleToTuple(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewTuple(vm.PeekSymbolTable(), append([]IObject{}, self.(*Tuple).Content...)), nil
}
func TupleHash(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	var tupleHash int64 = 0
	var objectHashFunc IObject
	for _, object := range self.(*Tuple).Content {
		objectHashFunc, getError = object.Get(Hash)
		if getError != nil {
			return nil, getError
		}
		if _, ok := objectHashFunc.(*Function); !ok {
			return nil, errors.NewTypeError(objectHashFunc.TypeName(), FunctionName)
		}
		objectHash, callError := CallFunction(objectHashFunc.(*Function), vm, self.SymbolTable())
		if callError != nil {
			return nil, callError
		}
		if _, ok := objectHash.(*Integer); !ok {
			return nil, errors.NewTypeError(objectHash.TypeName(), IntegerName)
		}
		if tupleHash == 0 {
			tupleHash = objectHash.(*Integer).Value
		} else {
			tupleHash <<= objectHash.(*Integer).Value
		}
	}
	return NewInteger(vm.PeekSymbolTable(), tupleHash), nil
}

/*
	// Binary Operations
	//// Comparison Binary
	Equals            - (Done)
	RightEquals       - (Done)
	NotEquals         - (Done)
	RightNotEquals    - (Done)
	// Behavior
	Hash						- (Done)
	Copy						- (Done)
	Index   Integer or Tuple	- (Done)
	Iter
	// Transformation
	ToString      - (Done)
	ToBool        - (Done)
	ToArray       - (Done)
	ToTuple       - (Done)
*/
func TupleInitialize(_ VirtualMachine, object IObject) *errors.Error {
	object.Set(Equals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, TupleEquals)))
	object.Set(RightEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, TupleRightEquals)))
	object.Set(NotEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, TupleNotEquals)))
	object.Set(RightNotEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, TupleRightNotEquals)))
	object.Set(Hash, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, TupleHash)))
	object.Set(Copy, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, TupleCopy)))
	object.Set(Index, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, TupleIndex)))
	// object.Set(Iter, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, TupleIter)))
	object.Set(ToString, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, TupleToString)))
	object.Set(ToBool, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, TupleToBool)))
	object.Set(ToArray, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, TupleToArray)))
	object.Set(ToTuple, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, TupleToTuple)))
	return nil
}

func NewTuple(parentSymbols *SymbolTable, content []IObject) *Tuple {
	tuple := &Tuple{
		Object:  NewObject(TupleName, nil, parentSymbols),
		Content: content,
		Length:  len(content),
	}
	TupleInitialize(nil, tuple)
	return tuple
}

type Bool struct {
	*Object
	Value bool
}

func BoolEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	if _, ok := right.(*Bool); !ok {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	return NewBool(vm.PeekSymbolTable(), self.(*Bool).Value == right.(*Bool).Value), nil
}

func BoolRightEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	if _, ok := left.(*Bool); !ok {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	return NewBool(vm.PeekSymbolTable(), left.(*Bool).Value == self.(*Bool).Value), nil
}

func BoolNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	if _, ok := right.(*Bool); !ok {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	return NewBool(vm.PeekSymbolTable(), self.(*Bool).Value != right.(*Bool).Value), nil
}

func BoolRightNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	if _, ok := left.(*Bool); !ok {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	return NewBool(vm.PeekSymbolTable(), left.(*Bool).Value != self.(*Bool).Value), nil
}

func BoolCopy(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewBool(vm.PeekSymbolTable(), self.(*Bool).Value), nil
}

func BoolToInteger(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	if self.(*Bool).Value {
		return NewInteger(vm.PeekSymbolTable(), 1), nil
	}
	return NewInteger(vm.PeekSymbolTable(), 0), nil
}

func BoolToFloat(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	if self.(*Bool).Value {
		return NewFloat(vm.PeekSymbolTable(), 1), nil
	}
	return NewFloat(vm.PeekSymbolTable(), 0), nil
}

func BoolToString(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	if self.(*Bool).Value {
		return NewString(vm.PeekSymbolTable(), TrueName), nil
	}
	return NewString(vm.PeekSymbolTable(), FalseName), nil
}

func BoolHash(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	if self.GetHash() == 0 {
		boolHash, hashingError := vm.HashString(fmt.Sprintf("%t-%s", self.(*Bool).Value, BoolName))
		if hashingError != nil {
			return nil, hashingError
		}
		self.SetHash(boolHash)
	}
	return NewInteger(vm.PeekSymbolTable(), self.GetHash()), nil
}

/*
	// Binary Operations
	//// Comparison Binary
	Equals - (Done)
	RightEquals - (Done)
	NotEquals           - (Done)
	RightNotEquals      - (Done)
	// Behavior
	Hash			   - (Done)
	Copy               - (Done)
	// Transformation
	ToInteger    - (Done)
	ToFloat      - (Done)
	ToString   - (Done)
	ToBool      - (Done)
*/
func BoolInitialize(_ VirtualMachine, object IObject) *errors.Error {
	object.Set(Equals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, BoolEquals)))
	object.Set(RightEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, BoolRightEquals)))
	object.Set(NotEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, BoolNotEquals)))
	object.Set(RightNotEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, BoolRightNotEquals)))
	object.Set(Copy, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, BoolHash)))
	object.Set(Copy, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, BoolCopy)))
	object.Set(ToInteger, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, BoolToInteger)))
	object.Set(ToFloat, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, BoolToFloat)))
	object.Set(ToString, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, BoolToString)))
	object.Set(ToBool, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, BoolCopy)))
	return nil
}

func NewBool(parentSymbols *SymbolTable, value bool) *Bool {
	bool_ := &Bool{
		Object: NewObject(BoolName, nil, parentSymbols),
		Value:  value,
	}
	BoolInitialize(nil, bool_)
	return bool_
}

type KeyValue struct {
	Key   IObject
	Value IObject
}
type HashTable struct {
	*Object
	Content map[int64][]*KeyValue
	Length  int
}

func HashTableEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	rawRight := arguments[0]
	if _, ok := rawRight.(*HashTable); !ok {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	right := rawRight.(*HashTable)
	if self.(*HashTable).Length != right.Length {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var rightIndex IObject
	rightIndex, getError = right.Get(Index)
	if getError != nil {
		return nil, getError
	}
	if _, ok := rightIndex.(*Function); !ok {
		return nil, errors.NewTypeError(rightIndex.TypeName(), FunctionName)
	}
	for key, leftValue := range self.(*HashTable).Content {
		// Check if other has the key
		rightValue, ok := right.Content[key]
		if !ok {
			return NewBool(vm.PeekSymbolTable(), false), nil
		}
		// Check if the each entry one has the same length
		if len(leftValue) != len(rightValue) {
			return NewBool(vm.PeekSymbolTable(), false), nil
		}
		// Start comparing the entries
		for _, entry := range leftValue {
			_, indexingError := CallFunction(rightIndex.(*Function), vm, vm.PeekSymbolTable(), entry.Key)
			if indexingError != nil {
				return NewBool(vm.PeekSymbolTable(), false), nil
			}
		}
	}
	return NewBool(vm.PeekSymbolTable(), true), nil
}

func HashTableRightEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	rawLeft := arguments[0]
	if _, ok := rawLeft.(*HashTable); !ok {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	left := rawLeft.(*HashTable)
	if self.(*HashTable).Length != left.Length {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftIndex IObject
	leftIndex, getError = left.Get(Index)
	if getError != nil {
		return nil, getError
	}
	if _, ok := leftIndex.(*Function); !ok {
		return nil, errors.NewTypeError(leftIndex.TypeName(), FunctionName)
	}
	for key, leftValue := range left.Content {
		// Check if other has the key
		rightValue, ok := self.(*HashTable).Content[key]
		if !ok {
			return NewBool(vm.PeekSymbolTable(), false), nil
		}
		// Check if the each entry one has the same length
		if len(leftValue) != len(rightValue) {
			return NewBool(vm.PeekSymbolTable(), false), nil
		}
		// Start comparing the entries
		for _, entry := range leftValue {
			_, indexingError := CallFunction(leftIndex.(*Function), vm, vm.PeekSymbolTable(), entry.Key)
			if indexingError != nil {
				return NewBool(vm.PeekSymbolTable(), false), nil
			}
		}
	}
	return NewBool(vm.PeekSymbolTable(), true), nil
}

func HashTableNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	rawRight := arguments[0]
	if _, ok := rawRight.(*HashTable); !ok {
		return NewBool(vm.PeekSymbolTable(), true), nil
	}
	right := rawRight.(*HashTable)
	if self.(*HashTable).Length != right.Length {
		return NewBool(vm.PeekSymbolTable(), true), nil
	}
	var rightIndex IObject
	rightIndex, getError = right.Get(Index)
	if getError != nil {
		return nil, getError
	}
	if _, ok := rightIndex.(*Function); !ok {
		return nil, errors.NewTypeError(rightIndex.TypeName(), FunctionName)
	}
	for key, leftValue := range self.(*HashTable).Content {
		// Check if other has the key
		rightValue, ok := right.Content[key]
		if !ok {
			return NewBool(vm.PeekSymbolTable(), true), nil
		}
		// Check if the each entry one has the same length
		if len(leftValue) != len(rightValue) {
			return NewBool(vm.PeekSymbolTable(), true), nil
		}
		// Start comparing the entries
		for _, entry := range leftValue {
			_, indexingError := CallFunction(rightIndex.(*Function), vm, vm.PeekSymbolTable(), entry.Key)
			if indexingError != nil {
				return NewBool(vm.PeekSymbolTable(), true), nil
			}
		}
	}
	return NewBool(vm.PeekSymbolTable(), false), nil
}

func HashTableRightNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	rawLeft := arguments[0]
	if _, ok := rawLeft.(*HashTable); !ok {
		return NewBool(vm.PeekSymbolTable(), true), nil
	}
	left := rawLeft.(*HashTable)
	if self.(*HashTable).Length != left.Length {
		return NewBool(vm.PeekSymbolTable(), true), nil
	}
	var leftIndex IObject
	leftIndex, getError = left.Get(Index)
	if getError != nil {
		return nil, getError
	}
	if _, ok := leftIndex.(*Function); !ok {
		return nil, errors.NewTypeError(leftIndex.TypeName(), FunctionName)
	}
	for key, leftValue := range left.Content {
		// Check if other has the key
		rightValue, ok := self.(*HashTable).Content[key]
		if !ok {
			return NewBool(vm.PeekSymbolTable(), true), nil
		}
		// Check if the each entry one has the same length
		if len(leftValue) != len(rightValue) {
			return NewBool(vm.PeekSymbolTable(), true), nil
		}
		// Start comparing the entries
		for _, entry := range leftValue {
			_, indexingError := CallFunction(leftIndex.(*Function), vm, vm.PeekSymbolTable(), entry.Key)
			if indexingError != nil {
				return NewBool(vm.PeekSymbolTable(), true), nil
			}
		}
	}
	return NewBool(vm.PeekSymbolTable(), false), nil
}

func HashTableHash(_ VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	return nil, errors.NewUnhashableTypeError(errors.UnknownLine)
}

func HashTableIndex(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	indexObject := arguments[0]
	var indexObjectHash IObject
	indexObjectHash, getError = indexObject.Get(Hash)
	if _, ok := indexObjectHash.(*Function); !ok {
		return nil, errors.NewTypeError(indexObjectHash.TypeName(), FunctionName)
	}
	indexHash, callError := CallFunction(indexObjectHash.(*Function), vm, indexObject.SymbolTable())
	if callError != nil {
		return nil, callError
	}
	if _, ok := indexHash.(*Integer); !ok {
		return nil, errors.NewTypeError(indexHash.TypeName(), IntegerName)
	}
	keyValues, found := self.(*HashTable).Content[indexHash.(*Integer).Value]
	if !found {
		return nil, errors.NewKeyNotFoundError()
	}
	var indexObjectEquals IObject
	indexObjectEquals, getError = indexObject.Get(Equals)
	if _, ok := indexObjectEquals.(*Function); !ok {
		return nil, errors.NewTypeError(indexObjectEquals.TypeName(), FunctionName)
	}
	var equals IObject
	for _, keyValue := range keyValues {
		equals, callError = CallFunction(indexObjectEquals.(*Function), vm, indexObject.SymbolTable(), keyValue.Key)
		if callError != nil {
			return nil, callError
		}
		if _, ok := equals.(*Bool); !ok {
			return nil, errors.NewTypeError(equals.TypeName(), BoolName)
		}
		if equals.(*Bool).Value {
			return keyValue.Value, nil
		}
	}
	return nil, errors.NewKeyNotFoundError()
}

func HashTableAssign(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	indexObject := arguments[0]
	newValue := arguments[1]
	var indexObjectHash IObject
	indexObjectHash, getError = indexObject.Get(Hash)
	if _, ok := indexObjectHash.(*Function); !ok {
		return nil, errors.NewTypeError(indexObjectHash.TypeName(), FunctionName)
	}
	indexHash, callError := CallFunction(indexObjectHash.(*Function), vm, indexObject.SymbolTable())
	if callError != nil {
		return nil, callError
	}
	if _, ok := indexHash.(*Integer); !ok {
		return nil, errors.NewTypeError(indexHash.TypeName(), IntegerName)
	}
	keyValues, found := self.(*HashTable).Content[indexHash.(*Integer).Value]
	if !found {
		return nil, errors.NewKeyNotFoundError()
	}
	var indexObjectEquals IObject
	indexObjectEquals, getError = indexObject.Get(Equals)
	if _, ok := indexObjectEquals.(*Function); !ok {
		return nil, errors.NewTypeError(indexObjectEquals.TypeName(), FunctionName)
	}
	var equals IObject
	for index, keyValue := range keyValues {
		equals, callError = CallFunction(indexObjectEquals.(*Function), vm, indexObject.SymbolTable(), keyValue.Key)
		if callError != nil {
			return nil, callError
		}
		if _, ok := equals.(*Bool); !ok {
			return nil, errors.NewTypeError(equals.TypeName(), BoolName)
		}
		if equals.(*Bool).Value {
			self.(*HashTable).Content[indexHash.(*Integer).Value][index].Value = newValue
			self.(*HashTable).Length++
			return vm.PeekSymbolTable().GetAny(None)
		}
	}
	self.(*HashTable).Length++
	self.(*HashTable).Content[indexHash.(*Integer).Value] = append(self.(*HashTable).Content[indexHash.(*Integer).Value], &KeyValue{
		Key:   indexObject,
		Value: newValue,
	})
	return vm.PeekSymbolTable().GetAny(None)
}

func HashTableToString(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	result := "{"
	var (
		keyToString   IObject
		keyString     IObject
		valueToString IObject
		valueString   IObject
		callError     *errors.Error
	)
	for _, keyValues := range self.(*HashTable).Content {
		for _, keyValue := range keyValues {
			keyToString, getError = keyValue.Key.Get(ToString)
			if getError != nil {
				return nil, getError
			}
			keyString, callError = CallFunction(keyToString.(*Function), vm, keyValue.Key.SymbolTable())
			if callError != nil {
				return nil, callError
			}
			if _, ok := keyString.(*String); !ok {
				return nil, errors.NewTypeError(keyString.TypeName(), StringName)
			}
			result += keyString.(*String).Value
			valueToString, getError = keyValue.Value.Get(ToString)
			if getError != nil {
				return nil, getError
			}
			valueString, callError = CallFunction(valueToString.(*Function), vm, keyValue.Value.SymbolTable())
			if callError != nil {
				return nil, callError
			}
			if _, ok := valueString.(*String); !ok {
				return nil, errors.NewTypeError(keyString.TypeName(), StringName)
			}
			result += ":" + valueString.(*String).Value + ","
		}
	}
	if len(result) > 1 {
		result = result[:len(result)-1]
	}
	return NewString(vm.PeekSymbolTable(), result+"}"), nil
}

func HashTableToBool(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	if self.(*HashTable).Length > 0 {
		return NewBool(vm.PeekSymbolTable(), true), nil
	}
	return NewBool(vm.PeekSymbolTable(), false), nil
}

func HashTableToArray(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	var keys []IObject
	for _, keyValues := range self.(*HashTable).Content {
		for _, keyValue := range keyValues {
			keys = append(keys, keyValue.Key)
		}
	}
	return NewArray(vm.PeekSymbolTable(), keys), nil
}

func HashTableToTuple(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	var keys []IObject
	for _, keyValues := range self.(*HashTable).Content {
		for _, keyValue := range keyValues {
			keys = append(keys, keyValue.Key)
		}
	}
	return NewTuple(vm.PeekSymbolTable(), keys), nil
}

/*
	// Binary Operations
	//// Comparison Binary
	Equals                 - (Done)
	RightEquals            - (Done)
	NotEquals              - (Done)
	RightNotEquals         - (Done)
	// Behavior
	Hash  Unhashable type  - (Done)
	Copy                   - ()
	Index                  - (Done)
	Assign                 - (Done)
	// Transformation
	ToString  			- (Done)
	ToBool    			- (Done)
	ToArray   Keys only - (Done)
	ToTuple   Keys only - (Done)
*/
func HashTableInitialize(_ VirtualMachine, object IObject) *errors.Error {
	object.Set(Equals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, HashTableEquals)))
	object.Set(RightEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, HashTableRightEquals)))
	object.Set(NotEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, HashTableNotEquals)))
	object.Set(RightNotEquals, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, HashTableRightNotEquals)))

	object.Set(Hash, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, HashTableHash)))
	// object.Set(Copy, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, HashTableHash)))
	object.Set(Index, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, HashTableIndex)))
	object.Set(Assign, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 2, HashTableAssign)))

	object.Set(ToString, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, HashTableToString)))
	object.Set(ToBool, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, HashTableToBool)))
	object.Set(ToArray, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, HashTableToArray)))
	object.Set(ToTuple, NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, HashTableToTuple)))
	return nil
}

func NewHashTable(parent *SymbolTable, entries map[int64][]*KeyValue, entriesLength int) *HashTable {
	hashTable := &HashTable{
		Object:  NewObject(HashName, nil, parent),
		Content: entries,
		Length:  entriesLength,
	}
	HashTableInitialize(nil, hashTable)
	return hashTable
}

type Type struct {
	*Object
	Constructor Constructor
}

func NewType(parent *SymbolTable, subclasses []*Type, constructor Constructor) *Type {
	return &Type{
		Object:      NewObject(TypeName, subclasses, parent),
		Constructor: constructor,
	}
}

/*
	SetDefaultSymbolTable
	// Types
	Type       - (Done)
	Function   - (Done)
	Object     - (Done)
	Bool       - (Done)
	Bytes      - (Done)
	String     - (Done)
	HashTable  - (Done)
	Integer    - (Done)
	Array      - (Done)
	Tuple      - (Done)
	// Names
	None 	   - ()
	// Functions
	Hash       - ()
	Id         - ()
	Range      - ()
	Len        - ()
	// Direct To... (Transformations)
	ToString     - (Done)
	ToTuple      - (Done)
	ToArray      - (Done)
	ToInteger    - (Done)
	ToFloat      - (Done)
	ToBool       - (Done)
	ToHashTable  - ()
	ToIter       - ()
	ToBytes      - ()
	ToObject     - ()
*/
func SetDefaultSymbolTable() *SymbolTable {
	symbolTable := NewSymbolTable(nil)
	// Types
	type_ := NewType(nil, nil,
		NewBuiltInConstructor(ObjectInitialize),
	)
	symbolTable.Set(TypeName, type_)
	symbolTable.Set(BoolName,
		NewType(nil, []*Type{type_},
			NewBuiltInConstructor(BoolInitialize),
		),
	)
	symbolTable.Set(ObjectName,
		NewType(nil, []*Type{type_},
			NewBuiltInConstructor(ObjectInitialize),
		),
	)
	symbolTable.Set(FunctionName,
		NewType(nil, []*Type{type_},
			NewBuiltInConstructor(func(machine VirtualMachine, object IObject) *errors.Error {
				return nil
			}),
		),
	)
	symbolTable.Set(IntegerName,
		NewType(nil, []*Type{type_},
			NewBuiltInConstructor(IntegerInitialize),
		),
	)
	symbolTable.Set(StringName,
		NewType(nil, []*Type{type_},
			NewBuiltInConstructor(StringInitialize),
		),
	)
	symbolTable.Set(BytesName,
		NewType(nil, []*Type{type_},
			NewBuiltInConstructor(BytesInitialize),
		),
	)
	symbolTable.Set(TupleName,
		NewType(nil, []*Type{type_},
			NewBuiltInConstructor(TupleInitialize),
		),
	)
	symbolTable.Set(ArrayName,
		NewType(nil, []*Type{type_},
			NewBuiltInConstructor(ArrayInitialize),
		),
	)
	symbolTable.Set(HashName,
		NewType(nil, []*Type{type_},
			NewBuiltInConstructor(HashTableInitialize),
		),
	)
	// Names
	symbolTable.Set(None,
		NewObject(NoneName, nil, symbolTable),
	)
	// To... (Transformations)
	symbolTable.Set(ToFloat,
		NewFunction(symbolTable,
			NewBuiltInFunction(1,
				func(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
					toFloat, getError := arguments[0].Get(ToFloat)
					if getError != nil {
						return nil, getError
					}
					if _, ok := toFloat.(*Function); !ok {
						return nil, errors.NewTypeError(toFloat.(IObject).TypeName(), FunctionName)
					}
					return CallFunction(toFloat.(*Function), vm, arguments[0].SymbolTable().Parent)
				},
			),
		),
	)
	symbolTable.Set(ToString,
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
					return CallFunction(toString.(*Function), vm, arguments[0].SymbolTable().Parent)
				},
			),
		),
	)
	symbolTable.Set(ToInteger,
		NewFunction(symbolTable,
			NewBuiltInFunction(1,
				func(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
					toInteger, getError := arguments[0].Get(ToInteger)
					if getError != nil {
						return nil, getError
					}
					if _, ok := toInteger.(*Function); !ok {
						return nil, errors.NewTypeError(toInteger.(IObject).TypeName(), FunctionName)
					}
					return CallFunction(toInteger.(*Function), vm, arguments[0].SymbolTable().Parent)
				},
			),
		),
	)
	symbolTable.Set(ToArray,
		NewFunction(symbolTable,
			NewBuiltInFunction(1,
				func(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
					toArray, getError := arguments[0].Get(ToArray)
					if getError != nil {
						return nil, getError
					}
					if _, ok := toArray.(*Function); !ok {
						return nil, errors.NewTypeError(toArray.(IObject).TypeName(), FunctionName)
					}
					return CallFunction(toArray.(*Function), vm, arguments[0].SymbolTable().Parent)
				},
			),
		),
	)
	symbolTable.Set(ToTuple,
		NewFunction(symbolTable,
			NewBuiltInFunction(1,
				func(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
					toTuple, getError := arguments[0].Get(ToTuple)
					if getError != nil {
						return nil, getError
					}
					if _, ok := toTuple.(*Function); !ok {
						return nil, errors.NewTypeError(toTuple.(IObject).TypeName(), FunctionName)
					}
					return CallFunction(toTuple.(*Function), vm, arguments[0].SymbolTable().Parent)
				},
			),
		),
	)
	symbolTable.Set(ToBool,
		NewFunction(symbolTable,
			NewBuiltInFunction(1,
				func(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
					toBool, getError := arguments[0].Get(ToBool)
					if getError != nil {
						return nil, getError
					}
					if _, ok := toBool.(*Function); !ok {
						return nil, errors.NewTypeError(toBool.(IObject).TypeName(), FunctionName)
					}
					return CallFunction(toBool.(*Function), vm, arguments[0].SymbolTable().Parent)
				},
			),
		),
	)
	return symbolTable
}
