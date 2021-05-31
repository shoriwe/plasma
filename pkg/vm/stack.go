package vm

import (
	"math/bits"
)

type stackNode struct {
	value interface{}
	next  *stackNode
}

func NewStackNode(value interface{}, next *stackNode) *stackNode {
	return &stackNode{
		value: value,
		next:  next,
	}
}

type CodeStack struct {
	head   *stackNode
	length uint
}

func (stack *CodeStack) Pop() *Bytecode {
	result := stack.head.value
	stack.head = stack.head.next
	stack.length--
	return result.(*Bytecode)
}

func (stack *CodeStack) Peek() *Bytecode {
	return stack.head.value.(*Bytecode)
}

func (stack *CodeStack) Push(code *Bytecode) {
	if stack.length == bits.UintSize {
		panic("Memory Stack is Full")
	}
	stack.length++
	stack.head = NewStackNode(code, stack.head)
}

func (stack *CodeStack) HasNext() bool {
	return stack.length > 0
}

func (stack *CodeStack) Clear() {
	stack.head = nil
	stack.length = 0
}

func NewCodeStack() *CodeStack {
	return &CodeStack{
		head:   nil,
		length: 0,
	}
}

type ObjectStack struct {
	head   *stackNode
	length uint
}

func (stack *ObjectStack) Pop() IObject {
	result := stack.head.value
	stack.head = stack.head.next
	stack.length--
	return result.(IObject)
}

func (stack *ObjectStack) Peek() IObject {
	return stack.head.value.(IObject)
}

func (stack *ObjectStack) Push(object IObject) {
	stack.length++
	stack.head = NewStackNode(object, stack.head)
}

func (stack *ObjectStack) HasNext() bool {
	return stack.length > 0
}

func (stack *ObjectStack) Clear() {
	stack.head = nil
	stack.length = 0
}

func NewObjectStack() *ObjectStack {
	return &ObjectStack{
		head:   nil,
		length: 0,
	}
}

type SymbolStack struct {
	head   *stackNode
	length uint
}

func (stack *SymbolStack) Pop() *SymbolTable {
	result := stack.head.value
	stack.head = stack.head.next
	stack.length--
	return result.(*SymbolTable)
}

func (stack *SymbolStack) Peek() *SymbolTable {
	return stack.head.value.(*SymbolTable)
}

func (stack *SymbolStack) Push(symbolTable *SymbolTable) {
	if stack.length == bits.UintSize {
		panic("Memory Stack is Full")
	}
	stack.length++
	stack.head = NewStackNode(symbolTable, stack.head)
}

func (stack *SymbolStack) HasNext() bool {
	return stack.length > 0
}

func (stack *SymbolStack) Clear() {
	stack.head = nil
	stack.length = 0
}

func NewSymbolStack() *SymbolStack {
	return &SymbolStack{
		head:   nil,
		length: 0,
	}
}
