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
}

func (parser *Parser) hasNext() bool {
	return !parser.complete
}

func (parser *Parser) next() error {
	token, tokenizingError := parser.lexer.Next()
	if tokenizingError != nil {
		return tokenizingError
	}
	if token.Kind == lexer.EOF {
		parser.complete = true
	}
	parser.currentToken = token
	return nil
}

func (parser *Parser) matchDirect(directValue int) bool {
	if parser.currentToken == nil {
		return false
	}
	return parser.currentToken.DirectValue == directValue
}

func (parser *Parser) matchKind(kind int) bool {
	if parser.currentToken == nil {
		return false
	}
	return parser.currentToken.Kind == kind
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
	if !parser.matchKind(lexer.Literal) &&
		!parser.matchKind(lexer.Boolean) &&
		!parser.matchKind(lexer.NoneType) {
		return nil, errors.New(fmt.Sprintf("invalid kind of token %s at line %d", parser.currentToken.String, parser.currentToken.Line))
	}
	switch parser.currentToken.DirectValue {
	case lexer.SingleQuoteString, lexer.DoubleQuoteString, lexer.ByteString,
		lexer.Integer, lexer.HexadecimalInteger, lexer.BinaryInteger, lexer.OctalInteger,
		lexer.Float, lexer.ScientificFloat,
		lexer.True, lexer.False, lexer.None:
		currentToken := parser.currentToken
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		return &ast.BasicLiteralExpression{
			Token:       currentToken,
			Kind:        currentToken.Kind,
			DirectValue: currentToken.DirectValue,
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
	for ; !parser.complete; {
		if !parser.matchKind(lexer.Operator) &&
			!parser.matchKind(lexer.Comparator) {
			break
		}
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
			Operator:      operator,
			RightHandSide: rightHandSide,
		}
	}
	return leftHandSide, nil
}

func (parser *Parser) parseUnaryExpression() (ast.Node, error) {
	// Do something to parse Unary
	if parser.matchKind(lexer.Operator) {
		switch parser.currentToken.DirectValue {
		case lexer.Sub, lexer.Add, lexer.NegateBits, lexer.SignNot, lexer.Not:
			operator := parser.currentToken
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
		case lexer.BitWiseAnd:
			tokenizingError := parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			x, parsingError := parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			return &ast.PointerExpression{
				X:        x,
			}, nil
		case lexer.Star:
			tokenizingError := parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			x, parsingError := parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			return &ast.StarExpression{
				X:        x,
			}, nil
		}
	}
	// Do something to parse Lambda
	// What about selectors?
	return parser.parsePrimaryExpression()
}

func (parser *Parser) parseLambdaExpression() (ast.Node, error) {
	var arguments []*ast.Identifier
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	for ; !parser.complete; {
		if !parser.matchKind(lexer.IdentifierKind) {
			break
		}
		arguments = append(arguments, &ast.Identifier{
			Token: parser.currentToken,
		})
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		if parser.matchDirect(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if !parser.matchDirect(lexer.Colon) {
			return nil, errors.New(fmt.Sprintf("invalid lambda definition at line %d", parser.currentToken.Line))
		}
	}
	if !parser.matchDirect(lexer.Colon) {
		return nil, errors.New(fmt.Sprintf("invalid lambda definition at line %d", parser.currentToken.Line))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	code, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	return &ast.LambdaExpression{
		Arguments: arguments,
		Code:      code,
	}, nil
}

func (parser *Parser) parseParentheses() (ast.Node, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if parser.matchDirect(lexer.CloseParentheses) {
		return nil, errors.New(fmt.Sprintf("syntax error: empty parentheses expression at line %d", parser.currentToken.Line))
	}
	firstExpression, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if parser.matchDirect(lexer.CloseParentheses) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		return &ast.ParenthesesExpression{
			X: firstExpression,
		}, nil
	}
	if !parser.matchDirect(lexer.Comma) {
		return nil, errors.New(fmt.Sprintf("syntax error: empty parentheses expression at line %d", parser.currentToken.Line))
	}
	var values []ast.Expression
	values = append(values, firstExpression)
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var nextValue ast.Node
	for ; !parser.complete; {
		if parser.matchDirect(lexer.CloseParentheses) {
			break
		}
		nextValue, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		values = append(values, nextValue)
		if parser.matchDirect(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		}
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.TupleExpression{
		Values: values,
	}, nil
}
func (parser *Parser) parseArrayExpression() (ast.Node, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, nil
	}
	var values []ast.Expression
	for ; !parser.complete; {
		if parser.matchDirect(lexer.CloseSquareBracket) {
			break
		}
		value, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		values = append(values, value)
		if parser.matchDirect(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, nil
			}
		}
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, nil
	}
	return &ast.ArrayExpression{
		Values: values,
	}, nil
}

func (parser *Parser) parseHashExpression() (ast.Node, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, nil
	}
	var values []*ast.KeyValue
	var leftHandSide ast.Node
	var rightHandSide ast.Node
	var parsingError error
	for ; !parser.complete; {
		if parser.matchDirect(lexer.CloseBrace) {
			break
		}
		leftHandSide, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if !parser.matchDirect(lexer.Colon) {
			return nil, errors.New(fmt.Sprintf("syntax error: invalid hash definition at line %d", parser.currentToken.Line))
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, nil
		}
		rightHandSide, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		values = append(values, &ast.KeyValue{
			Key:   leftHandSide,
			Value: rightHandSide,
		})
		if parser.matchDirect(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, nil
			}
		}
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, nil
	}
	return &ast.HashExpression{
		Values: values,
	}, nil
}

func (parser *Parser) parseWhileLoop() (ast.Statement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	condition, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.Expression); !ok {
		return nil, errors.New(fmt.Sprintf("invalid while loop declaration at line %d", parser.currentToken.Line))
	}
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid while loop declaration at line %d", parser.currentToken.Line))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var whileChild ast.Node
	var body []ast.Node
	for ; !parser.complete; {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirect(lexer.End) {
				break
			}
			continue
		}
		whileChild, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		body = append(body, whileChild)
	}
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("while statement never closed at line %d", parser.currentToken.Line))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.WhileLoopStatement{
		Condition: condition,
		Body:      body,
	}, nil
}

func (parser *Parser) parseForLoop() (ast.Statement, error) {
	return nil, nil
}

func (parser *Parser) parseUntilLoop() (ast.Statement, error) {
	return nil, nil
}

func (parser *Parser) parseIfStatement() (ast.Statement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	condition, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid if statement declaration at line %d", parser.currentToken.Line))
	}
	var body []ast.Node
	var elifBlocks []*ast.ElifBlock
	var elseBody []ast.Node
	var node ast.Node
	var elifCondition ast.Node
	for ; !parser.complete; {
		if parser.matchDirect(lexer.End) {
			break
		}
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			continue
		}
		switch parser.currentToken.DirectValue {
		case lexer.Elif:
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			elifCondition, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if !parser.matchDirect(lexer.NewLine) {
				return nil, errors.New(fmt.Sprintf("invalid if statement declaration at line %d", parser.currentToken.Line))
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			var elifBody []ast.Node
			for ; !parser.complete; {
				if parser.matchDirect(lexer.Else) ||
					parser.matchDirect(lexer.Elif) ||
					parser.matchDirect(lexer.End) {
					break
				}
				if parser.matchKind(lexer.Separator) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
					continue
				}
				node, parsingError = parser.parseBinaryExpression(0)
				if parsingError != nil {
					return nil, parsingError
				}
				elifBody = append(elifBody, node)
			}
			elifBlocks = append(elifBlocks, &ast.ElifBlock{
				Condition: elifCondition,
				Body:      elifBody,
			})
		case lexer.Else:
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if !parser.matchDirect(lexer.NewLine) {
				return nil, errors.New(fmt.Sprintf("invalid else statement declaration in if statement at line %d", parser.currentToken.Line))
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			var elseNode ast.Node
			for ; !parser.complete; {
				if parser.matchDirect(lexer.End) {
					break
				}
				if parser.matchKind(lexer.Separator) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
					continue
				}
				elseNode, parsingError = parser.parseBinaryExpression(0)
				if parsingError != nil {
					return nil, parsingError
				}
				elseBody = append(elseBody, elseNode)
			}
			break
		default:
			node, parsingError = parser.parseBinaryExpression(0)
			body = append(body, node)
		}
	}
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("never closed if statement at line %d", parser.currentToken.Line))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.IfStatement{
		Condition:  condition,
		Body:       body,
		ElifBlocks: elifBlocks,
		Else:       elseBody,
	}, nil
}

func (parser *Parser) parseUnlessStatement() (ast.Statement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	condition, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid if statement declaration at line %d", parser.currentToken.Line))
	}
	var body []ast.Node
	var elifBlocks []*ast.ElifBlock
	var elseBody []ast.Node
	var node ast.Node
	var elifCondition ast.Node
	for ; !parser.complete; {
		if parser.matchDirect(lexer.End) {
			break
		}
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			continue
		}
		switch parser.currentToken.DirectValue {
		case lexer.Elif:
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			elifCondition, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if !parser.matchDirect(lexer.NewLine) {
				return nil, errors.New(fmt.Sprintf("invalid if statement declaration at line %d", parser.currentToken.Line))
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			var elifBody []ast.Node
			for ; !parser.complete; {
				if parser.matchDirect(lexer.Else) ||
					parser.matchDirect(lexer.Elif) ||
					parser.matchDirect(lexer.End) {
					break
				}
				if parser.matchKind(lexer.Separator) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
					continue
				}
				node, parsingError = parser.parseBinaryExpression(0)
				if parsingError != nil {
					return nil, parsingError
				}
				elifBody = append(elifBody, node)
			}
			elifBlocks = append(elifBlocks, &ast.ElifBlock{
				Condition: elifCondition,
				Body:      elifBody,
			})
		case lexer.Else:
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if !parser.matchDirect(lexer.NewLine) {
				return nil, errors.New(fmt.Sprintf("invalid else statement declaration in if statement at line %d", parser.currentToken.Line))
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			var elseNode ast.Node
			for ; !parser.complete; {
				if parser.matchDirect(lexer.End) {
					break
				}
				if parser.matchKind(lexer.Separator) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
					continue
				}
				elseNode, parsingError = parser.parseBinaryExpression(0)
				if parsingError != nil {
					return nil, parsingError
				}
				elseBody = append(elseBody, elseNode)
			}
			break
		default:
			node, parsingError = parser.parseBinaryExpression(0)
			body = append(body, node)
		}
	}
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("never closed if statement at line %d", parser.currentToken.Line))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.UnlessStatement{
		Condition:  condition,
		Body:       body,
		ElifBlocks: elifBlocks,
		Else:       elseBody,
	}, nil
}

func (parser *Parser) parseSwitchStatement() (ast.Statement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	target, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid switch statement at line %d", parser.currentToken.Line))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var caseBlocks []*ast.CaseBlock
	var elseBody []ast.Node
	var elseChild ast.Node
	var caseChild ast.Node
	var caseTarget ast.Expression
	// Parse Body
	for ; !parser.complete; {
		if parser.matchDirect(lexer.End) {
			break
		} else if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			continue
		}
		if parser.matchDirect(lexer.Case) {
			// Parse Case
			var cases []ast.Expression
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			for ; !parser.complete; {
				if parser.matchDirect(lexer.NewLine) {
					break
				}
				caseTarget, parsingError = parser.parseBinaryExpression(0)
				if parsingError != nil {
					return nil, parsingError
				}
				cases = append(cases, caseTarget)
				if parser.matchDirect(lexer.Comma) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
				}
			}
			if !parser.matchDirect(lexer.NewLine) {
				return nil, errors.New(fmt.Sprintf("invalid struct statement at line %d", parser.currentToken.Line))
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			var caseBody []ast.Node
			for ; !parser.complete; {
				if parser.matchDirect(lexer.End) ||
					parser.matchDirect(lexer.Case) {
					break
				} else if parser.matchKind(lexer.Separator) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
					continue
				}
				caseChild, parsingError = parser.parseBinaryExpression(0)
				if parsingError != nil {
					return nil, parsingError
				}
				caseBody = append(caseBody, caseChild)
			}
			caseBlocks = append(caseBlocks, &ast.CaseBlock{
				Cases: cases,
				Body:  caseBody,
			})
		} else if parser.matchDirect(lexer.Else) {
			// Parse Else
			if !parser.matchDirect(lexer.NewLine) {
				return nil, errors.New(fmt.Sprintf("invalid struct statement at line %d", parser.currentToken.Line))
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			for ; !parser.complete; {
				if parser.matchDirect(lexer.End) {
					break
				} else if parser.matchKind(lexer.Separator) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
					continue
				}
				elseChild, parsingError = parser.parseBinaryExpression(0)
				if parsingError != nil {
					return nil, parsingError
				}
				elseBody = append(elseBody, elseChild)
			}
			break
		} else {
			return nil, errors.New(fmt.Sprintf("invalid declaration of switch statement at line %d", parser.currentToken.Line))
		}
	}
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("Switch declaration never ended at line %d", parser.currentToken.Line))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.SwitchStatement{
		Target:     target,
		CaseBlocks: caseBlocks,
		Else:       elseBody,
	}, nil
}

func (parser *Parser) parseModuleStatement() (ast.Statement, error) {
	return nil, nil
}

func (parser *Parser) parseFunctionDefinitionStatement() (ast.Statement, error) {
	return nil, nil
}

func (parser *Parser) parseStructStatement() (ast.Statement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchKind(lexer.IdentifierKind) {
		return nil, errors.New(fmt.Sprintf("invalid struct definition at line %d", parser.currentToken.Line))
	}
	name := &ast.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid struct definition at line %d", parser.currentToken.Line))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var fields []*ast.Identifier
	for ; !parser.complete; {
		if parser.matchDirect(lexer.End) {
			break
		} else if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			continue
		} else if !parser.matchKind(lexer.IdentifierKind) {
			return nil, errors.New(fmt.Sprintf("invalid struct definition at line %d", parser.currentToken.Line))
		}
		fields = append(fields, &ast.Identifier{
			Token: parser.currentToken,
		})
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("invalid struct definition at line %d", parser.currentToken.Line))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.StructStatement{
		Name:   name,
		Fields: fields,
	}, nil
}

func (parser *Parser) parseInterfaceStatement() (ast.Statement, error) {
	return nil, nil
}

func (parser *Parser) parseDeferStatement() (ast.Statement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	methodInvocation, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	switch methodInvocation.(type) {
	case *ast.MethodInvocationExpression:
		return &ast.DeferStatement{
			X: methodInvocation.(*ast.MethodInvocationExpression),
		}, nil
	default:
		return nil, errors.New(fmt.Sprintf("no function call passed to go statement at line %d", parser.currentToken.Line))
	}
}

func (parser *Parser) parseClassStatement() (ast.Statement, error) {
	return nil, nil
}

func (parser *Parser) parseTryStatement() (ast.Statement, error) {
	return nil, nil
}

func (parser *Parser) parseBeginStatement() (ast.Statement, error) {
	return nil, nil
}

func (parser *Parser) parseEndStatement() (ast.Statement, error) {
	return nil, nil
}

func (parser *Parser) parseGoStatement() (ast.Statement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	methodInvocation, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	switch methodInvocation.(type) {
	case *ast.MethodInvocationExpression:
		return &ast.GoStatement{
			X: methodInvocation.(*ast.MethodInvocationExpression),
		}, nil
	default:
		return nil, errors.New(fmt.Sprintf("no function call passed to go statement at line %d", parser.currentToken.Line))
	}
}

func (parser *Parser) parseReturnStatement() (ast.Statement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var results []ast.Expression
	for ; !parser.complete; {
		result, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		results = append(results, result)
		if !parser.matchDirect(lexer.Comma) {
			break
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	return &ast.ReturnStatement{
		Results: results,
	}, nil
}

func (parser *Parser) parseYieldStatement() (ast.Statement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var results []ast.Expression
	for ; !parser.complete; {
		result, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		results = append(results, result)
		if !parser.matchDirect(lexer.Comma) {
			break
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	return &ast.YieldStatement{
		Results: results,
	}, nil
}

func (parser *Parser) parseSuperStatement() (ast.Statement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirect(lexer.OpenParentheses) {
		return nil, errors.New(fmt.Sprintf("invalid super call at line %d", parser.currentToken.Line))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var arguments []ast.Expression
	for ; !parser.complete; {
		if parser.matchDirect(lexer.CloseParentheses) {
			break
		}
		argument, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		arguments = append(arguments, argument)
		if !parser.matchDirect(lexer.Comma) {
			break
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.SuperInvocationStatement{
		Arguments: arguments,
	}, nil
}

func (parser *Parser) parseRetryStatement() (ast.Statement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.RetryStatement{}, nil
}

func (parser *Parser) parseBreakStatement() (ast.Statement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.BreakStatement{}, nil
}

func (parser *Parser) parseRedoStatement() (ast.Statement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.RedoStatement{}, nil
}

func (parser *Parser) parsePassStatement() (ast.Statement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.PassStatement{}, nil
}

func (parser *Parser) parseEnumStatement() (ast.Statement, error) { // What about initializing it's identifiers with an specific value?
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchKind(lexer.IdentifierKind) {
		return nil, errors.New(fmt.Sprintf("invalid declaration of enum statement at line %d", parser.currentToken.Line))
	}
	namespace := &ast.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid declaration of enum statement at line %d", parser.currentToken.Line))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var identifiers []*ast.Identifier
	for ; !parser.complete; {
		if parser.matchDirect(lexer.End) {
			break
		} else if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			continue
		} else if !parser.matchKind(lexer.IdentifierKind) {
			return nil, errors.New(fmt.Sprintf("invalid declaration of enum statement at line %d", parser.currentToken.Line))
		}
		identifiers = append(identifiers, &ast.Identifier{
			Token: parser.currentToken,
		})
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("enum never ended at line %d", parser.currentToken.Line))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.EnumStatement{
		Name:            namespace,
		EnumIdentifiers: identifiers,
	}, nil
}

func (parser *Parser) parseOperand() (ast.Node, error) {
	switch parser.currentToken.Kind {
	case lexer.Literal, lexer.Boolean, lexer.NoneType:
		return parser.parseLiteral()
	case lexer.IdentifierKind:
		identifier := parser.currentToken
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		return &ast.Identifier{
			Token: identifier,
		}, nil
	case lexer.Keyboard:
		switch parser.currentToken.DirectValue {
		case lexer.Lambda:
			return parser.parseLambdaExpression()
		case lexer.While:
			return parser.parseWhileLoop()
		case lexer.For:
			return parser.parseForLoop()
		case lexer.Until:
			return parser.parseUntilLoop()
		case lexer.If:
			return parser.parseIfStatement()
		case lexer.Unless:
			return parser.parseUnlessStatement()
		case lexer.Switch:
			return parser.parseSwitchStatement()
		case lexer.Module:
			return parser.parseModuleStatement()
		case lexer.Def:
			return parser.parseFunctionDefinitionStatement()
		case lexer.Struct:
			return parser.parseStructStatement()
		case lexer.Interface:
			return parser.parseInterfaceStatement()
		case lexer.Defer:
			return parser.parseDeferStatement()
		case lexer.Class:
			return parser.parseClassStatement()
		case lexer.Try:
			return parser.parseTryStatement()
		case lexer.BEGIN:
			return parser.parseBeginStatement()
		case lexer.END:
			return parser.parseEndStatement()
		case lexer.Go:
			return parser.parseGoStatement()
		case lexer.Return:
			return parser.parseReturnStatement()
		case lexer.Yield:
			return parser.parseYieldStatement()
		case lexer.Super:
			return parser.parseSuperStatement()
		case lexer.Retry:
			return parser.parseRetryStatement()
		case lexer.Break:
			return parser.parseBreakStatement()
		case lexer.Redo:
			return parser.parseRedoStatement()
		case lexer.Pass:
			return parser.parsePassStatement()
		case lexer.Enum:
			return parser.parseEnumStatement()
		}
	case lexer.Punctuation:
		switch parser.currentToken.DirectValue {
		case lexer.OpenParentheses:
			return parser.parseParentheses()
		case lexer.OpenSquareBracket: // Parse Arrays
			return parser.parseArrayExpression()
		case lexer.OpenBrace: // Parse Dictionaries
			return parser.parseHashExpression()
		}
	}
	return nil, errors.New(fmt.Sprintf("unknown expression with token at line %d", parser.currentToken.Line))
}

func (parser *Parser) parseSelectorExpression(expression ast.Node) (ast.Node, error) {
	selector := expression
	for ; !parser.complete; {
		if !parser.matchDirect(lexer.Dot) {
			break
		}
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		identifier := parser.currentToken
		if identifier.Kind != lexer.IdentifierKind {
			return nil, errors.New(fmt.Sprintf("invalid selector at token in line %d", identifier.Line))
		}
		selector = &ast.SelectorExpression{
			X: selector,
			Identifier: &ast.Identifier{
				Token: identifier,
			},
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	return selector, nil
}

func (parser *Parser) parseMethodInvocationExpression(expression ast.Node) (ast.Node, error) {
	var arguments []ast.Expression
	// The first token is open parentheses
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	for ; !parser.complete; {
		if parser.matchDirect(lexer.CloseParentheses) {
			break
		}
		argument, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		arguments = append(arguments, argument)
		if parser.matchDirect(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		}
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.MethodInvocationExpression{
		Function:  expression,
		Arguments: arguments,
	}, nil
}

func (parser *Parser) parseIndexExpression(expression ast.Node) (ast.Node, error) {
	tokenizationError := parser.next()
	if tokenizationError != nil {
		return nil, tokenizationError
	}
	var rightIndex ast.Expression
	leftIndex, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if parser.matchDirect(lexer.Colon) {
		tokenizationError = parser.next()
		if tokenizationError != nil {
			return nil, tokenizationError
		}
		rightIndex, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
	}
	tokenizationError = parser.next()
	if tokenizationError != nil {
		return nil, tokenizationError
	}
	return &ast.IndexExpression{
		Source: expression,
		Index: [2]ast.Expression{
			leftIndex,
			rightIndex,
		},
	}, nil
}

func (parser *Parser) parseIfOneLiner(result ast.Expression) (ast.Node, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	condition, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if !parser.matchDirect(lexer.Else) {
		return &ast.OneLineIfExpression{
			Result:     result,
			Condition:  condition,
			ElseResult: nil,
		}, nil
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var elseResult ast.Expression
	elseResult, parsingError = parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	return &ast.OneLineIfExpression{
		Result:     result,
		Condition:  condition,
		ElseResult: elseResult,
	}, nil
}

func (parser *Parser) parseUnlessOneLiner(result ast.Expression) (ast.Node, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	condition, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if !parser.matchDirect(lexer.Else) {
		return &ast.OneLineUnlessExpression{
			Result:     result,
			Condition:  condition,
			ElseResult: nil,
		}, nil
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var elseResult ast.Expression
	elseResult, parsingError = parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	return &ast.OneLineUnlessExpression{
		Result:     result,
		Condition:  condition,
		ElseResult: elseResult,
	}, nil
}

func (parser *Parser) parseGeneratorExpression(operation ast.Expression) (ast.Node, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var variables []*ast.Identifier
	numberOfVariables := 0
	for ; !parser.complete; {
		if parser.matchDirect(lexer.In) {
			break
		}
		variables = append(variables, &ast.Identifier{
			Token: parser.currentToken,
		})
		numberOfVariables++
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	if numberOfVariables == 0 {
		return nil, errors.New(fmt.Sprintf("syntax error: no receivers in generator defined at line %d", parser.currentToken.Line))
	}
	if !parser.matchDirect(lexer.In) {
		return nil, errors.New(fmt.Sprintf("syntax error: invalid generator syntax at line %d", parser.currentToken.Line))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	source, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	return &ast.GeneratorExpression{
		Operation: operation,
		Variables: variables,
		Source:    source,
	}, nil
}

func (parser *Parser) parseAssignmentStatement(leftHandSide ast.Expression) (ast.Node, error) {
	assignmentToken := parser.currentToken
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	rightHandSide, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	return &ast.AssignStatement{
		LeftHandSide:   leftHandSide,
		AssignOperator: assignmentToken,
		RightHandSide:  rightHandSide,
	}, nil
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
			expression, parsingError = parser.parseMethodInvocationExpression(expression)
		case lexer.OpenSquareBracket: // Is indexing
			expression, parsingError = parser.parseIndexExpression(expression)
		case lexer.For: // Generators
			expression, parsingError = parser.parseGeneratorExpression(expression)
		case lexer.If: // One line If
			expression, parsingError = parser.parseIfOneLiner(expression)
		case lexer.Unless: // One line Unless
			expression, parsingError = parser.parseUnlessOneLiner(expression)
		default:
			if parser.matchKind(lexer.Assignment) {
				expression, parsingError = parser.parseAssignmentStatement(expression)
			}
			break expressionPendingLoop
		}
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
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	for ; !parser.complete; {
		if parser.matchKind(lexer.Separator) {
			for ; !parser.complete; {
				if !parser.matchKind(lexer.Separator) {
					break
				}
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
			}
			if parser.matchKind(lexer.EOF) {
				break
			}
		}

		parsedExpression, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		result.Body = append(result.Body, parsedExpression)
	}
	parser.complete = true
	return result, nil
}

func NewParser(lexer_ *lexer.Lexer) (*Parser, error) {
	parser := &Parser{
		lexer:        lexer_,
		complete:     false,
		currentToken: nil,
	}
	return parser, nil
}
