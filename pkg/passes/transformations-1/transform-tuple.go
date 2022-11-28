package transformations_1

import (
	"github.com/shoriwe/plasma/pkg/ast2"
	"github.com/shoriwe/plasma/pkg/ast3"
)

func (transform *transformPass) Tuple(tuple *ast2.Tuple) *ast3.Tuple {
	values := make([]ast3.Expression, 0, len(tuple.Values))
	for _, value := range tuple.Values {
		values = append(values, transform.Expression(value))
	}
	return &ast3.Tuple{
		Values: values,
	}
}
