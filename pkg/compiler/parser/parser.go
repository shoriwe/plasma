package parser

import (
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

func (parser *Parser) updateState() {
	parser.complete = true
}

func (parser *Parser) parseKeyboardStatement() (ast.Statement, error) {
	return nil, nil
}

func (parser *Parser) parseLiteral() (ast.Expression, error) {
	return nil, nil
}

func (parser *Parser) parseAssignStatementOrExpression() (ast.Node, error) {
	return nil, nil
}

func (parser *Parser) Parse() (*ast.Program, error) {
	return nil, nil
}

func NewParser(lexer_ *lexer.Lexer) *Parser {
	return &Parser{lexer_, false, nil, nil}
}
