package runtime

type Stack struct {
	content       []interface{}
	contentLength uint
}

func (stack *Stack) Pop() interface{} {
	stack.contentLength--
	result := stack.content[stack.contentLength]
	stack.content = stack.content[:stack.contentLength]
	return result
}

func (stack *Stack) Peek() interface{} {
	return stack.content[stack.contentLength-1]
}

func (stack *Stack) Push(value interface{}) {
	stack.content = append(stack.content, value)
	stack.contentLength++
}

func (stack *Stack) IsEmpty() bool {
	return stack.contentLength == 0
}

func NewStack() *Stack {
	return &Stack{
		content:       make([]interface{}, 0),
		contentLength: 0,
	}
}
