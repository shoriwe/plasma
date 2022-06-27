package lexer

import "github.com/shoriwe/gplasma/pkg/errors"

func (lexer *Lexer) tokenizeNumeric(firstDigit rune) ([]rune, Kind, DirectValue, *errors.Error) {
	if !lexer.reader.HasNext() {
		return []rune{firstDigit}, Literal, Integer, nil
	}
	nextChar := lexer.reader.Char()
	lexer.reader.Next()
	if firstDigit == '0' {
		switch nextChar {
		case 'x', 'X': // Hexadecimal
			return lexer.tokenizeHexadecimal(nextChar)
		case 'b', 'B': // Binary
			return lexer.tokenizeBinary(nextChar)
		case 'o', 'O': // Octal
			return lexer.tokenizeOctal(nextChar)
		case 'e', 'E': // Scientific float
			return lexer.tokenizeScientificFloat([]rune{firstDigit, nextChar})
		case '.': // Maybe a float
			return lexer.tokenizeFloat([]rune{firstDigit, nextChar}) // Integer, Float Or Scientific Float
		default:
			if ('0' <= nextChar && nextChar <= '9') || nextChar == '_' {
				return lexer.tokenizeInteger([]rune{firstDigit, nextChar}) // Integer, Float or Scientific Float
			}
		}
	} else {
		switch nextChar {
		case 'e', 'E': // Scientific float
			return lexer.tokenizeScientificFloat([]rune{firstDigit, nextChar})
		case '.': // Maybe a float
			return lexer.tokenizeFloat([]rune{firstDigit, nextChar}) // Integer, Float Or Scientific Float
		default:
			if ('0' <= nextChar && nextChar <= '9') || nextChar == '_' {
				return lexer.tokenizeInteger([]rune{firstDigit, nextChar}) // Integer, Float or Scientific Float
			}
		}
	}
	lexer.reader.Redo()
	return []rune{firstDigit}, Literal, Integer, nil
}
