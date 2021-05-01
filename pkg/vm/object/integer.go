package object

import (
	"github.com/shoriwe/gruby/pkg/errors"
	vmerrors "github.com/shoriwe/gruby/pkg/vm/vm-errors"
	"math/big"
	"reflect"
)

func (integer *Integer) Addition(other Object) (Object, *errors.Error) {
	switch other.(type) {
	case *Float:
		break
	case *Integer:
		result := big.NewInt(0)
		result.Add(result, integer.value)
		result.Add(result, other.(*Integer).value)
		return &Integer{
			value: result,
		}, nil
	}
	return nil, vmerrors.NewTypeError(IntegerTypeString, reflect.TypeOf(other).String())
}

func NewInteger(number string, base int) *Integer {
	n := big.NewInt(0)
	n.SetString(number, base)
	return &Integer{
		value: n,
	}
}
