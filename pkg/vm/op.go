package vm

import (
	"github.com/shoriwe/gruby/pkg/errors"
	"github.com/shoriwe/gruby/pkg/vm/utils"
)

type InstructionOP func(*utils.Stack) *errors.Error

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
