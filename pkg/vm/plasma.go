package vm

import (
	"github.com/shoriwe/gruby/pkg/errors"
	"github.com/shoriwe/gruby/pkg/vm/object"
	"github.com/shoriwe/gruby/pkg/vm/utils"
)

type Plasma struct {
	Code              []interface{}
	Cursor            int
	CodeLength        int
	MemoryStack       utils.Stack
	masterSymbolTable *utils.SymbolTable
}

func (p *Plasma) GetStack() utils.Stack {
	return p.MemoryStack
}

func (p *Plasma) Initialize(code []interface{}) *errors.Error {
	p.Code = code
	p.CodeLength = len(code)
	p.MemoryStack.Clear()
	return nil
}

func (p *Plasma) newOp() *errors.Error {
	p.Cursor++
	type_ := p.MemoryStack.Pop().(*object.Type)
	numberOfArguments := p.MemoryStack.Pop().(int)
	var arguments []interface{}
	for i := 0; i < numberOfArguments; i++ {
		arguments = append(arguments, p.MemoryStack.Pop())
	}
	instance, creationError := type_.Constructor(p, p.masterSymbolTable, p.masterSymbolTable, arguments)
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

func (p *Plasma) pushOP() *errors.Error {
	p.Cursor++
	value := p.Code[p.Cursor]
	p.Cursor++
	return p.MemoryStack.Push(value)
}

func (p *Plasma) Execute() (interface{}, *errors.Error) {
	var executionError *errors.Error
	for ; p.Cursor < p.CodeLength; {
		switch p.Code[p.Cursor].(uint16) {
		case NewOP:
			executionError = p.newOp()
		case GetOP:
			executionError = p.getOP()
		case PushOP:
			executionError = p.pushOP()
		case ReturnOP:
			if p.MemoryStack.HashNext() {
				return p.MemoryStack.Pop(), nil
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

func (p *Plasma) New() utils.VirtualMachine {
	return &Plasma{
		Code:              nil,
		Cursor:            0,
		CodeLength:        0,
		MemoryStack:       utils.NewArrayStack(),
		masterSymbolTable: p.masterSymbolTable,
	}
}

func (p *Plasma) MasterSymbolTable() *utils.SymbolTable {
	return p.masterSymbolTable
}

func NewPlasmaVM() *Plasma {
	return &Plasma{
		Code:              nil,
		Cursor:            0,
		MemoryStack:       utils.NewArrayStack(),
		masterSymbolTable: utils.NewSymbolTable(nil, nil),
	}
}
