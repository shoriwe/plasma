package assembler

import (
	"github.com/shoriwe/plasma/pkg/ast3"
	"github.com/shoriwe/plasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/plasma/pkg/common"
)

func (a *assembler) Identifier(ident *ast3.Identifier) []byte {
	var result []byte
	result = append(result, opcodes.Identifier)
	result = append(result, common.IntToBytes(len(ident.Symbol))...)
	result = append(result, []byte(ident.Symbol)...)
	return result
}
