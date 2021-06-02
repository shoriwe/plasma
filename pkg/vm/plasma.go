package vm

import (
	"crypto/rand"
	"github.com/shoriwe/gruby/pkg/errors"
	"hash"
	"hash/crc32"
	"io"
	"math/big"
)

const (
	polySize = 0xffffffff
)

type Plasma struct {
	IterStack                *IterStack
	programMasterSymbolTable *SymbolTable
	Code                     *CodeStack
	MemoryStack              *ObjectStack
	Context                  *SymbolStack
	Crc32Hash                hash.Hash32
	seed                     uint64
	stdin                    io.Reader
	stdout                   io.Writer
	stderr                   io.Writer
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

func (p *Plasma) Initialize(code []Code) {
	p.InitializeByteCode(NewBytecodeFromArray(code))
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
		result, callError = CallFunction(function.(*Function), p, function.SymbolTable(), arguments...)
	case *Type:
		result, callError = ConstructObject(function.(*Type), p, NewSymbolTable(p.PeekSymbolTable()))
	default:
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

func (p *Plasma) ifJumpOP(code Code) *errors.Error {
	condition := p.PopObject()
	toBool, getError := condition.Get(ToBool)
	if getError != nil {
		return getError
	}
	if _, ok := toBool.(*Function); !ok {
		return errors.NewTypeError(toBool.TypeName(), FunctionName)
	}
	conditionBool, callError := CallFunction(toBool.(*Function), p, toBool.SymbolTable())
	if callError != nil {
		return callError
	}
	if !conditionBool.GetBool() {
		p.PeekCode().index += code.Value.(int)
	}
	return nil
}

func (p *Plasma) unlessJumpOP(code Code) *errors.Error {
	condition := p.PopObject()
	toBool, getError := condition.Get(ToBool)
	if getError != nil {
		return getError
	}
	if _, ok := toBool.(*Function); !ok {
		return errors.NewTypeError(toBool.TypeName(), FunctionName)
	}
	conditionBool, callError := CallFunction(toBool.(*Function), p, toBool.SymbolTable())
	if callError != nil {
		return callError
	}
	if conditionBool.GetBool() {
		p.PeekCode().index += code.Value.(int)
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

func (p *Plasma) jumpOP(code Code) *errors.Error {
	p.PeekCode().index += code.Value.(int)
	return nil
}

func (p *Plasma) breakOP(code Code) *errors.Error {
	p.PeekCode().index += code.Value.(int)
	return nil
}

func (p *Plasma) redoOP(code Code) *errors.Error {
	p.PeekCode().index += code.Value.(int)
	return nil
}

func (p *Plasma) continueOP(code Code) *errors.Error {
	p.PeekCode().index += code.Value.(int) // This could be positive or negative (positive only for do-while loop)
	return nil
}

func (p *Plasma) setupForLoopOP() *errors.Error {
	// First check if it is a iterator
	value := p.PopObject()
	if _, ok := value.(*Iterator); ok {
		p.IterStack.Push(&IterEntry{
			Iterable:  value,
			LastValue: nil,
		})
		return nil
	}
	// Then check if the type implements iterator
	_, getError := value.Get(HasNext)
	if getError == nil {
		_, getError = value.Get(Next)
		if getError == nil {
			p.IterStack.Push(&IterEntry{
				Iterable:  value,
				LastValue: nil,
			})
			return nil
		}
	}
	// Finally transform it to iterable
	valueIterFunc, getError := value.Get(Iter)
	if getError != nil {
		return getError
	}
	if _, ok := valueIterFunc.(*Function); !ok {
		return errors.NewTypeError(valueIterFunc.TypeName(), FunctionName)
	}
	valueIter, callError := CallFunction(valueIterFunc.(*Function), p, value.SymbolTable())
	if callError != nil {
		return callError
	}
	p.IterStack.Push(&IterEntry{
		Iterable:  valueIter,
		LastValue: nil,
	})
	return nil
}

func (p *Plasma) hasNextOP(code Code) *errors.Error {
	hasNext, getError := p.IterStack.Peek().Iterable.Get("HasNext")
	if getError != nil {
		return getError
	}
	if _, ok := hasNext.(*Function); !ok {
		return errors.NewTypeError(hasNext.TypeName(), FunctionName)
	}
	result, callError := CallFunction(hasNext.(*Function), p, p.IterStack.Peek().Iterable.SymbolTable())
	if callError != nil {
		return callError
	}
	if _, ok := result.(*Bool); !ok {
		var resultToBool IObject
		resultToBool, getError = hasNext.Get(ToBool)
		if getError != nil {
			return getError
		}
		if _, ok = resultToBool.(*Function); !ok {
			return errors.NewTypeError(resultToBool.TypeName(), FunctionName)
		}
		var resultBool IObject
		resultBool, callError = CallFunction(resultToBool.(*Function), p, hasNext.SymbolTable())
		if callError != nil {
			return callError
		}
		if _, ok = resultBool.(*Bool); !ok {
			return errors.NewTypeError(resultBool.TypeName(), BoolName)
		}
		result = resultBool
	}
	if !result.GetBool() {
		p.PeekCode().index += code.Value.(int)
	}
	return nil
}

func (p *Plasma) unpackReceiversPopOP() *errors.Error {
	next, getError := p.IterStack.Peek().Iterable.Get(Next)
	if getError != nil {
		return getError
	}
	if _, ok := next.(*Function); !ok {
		return errors.NewTypeError(next.TypeName(), FunctionName)
	}
	nextValue, callError := CallFunction(next.(*Function), p, p.IterStack.Peek().Iterable.SymbolTable())
	if callError != nil {
		return callError
	}
	p.IterStack.Peek().LastValue = nextValue
	return nil
}

func (p *Plasma) unpackReceiversPeekOP(code Code) *errors.Error {
	receivers := code.Value.([]string)
	if len(receivers) == 1 {
		p.PeekSymbolTable().Set(receivers[0], p.IterStack.Peek().LastValue)
		return nil
	}
	// First try to unpack iterators
	hasNext, getError := p.IterStack.Peek().LastValue.Get(HasNext)
	if _, ok := hasNext.(*Function); getError == nil && ok {
		var next IObject
		next, getError = p.IterStack.Peek().LastValue.Get(Next)
		if _, ok = next.(*Function); getError == nil && ok {
			for _, receiver := range receivers {
				// First check if there is next value
				hasNextResult, callError := CallFunction(hasNext.(*Function), p, p.IterStack.Peek().LastValue.SymbolTable())
				if callError != nil {
					return callError
				}
				var hasNextResultBool IObject
				if _, ok = hasNextResult.(*Bool); !ok {
					var hasNextResultToBool IObject
					hasNextResultToBool, getError = hasNextResult.Get(ToBool)
					if getError != nil {
						return getError
					}
					if _, ok = hasNextResultToBool.(*Function); !ok {
						return errors.NewTypeError(hasNextResultToBool.TypeName(), FunctionName)
					}
					hasNextResultBool, callError = CallFunction(hasNextResultToBool.(*Function), p, hasNextResult.SymbolTable())
					if callError != nil {
						return callError
					}
					if _, ok = hasNextResultBool.(*Bool); !ok {
						return errors.NewTypeError(hasNextResultBool.TypeName(), BoolName)
					}
					hasNextResult = hasNextResultBool
				}
				if hasNextResult.GetBool() {
					var value IObject
					value, callError = CallFunction(next.(*Function), p, p.IterStack.Peek().LastValue.SymbolTable())
					if callError != nil {
						return callError
					}
					p.PeekSymbolTable().Set(receiver, value)
				}
			}
			return nil
		}
	}
	// Then try to unpack index-ables
	var lastValue IObject
	if _, ok := p.IterStack.Peek().LastValue.(*Tuple); !ok {
		var toTuple IObject
		toTuple, getError = p.IterStack.Peek().LastValue.Get(ToTuple)
		if getError != nil {
			return getError
		}
		if _, ok = toTuple.(*Function); !ok {
			return errors.NewTypeError(toTuple.TypeName(), FunctionName)
		}
		var callError *errors.Error
		lastValue, callError = CallFunction(toTuple.(*Function), p, p.IterStack.Peek().LastValue.SymbolTable())
		if callError != nil {
			return callError
		}
		if _, ok = lastValue.(*Tuple); !ok {
			return errors.NewTypeError(lastValue.TypeName(), TupleName)
		}
	} else {
		lastValue = p.IterStack.Peek().LastValue
	}
	if len(receivers) != lastValue.GetLength() {
		return errors.NewInvalidNumberOfArguments(lastValue.GetLength(), len(receivers))
	}
	for i, receiver := range receivers {
		p.PeekSymbolTable().Set(receiver, lastValue.GetContent()[i])
	}
	return nil
}

func (p *Plasma) newIteratorOP(code Code) *errors.Error {
	source := p.PopObject()
	var iterSource IObject
	var callError *errors.Error
	if _, ok := source.(*Iterator); ok {
		iterSource = source
	} else {
		iter, getError := source.Get(Iter)
		if getError != nil {
			return getError
		}
		if _, ok := iter.(*Function); !ok {
			return errors.NewTypeError(iter.TypeName(), FunctionName)
		}
		iterSource, callError = CallFunction(iter.(*Function), p, source.SymbolTable())
		if callError != nil {
			return callError
		}
	}
	generatorIterator := NewIterator(p.PeekSymbolTable())
	generatorIterator.Set(Source, iterSource)

	hasNextCodeLength, nextCodeLength := code.Value.([2]int)[0], code.Value.([2]int)[1]
	var hasNextCode []Code
	for i := 0; i < hasNextCodeLength; i++ {
		hasNextCode = append(hasNextCode, p.PeekCode().Next())
	}
	var nextCode []Code
	for i := 0; i < nextCodeLength; i++ {
		nextCode = append(nextCode, p.PeekCode().Next())
	}
	generatorIterator.Set(Next,
		NewFunction(generatorIterator.symbols,
			NewPlasmaClassFunction(generatorIterator, 0, nextCode),
		),
	)
	generatorIterator.Set(HasNext,
		NewFunction(generatorIterator.symbols,
			NewPlasmaClassFunction(generatorIterator, 0, hasNextCode),
		),
	)
	p.PushObject(generatorIterator)
	return nil
}

func (p *Plasma) Execute() (IObject, *errors.Error) {
	var executionError *errors.Error
	for ; p.PeekCode().HasNext(); {
		code := p.PeekCode().Next()
		// fmt.Println("Exec:", code.Instruction, code.Value)
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
			if executionError != nil {
				return nil, executionError
			}
			if p.MemoryStack.HasNext() {
				return p.PopObject(), nil
			}
			return p.PeekSymbolTable().GetAny(None)
		case IfJumpOP:
			executionError = p.ifJumpOP(code)
		case UnlessJumpOP:
			executionError = p.unlessJumpOP(code)
		// Special Instructions
		case LoadFunctionArgumentsOP:
			executionError = p.loadFunctionArgumentsOP(code)
		case NewFunctionOP:
			executionError = p.newFunctionOP(code)
		case JumpOP:
			executionError = p.jumpOP(code)
		case RedoOP:
			executionError = p.redoOP(code)
		case BreakOP:
			executionError = p.breakOP(code)
		case ContinueOP:
			executionError = p.continueOP(code)
		case PopOP:
			p.MemoryStack.Pop()
		case NOP:
			break
		case SetupForLoopOP:
			executionError = p.setupForLoopOP()
		case HasNextOP:
			executionError = p.hasNextOP(code)
		case UnpackReceiversPopOP:
			executionError = p.unpackReceiversPopOP()
		case UnpackReceiversPeekOP:
			executionError = p.unpackReceiversPeekOP(code)
		case PopIterOP:
			p.IterStack.Pop()
		case NewIteratorOP:
			executionError = p.newIteratorOP(code)
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
	p.programMasterSymbolTable = SetDefaultSymbolTable()
	p.Context.Push(p.programMasterSymbolTable)
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

func (p *Plasma) StdIn() io.Reader {
	return p.stdin
}

func (p *Plasma) StdOut() io.Writer {
	return p.stdout
}

func (p *Plasma) StdErr() io.Writer {
	return p.stderr
}

func (p *Plasma) MasterSymbolTable() *SymbolTable {
	return p.programMasterSymbolTable
}

func NewPlasmaVM(stdin io.Reader, stdout io.Writer, stderr io.Writer) *Plasma {
	number, randError := rand.Int(rand.Reader, big.NewInt(polySize))
	if randError != nil {
		panic(randError)
	}
	return &Plasma{
		IterStack:   NewIterStack(),
		Code:        NewCodeStack(),
		MemoryStack: NewObjectStack(),
		Context:     NewSymbolStack(),
		Crc32Hash:   crc32.New(crc32.MakeTable(polySize)),
		seed:        number.Uint64(),
		stdin:       stdin,
		stdout:      stdout,
		stderr:      stderr,
	}
}
