package vm

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

type IterEntry struct {
	Iterable  IObject
	LastValue IObject // Used when redo is called
}

type IterStack struct {
	head *stackNode
}

func (stack *IterStack) Pop() *IterEntry {
	result := stack.head.value
	stack.head = stack.head.next
	return result.(*IterEntry)
}

func (stack *IterStack) Peek() *IterEntry {
	return stack.head.value.(*IterEntry)
}

func (stack *IterStack) Push(iter *IterEntry) {
	stack.head = NewStackNode(iter, stack.head)
}

func (stack *IterStack) HasNext() bool {
	return stack.head != nil
}

func (stack *IterStack) Clear() {
	stack.head = nil
}

func NewIterStack() *IterStack {
	return &IterStack{
		head: nil,
	}
}

type CodeStack struct {
	head *stackNode
}

func (stack *CodeStack) Pop() *Bytecode {
	result := stack.head.value
	stack.head = stack.head.next
	return result.(*Bytecode)
}

func (stack *CodeStack) Peek() *Bytecode {
	return stack.head.value.(*Bytecode)
}

func (stack *CodeStack) Push(code *Bytecode) {
	stack.head = NewStackNode(code, stack.head)
}

func (stack *CodeStack) HasNext() bool {
	return stack.head != nil
}

func (stack *CodeStack) Clear() {
	stack.head = nil
}

func NewCodeStack() *CodeStack {
	return &CodeStack{
		head: nil,
	}
}

type ObjectStack struct {
	head *stackNode
}

func (stack *ObjectStack) Pop() IObject {
	result := stack.head.value
	stack.head = stack.head.next
	return result.(IObject)
}

func (stack *ObjectStack) Peek() IObject {
	return stack.head.value.(IObject)
}

func (stack *ObjectStack) Push(object IObject) {
	stack.head = NewStackNode(object, stack.head)
}

func (stack *ObjectStack) HasNext() bool {
	return stack.head != nil
}

func (stack *ObjectStack) Clear() {
	stack.head = nil
}

func NewObjectStack() *ObjectStack {
	return &ObjectStack{
		head: nil,
	}
}

type SymbolStack struct {
	head *stackNode
}

func (stack *SymbolStack) Pop() *SymbolTable {
	result := stack.head.value
	stack.head = stack.head.next
	return result.(*SymbolTable)
}

func (stack *SymbolStack) Peek() *SymbolTable {
	return stack.head.value.(*SymbolTable)
}

func (stack *SymbolStack) Push(symbolTable *SymbolTable) {
	stack.head = NewStackNode(symbolTable, stack.head)
}

func (stack *SymbolStack) HasNext() bool {
	return stack.head != nil
}

func (stack *SymbolStack) Clear() {
	stack.head = nil
}

func NewSymbolStack() *SymbolStack {
	return &SymbolStack{
		head: nil,
	}
}
