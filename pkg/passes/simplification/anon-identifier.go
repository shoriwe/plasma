package simplification

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simp *simplify) nextAnonIdentifier() *ast2.Identifier {
	ident := simp.currentAnonIdent
	simp.currentAnonIdent++
	return &ast2.Identifier{
		Symbol: fmt.Sprintf("____%d", ident),
	}
}
