package vm

import (
	"github.com/shoriwe/gruby/pkg/errors"
	"hash"
	"hash/crc64"
)

type Plasma struct {
	Code        *Bytecode
	MemoryStack *ObjectStack
	Context     *SymbolStack
	Crc64Hash   hash.Hash64
}

func (p *Plasma) LoadCode(codes []Code) {
	for _, code := range codes {
		p.Code.Push(code)
	}
}

func (p *Plasma) PushSymbolTable(table *SymbolTable) {
	p.Context.Push(table)
}

func (p *Plasma) PopSymbolTable() *SymbolTable {
	return p.Context.Pop()
}

func (p *Plasma) Initialize(code []Code) *errors.Error {
	p.Code = NewBytecodeFromArray(code)
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

func (p *Plasma) newIntegerOP(code Code) *errors.Error {
	value := code.Value.(int64)
	integer := NewInteger(p.Context.Peek(), value)
	p.MemoryStack.Push(integer)
	return nil
}

func (p *Plasma) constructObject(type_ *Type) *errors.Error {
	object, constructionError := ConstructObject(type_, p, p.PeekSymbolTable())
	if constructionError != nil {
		return constructionError
	}
	initFunc, getError := object.Get(Initialize)
	if getError != nil {
		return getError
	}
	if _, ok := initFunc.(*Function); !ok {
		return errors.NewTypeError(initFunc.TypeName(), FunctionName)
	}
	numberOfArguments := p.Code.Next().Value.(int)
	var arguments []IObject
	for i := 0; i < numberOfArguments; i++ {
		arguments = append(arguments, p.MemoryStack.Pop())
	}
	_, callError := CallFunction(initFunc.(*Function), p, object.SymbolTable(), arguments...)
	if callError != nil {
		return callError
	}
	p.MemoryStack.Push(object)
	return nil
}

func (p *Plasma) callOP() *errors.Error {
	function := p.MemoryStack.Pop()
	if _, ok := function.(*Function); !ok {
		if _, ok = function.(*Type); ok {
			return p.constructObject(function.(*Type))
		} else {
			var getError *errors.Error
			function, getError = function.Get(Call)
			if getError != nil {
				return getError
			}
			if _, ok2 := function.(*Function); !ok2 {
				return errors.New(errors.UnknownLine, "Expecting Function", "NonFunctionObjectReceived")
			}
		}
	}
	var parent *SymbolTable
	if p.Code.Next().Value.(bool) {
		parent = function.SymbolTable().Parent
	} else {
		parent = p.Context.Peek()
	}
	numberOfArguments := p.Code.Next().Value.(int)
	var arguments []IObject
	for i := 0; i < numberOfArguments; i++ {
		arguments = append(arguments, p.MemoryStack.Pop())
	}
	var result IObject
	var callError *errors.Error
	if _, ok3 := function.(*Function).Callable.(Constructor); ok3 {
		result, callError = CallFunction(function.(*Function), p, parent, nil)
		if callError != nil {
			return callError
		}
		resultInit, getError := result.Get(Initialize)
		if getError != nil {
			return getError
		}
		if _, ok4 := resultInit.(*Function); !ok4 {
			return errors.New(errors.UnknownLine, "Expecting Function", "NonFunctionObjectReceived")
		}
		_, callError = CallFunction(resultInit.(*Function), p, result.SymbolTable(), arguments...)
		if callError != nil {
			return callError
		}
	} else {
		result, callError = CallFunction(function.(*Function), p, parent, arguments...)
		if callError != nil {
			return callError
		}
	}
	p.MemoryStack.Push(result)
	return nil
}

func (p *Plasma) getOP() *errors.Error {
	name := p.Code.Next().Value.(string)
	result, getError := p.Context.Peek().GetAny(name)
	if getError != nil {
		return getError
	}
	p.MemoryStack.Push(result)
	return nil
}

func (p *Plasma) getFromOP() *errors.Error {
	name := p.Code.Next().Value.(string)
	result, getError := p.Context.Pop().GetSelf(name)
	if getError != nil {
		return getError
	}
	p.MemoryStack.Push(result)
	return nil
}

func (p *Plasma) Execute() (IObject, *errors.Error) {
	var executionError *errors.Error
	defer p.PopSymbolTable()
	for ; p.Code.HasNext(); {
		code := p.Code.Next()
		switch code.Instruction.OpCode {
		case NewStringOP:
			executionError = p.newStringOP(code)
		case NewIntegerOP:
			executionError = p.newIntegerOP(code)
		case CallOP:
			executionError = p.callOP()
		case GetOP:
			executionError = p.getOP()
		case GetFromOP:
			executionError = p.getFromOP()
		case ReturnOP:
			if p.MemoryStack.HasNext() {
				return p.MemoryStack.Pop(), nil
			}
			return p.PeekSymbolTable().GetAny(None)
		default:
			return nil, errors.NewUnknownVMOperationError(code.Instruction.OpCode)
		}
		if executionError != nil {
			return nil, executionError
		}
	}

	return p.PeekSymbolTable().GetAny(None)
}

func (p *Plasma) HashString(s string) (int64, *errors.Error) {
	_, hashingError := p.Crc64Hash.Write([]byte(s))
	if hashingError != nil {
		return 0, errors.NewHashingStringError()
	}
	hashValue := p.Crc64Hash.Sum64()
	p.Crc64Hash.Reset()
	return int64(hashValue), nil
}

func (p *Plasma) InitializeByteCode(bytecode *Bytecode) {
	p.Code = bytecode
	p.MemoryStack.Clear()
	p.Context.Clear()
	p.Context.Push(SetDefaultSymbolTable())
}

func NewPlasmaVM() *Plasma {
	return &Plasma{
		Code:        nil,
		MemoryStack: NewObjectStack(),
		Context:     NewSymbolStack(),
		Crc64Hash:   crc64.New(crc64.MakeTable(9223372036854775807)),
	}
}
