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

type stateEntry struct {
	Action uint8
}

type StateStack struct {
	head *stackNode
}

func (stack *StateStack) Pop() *stateEntry {
	result := stack.head.value
	stack.head = stack.head.next
	return result.(*stateEntry)
}

func (stack *StateStack) Peek() *stateEntry {
	return stack.head.value.(*stateEntry)
}

func (stack *StateStack) Push(tryStackEntry *stateEntry) {
	stack.head = NewStackNode(tryStackEntry, stack.head)
}

func (stack *StateStack) HasNext() bool {
	return stack.head != nil
}

func (stack *StateStack) Clear() {
	stack.head = nil
}

func NewStateStack() *StateStack {
	return &StateStack{
		head: nil,
	}
}

type propagationEntry struct {
	PropagationLevel int
}

func (p *propagationEntry) Decrement() {
	p.PropagationLevel--
}

type PropagationStack struct {
	head *stackNode
}

func (stack *PropagationStack) Pop() *propagationEntry {
	result := stack.head.value
	stack.head = stack.head.next
	return result.(*propagationEntry)
}

func (stack *PropagationStack) Peek() *propagationEntry {
	return stack.head.value.(*propagationEntry)
}

func (stack *PropagationStack) Push(initialPropagation int) {
	stack.head = NewStackNode(&propagationEntry{
		PropagationLevel: initialPropagation,
	}, stack.head)
}

func (stack *PropagationStack) HasNext() bool {
	return stack.head != nil
}

func (stack *PropagationStack) Clear() {
	stack.head = nil
}

func NewPropagationStack() *PropagationStack {
	return &PropagationStack{
		head: nil,
	}
}
