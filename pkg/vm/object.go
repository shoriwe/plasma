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
	// Built-In Getters and Setters
	//// Getters
	GetInteger64 = "GetInteger64"
	GetBool      = "GetBool"
	GetBytes     = "GetBytes"
	GetString    = "GetString"
	GetFloat64   = "GetFloat64"
	GetContent   = "GetContent"
	GetKeyValues = "GetKeyValues"
	GetLength    = "GetLength"
	//// Setters
	SetBool      = "SetBool"
	SetBytes     = "SetBytes"
	SetString    = "SetString"
	SetInteger64 = "SetInteger64"
	SetFloat64   = "SetFloat64"
	SetContent   = "SetContent"
	SetKeyValues = "SetKeyValues"
	SetLength    = "SetLength"
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
	index := int(object.GetInteger64())
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

	GetBool() bool
	GetBytes() []*Integer
	GetString() string
	GetInteger64() int64
	GetFloat64() float64
	GetContent() []IObject
	GetKeyValues() map[int64][]*KeyValue
	GetLength() int

	SetBool(bool)
	SetBytes([]*Integer)
	SetString(string)
	SetInteger64(int64)
	SetFloat64(float64)
	SetContent([]IObject)
	SetKeyValues(map[int64][]*KeyValue)
	SetLength(int)
	IncreaseLength()
}

// MetaClass for IObject
type Object struct {
	id             uint
	typeName       string
	subClasses     []*Type
	symbols        *SymbolTable
	virtualMachine VirtualMachine
	hash           int64
	// Stuff Related to built in objects
	Bool      bool
	String    string
	Bytes     []*Integer
	Integer64 int64
	Float64   float64
	Content   []IObject
	KeyValues map[int64][]*KeyValue
	//
	Length int
}

func (o *Object) IncreaseLength() {
	o.Length++
}

func (o *Object) GetBool() bool {
	return o.Bool
}

func (o *Object) GetBytes() []*Integer {
	return o.Bytes
}

func (o *Object) GetString() string {
	return o.String
}

func (o *Object) GetInteger64() int64 {
	return o.Integer64
}

func (o *Object) GetFloat64() float64 {
	return o.Float64
}

func (o *Object) GetContent() []IObject {
	return o.Content
}

func (o *Object) GetKeyValues() map[int64][]*KeyValue {
	return o.KeyValues
}

func (o *Object) GetLength() int {
	return o.Length
}

func (o *Object) SetString(s string) {
	o.String = s
}

func (o *Object) SetInteger64(i int64) {
	o.Integer64 = i
}

func (o *Object) SetFloat64(f float64) {
	o.Float64 = f
}

func (o *Object) SetContent(objects []IObject) {
	o.Content = objects
}

func (o *Object) SetKeyValues(m map[int64][]*KeyValue) {
	o.KeyValues = m
}

func (o *Object) SetLength(i int) {
	o.Length = i
}

func (o *Object) SetBool(b bool) {
	o.Bool = b
}

func (o *Object) SetBytes(b []*Integer) {
	o.Bytes = b
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
	return NewBool(vm.PeekSymbolTable(), leftBool.GetBool() && rightBool.GetBool()), nil
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
	return NewBool(vm.PeekSymbolTable(), leftBool.GetBool() && rightBool.GetBool()), nil
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
	return NewBool(vm.PeekSymbolTable(), leftBool.GetBool() || rightBool.GetBool()), nil
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
	return NewBool(vm.PeekSymbolTable(), leftBool.GetBool() || rightBool.GetBool()), nil
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
	return NewBool(vm.PeekSymbolTable(), leftBool.GetBool() != rightBool.GetBool()), nil
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
	return NewBool(vm.PeekSymbolTable(), rightBool.GetBool() != leftBool.GetBool()), nil
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
	return NewBool(vm.PeekSymbolTable(), !selfBool.GetBool()), nil
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

func ObjGetInteger64(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewInteger(vm.PeekSymbolTable(), self.GetInteger64()), nil
}
func ObjGetBool(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewBool(vm.PeekSymbolTable(), self.GetBool()), nil
}
func ObjGetBytes(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewBytes(vm.PeekSymbolTable(), self.GetBytes()), nil
}
func ObjGetString(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewString(vm.PeekSymbolTable(), self.GetString()), nil
}
func ObjGetFloat64(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewFloat(vm.PeekSymbolTable(), self.GetFloat64()), nil
}
func ObjGetContent(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewArray(vm.PeekSymbolTable(), self.GetContent()), nil
}
func ObjGetKeyValues(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewHashTable(vm.PeekSymbolTable(), self.GetKeyValues(), self.GetLength()), nil
}
func ObjGetLength(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewInteger(vm.PeekSymbolTable(), int64(self.GetLength())), nil
}

func ObjSetBool(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	self.SetBool(arguments[0].GetBool())
	return vm.PeekSymbolTable().GetAny(None)
}
func ObjSetBytes(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	self.SetBytes(arguments[0].GetBytes())
	self.SetLength(arguments[0].GetLength())
	return vm.PeekSymbolTable().GetAny(None)
}
func ObjSetString(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	self.SetString(arguments[0].GetString())
	self.SetLength(arguments[0].GetLength())
	return vm.PeekSymbolTable().GetAny(None)
}
func ObjSetInteger64(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	self.SetInteger64(arguments[0].GetInteger64())
	return vm.PeekSymbolTable().GetAny(None)
}
func ObjSetFloat64(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	self.SetFloat64(arguments[0].GetFloat64())
	return vm.PeekSymbolTable().GetAny(None)
}
func ObjSetContent(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	self.SetContent(arguments[0].GetContent())
	self.SetLength(arguments[0].GetLength())
	return vm.PeekSymbolTable().GetAny(None)
}
func ObjSetKeyValues(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	self.SetKeyValues(arguments[0].GetKeyValues())
	self.SetLength(arguments[0].GetLength())
	return vm.PeekSymbolTable().GetAny(None)
}
func ObjSetLength(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	self.SetLength(arguments[0].GetLength())
	return vm.PeekSymbolTable().GetAny(None)
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
		// Getters and Setters for Built in properties
		//// Getters
		GetInteger64: NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, ObjGetInteger64)),
		GetBool:      NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, ObjGetBool)),
		GetBytes:     NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, ObjGetBytes)),
		GetString:    NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, ObjGetString)),
		GetFloat64:   NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, ObjGetFloat64)),
		GetContent:   NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, ObjGetContent)),
		GetKeyValues: NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, ObjGetKeyValues)),
		GetLength:    NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 0, ObjGetLength)),
		//// Setters
		SetBool:      NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ObjSetBool)),
		SetBytes:     NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ObjSetBytes)),
		SetString:    NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ObjSetString)),
		SetInteger64: NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ObjSetInteger64)),
		SetFloat64:   NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ObjSetFloat64)),
		SetContent:   NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ObjSetContent)),
		SetKeyValues: NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ObjSetKeyValues)),
		SetLength:    NewFunction(object.SymbolTable(), NewBuiltInClassFunction(object, 1, ObjSetLength)),
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
	result.Length = 0
	result.Bool = true
	result.String = ""
	result.Integer64 = 0
	result.Float64 = 0
	result.Content = []IObject{}
	result.KeyValues = map[int64][]*KeyValue{}
	result.Bytes = []*Integer{}
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
}

func StringAdd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[1]
	return NewString(
		vm.PeekSymbolTable(),
		self.GetString()+right.GetString(),
	), nil
}

func StringRightAdd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[1]
	return NewString(
		vm.PeekSymbolTable(),
		left.GetString()+self.GetString(),
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
		Repeat(self.GetString(), right.GetInteger64()),
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
		Repeat(self.GetString(), left.GetInteger64()),
	), nil
}

func StringCopy(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewString(
		vm.PeekSymbolTable(),
		self.GetString(),
	), nil
}

func StringToBool(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewBool(vm.PeekSymbolTable(), self.GetLength() != 0), nil
}

func StringEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	return NewBool(vm.PeekSymbolTable(), self.GetString() == right.GetString()), nil
}

func StringRightEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	return NewBool(vm.PeekSymbolTable(), left.GetString() == self.GetString()), nil
}

func StringNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	return NewBool(vm.PeekSymbolTable(), self.GetString() != right.GetString()), nil
}

func StringRightNotEquals(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	left := arguments[0]
	return NewBool(vm.PeekSymbolTable(), left.GetString() != self.GetString()), nil
}

func StringIndex(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	indexObject := arguments[0]
	if _, ok := indexObject.(*Integer); ok {
		index, getIndexError := CalcIndex(indexObject, self.GetLength())
		if getIndexError != nil {
			return nil, getIndexError
		}
		return NewString(vm.PeekSymbolTable(), string(self.GetString()[index])), nil
	} else if _, ok = indexObject.(*Tuple); ok {
		if len(indexObject.GetContent()) != 2 {
			return nil, errors.NewInvalidNumberOfArguments(len(indexObject.GetContent()), 2)
		}
		startIndex, calcError := CalcIndex(indexObject.GetContent()[0], self.GetLength())
		if calcError != nil {
			return nil, calcError
		}
		targetIndex, calcError := CalcIndex(indexObject.GetContent()[1], self.GetLength())
		if calcError != nil {
			return nil, calcError
		}
		return NewString(vm.PeekSymbolTable(), self.GetString()[startIndex:targetIndex]), nil
	} else {
		return nil, errors.NewTypeError(indexObject.TypeName(), IntegerName, TupleName)
	}
}

func StringToInteger(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	content, base, valid := ParseInteger(self.GetString())
	if !valid {
		return nil, errors.NewInvalidIntegerDefinition(errors.UnknownLine, self.GetString())
	}
	number, parsingError := strconv.ParseInt(content, base, 64)
	if parsingError != nil {
		return nil, errors.NewInvalidIntegerDefinition(errors.UnknownLine, self.GetString())
	}
	return NewInteger(vm.PeekSymbolTable(), number), nil
}

func StringToFloat(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	content, valid := ParseFloat(self.GetString())
	if !valid {
		return nil, errors.NewInvalidFloatDefinition(errors.UnknownLine, self.GetString())
	}
	number, parsingError := strconv.ParseFloat(content, 64)
	if parsingError != nil {
		return nil, errors.NewInvalidFloatDefinition(errors.UnknownLine, self.GetString())
	}
	return NewFloat(vm.PeekSymbolTable(), number), nil
}

func StringToArray(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	var content []IObject
	for _, char := range self.GetString() {
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
	for _, char := range self.GetString() {
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
		stringHash, hashingError := vm.HashString(fmt.Sprintf("%s-%s", self.GetString(), StringName))
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
	}
	string_.SetString(value)
	string_.SetLength(len(value))
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
	for _, byte_ := range self.GetBytes() {
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
	for _, byte_ := range right.GetBytes() {
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
	for _, byte_ := range left.GetBytes() {
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
	for _, byte_ := range self.GetBytes() {
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
	for i = 0; i < right.GetInteger64(); i++ {
		for _, byte_ := range self.GetBytes() {
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
	for i = 0; i < left.GetInteger64(); i++ {
		for _, byte_ := range self.GetBytes() {
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
	if self.GetLength() != right.GetLength() {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftEquals IObject
	for index, byte_ := range self.GetBytes() {
		leftEquals, getError = byte_.Get(Equals)
		if getError != nil {
			return nil, getError
		}
		if _, ok := leftEquals.(*Function); !ok {
			return nil, errors.NewTypeError(leftEquals.TypeName(), FunctionName)
		}
		result, callError := CallFunction(leftEquals.(*Function), vm, vm.PeekSymbolTable(), right.GetBytes()[index])
		if callError != nil {
			return nil, callError
		}
		if _, ok := result.(*Bool); !ok {
			return nil, errors.NewTypeError(result.TypeName(), BoolName)
		}
		if !result.GetBool() {
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
	if left.GetLength() != self.GetLength() {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftEquals IObject
	for index, byte_ := range left.GetBytes() {
		leftEquals, getError = byte_.Get(RightEquals)
		if getError != nil {
			return nil, getError
		}
		if _, ok := leftEquals.(*Function); !ok {
			return nil, errors.NewTypeError(leftEquals.TypeName(), FunctionName)
		}
		result, callError := CallFunction(leftEquals.(*Function), vm, vm.PeekSymbolTable(), self.GetBytes()[index])
		if callError != nil {
			return nil, callError
		}
		if _, ok := result.(*Bool); !ok {
			return nil, errors.NewTypeError(result.TypeName(), BoolName)
		}
		if !result.GetBool() {
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
	if self.GetLength() != right.GetLength() {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftEquals IObject
	for index, byte_ := range self.GetBytes() {
		leftEquals, getError = byte_.Get(NotEquals)
		if getError != nil {
			return nil, getError
		}
		if _, ok := leftEquals.(*Function); !ok {
			return nil, errors.NewTypeError(leftEquals.TypeName(), FunctionName)
		}
		result, callError := CallFunction(leftEquals.(*Function), vm, vm.PeekSymbolTable(), right.GetBytes()[index])
		if callError != nil {
			return nil, callError
		}
		if _, ok := result.(*Bool); !ok {
			return nil, errors.NewTypeError(result.TypeName(), BoolName)
		}
		if !result.GetBool() {
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
	if left.GetLength() != self.GetLength() {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftEquals IObject
	for index, byte_ := range left.GetBytes() {
		leftEquals, getError = byte_.Get(RightNotEquals)
		if getError != nil {
			return nil, getError
		}
		if _, ok := leftEquals.(*Function); !ok {
			return nil, errors.NewTypeError(leftEquals.TypeName(), FunctionName)
		}
		result, callError := CallFunction(leftEquals.(*Function), vm, vm.PeekSymbolTable(), self.GetBytes()[index])
		if callError != nil {
			return nil, callError
		}
		if _, ok := result.(*Bool); !ok {
			return nil, errors.NewTypeError(result.TypeName(), BoolName)
		}
		if !result.GetBool() {
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
	for index, byte_ := range self.GetBytes() {
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
	for _, byte_ := range self.GetBytes() {
		newBytes = append(newBytes, NewInteger(vm.PeekSymbolTable(), byte_.Integer64))
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
		index, calcError := CalcIndex(indexObject, self.GetLength())
		if calcError != nil {
			return nil, calcError
		}
		return self.GetBytes()[index], nil
	} else if _, ok = indexObject.(*Tuple); ok {
		if len(indexObject.GetContent()) != 2 {
			return nil, errors.NewInvalidNumberOfArguments(len(indexObject.GetContent()), 2)
		}
		startIndex, calcError := CalcIndex(indexObject.GetContent()[0], self.GetLength())
		if calcError != nil {
			return nil, calcError
		}
		targetIndex, calcError := CalcIndex(indexObject.GetContent()[1], self.GetLength())
		if calcError != nil {
			return nil, calcError
		}
		return NewBytes(vm.PeekSymbolTable(), self.GetBytes()[startIndex:targetIndex]), nil
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
	for _, byte_ := range self.GetBytes() {
		numbers = append(numbers, uint8(byte_.Integer64))
	}
	return NewInteger(vm.PeekSymbolTable(), int64(binary.BigEndian.Uint32(numbers))), nil
}

func BytesToString(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	var numbers []uint8
	for _, byte_ := range self.GetBytes() {
		numbers = append(numbers, uint8(byte_.Integer64))
	}
	return NewString(vm.PeekSymbolTable(), string(numbers)), nil
}

func BytesToBool(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewBool(vm.PeekSymbolTable(), self.GetLength() != 0), nil
}

func BytesToArray(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	var newContent []IObject
	var byteCopyFunc IObject
	for _, byte_ := range self.GetBytes() {
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
	for _, byte_ := range self.GetBytes() {
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
		Object: NewObject(BytesName, nil, parent),
	}
	bytes_.SetBytes(content)
	bytes_.SetLength(len(content))
	BytesInitialize(nil, bytes_)
	return bytes_
}

type Integer struct {
	*Object
}

func IntegerToBool(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewBool(vm.PeekSymbolTable(), self.GetInteger64() != 0), nil
}

func IntegerCopy(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewInteger(vm.PeekSymbolTable(), self.GetInteger64()), nil
}

func IntegerToString(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewString(vm.PeekSymbolTable(), fmt.Sprint(self.GetInteger64())), nil
}

func IntegerToFloat(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewFloat(vm.PeekSymbolTable(), float64(self.GetInteger64())), nil
}

func IntegerNegBits(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewInteger(vm.PeekSymbolTable(), ^self.(Integer).Integer64), nil
}

func IntegerAdd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	switch right.(type) {
	case *Integer:
		return NewInteger(vm.PeekSymbolTable(), self.GetInteger64()+right.GetInteger64()), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), float64(self.GetInteger64())+right.GetFloat64()), nil
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
		return NewInteger(vm.PeekSymbolTable(), left.GetInteger64()+self.GetInteger64()), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), left.GetFloat64()+float64(self.GetInteger64())), nil
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
		return NewInteger(vm.PeekSymbolTable(), self.GetInteger64()-right.GetInteger64()), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), float64(self.GetInteger64())-right.GetFloat64()), nil
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
		return NewInteger(vm.PeekSymbolTable(), left.GetInteger64()-self.GetInteger64()), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), left.GetFloat64()-float64(self.GetInteger64())), nil
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
		return NewInteger(vm.PeekSymbolTable(), self.GetInteger64()*right.GetInteger64()), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), float64(self.GetInteger64())*right.GetFloat64()), nil
	case *String:
		return NewString(vm.PeekSymbolTable(), Repeat(right.GetString(), self.GetInteger64())), nil
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
		return NewInteger(vm.PeekSymbolTable(), left.GetInteger64()*self.GetInteger64()), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), left.GetFloat64()*float64(self.GetInteger64())), nil
	case *String:
		return NewString(vm.PeekSymbolTable(), Repeat(left.GetString(), self.GetInteger64())), nil
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
		return NewInteger(vm.PeekSymbolTable(), self.GetInteger64()/right.GetInteger64()), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), float64(self.GetInteger64())/right.GetFloat64()), nil
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
		return NewInteger(vm.PeekSymbolTable(), left.GetInteger64()/self.GetInteger64()), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), left.GetFloat64()/float64(self.GetInteger64())), nil
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
		return NewInteger(vm.PeekSymbolTable(), self.GetInteger64()%right.GetInteger64()), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), math.Mod(float64(self.GetInteger64()), right.GetFloat64())), nil
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
		return NewInteger(vm.PeekSymbolTable(), left.GetInteger64()%self.GetInteger64()), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), math.Mod(left.GetFloat64(), float64(self.GetInteger64()))), nil
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
		return NewFloat(vm.PeekSymbolTable(), math.Pow(float64(self.GetInteger64()), float64(right.GetInteger64()))), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), math.Pow(float64(self.GetInteger64()), right.GetFloat64())), nil
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
		return NewFloat(vm.PeekSymbolTable(), math.Pow(float64(left.GetInteger64()), float64(self.GetInteger64()))), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), math.Pow(left.GetFloat64(), float64(self.GetInteger64()))), nil
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
	return NewInteger(vm.PeekSymbolTable(), self.GetInteger64()^right.GetInteger64()), nil
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
	return NewInteger(vm.PeekSymbolTable(), left.GetInteger64()^self.GetInteger64()), nil
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
	return NewInteger(vm.PeekSymbolTable(), self.GetInteger64()&right.GetInteger64()), nil
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
	return NewInteger(vm.PeekSymbolTable(), left.GetInteger64()&self.GetInteger64()), nil
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
	return NewInteger(vm.PeekSymbolTable(), self.GetInteger64()|right.GetInteger64()), nil
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
	return NewInteger(vm.PeekSymbolTable(), left.GetInteger64()|self.GetInteger64()), nil
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
	return NewInteger(vm.PeekSymbolTable(), self.GetInteger64()<<right.GetInteger64()), nil
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
	return NewInteger(vm.PeekSymbolTable(), left.GetInteger64()<<self.GetInteger64()), nil
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
	return NewInteger(vm.PeekSymbolTable(), self.GetInteger64()>>right.GetInteger64()), nil
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
	return NewInteger(vm.PeekSymbolTable(), left.GetInteger64()>>self.GetInteger64()), nil
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
		floatRight = float64(right.GetInteger64())
	case *Float:
		floatRight = right.GetFloat64()
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), float64(self.GetInteger64()) == floatRight), nil
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
		floatLeft = float64(left.GetInteger64())
	case *Float:
		floatLeft = left.GetFloat64()
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft == float64(self.GetInteger64())), nil
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
		floatRight = float64(right.GetInteger64())
	case *Float:
		floatRight = right.GetFloat64()
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), float64(self.GetInteger64()) != floatRight), nil
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
		floatLeft = float64(left.GetInteger64())
	case *Float:
		floatLeft = left.GetFloat64()
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft != float64(self.GetInteger64())), nil
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
		floatRight = float64(right.GetInteger64())
	case *Float:
		floatRight = right.GetFloat64()
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), float64(self.GetInteger64()) > floatRight), nil
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
		floatLeft = float64(left.GetInteger64())
	case *Float:
		floatLeft = left.GetFloat64()
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft > float64(self.GetInteger64())), nil
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
		floatRight = float64(right.GetInteger64())
	case *Float:
		floatRight = right.GetFloat64()
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), float64(self.GetInteger64()) < floatRight), nil
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
		floatLeft = float64(left.GetInteger64())
	case *Float:
		floatLeft = left.GetFloat64()
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft < float64(self.GetInteger64())), nil
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
		floatRight = float64(right.GetInteger64())
	case *Float:
		floatRight = right.GetFloat64()
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), float64(self.GetInteger64()) >= floatRight), nil
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
		floatLeft = float64(left.GetInteger64())
	case *Float:
		floatLeft = left.GetFloat64()
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft >= float64(self.GetInteger64())), nil
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
		floatRight = float64(right.GetInteger64())
	case *Float:
		floatRight = right.GetFloat64()
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), float64(self.GetInteger64()) <= floatRight), nil
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
		floatLeft = float64(left.GetInteger64())
	case *Float:
		floatLeft = left.GetFloat64()
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft <= float64(self.GetInteger64())), nil
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
	}
	integer.SetInteger64(value)
	IntegerInitialize(nil, integer)
	return integer
}

type Float struct {
	*Object
}

func FloatToBool(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewBool(vm.PeekSymbolTable(), self.GetFloat64() != 0), nil
}

func FloatCopy(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewFloat(vm.PeekSymbolTable(), self.GetFloat64()), nil
}

func FloatToString(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewString(vm.PeekSymbolTable(), fmt.Sprint(self.GetFloat64())), nil
}

func FloatToInteger(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewInteger(vm.PeekSymbolTable(), int64(self.GetFloat64())), nil
}

func FloatAdd(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	right := arguments[0]
	switch right.(type) {
	case *Integer:
		return NewFloat(vm.PeekSymbolTable(), self.GetFloat64()+float64(right.GetInteger64())), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), self.GetFloat64()+right.GetFloat64()), nil
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
		return NewFloat(vm.PeekSymbolTable(), float64(left.GetInteger64())+self.GetFloat64()), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), left.GetFloat64()+self.GetFloat64()), nil
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
		return NewFloat(vm.PeekSymbolTable(), self.GetFloat64()-float64(right.GetInteger64())), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), self.GetFloat64()-right.GetFloat64()), nil
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
		return NewFloat(vm.PeekSymbolTable(), float64(left.GetInteger64())-self.GetFloat64()), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), left.GetFloat64()-self.GetFloat64()), nil
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
		return NewFloat(vm.PeekSymbolTable(), self.GetFloat64()*float64(right.GetInteger64())), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), self.GetFloat64()*right.GetFloat64()), nil
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
		return NewFloat(vm.PeekSymbolTable(), float64(left.GetInteger64())*self.GetFloat64()), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), left.GetFloat64()*self.GetFloat64()), nil
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
		return NewFloat(vm.PeekSymbolTable(), self.GetFloat64()/float64(right.GetInteger64())), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), self.GetFloat64()/right.GetFloat64()), nil
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
		return NewFloat(vm.PeekSymbolTable(), float64(left.GetInteger64())/self.GetFloat64()), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), left.GetFloat64()/self.GetFloat64()), nil
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
		return NewFloat(vm.PeekSymbolTable(), math.Mod(self.GetFloat64(), float64(right.GetInteger64()))), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), math.Mod(self.GetFloat64(), right.GetFloat64())), nil
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
		return NewFloat(vm.PeekSymbolTable(), math.Mod(float64(left.GetInteger64()), self.GetFloat64())), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), math.Mod(left.GetFloat64(), self.GetFloat64())), nil
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
		return NewFloat(vm.PeekSymbolTable(), math.Pow(self.GetFloat64(), float64(right.GetInteger64()))), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), math.Pow(self.GetFloat64(), right.GetFloat64())), nil
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
		return NewFloat(vm.PeekSymbolTable(), math.Pow(float64(left.GetInteger64()), self.GetFloat64())), nil
	case *Float:
		return NewFloat(vm.PeekSymbolTable(), math.Pow(left.GetFloat64(), self.GetFloat64())), nil
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
		floatRight = float64(right.GetInteger64())
	case *Float:
		floatRight = right.GetFloat64()
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), self.GetFloat64() == floatRight), nil
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
		floatLeft = float64(left.GetInteger64())
	case *Float:
		floatLeft = left.GetFloat64()
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft == self.GetFloat64()), nil
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
		floatRight = float64(right.GetInteger64())
	case *Float:
		floatRight = right.GetFloat64()
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), self.GetFloat64() != floatRight), nil
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
		floatLeft = float64(left.GetInteger64())
	case *Float:
		floatLeft = left.GetFloat64()
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft != self.GetFloat64()), nil
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
		floatRight = float64(right.GetInteger64())
	case *Float:
		floatRight = right.GetFloat64()
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), self.GetFloat64() > floatRight), nil
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
		floatLeft = float64(left.GetInteger64())
	case *Float:
		floatLeft = left.GetFloat64()
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft > self.GetFloat64()), nil
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
		floatRight = float64(right.GetInteger64())
	case *Float:
		floatRight = right.GetFloat64()
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), self.GetFloat64() < floatRight), nil
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
		floatLeft = float64(left.GetInteger64())
	case *Float:
		floatLeft = left.GetFloat64()
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft < self.GetFloat64()), nil
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
		floatRight = float64(right.GetInteger64())
	case *Float:
		floatRight = right.GetFloat64()
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), self.GetFloat64() >= floatRight), nil
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
		floatLeft = float64(left.GetInteger64())
	case *Float:
		floatLeft = left.GetFloat64()
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft >= self.GetFloat64()), nil
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
		floatRight = float64(right.GetInteger64())
	case *Float:
		floatRight = right.GetFloat64()
	default:
		return nil, errors.NewTypeError(right.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), self.GetFloat64() <= floatRight), nil
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
		floatLeft = float64(left.GetInteger64())
	case *Float:
		floatLeft = left.GetFloat64()
	default:
		return nil, errors.NewTypeError(left.TypeName(), IntegerName, FloatName)
	}
	return NewBool(vm.PeekSymbolTable(), floatLeft <= self.GetFloat64()), nil
}

func FloatHash(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	if self.GetHash() == 0 {
		floatHash, hashingError := vm.HashString(fmt.Sprintf("%f-%s", self.GetFloat64(), FloatName))
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
	}
	float_.SetFloat64(value)
	FloatInitialize(nil, float_)
	return float_
}

type Array struct {
	*Object
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
	if self.GetLength() != right.GetLength() {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftEquals IObject
	var rightEquals IObject
	var comparisonResult IObject
	var callError *errors.Error
	var comparisonResultToBool IObject
	var comparisonBool IObject

	for i := 0; i < self.GetLength(); i++ {
		leftEquals, getError = self.GetContent()[i].Get(Equals)
		if getError != nil {
			rightEquals, getError = right.GetContent()[i].Get(RightEquals)
			if getError != nil {
				return nil, getError
			}
			if _, ok := rightEquals.(*Function); !ok {
				return nil, errors.NewTypeError(rightEquals.TypeName(), FunctionName)
			}
			comparisonResult, callError = CallFunction(rightEquals.(*Function), vm, vm.PeekSymbolTable(), self.GetContent()[i])
		} else {
			comparisonResult, callError = CallFunction(leftEquals.(*Function), vm, vm.PeekSymbolTable(), right.GetContent()[i])
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
		if !comparisonBool.GetBool() {
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
	if self.GetLength() != left.GetLength() {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftEquals IObject
	var rightEquals IObject
	var comparisonResult IObject
	var callError *errors.Error
	var comparisonResultToBool IObject
	var comparisonBool IObject

	for i := 0; i < self.GetLength(); i++ {
		leftEquals, getError = left.GetContent()[i].Get(Equals)
		if getError != nil {
			rightEquals, getError = self.GetContent()[i].Get(RightEquals)
			if getError != nil {
				return nil, getError
			}
			if _, ok := rightEquals.(*Function); !ok {
				return nil, errors.NewTypeError(rightEquals.TypeName(), FunctionName)
			}
			comparisonResult, callError = CallFunction(rightEquals.(*Function), vm, vm.PeekSymbolTable(), left.GetContent()[i])
		} else {
			comparisonResult, callError = CallFunction(leftEquals.(*Function), vm, vm.PeekSymbolTable(), self.GetContent()[i])
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
		if !comparisonBool.GetBool() {
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
	if self.GetLength() != right.GetLength() {
		return NewBool(vm.PeekSymbolTable(), true), nil
	}
	var leftNotEquals IObject
	var rightNotEquals IObject
	var comparisonResult IObject
	var callError *errors.Error
	var comparisonResultToBool IObject
	var comparisonBool IObject

	for i := 0; i < self.GetLength(); i++ {
		leftNotEquals, getError = self.GetContent()[i].Get(NotEquals)
		if getError != nil {
			rightNotEquals, getError = right.GetContent()[i].Get(RightNotEquals)
			if getError != nil {
				return nil, getError
			}
			if _, ok := rightNotEquals.(*Function); !ok {
				return nil, errors.NewTypeError(rightNotEquals.TypeName(), FunctionName)
			}
			comparisonResult, callError = CallFunction(rightNotEquals.(*Function), vm, vm.PeekSymbolTable(), self.GetContent()[i])
		} else {
			comparisonResult, callError = CallFunction(leftNotEquals.(*Function), vm, vm.PeekSymbolTable(), right.GetContent()[i])
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
		if !comparisonBool.GetBool() {
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
	if self.GetLength() != left.GetLength() {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftEquals IObject
	var rightEquals IObject
	var comparisonResult IObject
	var callError *errors.Error
	var comparisonResultToBool IObject
	var comparisonBool IObject

	for i := 0; i < self.GetLength(); i++ {
		leftEquals, getError = left.GetContent()[i].Get(NotEquals)
		if getError != nil {
			rightEquals, getError = self.GetContent()[i].Get(RightNotEquals)
			if getError != nil {
				return nil, getError
			}
			if _, ok := rightEquals.(*Function); !ok {
				return nil, errors.NewTypeError(rightEquals.TypeName(), FunctionName)
			}
			comparisonResult, callError = CallFunction(rightEquals.(*Function), vm, vm.PeekSymbolTable(), left.GetContent()[i])
		} else {
			comparisonResult, callError = CallFunction(leftEquals.(*Function), vm, vm.PeekSymbolTable(), self.GetContent()[i])
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
		if !comparisonBool.GetBool() {
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
	for _, object := range self.GetContent() {
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
		index, calcError := CalcIndex(indexObject, self.GetLength())
		if calcError != nil {
			return nil, calcError
		}
		return self.GetContent()[index], nil
	} else if _, ok = indexObject.(*Tuple); ok {
		if len(indexObject.GetContent()) != 2 {
			return nil, errors.NewInvalidNumberOfArguments(len(indexObject.GetContent()), 2)
		}
		startIndex, calcError := CalcIndex(indexObject.GetContent()[0], self.GetLength())
		if calcError != nil {
			return nil, calcError
		}
		targetIndex, calcError := CalcIndex(indexObject.GetContent()[1], self.GetLength())
		if calcError != nil {
			return nil, calcError
		}
		return NewArray(vm.PeekSymbolTable(), self.GetContent()[startIndex:targetIndex]), nil
	} else {
		return nil, errors.NewTypeError(indexObject.TypeName(), IntegerName, TupleName)
	}
}

func ArrayAssign(vm VirtualMachine, arguments ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	index, calcError := CalcIndex(arguments[0], self.GetLength())
	if calcError != nil {
		return nil, calcError
	}
	self.GetContent()[index] = arguments[1]
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
	for index, object := range self.GetContent() {
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
		result += objectString.GetString()
	}
	return NewString(vm.PeekSymbolTable(), result+")"), nil
}

func ArrayToBool(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewBool(vm.PeekSymbolTable(), self.GetLength() != 0), nil
}

func ArrayToArray(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewArray(vm.PeekSymbolTable(), append([]IObject{}, self.GetContent()...)), nil
}

func ArrayToTuple(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewTuple(vm.PeekSymbolTable(), append([]IObject{}, self.GetContent()...)), nil
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
		Object: NewObject(ArrayName, nil, parentSymbols),
	}
	array.SetContent(content)
	array.SetLength(len(content))
	ArrayInitialize(nil, array)
	return array
}

type Tuple struct {
	*Object
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
	if self.GetLength() != right.GetLength() {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftEquals IObject
	var rightEquals IObject
	var comparisonResult IObject
	var callError *errors.Error
	var comparisonResultToBool IObject
	var comparisonBool IObject

	for i := 0; i < self.GetLength(); i++ {
		leftEquals, getError = self.GetContent()[i].Get(Equals)
		if getError != nil {
			rightEquals, getError = right.GetContent()[i].Get(RightEquals)
			if getError != nil {
				return nil, getError
			}
			if _, ok := rightEquals.(*Function); !ok {
				return nil, errors.NewTypeError(rightEquals.TypeName(), FunctionName)
			}
			comparisonResult, callError = CallFunction(rightEquals.(*Function), vm, vm.PeekSymbolTable(), self.GetContent()[i])
		} else {
			comparisonResult, callError = CallFunction(leftEquals.(*Function), vm, vm.PeekSymbolTable(), right.GetContent()[i])
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
		if !comparisonBool.GetBool() {
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
	if self.GetLength() != left.GetLength() {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftEquals IObject
	var rightEquals IObject
	var comparisonResult IObject
	var callError *errors.Error
	var comparisonResultToBool IObject
	var comparisonBool IObject

	for i := 0; i < self.GetLength(); i++ {
		leftEquals, getError = left.GetContent()[i].Get(Equals)
		if getError != nil {
			rightEquals, getError = self.GetContent()[i].Get(RightEquals)
			if getError != nil {
				return nil, getError
			}
			if _, ok := rightEquals.(*Function); !ok {
				return nil, errors.NewTypeError(rightEquals.TypeName(), FunctionName)
			}
			comparisonResult, callError = CallFunction(rightEquals.(*Function), vm, vm.PeekSymbolTable(), left.GetContent()[i])
		} else {
			comparisonResult, callError = CallFunction(leftEquals.(*Function), vm, vm.PeekSymbolTable(), self.GetContent()[i])
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
		if !comparisonBool.GetBool() {
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
	if self.GetLength() != right.GetLength() {
		return NewBool(vm.PeekSymbolTable(), true), nil
	}
	var leftNotEquals IObject
	var rightNotEquals IObject
	var comparisonResult IObject
	var callError *errors.Error
	var comparisonResultToBool IObject
	var comparisonBool IObject

	for i := 0; i < self.GetLength(); i++ {
		leftNotEquals, getError = self.GetContent()[i].Get(NotEquals)
		if getError != nil {
			rightNotEquals, getError = right.GetContent()[i].Get(RightNotEquals)
			if getError != nil {
				return nil, getError
			}
			if _, ok := rightNotEquals.(*Function); !ok {
				return nil, errors.NewTypeError(rightNotEquals.TypeName(), FunctionName)
			}
			comparisonResult, callError = CallFunction(rightNotEquals.(*Function), vm, vm.PeekSymbolTable(), self.GetContent()[i])
		} else {
			comparisonResult, callError = CallFunction(leftNotEquals.(*Function), vm, vm.PeekSymbolTable(), right.GetContent()[i])
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
		if !comparisonBool.GetBool() {
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
	if self.GetLength() != left.GetLength() {
		return NewBool(vm.PeekSymbolTable(), false), nil
	}
	var leftEquals IObject
	var rightEquals IObject
	var comparisonResult IObject
	var callError *errors.Error
	var comparisonResultToBool IObject
	var comparisonBool IObject

	for i := 0; i < self.GetLength(); i++ {
		leftEquals, getError = left.GetContent()[i].Get(NotEquals)
		if getError != nil {
			rightEquals, getError = self.GetContent()[i].Get(RightNotEquals)
			if getError != nil {
				return nil, getError
			}
			if _, ok := rightEquals.(*Function); !ok {
				return nil, errors.NewTypeError(rightEquals.TypeName(), FunctionName)
			}
			comparisonResult, callError = CallFunction(rightEquals.(*Function), vm, vm.PeekSymbolTable(), left.GetContent()[i])
		} else {
			comparisonResult, callError = CallFunction(leftEquals.(*Function), vm, vm.PeekSymbolTable(), self.GetContent()[i])
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
		if !comparisonBool.GetBool() {
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
	for _, object := range self.GetContent() {
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
		index, calcError := CalcIndex(indexObject, self.GetLength())
		if calcError != nil {
			return nil, calcError
		}
		return self.GetContent()[index], nil
	} else if _, ok = indexObject.(*Tuple); ok {
		if len(indexObject.GetContent()) != 2 {
			return nil, errors.NewInvalidNumberOfArguments(len(indexObject.GetContent()), 2)
		}
		startIndex, calcError := CalcIndex(indexObject.GetContent()[0], self.GetLength())
		if calcError != nil {
			return nil, calcError
		}
		targetIndex, calcError := CalcIndex(indexObject.GetContent()[1], self.GetLength())
		if calcError != nil {
			return nil, calcError
		}
		return NewTuple(vm.PeekSymbolTable(), self.GetContent()[startIndex:targetIndex]), nil
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
	for index, object := range self.GetContent() {
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
		result += objectString.GetString()
	}
	return NewString(vm.PeekSymbolTable(), result+")"), nil
}

func TupleToBool(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewBool(vm.PeekSymbolTable(), self.GetLength() != 0), nil
}

func TupleToArray(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewArray(vm.PeekSymbolTable(), append([]IObject{}, self.GetContent()...)), nil
}

func TupleToTuple(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewTuple(vm.PeekSymbolTable(), append([]IObject{}, self.GetContent()...)), nil
}
func TupleHash(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	var tupleHash int64 = 0
	var objectHashFunc IObject
	for _, object := range self.GetContent() {
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
			tupleHash = objectHash.GetInteger64()
		} else {
			tupleHash <<= objectHash.GetInteger64()
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
		Object: NewObject(TupleName, nil, parentSymbols),
	}
	tuple.SetContent(content)
	tuple.SetLength(len(content))
	TupleInitialize(nil, tuple)
	return tuple
}

type Bool struct {
	*Object
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
	return NewBool(vm.PeekSymbolTable(), self.GetBool() == right.GetBool()), nil
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
	return NewBool(vm.PeekSymbolTable(), left.GetBool() == self.GetBool()), nil
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
	return NewBool(vm.PeekSymbolTable(), self.GetBool() != right.GetBool()), nil
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
	return NewBool(vm.PeekSymbolTable(), left.GetBool() != self.GetBool()), nil
}

func BoolCopy(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	return NewBool(vm.PeekSymbolTable(), self.GetBool()), nil
}

func BoolToInteger(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	if self.GetBool() {
		return NewInteger(vm.PeekSymbolTable(), 1), nil
	}
	return NewInteger(vm.PeekSymbolTable(), 0), nil
}

func BoolToFloat(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	if self.GetBool() {
		return NewFloat(vm.PeekSymbolTable(), 1), nil
	}
	return NewFloat(vm.PeekSymbolTable(), 0), nil
}

func BoolToString(vm VirtualMachine, _ ...IObject) (IObject, *errors.Error) {
	self, getError := vm.PeekSymbolTable().GetSelf(Self)
	if getError != nil {
		return nil, getError
	}
	if self.GetBool() {
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
		boolHash, hashingError := vm.HashString(fmt.Sprintf("%t-%s", self.GetBool(), BoolName))
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
	}
	bool_.SetBool(value)
	BoolInitialize(nil, bool_)
	return bool_
}

type KeyValue struct {
	Key   IObject
	Value IObject
}
type HashTable struct {
	*Object
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
	if self.GetLength() != right.Length {
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
	for key, leftValue := range self.GetKeyValues() {
		// Check if other has the key
		rightValue, ok := right.KeyValues[key]
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
	if self.GetLength() != left.Length {
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
	for key, leftValue := range left.KeyValues {
		// Check if other has the key
		rightValue, ok := self.GetKeyValues()[key]
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
	if self.GetLength() != right.Length {
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
	for key, leftValue := range self.GetKeyValues() {
		// Check if other has the key
		rightValue, ok := right.KeyValues[key]
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
	if self.GetLength() != left.Length {
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
	for key, leftValue := range left.KeyValues {
		// Check if other has the key
		rightValue, ok := self.GetKeyValues()[key]
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
	keyValues, found := self.GetKeyValues()[indexHash.GetInteger64()]
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
		if equals.GetBool() {
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
	keyValues, found := self.GetKeyValues()[indexHash.GetInteger64()]
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
		if equals.GetBool() {
			self.GetKeyValues()[indexHash.GetInteger64()][index].Value = newValue
			self.IncreaseLength()
			return vm.PeekSymbolTable().GetAny(None)
		}
	}
	self.IncreaseLength()
	self.GetKeyValues()[indexHash.GetInteger64()] = append(self.GetKeyValues()[indexHash.GetInteger64()], &KeyValue{
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
	for _, keyValues := range self.GetKeyValues() {
		for _, keyValue := range keyValues {
			keyToString, getError = keyValue.Key.Get(ToString)
			if getError != nil {
				return nil, getError
			}
			keyString, callError = CallFunction(keyToString.(*Function), vm, keyValue.Key.SymbolTable())
			if callError != nil {
				return nil, callError
			}
			result += keyString.GetString()
			valueToString, getError = keyValue.Value.Get(ToString)
			if getError != nil {
				return nil, getError
			}
			valueString, callError = CallFunction(valueToString.(*Function), vm, keyValue.Value.SymbolTable())
			if callError != nil {
				return nil, callError
			}
			result += ":" + valueString.GetString() + ","
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
	if self.GetLength() > 0 {
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
	for _, keyValues := range self.GetKeyValues() {
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
	for _, keyValues := range self.GetKeyValues() {
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
		Object: NewObject(HashName, nil, parent),
	}
	hashTable.SetKeyValues(entries)
	hashTable.SetLength(entriesLength)
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
	None 	   - (Done)
	True	   - (Done)
	False	   - (Done)
	// Functions
	Hash       - ()
	Id         - ()
	Range      - ()
	Len        - ()
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
	symbolTable.Set(TrueName,
		NewBool(nil, true),
	)
	symbolTable.Set(FalseName,
		NewBool(nil, false),
	)
	// Functions

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
