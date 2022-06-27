package lexer

import (
	"github.com/shoriwe/gplasma/pkg/errors"
	"github.com/shoriwe/gplasma/pkg/reader"
	"regexp"
)

var (
	identifierCheck = regexp.MustCompile("(?m)^[a-zA-Z_]+[a-zA-Z0-9_]*$")
	junkKindCheck   = regexp.MustCompile("(?m)^\\00+$")
)

type Lexer struct {
	lastToken *Token
	line      int
	reader    reader.Reader
	complete  bool
}

func (lexer *Lexer) HasNext() bool {
	return !lexer.complete
}

func (lexer *Lexer) next() (*Token, *errors.Error) {
	if !lexer.reader.HasNext() {
		lexer.complete = true
		return &Token{
			String:      "EOF",
			DirectValue: InvalidDirectValue,
			Kind:        EOF,
			Line:        lexer.line,
			Index:       lexer.reader.Index(),
		}, nil
	}
	var tokenizingError *errors.Error
	var kind Kind
	var content []rune
	directValue := InvalidDirectValue
	line := lexer.line
	index := lexer.reader.Index()
	char := lexer.reader.Char()
	lexer.reader.Next()
	switch char {
	case '\r':
		if lexer.reader.Char() == NewLineChar {
			lexer.reader.Next()
			lexer.line++
			content = []rune{char}
			kind = Separator
			directValue = NewLine
		}
	case NewLineChar:
		lexer.line++
		content = []rune{char}
		kind = Separator
		directValue = NewLine
	case SemiColonChar:
		content = []rune{char}
		kind = Separator
		directValue = SemiColon
	case ColonChar:
		content = []rune{char}
		directValue = Colon
		kind = Punctuation
	case CommaChar:
		content = []rune{char}
		directValue = Comma
		kind = Punctuation
	case OpenParenthesesChar:
		content = []rune{char}
		directValue = OpenParentheses
		kind = Punctuation
	case CloseParenthesesChar:
		content = []rune{char}
		directValue = CloseParentheses
		kind = Punctuation
	case OpenSquareBracketChar:
		content = []rune{char}
		directValue = OpenSquareBracket
		kind = Punctuation
	case CloseSquareBracketChar:
		content = []rune{char}
		directValue = CloseSquareBracket
		kind = Punctuation
	case OpenBraceChar:
		content = []rune{char}
		directValue = OpenBrace
		kind = Punctuation
	case CloseBraceChar:
		content = []rune{char}
		directValue = CloseBrace
		kind = Punctuation
	case DollarSignChar:
		content = []rune{char}
		directValue = DollarSign
		kind = Punctuation
	case DotChar:
		content = []rune{char}
		directValue = Dot
		kind = Punctuation
	case WhiteSpaceChar:
		directValue = Blank
		kind = Whitespace
		content = []rune{char}
	case TabChar:
		directValue = Blank
		kind = Whitespace
		content = []rune{char}
	case CommentChar:
		content, kind, tokenizingError = lexer.tokenizeComment()
		content = append([]rune{'#'}, content...)
	case '\'', '"': // String1
		content, kind, directValue, tokenizingError = lexer.tokenizeStringLikeExpressions(char)
	case '`':
		content, kind, directValue, tokenizingError = lexer.tokenizeStringLikeExpressions(char)
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		content, kind, directValue, tokenizingError = lexer.tokenizeNumeric(char)
	case StarChar:
		content, kind, directValue = lexer.tokenizeRepeatableOperator(char, Star, Operator, PowerOf, Operator, StarAssign, Assignment, PowerOfAssign, Assignment)
	case DivChar:
		content, kind, directValue = lexer.tokenizeRepeatableOperator(char, Div, Operator, FloorDiv, Operator, DivAssign, Assignment, FloorDivAssign, Assignment)
	case LessThanChar:
		content, kind, directValue = lexer.tokenizeRepeatableOperator(char, LessThan, Comparator, BitwiseLeft, Operator, LessOrEqualThan, Comparator, BitwiseLeftAssign, Assignment)
	case GreatThanChar:
		content, kind, directValue = lexer.tokenizeRepeatableOperator(char, GreaterThan, Comparator, BitwiseRight, Operator, GreaterOrEqualThan, Comparator, BitwiseRightAssign, Assignment)
	case AddChar:
		content, kind, directValue = lexer.tokenizeSingleOperator(char, Add, Operator, AddAssign, Assignment)
	case SubChar:
		content, kind, directValue = lexer.tokenizeSingleOperator(char, Sub, Operator, SubAssign, Assignment)
	case ModulusChar:
		content, kind, directValue = lexer.tokenizeSingleOperator(char, Modulus, Operator, ModulusAssign, Assignment)
	case BitwiseXorChar:
		content, kind, directValue = lexer.tokenizeSingleOperator(char, BitwiseXor, Operator, BitwiseXorAssign, Assignment)
	case BitWiseAndChar:
		content, kind, directValue = lexer.tokenizeSingleOperator(char, BitWiseAnd, Operator, BitWiseAndAssign, Assignment)
	case BitwiseOrChar:
		content, kind, directValue = lexer.tokenizeSingleOperator(char, BitwiseOr, Operator, BitwiseOrAssign, Assignment)
	case SignNotChar:
		content, kind, directValue = lexer.tokenizeSingleOperator(char, SignNot, Operator, NotEqual, Comparator)
	case NegateBitsChar:
		content, kind, directValue = lexer.tokenizeSingleOperator(char, NegateBits, Operator, NegateBitsAssign, Assignment)
	case EqualsChar:
		content, kind, directValue = lexer.tokenizeSingleOperator(char, Assign, Assignment, Equals, Comparator)
	case BackSlashChar:
		content = []rune{char}
		if lexer.reader.HasNext() {
			nextChar := lexer.reader.Char()
			if nextChar != '\n' {
				return nil, errors.New(lexer.line, "line escape not followed by a new line", errors.LexingError)
			}
			content = append(content, '\n')
			lexer.reader.Next()
		}
		kind = PendingEscape
	default:
		if char == 'b' {
			if lexer.reader.HasNext() {
				nextChar := lexer.reader.Char()
				if nextChar == '\'' || nextChar == '"' {
					var byteStringPart []rune
					lexer.reader.Next()
					byteStringPart, kind, directValue, tokenizingError = lexer.tokenizeStringLikeExpressions(nextChar)
					content = append([]rune{char}, byteStringPart...)
					if directValue != InvalidDirectValue {
						directValue = ByteString
					}
					break
				}
			}
		}
		content, kind, directValue, tokenizingError = lexer.tokenizeWord(char)
	}
	return &Token{
		DirectValue: directValue,
		String:      string(content),
		Kind:        kind,
		Line:        line,
		Index:       index,
	}, tokenizingError
}

/*
	This function will yield just the necessary token, this means not repeated separators
*/
func (lexer *Lexer) Next() (*Token, *errors.Error) {
	token, tokenizingError := lexer.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if token.Kind == JunkKind {
		return lexer.Next()
	}
	if token.Kind == Comment {
		return lexer.Next()
	}
	if token.Kind == Whitespace {
		return lexer.Next()
	}
	if token.Kind == Separator {
		if lexer.lastToken == nil {
			return lexer.Next()
		}
		switch lexer.lastToken.Kind {
		case Separator:
			return lexer.Next()
		case Operator, Comparator:
			return lexer.Next()
		default:
			break
		}
		switch lexer.lastToken.DirectValue {
		case Comma, OpenParentheses, OpenSquareBracket, OpenBrace:
			return lexer.Next()
		default:
			break
		}
	}
	lexer.lastToken = token
	return token, nil
}

func NewLexer(codeReader reader.Reader) *Lexer {
	return &Lexer{
		lastToken: nil,
		line:      1,
		reader:    codeReader,
		complete:  false,
	}
}
