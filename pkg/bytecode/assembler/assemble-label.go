package assembler

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (a *assembler) Label(label *ast3.Label) []byte {
	oldIndex, found := a.labels[label.Code]
	if found && oldIndex != -1 {
		// fmt.Printf("Repeated label %d - Old: %d and New: %d\n", label.Code, oldIndex, a.bytecodeIndex)
		panic(fmt.Sprintf("Label %d should be unique", label.Code))
	}
	a.labels[label.Code] = a.bytecodeIndex
	return nil
}
