package vm

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/gplasma/pkg/common"
	"sync"
)

type (
	Context struct {
		Bytecode  []byte
		Index     int
		VM        *VM
		Namespace *Symbols
		TempMem   *Value
		Stack     *ArrayStack
	}
	VM struct {
		mutex                            *sync.Mutex
		initialized                      bool
		RootNamespace                    *Symbols
		TrueValue, FalseValue, NoneValue *Value
	}
)

func (vm *VM) init() {
	/*
		- Value class TODO
		- Bool class TODO
		- Integer class TODO
		- Float class TODO
		- String class TODO
		- Bytes class TODO
		- Array class TODO
		- Tuple class TODO
		- Hash class TODO
		- True
		- False
		- None
		- All global magic methods TODO
		- All global special symbols TODO
	*/
	panic("implement me!")
}

func (vm *VM) Push(ctx *Context) {
	ctx.Index++
	ctx.Stack.Push(ctx.TempMem)
}

func (vm *VM) Pop(ctx *Context) {
	ctx.Index++
	ctx.Stack.Pop()
}

func (vm *VM) IdentifierAssign(ctx *Context) {
	ctx.Index++
	symbolLength := int(common.BytesToInt(ctx.Bytecode[ctx.Index : ctx.Index+8]))
	ctx.Index += 8
	symbol := string(ctx.Bytecode[ctx.Index : ctx.Index+symbolLength])
	ctx.Index += symbolLength
	ctx.Namespace.Set(symbol, ctx.Stack.Pop())
}

func (vm *VM) SelectorAssign(ctx *Context) {
	target := ctx.Stack.Pop()
	ctx.Index++
	symbolLength := int(common.BytesToInt(ctx.Bytecode[ctx.Index : ctx.Index+8]))
	ctx.Index += 8
	symbol := string(ctx.Bytecode[ctx.Index : ctx.Index+symbolLength])
	ctx.Index += symbolLength
	target.VirtualTable.Set(symbol, ctx.Stack.Pop())
}

func (vm *VM) Label(ctx *Context) {
	ctx.Index += 9
}

func (vm *VM) Jump(ctx *Context) {
	ctx.Index += int(common.BytesToInt(ctx.Bytecode[ctx.Index+1 : ctx.Index+9]))
}

func (vm *VM) IfJump(ctx *Context) {
	condition := ctx.Stack.Pop()
	jump := int(common.BytesToInt(ctx.Bytecode[ctx.Index+1 : ctx.Index+9]))
	if condition.Bool() {
		ctx.Index += jump
	}
}

func (vm *VM) Return(ctx *Context) { panic("implement me!") }

func (vm *VM) Require(ctx *Context) { panic("implement me!") }

func (vm *VM) DeleteIdentifier(ctx *Context) { panic("implement me!") }

func (vm *VM) DeleteSelector(ctx *Context) { panic("implement me!") }

func (vm *VM) Defer(ctx *Context) { panic("implement me!") }

func (vm *VM) EnterBlock(ctx *Context) {
	ctx.Index++
	ctx.Namespace = NewSymbols(ctx.Namespace)
}

func (vm *VM) ExitBlock(ctx *Context) {
	ctx.Index++
	ctx.Namespace = ctx.Namespace.Parent
}

func (vm *VM) NewFunction(ctx *Context) { panic("implement me!") }

func (vm *VM) NewClass(ctx *Context) { panic("implement me!") }

func (vm *VM) Call(ctx *Context) { panic("implement me!") }

func (vm *VM) IfOneLiner(ctx *Context) { panic("implement me!") }

func (vm *VM) NewArray(ctx *Context) {
	ctx.Index++
	numberOfValues := common.BytesToInt(ctx.Bytecode[ctx.Index : ctx.Index+8])
	ctx.Index += 8
	values := make([]*Value, numberOfValues)
	for i := numberOfValues - 1; i >= 0; i-- {
		values[i] = ctx.Stack.Pop()
	}
	ctx.TempMem = ctx.ArrayValue(values)
}

func (vm *VM) NewTuple(ctx *Context) {
	ctx.Index++
	numberOfValues := common.BytesToInt(ctx.Bytecode[ctx.Index : ctx.Index+8])
	ctx.Index += 8
	values := make([]*Value, numberOfValues)
	for i := numberOfValues - 1; i >= 0; i-- {
		values[i] = ctx.Stack.Pop()
	}
	ctx.TempMem = ctx.TupleValue(values)
}

func (vm *VM) NewHash(ctx *Context) { panic("implement me!") }

func (vm *VM) Identifier(ctx *Context) {
	ctx.Index++
	symbolLength := int(common.BytesToInt(ctx.Bytecode[ctx.Index : ctx.Index+8]))
	ctx.Index += 8
	symbol := string(ctx.Bytecode[ctx.Index : ctx.Index+symbolLength])
	ctx.Index += symbolLength
	var getError error
	ctx.TempMem, getError = ctx.Namespace.Get(symbol)
	if getError != nil {
		panic(getError)
	}
}

func (vm *VM) Integer(ctx *Context) {
	ctx.Index++
	value := common.BytesToInt(ctx.Bytecode[ctx.Index : ctx.Index+8])
	ctx.Index += 8
	ctx.TempMem = ctx.IntegerValue(value)
}

func (vm *VM) Float(ctx *Context) {
	ctx.Index++
	value := common.BytesToFloat(ctx.Bytecode[ctx.Index : ctx.Index+8])
	ctx.Index += 8
	ctx.TempMem = ctx.FloatValue(value)
}

func (vm *VM) String(ctx *Context) {
	ctx.Index++
	contentsLength := int(common.BytesToInt(ctx.Bytecode[ctx.Index : ctx.Index+8]))
	ctx.Index += 8
	contents := ctx.Bytecode[ctx.Index : ctx.Index+contentsLength]
	ctx.Index += contentsLength
	ctx.TempMem = ctx.StringValue(contents)
}

func (vm *VM) Bytes(ctx *Context) {
	ctx.Index++
	contentsLength := int(common.BytesToInt(ctx.Bytecode[ctx.Index : ctx.Index+8]))
	ctx.Index += 8
	contents := ctx.Bytecode[ctx.Index : ctx.Index+contentsLength]
	ctx.Index += contentsLength
	ctx.TempMem = ctx.BytesValue(contents)
}

func (vm *VM) True(ctx *Context) {
	ctx.Index++
	ctx.TempMem = ctx.TrueValue()
}

func (vm *VM) False(ctx *Context) {
	ctx.Index++
	ctx.TempMem = ctx.FalseValue()
}

func (vm *VM) None(ctx *Context) {
	ctx.Index++
	ctx.TempMem = ctx.NoneValue()
}

func (vm *VM) Selector(ctx *Context) {
	x := ctx.Stack.Pop()
	ctx.Index++
	symbolLength := int(common.BytesToInt(ctx.Bytecode[ctx.Index : ctx.Index+8]))
	ctx.Index += 8
	symbol := string(ctx.Bytecode[ctx.Index : ctx.Index+symbolLength])
	ctx.Index += symbolLength
	var getError error
	ctx.TempMem, getError = x.Get(symbol)
	if getError != nil {
		panic(getError)
	}
}

func (vm *VM) Super(ctx *Context) {
	// TODO: Implement me!
	panic("implement me!")
}

func (vm *VM) ExecuteContext(ctx *Context) {
	if !vm.initialized {
		vm.init()
	}
	for ctx.Index < len(ctx.Bytecode) {
		switch ctx.Bytecode[ctx.Index] {
		case opcodes.Push:
			vm.Push(ctx)
		case opcodes.Pop:
			vm.Pop(ctx)
		case opcodes.IdentifierAssign:
			vm.IdentifierAssign(ctx)
		case opcodes.SelectorAssign:
			vm.SelectorAssign(ctx)
		case opcodes.Label:
			vm.Label(ctx)
		case opcodes.Jump:
			vm.Jump(ctx)
		case opcodes.IfJump:
			vm.IfJump(ctx)
		case opcodes.Return:
			vm.Return(ctx)
		case opcodes.Require:
			vm.Require(ctx)
		case opcodes.DeleteIdentifier:
			vm.DeleteIdentifier(ctx)
		case opcodes.DeleteSelector:
			vm.DeleteSelector(ctx)
		case opcodes.Defer:
			vm.Defer(ctx)
		case opcodes.EnterBlock:
			vm.EnterBlock(ctx)
		case opcodes.ExitBlock:
			vm.ExitBlock(ctx)
		case opcodes.NewFunction:
			vm.NewFunction(ctx)
		case opcodes.NewClass:
			vm.NewClass(ctx)
		case opcodes.Call:
			vm.Call(ctx)
		case opcodes.IfOneLiner:
			vm.IfOneLiner(ctx)
		case opcodes.NewArray:
			vm.NewArray(ctx)
		case opcodes.NewTuple:
			vm.NewTuple(ctx)
		case opcodes.NewHash:
			vm.NewHash(ctx)
		case opcodes.Identifier:
			vm.Identifier(ctx)
		case opcodes.Integer:
			vm.Integer(ctx)
		case opcodes.Float:
			vm.Float(ctx)
		case opcodes.String:
			vm.String(ctx)
		case opcodes.Bytes:
			vm.Bytes(ctx)
		case opcodes.True:
			vm.True(ctx)
		case opcodes.False:
			vm.False(ctx)
		case opcodes.None:
			vm.None(ctx)
		case opcodes.Selector:
			vm.Selector(ctx)
		case opcodes.Super:
			vm.Super(ctx)
		default:
			panic(fmt.Sprintf("unknown opcode %d", ctx.Bytecode[ctx.Index]))
		}
	}
}
