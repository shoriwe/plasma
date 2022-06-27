package parser

import (
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseSwitchStatement() (*ast2.SwitchStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}

	target, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := target.(ast2.IExpression); !ok {
		return nil, parser.expectingExpressionError(SwitchStatement)
	}
	if !parser.matchDirectValue(lexer2.NewLine) {
		return nil, parser.newSyntaxError(SwitchStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	// parse Cases
	var caseBlocks []*ast2.CaseBlock
	if parser.matchDirectValue(lexer2.Case) {
		for parser.hasNext() {
			if parser.matchDirectValue(lexer2.Default) ||
				parser.matchDirectValue(lexer2.End) {
				break
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			newLinesRemoveError = parser.removeNewLines()
			if newLinesRemoveError != nil {
				return nil, newLinesRemoveError
			}
			var cases []ast2.IExpression
			var caseTarget ast2.Node
			for parser.hasNext() {
				caseTarget, parsingError = parser.parseBinaryExpression(0)
				if parsingError != nil {
					return nil, parsingError
				}
				if _, ok := caseTarget.(ast2.IExpression); !ok {
					return nil, parser.expectingExpressionError(CaseBlock)
				}
				cases = append(cases, caseTarget.(ast2.IExpression))
				if parser.matchDirectValue(lexer2.NewLine) {
					break
				} else if parser.matchDirectValue(lexer2.Comma) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
				} else {
					return nil, parser.newSyntaxError(CaseBlock)
				}
			}
			if !parser.matchDirectValue(lexer2.NewLine) {
				return nil, parser.newSyntaxError(CaseBlock)
			}
			// Targets Body
			var caseBody []ast2.Node
			var caseBodyNode ast2.Node
			for parser.hasNext() {
				if parser.matchKind(lexer2.Separator) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
					if parser.matchDirectValue(lexer2.Case) ||
						parser.matchDirectValue(lexer2.Default) ||
						parser.matchDirectValue(lexer2.End) {
						break
					}
					continue
				}
				caseBodyNode, parsingError = parser.parseBinaryExpression(0)
				if parsingError != nil {
					return nil, parsingError
				}
				caseBody = append(caseBody, caseBodyNode)
			}
			// Targets block
			caseBlocks = append(caseBlocks, &ast2.CaseBlock{
				Cases: cases,
				Body:  caseBody,
			})
		}
	}
	// Parse Default
	var defaultBody []ast2.Node
	if parser.matchDirectValue(lexer2.Default) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		if !parser.matchDirectValue(lexer2.NewLine) {
			return nil, parser.newSyntaxError(DefaultBlock)
		}
		var defaultBodyNode ast2.Node
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
			defaultBodyNode, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			defaultBody = append(defaultBody, defaultBodyNode)
		}
	}
	// Finally detect valid end
	if !parser.matchDirectValue(lexer2.End) {
		return nil, parser.statementNeverEndedError(SwitchStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast2.SwitchStatement{
		Target:     target.(ast2.IExpression),
		CaseBlocks: caseBlocks,
		Default:    defaultBody,
	}, nil
}
