package op

import (
	"github.com/shoriwe/gruby/pkg/errors"
	"github.com/shoriwe/gruby/pkg/vm/utils"
)

type InstructionOP func(*utils.Stack) *errors.Error

const (
	AddOP uint16 = iota
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
)
