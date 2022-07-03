package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/common"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func simplifyBasicLiteral(literal *ast.BasicLiteralExpression) *ast.BasicLiteralExpression {
	var (
		result = &ast.BasicLiteralExpression{}
		token  *lexer.Token
	)
	switch {
	case literalIsString(literal):
		token = common.ResolveString(literal.Token)
	case literalIsBytes(literal):
		token = common.ResolveBytesString(literal.Token)
	case literalIsInteger(literal):
		token = common.ResolveInteger(literal.Token)
	case literalIsFloat(literal):
		token = common.ResolveFloat(literal.Token)
	}
	result.Kind = token.Kind
	result.DirectValue = token.DirectValue
	return result
}
