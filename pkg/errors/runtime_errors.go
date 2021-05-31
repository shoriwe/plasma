package errors

import "fmt"

// This should be changed
const (
	UnknownLine = 0
)

// Errors Types
const (
	UnknownVmOperationError      = "UnknownVMOperationError"
	StackOverflowError           = "StackOverflowError"
	TypeError                    = "TypeError"
	OperationNotSupportedError   = "OperationNotSupportedError"
	NameNotFoundError            = "NameNotFoundError"
	IndexError                   = "IndexError"
	InvalidNumberDefinitionError = "InvalidNumberDefinitionError"
	StringHashingError           = "StringHashingError"
	UnhashableTypeError          = "UnhashableTypeError"
	KeyNotFoundError             = "KeyNotFoundError"
	GoRuntimeError               = "GoRuntimeError"
)

// Errors Messages
const (
	StackOverflowMessage         = "Memory stack has reached it's maximum size (platform uint size)"
	ExpectingNArgumentsMessage   = "Expecting %d arguments but received %d"
	OperationNotSupportedMessage = "Operation not supported"
	NameNotFoundMessage          = "\"Name not found\""
)

func NewGoRuntimeError(e error) *Error {
	return New(UnknownLine, e.Error(), GoRuntimeError)
}
func NewKeyNotFoundError() *Error {
	return New(UnknownLine, "Key not found in hash table", KeyNotFoundError)
}

func NewUnhashableTypeError(line int) *Error {
	return New(line, "Object Type is unhashable", UnhashableTypeError)
}
func NewHashingStringError() *Error {
	return New(UnknownLine, "Unknown cause for hashing error", StringHashingError)
}
func NewInvalidFloatDefinition(line int, s string) *Error {
	return New(line, fmt.Sprintf("Invalid Float definition: %s", s), InvalidNumberDefinitionError)
}

func NewInvalidIntegerDefinition(line int, s string) *Error {
	return New(line, fmt.Sprintf("Invalid Integer definition: %s", s), InvalidNumberDefinitionError)
}

func NewIndexOutOfRange(line int, length int, index int) *Error {
	return New(line, fmt.Sprintf("index %d out of bound for a %d container", index, length), IndexError)
}

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
