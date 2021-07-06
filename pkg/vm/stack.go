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

type ObjectStack struct {
	head *stackNode
}

func (stack *ObjectStack) Pop() Value {
	result := stack.head.value
	stack.head = stack.head.next
	return result.(Value)
}

func (stack *ObjectStack) Peek() Value {
	return stack.head.value.(Value)
}

func (stack *ObjectStack) Push(object Value) {
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

type loopEntry struct {
	Action uint8
}

type LoopStack struct {
	head *stackNode
}

func (stack *LoopStack) Pop() *loopEntry {
	result := stack.head.value
	stack.head = stack.head.next
	return result.(*loopEntry)
}

func (stack *LoopStack) Peek() *loopEntry {
	return stack.head.value.(*loopEntry)
}

func (stack *LoopStack) Push(tryStackEntry *loopEntry) {
	stack.head = NewStackNode(tryStackEntry, stack.head)
}

func (stack *LoopStack) HasNext() bool {
	return stack.head != nil
}

func (stack *LoopStack) Clear() {
	stack.head = nil
}

func NewLoopStack() *LoopStack {
	return &LoopStack{
		head: nil,
	}
}
