package vm

import (
	"github.com/shoriwe/gruby/pkg/errors"
)

type InstructionOP func(*Stack) *errors.Error

const (
	// Binary Operations
	AddOP uint = iota
	SubOP
	DivOP
	MulOP
	PowOP
	ModOP
	NegateBitsOP
	BitAndOP
	BitOrOP
	BitXorOP
	BitLeftOP
	BitRightOP

	// Memory Operations
	PushOP
)
