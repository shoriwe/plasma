package parser

import (
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseTryStatement() (*ast2.TryStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirectValue(lexer2.NewLine) {
		return nil, parser.newSyntaxError(TryStatement)
	}
	var body []ast2.Node
	var bodyNode ast2.Node
	var parsingError error
	for parser.hasNext() {
		if parser.matchKind(lexer2.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirectValue(lexer2.End) ||
				parser.matchDirectValue(lexer2.Except) ||
				parser.matchDirectValue(lexer2.Else) ||
				parser.matchDirectValue(lexer2.Finally) {
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
	var exceptBlocks []*ast2.ExceptBlock
	for parser.hasNext() {
		if !parser.matchDirectValue(lexer2.Except) {
			break
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		var targets []ast2.IExpression
		var target ast2.Node
		for parser.hasNext() {
			if parser.matchDirectValue(lexer2.NewLine) ||
				parser.matchDirectValue(lexer2.As) {
				break
			}
			target, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := target.(ast2.IExpression); !ok {
				return nil, parser.newSyntaxError(ExceptBlock)
			}
			targets = append(targets, target.(ast2.IExpression))
			if parser.matchDirectValue(lexer2.Comma) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
			}
		}
		var captureName *ast2.Identifier
		if parser.matchDirectValue(lexer2.As) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			newLinesRemoveError := parser.removeNewLines()
			if newLinesRemoveError != nil {
				return nil, newLinesRemoveError
			}
			if !parser.matchKind(lexer2.IdentifierKind) {
				return nil, parser.newSyntaxError(ExceptBlock)
			}
			captureName = &ast2.Identifier{
				Token: parser.currentToken,
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		}
		if !parser.matchDirectValue(lexer2.NewLine) {
			return nil, parser.newSyntaxError(TryStatement)
		}
		var exceptBody []ast2.Node
		var exceptBodyNode ast2.Node
		for parser.hasNext() {
			if parser.matchKind(lexer2.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.matchDirectValue(lexer2.End) ||
					parser.matchDirectValue(lexer2.Except) ||
					parser.matchDirectValue(lexer2.Else) ||
					parser.matchDirectValue(lexer2.Finally) {
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
		exceptBlocks = append(exceptBlocks, &ast2.ExceptBlock{
			Targets:     targets,
			CaptureName: captureName,
			Body:        exceptBody,
		})
	}
	var elseBody []ast2.Node
	var elseBodyNode ast2.Node
	if parser.matchDirectValue(lexer2.Else) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		if !parser.matchDirectValue(lexer2.NewLine) {
			return nil, parser.newSyntaxError(ElseBlock)
		}
		for parser.hasNext() {
			if parser.matchKind(lexer2.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.matchDirectValue(lexer2.End) ||
					parser.matchDirectValue(lexer2.Finally) {
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
	var finallyBody []ast2.Node
	var finallyBodyNode ast2.Node
	if parser.matchDirectValue(lexer2.Finally) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		if !parser.matchDirectValue(lexer2.NewLine) {
			return nil, parser.newSyntaxError(FinallyBlock)
		}
		for parser.hasNext() {
			if parser.matchKind(lexer2.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.matchDirectValue(lexer2.End) {
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
	if !parser.matchDirectValue(lexer2.End) {
		return nil, parser.newSyntaxError(TryStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast2.TryStatement{
		Body:         body,
		ExceptBlocks: exceptBlocks,
		Else:         elseBody,
		Finally:      finallyBody,
	}, nil
}
