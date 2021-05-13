package errors

import "fmt"

// This should be changed
const (
	UnknownLine = 0
)

// Errors Types
const (
	UnknownVmOperationError    = "UnknownVMOperationError"
	StackOverflowError         = "StackOverflowError"
	TypeError                  = "TypeError"
	OperationNotSupportedError = "OperationNotSupportedError"
	NameNotFoundError          = "NameNotFoundError"
)

// Errors Messages
const (
	StackOverflowMessage         = "Memory stack has reached it's maximum size (platform uint size)"
	ExpectingNArgumentsMessage   = "Expecting %d arguments but received %d"
	OperationNotSupportedMessage = "Operation not supported"
	NameNotFoundMessage          = "\"Name not found\""
)

func NewOperationNotSupportedError() *Error {
	return New(UnknownLine, OperationNotSupportedMessage, OperationNotSupportedError)
}

func NewInvalidNumberOfArguments(received int, expecting int) *Error {
	return New(UnknownLine, fmt.Sprintf(ExpectingNArgumentsMessage, expecting, received), TypeError)
}

func NewTypeError(received string, expecting ...string) *Error {
	return New(UnknownLine, fmt.Sprintf("Expecting %s but received %s", expecting, received), TypeError)
}

func NewNameNotFoundError() *Error {
	return New(UnknownLine, NameNotFoundMessage, NameNotFoundError)
}

func NewStackOverflowError() *Error {
	return New(UnknownLine, StackOverflowMessage, StackOverflowError)
}

func NewUnknownVMOperationError(operation uint8) *Error {
	return New(UnknownLine, fmt.Sprintf("unknown operation with value %d", operation), UnknownVmOperationError)
}
