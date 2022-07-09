package simplification

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

var current uint64 = 0

func nextAnonIdentifier() *ast2.Identifier {
	return &ast2.Identifier{
		Symbol: fmt.Sprintf("%d", current),
	}
}
