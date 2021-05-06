package vm

import (
	"github.com/shoriwe/gruby/pkg/errors"
)

type Plasma struct {
	Code              []interface{}
	Cursor            int
	CodeLength        int
	MemoryStack       Stack
	masterSymbolTable *SymbolTable
}

func (p *Plasma) GetStack() Stack {
	return p.MemoryStack
}

func (p *Plasma) Initialize(code []interface{}) *errors.Error {
	p.Code = code
	p.CodeLength = len(code)
	p.MemoryStack.Clear()
	return nil
}

func (p *Plasma) newStringOP() *errors.Error {
	p.Cursor++
	value := p.MemoryStack.Pop().(string)
	stringObject := NewString(Empty, nil, p.masterSymbolTable, value)
	p.MemoryStack.Push(stringObject)
	return nil
}

func (p *Plasma) newOp() *errors.Error {
	p.Cursor++
	type_ := p.MemoryStack.Pop().(*Function)
	numberOfArguments := p.MemoryStack.Pop().(int)
	var arguments []IObject
	for i := 0; i < numberOfArguments; i++ {
		arguments = append(arguments, p.MemoryStack.Pop().(IObject))
	}
	// Create the object
	instance, creationError := type_.Callable.Call(p.masterSymbolTable, p, arguments...)
	if creationError != nil {
		return creationError
	}
	p.MemoryStack.Push(instance)
	return nil
}

func (p *Plasma) getOP() *errors.Error {
	p.Cursor++
	name := p.MemoryStack.Pop().(string)
	obj, getError := p.masterSymbolTable.GetSelf(name)
	if getError != nil {
		return getError
	}
	p.MemoryStack.Push(obj)
	return nil
}

func (p *Plasma) getFromOP() *errors.Error {
	p.Cursor++
	name := p.MemoryStack.Pop().(string)
	obj := p.MemoryStack.Pop().(IObject)
	target, getError := obj.Get(name)
	if getError != nil {
		return getError
	}
	p.MemoryStack.Push(target)
	return nil
}

func (p *Plasma) pushOP() *errors.Error {
	p.Cursor++
	value := p.Code[p.Cursor]
	p.Cursor++
	return p.MemoryStack.Push(value)
}

func (p *Plasma) Execute() (IObject, *errors.Error) {
	var executionError *errors.Error
	for ; p.Cursor < p.CodeLength; {
		switch p.Code[p.Cursor].(uint16) {
		case NewOP:
			executionError = p.newOp()
		case NewStringOP:
			executionError = p.newStringOP()
		case GetOP:
			executionError = p.getOP()
		case GetFromOP:
			executionError = p.getFromOP()
		case PushOP:
			executionError = p.pushOP()
		case ReturnOP:
			if p.MemoryStack.HashNext() {
				return p.MemoryStack.Pop().(IObject), nil
			}
			return nil, nil
		default:
			return nil, errors.NewUnknownVMOperationError(p.Code[p.Cursor].(uint16))
		}
		if executionError != nil {
			return nil, executionError
		}
	}
	return nil, nil
}

func (p *Plasma) New(symbolTable *SymbolTable) VirtualMachine {
	return &Plasma{
		Code:              nil,
		Cursor:            0,
		CodeLength:        0,
		MemoryStack:       NewArrayStack(),
		masterSymbolTable: symbolTable,
	}
}

func (p *Plasma) MasterSymbolTable() *SymbolTable {
	return p.masterSymbolTable
}

func NewPlasmaVM() *Plasma {
	return &Plasma{
		Code:              nil,
		Cursor:            0,
		MemoryStack:       NewArrayStack(),
		masterSymbolTable: NewSymbolTable(nil),
	}
}
