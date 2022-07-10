package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) IfOneLiner(iol *ast2.IfOneLiner) *ast3.IfOneLiner {
	return &ast3.IfOneLiner{
		Condition: transform.Expression(iol.Condition),
		Result:    transform.Expression(iol.Result),
		Else:      transform.Expression(iol.Else),
	}
}
