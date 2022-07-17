package gplasma

import (
	"github.com/shoriwe/gplasma/pkg/compiler"
	"github.com/shoriwe/gplasma/pkg/vm"
	"io"
)

func NewVM(stdin io.Reader, stdout io.Writer, stderr io.Writer) *vm.Plasma {
	return vm.NewVM(stdin, stdout, stderr)
}

func Compile(scriptCode string) ([]byte, error) {
	return compiler.Compile(scriptCode)
}
