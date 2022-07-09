package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func simplifyRaise(raise *ast.RaiseStatement) *ast2.Raise {
	return &ast2.Raise{
		Statement: nil,
		X:         simplifyExpression(raise.X),
	}
}

func simplifyTry(try *ast.TryStatement) *ast2.Try {
	var (
		body         = make([]ast2.Node, 0, len(try.Body))
		exceptBlocks = make([]*ast2.ExceptBlock, 0, len(try.ExceptBlocks))
		else_        = make([]ast2.Node, 0, len(try.Else))
		finally      = make([]ast2.Node, 0, len(try.Finally))
	)
	for _, node := range try.Body {
		body = append(body, simplifyNode(node))
	}
	for _, except := range try.ExceptBlocks {
		targets := make([]ast2.Expression, 0, len(except.Targets))
		for _, target := range except.Targets {
			targets = append(targets, simplifyExpression(target))
		}
		exceptBody := make([]ast2.Node, 0, len(except.Body))
		for _, node := range except.Body {
			exceptBody = append(exceptBody, simplifyNode(node))
		}
		exceptBlocks = append(exceptBlocks, &ast2.ExceptBlock{
			Targets:     targets,
			CaptureName: simplifyIdentifier(except.CaptureName),
			Body:        exceptBody,
		})
	}
	for _, node := range try.Else {
		else_ = append(else_, simplifyNode(node))
	}
	for _, node := range try.Finally {
		finally = append(finally, simplifyNode(node))
	}
	return &ast2.Try{
		Body:         body,
		ExceptBlocks: exceptBlocks,
		Else:         else_,
		Finally:      finally,
	}
}
