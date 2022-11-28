package parser

import (
	"github.com/shoriwe/plasma/pkg/ast"
	"github.com/shoriwe/plasma/pkg/lexer"
)

func (parser *Parser) parseIndexExpression(expression ast.Expression) (*ast.IndexExpression, error) {
	tokenizationError := parser.next()
	if tokenizationError != nil {
		return nil, tokenizationError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	// var rightIndex ast.Node

	var indexRange []ast.Expression
	switch {
	case parser.matchDirectValue(lexer.Colon):
		indexRange = append(indexRange, &ast.BasicLiteralExpression{
			Token: &lexer.Token{
				Contents:    []rune(lexer.NoneString),
				DirectValue: lexer.None,
				Kind:        lexer.NoneType,
			},
			Kind:        lexer.NoneType,
			DirectValue: lexer.None,
		})
		tokenizationError = parser.next()
		if tokenizationError != nil {
			return nil, tokenizationError
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
	default:
		break
	}
	//
	index, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := index.(ast.Expression); !ok {
		return nil, parser.expectingExpressionError(IndexExpression)
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	//
	if indexRange != nil {
		indexRange = append(indexRange, index.(ast.Expression))
		if !parser.matchDirectValue(lexer.CloseSquareBracket) {
			return nil, parser.newSyntaxError(IndexExpression)
		}
		tokenizationError = parser.next()
		if tokenizationError != nil {
			return nil, tokenizationError
		}
		return &ast.IndexExpression{
			Source: expression,
			Index: &ast.TupleExpression{
				Values: indexRange,
			},
		}, nil
	}
	switch {
	case parser.matchDirectValue(lexer.Colon):
		tokenizationError = parser.next()
		if tokenizationError != nil {
			return nil, tokenizationError
		}
		indexRange = append(indexRange, index.(ast.Expression))
		if parser.matchDirectValue(lexer.CloseSquareBracket) {
			tokenizationError = parser.next()
			if tokenizationError != nil {
				return nil, tokenizationError
			}
			indexRange = append(indexRange, &ast.BasicLiteralExpression{
				Token: &lexer.Token{
					Contents:    []rune(lexer.NoneString),
					DirectValue: lexer.None,
					Kind:        lexer.NoneType,
				},
				Kind:        lexer.NoneType,
				DirectValue: lexer.None,
			})
			return &ast.IndexExpression{
				Source: expression,
				Index: &ast.TupleExpression{
					Values: indexRange,
				},
			}, nil
		}
		// Parse next number
		index2, parsingError2 := parser.parseBinaryExpression(0)
		if parsingError2 != nil {
			return nil, parsingError2
		}
		if _, ok := index2.(ast.Expression); !ok {
			return nil, parser.expectingExpressionError(IndexExpression)
		}
		indexRange = append(indexRange, index2.(ast.Expression))
		if !parser.matchDirectValue(lexer.CloseSquareBracket) {
			return nil, parser.newSyntaxError(IndexExpression)
		}
		tokenizationError = parser.next()
		if tokenizationError != nil {
			return nil, tokenizationError
		}
		return &ast.IndexExpression{
			Source: expression,
			Index: &ast.TupleExpression{
				Values: indexRange,
			},
		}, nil
	}

	if !parser.matchDirectValue(lexer.CloseSquareBracket) {
		return nil, parser.newSyntaxError(IndexExpression)
	}

	tokenizationError = parser.next()
	if tokenizationError != nil {
		return nil, tokenizationError
	}
	return &ast.IndexExpression{
		Source: expression,
		Index:  index.(ast.Expression),
	}, nil
}
