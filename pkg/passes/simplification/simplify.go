package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

type simplify struct {
	currentAnonIdent uint
}

func (simp *simplify) simplifyNode(node ast.Node) ast2.Node {
	switch n := node.(type) {
	case ast.Statement:
		return simp.simplifyStatement(n)
	case ast.Expression:
		return simp.simplifyExpression(n)
	default:
		panic("unknown node type")
	}
}

func Simplify(program *ast.Program) ast2.Program {
	var (
		begin []ast2.Node
		body  []ast2.Node
		end   []ast2.Node
	)
	simp := simplify{currentAnonIdent: 1}
	if program.Begin != nil {
		begin = make([]ast2.Node, 0, len(program.Begin.Body))
		for _, node := range program.Begin.Body {
			begin = append(begin, simp.simplifyNode(node))
		}
	}
	body = make([]ast2.Node, 0, len(program.Body))
	for _, node := range program.Body {
		body = append(body, simp.simplifyNode(node))
	}
	if program.End != nil {
		begin = make([]ast2.Node, 0, len(program.End.Body))
		for _, node := range program.End.Body {
			begin = append(begin, simp.simplifyNode(node))
		}
	}
	result := make(ast2.Program, 0, len(begin)+len(body)+len(end))
	result = append(result, begin...)
	result = append(result, body...)
	result = append(result, end...)
	return result
}
