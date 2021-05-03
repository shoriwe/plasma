package runtime

import (
	"github.com/shoriwe/gruby/pkg/errors"
)

type Iterable interface {
	Object
	Next() (Object, *errors.Error)
	HasNext() (*Bool, *errors.Error)
}
