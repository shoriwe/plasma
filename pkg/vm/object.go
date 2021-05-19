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
	None         = "None"
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
	Dir        = "Dir"
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
	GetHash() int64
	SetHash(int64)
}

// MetaClass for IObject
type Object struct {
	id             uint
	typeName       string
	subClasses     []*Function
	symbols        *SymbolTable
	virtualMachine VirtualMachine
	hash           int64
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

// NewObject Creates an Empty Object
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
func NewObject(
	typeName string,
	subClasses []*Function,
	parentSymbols *SymbolTable,
) *Object {
	result := &Object{
		id:         counter.NextId(),
		typeName:   typeName,
		subClasses: subClasses,
		hash:       0,
	}
	result.symbols = NewSymbolTable(parentSymbols)
	result.symbols.Symbols = map[string]IObject{
		// IObject Creation
		Initialize: NewFunction(result.symbols, NewBuiltInClassFunction(result, 0, ObjInitialize)),
		// Unary Operations
		NegBits: NewFunction(result.symbols, NewNotImplementedCallable(0)),
		Negate:  NewFunction(result.symbols, NewBuiltInClassFunction(result, 0, ObjNegate)),
		// Binary Operations
		//// Math binary
		Add:           NewFunction(result.symbols, NewNotImplementedCallable(1)),
		RightAdd:      NewFunction(result.symbols, NewNotImplementedCallable(1)),
		Sub:           NewFunction(result.symbols, NewNotImplementedCallable(1)),
		RightSub:      NewFunction(result.symbols, NewNotImplementedCallable(1)),
		Mul:           NewFunction(result.symbols, NewNotImplementedCallable(1)),
		RightMul:      NewFunction(result.symbols, NewNotImplementedCallable(1)),
		Div:           NewFunction(result.symbols, NewNotImplementedCallable(1)),
		RightDiv:      NewFunction(result.symbols, NewNotImplementedCallable(1)),
		Mod:           NewFunction(result.symbols, NewNotImplementedCallable(1)),
		RightMod:      NewFunction(result.symbols, NewNotImplementedCallable(1)),
		Pow:           NewFunction(result.symbols, NewNotImplementedCallable(1)),
		RightPow:      NewFunction(result.symbols, NewNotImplementedCallable(1)),
		BitXor:        NewFunction(result.symbols, NewNotImplementedCallable(1)),
		RightBitXor:   NewFunction(result.symbols, NewNotImplementedCallable(1)),
		BitAnd:        NewFunction(result.symbols, NewNotImplementedCallable(1)),
		RightBitAnd:   NewFunction(result.symbols, NewNotImplementedCallable(1)),
		BitOr:         NewFunction(result.symbols, NewNotImplementedCallable(1)),
		RightBitOr:    NewFunction(result.symbols, NewNotImplementedCallable(1)),
		BitLeft:       NewFunction(result.symbols, NewNotImplementedCallable(1)),
		RightBitLeft:  NewFunction(result.symbols, NewNotImplementedCallable(1)),
		BitRight:      NewFunction(result.symbols, NewNotImplementedCallable(1)),
		RightBitRight: NewFunction(result.symbols, NewNotImplementedCallable(1)),
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
		GreaterThan:             NewFunction(result.symbols, NewNotImplementedCallable(1)),
		RightGreaterThan:        NewFunction(result.symbols, NewNotImplementedCallable(1)),
		LessThan:                NewFunction(result.symbols, NewNotImplementedCallable(1)),
		RightLessThan:           NewFunction(result.symbols, NewNotImplementedCallable(1)),
		GreaterThanOrEqual:      NewFunction(result.symbols, NewNotImplementedCallable(1)),
		RightGreaterThanOrEqual: NewFunction(result.symbols, NewNotImplementedCallable(1)),
		LessThanOrEqual:         NewFunction(result.symbols, NewNotImplementedCallable(1)),
		RightLessThanOrEqual:    NewFunction(result.symbols, NewNotImplementedCallable(1)),
		// Behavior
		Hash:       NewFunction(result.symbols, NewBuiltInClassFunction(result, 0, ObjHash)),
		Copy:       NewFunction(result.symbols, NewNotImplementedCallable(0)),
		Dir:        NewFunction(result.symbols, NewNotImplementedCallable(0)),
		Index:      NewFunction(result.symbols, NewNotImplementedCallable(1)),
		Assign:     NewFunction(result.symbols, NewNotImplementedCallable(2)),
		Call:       NewFunction(result.symbols, NewNotImplementedCallable(0)),
		Iter:       NewFunction(result.symbols, NewNotImplementedCallable(0)),
		Class:      NewFunction(result.symbols, NewNotImplementedCallable(0)),
		SubClasses: NewFunction(result.symbols, NewNotImplementedCallable(0)),
		// Transformation
		ToInteger: NewFunction(result.symbols, NewNotImplementedCallable(0)),
		ToFloat:   NewFunction(result.symbols, NewNotImplementedCallable(0)),
		ToString:  NewFunction(result.symbols, NewBuiltInClassFunction(result, 0, ObjToString)),
		ToBool:    NewFunction(result.symbols, NewBuiltInClassFunction(result, 0, ObjToBool)),
		ToArray:   NewFunction(result.symbols, NewNotImplementedCallable(0)),
		ToTuple:   NewFunction(result.symbols, NewNotImplementedCallable(0)),
	}
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
func NewString(
	parentSymbols *SymbolTable,
	value string,
) *String {
	string_ := &String{
		Object: NewObject(StringName, nil, parentSymbols),
		Value:  value,
		Length: len(value),
	}
	string_.Set(Add, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 1, StringAdd)))
	string_.Set(RightAdd, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 1, StringRightAdd)))
	string_.Set(Mul, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 1, StringMul)))
	string_.Set(RightMul, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 1, StringRightMul)))

	string_.Set(Equals, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringEquals)))
	string_.Set(RightEquals, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringRightEquals)))
	string_.Set(NotEquals, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringNotEquals)))
	string_.Set(RightNotEquals, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringRightNotEquals)))

	string_.Set(Hash, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringHash)))
	string_.Set(Copy, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringCopy)))
	string_.Set(Index, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringIndex)))
	// string_.Set(Iter, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringIter)))

	string_.Set(ToInteger, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringToInteger)))
	string_.Set(ToFloat, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringToFloat)))
	string_.Set(ToString, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringCopy)))
	string_.Set(ToBool, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringToBool)))
	string_.Set(ToArray, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringToArray)))
	string_.Set(ToTuple, NewFunction(string_.symbols, NewBuiltInClassFunction(string_, 0, StringToTuple)))
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
func NewBytes(parent *SymbolTable, content []*Integer) IObject {
	bytes_ := &Bytes{
		Object:  NewObject(BytesName, nil, parent),
		Content: content,
		Length:  len(content),
	}
	bytes_.Set(Add, NewFunction(bytes_.symbols, NewBuiltInClassFunction(bytes_, 1, BytesAdd)))
	bytes_.Set(RightAdd, NewFunction(bytes_.symbols, NewBuiltInClassFunction(bytes_, 1, BytesRightAdd)))
	bytes_.Set(Mul, NewFunction(bytes_.symbols, NewBuiltInClassFunction(bytes_, 1, BytesMul)))
	bytes_.Set(RightMul, NewFunction(bytes_.symbols, NewBuiltInClassFunction(bytes_, 1, BytesRightMul)))

	bytes_.Set(Equals, NewFunction(bytes_.symbols, NewBuiltInClassFunction(bytes_, 1, BytesEquals)))
	bytes_.Set(RightEquals, NewFunction(bytes_.symbols, NewBuiltInClassFunction(bytes_, 1, BytesRightEquals)))
	bytes_.Set(NotEquals, NewFunction(bytes_.symbols, NewBuiltInClassFunction(bytes_, 1, BytesNotEquals)))
	bytes_.Set(RightNotEquals, NewFunction(bytes_.symbols, NewBuiltInClassFunction(bytes_, 1, BytesRightNotEquals)))

	bytes_.Set(Hash, NewFunction(bytes_.symbols, NewBuiltInClassFunction(bytes_, 0, BytesHash)))
	bytes_.Set(Copy, NewFunction(bytes_.symbols, NewBuiltInClassFunction(bytes_, 0, BytesCopy)))
	bytes_.Set(Index, NewFunction(bytes_.symbols, NewBuiltInClassFunction(bytes_, 1, BytesIndex)))

	bytes_.Set(ToInteger, NewFunction(bytes_.symbols, NewBuiltInClassFunction(bytes_, 0, BytesToInteger)))
	bytes_.Set(ToString, NewFunction(bytes_.symbols, NewBuiltInClassFunction(bytes_, 0, BytesToString)))
	bytes_.Set(ToBool, NewFunction(bytes_.symbols, NewBuiltInClassFunction(bytes_, 0, BytesToBool)))
	bytes_.Set(ToArray, NewFunction(bytes_.symbols, NewBuiltInClassFunction(bytes_, 0, BytesToArray)))
	bytes_.Set(ToTuple, NewFunction(bytes_.symbols, NewBuiltInClassFunction(bytes_, 0, BytesToTuple)))
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

	integer.Set(Hash, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerHash)))
	integer.Set(Copy, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerCopy)))

	integer.Set(ToInteger, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerCopy)))
	integer.Set(ToFloat, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerToFloat)))
	integer.Set(ToString, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerToString)))
	integer.Set(ToBool, NewFunction(integer.symbols, NewBuiltInClassFunction(integer, 0, IntegerToBool)))
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
	Hash          - Done
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

	float_.Set(Hash, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatHash)))
	float_.Set(Copy, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatCopy)))

	float_.Set(ToInteger, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatToInteger)))
	float_.Set(ToFloat, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatCopy)))
	float_.Set(ToString, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatToString)))
	float_.Set(ToBool, NewFunction(float_.symbols, NewBuiltInClassFunction(float_, 0, FloatToBool)))
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

func NewArray(parentSymbols *SymbolTable, content []IObject) *Array {
	array := &Array{
		Object:  NewObject(ArrayName, nil, parentSymbols),
		Content: content,
		Length:  len(content),
	}
	array.Set(Equals, NewFunction(array.symbols, NewBuiltInClassFunction(array, 1, ArrayEquals)))
	array.Set(RightEquals, NewFunction(array.symbols, NewBuiltInClassFunction(array, 1, ArrayRightEquals)))
	array.Set(NotEquals, NewFunction(array.symbols, NewBuiltInClassFunction(array, 1, ArrayNotEquals)))
	array.Set(RightNotEquals, NewFunction(array.symbols, NewBuiltInClassFunction(array, 1, ArrayRightNotEquals)))
	array.Set(Hash, NewFunction(array.symbols, NewBuiltInClassFunction(array, 0, ArrayHash)))
	array.Set(Copy, NewFunction(array.symbols, NewBuiltInClassFunction(array, 0, ArrayCopy)))
	array.Set(Index, NewFunction(array.symbols, NewBuiltInClassFunction(array, 1, ArrayIndex)))
	array.Set(Assign, NewFunction(array.symbols, NewBuiltInClassFunction(array, 2, ArrayAssign)))
	// array.Set(Iter, NewFunction(array.symbols, NewBuiltInClassFunction(array, 0, ArrayIter)))
	array.Set(ToString, NewFunction(array.symbols, NewBuiltInClassFunction(array, 0, ArrayToString)))
	array.Set(ToBool, NewFunction(array.symbols, NewBuiltInClassFunction(array, 0, ArrayToBool)))
	array.Set(ToArray, NewFunction(array.symbols, NewBuiltInClassFunction(array, 0, ArrayToArray)))
	array.Set(ToTuple, NewFunction(array.symbols, NewBuiltInClassFunction(array, 0, ArrayToTuple)))
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
	Hash      - (Done)
	Copy      - (Done)
	Index   Integer or Tuple
	Iter
	// Transformation
	ToString      - (Done)
	ToBool        - (Done)
	ToArray       - (Done)
	ToTuple       - (Done)
*/
func NewTuple(parentSymbols *SymbolTable, content []IObject) *Tuple {
	tuple := &Tuple{
		Object:  NewObject(TupleName, nil, parentSymbols),
		Content: content,
		Length:  len(content),
	}
	tuple.Set(Equals, NewFunction(tuple.symbols, NewBuiltInClassFunction(tuple, 1, TupleEquals)))
	tuple.Set(RightEquals, NewFunction(tuple.symbols, NewBuiltInClassFunction(tuple, 1, TupleRightEquals)))
	tuple.Set(NotEquals, NewFunction(tuple.symbols, NewBuiltInClassFunction(tuple, 1, TupleNotEquals)))
	tuple.Set(RightNotEquals, NewFunction(tuple.symbols, NewBuiltInClassFunction(tuple, 1, TupleRightNotEquals)))
	tuple.Set(Hash, NewFunction(tuple.symbols, NewBuiltInClassFunction(tuple, 0, TupleHash)))
	tuple.Set(Copy, NewFunction(tuple.symbols, NewBuiltInClassFunction(tuple, 0, TupleCopy)))
	tuple.Set(Index, NewFunction(tuple.symbols, NewBuiltInClassFunction(tuple, 1, TupleIndex)))
	// tuple.Set(Iter, NewFunction(tuple.symbols, NewBuiltInClassFunction(tuple, 0, TupleIter)))
	tuple.Set(ToString, NewFunction(tuple.symbols, NewBuiltInClassFunction(tuple, 0, TupleToString)))
	tuple.Set(ToBool, NewFunction(tuple.symbols, NewBuiltInClassFunction(tuple, 0, TupleToBool)))
	tuple.Set(ToArray, NewFunction(tuple.symbols, NewBuiltInClassFunction(tuple, 0, TupleToArray)))
	tuple.Set(ToTuple, NewFunction(tuple.symbols, NewBuiltInClassFunction(tuple, 0, TupleToTuple)))
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

func NewBool(parentSymbols *SymbolTable, value bool) *Bool {
	bool_ := &Bool{
		Object: NewObject(BoolName, nil, parentSymbols),
		Value:  value,
	}
	bool_.Set(Equals, NewFunction(bool_.symbols, NewBuiltInClassFunction(bool_, 1, BoolEquals)))
	bool_.Set(RightEquals, NewFunction(bool_.symbols, NewBuiltInClassFunction(bool_, 1, BoolRightEquals)))
	bool_.Set(NotEquals, NewFunction(bool_.symbols, NewBuiltInClassFunction(bool_, 1, BoolNotEquals)))
	bool_.Set(RightNotEquals, NewFunction(bool_.symbols, NewBuiltInClassFunction(bool_, 1, BoolRightNotEquals)))
	bool_.Set(Copy, NewFunction(bool_.symbols, NewBuiltInClassFunction(bool_, 0, BoolHash)))
	bool_.Set(Copy, NewFunction(bool_.symbols, NewBuiltInClassFunction(bool_, 0, BoolCopy)))
	bool_.Set(ToInteger, NewFunction(bool_.symbols, NewBuiltInClassFunction(bool_, 0, BoolToInteger)))
	bool_.Set(ToFloat, NewFunction(bool_.symbols, NewBuiltInClassFunction(bool_, 0, BoolToFloat)))
	bool_.Set(ToString, NewFunction(bool_.symbols, NewBuiltInClassFunction(bool_, 0, BoolToString)))
	bool_.Set(ToBool, NewFunction(bool_.symbols, NewBuiltInClassFunction(bool_, 0, BoolCopy)))
	return bool_
}

/*
	SetDefaultSymbolTable
	// Names
	None 	   - (Done)
	// Functions
	Hash       - ()
	Id         - ()
	Range      - ()
	Len        - ()
	// Types
	String     - (Done)
	Tuple      - (Done)
	Array      - (Done)
	Integer    - (Done)
	Float      - (Done)
	Bool       - (Done)
	HashTable  - ()
	Iter       - ()
	Bytes      - ()
	Object     - (Done)
*/
func SetDefaultSymbolTable() *SymbolTable {
	symbolTable := NewSymbolTable(nil)
	symbolTable.Set(None,
		NewObject(NoneName, nil, symbolTable),
	)
	symbolTable.Set(FloatName,
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
	symbolTable.Set(ObjectName,
		NewFunction(symbolTable,
			NewBuiltInFunction(0,
				func(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
					return NewObject(ObjectName, nil, vm.PeekSymbolTable()), nil
				},
			),
		),
	)
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
					return CallFunction(toString.(*Function), vm, arguments[0].SymbolTable().Parent)
				},
			),
		),
	)
	symbolTable.Set(IntegerName,
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
	symbolTable.Set(ArrayName,
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
	symbolTable.Set(TupleName,
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
	symbolTable.Set(BoolName,
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
