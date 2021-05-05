package errors

import "fmt"

// This should be changed
const (
	UnknownLine = 0
)

// Errors Types
const (
	TypeError                  = "TypeError"
	OperationNotSupportedError = "OperationNotSupportedError"
	NameNotFoundError          = "NameNotFoundError"
)

// Errors Messages
const (
	ExpectingNArgumentsMessage   = "Expecting %d arguments but received %d"
	OperationNotSupportedMessage = "Operation not supported"
	NameNotFoundMessage          = "\"Name not found\""
)

func NewOperationNotSupportedError(message string, typeReceived string) *Error {
	return New(UnknownLine, fmt.Sprintf("%s for type %s", message, typeReceived), OperationNotSupportedError)
}

func NewInvalidNumberOfArguments(received int, expecting int) *Error {
	return New(UnknownLine, fmt.Sprintf(ExpectingNArgumentsMessage, expecting, received), TypeError)
}

func NewTypeError(expecting []string, received string) *Error {
	return New(UnknownLine, fmt.Sprintf("Expecting %s but received %s", expecting, received), TypeError)
}

func NewNameNotFoundError() *Error {
	return New(-1, NameNotFoundMessage, NameNotFoundError)
}
