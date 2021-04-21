package parser

import (
	"errors"
	"fmt"
	"github.com/shoriwe/gruby/pkg/compiler/ast"
	"github.com/shoriwe/gruby/pkg/compiler/lexer"
)

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
	if parser.currentToken == nil {
		return false
	}
	return parser.currentToken.Kind == kind
}

func (parser *Parser) matchDirect(directValue int) bool {
	if parser.currentToken == nil {
		return false
	}
	return parser.currentToken.DirectValue == directValue
}

func (parser *Parser) matchString(value string) bool {
	if parser.currentToken == nil {
		return false
	}
	return parser.currentToken.String == value
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
	case lexer.SingleQuoteString, lexer.DoubleQuoteString, lexer.ByteString,
		lexer.Integer, lexer.HexadecimalInteger, lexer.BinaryInteger, lexer.OctalInteger,
		lexer.Float, lexer.ScientificFloat,
		lexer.Boolean, lexer.NoneType:
		currentToken := parser.currentToken
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		return &ast.BasicLiteralExpression{
			String: currentToken.String,
			Kind:   currentToken.Kind,
		}, nil
	}
	return nil, errors.New(fmt.Sprintf("could not determine the directValue of token %s at line %d", parser.currentToken.String, parser.currentToken.Line))
}

func (parser *Parser) parseBinaryExpression(precedence int) (ast.Node, error) {
	var leftHandSide ast.Node
	var rightHandSide ast.Node
	var parsingError error
	leftHandSide, parsingError = parser.parseUnaryExpression()
	if parsingError != nil {
		return nil, parsingError
	}
	for ; parser.currentToken.Kind == lexer.Operator ||
		parser.currentToken.Kind == lexer.Comparator; {
		operator := parser.currentToken
		operatorPrecedence := parser.currentToken.DirectValue
		if operatorPrecedence < precedence {
			return leftHandSide, nil
		}
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		rightHandSide, parsingError = parser.parseBinaryExpression(operatorPrecedence + 1)
		if parsingError != nil {
			return nil, parsingError
		}
		leftHandSide = &ast.BinaryExpression{
			LeftHandSide:  leftHandSide,
			Operator:      operator.String,
			RightHandSide: rightHandSide,
		}
	}
	return leftHandSide, nil
}

func (parser *Parser) parseUnaryExpression() (ast.Node, error) {
	// Do something to parse Unary
	if parser.check(lexer.Operator) {
		switch parser.currentToken.DirectValue {
		case lexer.Sub, lexer.Add, lexer.NegateBits, lexer.SignNot:
			operator := parser.currentToken.String
			tokenizingError := parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			x, parsingError := parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			return &ast.UnaryExpression{
				Operator: operator,
				X:        x,
			}, nil
		}
	}
	// Do something to parse Lambda
	// What about selectors?
	return parser.parsePrimaryExpression()
}

func (parser *Parser) parseOperand() (ast.Node, error) {
	switch parser.currentToken.Kind {
	case lexer.Literal:
		return parser.parseLiteral()
	case lexer.IdentifierKind:
		identifier := parser.currentToken.String
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		return &ast.Identifier{
			String: identifier,
		}, nil
	case lexer.Keyboard:
		switch parser.currentToken.DirectValue {
		case lexer.Lambda:
			break
			// return parser.parseLambda()
		}
	case lexer.Punctuation:
		switch parser.currentToken.DirectValue {
		case lexer.OpenParentheses:
			break
			// return parser.parseOpenParentheses()
		}
	}
	return nil, errors.New(fmt.Sprintf("unknown expression with token at line %d", parser.currentToken.Line))
}

func (parser *Parser) parseSelectorExpression(expression ast.Node) (ast.Node, error) {
	var selector ast.Node
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	identifier := parser.currentToken
	if identifier.Kind != lexer.IdentifierKind {
		return nil, errors.New(fmt.Sprintf("invalid selector at token in line %d", identifier.Line))
	}
	selector = &ast.SelectorExpression{
		X: expression,
		Identifier: &ast.Identifier{
			String: identifier.String,
		},
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	for ; parser.currentToken.DirectValue == lexer.Dot; {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		identifier = parser.currentToken
		if identifier.Kind != lexer.IdentifierKind {
			return nil, errors.New(fmt.Sprintf("invalid selector at token in line %d", identifier.Line))
		}
		selector = &ast.SelectorExpression{
			X: selector,
			Identifier: &ast.Identifier{
				String: identifier.String,
			},
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	return selector, nil
}

func (parser *Parser) parseCallExpression(expression ast.Node) (ast.Node, error) {
	return nil, nil
}

func (parser *Parser) parseIndexExpression(expression ast.Node) (ast.Node, error) {
	return nil, nil
}

func (parser *Parser) parsePrimaryExpression() (ast.Node, error) {
	var expression ast.Node
	var parsingError error
	expression, parsingError = parser.parseOperand()
	if parsingError != nil {
		return nil, parsingError
	}
expressionPendingLoop:
	for {
		switch parser.currentToken.DirectValue {
		case lexer.Dot: // Is selector
			expression, parsingError = parser.parseSelectorExpression(expression)
		case lexer.OpenParentheses: // Is function Call
			expression, parsingError = parser.parseCallExpression(expression)
		case lexer.OpenSquareBracket: // Is indexing
			expression, parsingError = parser.parseIndexExpression(expression)
		default:
			break expressionPendingLoop
		}
	}
	return expression, nil
}

func (parser *Parser) Parse() (*ast.Program, error) {
	result := &ast.Program{
		Begin: nil,
		End:   nil,
		Body:  nil,
	}
	for ; !parser.complete; {
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
			parsedExpression, parsingError := parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			result.Body = append(result.Body, parsedExpression)
		}
	}
	return result, nil
}

func NewParser(lexer_ *lexer.Lexer) (*Parser, error) {
	parser := &Parser{lexer_, false, nil, nil}
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return parser, nil
}
