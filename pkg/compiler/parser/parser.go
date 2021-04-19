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

func getPrecedenceWeight(operator string) int {
	for index, otherOperator := range binaryPrecedence {
		if otherOperator == operator {
			return index
		}
	}
	panic("invalid operator received")
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

		result := &ast.BasicLiteralExpression{
			String: parser.currentToken.String,
			Kind:   parser.currentToken.Kind,
		}
		tokenizingError := parser.next()
		if parser.check(lexer.Operator) {
			return parser.parseBinaryExpression(result)
		} else if parser.matchDirect(lexer.Dot) {
			return nil, nil // parser.parseSelectorExpression(result)
		}
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		return result, nil
	}
	return nil, errors.New(fmt.Sprintf("could not determine the directValue of token %s at line %d", parser.currentToken.String, parser.currentToken.Line))
}

func (parser *Parser) parseAssignStatementOrExpression() (ast.Node, error) {
	return nil, nil
}

func (parser *Parser) parseBinaryExpression(leftHandSide ast.Node) (ast.Node, error) {
	operator := parser.currentToken.String
	result := &ast.BinaryExpression{
		LeftHandSide:  leftHandSide,
		Operator:      operator,
		RightHandSide: nil,
	}
	switch leftHandSide.(type) {
	case *ast.BinaryExpression:
		leftWeight := getPrecedenceWeight(leftHandSide.(*ast.BinaryExpression).Operator)
		currentWeight := getPrecedenceWeight(operator)
		if leftWeight < currentWeight {
			result.LeftHandSide = leftHandSide.(*ast.BinaryExpression).RightHandSide
			leftHandSide.(*ast.BinaryExpression).RightHandSide = result
			result = leftHandSide.(*ast.BinaryExpression)
		}
	default:
		break
	}
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	rightHandSide, parsingError := parser.parseExpression()
	if parsingError != nil {
		return nil, parsingError
	}
	switch rightHandSide.(type) {
	case *ast.BinaryExpression:
		rightWeight := getPrecedenceWeight(rightHandSide.(*ast.BinaryExpression).Operator)
		currentWeight := getPrecedenceWeight(operator)
		if rightWeight < currentWeight {
			mustLeftParent := rightHandSide.(*ast.BinaryExpression)
			mustLeft := mustLeftParent.LeftHandSide
		mustLeftLoop:
			for ; ; {
				switch mustLeft.(type) {
				case (*ast.BinaryExpression):
					mustLeftParent = mustLeft.(*ast.BinaryExpression)
					mustLeft = mustLeft.(*ast.BinaryExpression).LeftHandSide
				default:
					break mustLeftLoop
				}
			}
			result.RightHandSide = mustLeft
			mustLeftParent.LeftHandSide = result
			result = rightHandSide.(*ast.BinaryExpression)
		} else {
			result.RightHandSide = rightHandSide
		}
	default:
		result.RightHandSide = rightHandSide
	}
	return result, nil
}

func (parser *Parser) parseUnaryExpression() (ast.Node, error) {
	return nil, nil
}

func (parser *Parser) parseExpression() (ast.Node, error) {
	var expression ast.Node
	var parsingError error
	switch parser.currentToken.Kind {
	case lexer.Literal:
		expression, parsingError = parser.parseLiteral()
	case lexer.IdentifierKind: // Only here assign can happen
		break
	case lexer.OpenParentheses:
		break
	case lexer.OpenSquareBracket:
		break
	case lexer.OpenBrace:
		break
	case lexer.Operator: // Unary
		expression, parsingError = parser.parseUnaryExpression()
	}
	if parsingError != nil {
		return nil, parsingError
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
			parsedExpression, parsingError := parser.parseExpression()
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
