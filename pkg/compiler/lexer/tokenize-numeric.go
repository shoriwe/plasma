package lexer

import "github.com/shoriwe/gplasma/pkg/errors"

func (lexer *Lexer) tokenizeNumeric() *errors.Error {
	if !lexer.reader.HasNext() {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = Integer
		return nil
	}
	nextChar := lexer.reader.Char()
	lexer.reader.Next()
	if lexer.currentToken.Contents[0] == '0' {
		switch nextChar {
		case 'x', 'X': // Hexadecimal
			lexer.currentToken.append(nextChar)
			return lexer.tokenizeHexadecimal()
		case 'b', 'B': // Binary
			lexer.currentToken.append(nextChar)
			return lexer.tokenizeBinary()
		case 'o', 'O': // Octal
			lexer.currentToken.append(nextChar)
			return lexer.tokenizeOctal()
		case 'e', 'E': // Scientific float
			lexer.currentToken.append(nextChar)
			return lexer.tokenizeScientificFloat()
		case '.': // Maybe a float
			lexer.currentToken.append(nextChar)
			return lexer.tokenizeFloat() // Integer, Float Or Scientific Float
		default:
			if ('0' <= nextChar && nextChar <= '9') || nextChar == '_' {
				lexer.currentToken.append(nextChar)
				return lexer.tokenizeInteger() // Integer, Float or Scientific Float
			}
		}
	} else {
		switch nextChar {
		case 'e', 'E': // Scientific float
			lexer.currentToken.append(nextChar)
			return lexer.tokenizeScientificFloat()
		case '.': // Maybe a float
			lexer.currentToken.append(nextChar)
			return lexer.tokenizeFloat() // Integer, Float Or Scientific Float
		default:
			if ('0' <= nextChar && nextChar <= '9') || nextChar == '_' {
				lexer.currentToken.append(nextChar)
				return lexer.tokenizeInteger() // Integer, Float or Scientific Float
			}
		}
	}
	lexer.reader.Redo()
	lexer.currentToken.Kind = Literal
	lexer.currentToken.DirectValue = Integer
	return nil
}
