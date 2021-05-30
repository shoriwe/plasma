package vm

import (
	"crypto/rand"
	"github.com/shoriwe/gruby/pkg/errors"
	"hash"
	"hash/crc32"
	"math/big"
)

const (
	polySize = 4294967295
)

type Plasma struct {
	Code        *CodeStack
	MemoryStack *ObjectStack
	Context     *SymbolStack
	Crc32Hash   hash.Hash32
	seed        uint64
}

func (p *Plasma) PushObject(object IObject) {
	p.MemoryStack.Push(object)
}
func (p *Plasma) PeekObject() IObject {
	return p.MemoryStack.Peek()
}
func (p *Plasma) PopObject() IObject {
	return p.MemoryStack.Pop()
}

func (p *Plasma) PushSymbolTable(table *SymbolTable) {
	p.Context.Push(table)
}

func (p *Plasma) PopSymbolTable() *SymbolTable {
	return p.Context.Pop()
}

func (p *Plasma) Initialize(code []Code) *errors.Error {
	p.PushCode(NewBytecodeFromArray(code))
	p.MemoryStack.Clear()
	p.Context.Clear()
	p.Context.Push(SetDefaultSymbolTable())
	return nil
}

func (p *Plasma) PeekSymbolTable() *SymbolTable {
	return p.Context.Peek()
}

func (p *Plasma) newStringOP(code Code) *errors.Error {
	value := code.Value.(string)
	stringObject := NewString(p.Context.Peek(), value)
	p.PushObject(stringObject)
	return nil
}

func (p *Plasma) newBytesOP(code Code) *errors.Error {
	value := code.Value.([]byte)
	bytesObject := NewBytes(p.Context.Peek(), value)
	p.PushObject(bytesObject)
	return nil
}

func (p *Plasma) newIntegerOP(code Code) *errors.Error {
	value := code.Value.(int64)
	integer := NewInteger(p.Context.Peek(), value)
	p.PushObject(integer)
	return nil
}

func (p *Plasma) newFloatOP(code Code) *errors.Error {
	value := code.Value.(float64)
	float := NewFloat(p.Context.Peek(), value)
	p.PushObject(float)
	return nil
}

func (p *Plasma) newTrueBoolOP() *errors.Error {
	p.PushObject(NewBool(p.PeekSymbolTable(), true))
	return nil
}

func (p *Plasma) newFalseBoolOP() *errors.Error {
	p.PushObject(NewBool(p.PeekSymbolTable(), false))
	return nil
}

func (p *Plasma) getNoneOP() *errors.Error {
	none, getError := p.PeekSymbolTable().GetAny(None)
	if getError != nil {
		return getError
	}
	p.PushObject(none)
	return nil
}

func (p *Plasma) newTupleOP(code Code) *errors.Error {
	numberOfValues := code.Value.(int)
	var values []IObject
	for i := 0; i < numberOfValues; i++ {
		if !p.MemoryStack.HasNext() {
			return errors.NewInvalidNumberOfArguments(i, numberOfValues)
		}
		values = append(values, p.PopObject())
	}
	p.PushObject(NewTuple(p.PeekSymbolTable(), values))
	return nil
}

func (p *Plasma) newArrayOP(code Code) *errors.Error {
	numberOfValues := code.Value.(int)
	var values []IObject
	for i := 0; i < numberOfValues; i++ {
		if !p.MemoryStack.HasNext() {
			return errors.NewInvalidNumberOfArguments(i, numberOfValues)
		}
		values = append(values, p.PopObject())
	}
	p.PushObject(NewArray(p.PeekSymbolTable(), values))
	return nil
}

func (p *Plasma) newHashOP(code Code) *errors.Error {
	numberOfValues := code.Value.(int)
	var keyValues []*KeyValue
	for i := 0; i < numberOfValues; i++ {
		if !p.MemoryStack.HasNext() {
			return errors.NewInvalidNumberOfArguments(i, numberOfValues)
		}
		key := p.PopObject()
		if !p.MemoryStack.HasNext() {
			return errors.NewInvalidNumberOfArguments(i, numberOfValues)
		}
		value := p.PopObject()
		keyValues = append(keyValues, &KeyValue{
			Key:   key,
			Value: value,
		})
	}
	hashTable := NewHashTable(p.PeekSymbolTable(), map[int64][]*KeyValue{}, numberOfValues)
	hashTableAssign, getError := hashTable.Get(Assign)
	if getError != nil {
		return getError
	}
	for _, keyValue := range keyValues {
		_, assignError := CallFunction(hashTableAssign.(*Function), p, hashTable.SymbolTable(), keyValue.Key, keyValue.Value)
		if assignError != nil {
			return assignError
		}
	}
	p.PushObject(hashTable)
	return nil
}

// Useful function to call those built ins that doesn't receive arguments of an object
func (p *Plasma) noArgsGetAndCall(operationName string) *errors.Error {
	object := p.PopObject()
	operation, getError := object.Get(operationName)
	if getError != nil {
		return getError
	}
	if _, ok := operation.(*Function); !ok {
		return errors.NewTypeError(operation.TypeName(), FunctionName)
	}
	result, callError := CallFunction(operation.(*Function), p, object.SymbolTable())
	if callError != nil {
		return callError
	}
	p.PushObject(result)
	return nil
}

// Function useful to cal object built-ins binary expression functions
func (p *Plasma) leftBinaryExpressionFuncCall(operationName string) *errors.Error {
	leftHandSide := p.PopObject()
	rightHandSide := p.PopObject()
	operation, getError := leftHandSide.Get(operationName)
	if getError != nil {
		return p.rightBinaryExpressionFuncCall(leftHandSide, rightHandSide, operationName)
	}
	if _, ok := operation.(*Function); !ok {
		return p.rightBinaryExpressionFuncCall(leftHandSide, rightHandSide, operationName)
	}
	result, callError := CallFunction(operation.(*Function), p, leftHandSide.SymbolTable(), rightHandSide)
	if callError != nil {
		return p.rightBinaryExpressionFuncCall(leftHandSide, rightHandSide, operationName)
	}
	p.PushObject(result)
	return nil
}

func (p *Plasma) rightBinaryExpressionFuncCall(leftHandSide IObject, rightHandSide IObject, operationName string) *errors.Error {
	operation, getError := rightHandSide.Get("Right" + operationName)
	if getError != nil {
		return getError
	}
	if _, ok := operation.(*Function); !ok {
		return errors.NewTypeError(operation.TypeName(), FunctionName)
	}
	result, callError := CallFunction(operation.(*Function), p, rightHandSide.SymbolTable(), leftHandSide)
	if callError != nil {
		return callError
	}
	p.PushObject(result)
	return nil
}

func (p *Plasma) indexOP() *errors.Error {
	index := p.PopObject()
	source := p.PopObject()
	indexOperation, getError := source.Get(Index)
	if getError != nil {
		return getError
	}
	if _, ok := indexOperation.(*Function); !ok {
		return errors.NewTypeError(indexOperation.TypeName(), FunctionName)
	}
	result, callError := CallFunction(indexOperation.(*Function), p, source.SymbolTable(), index)
	if callError != nil {
		return callError
	}
	p.PushObject(result)
	return nil
}

func (p *Plasma) selectNameFromObjectOP(code Code) *errors.Error {
	name := code.Value.(string)
	source := p.PopObject()
	value, getError := source.Get(name)
	if getError != nil {
		return getError
	}
	p.PushObject(value)
	return nil
}

func (p *Plasma) methodInvocationOP(code Code) *errors.Error {
	numberOfArguments := code.Value.(int)
	function := p.PopObject()
	var arguments []IObject
	for i := 0; i < numberOfArguments; i++ {
		if !p.MemoryStack.HasNext() {
			return errors.NewInvalidNumberOfArguments(i, numberOfArguments)
		}
		arguments = append(arguments, p.PopObject())
	}
	var result IObject
	var callError *errors.Error
	switch function.(type) {
	case *Function:
		result, callError = CallFunction(function.(*Function), p, p.PeekSymbolTable(), arguments...)
	default:
		// ToDo: Add Support for Types too
		return errors.NewTypeError(function.TypeName(), FunctionName)
	}
	if callError != nil {
		return callError
	}
	p.PushObject(result)
	return nil
}

func (p *Plasma) getIdentifierOP(code Code) *errors.Error {
	value, getError := p.PeekSymbolTable().GetAny(code.Value.(string))
	if getError != nil {
		return getError
	}
	p.PushObject(value)
	return nil
}

// Assign Statement

func (p *Plasma) assignIdentifierOP(code Code) *errors.Error {
	identifier := code.Value.(string)
	p.PeekSymbolTable().Set(identifier, p.PopObject())
	return nil
}

func (p *Plasma) assignSelectorOP(code Code) *errors.Error {
	target := p.PopObject()
	value := p.PopObject()
	identifier := code.Value.(string)
	target.Set(identifier, value)
	return nil
}

func (p *Plasma) assignIndexOP() *errors.Error {
	index := p.PopObject()
	source := p.PopObject()
	value := p.PopObject()
	sourceAssign, getError := source.Get(Assign)
	if getError != nil {
		return getError
	}
	if _, ok := sourceAssign.(*Function); !ok {
		return errors.NewTypeError(sourceAssign.TypeName(), FunctionName)
	}
	_, callError := CallFunction(sourceAssign.(*Function), p, p.PeekSymbolTable(), index, value)
	if callError != nil {
		return callError
	}
	return nil
}

func (p *Plasma) returnOP(code Code) *errors.Error {
	numberOfReturnValues := code.Value.(int)
	if numberOfReturnValues == 0 {
		noneValue, getError := p.PeekSymbolTable().GetAny(None)
		if getError != nil {
			return getError
		}
		p.PushObject(noneValue)
		return nil
	}

	var values []IObject
	for i := 0; i < numberOfReturnValues; i++ {
		if !p.MemoryStack.HasNext() {
			return errors.NewInvalidNumberOfArguments(i, numberOfReturnValues)
		}
		values = append(values, p.PopObject())
	}
	if len(values) == 1 {
		p.PushObject(values[0])
	} else {
		p.PushObject(NewTuple(p.PeekSymbolTable(), values))
	}
	return nil
}

// Special Instructions

func (p *Plasma) loadFunctionArgumentsOP(code Code) *errors.Error {
	for _, argument := range code.Value.([]string) {
		value := p.PopObject()
		p.PeekSymbolTable().Set(argument, value)
	}
	return nil
}

func (p *Plasma) newFunctionOP(code Code) *errors.Error {
	functionInformation := code.Value.([2]int)
	codeLength := functionInformation[0]
	numberOfArguments := functionInformation[1]
	start := p.PeekCode().index
	p.PeekCode().index += codeLength
	end := p.PeekCode().index
	functionCode := make([]Code, codeLength)
	copy(functionCode, p.PeekCode().instructions[start:end])
	p.PushObject(NewFunction(p.PeekSymbolTable(), NewPlasmaFunction(numberOfArguments, functionCode)))
	return nil
}

func (p *Plasma) Execute() (IObject, *errors.Error) {
	var executionError *errors.Error
	for ; p.PeekCode().HasNext(); {
		code := p.PeekCode().Next()
		switch code.Instruction.OpCode {
		// Literals
		case NewStringOP:
			executionError = p.newStringOP(code)
		case NewBytesOP:
			executionError = p.newBytesOP(code)
		case NewIntegerOP:
			executionError = p.newIntegerOP(code)
		case NewFloatOP:
			executionError = p.newFloatOP(code)
		case NewTrueBoolOP:
			executionError = p.newTrueBoolOP()
		case NewFalseBoolOP:
			executionError = p.newFalseBoolOP()
		case GetNoneOP:
			executionError = p.getNoneOP()
		// Composite creation
		case NewTupleOP:
			executionError = p.newTupleOP(code)
		case NewArrayOP:
			executionError = p.newArrayOP(code)
		case NewHashOP:
			executionError = p.newHashOP(code)
		// Unary Expressions
		case NegateBitsOP:
			executionError = p.noArgsGetAndCall(NegBits)
		case BoolNegateOP:
			executionError = p.noArgsGetAndCall(Negate)
		case NegativeOP:
			executionError = p.noArgsGetAndCall(Negative)
		// Binary Expressions
		case AddOP:
			executionError = p.leftBinaryExpressionFuncCall(Add)
		case SubOP:
			executionError = p.leftBinaryExpressionFuncCall(Sub)
		case MulOP:
			executionError = p.leftBinaryExpressionFuncCall(Mul)
		case DivOP:
			executionError = p.leftBinaryExpressionFuncCall(Div)
		case ModOP:
			executionError = p.leftBinaryExpressionFuncCall(Mod)
		case PowOP:
			executionError = p.leftBinaryExpressionFuncCall(Pow)
		case BitXorOP:
			executionError = p.leftBinaryExpressionFuncCall(BitXor)
		case BitAndOP:
			executionError = p.leftBinaryExpressionFuncCall(BitAnd)
		case BitOrOP:
			executionError = p.leftBinaryExpressionFuncCall(BitOr)
		case BitLeftOP:
			executionError = p.leftBinaryExpressionFuncCall(BitLeft)
		case BitRightOP:
			executionError = p.leftBinaryExpressionFuncCall(BitRight)
		case AndOP:
			executionError = p.leftBinaryExpressionFuncCall(And)
		case OrOP:
			executionError = p.leftBinaryExpressionFuncCall(Or)
		case XorOP:
			executionError = p.leftBinaryExpressionFuncCall(Xor)
		case EqualsOP:
			executionError = p.leftBinaryExpressionFuncCall(Equals)
		case NotEqualsOP:
			executionError = p.leftBinaryExpressionFuncCall(NotEquals)
		case GreaterThanOP:
			executionError = p.leftBinaryExpressionFuncCall(GreaterThan)
		case LessThanOP:
			executionError = p.leftBinaryExpressionFuncCall(LessThan)
		case GreaterThanOrEqualOP:
			executionError = p.leftBinaryExpressionFuncCall(GreaterThanOrEqual)
		case LessThanOrEqualOP:
			executionError = p.leftBinaryExpressionFuncCall(LessThanOrEqual)
		// Other Expressions
		case GetIdentifierOP:
			executionError = p.getIdentifierOP(code)
		case IndexOP:
			executionError = p.indexOP()
		case SelectNameFromObjectOP:
			executionError = p.selectNameFromObjectOP(code)
		case MethodInvocationOP:
			executionError = p.methodInvocationOP(code)
		// Assign Statement
		case AssignIdentifierOP:
			executionError = p.assignIdentifierOP(code)
		case AssignSelectorOP:
			executionError = p.assignSelectorOP(code)
		case AssignIndexOP:
			executionError = p.assignIndexOP()
		case ReturnOP:
			executionError = p.returnOP(code)
		// Special Instructions
		case LoadFunctionArgumentsOP:
			executionError = p.loadFunctionArgumentsOP(code)
		case NewFunctionOP:
			executionError = p.newFunctionOP(code)
		default:
			return nil, errors.NewUnknownVMOperationError(code.Instruction.OpCode)
		}
		if executionError != nil {
			return nil, executionError
		}
	}
	if p.MemoryStack.HasNext() {
		return p.PopObject(), nil
	}
	return p.PeekSymbolTable().GetAny(None)
}

func (p *Plasma) HashString(s string) (int64, *errors.Error) {
	_, hashingError := p.Crc32Hash.Write([]byte(s))
	if hashingError != nil {
		return 0, errors.NewHashingStringError()
	}
	hashValue := p.Crc32Hash.Sum32()
	p.Crc32Hash.Reset()
	return int64(hashValue), nil
}

func (p *Plasma) HashBytes(s []byte) (int64, *errors.Error) {
	_, hashingError := p.Crc32Hash.Write(s)
	if hashingError != nil {
		return 0, errors.NewHashingStringError()
	}
	hashValue := p.Crc32Hash.Sum32()
	p.Crc32Hash.Reset()
	return int64(hashValue), nil
}

func (p *Plasma) Seed() uint64 {
	return p.seed
}

func (p *Plasma) InitializeByteCode(bytecode *Bytecode) {
	p.PushCode(bytecode)
	p.MemoryStack.Clear()
	p.Context.Clear()
	p.Context.Push(SetDefaultSymbolTable())
}

func (p *Plasma) PushCode(code *Bytecode) {
	p.Code.Push(code)
}

func (p *Plasma) PopCode() *Bytecode {
	return p.Code.Pop()
}

func (p *Plasma) PeekCode() *Bytecode {
	return p.Code.Peek()
}

func NewPlasmaVM() *Plasma {
	number, randError := rand.Int(rand.Reader, big.NewInt(polySize))
	if randError != nil {
		panic(randError)
	}
	return &Plasma{
		Code:        NewCodeStack(),
		MemoryStack: NewObjectStack(),
		Context:     NewSymbolStack(),
		Crc32Hash:   crc32.New(crc32.MakeTable(polySize)),
		seed:        number.Uint64(),
	}
}
