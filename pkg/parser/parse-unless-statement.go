package parser

import (
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseUnlessStatement() (*ast2.UnlessStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}

	condition, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast2.IExpression); !ok {
		return nil, parser.expectingExpressionError(UnlessStatement)
	}
	if !parser.matchDirectValue(lexer2.NewLine) {
		return nil, parser.newSyntaxError(UnlessStatement)
	}
	// Parse Unless
	root := &ast2.UnlessStatement{
		Condition: condition.(ast2.IExpression),
		Body:      []ast2.Node{},
		Else:      []ast2.Node{},
	}
	var bodyNode ast2.Node
	for parser.hasNext() {
		if parser.matchKind(lexer2.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirectValue(lexer2.Elif) ||
				parser.matchDirectValue(lexer2.Else) ||
				parser.matchDirectValue(lexer2.End) {
				break
			}
			continue
		}
		bodyNode, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		root.Body = append(root.Body, bodyNode)
	}
	// Parse Elifs
	lastCondition := root
	if parser.matchDirectValue(lexer2.Elif) {
	elifBlocksParsingLoop:
		for parser.hasNext() {
			block := ast2.ElifBlock{
				Condition: nil,
				Body:      nil,
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			condition, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := condition.(ast2.IExpression); !ok {
				return nil, parser.expectingExpressionError(ElifBlock)
			}
			block.Condition = condition.(ast2.IExpression)
			if !parser.matchDirectValue(lexer2.NewLine) {
				return nil, parser.newSyntaxError(IfStatement)
			}
			for parser.hasNext() {
				if parser.matchKind(lexer2.Separator) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
					if parser.matchDirectValue(lexer2.Else) ||
						parser.matchDirectValue(lexer2.End) {
						root.ElifBlocks = append(root.ElifBlocks, block)
						break elifBlocksParsingLoop
					} else if parser.matchDirectValue(lexer2.Elif) {
						break
					}
					continue
				}
				bodyNode, parsingError = parser.parseBinaryExpression(0)
				if parsingError != nil {
					return nil, parsingError
				}
				block.Body = append(block.Body, bodyNode)
			}
			if !parser.matchDirectValue(lexer2.Elif) {
				return nil, parser.newSyntaxError(ElifBlock)
			}
			root.ElifBlocks = append(root.ElifBlocks, block)
		}
	}
	// Parse Default
	if parser.matchDirectValue(lexer2.Else) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		var elseBodyNode ast2.Node
		if !parser.matchDirectValue(lexer2.NewLine) {
			return nil, parser.newSyntaxError(ElseBlock)
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
			elseBodyNode, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			lastCondition.Else = append(lastCondition.Else, elseBodyNode)
		}
	}
	if !parser.matchDirectValue(lexer2.End) {
		return nil, parser.statementNeverEndedError(UnlessStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return root, nil
}
