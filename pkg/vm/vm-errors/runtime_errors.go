package vm_errors

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/errors"
)

// Runtime Errors
const (
	TypeError              = "TypeError"
	AttributeNotFoundError = "AttributeNotFoundError"
	NilObjectError         = "NilObjectError"
)

const UnknownLine = -1

func NewRuntimeError(errorType string, message string) *errors.Error {
	return errors.New(UnknownLine, message, errorType)
}

func NewTypeError(expecting string, received string) *errors.Error {
	return NewRuntimeError(TypeError, fmt.Sprintf("Expecting object of type: %s but Recevied: %s", expecting, received))
}

func NewAttributeNotFound(ownerName string, name string) *errors.Error {
	return NewRuntimeError(AttributeNotFoundError, fmt.Sprintf("%s has no attribute with name %s", ownerName, name))
}
