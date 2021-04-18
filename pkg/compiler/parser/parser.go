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

func (parser *Parser) checkKind(kind int) bool {
	if parser.nextToken == nil {
		return false
	}
	return parser.nextToken.Kind == kind
}

func (parser *Parser) checkDirectValue(directValue int) bool {
	if parser.nextToken == nil {
		return false
	}
	return parser.nextToken.DirectValue == directValue
}

func (parser *Parser) updateState() {
	parser.complete = true
}

func (parser *Parser) parseKeyboardStatement() (ast.Statement, error) {
	return nil, nil
}

func (parser *Parser) parseLiteral() (ast.Expression, error) {
	switch parser.currentToken.Kind {
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
	return nil, errors.New(fmt.Sprintf("could not determine the kind of token %s at line %d", parser.currentToken.String, parser.currentToken.Line))
}

func (parser *Parser) parseAssignStatementOrExpression() (ast.Node, error) {
	return nil, nil
}

func (parser *Parser) Parse() (*ast.Program, error) {
	return nil, nil
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
