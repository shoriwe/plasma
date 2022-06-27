package parser

import (
	"github.com/shoriwe/gplasma/pkg/compiler/ast"
	"github.com/shoriwe/gplasma/pkg/compiler/lexer"
)

func (parser *Parser) parseTryStatement() (*ast.TryStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(TryStatement)
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError error
	for parser.hasNext() {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirectValue(lexer.End) ||
				parser.matchDirectValue(lexer.Except) ||
				parser.matchDirectValue(lexer.Else) ||
				parser.matchDirectValue(lexer.Finally) {
				break
			}
			continue
		}
		bodyNode, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		body = append(body, bodyNode)
	}
	var exceptBlocks []*ast.ExceptBlock
	for parser.hasNext() {
		if !parser.matchDirectValue(lexer.Except) {
			break
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		var targets []ast.IExpression
		var target ast.Node
		for parser.hasNext() {
			if parser.matchDirectValue(lexer.NewLine) ||
				parser.matchDirectValue(lexer.As) {
				break
			}
			target, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := target.(ast.IExpression); !ok {
				return nil, parser.newSyntaxError(ExceptBlock)
			}
			targets = append(targets, target.(ast.IExpression))
			if parser.matchDirectValue(lexer.Comma) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
			}
		}
		var captureName *ast.Identifier
		if parser.matchDirectValue(lexer.As) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			newLinesRemoveError := parser.removeNewLines()
			if newLinesRemoveError != nil {
				return nil, newLinesRemoveError
			}
			if !parser.matchKind(lexer.IdentifierKind) {
				return nil, parser.newSyntaxError(ExceptBlock)
			}
			captureName = &ast.Identifier{
				Token: parser.currentToken,
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		}
		if !parser.matchDirectValue(lexer.NewLine) {
			return nil, parser.newSyntaxError(TryStatement)
		}
		var exceptBody []ast.Node
		var exceptBodyNode ast.Node
		for parser.hasNext() {
			if parser.matchKind(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.matchDirectValue(lexer.End) ||
					parser.matchDirectValue(lexer.Except) ||
					parser.matchDirectValue(lexer.Else) ||
					parser.matchDirectValue(lexer.Finally) {
					break
				}
				continue
			}
			exceptBodyNode, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			exceptBody = append(exceptBody, exceptBodyNode)
		}
		exceptBlocks = append(exceptBlocks, &ast.ExceptBlock{
			Targets:     targets,
			CaptureName: captureName,
			Body:        exceptBody,
		})
	}
	var elseBody []ast.Node
	var elseBodyNode ast.Node
	if parser.matchDirectValue(lexer.Else) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		if !parser.matchDirectValue(lexer.NewLine) {
			return nil, parser.newSyntaxError(ElseBlock)
		}
		for parser.hasNext() {
			if parser.matchKind(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.matchDirectValue(lexer.End) ||
					parser.matchDirectValue(lexer.Finally) {
					break
				}
				continue
			}
			elseBodyNode, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			elseBody = append(elseBody, elseBodyNode)
		}
	}
	var finallyBody []ast.Node
	var finallyBodyNode ast.Node
	if parser.matchDirectValue(lexer.Finally) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		if !parser.matchDirectValue(lexer.NewLine) {
			return nil, parser.newSyntaxError(FinallyBlock)
		}
		for parser.hasNext() {
			if parser.matchKind(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.matchDirectValue(lexer.End) {
					break
				}
				continue
			}
			finallyBodyNode, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			finallyBody = append(finallyBody, finallyBodyNode)
		}
	}
	if !parser.matchDirectValue(lexer.End) {
		return nil, parser.newSyntaxError(TryStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.TryStatement{
		Body:         body,
		ExceptBlocks: exceptBlocks,
		Else:         elseBody,
		Finally:      finallyBody,
	}, nil
}
