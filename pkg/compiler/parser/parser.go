package parser

import (
	"errors"
	"fmt"
	"github.com/shoriwe/gruby/pkg/compiler/ast"
	"github.com/shoriwe/gruby/pkg/compiler/lexer"
)

var binaryPrecedence = []string{
	"or",
	"and",
	"not",
	"==", "!=", ">", ">=", "<", "<=", "isinstanceof", "in",
	"|",
	"^",
	"&",
	"<<", ">>",
	"+", "-",
	"*", "/", "//", "%",
	"**",
}

type Parser struct {
	lexer        *lexer.Lexer
	complete     bool
	currentToken *lexer.Token
	nextToken    *lexer.Token
}

func (parser *Parser) hasNext() bool {
	return !parser.complete
}

func (parser *Parser) next() error {
	if parser.lexer.HasNext() {
		token_, tokenizingError := parser.lexer.Next()
		if tokenizingError != nil {
			return tokenizingError
		}
		parser.currentToken = token_
		if parser.lexer.HasNext() {
			token_, tokenizingError = parser.lexer.Peek()
			if tokenizingError != nil {
				return tokenizingError
			}
			parser.nextToken = token_
		}
	} else {
		parser.complete = true
	}
	return nil
}

func (parser *Parser) peek() *lexer.Token {
	return parser.nextToken
}

func (parser *Parser) check(kind int) bool {
	if parser.nextToken == nil {
		return false
	}
	return parser.nextToken.Kind == kind
}

func (parser *Parser) matchDirect(directValue int) bool {
	if parser.nextToken == nil {
		return false
	}
	return parser.nextToken.DirectValue == directValue
}

func (parser *Parser) matchString(value string) bool {
	if parser.nextToken == nil {
		return false
	}
	return parser.nextToken.String == value
}

func (parser *Parser) updateState() {
	parser.complete = true
}

func (parser *Parser) parseKeyboardStatement() (ast.Statement, error) {
	return nil, nil
}

func (parser *Parser) parseLiteral() (ast.Expression, error) {
	if parser.currentToken.Kind != lexer.Literal {
		return nil, errors.New(fmt.Sprintf("invalid kind of token %s at line %d", parser.currentToken.String, parser.currentToken.Line))
	}
	switch parser.currentToken.DirectValue {
	case lexer.SingleQuoteString, lexer.DoubleQuoteString, lexer.ByteString:
		return ast.BasicLiteralExpression{
			String: parser.currentToken.String,
			Kind:   parser.currentToken.Kind,
		}, nil
	case lexer.Integer, lexer.HexadecimalInteger, lexer.BinaryInteger, lexer.OctalInteger:
		return ast.BasicLiteralExpression{
			String: parser.currentToken.String,
			Kind:   parser.currentToken.Kind,
		}, nil
	case lexer.Float, lexer.ScientificFloat:
		return ast.BasicLiteralExpression{
			String: parser.currentToken.String,
			Kind:   parser.currentToken.Kind,
		}, nil
	case lexer.Boolean, lexer.NoneType:
		return ast.BasicLiteralExpression{
			String: parser.currentToken.String,
			Kind:   parser.currentToken.Kind,
		}, nil
	}
	return nil, errors.New(fmt.Sprintf("could not determine the directValue of token %s at line %d", parser.currentToken.String, parser.currentToken.Line))
}

func (parser *Parser) parseAssignStatementOrExpression() (ast.Node, error) {
	return nil, nil
}

func (parser *Parser) parseBinaryExpression() (ast.Node, error) {
	return nil, nil
}

func (parser *Parser) Parse() (*ast.Program, error) {
	result := &ast.Program{
		Begin: nil,
		End:   nil,
		Body:  nil,
	}
	switch parser.currentToken.Kind {
	case lexer.Keyboard: // Lambda and all statements
		switch parser.currentToken.DirectValue {
		case lexer.Lambda:
			break
		case lexer.BEGIN:
			break
		case lexer.END:
			break
		default:
			break
		}
	default: // Here it will be an assign statement or any expression
		parsedExpression, parsingError := parser.parseBinaryExpression()
		if parsingError != nil {
			return nil, parsingError
		}
		result.Body = append(result.Body, parsedExpression)
	}
	return result, nil
}

func NewParser(lexer_ *lexer.Lexer) (*Parser, error) {
	parser := &Parser{lexer_, false, nil, nil}
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return parser, nil
}
