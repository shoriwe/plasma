package utils

type Stack struct {
	content       []interface{}
	contentLength int
}

func (stack *Stack) Pop() interface{} {
	stack.contentLength--
	result := stack.content[stack.contentLength]
	stack.content = stack.content[:stack.contentLength]
	return result
}

func (stack *Stack) Push(value interface{}) {
	stack.content = append(stack.content, value)
	stack.contentLength++
}

func (stack *Stack) IsEmpty() bool {
	return stack.contentLength == 0
}
