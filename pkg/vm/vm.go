package vm

import (
	"github.com/shoriwe/gruby/pkg/errors"
	"github.com/shoriwe/gruby/pkg/vm/object"
)

type Plasma struct {
	cursor int
}

func (plasma *Plasma) Execute(code []interface{}) (object.Object, *errors.Error) {
	return nil, nil
}

func NewPlasmaVM() *Plasma {
	return &Plasma{
		cursor: 0,
	}
}