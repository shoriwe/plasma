package runtime

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/errors"
)

// Runtime Errors
const (
	TypeError              = "TypeError"
	AttributeNotFoundError = "AttributeNotFoundError"
	NilObjectError         = "NilObjectError"
	UnknownOP              = "UnknownOP"
	MethodNotImplemented   = "MethodNotImplemented"
)

const UnknownLine = -1

func NewRuntimeError(errorType string, message string) *errors.Error {
	return errors.New(UnknownLine, message, errorType)
}

func NewTypeError(received string, expecting ...string) *errors.Error {
	if len(expecting) == 1 {
		return NewRuntimeError(TypeError, fmt.Sprintf("Expecting object of type: %s but recevied object of type: %s", expecting[0], received))
	}
	return NewRuntimeError(TypeError, fmt.Sprintf("Expecting object of types: %s but received object of type: %s", expecting, received))
}

func NewAttributeNotFound(ownerName string, name string) *errors.Error {
	return NewRuntimeError(AttributeNotFoundError, fmt.Sprintf("%s has no attribute with name %s", ownerName, name))
}

func NewMethodNotImplemented(methodName string) *errors.Error {
	return NewRuntimeError(MethodNotImplemented, fmt.Sprintf("method %s not implemented", methodName))
}
