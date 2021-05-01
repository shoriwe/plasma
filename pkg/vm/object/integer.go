package object

import (
	"github.com/shoriwe/gruby/pkg/errors"
	"math/big"
)

func (integer *Integer) Addition(other Object) (Object, *errors.Error) {
	switch other.(type) {
	case *Float:
		break
	case *Integer:
		break
	}
	return nil, nil
}

func NewInteger(number string, base int) *Integer {
	n := big.NewInt(0)
	n.SetString(number, base)
	return &Integer{
		value: n,
	}
}
