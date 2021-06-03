package tools

import "github.com/shoriwe/gruby/pkg/errors"

func CalcIndex(index64 int64, length int) (int, *errors.Error) {
	index := int(index64)
	if length <= index {
		return 0, errors.NewIndexOutOfRange(errors.UnknownLine, length, index)
	}
	if index < 0 {
		index = length + index
		if index < 0 {
			return 0, errors.NewIndexOutOfRange(errors.UnknownLine, length, index)
		}
	}
	return index, nil
}
