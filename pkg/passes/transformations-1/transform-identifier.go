package transformations_1

import (
	"github.com/shoriwe/plasma/pkg/ast2"
	"github.com/shoriwe/plasma/pkg/ast3"
)

func (transform *transformPass) Identifier(ident *ast2.Identifier) *ast3.Identifier {
	return &ast3.Identifier{
		Symbol: ident.Symbol,
	}
}
