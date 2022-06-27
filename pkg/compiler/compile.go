package compiler

import (
	"github.com/shoriwe/gplasma/pkg/compiler/lexer"
	"github.com/shoriwe/gplasma/pkg/compiler/parser"
	"github.com/shoriwe/gplasma/pkg/errors"
	"github.com/shoriwe/gplasma/pkg/reader"
	"github.com/shoriwe/gplasma/pkg/vm"
)

func Compile(r reader.Reader) (*vm.Bytecode, *errors.Error) {
	l := lexer.NewLexer(r)
	p := parser.NewParser(l)
	program, parsingError := p.Parse()
	if parsingError != nil {
		return nil, errors.New(0, parsingError.Error(), "SI")
	}
	sourceCode, compilationError := program.Compile()
	if compilationError != nil {
		return nil, compilationError
	}
	return vm.NewBytecodeFromArray(sourceCode), nil
}
