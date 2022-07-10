package transformations_1

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
	"reflect"
)

type transformPass struct {
	currentLabel int
}

func (transform *transformPass) Node(node ast2.Node) []ast3.Node {
	switch n := node.(type) {
	case ast2.Statement:
		return transform.Statement(n)
	case ast2.Expression:
		return transform.Expression(n)
	default:
		panic(fmt.Sprintf("unknown node type %s", reflect.TypeOf(n).String()))
	}
}

func Transform(program ast2.Program) ast3.Program {
	result := make(ast3.Program, 0, len(program))
	transform := transformPass{currentLabel: 1}
	for _, node := range program {
		result = append(result, transform.Node(node)...)
	}
	return result
}
