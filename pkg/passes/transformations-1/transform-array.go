package transformations_1

import (
	"github.com/shoriwe/plasma/pkg/ast2"
	"github.com/shoriwe/plasma/pkg/ast3"
)

func (transform *transformPass) Array(array *ast2.Array) *ast3.Array {
	values := make([]ast3.Expression, 0, len(array.Values))
	for _, value := range array.Values {
		values = append(values, transform.Expression(value))
	}
	return &ast3.Array{
		Values: values,
	}
}
