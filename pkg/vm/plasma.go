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
	p.MemoryStack.Push(stringObject)
	return nil
}

func (p *Plasma) newBytesOP(code Code) *errors.Error {
	value := code.Value.([]byte)
	bytesObject := NewBytes(p.Context.Peek(), value)
	p.MemoryStack.Push(bytesObject)
	return nil
}

func (p *Plasma) newIntegerOP(code Code) *errors.Error {
	value := code.Value.(int64)
	integer := NewInteger(p.Context.Peek(), value)
	p.MemoryStack.Push(integer)
	return nil
}

func (p *Plasma) newFloatOP(code Code) *errors.Error {
	value := code.Value.(float64)
	float := NewFloat(p.Context.Peek(), value)
	p.MemoryStack.Push(float)
	return nil
}

func (p *Plasma) newTrueBoolOP() *errors.Error {
	p.MemoryStack.Push(NewBool(p.PeekSymbolTable(), true))
	return nil
}

func (p *Plasma) newFalseBoolOP() *errors.Error {
	p.MemoryStack.Push(NewBool(p.PeekSymbolTable(), false))
	return nil
}

func (p *Plasma) getNoneOP() *errors.Error {
	none, getError := p.PeekSymbolTable().GetAny(None)
	if getError != nil {
		return getError
	}
	p.MemoryStack.Push(none)
	return nil
}

func (p *Plasma) newTupleOP(code Code) *errors.Error {
	numberOfValues := code.Value.(int)
	var values []IObject
	for i := 0; i < numberOfValues; i++ {
		if !p.MemoryStack.HasNext() {
			return errors.NewInvalidNumberOfArguments(i, numberOfValues)
		}
		values = append(values, p.MemoryStack.Pop())
	}
	p.MemoryStack.Push(NewTuple(p.PeekSymbolTable(), values))
	return nil
}

func (p *Plasma) newArrayOP(code Code) *errors.Error {
	numberOfValues := code.Value.(int)
	var values []IObject
	for i := 0; i < numberOfValues; i++ {
		if !p.MemoryStack.HasNext() {
			return errors.NewInvalidNumberOfArguments(i, numberOfValues)
		}
		values = append(values, p.MemoryStack.Pop())
	}
	p.MemoryStack.Push(NewArray(p.PeekSymbolTable(), values))
	return nil
}

func (p *Plasma) newHashOP(code Code) *errors.Error {
	numberOfValues := code.Value.(int)
	var keyValues []*KeyValue
	for i := 0; i < numberOfValues; i++ {
		if !p.MemoryStack.HasNext() {
			return errors.NewInvalidNumberOfArguments(i, numberOfValues)
		}
		key := p.MemoryStack.Pop()
		if !p.MemoryStack.HasNext() {
			return errors.NewInvalidNumberOfArguments(i, numberOfValues)
		}
		value := p.MemoryStack.Pop()
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
	p.MemoryStack.Push(hashTable)
	return nil
}

// Useful function to call those built ins that doesn't receive arguments of an object
func (p *Plasma) noArgsGetAndCall(operationName string) *errors.Error {
	object := p.MemoryStack.Pop()
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
	p.MemoryStack.Push(result)
	return nil
}

// Function useful to cal object built-ins binary expression functions
func (p *Plasma) leftBinaryExpressionFuncCall(operationName string) *errors.Error {
	leftHandSide := p.MemoryStack.Pop()
	rightHandSide := p.MemoryStack.Pop()
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
	p.MemoryStack.Push(result)
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
	p.MemoryStack.Push(result)
	return nil
}

func (p *Plasma) indexOP() *errors.Error {
	index := p.MemoryStack.Pop()
	source := p.MemoryStack.Pop()
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
	p.MemoryStack.Push(result)
	return nil
}

func (p *Plasma) returnOP(code Code) *errors.Error {
	numberOfReturnValues := code.Value.(int)
	if numberOfReturnValues == 0 {
		noneValue, getError := p.PeekSymbolTable().GetAny(None)
		if getError != nil {
			return getError
		}
		p.MemoryStack.Push(noneValue)
		return nil
	}

	var values []IObject
	for i := 0; i < numberOfReturnValues; i++ {
		if !p.MemoryStack.HasNext() {
			return errors.NewInvalidNumberOfArguments(i, numberOfReturnValues)
		}
		values = append(values, p.MemoryStack.Pop())
	}
	p.MemoryStack.Push(NewTuple(p.PeekSymbolTable(), values))
	return nil
}

func (p *Plasma) Execute() (IObject, *errors.Error) {
	var executionError *errors.Error
	defer func() {
		p.PopSymbolTable()
		p.PopCode()
	}()
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
		//
		case IndexOP:
			executionError = p.indexOP()
		//
		case ReturnOP:
			executionError = p.returnOP(code)
		default:
			return nil, errors.NewUnknownVMOperationError(code.Instruction.OpCode)
		}
		if executionError != nil {
			return nil, executionError
		}
	}
	if p.MemoryStack.HasNext() {
		return p.MemoryStack.Pop(), nil
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
