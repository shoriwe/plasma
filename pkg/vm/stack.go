package vm

type ArrayStack struct {
	Contents []*Value
	length   int
}

func (as *ArrayStack) Push(value *Value) {
	if as.length >= len(as.Contents) {
		as.length++
		as.Contents = append(as.Contents, value)
		return
	}
	as.Contents[as.length] = value
	as.length++
}

func (as *ArrayStack) Pop() *Value {
	if as.length <= 0 {
		panic("pop from empty stack")
	}
	as.length--
	return as.Contents[as.length]
}

func (as *ArrayStack) Peek() *Value {
	if as.length <= 0 {
		panic("peek from empty stack")
	}
	return as.Contents[as.length-1]
}

func NewArrayStack(reserved int) *ArrayStack {
	return &ArrayStack{
		Contents: make([]*Value, reserved),
		length:   0,
	}
}
