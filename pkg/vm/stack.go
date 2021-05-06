package vm

import (
	"github.com/shoriwe/gruby/pkg/errors"
	"math/bits"
)

type Stack interface {
	Push(interface{}) *errors.Error
	Pop() interface{}
	HashNext() bool
	Clear()
	Peek() interface{}
}

type ArrayStack struct {
	content       []interface{}
	currentLength uint
}

func (a *ArrayStack) Push(value interface{}) *errors.Error {
	if a.currentLength == bits.UintSize {
		return errors.NewStackOverflowError()
	}
	a.content = append(a.content, value)
	a.currentLength++
	return nil
}

func (a *ArrayStack) Pop() interface{} {
	a.currentLength--
	result := a.content[a.currentLength]
	a.content = a.content[:a.currentLength]
	return result
}

func (a *ArrayStack) HashNext() bool {
	return a.currentLength > 0
}

func (a *ArrayStack) Peek() interface{} {
	return a.content[a.currentLength-1]
}

func (a *ArrayStack) Clear() {
	a.currentLength = 0
	a.content = make([]interface{}, 0)
}

func NewArrayStack() *ArrayStack {
	return &ArrayStack{
		content:       make([]interface{}, 0),
		currentLength: 0,
	}
}
